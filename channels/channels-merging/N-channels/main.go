package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	a := asChan(0, 1, 2, 3, 4)
	b := asChan(5, 6, 7, 8, 9)

	c := merge(a, b)
	for v := range c {
		fmt.Println(v)
	}
}

func merge(a, b <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for a != nil || b != nil {
			select {
			case v, ok := <-a:
				if !ok {
					a = nil
					fmt.Println("a is done")
					continue
				}
				out <- v
			case v, ok := <-b:
				if !ok {
					b = nil
					fmt.Println("b is done")
					continue
				}
				out <- v
			}
		}
	}()

	return out
}

func asChan(n ...int) <-chan int {
	c := make(chan int)
	rand.Seed(time.Now().UnixNano())

	go func() {
		for _, v := range n {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
		close(c)
	}()

	return c
}
