package scheduler

import "time"

var MaximumSleep = time.Minute * 30

// minDuration returns the lowest of (a, b)
func minDuration(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}

// New returns an initialized Scheduler, the scheduler
// is automatically started before NewScheduler returns.
//
// The returned scheduler has to be stopped by calling Stop or
// it will leak resources.
func New(name string) *Scheduler {
	s := &Scheduler{
		Name:       name,
		Stats:      NewStats(),
		stop:       make(chan struct{}),
		stopped:    make(chan struct{}),
		newTask:    make(chan Task),
		removeTask: make(chan Task),
		queue:      new(sortedQueue),
	}

	go s.manage()
	return s
}

// Scheduler implements a basic task scheduler
type Scheduler struct {
	// Name of the scheduler instance, this is used for storing
	// the schedule and any other persistent information.
	Name string

	// Stats is a collection of various statistics for this scheduler
	Stats *Stats

	// channel to signal a "you should stop running"
	stop chan struct{}
	// channel to signal "we have stopped running"
	stopped chan struct{}
	// channel to signal a new task is to be scheduled
	newTask chan Task
	// channel to signal removal of a task (before it has ran)
	removeTask chan Task

	// queue of tasks to be run
	queue *sortedQueue
}

// ScheduleTask schedules a task to be run
func (s *Scheduler) ScheduleTask(tsk Task) Task {
	s.newTask <- tsk
	tsk.schedule = s
	return tsk
}

// Schedule schedules j, the Job given according to the Planner p
func (s *Scheduler) Schedule(p Planner, j Job) Task {
	return s.ScheduleTask(Task{
		Job:     j,
		Planner: p,
	})
}

// ScheduleJob schedules a Job to be run once at the specified time
func (s *Scheduler) ScheduleJob(t time.Time, j Job) Task {
	return s.ScheduleTask(Task{
		Job:     j,
		Planner: newTimePlanner(t),
	})
}

// manage manages the scheduling process
func (s *Scheduler) manage() {
	var wait = time.NewTimer(MaximumSleep)
	var nextTaskTime time.Time
	var estimate time.Duration

stopScheduler:
	for {
		select {
		case t := <-s.newTask:
			nextTaskTime = s.queueTask(t)
		case t := <-s.removeTask:
			nextTaskTime = s.unqueueTask(t)
		case <-wait.C:
			nextTaskTime = s.processQueue()
		case <-s.stop:
			break stopScheduler
		}

		// sync our timer to the schedule changes, we don't ever sleep
		// longer than MaximumSleep.
		estimate = nextTaskTime.Sub(time.Now())
		wait.Reset(minDuration(estimate, MaximumSleep))
	}

	close(s.stopped)
	wait.Stop()
}

// queueTask adds a task to the scheduling queue, this function is
// only safe to call from the managing goroutine.
func (s *Scheduler) queueTask(tsk Task) time.Time {
	var taskTime = tsk.PlanJob(tsk.Job)

	if !taskTime.Equal(NoPlan) {
		s.queue.put(taskTime, tsk)
	}

	return s.queue.first()
}

// unqueueTask removes a task from the scheduling queue, this function
// is only safe to call from the managing goroutine.
func (s *Scheduler) unqueueTask(tsk Task) time.Time {
	s.queue.removeTask(tsk)
	return s.queue.first()
}

// processSchedule processes the schedule, this involves a few steps:
// 1. runs all tasks that are ready to run
// 2. returns the time of the first task that isn't ready to run
// yet.
func (s *Scheduler) processQueue() time.Time {
	var task Task
	var taskTime time.Time

	now := time.Now()
	for taskTime, task = s.queue.pop(now); !taskTime.Equal(noTask); taskTime, task = s.queue.pop(now) {
		go s.runTask(task)
	}

	// decide when we want to be called again, if there is stuff in the queue
	// we choose the time of the first task that wants to be ran. Otherwise we
	// use a very liberal 24 hours from the time this function started running.
	if taskTime = s.queue.first(); !taskTime.Equal(noTask) {
		return taskTime
	}

	return now.Add(time.Hour * 24)
}

func (s *Scheduler) runTask(task Task) {
	task.Run()
	s.ScheduleTask(task)
}

// Stop stops the scheduler, Stop waits until an acknowledgement of stopping
// has been received. Calling Stop multiple times does nothing
func (s *Scheduler) Stop() {
	select {
	case <-s.stopped: // already closed check
		return
	case s.stop <- struct{}{}: // send stop signal
	}

	select { // check that we've actually stopped
	case <-s.stopped:
	}
}
