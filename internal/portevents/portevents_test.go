package portevents_test

import (
	"sync"
	"testing"

	"github.com/example/portwatch/internal/portevents"
)

func TestSubscribe_And_Publish(t *testing.T) {
	bus := portevents.New()
	var got []portevents.Event
	bus.Subscribe(func(e portevents.Event) {
		got = append(got, e)
	})
	bus.Publish(portevents.Event{Host: "localhost", Port: 80, Type: portevents.EventOpened})
	if len(got) != 1 {
		t.Fatalf("expected 1 event, got %d", len(got))
	}
	if got[0].Port != 80 || got[0].Type != portevents.EventOpened {
		t.Errorf("unexpected event: %+v", got[0])
	}
}

func TestPublishAll(t *testing.T) {
	bus := portevents.New()
	var got []portevents.Event
	bus.Subscribe(func(e portevents.Event) { got = append(got, e) })
	bus.PublishAll("host1", []int{22, 443, 8080}, portevents.EventClosed)
	if len(got) != 3 {
		t.Fatalf("expected 3 events, got %d", len(got))
	}
	for _, e := range got {
		if e.Type != portevents.EventClosed {
			t.Errorf("expected EventClosed, got %s", e.Type)
		}
		if e.Host != "host1" {
			t.Errorf("expected host1, got %s", e.Host)
		}
	}
}

func TestMultipleSubscribers(t *testing.T) {
	bus := portevents.New()
	var mu sync.Mutex
	counts := make([]int, 3)
	for i := 0; i < 3; i++ {
		i := i
		bus.Subscribe(func(e portevents.Event) {
			mu.Lock()
			counts[i]++
			mu.Unlock()
		})
	}
	bus.Publish(portevents.Event{Host: "h", Port: 9, Type: portevents.EventOpened})
	for i, c := range counts {
		if c != 1 {
			t.Errorf("subscriber %d: expected 1 call, got %d", i, c)
		}
	}
}

func TestPublish_NoSubscribers(t *testing.T) {
	bus := portevents.New()
	// should not panic
	bus.Publish(portevents.Event{Host: "x", Port: 1, Type: portevents.EventOpened})
}
