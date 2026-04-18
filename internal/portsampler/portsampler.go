// Package portsampler provides periodic port sampling with configurable
// jitter to avoid thundering-herd effects when many hosts are monitored.
package portsampler

import (
	"context"
	"math/rand"
	"sync"
	"time"
)

// Sample holds the result of a single port sampling run.
type Sample struct {
	Host      string
	Ports     []int
	SampledAt time.Time
	Err       error
}

// ScanFunc is the function called to obtain open ports for a host.
type ScanFunc func(ctx context.Context, host string) ([]int, error)

// Sampler periodically invokes a ScanFunc and sends results on a channel.
type Sampler struct {
	host     string
	interval time.Duration
	jitter   time.Duration
	scan     ScanFunc
	mu       sync.Mutex
	last     *Sample
}

// New creates a Sampler for the given host.
// interval is the base period between samples; jitter adds a random
// offset in [0, jitter) to each interval to spread load.
func New(host string, interval, jitter time.Duration, fn ScanFunc) (*Sampler, error) {
	if interval <= 0 {
		return nil, ErrInvalidInterval
	}
	if jitter < 0 {
		return nil, ErrInvalidJitter
	}
	if fn == nil {
		return nil, ErrNilScanFunc
	}
	return &Sampler{host: host, interval: interval, jitter: jitter, scan: fn}, nil
}

// Last returns the most recent sample, or nil if none has been taken.
func (s *Sampler) Last() *Sample {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.last
}

// Run starts the sampler and sends each Sample on the returned channel.
// It stops when ctx is cancelled, then closes the channel.
func (s *Sampler) Run(ctx context.Context) <-chan Sample {
	ch := make(chan Sample, 1)
	go func() {
		defer close(ch)
		for {
			sample := s.collect(ctx)
			s.mu.Lock()
			s.last = &sample
			s.mu.Unlock()
			select {
			case ch <- sample:
			case <-ctx.Done():
				return
			}
			wait := s.interval
			if s.jitter > 0 {
				wait += time.Duration(rand.Int63n(int64(s.jitter)))
			}
			select {
			case <-time.After(wait):
			case <-ctx.Done():
				return
			}
		}
	}()
	return ch
}

func (s *Sampler) collect(ctx context.Context) Sample {
	ports, err := s.scan(ctx, s.host)
	return Sample{Host: s.host, Ports: ports, SampledAt: time.Now(), Err: err}
}
