package main

import (
	"fmt"

	"github.com/xxjwxc/gowp/workpool"
)

func tets() {
	wp := workpool.New(5) //设置最大线程数
	fmt.Println(wp.IsDone())
	wp.DoWait(func() error {
		for j := 0; j < 10; j++ {
			fmt.Println(fmt.Sprintf("%v->\t%v", 000, j))
		}

		return nil
		// time.Sleep(1 * time.Second)
		// return errors.New("my test err")
	})

	for i := 0; i < 10; i++ { //开启10个请求
		ii := i
		wp.Do(func() error {
			for j := 0; j < 5; j++ { //每次打印0-10的值
				fmt.Println(fmt.Sprintf("%v->\t%v", ii, j))
				// if ii == 1 {
				// 	return errors.Cause(errors.New("my test err")) //有err 立即返回
				// }
			}
			return nil
		})

		fmt.Println(wp.IsDone())
	}

	wp.Wait()
	fmt.Println(wp.IsDone())
	fmt.Println(wp.IsClosed())
	fmt.Println("down")
}
func main() {
	tets()
}
