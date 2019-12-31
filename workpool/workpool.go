package workpool

import (
	"context"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/xxjwxc/public/myqueue"
)

// New new workpool and set the max number of concurrencies
func New(max int) *WorkPool { // 注册工作池，并设置最大并发数
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
func (p *WorkPool) SetTimeout(timeout time.Duration) { // 设置超时时间
	p.timeout = timeout
}

// Do Add to the workpool and return immediately
func (p *WorkPool) Do(fn TaskHandler) { // 添加到工作池，并立即返回
	if p.IsClosed() { // 已关闭
		return
	}
	p.waitingQueue.Push(fn)
	// p.task <- fn
}

// DoWait Add to the workpool and wait for execution to complete before returning
func (p *WorkPool) DoWait(task TaskHandler) { // 添加到工作池，并等待执行完成之后再返回
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
func (p *WorkPool) Wait() error { // 等待工作线程执行结束
	p.waitingQueue.Wait()  // 等待队列结束
	p.waitingQueue.Close() //
	p.waitTask()           // wait que down
	close(p.task)
	p.wg.Wait() // 等待结束
	select {
	case err := <-p.errChan:
		return err
	default:
		return nil
	}
}

// IsDone Determine whether it is complete (non-blocking)
func (p *WorkPool) IsDone() bool { // 判断是否完成 (非阻塞)
	if p == nil || p.task == nil {
		return true
	}

	return p.waitingQueue.Len() == 0 && len(p.task) == 0
}

// IsClosed Has it been closed?
func (p *WorkPool) IsClosed() bool { // 是否已经关闭
	if atomic.LoadInt32(&p.closed) == 1 { // closed
		return true
	}
	return false
}

func (p *WorkPool) startQueue() {
	p.isQueTask = 1
	for {
		tmp := p.waitingQueue.Pop()
		if p.IsClosed() { // closed
			p.waitingQueue.Close()
			break
		}
		if tmp != nil {
			fn := tmp.(TaskHandler)
			if fn != nil {
				p.task <- fn
			}
		} else {
			break
		}

	}
	atomic.StoreInt32(&p.isQueTask, 0)
}

func (p *WorkPool) waitTask() {
	for {
		runtime.Gosched() // 出让时间片
		if p.IsDone() {
			if atomic.LoadInt32(&p.isQueTask) == 0 {
				break
			}
		}
	}
}

func (p *WorkPool) loop(maxWorkersCount int) {
	go p.startQueue() // Startup queue , 启动队列

	p.wg.Add(maxWorkersCount) // Maximum number of work cycles,最大的工作协程数
	// Start Max workers, 启动max个worker
	for i := 0; i < maxWorkersCount; i++ {
		go func() {
			defer p.wg.Done()
			// worker 开始干活
			for wt := range p.task {
				if wt == nil || atomic.LoadInt32(&p.closed) == 1 { // returns immediately,有err 立即返回
					continue // It needs to be consumed before returning.需要先消费完了之后再返回，
				}

				closed := make(chan struct{}, 1)
				// Set timeout, priority task timeout.有设置超时,优先task 的超时
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

				err := wt() // Points of Execution.真正执行的点
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
