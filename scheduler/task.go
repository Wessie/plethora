package scheduler

// Task consists of a Job and a Planner, the job will be
// scheduled and ran according to the planners timing
type Task struct {
	Job
	Planner
}

// ScheduleTask schedules a task to be run
func (s Scheduler) ScheduleTask(t Task) {
	s.newTask <- t
}
