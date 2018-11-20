// Share completion by completing communication

// Allocate the channel to be closed when the event starts, or when the first waiter appear.
package main

import (
	"context"
)

type Idler struct {
	next chan chan struct{}
}

func (i *Idler) AwaitIdle(ctx context.Context) error {
	idle := <-i.next
	i.next <- idle
	if idle != nil {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-idle:
		}
	}
	return nil
}

func (i *Idler) SetBusy(b bool) {
	idle := <-i.next
	if b && (idle == nil) { // idle to busy transition
		idle = make(chan struct{})
	} else if !b && (idle != nil) { // busy to idle transition
		close(idle) // idle now
		idle = nil
	}
}

func NewIdler() *Idler {
	next := make(chan chan struct{}, 1)
	next <- nil
	return &Idler{next}
}
