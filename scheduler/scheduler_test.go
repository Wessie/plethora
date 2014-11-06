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

	s.ScheduleJob(time.Now(), job)

	select {
	case <-ran:
	case <-time.After(time.Second * 5):
		t.Error("failed to execute job in allowed time")
	}
}

func TestSchedulerDelay(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	var delta = time.Second / 4
	var delay = time.Second
	var s = NewScheduler("test")
	var when = make(chan time.Time)

	job := FuncJob(func() error {
		when <- time.Now()
		return nil
	})

	now := time.Now()
	s.ScheduleJob(time.Now().Add(delay), job)

	var rt time.Time
	select {
	case rt = <-when:
	case <-time.After(time.Second * 3):
		t.Fatal("failed to execute job in allowed time")
	}

	val := rt.Sub(now)
	if delay-delta > val || val > delay+delta {
		t.Errorf("expected execution in %s-%s range, actual was %s",
			delay-delta, delay+delta, val)
	}

	t.Logf("expected range: %s-%s, actual: %s", delay-delta, delay+delta, val)
}
