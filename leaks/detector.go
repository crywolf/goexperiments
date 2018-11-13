package main

import (
	"bytes"
	"fmt"
	"runtime"
	"time"
)

func leakDetector() {
	buf := make([]byte, 2048)
	for {
		n := runtime.Stack(buf, true)
		//splits := strings.Split(string(buf[:n]), "\n\n")
		splits := bytes.Split(buf[:n], []byte("\n\n"))
		for i, split := range splits {
			fmt.Printf("#%d:\n%s\n", i, split)
		}
		<-time.After(3 * time.Second)
		// Reuse buffer for next refresh
		copy(buf[:n], bytes.Repeat([]byte{0x00}, n))
	}
}

func doWork() {
	for {
		<-time.After(1 * time.Second)
		_ = 1 + 1
	}
}

func main() {
	for i := 0; i < 5; i++ {
		go doWork()
	}
	leakDetector()
}
