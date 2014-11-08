package scheduler

import (
	"fmt"
	"reflect"
)

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

func (t Task) String() string {
	var j, p fmt.Stringer
	var ok bool

	if j, ok = t.Job.(fmt.Stringer); !ok {
		j = reflect.TypeOf(t.Job)
	}

	if p, ok = t.Planner.(fmt.Stringer); !ok {
		p = reflect.TypeOf(t.Planner)
	}

	return fmt.Sprintf("%s/%s", j, p)
}
