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
	closed       int32
	errChan      chan error    //error chan
	timeout      time.Duration //max timeout
	wg           sync.WaitGroup
	task         chan TaskHandler
	waitingQueue *myqueue.MyQueue
}
