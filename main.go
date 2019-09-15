package main

import (
	"fmt"
	"time"

	"github.com/xxjwxc/gowp/workerpool"
)

func main() {
	wp := workerpool.New(5)   //设置最大线程数
	for i := 0; i < 10; i++ { //开启10个请求
		ii := i
		wp.Do(func() error {
			for j := 0; j < 5; j++ { //每次打印0-10的值
				fmt.Println(fmt.Sprintf("%v->\t%v", ii, j))
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
