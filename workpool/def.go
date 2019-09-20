package workpool

import (
	"sync"
	"time"

	"github.com/xxjwxc/public/myqueue"
)

// TaskHandler Define function callbacks
type TaskHandler func() error

//WorkPool serves incoming connections via a pool of workers
type WorkPool struct {
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
