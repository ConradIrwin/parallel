// Package parallel is an implementation of structured concurrency for go.
// https://vorpus.org/blog/notes-on-structured-concurrency-or-go-statement-considered-harmful/
//
// It is designed to help reason about parallel code by ensuring that
// go-routines are started and stopped in a strictly nested pattern: a child
// goroutine will never outlive its parent.
package parallel

import (
	"sync"
	"sync/atomic"
)

// P represents the parallel execution of a set of goroutines.
type P struct {
	// OnPanic is called when a goroutine panics. You should return
	// false from this if you don't wish the panic to propagate.
	// This callback must be safe to call from multiple goroutines.
	OnPanic func(p any) bool

	panicked  atomic.Value
	finished  atomic.Bool
	waitgroup sync.WaitGroup
}

// Do starts a new parallel execution context and runs it to completion.
//
// After f has run, and after any goroutines started by p.Go have finished,
// Do will mark p as finished and then return. Any further calls to p.Go will panic.
//
// In practice this means that you can only safely call p.Go from within f,
// or within the goroutines started by p.Go.
//
// If f, or any goroutine started by p.Go panics, then Do will panic.
//
// The panic behaviour can be overwritten by setting p.OnPanic
// from within the callback passed to .Do() before any calls to .Go().
func Do(f func(p *P)) {
	p := &P{OnPanic: func(any) bool { return true }}
	defer p.wait()
	defer p.recover()
	f(p)
}

// Go starts a new goroutine. If p is already marked as finished, Go will panic.
func (p *P) Go(f func()) {
	if p.finished.Load() {
		panic("parallel: cannot call Go after Do has returned")
	}
	p.waitgroup.Add(1)
	go func() {
		defer p.waitgroup.Done()
		defer p.recover()

		f()
	}()
}

func (p *P) recover() {
	if r := recover(); r != nil && p.OnPanic(r) {
		if p.panicked.CompareAndSwap(nil, r) {
			return
		}
	}
}

func (p *P) wait() {
	p.waitgroup.Wait()
	// note: race here if someone calls p.Go during this comment.
	// that's ok – they'll likely get a panic eventually to flag the API misuse.
	p.finished.Store(true)
	if r := p.panicked.Load(); r != nil {
		panic(p.panicked.Load())
	}
}

// Each runs the callback for each item in the slice in parallel.
func Each[T any](items []T, f func(T)) {
	Do(func(p *P) {
		for _, v := range items {
			v2 := v
			p.Go(func() { f(v2) })
		}
	})
}
