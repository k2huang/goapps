// 有问题版本，可能会出现死锁
package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

//叉子, 哲学家要争抢的资源
//需要用锁保护
type Fork struct {
	sync.Mutex
	ID int
}

//哲学家
type Philosopher struct {
	name   string
	first  *Fork
	second *Fork
}

func New(name string, first, second *Fork) *Philosopher {
	return &Philosopher{
		name:   name,
		first:  first,
		second: second,
	}
}

func (p *Philosopher) ThinkOrEat() {

	for {
		//模拟思考过程
		log.Println(p.name, "- 开始思考...")
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
		log.Println(p.name, "- 思考结束")

		//思考过后开始进餐
		p.first.Lock()  //拿起第一把叉子
		p.second.Lock() //拿起第二把叉子
		{
			log.Println(p.name, "- 开始进餐...")
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
			log.Println(p.name, "- 进餐结束")
		}
		p.second.Unlock()
		p.first.Unlock()
	}
}

func main() {
	var (
		forks        [5]*Fork
		philosophers [5]*Philosopher

		wg sync.WaitGroup
	)

	for i := 0; i < 5; i++ {
		forks[i] = &Fork{ID: i}
	}

	for i := 0; i < 5; i++ {
		//所有哲学家要拿起的第一把叉子正好是左手边的，
		//要拿起的第二把叉子正好是右手边的
		philosophers[i] = New(fmt.Sprintf("P%d", i), forks[i], forks[(i+1)%5])

		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			philosophers[i].ThinkOrEat()
		}(i)
	}

	wg.Wait()
}
