package workerpool

import (
	"fmt"
	"testing"
	"time"

	"github.com/xxjwxc/public/errors"
)

//template
func TestWorkerPoolStart(t *testing.T) {
	wp := New(10) //Set the maximum number of threads，设置最大线程数
	wp.SetTimeout(time.Millisecond)
	for i := 0; i < 20; i++ { //Open 20 requests 开启20个请求
		ii := i
		wp.Do(func() error {
			for j := 0; j < 10; j++ {
				fmt.Println(fmt.Sprintf("%v->\t%v", ii, j))
				time.Sleep(1 * time.Millisecond)
			}
			//time.Sleep(1 * time.Second)
			return nil
		})
	}

	wp.Wait()
	fmt.Println("down")
}

//Support for error return
//支持错误返回
func TestWorkerPoolError(t *testing.T) {
	wp := New(10) //Set the maximum number of threads，设置最大线程数
	for i := 0; i < 20; i++ {
		ii := i
		wp.Do(func() error {
			for j := 0; j < 10; j++ {
				fmt.Println(fmt.Sprintf("%v->\t%v", ii, j))
				if ii == 1 {
					return errors.Cause(errors.New("my test err")) //有err 立即返回
				}
				time.Sleep(1 * time.Millisecond)
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

//Determine whether completion (non-blocking) is placed in the workpool and wait for execution results
//放到工作池里面 且等待执行结果
func TestWorkerPoolDoWait(t *testing.T) {
	wp := New(5) //Set the maximum number of threads，设置最大线程数
	for i := 0; i < 10; i++ {
		ii := i
		wp.DoWait(func() error {
			for j := 0; j < 5; j++ { //每次打印0-10的值
				fmt.Println(fmt.Sprintf("%v->\t%v", ii, j))
				// if ii == 1 {
				// 	return errors.New("my test err")
				// }
				time.Sleep(1 * time.Millisecond)
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

//Determine whether it is complete (non-blocking)
//判断是否完成 (非阻塞)
func TestWorkerPoolIsDone(t *testing.T) {
	wp := New(5) //Set the maximum number of threads，设置最大线程数
	for i := 0; i < 10; i++ {
		//	ii := i
		wp.Do(func() error {
			for j := 0; j < 5; j++ {
				//fmt.Println(fmt.Sprintf("%v->\t%v", ii, j))
				time.Sleep(1 * time.Millisecond)
			}
			return nil
		})

		fmt.Println(wp.IsDone())
	}
	wp.Wait()
	fmt.Println(wp.IsDone())
	fmt.Println("down")
}
