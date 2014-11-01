package scheduler

func runTask(s Scheduler, t Task) {
	t.Run()

	if t.Next() == NoMore {
		return
	}

	s.ScheduleTask(t)
	return
}
