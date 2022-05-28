[![Build Status](https://travis-ci.org/xxjwxc/gowp.svg?branch=master)](https://travis-ci.org/xxjwxc/gowp)
[![Go Report Card](https://goreportcard.com/badge/github.com/xxjwxc/gowp)](https://goreportcard.com/report/github.com/xxjwxc/gowp)
[![codecov](https://codecov.io/gh/xxjwxc/gowp/branch/master/graph/badge.svg)](https://codecov.io/gh/xxjwxc/gowp)
[![GoDoc](https://godoc.org/github.com/xxjwxc/gowp?status.svg)](https://godoc.org/github.com/xxjwxc/gowp)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)


## golang worker pool ,线程池 , 工作池
### [English](README_cn.md)
- 并发限制goroutine池。
- 限制任务执行的并发性，而不是排队的任务数。
- 无论排队多少任务，都不会阻止提交任务。
- 通过队列支持[queue](https://github.com/xxjwxc/public/tree/master/myqueue)
- 限流器 限制并发数

### golang 工作池公共库

## 安装

安装最简单方法:

```
$ go get github.com/xxjwxc/gowp
```

### 支持最大任务数, 放到工作池里面 并等待全部完成

```go
package main

import (
	"fmt"
	"time"

	"github.com/xxjwxc/gowp/workpool"
)

func main() {
	wp := workpool.New(10)             //设置最大线程数
	for i := 0; i < 20; i++ { //开启20个请求
		ii := i
		wp.Do(func() error {
			for j := 0; j < 10; j++ { //每次打印0-10的值
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

### 支持错误返回

```go
package main

import (
	"fmt"
	"time"

	"github.com/xxjwxc/gowp/workpool"
)

func main() {
	wp := workpool.New(10)             //设置最大线程数
	for i := 0; i < 20; i++ { //开启20个请求
		ii := i
		wp.Do(func() error {
			for j := 0; j < 10; j++ { //每次打印0-10的值
				fmt.Println(fmt.Sprintf("%v->\t%v", ii, j))
				if ii == 1 {
					return errors.Cause(errors.New("my test err")) //有err 立即返回
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

### 支持判断是否完成 (非阻塞)

```go
package main

import (
	"fmt"
	"time"

	"github.com/xxjwxc/gowp/workpool"
)

func main() {
	wp := workpool.New(5)              //设置最大线程数
	for i := 0; i < 10; i++ { //开启20个请求
		//	ii := i
		wp.Do(func() error {
			for j := 0; j < 5; j++ { 
				//fmt.Println(fmt.Sprintf("%v->\t%v", ii, j))
				time.Sleep(1 * time.Second)
			}
			return nil
		})

		fmt.Println(wp.IsDone())//判断是否完成
	}
	wp.Wait()
	fmt.Println(wp.IsDone())
	fmt.Println("down")
}
```

### 支持同步等待结果

```go
package main

import (
	"fmt"
	"time"

	"github.com/xxjwxc/gowp/workpool"
)

func main() {
	wp := workpool.New(5)              //设置最大线程数
	for i := 0; i < 10; i++ { //开启20个请求
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

### 支持超时退出

```go
package main

import (
	"fmt"
	"time"
	"time"
	"github.com/xxjwxc/gowp/workpool"
)

func main() {
	wp := workpool.New(5)              // 设置最大线程数
		wp.SetTimeout(time.Millisecond) // 设置超时时间
	for i := 0; i < 10; i++ { // 开启20个请求
		ii := i
		wp.DoWait(func() error {
			for j := 0; j < 5; j++ {
				fmt.Println(fmt.Sprintf("%v->\t%v", ii, j))
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

## 限流器(cache)

```go
package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/xxjwxc/gowp/limiter"
)

func main() {
	limiter := limiter.NewLimiter(limiter.WithLimit(10), limiter.WithNamespace("test"), limiter.WithTsTimeout(true) /*, limiter.WithRedis(res)*/)
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			token, _ := limiter.Acquire(10) // 获取一个
			fmt.Println(token)

			time.Sleep(1 * time.Second)
			limiter.Release(token) // 回退
		}()
	}
	wg.Wait()
	fmt.Println("down")
}
```
## 限流器(redis)

```go
package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/xxjwxc/gowp/limiter"
	"github.com/xxjwxc/public/myredis"
)

func main() {
	conf := myredis.InitRedis(myredis.WithAddr("127.0.0.1:6379"), myredis.WithPwd("123456"), myredis.WithGroupName("test"))
	res, err := myredis.NewRedis(conf)
	if err != nil {
		fmt.Println(err)
		return
	}

	limiter := limiter.NewLimiter(limiter.WithRedis(res), limiter.WithLimit(10), limiter.WithNamespace("test") /*, limiter.WithRedis(res)*/)
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			token, _ := limiter.Acquire(10) // 获取一个
			fmt.Println(token)

			time.Sleep(1 * time.Second)
			limiter.Release(token) // 回退
		}()
	}
	wg.Wait()
	fmt.Println("down")
}
```
