/*******
https://www.acwing.com/problem/content/150/
在一个果园里，达达已经将所有的果子打了下来，而且按果子的不同种类分成了不同的 n
𝑛
 堆。

达达决定把所有的果子合成一堆。

每一次合并，达达可以把两堆果子合并到一起，消耗的体力等于两堆果子的重量之和。

可以看出，所有的果子经过 n−1
𝑛
−
1
 次合并之后，就只剩下一堆了。

达达在合并果子时总共消耗的体力等于每次合并所耗体力之和。

因为还要花大力气把这些果子搬回家，所以达达在合并果子时要尽可能地节省体力。

假定每个果子重量都为 1
1
，并且已知果子的种类数和每种果子的数目，

你的任务是设计出合并的次序方案，使达达耗费的体力最少，并输出这个最小的体力耗费值。

输入格式
输入包括两行，第一行是一个整数 n
𝑛
，表示果子的种类数。

第二行包含 n
𝑛
 个整数，用空格分隔，第 i
𝑖
 个整数 hi
ℎ
𝑖
 是第 i
𝑖
 种果子的数目。

输出格式
输出包括一行，这一行只包含一个整数，也就是最小的体力耗费值。

*******/

package cookbook

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type MinHeap []int

func (h MinHeap) Len() int {
	return len(h)
}

func (h MinHeap) Less(i, j int) bool {
	return h[i] < h[j]
}

func (h MinHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MinHeap) Pop() interface{} {
	node := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]
	return node
}

func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(reader, &n)
	if n <= 1 {
		fmt.Println(0)
		return
	}
	h := &MinHeap{}
	heap.Init(h)
	for i := 0; i < n; i++ {
		var ai int
		fmt.Fscan(reader, &ai)
		heap.Push(h, ai)
	}
	cost := 0
	tmp := 0
	for h.Len() > 1 {
		//因为这里Pop出来的是接口类型，所以要进行类型断言
		//所有“通过”接口的变量都要进行类型断言
		a := heap.Pop(h).(int)
		b := heap.Pop(h).(int)
		tmp = a + b
		heap.Push(h, tmp)
		cost += tmp
	}
	fmt.Println(cost)
}
