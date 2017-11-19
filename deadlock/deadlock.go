package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	lockA sync.Mutex
	toyA  = "toyA"
)

var (
	lockB sync.Mutex
	toyB  = "toyB"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() { //小红
		defer wg.Done()
		fmt.Println("小红抢:", toyA)
		lockA.Lock()
		{
			fmt.Println("小红拿到:", toyA)
			time.Sleep(100 * time.Millisecond) //模拟业务
			{
				fmt.Println("小红抢", toyB)
				lockB.Lock()
				{
					fmt.Println("小红拿到:", toyB)
					time.Sleep(100 * time.Millisecond) //模拟业务
				}
				lockB.Unlock()
			}
		}
		lockA.Unlock()
	}()
	/*
		wg.Add(1)
		go func() { //小明 - 死锁
			defer wg.Done()
			fmt.Println("小明抢:", toyB)
			lockB.Lock()
			{
				fmt.Println("小明拿到:", toyB)
				time.Sleep(100 * time.Millisecond) //模拟业务
				{
					fmt.Println("小明抢", toyA)
					lockA.Lock()
					{
						fmt.Println("小明拿到:", toyA)
						time.Sleep(100 * time.Millisecond) //模拟业务
					}
					lockA.Unlock()
				}
			}
			lockB.Unlock()
		}()
	*/

	wg.Add(1)
	go func() { //小明 - 与小红一样: 按照 A->B 的顺序抢玩具
		defer wg.Done()
		fmt.Println("小明抢:", toyA)
		lockA.Lock()
		{
			fmt.Println("小明拿到:", toyA)
			time.Sleep(100 * time.Millisecond) //模拟业务
			{
				fmt.Println("小明抢", toyB)
				lockB.Lock()
				{
					fmt.Println("小明拿到:", toyB)
					time.Sleep(100 * time.Millisecond) //模拟业务
				}
				lockB.Unlock()
			}
		}
		lockA.Unlock()
	}()

	wg.Wait()
}
