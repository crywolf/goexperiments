package contextimpl

import (
	"errors"
	"reflect"
	"sync"
	"time"
)

// A Context carries a deadline, a cancelation signal, and other values across API boundaries.
// Context's methods may be called by multiple goroutines simultaneously.
type Context interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key interface{}) interface{}
}

type emptyCtx int

func (emptyCtx) Deadline() (deadline time.Time, ok bool) { return }
func (emptyCtx) Done() <-chan struct{}                   { return nil }
func (emptyCtx) Err() error                              { return nil }
func (emptyCtx) Value(key interface{}) interface{}       { return nil }

var (
	background = new(emptyCtx)
	todo       = new(emptyCtx)
)

// Background returns a non-nil, empty Context. It is never canceled, has no values, and has no deadline. It is typically used by the main function, initialization, and tests, and as the top-level Context for incoming requests.
func Background() Context { return background }

// TODO returns a non-nil, empty Context. Code should use context.TODO when it's unclear which Context to use or it is not yet available (because the surrounding function has not yet been extended to accept a Context parameter). TODO is recognized by static analysis tools that determine whether Contexts are propagated correctly in a program.
func TODO() Context { return todo }

type cancelCtx struct {
	Context
	done chan struct{}
	err  error
	sync.Mutex
}

func (ctx *cancelCtx) Done() <-chan struct{} { return ctx.done }
func (ctx *cancelCtx) Err() error {
	ctx.Lock()
	defer ctx.Unlock()
	return ctx.err
}

// ErrCanceled is the error returned by Context.Err when the context is canceled.
var ErrCanceled = errors.New("context canceled")

type deadlineExceededError struct{}

func (deadlineExceededError) Error() string {
	return "deadline exceeded"
}
func (deadlineExceededError) Timeout() bool {
	return true
}

// ErrDeadlineExceeded is the error returned by Context.Err when the context's deadline passes.
var ErrDeadlineExceeded error = deadlineExceededError{}

// A CancelFunc tells an operation to abandon its work. A CancelFunc does not wait for the work to stop. After the first call, subsequent calls to a CancelFunc do nothing.
type CancelFunc func()

// WithCancel returns a copy of parent with a new Done channel. The returned context's Done channel is closed when the returned cancel function is called or when the parent context's Done channel is closed, whichever happens first.
// Canceling this context releases resources associated with it, so code should call cancel as soon as the operations running in this Context complete.
func WithCancel(parent Context) (Context, CancelFunc) {
	ctx := &cancelCtx{
		Context: parent,
		done:    make(chan struct{}),
	}

	cancel := func() {
		ctx.cancel(ErrCanceled)
	}

	go func() {
		select {
		case <-parent.Done():
			ctx.cancel(parent.Err())
		case <-ctx.Done():
		}
	}()

	return ctx, cancel
}

func (ctx *cancelCtx) cancel(err error) {
	ctx.Lock()
	defer ctx.Unlock()
	if ctx.err != nil {
		return
	}
	ctx.err = err
	close(ctx.done)
}

type deadlineContext struct {
	*cancelCtx
	deadline time.Time
}

func (ctx *deadlineContext) Deadline() (deadline time.Time, ok bool) {
	return ctx.deadline, true
}

// WithDeadline returns a copy of the parent context with the deadline adjusted to be no later than d. If the parent's deadline is already earlier than d, WithDeadline(parent, d) is semantically equivalent to parent. The returned context's Done channel is closed when the deadline expires, when the returned cancel function is called, or when the parent context's Done channel is closed, whichever happens first.
// Canceling this context releases resources associated with it, so code should call cancel as soon as the operations running in this Context complete.
func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc) {
	cctx, cancel := WithCancel(parent)

	ctx := &deadlineContext{
		cancelCtx: cctx.(*cancelCtx),
		deadline:  deadline,
	}

	t := time.AfterFunc(time.Until(deadline), func() {
		ctx.cancel(ErrDeadlineExceeded)
	})

	stop := func() {
		t.Stop()
		cancel()
	}

	return ctx, stop
}

// WithTimeout returns WithDeadline(parent, time.Now().Add(timeout)).
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
	return WithDeadline(parent, time.Now().Add(timeout))
}

type valueCtx struct {
	Context
	key, value interface{}
}

func (ctx *valueCtx) Value(key interface{}) interface{} {
	if key == nil {
		panic("key is nil")
	}
	if !reflect.TypeOf(key).Comparable() {
		panic("key is not comparable")
	}
	if ctx.key == key {
		return ctx.value
	}
	return ctx.Context.Value(key)
}

// WithValue returns a copy of parent in which the value associated with key is val.
func WithValue(parent Context, key, val interface{}) Context {
	return &valueCtx{
		Context: parent,
		key:     key,
		value:   val,
	}
}
