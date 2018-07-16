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

type LoadBalancer struct {
}

func (l *LoadBalancer) start(RequestChan *chan Request) {
	for r := range *RequestChan {
		go worker(r)
	}
	fmt.Println("Closing channel")
}

func sendTask(RequestChan chan Request, wg *sync.WaitGroup, request *Request) {
	RequestChan <- *request
	go func(request *Request) {
		switch request.result.(type) {
		case chan int:
			ResChan := request.result.(chan int)
			fmt.Println(<-ResChan)
		case chan string:
			ResChan := request.result.(chan string)
			fmt.Println(<-ResChan)
		}
		wg.Done()
	}(request)
}

func worker(request Request) {
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
}
func main() {
	fmt.Println("Starting")

	RequestChan := make(chan Request)
	var wg sync.WaitGroup
	loadBalancer := LoadBalancer{}

	go loadBalancer.start(&RequestChan)

	req1 := Request{
		fn: func() int {
			time.Sleep(time.Second)
			return 2 * 4
		},
		result: make(chan int),
	}

	wg.Add(1)
	go sendTask(RequestChan, &wg, &req1)

	req2 := Request{
		fn: func() string {
			return "hello"
		},
		result: make(chan string),
	}

	wg.Add(1)
	go sendTask(RequestChan, &wg, &req2)

	wg.Wait()
	close(RequestChan)
}
