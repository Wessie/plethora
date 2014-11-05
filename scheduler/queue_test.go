package scheduler

import (
	"errors"
	"testing"
	"time"
)

func fillFrom(sl *sortedQueue, from time.Time) []time.Time {
	var filling []time.Time
	var i time.Duration

	for i = 0; i < 10; i++ {
		when := from.Add(time.Second * i)

		filling = append(filling, when)
		sl.put(when, Task{})
	}

	return filling
}

func TestQueuePop(t *testing.T) {
	var sl sortedQueue
	var now = time.Now()

	filling := fillFrom(&sl, now)
	till := now.Add(time.Second * 5)

	var j int
	for tm, _ := sl.pop(till); tm != NoMore; tm, _ = sl.pop(till) {
		if tm != filling[j] {
			t.Errorf("expected equal times, wanted %s got %s", filling[j], tm)
		}
		j++
	}
}

func TestQueuePutFront(t *testing.T) {
	var q sortedQueue
	var now = time.Now()

	q.put(now, Task{})
	q.put(now.Add(-time.Second), Task{})
}

func TestQueuePopNone(t *testing.T) {
	var sl sortedQueue
	var now = time.Now()

	fillFrom(&sl, now)

	tm, _ := sl.pop(now.Add(-time.Second))

	if tm != NoMore {
		t.Errorf("expected nothing to pop, got: %s", tm)
	}
}

func TestQueueFirst(t *testing.T) {
	var sl sortedQueue
	var now = time.Now()

	fillFrom(&sl, now)

	f := sl.first()
	if f != now {
		t.Errorf("expected now, got: %s != %s", now, f)
	}
}

func TestQueueFillEmptyFillCycle(t *testing.T) {
	var sl sortedQueue
	var now = time.Now()
	var j int
	var everything = now.Add(time.Minute)

	for k := 0; k < 20; k++ {
		t.Logf("entering cycle: %d", k)
		filling := fillFrom(&sl, now)

		j = 0
		for tm, _ := sl.pop(everything); tm != NoMore; tm, _ = sl.pop(everything) {
			t.Logf("entering cycle: %d.%d", k, j)
			if tm != filling[j] {
				t.Errorf("expected equal times, wanted %s got %s", filling[j], tm)
			}
			j++
		}
	}
}

func TestQueueRemove(t *testing.T) {
	var sl sortedQueue
	var now = time.Now()

	fillFrom(&sl, now)

	err := errors.New("queue removal")
	dummyJob := FuncJob(func() error {
		return err
	})

	dummyTime := now.Add(time.Minute)
	dummy := Task{Job: dummyJob}

	// add our dummy one at the end
	sl.put(dummyTime, dummy)
	// remove all others
	sl.removeTask(Task{})

	// now try to grab our dummy again, it should be the head
	tm, tsk := sl.pop(dummyTime)
	if tm != dummyTime {
		t.Errorf("expected dummyTime %s, got %s", dummyTime, tm)
	}

	if dummy.Run() != err {
		t.Errorf("expected dummyJob %v, got %v", dummy, tsk)
	}
}
