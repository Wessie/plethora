package scheduler

// Task is a task to be ran by the scheduler, it requires
// a Job to execute, and a Planner to plan said execution.
type Task struct {
	Job
	Planner
}

func (t *Task) Cancel() error {
	return nil
}

// ScheduleTask schedules a task to be run
func (s Scheduler) ScheduleTask(t Task) {
	s.newTask <- t
}
