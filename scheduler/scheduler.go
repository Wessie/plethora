package scheduler

import (
	"strings"
	"time"
)

const DatabaseName = "plethora-scheduler"

func NewScheduler(name string) (*Scheduler, error) {
	s := &Scheduler{
		Name:    name,
		stop:    make(chan struct{}),
		stopped: make(chan struct{}),
		newTask: make(chan Task),
	}

	go s.manage()
	return s, nil
}

type Scheduler struct {
	// Name of the scheduler instance, this is used for storing
	// the schedule and any other persistent information.
	Name string

	// channel to signal a "you should stop"
	stop chan struct{}
	// channel to signal "we have stopped"
	stopped chan struct{}
	// channel to signal a new task being scheduled
	newTask chan Task

	// sorted linked list of the schedule
	schedule *scheduleList
}

// manage manages the scheduling process
func (s Scheduler) manage() {
	var wait = time.NewTimer(time.Hour * 24)
	var nextTaskTime time.Time

stopScheduler:
	for {
		select {
		case t := <-s.newTask:
			nextTaskTime = s.scheduleTask(t)

		case <-wait.C:
			nextTaskTime = s.processSchedule()

		case <-s.stop:
			break stopScheduler
		}

		// sync our timer to the schedule changes
		wait.Reset(nextTaskTime.Sub(time.Now()))
	}

	close(s.stopped)
	wait.Stop()
}

// scheduleTask adds a task to the schedule, this function is only
// safe to call from the managing goroutine.
func (s Scheduler) scheduleTask(tsk Task) time.Time {
	var taskTime = tsk.PlanJob(tsk.Job)

	if taskTime == NoMore {
		return s.schedule.first()
	}

	return s.schedule.put(taskTime, tsk)
}

// processSchedule processes the schedule, this involves a few steps:
// 1. runs all tasks that are ready to run
// 2. returns the time of the first task that isn't ready to run
// yet.
func (s Scheduler) processSchedule() time.Time {
	return time.Now().Add(time.Hour * 24)
}

// Stop stops the scheduler, Stop waits until an acknowledgement of stopping
// has been received. Calling Stop multiple times does nothing
func (s Scheduler) Stop() {
	// if we're already closed there is nothing to do here
	if _, ok := <-s.stopped; !ok {
		return
	}

	close(s.stop)
	<-s.stopped
}

func bucketName(name ...string) []byte {
	return []byte(strings.Join(name, "-"))
}
