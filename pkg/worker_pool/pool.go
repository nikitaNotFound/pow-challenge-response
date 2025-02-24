package worker_pool

import "context"

type work func()

type WorkerPool struct {
	workersAmount int
	ctx           context.Context

	workChannel chan work
}

func NewWorkerPool(workersAmount int, ctx context.Context) *WorkerPool {
	return &WorkerPool{
		workersAmount: workersAmount,
		ctx:           ctx,
		workChannel:   make(chan work),
	}
}

func (p *WorkerPool) Start() {
	for i := 0; i < p.workersAmount; i++ {
		go p.worker(p.ctx)
	}
}

func (p *WorkerPool) RunWork(worker work) {
	select {
	case <-p.ctx.Done():
		return
	case p.workChannel <- worker:
	}
}

func (p *WorkerPool) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case work := <-p.workChannel:
			work()
		}
	}
}
