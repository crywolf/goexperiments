package main

import (
	"fmt"
	"sync"
)

func main() {
	a := asChan(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
	b := asChan(10, 11, 12, 13, 14, 15, 16, 17, 18, 19)
	c := asChan(20, 21, 22, 23, 24, 25, 26, 27, 28, 29)

	for v := range merge(a, b, c) {
		fmt.Println(v)
	}
}

func merge(chans ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	wg.Add(len(chans))

	go func() {
		for _, c := range chans {
			go func(c <-chan int) {
				for v := range c {
					out <- v
				}
				wg.Done()
			}(c)
		}

		wg.Wait()
		close(out)
	}()

	return out
}

func asChan(n ...int) <-chan int {
	c := make(chan int)

	go func() {
		for _, v := range n {
			c <- v
		}
		close(c)
	}()

	return c
}
