package workerpool

import (
	"fmt"
	"testing"
	"time"

	"github.com/xxjwxc/public/errors"
)

func TestWorkerPoolStart(t *testing.T) {
	wp := New(10)             //设置最大线程数
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

func TestWorkerPoolError(t *testing.T) {
	wp := New(10)             //设置最大线程数
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

//放到工作池里面 且等待执行结果
func TestWorkerPoolDoWait(t *testing.T) {
	wp := New(5)              //设置最大线程数
	for i := 0; i < 10; i++ { //开启20个请求
		ii := i
		wp.DoWait(func() error {
			for j := 0; j < 5; j++ { //每次打印0-10的值
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

//判断是否完成 (非阻塞)
func TestWorkerPoolIsDone(t *testing.T) {
	wp := New(5)              //设置最大线程数
	for i := 0; i < 10; i++ { //开启20个请求
		//	ii := i
		wp.Do(func() error {
			for j := 0; j < 5; j++ { //每次打印0-10的值
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
