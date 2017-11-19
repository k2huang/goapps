package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var (
		mtx     sync.Mutex
		counter int
		wg      sync.WaitGroup // 用于统计正在运行Goroutine的数量
	)

	start := time.Now()
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 10000; i++ {
				mtx.Lock()
				counter++
				mtx.Unlock()
			}
		}()
	}
	wg.Wait() // 等待所以Goroutine都执行完
	elapsed := time.Since(start)
	fmt.Println("counter:", counter)
	fmt.Println("time consume:", elapsed.Seconds(), "s")
}
