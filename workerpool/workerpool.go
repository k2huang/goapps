package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type Runnable interface {
	Run()
}

// 适配器: func() -> Runnable
var _ Runnable = fnRunner(nil)

type fnRunner func()

func (f fnRunner) Run() {
	f()
}

// goroutine pool
type WorkerPool struct {
	maxWorkersCount       int
	MaxIdleWorkerDuration time.Duration

	lock         sync.Mutex
	workersCount int
	mustStop     bool

	ready []*workerChan

	stopCh chan struct{}

	workerChanPool sync.Pool

	wg sync.WaitGroup //记录当前WorkerPool的所有worker goroutine
}

// 用来管理WorkerPool中的goroutine
// 一个workerChan对应goroutine
type workerChan struct {
	lastUseTime time.Time
	taskCh      chan Runnable //用于向对应的goroutine传递任务
}

func New(poolSize int) *WorkerPool {
	return &WorkerPool{
		maxWorkersCount: poolSize,
	}
}

func (wp *WorkerPool) Start() {
	if wp.stopCh != nil {
		panic("WorkerPool already started")
	}

	wp.stopCh = make(chan struct{})

	// 定期清除(关闭)长时间空闲的goroutine
	go func() {
		for {
			wp.cleanup()

			select {
			case <-wp.stopCh:
				return
			default:
				time.Sleep(wp.getMaxIdleWorkerDuration())
			}
		}
	}()
}

func (wp *WorkerPool) Stop() {
	if wp.stopCh == nil {
		panic("WorkerPool wasn't started")
	}

	close(wp.stopCh)
	// 通过wp.ready停止所有空闲goroutine
	// 正在工作的goroutine将在工作完成之后，通过判断wp.mustStop == true自动退出
	wp.lock.Lock()
	ready := wp.ready
	for i, wc := range ready {
		wc.taskCh <- nil
		ready[i] = nil
	}
	wp.ready = ready[:0]
	wp.mustStop = true
	wp.lock.Unlock()

	// 等待当前WorkerPool所有worker goroutine都退出
	wp.wg.Wait()
}

func (wp *WorkerPool) getMaxIdleWorkerDuration() time.Duration {
	if wp.MaxIdleWorkerDuration <= 0 {
		return 10 * time.Second
	}
	return wp.MaxIdleWorkerDuration
}

func (wp *WorkerPool) cleanup() {
	maxIdleWorkerDuration := wp.getMaxIdleWorkerDuration()

	// 找出所有空闲时间大于maxIdleWorkerDuration
	// 的goroutine所对应的workerChan
	var idleList []*workerChan
	currentTime := time.Now()

	wp.lock.Lock()
	fmt.Printf("there are (%d) goroutine(s) in workerpool\n", wp.workersCount)
	ready := wp.ready
	n := len(ready)
	i := 0
	for i < n && currentTime.Sub(ready[i].lastUseTime) > maxIdleWorkerDuration {
		i++
	}
	idleList = append(idleList, ready[:i]...)
	if i > 0 {
		m := copy(ready, ready[i:])
		for i = m; i < n; i++ {
			ready[i] = nil
		}
		wp.ready = ready[:m]
	}
	wp.lock.Unlock()

	// 通过workerChan停止对应的goroutine
	for _, wc := range idleList {
		wc.taskCh <- nil
	}
}

func (wp *WorkerPool) Execute(r Runnable) bool {
	ch := wp.getCh()
	if ch == nil { //WorkerPool中的goroutine数量到已达最大，且没有空闲的
		return false
	}

	ch.taskCh <- r
	return true
}

func (wp *WorkerPool) ExecuteFunc(f func()) bool {
	return wp.Execute(fnRunner(f))
}

var workerChanCap = func() int {
	if runtime.GOMAXPROCS(0) == 1 {
		return 0
	}

	return 1
}()

func (wp *WorkerPool) getCh() *workerChan {
	var ch *workerChan
	createWorker := false

	wp.lock.Lock()
	ready := wp.ready
	n := len(ready) - 1
	if n < 0 {
		if wp.workersCount < wp.maxWorkersCount {
			createWorker = true
			wp.workersCount++
		}
	} else {
		ch = ready[n]
		ready[n] = nil
		wp.ready = ready[:n]
	}
	wp.lock.Unlock()

	if ch == nil {
		if !createWorker {
			return nil
		}

		newCh := wp.workerChanPool.Get()
		if newCh == nil {
			ch = &workerChan{
				taskCh: make(chan Runnable, workerChanCap),
			}
		} else {
			ch = newCh.(*workerChan)
		}

		// 一个goroutine对应一个workerChan
		wp.wg.Add(1)
		go func() {
			defer wp.wg.Done()
			defer wp.workerChanPool.Put(ch)

			wp.workerFunc(ch)
		}()
	}

	return ch
}

func (wp *WorkerPool) release(ch *workerChan) bool {
	ch.lastUseTime = time.Now()
	wp.lock.Lock()
	if wp.mustStop {
		wp.lock.Unlock()
		return false
	}

	wp.ready = append(wp.ready, ch)
	wp.lock.Unlock()
	return true
}

func (wp *WorkerPool) workerFunc(ch *workerChan) {
	for r := range ch.taskCh { //等待有任务进来
		if r == nil {
			break
		}

		// 执行任务
		r.Run()

		// 任务执行结束，将ch放回wp.ready中，等待新的任务使用
		if !wp.release(ch) {
			break
		}
	}

	wp.lock.Lock()
	wp.workersCount--
	wp.lock.Unlock()
}
