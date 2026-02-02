package testutil

import "time"

// Clock is an interface for time operations to enable testing.
type Clock interface {
	Now() time.Time
}

// RealClock implements Clock using the real system time.
type RealClock struct{}

// Now returns the current time.
func (RealClock) Now() time.Time {
	return time.Now()
}

// FixedClock implements Clock with a fixed time for testing.
type FixedClock struct {
	FixedTime time.Time
}

// Now returns the fixed time.
func (c FixedClock) Now() time.Time {
	return c.FixedTime
}

// NewFixedClock creates a new FixedClock with the given time.
func NewFixedClock(t time.Time) *FixedClock {
	return &FixedClock{FixedTime: t}
}

// NewFixedClockAt creates a new FixedClock at a specific date/time.
// This is a convenience function for tests.
func NewFixedClockAt(year, month, day, hour, min, sec int) *FixedClock {
	return &FixedClock{
		FixedTime: time.Date(year, time.Month(month), day, hour, min, sec, 0, time.UTC),
	}
}

// Advance moves the fixed clock forward by the given duration.
func (c *FixedClock) Advance(d time.Duration) {
	c.FixedTime = c.FixedTime.Add(d)
}

// Set sets the fixed clock to a specific time.
func (c *FixedClock) Set(t time.Time) {
	c.FixedTime = t
}
