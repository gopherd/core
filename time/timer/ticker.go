package timer

import "time"

type Ticker struct {
	interval time.Duration
	next     time.Time
}

func NewTicker(interval time.Duration) *Ticker {
	if interval <= 0 {
		panic("NewTicker: interval must be greater than zero")
	}
	return &Ticker{
		interval: interval,
		next:     time.Now().Add(interval),
	}
}

func (t *Ticker) Interval() time.Duration { return t.interval }

func (t *Ticker) Next(now time.Time) bool {
	diff := now.Sub(t.next)
	if diff < 0 {
		return false
	}
	t.next = t.next.Add(t.interval)
	return true
}
