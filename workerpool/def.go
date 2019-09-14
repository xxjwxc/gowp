package workerpool

import (
	"sync"
	"time"
)

// // CallHandler process .定义调用回调体(可修改)
// type CallHandler func()

// TaskHandler process .定义函数回调体
type TaskHandler func() error

// ServeHandler must process tls.Config.NextProto negotiated requests.
//type ServeHandler func(c net.Conn) error

// workerPool serves incoming connections via a pool of workers
// in FILO order, i.e. the most recently stopped worker will serve the next
// incoming connection.
//
// Such a scheme keeps CPU caches hot (in theory).
type WorkerPool struct {
	//sync.Mutex
	maxWorkersCount int //最大的工作协程数
	closed          int32
	errChan         chan error    //错误chan
	timeout         time.Duration //最大超时时间
	wg              sync.WaitGroup
	task            chan TaskHandler
	start           sync.Once
}
