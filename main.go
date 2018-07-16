package main

import (
	"fmt"
	"time"
)

type Request struct {
	//make()
	// getFn()
	// getResChan()
	fn     interface{}
	result interface{}
}

type IntReq struct {
	fn     func() int
	result chan int
}

// func (r *IntReq) getFn() {
// 	return r.fn
// }

type StringReq struct {
	fn     func() string
	result chan string
}

// func (r *Request) IntReq(fn func() int, result chan int){
// 	fn     func() int
// 	result chan int
// }
type LoadBalancer struct {
}

func (l *LoadBalancer) start(RequestChan chan Request) {
	for {
		select {
		case r := <-RequestChan:
			{
				go worker(r)
			}
		}
	}
}

func worker(request Request) {
	//request.result <- request.fn()
	switch request.result.(type) {
	case chan int:
		ResChan := request.result.(chan int)
		fn := request.fn.(func() int)
		ResChan <- fn()
		fmt.Println("spinnig a worker!")
	case chan string:
		ResChan := request.result.(chan string)
		fn := request.fn.(func() string)
		ResChan <- fn()
	}
}
func main() {
	fmt.Println("Starting")

	RequestChan := make(chan Request)

	loadBalancer := LoadBalancer{}
	go loadBalancer.start(RequestChan)

	//var req1 Request

	req1 := Request{
		fn: func() int {
			time.Sleep(time.Second * 4)
			return 2 * 4
		},
		result: make(chan int),
	}

	RequestChan <- req1

	switch req1.result.(type) {
	case chan int:
		ResChan := req1.result.(chan int)
		fmt.Println(<-ResChan)
	case chan string:
		ResChan := req1.result.(chan string)
		fmt.Println(<-ResChan)
	}

	req2 := Request{
		fn: func() string {
			return "hello"
		},
		result: make(chan string),
	}

	RequestChan <- req2

	switch req2.result.(type) {
	case chan int:
		ResChan := req2.result.(chan int)
		fmt.Println(<-ResChan)
	case chan string:
		ResChan := req2.result.(chan string)
		fmt.Println(<-ResChan)
	}

}
