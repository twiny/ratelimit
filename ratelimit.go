package ratelimit

import (
	"strconv"
	"sync"
	"time"

	"github.com/benbjohnson/clock"
)

// Limiter
type Limiter struct {
	mu         *sync.Mutex
	rate       int
	interval   time.Duration
	sleepFor   time.Duration
	last       time.Time
	perRequest time.Duration
	maxSlack   time.Duration
	clock      Clock
}

// Clock
type Clock interface {
	Now() time.Time
	Sleep(time.Duration)
}

// config limiter.
type config struct {
	maxSlack time.Duration
	clock    Clock
}

// RateLimiter New returns a Limiter that will limit to the given RPS.
func NewLimiter(rate int, interval time.Duration) *Limiter {
	config := config{
		maxSlack: 10,
		clock:    clock.New(),
	}

	l := &Limiter{
		mu:       &sync.Mutex{},
		rate:     rate,
		interval: interval,
		//
		perRequest: interval / time.Duration(rate),
		maxSlack:   -1 * config.maxSlack * interval / time.Duration(rate),
		clock:      config.clock,
	}

	return l
}

// Take blocks to ensure that the time spent between multiple
// Take calls is on average time.Duration/rate.
func (l *Limiter) Take() time.Time {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.clock.Now()

	// If this is our first request, then we allow it.
	if l.last.IsZero() {
		l.last = now
		return l.last
	}

	// sleepFor calculates how much time we should sleep based on
	// the perRequest budget and how long the last request took.
	l.sleepFor += l.perRequest - now.Sub(l.last)

	// We shouldn't allow sleepFor to get too negative.
	if l.sleepFor < l.maxSlack {
		l.sleepFor = l.maxSlack
	}

	// If sleepFor is positive, then we should sleep now.
	switch {
	case l.sleepFor > 0:
		l.clock.Sleep(l.sleepFor)
		l.last = now.Add(l.sleepFor)
		l.sleepFor = 0
	default:
		l.last = now
	}

	return l.last
}

// Rate
func (l *Limiter) Rate() int {
	return l.rate
}

// Duration
func (l *Limiter) Duration() time.Duration {
	return l.interval
}

// String
func (l *Limiter) String() string {
	return strconv.Itoa(l.rate) + "/" + l.Duration().String()
}
