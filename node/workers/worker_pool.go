package workers

import (
	"fmt"
	"go.uber.org/zap"
	e "node/events"
	t "node/types"
	u "node/util"
	"os"
	"strconv"
	"sync"
)

type WorkerPool struct {
	JobQueue chan t.NodeActionEvent
	Wg       sync.WaitGroup
	Logger   *zap.SugaredLogger
}

var pool *WorkerPool

func GetWorkerPool() (*WorkerPool, error) {
	if pool != nil {
		return pool, nil
	}

	l, err := u.GetLogger()
	if err != nil {
		return nil, err
	}

	numWorkersEnv := os.Getenv("WORKER_POOL_SIZE")
	if numWorkersEnv == "" {
		numWorkersEnv = "4"
	}

	var numWorkers int
	numWorkers, err = strconv.Atoi(numWorkersEnv)
	if err != nil {
		return nil, err
	}

	pool := &WorkerPool{
		JobQueue: make(chan t.NodeActionEvent, numWorkers),
		Logger:   l,
	}

	// Start the workers
	for i := 0; i < numWorkers; i++ {
		pool.Wg.Add(1)
		go pool.worker(i)
	}

	return pool, nil
}

func (p *WorkerPool) Submit(event t.NodeActionEvent) {
	p.JobQueue <- event
}

func (p *WorkerPool) worker(workerID int) {
	defer p.Wg.Done()

	for job := range p.JobQueue {
		err := e.Push(job)
		if err != nil {
			p.Logger.Error(fmt.Sprintf("Worker %d failed to write event of type %s\n", workerID, job.Type.String()))
		} else {
			p.Logger.Info(fmt.Sprintf("Worker %d successfully wrote event of type %s\n", workerID, job.Type.String()))
		}
	}
}

func (p *WorkerPool) Close() {
	close(p.JobQueue)
	p.Wg.Wait()
}
