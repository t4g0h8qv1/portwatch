// Package portevents implements a lightweight publish/subscribe event bus
// for port state changes within portwatch.
//
// Subscribers register a Handler via Subscribe and will receive an Event
// each time Publish or PublishAll is called. The bus is safe for concurrent
// use from multiple goroutines.
//
// Typical usage:
//
//	bus := portevents.New()
//	bus.Subscribe(func(e portevents.Event) {
//		fmt.Printf("port %d %s on %s\n", e.Port, e.Type, e.Host)
//	})
//	bus.PublishAll("192.168.1.1", newPorts, portevents.EventOpened)
package portevents
