package main

import (
	"fmt"
)

type Request struct {
	fn     func() int
	result chan int
}

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
	request.result <- request.fn()
}
func main() {
	fmt.Println("Starting")

	RequestChan := make(chan Request)

	loadBalancer := LoadBalancer{}
	go loadBalancer.start(RequestChan)

	r1 := Request{
		fn: func() int {
			return 2 * 4
		},
		result: make(chan int),
	}

	RequestChan <- r1
	fmt.Println(<-r1.result)

}
