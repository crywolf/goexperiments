package main

import "sync"

type Idler struct {
	mu    sync.Mutex
	idle  sync.Cond
	busy  bool
	idles int64
}

func (i *Idler) AwaitIdle() {
	i.mu.Lock()
	defer i.mu.Unlock()
	idles := i.idles
	for i.busy && idles == i.idles {
		i.idle.Wait()
	}
}

func (i *Idler) SetBusy(b bool) {
	i.mu.Lock()
	defer i.mu.Unlock()
	wasBusy := i.busy
	i.busy = b
	if wasBusy && !i.busy {
		i.idles++
		i.idle.Broadcast()
	}
}

func NewIdler() *Idler {
	i := &Idler{}
	i.idle.L = &i.mu
	return i
}
