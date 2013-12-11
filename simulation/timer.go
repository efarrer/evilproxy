package simulation

import (
	"time"
)

/*
 * A simple timer for seeing how long an operation lasts
 */
type Timer time.Time

/*
 * Starts a timer
 */
func StartTimer() Timer {
	return Timer(time.Now())
}

/*
 * Returns the number of milliseconds that have elapsed since the timer started
 */
func (t Timer) ElapsedMilliseconds() time.Duration {
	return time.Now().Sub(time.Time(t)) * time.Millisecond
}

/*
 * Utility function for fuzzy comparisons of durations
 */
func FuzzyEquals(a, b, delta time.Duration) bool {
	diff := a - b
	if diff < 0 {
		diff *= -1
	}
	return diff <= delta
}
