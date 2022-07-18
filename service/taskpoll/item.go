package taskpoll

type TaskPollItem interface {
	ID() string
	TaskName() string
	Weight() uint
	Count() int

	Yield(i int, tpc *TaskPollController)
	OnExit(exitCode TaskPollExitCode)
	Init() TaskPollItem
}
