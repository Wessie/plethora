package scheduler

import "time"

type element struct {
	next *element
	prev *element

	key time.Time
	val Task
}

type sortedQueue struct {
	head *element
	tail *element
}

// first returns the earliest time.Time in the queue
func (sl *sortedQueue) first() time.Time {
	if sl.head == nil {
		return NoMore
	}

	return sl.head.key
}

// put adds the key and value given to the queue, it is sorted by key.
func (sl *sortedQueue) put(key time.Time, val Task) {
	el := &element{
		key: key,
		val: val,
	}

	if sl.head == nil {
		sl.head = el
		sl.tail = el
		return
	}

	if key.Before(sl.head.key) {
		sl.head.prev = el
		el.next = sl.head
		sl.head = el
		return
	}

	if key.After(sl.tail.key) {
		sl.tail.next = el
		el.prev = sl.tail
		sl.tail = el
		return
	}

	// find the element to insert before, we can omit nil checks because
	// we know the resulting element will be between head and tail
	var n *element
	for n = sl.head; n.key.Before(key); n = n.next {
	}

	el.prev = n.prev
	el.next = n

	n.prev.next = el
	n.prev = el
}

func (sl *sortedQueue) pop(key time.Time) (time.Time, Task) {
	if sl.head == nil {
		return NoMore, Task{}
	}

	if key.Before(sl.head.key) {
		return NoMore, Task{}
	}

	n := sl.head
	sl.remove(n)
	return n.key, n.val
}

// removeTask removes all elements that match el.val==val
func (sl *sortedQueue) removeTask(val Task) {
	for c := sl.head; c != nil; c = c.next {
		if c.val == val {
			sl.remove(c)
		}
	}
}

func (sl *sortedQueue) remove(n *element) {
	if n.prev == nil {
		sl.head = n.next
	} else {
		n.prev.next = n.next
	}

	if n.next == nil {
		sl.tail = n.prev
	} else {
		n.next.prev = n.prev
	}
}
