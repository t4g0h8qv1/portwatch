// Package portevents provides an event bus for port state change notifications.
package portevents

import "sync"

// EventType describes the kind of port event.
type EventType string

const (
	EventOpened EventType = "opened"
	EventClosed EventType = "closed"
)

// Event represents a single port state change.
type Event struct {
	Host  string
	Port  int
	Type  EventType
}

// Handler is a function that receives a port event.
type Handler func(e Event)

// Bus dispatches port events to registered handlers.
type Bus struct {
	mu       sync.RWMutex
	handlers []Handler
}

// New returns a new Bus.
func New() *Bus {
	return &Bus{}
}

// Subscribe registers a handler to receive all future events.
func (b *Bus) Subscribe(h Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers = append(b.handlers, h)
}

// Publish sends an event to all registered handlers.
func (b *Bus) Publish(e Event) {
	b.mu.RLock()
	handlers := make([]Handler, len(b.handlers))
	copy(handlers, b.handlers)
	b.mu.RUnlock()
	for _, h := range handlers {
		h(e)
	}
}

// PublishAll sends one event per port in the provided list.
func (b *Bus) PublishAll(host string, ports []int, t EventType) {
	for _, p := range ports {
		b.Publish(Event{Host: host, Port: p, Type: t})
	}
}
