package scheduler

// Job is a runnable job to be run by the scheduler
type Job interface {
	Run() error
}

// FuncJob implements the Job interface for a function that mimicks the
// Job.Run signature.
type FuncJob func() error

func (fn FuncJob) Run() error {
	return fn()
}
