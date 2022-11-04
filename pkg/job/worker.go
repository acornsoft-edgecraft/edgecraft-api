/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package job

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/logger"
)

var (
	ErrWorkerBusy = errors.New("workers are busy, try again later")
)

// TaskInfo - 실제 실행할 Task 정보
type TaskInfo struct {
	TaskData interface{}
	TaskFunc func(task string, taskData interface{})
}

// workType - Worker types
type workType struct {
	TaskID       string
	WorkDuration time.Duration
	Data         TaskInfo
}

// worker - Background worker info
type worker struct {
	workchan    chan workType
	workerCount int
	buffer      int
	wg          *sync.WaitGroup
	cancelFunc  context.CancelFunc
}

// Start - Worker start
func (w *worker) Start(pctx context.Context) {
	ctx, cancelFunc := context.WithCancel(pctx)
	w.cancelFunc = cancelFunc

	for i := 0; i < w.workerCount; i++ {
		w.wg.Add(1)
		go w.spawnWorkers(ctx)
	}

	logger.Infof("worker started [count: %d, buffer: %d]", w.workerCount, w.buffer)
}

// Stop - Worker stop
func (w *worker) Stop() {
	logger.Info("stop workers")
	close(w.workchan)
	w.cancelFunc()
	w.wg.Wait()
	logger.Info("all workers exited!")
}

// QueueTask - Queuing task
func (w *worker) QueueTask(task string, workDuration time.Duration, taskData TaskInfo) error {
	if len(w.workchan) >= w.buffer {
		return ErrWorkerBusy
	}

	w.workchan <- workType{TaskID: task, WorkDuration: workDuration, Data: taskData}
	return nil
}

// spawnWorkers - Spawn workers
func (w *worker) spawnWorkers(ctx context.Context) {
	defer w.wg.Done()

	for work := range w.workchan {
		select {
		case <-ctx.Done():
			return
		default:
			w.doWork(ctx, work.TaskID, work.WorkDuration, work.Data)
		}
	}
}

// doWork - Do work
func (w *worker) doWork(ctx context.Context, task string, workDuration time.Duration, taskInfo TaskInfo) {
	//sleepContext(ctx, workDuration)
	if taskInfo.TaskData == nil || taskInfo.TaskFunc == nil {
		logger.WithField("task", task).Info("Cannot start work. taskinfo not specified.")
	} else {
		logger.WithField("task", task).Info("Start requested work.")
		taskInfo.TaskFunc(task, taskInfo.TaskData)
		logger.WithField("task", task).Info("Complete requested work.")
	}
}

// IWorker - Background worker interface
type IWorker interface {
	Start(pctx context.Context)
	Stop()
	QueueTask(task string, workDuration time.Duration, taskData TaskInfo) error
}

// // sleepContext - Sleeping context
// func sleepContext(ctx context.Context, sleep time.Duration) {
// 	select {
// 	case <-ctx.Done():
// 	case <-time.After(sleep):
// 	}
// }

// NewWorker - Create instances of worker
func NewWorker(workerCount, buffer int) IWorker {
	w := worker{
		workchan:    make(chan workType, buffer),
		workerCount: workerCount,
		buffer:      buffer,
		wg:          new(sync.WaitGroup),
	}
	return &w
}
