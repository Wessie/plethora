package scheduler

import "time"

// Job is a runnable job to be run by the scheduler
type Job interface {
	Run() error
}

// ScheduleJob schedules a Job to be run once at the specified time
func (s Scheduler) ScheduleJob(t time.Time, j Job) {
	s.newTask <- Task{
		Job:     j,
		Planner: newTimePlanner(t),
	}
}

// ScheduleJobNow schedules a Job to be run at the earliest possible
// moment. This calls ScheduleJob(time.Now(), j)
func (s Scheduler) ScheduleJobNow(j Job) {
	s.ScheduleJob(time.Now(), j)
}

// FuncJob implements the Job interface for a function that mimicks the
// Job.Run signature.
type FuncJob func() error

func (f FuncJob) Run() error {
	return f()
}
