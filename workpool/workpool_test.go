package workpool

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xxjwxc/public/errors"
)

// template
func TestWorkerPoolStart(t *testing.T) {
	wp := New(10) // Set the maximum number of threads
	wp.SetTimeout(time.Millisecond)
	for i := 0; i < 20; i++ { // Open 20 requests
		ii := i
		wp.Do(func() error {
			for j := 0; j < 10; j++ {
				fmt.Printf("%v->\t%v", ii, j)
				time.Sleep(1 * time.Millisecond)
			}
			// time.Sleep(1 * time.Second)
			return nil
		})
	}

	err := wp.Wait()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("down")
}

// Support for error return
func TestWorkerPoolError(t *testing.T) {
	wp := New(10) // Set the maximum number of threads
	for i := 0; i < 20; i++ {
		ii := i
		wp.Do(func() error {
			for j := 0; j < 10; j++ {
				fmt.Printf("%v->\t%v", ii, j)
				if ii == 1 {
					return errors.Cause(errors.New("my test err"))
				}
				time.Sleep(1 * time.Millisecond)
			}

			return nil
			// time.Sleep(1 * time.Second)
			// return errors.New("my test err")
		})
	}

	err := wp.Wait()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("down")
}

// Determine whether completion (non-blocking) is placed in the workpool and wait for execution results
func TestWorkerPoolDoWait(t *testing.T) {
	wp := New(5) // Set the maximum number of threads
	for i := 0; i < 10; i++ {
		ii := i
		wp.DoWait(func() error {
			for j := 0; j < 5; j++ {
				fmt.Printf("%v->\t%v", ii, j)
				// if ii == 1 {
				// 	return errors.New("my test err")
				// }
				time.Sleep(1 * time.Millisecond)
			}

			return nil
			// time.Sleep(1 * time.Second)
			// return errors.New("my test err")
		})
	}

	err := wp.Wait()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("down")
}

// Determine whether it is complete (non-blocking)
func TestWorkerPoolIsDone(t *testing.T) {
	wp := New(5) // Set the maximum number of threads
	for i := 0; i < 10; i++ {
		//	ii := i
		wp.Do(func() error {
			for j := 0; j < 5; j++ {
				// fmt.Println(fmt.Sprintf("%v->\t%v", ii, j))
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

func TestDoBefore(t *testing.T) {
	wp := New(1) // Set the maximum number of threads

	var res int32 = 0

	wp.Do(func() error {
		time.Sleep(5 * time.Second)
		atomic.AddInt32(&res, 1)
		return nil
	})

	wp.DoBefore(func() error {
		atomic.AddInt32(&res, 1)
		return nil
	}, 5*time.Second)

	wp.Wait()
	assert.Equal(t, int32(1), res)
}
