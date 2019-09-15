package workerpool

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/xxjwxc/public/myqueue"

	"github.com/xxjwxc/public/mylog"
)

//New 注册工作池，并设置最大并发数
//new workpool and set the max number of concurrencies
func New(max int) *WorkerPool {
	if max < 1 {
		max = 1
	}

	p := &WorkerPool{
		task:         make(chan TaskHandler, 2*max),
		errChan:      make(chan error, 1),
		waitingQueue: myqueue.New(),
	}

	go p.loop(max)
	return p
}

//SetTimeout 设置超时时间
func (p *WorkerPool) SetTimeout(timeout time.Duration) {
	p.timeout = timeout
}

//Do 添加到工作池，并立即返回
func (p *WorkerPool) Do(fn TaskHandler) {
	if p.IsClosed() { // 已关闭
		return
	}
	p.waitingQueue.Push(fn)
	//p.task <- fn
}

//DoWait 添加到工作池，并等待执行完成之后再返回
func (p *WorkerPool) DoWait(task TaskHandler) {
	if p.IsClosed() { // 已关闭
		return
	}

	doneChan := make(chan struct{})
	p.waitingQueue.Push(TaskHandler(func() error {
		defer close(doneChan)
		return task()
	}))
	<-doneChan
}

//Wait 等待工作线程执行结束
func (p *WorkerPool) Wait() error {
	p.waitingQueue.Wait() //等待队列结束
	close(p.task)
	p.wg.Wait() //等待结束
	select {
	case err := <-p.errChan:
		return err
	default:
		return nil
	}
}

//IsDone 判断是否完成 (非阻塞)
func (p *WorkerPool) IsDone() bool {
	if p == nil || p.task == nil {
		return true
	}

	return len(p.task) == 0
}

//IsClosed 是否已经关闭
func (p *WorkerPool) IsClosed() bool {
	if atomic.LoadInt32(&p.closed) == 1 { // 已关闭
		return true
	}
	return false
}

func (p *WorkerPool) startQueue() {
	for {
		fn := p.waitingQueue.Pop().(TaskHandler)
		if p.IsClosed() { // 已关闭
			p.waitingQueue.Close()
			break
		}

		if fn != nil {
			p.task <- fn
		}
	}
}

func (p *WorkerPool) loop(maxWorkersCount int) {
	go p.startQueue() //启动队列

	p.wg.Add(maxWorkersCount) // 最大的工作协程数
	// 启动max个worker
	for i := 0; i < maxWorkersCount; i++ {
		go func() {
			defer p.wg.Done()
			// worker 开始干活
			for wt := range p.task {
				if wt == nil || atomic.LoadInt32(&p.closed) == 1 { //有err 立即返回
					continue //需要先消费完了之后再返回，
				}

				closed := make(chan struct{}, 1)
				// 有设置超时,优先task 的超时
				if p.timeout > 0 {
					ct, cancel := context.WithTimeout(context.Background(), p.timeout)
					go func() {
						select {
						case <-ct.Done():
							p.errChan <- ct.Err()
							//if atomic.LoadInt32(&p.closed) != 1 {
							mylog.Error(ct.Err())
							atomic.StoreInt32(&p.closed, 1)
							cancel()
						case <-closed:
						}
					}()
				}

				err := wt() //真正执行的点
				close(closed)
				if err != nil {
					select {
					case p.errChan <- err:
						//if atomic.LoadInt32(&p.closed) != 1 {
						mylog.Error(err)
						atomic.StoreInt32(&p.closed, 1)
					default:
					}
				}
			}
		}()
	}
}
