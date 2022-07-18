package taskpoll

import (
	"math/rand"
	"sync"
	"time"

	"github.com/miaokobot/miaospeed/utils"
	"github.com/miaokobot/miaospeed/utils/structs"
)

type TaskPollExitCode uint

const (
	TPExitSuccess TaskPollExitCode = iota
	TPExitError
	TPExitInterrupt
)

type taskPollItemWrapper struct {
	TaskPollItem
	counter int

	exitOnce sync.Once
}

func (tpw *taskPollItemWrapper) OnExit(exitCode TaskPollExitCode) {
	tpw.exitOnce.Do(func() {
		tpw.TaskPollItem.OnExit(exitCode)
	})
}

type TaskPollController struct {
	name        string
	concurrency uint
	interval    time.Duration
	emptyWait   time.Duration

	taskPoll    []*taskPollItemWrapper
	runningTask map[string]int

	current  uint
	pollLock sync.Mutex
}

func (tpc *TaskPollController) Name() string {
	return tpc.name
}

// single thread
func (tpc *TaskPollController) populate() (int, *taskPollItemWrapper) {
	tpc.pollLock.Lock()
	defer tpc.pollLock.Unlock()

	if tpc.current >= tpc.concurrency {
		return 0, nil
	}

	totalWeight := uint(0)
	totalCount := 0
	for _, tp := range tpc.taskPoll {
		totalWeight += tp.Weight()
		totalCount += tp.Count()
	}

	factor := 0
	if totalWeight > 0 {
		factor = rand.Intn(int(totalWeight))
	}

	for _, tp := range tpc.taskPoll {
		factor -= int(tp.Weight())
		if factor <= 0 {
			counter := tp.counter

			tp.counter += 1
			if tp.counter >= tp.Count() {
				tpc.remove_unsafe(tp.ID(), TPExitSuccess)
			}

			tpc.current += 1
			tpc.runningTask[tp.ID()] += 1
			return counter, tp
		}
	}

	// no task left
	time.Sleep(tpc.emptyWait)

	return 0, nil
}

func (tpc *TaskPollController) release(tpw *taskPollItemWrapper) {
	tpc.pollLock.Lock()
	defer tpc.pollLock.Unlock()

	tpc.runningTask[tpw.ID()] -= 1
	inWaitList := structs.MapContains(tpc.taskPoll, func(w *taskPollItemWrapper) string {
		return w.ID()
	}, tpw.ID())

	if !inWaitList && tpc.runningTask[tpw.ID()] == 0 {
		delete(tpc.runningTask, tpw.ID())
		tpw.OnExit(TPExitSuccess)
	}

	if tpc.current > 0 {
		tpc.current -= 1
	}
}

func (tpc *TaskPollController) AwaitingCount() int {
	tpc.pollLock.Lock()
	defer tpc.pollLock.Unlock()

	totalCount := 0
	for _, tp := range tpc.taskPoll {
		totalCount += tp.Count()
	}
	return totalCount
}

func (tpc *TaskPollController) Start() {
	sigTerm := utils.MakeSysChan()

	for {
		select {
		case <-sigTerm:
			utils.DInfo("task server shutted down.")
			return
		default:
			if itemIdx, tpw := tpc.populate(); tpw != nil {
				utils.DLogf("Task Poll | Task Populate, poll=%s type=%s id=%s index=%v", tpc.name, tpw.TaskName(), tpw.ID(), itemIdx)
				go func() {
					defer func() {
						utils.WrapErrorPure("Task population err", recover())
						tpc.release(tpw)
					}()
					tpw.Yield(itemIdx, tpc)
				}()
				if tpc.interval > 0 {
					time.Sleep(tpc.interval)
				}
			} else {
				// extra sleep for over-populated punishment
				time.Sleep(40 * time.Millisecond)
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (tpc *TaskPollController) Push(item TaskPollItem) TaskPollItem {
	tpc.pollLock.Lock()
	defer tpc.pollLock.Unlock()

	tpc.taskPoll = append(tpc.taskPoll, &taskPollItemWrapper{
		TaskPollItem: item,
	})

	return item
}

func (tpc *TaskPollController) remove_unsafe(id string, exitCode TaskPollExitCode) {
	var tp *taskPollItemWrapper = nil
	tpc.taskPoll = structs.Filter(tpc.taskPoll, func(w *taskPollItemWrapper) bool {
		if w.ID() == id {
			tp = w
			return false
		}
		return true
	})

	if tp != nil && exitCode != TPExitSuccess {
		tp.OnExit(exitCode)
	}
}

func (tpc *TaskPollController) Remove(id string, exitCode TaskPollExitCode) {
	tpc.pollLock.Lock()
	defer tpc.pollLock.Unlock()

	tpc.remove_unsafe(id, exitCode)
}

func NewTaskPollController(name string, concurrency uint, interval time.Duration, emptyWait time.Duration) *TaskPollController {
	if concurrency == 0 || concurrency > 256 {
		concurrency = 16
	}

	return &TaskPollController{
		name:        name,
		concurrency: concurrency,
		interval:    interval,
		emptyWait:   emptyWait,

		runningTask: make(map[string]int),
	}
}
