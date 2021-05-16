package limiter

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/xxjwxc/public/myredis"
)

// template
func TestLimiterRedis(t *testing.T) {
	conf := myredis.InitRedis(myredis.WithAddr("192.155.1.150:6379"), myredis.WithClientName(""),
		myredis.WithPool(2, 0),
		myredis.WithTimeout(10*time.Second), myredis.WithReadTimeout(10*time.Second), myredis.WithWriteTimeout(10*time.Second),
		myredis.WithPwd("Niren1015"), myredis.WithGroupName("gggg"), myredis.WithDB(0), myredis.WithLog(false))
	res, err := myredis.NewRedis(conf)
	if err != nil {
		fmt.Println(err)
		return
	}

	limiter := NewLimiter(WithRedis(res), WithLimit(10), WithNamespace("ttt"), WithTsTimeout(true))

	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			token, err := limiter.Acquire(10)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(token)
			time.Sleep(1 * time.Second)
			limiter.Release(token)
		}()
	}
	wg.Wait()
	fmt.Println("down")
}
