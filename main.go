package main

import (
	"fmt"
	"sync"
	"time"
)

type request struct {
	fn     interface{}
	result interface{}
}

type loadBalancer struct {
	pool        []*worker
	numWorkers  int
	RequestChan chan *request
}

type worker struct {
	RequestChan chan *request
	id          int
}

func (l *loadBalancer) start() {
	WorkerChan := make(chan *worker, l.numWorkers)
	for i := 1; i <= l.numWorkers; i++ {
		worker := &worker{
			RequestChan: make(chan *request),
			id:          i,
		}
		l.pool = append(l.pool, worker)
		WorkerChan <- worker
		go worker.start(WorkerChan)
	}

	for request := range l.RequestChan {
		select {
		case worker := <-WorkerChan:
			fmt.Println("selecting worker ", worker.id)
			worker.RequestChan <- request
		}
	}
	fmt.Println("Shutting down load balancer")

	for _, worker := range l.pool {
		close(worker.RequestChan)
	}
}

func (w *worker) start(WorkerChan chan *worker) {
	for request := range w.RequestChan {
		switch request.result.(type) {
		case chan int:
			ResChan := request.result.(chan int)
			fn := request.fn.(func() int)
			ResChan <- fn()
		case chan string:
			ResChan := request.result.(chan string)
			fn := request.fn.(func() string)
			ResChan <- fn()
		}
		WorkerChan <- w
	}
	fmt.Println("Closing worker ", w.id)
}

func waitForTask(wg *sync.WaitGroup, request *request) {
	switch request.result.(type) {
	case chan int:
		ResChan := request.result.(chan int)
		fmt.Println(<-ResChan)
	case chan string:
		ResChan := request.result.(chan string)
		fmt.Println(<-ResChan)
	}
	wg.Done()
}

func main() {
	fmt.Println("Starting...")

	RequestChan := make(chan *request)
	var wg sync.WaitGroup

	loadBalancer := loadBalancer{
		pool:        []*worker{},
		numWorkers:  10,
		RequestChan: RequestChan,
	}

	go loadBalancer.start()

	defer func() {
		close(loadBalancer.RequestChan)
		time.Sleep(time.Second * 4)
	}()

	for i := 1; i <= 10; i++ {
		//time.Sleep(time.Second)
		j := i
		req := request{
			fn: func() int {
				time.Sleep(time.Second * 2)
				return j
			},
			result: make(chan int),
		}

		RequestChan <- &req
		wg.Add(1)
		go waitForTask(&wg, &req)
	}

	req2 := request{
		fn: func() string {
			return "hello"
		},
		result: make(chan string),
	}

	RequestChan <- &req2
	wg.Add(1)
	go waitForTask(&wg, &req2)

	wg.Wait()

}
