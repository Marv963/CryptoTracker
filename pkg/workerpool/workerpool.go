package workerpool

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type Task struct {
	Ws      *websocket.Conn
	Data    interface{}
	Pair    string
	Execute func() error
	Error   error
}

type WorkerPool struct {
	Tasks    chan Task
	Shutdown chan struct{}
	Errors   chan Task // added an error channel
	Wait     sync.WaitGroup
}

func NewWorkerPool(maxWorkers int) *WorkerPool {
	pool := &WorkerPool{
		Tasks:    make(chan Task, 100), // buffer size of 100
		Shutdown: make(chan struct{}),
		Errors:   make(chan Task, 100), // buffer size of 100
	}

	pool.Wait.Add(maxWorkers)
	for i := 0; i < maxWorkers; i++ {
		go pool.worker()
	}

	return pool
}

func (pool *WorkerPool) worker() {
	defer pool.Wait.Done()

	for {
		select {
		case task := <-pool.Tasks:
			// Handle task
			if err := task.Execute(); err != nil {
				fmt.Println("Task execution error:", err)
				task.Error = err
				pool.Errors <- task // send error task back for handling
			}
		case <-pool.Shutdown:
			return
		}
	}
}

func (pool *WorkerPool) Stop() {
	close(pool.Shutdown)
	pool.Wait.Wait()
}
