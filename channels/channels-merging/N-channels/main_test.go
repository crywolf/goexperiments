package main

import (
	"fmt"
	"testing"
)

func TestMerge(t *testing.T) {
	c := merge(asChan(1, 2, 3), asChan(4, 5, 6), asChan(7, 8, 9))

	seen := make(map[int]bool)

	for v := range c {
		if seen[v] {
			t.Errorf("saw %d at least twice", v)
		}
		seen[v] = true
	}

	for i := 1; i <= 9; i++ {
		if !seen[i] {
			t.Errorf("did not see %d", i)
		}
	}
}

func BenchmarkMerge(b *testing.B) {
	for n := 1; n <= 1024; n *= 2 {
		chans := make([]<-chan int, n)

		b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {
			b.StopTimer()
			for i := 0; i < b.N; i++ {
				for j := range chans {
					chans[j] = asChan(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
				}
				b.StartTimer()
				c := merge(chans...)
				for range c {
				}
			}
		})

	}
}
