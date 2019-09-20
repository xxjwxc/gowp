package workpool

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/xxjwxc/public/myqueue"
)

// New new workpool and set the max number of concurrencies
func New(max int) *WorkPool {
	if max < 1 {
		max = 1
	}

	p := &WorkPool{
		task:         make(chan TaskHandler, 2*max),
		errChan:      make(chan error, 1),
		waitingQueue: myqueue.New(),
	}

	go p.loop(max)
	return p
}

// SetTimeout Setting timeout time
func (p *WorkPool) SetTimeout(timeout time.Duration) {
	p.timeout = timeout
}

// Do Add to the workpool and return immediately
func (p *WorkPool) Do(fn TaskHandler) {
	if p.IsClosed() { // Closed
		return
	}
	p.waitingQueue.Push(fn)
	// p.task <- fn
}

// DoWait Add to the workpool and wait for execution to complete before returning
func (p *WorkPool) DoWait(task TaskHandler) {
	if p.IsClosed() { // closed
		return
	}

	doneChan := make(chan struct{})
	p.waitingQueue.Push(TaskHandler(func() error {
		defer close(doneChan)
		return task()
	}))
	<-doneChan
}

// Wait Waiting for the worker thread to finish executing
func (p *WorkPool) Wait() error {
	p.waitingQueue.Wait() // Waiting for queue to end
	close(p.task)
	p.wg.Wait() // Waiting
	select {
	case err := <-p.errChan:
		return err
	default:
		return nil
	}
}

// IsDone Determine whether it is complete (non-blocking)
func (p *WorkPool) IsDone() bool {
	if p == nil || p.task == nil {
		return true
	}

	return p.waitingQueue.Len() == 0 && len(p.task) == 0
}

// IsClosed Has it been closed?
func (p *WorkPool) IsClosed() bool {
	if atomic.LoadInt32(&p.closed) == 1 { // closed
		return true
	}
	return false
}

func (p *WorkPool) startQueue() {
	for {
		fn := p.waitingQueue.Pop().(TaskHandler)
		if p.IsClosed() { // closed
			p.waitingQueue.Close()
			break
		}

		if fn != nil {
			p.task <- fn
		}
	}
}

func (p *WorkPool) loop(maxWorkersCount int) {
	go p.startQueue() // Startup queue

	p.wg.Add(maxWorkersCount) // Maximum number of work cycles
	// Start Max workers
	for i := 0; i < maxWorkersCount; i++ {
		go func() {
			defer p.wg.Done()
			// workering
			for wt := range p.task {
				if wt == nil || atomic.LoadInt32(&p.closed) == 1 { // returns immediately
					continue // It needs to be consumed before returning.
				}

				closed := make(chan struct{}, 1)
				// Set timeout, priority task timeout.
				if p.timeout > 0 {
					ct, cancel := context.WithTimeout(context.Background(), p.timeout)
					go func() {
						select {
						case <-ct.Done():
							p.errChan <- ct.Err()
							// if atomic.LoadInt32(&p.closed) != 1 {
							// mylog.Error(ct.Err())
							atomic.StoreInt32(&p.closed, 1)
							cancel()
						case <-closed:
						}
					}()
				}

				err := wt() // Points of Execution
				close(closed)
				if err != nil {
					select {
					case p.errChan <- err:
						// if atomic.LoadInt32(&p.closed) != 1 {
						// mylog.Error(err)
						atomic.StoreInt32(&p.closed, 1)
					default:
					}
				}
			}
		}()
	}
}
