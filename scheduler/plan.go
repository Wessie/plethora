package scheduler

import "time"

var NoMore time.Time

// Planner is an interface to be implemented by the planner of jobs,
// the planner is responsible for determining the next date and time
// a job is supposed to run.
type Planner interface {
	PlanJob(Job) time.Time
}

func newTimePlanner(t time.Time) Planner {
	st := timePlanner(t)
	return &st
}

type timePlanner time.Time

func (tp *timePlanner) PlanJob(Job) time.Time {
	return time.Time(*tp)
}
