package services

import (
	"sync"
	"time"
)

// ConfigEvent represents a mutation event broadcast to SSE subscribers.
type ConfigEvent struct {
	Type       string    `json:"type"`                 // route_created, route_updated, route_deleted, middleware_created, etc.
	Resource   string    `json:"resource"`             // route, middleware, instance
	ResourceID uint      `json:"resourceId,omitempty"` // ID of the affected resource
	InstanceID uint      `json:"instanceId,omitempty"` // associated instance ID
	Name       string    `json:"name,omitempty"`       // human-readable name of the resource
	UserID     string    `json:"userId,omitempty"`     // user who triggered the change
	Message    string    `json:"message"`              // human-readable description
	Timestamp  time.Time `json:"timestamp"`            // when the event occurred
}

// EventBus is a simple pub/sub hub for ConfigEvent.
type EventBus struct {
	subscribers map[chan ConfigEvent]struct{}
	mu          sync.Mutex
}

// NewEventBus creates a new EventBus.
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[chan ConfigEvent]struct{}),
	}
}

// Subscribe returns a buffered channel that will receive broadcast events.
func (eb *EventBus) Subscribe() chan ConfigEvent {
	ch := make(chan ConfigEvent, 32)
	eb.mu.Lock()
	eb.subscribers[ch] = struct{}{}
	eb.mu.Unlock()
	return ch
}

// Unsubscribe removes the channel from the subscriber set and closes it.
func (eb *EventBus) Unsubscribe(ch chan ConfigEvent) {
	eb.mu.Lock()
	delete(eb.subscribers, ch)
	eb.mu.Unlock()
	close(ch)
}

// Broadcast sends an event to every subscriber without blocking.
func (eb *EventBus) Broadcast(evt ConfigEvent) {
	if evt.Timestamp.IsZero() {
		evt.Timestamp = time.Now().UTC()
	}
	eb.mu.Lock()
	defer eb.mu.Unlock()
	for ch := range eb.subscribers {
		select {
		case ch <- evt:
		default:
			// drop if subscriber is slow
		}
	}
}
