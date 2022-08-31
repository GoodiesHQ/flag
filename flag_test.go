package flag

import (
	"testing"
	"time"
)

func TestNewFlag(t *testing.T) {
	f := NewFlag()
	if f.IsSet() != false {
		t.Fatal("flag initialization failure")
	}
}

func TestSetFlag(t *testing.T) {
	f := NewFlag()
	f.Set()
	if f.IsSet() != true {
		t.Fatal("flag set failure")
	}
}

func TestClearFlag(t *testing.T) {
	f := NewFlag()
	f.Set()
	f.Clear()
	if f.IsSet() != false {
		t.Fatal("flag clear failure")
	}
}

func TestChannelClosure(t *testing.T) {
	f := NewFlag()
	select {
	case _, ok := <-f.Channel():
		if ok {
			t.Fatal("channel should be empty")
		}
		t.Fatal("channel should be open")
	default:
		break
	}

	f.Set()
	select {
	case _, ok := <-f.Channel():
		if ok {
			t.Fatal("channel should be closed")
		}
		break
	default:
		t.Fatal("channel should be closed")
	}
}

func TestWait(t *testing.T) {
	f := NewFlag()
	if f.Wait(1*time.Second) != true {
		t.Fatal("wait should timeout")
	}
	go func() {
		time.Sleep(500 * time.Millisecond)
		f.Set()
	}()
	if f.Wait(1*time.Second) != false {
		t.Fatal("wait should be quick")
	}

	f.Clear()

	go func() {
		time.Sleep(500 * time.Millisecond)
		f.Set()
	}()
	if f.Wait(1*time.Second) != false {
		t.Fatal("wait should be quick")
	}
}
