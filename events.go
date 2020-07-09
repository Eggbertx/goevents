package goevents

import (
	"errors"
)

var (
	// ErrEventNotFound is returned when you call Emit or Dispose with a name that isn't registered
	ErrEventNotFound = errors.New("Requested event does not exist")
	// ErrEventExists is returned when you call AddEvent using an event name that is already registered
	// ErrEventExists = errors.New("Event already exists")
	defaultEmitter = &EventEmitter{}
	defaultOptions = &EventOptions{Once: false, Overwrite: false, AsGoroutine: false}
)

// EventOptions is used for passing options about how the event should be treated.
// AsGoroutine is currently not used
type EventOptions struct {
	Once        bool
	Overwrite   bool
	AsGoroutine bool
}

// Listener is created by AddListener, and sholdn't be created manually
type Listener struct {
	Name       string
	cb         func(...interface{})
	Options    *EventOptions
	alreadyRan bool
}

// EventEmitter ...
type EventEmitter map[string]*Listener

// AddListener creates a new event with the specified name and callback function. If once is true, Emit() will only call it the first time
func (emitter EventEmitter) AddListener(name string, cb func(...interface{}), options *EventOptions) *Listener {
	if options == nil {
		options = defaultOptions
	}
	emitter[name] = &Listener{
		Name:       name,
		cb:         cb,
		Options:    options,
		alreadyRan: false,
	}
	return emitter[name]
}

// Dispose deletes the event in the event map and returns ErrEventNotFound if it doesn't exist
func (emitter EventEmitter) Dispose(name string) error {
	if emitter[name] == nil {
		return ErrEventNotFound
	}
	delete(emitter, name)
	return nil
}

// DisposeAll clears the event map
func (emitter EventEmitter) DisposeAll() {
	for e := range emitter {
		delete(emitter, e)
	}
}

// Emit runs the callback function of the event with the specified name and returns ErrEventNotFound if it doesn't exist
func (emitter EventEmitter) Emit(name string, args ...interface{}) error {
	if emitter[name] == nil {
		return ErrEventNotFound
	}
	if emitter[name].Options.Once && emitter[name].alreadyRan {
		return nil
	}
	emitter[name].cb(args...)
	emitter[name].alreadyRan = true
	return nil
}

// AddListener creates a new event with the specified name and callback function. If once is true, Emit() will only call it the first time
func AddListener(name string, cb func(...interface{}), options *EventOptions) *Listener {
	return defaultEmitter.AddListener(name, cb, options)
}

// Dispose deletes the event in the event map and returns ErrEventNotFound if it doesn't exist
func Dispose(name string) error {
	return defaultEmitter.Dispose(name)
}

// DisposeAll clears the event map
func DisposeAll() {
	defaultEmitter.DisposeAll()
}

// Emit runs the callback function of the event with the specified name and returns ErrEventNotFound if it doesn't exist
func Emit(name string, args ...interface{}) error {
	return defaultEmitter.Emit(name, args...)
}
