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

func (sl *sortedQueue) first() time.Time {
	if sl.head == nil {
		return NoMore
	}

	return sl.head.key
}

func (sl *sortedQueue) put(key time.Time, val Task) time.Time {
	el := &element{
		key: key,
		val: val,
	}

	if sl.head == nil {
		sl.head = el
		sl.tail = el
		return sl.first()
	}

	if key.Before(sl.head.key) {
		sl.head.prev = el
		el.next = sl.head
		sl.head = el
		return sl.first()
	}

	if key.After(sl.tail.key) {
		sl.tail.next = el
		el.prev = sl.tail
		sl.tail = el
		return sl.first()
	}

	// find the element to insert after, we can omit nil checks because
	// we know the inserted element will be between the head and tail
	var n *element
	for n = sl.head; n.key.Before(key); n = n.next {
	}

	n.next.prev = el
	el.next = n.next
	n.next = el
	el.prev = n

	return sl.first()
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

// removeTask removes all elements that have the Task given
// as value.
func (sl *sortedQueue) removeTask(val Task) time.Time {
	for c := sl.head; c != nil; c = c.next {
		if c.val == val {
			sl.remove(c)
		}
	}

	return sl.first()
}

// remove removes the node from the queue
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
