//创建两个goroutine，交替打印数字和字母

package cookbook

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	// numChan 用于 num 协程发送信号给 char 协程
	// charChan 用于 char 协程发送信号给 num 协程
	numChan := make(chan struct{})
	charChan := make(chan struct{})

	// 协程1: 打印数字 (主导发送方，先发送1个信号给 charChan，然后等待 charChan 的回传)
	go func() {
		defer wg.Done()
		defer close(numChan) // 当 num 协程完成其所有发送后，关闭 numChan

		// 先发送一个信号，启动 char 协程打印 'A'
		// 在这里发送的信号是为了启动整个循环，并不是循环内的逻辑
		// charChan <- struct{}{} // 移除此行，因为 charChan 的第一次接收在循环内部
		// 我们需要从外部启动numChan，使num协程开始打印1

		for i := 1; i <= 10; i++ {
			// 等待 char 协程发送回来的信号
			// 这里的 ok 判断是必要的，因为 charChan 最终会被 char 协程关闭，以此作为退出信号
			_, ok := <-charChan
			if !ok {
				// 如果 charChan 在预期完成前关闭，则退出
				fmt.Println("num goroutine: charChan closed prematurely.")
				return
			}
			fmt.Println(i)

			// 如果不是最后一个数字，就继续发送信号给 char 协程
			// 如果是最后一个数字 (10)，则不发送，而是让 defer close(numChan) 来通知 char 协程结束
			if i < 10 {
				numChan <- struct{}{}
			}
		}
	}()

	// 协程2: 打印字母 (被动接收方，等待 numChan 的信号)
	go func() {
		defer wg.Done()
		defer close(charChan) // 当 char 协程完成其所有发送后，关闭 charChan

		for i := 'A'; i <= 'J'; i++ {
			// 等待 num 协程发送回来的信号
			_, ok := <-numChan
			if !ok {
				// 如果 numChan 在预期完成前关闭，则退出
				fmt.Println("char goroutine: numChan closed prematurely.")
				return
			}
			fmt.Println(string(i))

			// 如果不是最后一个字母，就继续发送信号给 num 协程
			if i < 'J' {
				charChan <- struct{}{}
			}
		}
	}()

	// 启动第一个信号，由外部 main 协程发起，让 num 协程开始第一个数字的打印
	// 注意：这里发送的信号是给 numChan 的，不是 charChan
	// 这样 num 协程可以打印 1，然后发送给 charChan
	numChan <- struct{}{}

	wg.Wait()
	fmt.Println("finished")
}
