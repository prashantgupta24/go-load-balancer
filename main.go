package main

import (
	"fmt"
)

type Request struct {
	fn     func() int
	result chan int
}

func loadBalancer(work chan Request) chan int {
	// go func() {
	// 	r.result <- r.fn()
	// }()
	// return r.result
	for {
		select {
		case w := <-work:
			{
				w.result <- w.fn()
			}
		}
	}
}

func main() {
	fmt.Println("hello")

	WorkChan := make(chan Request)
	go loadBalancer(WorkChan)
	//ResultChan := make(chan int)

	r1 := Request{
		fn: func() int {
			return 2 * 4
		},
		result: make(chan int),
	}

	WorkChan <- r1
	fmt.Println(<-r1.result)

}
