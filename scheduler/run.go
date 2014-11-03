package scheduler

func runTask(s Scheduler, t Task) {
	t.Run()
	s.ScheduleTask(t)
	return
}
