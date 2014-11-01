package scheduler

// Task consists of a Job and a Schedule to run the job on, the job
// will be rescheduled until Schedule.Next returns NoMore.
type Task struct {
	Job
	Schedule
}

// ScheduleTask schedules a task to be run
func (s Scheduler) ScheduleTask(t Task) {
	s.newTask <- t
}
