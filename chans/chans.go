package chans

import (
	"cmp"
	"context"
	"sync"
	"time"
)

var nopCancel = func() {}

// MakeChan make a channel with type T, The channel's buffer is initialized with the specified buffer capacity. If zero, or the size is omitted, the channel is unbuffered.
func MakeChan[T any](size ...int) (ch chan T, closech func()) {
	if len(size) > 0 {
		ch = make(chan T, size[0])
	} else {
		ch = make(chan T)
	}
	return ch, sync.OnceFunc(func() { close(ch) })
}

// StructChan make a unbuffered struct{} channel
func StructChan() (ch chan struct{}, closech func()) { return MakeChan[struct{}]() }

// AfterChan  will invoke fn when ch done, but cancel when ctx done or invoke cancel function
func AfterChan[T any](ctx context.Context, ch <-chan T, fn func(), cancelable ...bool) (cancel func()) {
	if ctx == nil {
		ctx = context.Background()
	}
	if cmp.Or(cancelable...) {
		ctx, cancel = context.WithCancel(ctx)
	} else {
		cancel = nopCancel
	}

	go func() {
		defer cancel()
		select {
		case <-ctx.Done():
			return
		case <-ch:
			fn()
		}
	}()

	return
}

// AfterContext  will invoke fn when ctx done, but cancel when invoke cancel function
func AfterContext(ctx context.Context, fn func(), cancelable ...bool) (cancel func()) {
	if cmp.Or(cancelable...) {
		done, closech := StructChan()
		go func() {
			defer closech()
			select {
			case <-done:
				return
			case <-ctx.Done():
				fn()
			}
		}()
		cancel = closech
	} else {
		go func() {
			<-ctx.Done()
			fn()
		}()
		cancel = nopCancel
	}

	return
}

// AfterTime will invoke fn after d, but cancel when ctx done or invoke cancel function
func AfterTime(ctx context.Context, d time.Duration, fn func(), cancelable ...bool) (cancel func()) {
	if ctx == nil {
		ctx = context.Background()
	}
	if cmp.Or(cancelable...) {
		ctx, cancel = context.WithCancel(ctx)
	} else {
		cancel = nopCancel
	}

	go func() {
		select {
		case <-ctx.Done():
			return
		case <-time.After(d):
			fn()
		}
	}()

	return
}

// Sleep sleeps for the specified duration. It returns false if the context is canceled.
func Sleep(ctx context.Context, d time.Duration) (r bool) {
	if ctx == nil {
		ctx = context.Background()
	}

	select {
	case <-ctx.Done():
		r = false
	case <-time.After(d):
		r = true
	}

	return
}
