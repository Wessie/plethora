package scheduler

import (
	"testing"
	"time"
)

func TestNewTimePlanner(t *testing.T) {
	now := time.Now()
	p := newTimePlanner(now)

	if got := p.PlanJob(nil); !got.Equal(now) {
		t.Errorf("expected %s, got %s", now, got)
	}

	if got := p.PlanJob(nil); !got.Equal(NoPlan) {
		t.Errorf("expected NoMore, got %s", got)
	}
}
