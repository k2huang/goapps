// 验证Go标准库 sync.Mutex的公平性
package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	var m sync.Mutex
	log.SetFlags(log.Lmicroseconds)

	go func() {
		log.Println("1, start to lock...")
		m.Lock()
		log.Println("1, locked")
		m.Unlock()
		log.Println("1, unlocked")
	}()

	go func() {
		log.Println("2, start to lock...")
		m.Lock()
		log.Println("2, locked")
		m.Unlock()
		log.Println("2, unlocked")
	}()

	log.Println("main, start to lock...")
	m.Lock()
	log.Println("main, locked")
	time.Sleep(3 * time.Millisecond)
	m.Unlock()
	log.Println("main, unlocked")

	log.Println("main, start to lock again...")
	m.Lock()
	log.Println("main, locked again")
	m.Unlock()
	log.Println("main, unlocked again")

	time.Sleep(time.Second)
	log.Println("main, end")
}
