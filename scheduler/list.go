package scheduler

import "time"

type element struct {
	next *element
	prev *element

	key time.Time
	val Task
}

type scheduleList struct {
	head *element
	tail *element
}

func (sl *scheduleList) first() time.Time {
	if sl.head == nil {
		return NoMore
	}

	return sl.head.key
}

func (sl *scheduleList) put(key time.Time, val Task) time.Time {
	el := &element{
		key: key,
		val: val,
	}

	// special case if this is the first item
	if sl.head == nil {
		sl.head = el
		sl.tail = el
		return key
	}

	var c *element
	for c = sl.head; c != nil && c.key.Before(key); c = c.next {
	}

	if c == nil { // insert at tail
		c = sl.tail

		el.prev = c
		sl.tail = el
		c.next = el
		return sl.first()
	}

	// otherwise it's a plain insert before node
	el.next = c
	el.prev = c.prev

	if c.prev == nil {
		sl.head = el
	} else {
		el.prev.next = el
	}

	c.prev = el

	return sl.first()
}

func (sl *scheduleList) pop(key time.Time) (time.Time, Task) {
	if sl.head == nil {
		return NoMore, Task{}
	}

	if sl.head.key.After(key) {
		return NoMore, Task{}
	}

	n := sl.head
	sl.head = sl.head.next
	if sl.head != nil {
		sl.head.prev = nil
	}
	return n.key, n.val
}
