package scheduler

import (
	"testing"
	"time"
)

func TestNewTimePlanner(t *testing.T) {
	now := time.Now()
	p := newTimePlanner(now)

	if got := p.PlanJob(nil); got != now {
		t.Errorf("expected %s, got %s", now, got)
	}

	if got := p.PlanJob(nil); got != NoMore {
		t.Errorf("expected NoMore, got %s", got)
	}
}
