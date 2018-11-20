// Semaphore channel: Limiting concurrency (limits work in flight)
// A semaphore channel takes place of the WaitGroup
package main

import "fmt"

type Task struct{}

func main() {
	var hugeSlice []Task
	type token struct{}

	limit := 8

	// start the work
	sem := make(chan token, limit)
	for _, task := range hugeSlice {
		sem <- token{}
		go func(task Task) {
			perform(task)
			<-sem
		}(task)
	}

	// wait for completion
	for i := 0; i < limit; i++ {
		sem <- token{}
	}
}

func perform(task Task) {
	fmt.Printf("%v\n", task)
}
