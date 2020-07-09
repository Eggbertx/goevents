package goevents

import (
	"fmt"
	"testing"
)

var (
	testOptions = &EventOptions{Once: true, AsGoroutine: false}
)

func TestListener(t *testing.T) {
	defer DisposeAll()
	var listener *Listener
	listener = AddListener("eventArgs", func(args ...interface{}) {
		t.Logf("fired %q event with args %v (only once: %t)\n", listener.Name, args, listener.Options.Once)
	}, testOptions)

	err := Emit("eventArgs", 1, 2, 3)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestListenerDispose(t *testing.T) {
	defer DisposeAll()
	err := Dispose("undefined1")
	if err == nil {
		t.Fatal(`Dispose("undefined1") should return an error because "undefined1" hasn't been created yet`)
	}
	AddListener("foo", func(args ...interface{}) {
		fmt.Println(`"foo" event emitted`)
	}, testOptions)
	if err = Dispose("foo"); err != nil {
		t.Fatal(`Dispose("foo") should return nil because "foo" was created and never disposed`)
	}
	if err = Dispose("foo"); err == nil {
		t.Fatal(`Dispose("foo") should return ErrEventNotFound because "foo" was already disposed`)
	}
}

func TestListenerOnce(t *testing.T) {
	defer DisposeAll()
	timesRan := 0
	AddListener("onceTest", func(...interface{}) {
		timesRan++
	}, testOptions)
	Emit("onceTest")
	Emit("onceTest") // this should be ignored
	Emit("onceTest") // this should be ignored
	if timesRan > 1 {
		t.Fatalf(`"onceTest" should only be running once (timesRan == %d)`, timesRan)
	}
}
