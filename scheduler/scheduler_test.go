package scheduler

import (
	"testing"
	"time"
)

func TestSchedulerStopping(t *testing.T) {
	s := NewScheduler("test")

	stopped := make(chan bool)
	go func() {
		s.Stop()
		close(stopped)
	}()

	select {
	case <-stopped:
	case <-time.After(time.Second * 3):
		t.Error("failed to stop scheduler in allowed timeslot")
	}
}

func TestScheduleJobNow(t *testing.T) {
	var ran = make(chan bool)
	var s = NewScheduler("test")

	job := FuncJob(func() error {
		close(ran)
		return nil
	})

	s.ScheduleJobNow(job)

	select {
	case <-ran:
	case <-time.After(time.Second * 5):
		t.Error("failed to execute job in allowed time")
	}
}

func TestSchedulerDelay(t *testing.T) {
	// TODO: very slow, and error prone test, find a way to fix it
	if testing.Short() {
		t.Skip()
	}
	var s = NewScheduler("test")
	var ran = make(chan bool)

	job := FuncJob(func() error {
		close(ran)
		return nil
	})

	// run our job in 3 seconds, we check 2 seconds before
	// and 2 seconds after the promised time to see if it
	// ran too early or too late.
	s.ScheduleJob(time.Now().Add(time.Second*3), job)

	before := time.After(time.Second * 1)
	after := time.After(time.Second * 5)

	for {
		select {
		case <-before:
			before = nil
		case <-ran:
			if before != nil {
				t.Fatal("job ran too early")
			}
			ran = nil
		case <-after:
			if before != nil || ran != nil {
				t.Error("failed to run job in 2-second window")
			}
			return
		}
	}
}
