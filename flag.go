package flag

import (
	"sync/atomic"
	"time"
)

type Flag struct {
	value atomic.Bool
	ch    chan struct{}
}

func NewFlag() *Flag {
	return &Flag{ch: make(chan struct{})}
}

func (f *Flag) IsSet() bool {
	return f.value.Load()
}

func (f *Flag) Set() {
	if f.value.CompareAndSwap(false, true) {
		close(f.ch)
	}
}

func (f *Flag) Clear() {
	if f.value.CompareAndSwap(true, false) {
		f.ch = make(chan struct{})
	}
}

// Returns the channel for the flag
func (f *Flag) Channel() <-chan struct{} {
	return f.ch
}

// Waits for the flag to be set, or duration.
// Returns true if a timeout occurred
// Returns false if the channel was closed
func (f *Flag) Wait(timeout time.Duration) bool {
	select {
	case <-time.After(timeout):
		return true
	case <-f.Channel():
		return false
	}
}
