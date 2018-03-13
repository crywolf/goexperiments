package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	//go func() {
	//	time.Sleep(time.Second)
	//	cancel()
	//}()

	time.AfterFunc(time.Second, cancel) // the same as commented out code above

	sleepAndTalk(ctx, 5*time.Second, "Ahoj!")
}

func sleepAndTalk(ctx context.Context, d time.Duration, msg string) {

	select {
	case <-time.After(d):
		fmt.Println(msg)
	case <-ctx.Done():
		log.Println(ctx.Err())
	}

}
