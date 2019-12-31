[![Build Status](https://travis-ci.org/xxjwxc/gowp.svg?branch=master)](https://travis-ci.org/xxjwxc/gowp)
[![Go Report Card](https://goreportcard.com/badge/github.com/xxjwxc/gowp)](https://goreportcard.com/report/github.com/xxjwxc/gowp)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)

## golang worker pool

### [中文说明](README_cn.md)

- Concurrency limiting goroutine pool. 
- Limits the concurrency of task execution, not the number of tasks queued. 
- Never blocks submitting tasks, no matter how many tasks are queued.
- Support through security queues [queue](https://github.com/xxjwxc/public/tree/master/myqueue)

- golang workpool common library

## Installation

The simplest way to install the library is to run:

```
$ go get github.com/xxjwxc/gowp
```


### Support the maximum number of tasks, put them in the workpool and wait for them to be completed

## Example

```
package main

import (
	"fmt"
	"time"

	"github.com/xxjwxc/gowp/workpool"
)

func main() {
	wp := workpool.New(10)     // Set the maximum number of threads
	for i := 0; i < 20; i++ { // Open 20 requests 
		ii := i
		wp.Do(func() error {
			for j := 0; j < 10; j++ { // 0-10 values per print
				fmt.Println(fmt.Sprintf("%v->\t%v", ii, j))
				time.Sleep(1 * time.Second)
			}
			//time.Sleep(1 * time.Second)
			return nil
		})
	}

	wp.Wait()
	fmt.Println("down")
}

```

### Support for error return

```
package main

import (
	"fmt"
	"time"

	"github.com/xxjwxc/gowp/workpool"
)

func main() {
	wp := workpool.New(10)             // Set the maximum number of threads
	for i := 0; i < 20; i++ { 
		ii := i
		wp.Do(func() error {
			for j := 0; j < 10; j++ { // 0-10 values per print
				fmt.Println(fmt.Sprintf("%v->\t%v", ii, j))
				if ii == 1 {
					return errors.Cause(errors.New("my test err")) // have err return
				}
				time.Sleep(1 * time.Second)
			}

			return nil
		})
	}

	err := wp.Wait()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("down")
	}
```

### Supporting judgement of completion (non-blocking)

```
package main

import (
	"fmt"
	"time"

	"github.com/xxjwxc/gowp/workpool"
)

func main() {
	wp := workpool.New(5)              // Set the maximum number of threads
	for i := 0; i < 10; i++ { 
		//	ii := i
		wp.Do(func() error {
			for j := 0; j < 5; j++ { 
				//fmt.Println(fmt.Sprintf("%v->\t%v", ii, j))
				time.Sleep(1 * time.Second)
			}
			return nil
		})

		fmt.Println(wp.IsDone())
	}
	wp.Wait()
	fmt.Println(wp.IsDone())
	fmt.Println("down")
}
```

### Support synchronous waiting for results

```
package main

import (
	"fmt"
	"time"

	"github.com/xxjwxc/gowp/workpool"
)

func main() {
	wp := workpool.New(5) //Set the maximum number of threads
	for i := 0; i < 10; i++ { 
		ii := i
		wp.DoWait(func() error {
			for j := 0; j < 5; j++ { 
				fmt.Println(fmt.Sprintf("%v->\t%v", ii, j))
				// if ii == 1 {
				// 	return errors.New("my test err")
				// }
				time.Sleep(1 * time.Second)
			}

			return nil
			//time.Sleep(1 * time.Second)
			//return errors.New("my test err")
		})
	}

	err := wp.Wait()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("down")
}
```
