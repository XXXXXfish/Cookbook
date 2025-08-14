package main

import (
	"container/list"
	"fmt"
	"time"
)

/***************
 * 实  现  部  分（惰性删除版 LRU + TTL）
 ***************/

type LRUCache struct {
	Cap  int
	Keys map[int]*list.Element
	List *list.List
}

type pair struct {
	k        int
	v        int
	expireAt time.Time // 绝对过期时间
}

func Constructor(capacity int) *LRUCache {
	return &LRUCache{
		Cap:  capacity,
		Keys: make(map[int]*list.Element),
		List: list.New(),
	}
}

// Get：命中则移动到头部；若已过期则惰性删除并返回 -1
func (c *LRUCache) Get(key int) int {
	if element, ok := c.Keys[key]; ok {
		p := element.Value.(pair)
		// 惰性删除：访问时检查是否过期
		if time.Now().After(p.expireAt) {
			c.List.Remove(element)
			delete(c.Keys, p.k)
			return -1
		}
		c.List.MoveToFront(element)
		return p.v
	}
	return -1
}

// Put：写入 key/value，并为该键设置 TTL（仅惰性删除，无后台清理）
func (c *LRUCache) Put(key int, value int, ttl time.Duration) {
	// 非正 TTL 视为立即过期（也可改为直接 return，不写入）
	if ttl <= 0 {
		ttl = time.Nanosecond
	}
	exp := time.Now().Add(ttl)

	if element, ok := c.Keys[key]; ok {
		// 键已存在：如果旧值已过期，当作新写入；否则更新并前移
		p := element.Value.(pair)
		if time.Now().After(p.expireAt) {
			c.List.Remove(element)
			delete(c.Keys, p.k)
			ne := c.List.PushFront(pair{k: key, v: value, expireAt: exp})
			c.Keys[key] = ne
		} else {
			element.Value = pair{k: key, v: value, expireAt: exp}
			c.List.MoveToFront(element)
		}
	} else {
		// 新键
		element := c.List.PushFront(pair{k: key, v: value, expireAt: exp})
		c.Keys[key] = element
	}

	// 容量控制：标准 LRU，从尾部淘汰
	for c.List.Len() > c.Cap {
		back := c.List.Back()
		if back == nil {
			break
		}
		p := back.Value.(pair)
		c.List.Remove(back)
		delete(c.Keys, p.k)
	}
}

/***************
 * 操  作  示  例
 ***************/

func main() {
	fmt.Println("== Demo 1: LRU + TTL（容量=2）")
	cache := Constructor(2)

	fmt.Println("Put(1,100, 500ms)")
	cache.Put(1, 100, 500*time.Millisecond)

	fmt.Println("Put(2,200, 2s)")
	cache.Put(2, 200, 2*time.Second)

	fmt.Println("Get(1) =>", cache.Get(1)) // 命中，刷新 LRU，输出 100

	fmt.Println("Put(3,300, 2s)  // 容量溢出，应淘汰最久未使用的 key=2")
	cache.Put(3, 300, 2*time.Second)

	fmt.Println("Get(2) =>", cache.Get(2)) // -1（被 LRU 淘汰）
	fmt.Println("Get(3) =>", cache.Get(3)) // 300

	fmt.Println("Sleep 600ms 等待 key=1 过期（它的 TTL=500ms）")
	time.Sleep(600 * time.Millisecond)

	fmt.Println("Get(1) =>", cache.Get(1)) // -1（惰性删除触发，发现已过期）
	fmt.Println("Get(3) =>", cache.Get(3)) // 300（仍然有效）

	fmt.Println("\n== Demo 2: 通过 Put 重置 TTL（续期）")
	cache2 := Constructor(2)

	fmt.Println("Put(10,1000, 120ms)")
	cache2.Put(10, 1000, 120*time.Millisecond)

	fmt.Println("Sleep 80ms")
	time.Sleep(80 * time.Millisecond)

	fmt.Println("Put(10,1001, 250ms)  // 同 key 再写，值=1001，并重置 TTL")
	cache2.Put(10, 1001, 250*time.Millisecond)

	fmt.Println("Sleep 150ms  // 现在距首次写入约 230ms，续期后的 TTL 仍未到")
	time.Sleep(150 * time.Millisecond)

	fmt.Println("Get(10) =>", cache2.Get(10)) // 1001（有效）

	fmt.Println("Sleep 120ms  // 230+120=350ms > 续期后的到期点，下一次访问应过期")
	time.Sleep(120 * time.Millisecond)

	fmt.Println("Get(10) =>", cache2.Get(10)) // -1（本次访问触发惰性删除）

	fmt.Println("\n== Demo 3: 过期后重新写入")
	fmt.Println("Put(1,100, 50ms)")
	cache3 := Constructor(2)
	cache3.Put(1, 100, 50*time.Millisecond)

	fmt.Println("Sleep 70ms  // 等它过期")
	time.Sleep(70 * time.Millisecond)

	fmt.Println("Put(1,101, 300ms)  // 过期后重写，相当于新键")
	cache3.Put(1, 101, 300*time.Millisecond)

	fmt.Println("Get(1) =>", cache3.Get(1)) // 101
}
