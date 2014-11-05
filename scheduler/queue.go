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

	if sl.head.key.After(key) {
		sl.head.prev = el
		el.next = sl.head
		sl.head = el
		return sl.first()
	}

	if sl.tail.key.Before(key) {
		sl.tail.next = el
		el.prev = sl.tail
		sl.tail = el
		return sl.first()
	}

	// find the element to insert after
	var after *element
	for after = sl.head; after.key.Before(key); after = after.next {
	}

	after.next.prev = el
	el.next = after.next
	after.next = el
	el.prev = after

	return sl.first()
}

func (sl *sortedQueue) pop(key time.Time) (time.Time, Task) {
	if sl.head == nil {
		return NoMore, Task{}
	}

	if sl.head.key.After(key) {
		return NoMore, Task{}
	}

	if sl.head.next != nil {
		sl.head.next.prev = nil
	} else {
		sl.tail = nil
	}

	n := sl.head
	sl.head = n.next
	return n.key, n.val
}

func (sl *sortedQueue) remove(val Task) time.Time {
	for c := sl.head; c != nil; c = c.next {
		if c.val != val {
			continue
		}

		if c.prev == nil {
			sl.head = c.next
		} else {
			c.prev.next = c.next
		}

		if c.next == nil {
			sl.tail = c.prev
		} else {
			c.next.prev = c.prev
		}
	}

	return sl.first()
}
