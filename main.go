package main

import (
	"fmt"
)

type Request struct {
	fn     func() int
	result chan int
}

func loadBalancer(r *Request) chan int {
	go func() {
		r.result <- r.fn()
	}()
	return r.result
}
func main() {
	fmt.Println("hello")

	res := loadBalancer(&Request{
		fn: func() int {
			return 2 * 4
		},
		result: make(chan int),
	})

	fmt.Println(<-res)

}
