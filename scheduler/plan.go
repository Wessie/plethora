package scheduler

import "time"

var noTask = time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
var NoPlan = time.Date(0, 0, 0, 1, 0, 0, 0, time.UTC)

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
	re := time.Time(*tp)
	*tp = timePlanner(NoPlan)
	return re
}
