package workerpool

import (
	"sync"
	"time"

	"github.com/xxjwxc/public/myqueue"
)

// TaskHandler process .定义函数回调体
type TaskHandler func() error

// workerPool serves incoming connections via a pool of workers
// in FILO order, i.e. the most recently stopped worker will serve the next
// incoming connection.
//
// Such a scheme keeps CPU caches hot (in theory).
type WorkerPool struct {
	//sync.Mutex
	//maxWorkersCount int //最大的工作协程数
	//start           sync.Once
	closed       int32
	errChan      chan error    //错误chan
	timeout      time.Duration //最大超时时间
	wg           sync.WaitGroup
	task         chan TaskHandler
	waitingQueue *myqueue.MyQueue
}
