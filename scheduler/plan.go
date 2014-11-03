package scheduler

import "time"

var NoMore time.Time

type Schedule interface {
	Next() time.Time
	JobRan() bool
}

func newScheduleTime(t time.Time) *scheduleTime {
	st := scheduleTime(t)
	return &st
}

type scheduleTime time.Time

func (st *scheduleTime) Next() time.Time {
	return time.Time(*st)
}

func (st *scheduleTime) JobRan() bool {
	*st = scheduleTime(NoMore)
	return false
}
