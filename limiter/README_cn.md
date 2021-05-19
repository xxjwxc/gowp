## 限流器


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
