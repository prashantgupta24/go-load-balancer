package main

import (
	"fmt"
	"sync"
	"time"
)

type Request struct {
	fn     interface{}
	result interface{}
}

type Pool struct {
	pool []*Worker
}

type LoadBalancer struct {
	pool        Pool
	numWorkers  int
	RequestChan chan *Request
	done        chan *Worker
}

type Worker struct {
	WorkerChan chan *Request
	id         int
}

func (l *LoadBalancer) Init() {
	for i := 0; i < l.numWorkers; i++ {
		WorkerChan := make(chan *Request)
		worker := &Worker{
			WorkerChan: WorkerChan,
			id:         i,
		}
		l.pool.pool = append(l.pool.pool, worker)
		go worker.start(l.done)
	}
}

func (l *LoadBalancer) start() {
	i := 0
	for {
		select {
		case request := <-l.RequestChan:
			fmt.Println("selecting worker ", i)
			worker := l.pool.pool[i]
			i++
			worker.WorkerChan <- request
		case w := <-l.done:
			fmt.Println("worker done!", w.id)
			i--
		}
	}
	// for request := range l.RequestChan {
	// 	fmt.Println("selecting worker ", l.i)
	// 	worker := l.pool[l.i%10]
	// 	l.i++
	// 	worker.WorkerChan <- request
	// 	<-worker.done
	// }
}

// for request := range l.RequestChan {
// 	fmt.Println("selecting worker ", i)
// 	worker := l.pool[i%10]
// 	i++
// 	worker.WorkerChan <- request
// }

func (w *Worker) start(done chan (*Worker)) {
	for request := range w.WorkerChan {
		switch request.result.(type) {
		case chan int:
			ResChan := request.result.(chan int)
			fn := request.fn.(func() int)
			ResChan <- fn()
			done <- w
		case chan string:
			ResChan := request.result.(chan string)
			fn := request.fn.(func() string)
			ResChan <- fn()
			done <- w
		}
	}
}

func WaitForTask(wg *sync.WaitGroup, request *Request) {
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

	RequestChan := make(chan *Request)
	var wg sync.WaitGroup
	loadBalancer := LoadBalancer{
		pool:        Pool{},
		numWorkers:  10,
		RequestChan: RequestChan,
		done:        make(chan *Worker, 10),
	}

	loadBalancer.Init()
	go loadBalancer.start()

	for i := 0; i < 15; i++ {
		//time.Sleep(time.Second)
		req := Request{
			fn: func() int {
				time.Sleep(time.Second * 5)
				return 2 * 4
			},
			result: make(chan int),
		}

		RequestChan <- &req
		wg.Add(1)
		go WaitForTask(&wg, &req)
	}

	req2 := Request{
		fn: func() string {
			return "hello"
		},
		result: make(chan string),
	}

	RequestChan <- &req2
	wg.Add(1)
	go WaitForTask(&wg, &req2)

	wg.Wait()
}
