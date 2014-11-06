package scheduler

// Task is a task to be ran by the scheduler, it requires
// a Job to execute, and a Planner to plan said execution.
type Task struct {
	Job
	Planner

	// schedule is the Scheduler this task is to be
	// run by
	schedule *Scheduler
}

// Cancel stops the task from being executed by the scheduler. If the
// task was not scheduled this will do nothing.
func (t Task) Cancel() error {
	if t.schedule == nil {
		return nil
	}

	t.schedule.removeTask <- t
	return nil
}

// ScheduleTask schedules a task to be run
func (s Scheduler) ScheduleTask(t Task) Task {
	s.newTask <- t
	t.schedule = &s
	return t
}
