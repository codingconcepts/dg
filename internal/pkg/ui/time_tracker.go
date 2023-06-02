package ui

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// Clock allows time.Now to be mocked in tests.
type Clock interface {
	Since(time.Time) time.Duration
}

// TimerFunc is a function whose execution can be deferred in order to
// time an event.
type TimerFunc func(time.Time, string)

// TimeTracker returns a new TimerFunc, given a maxMessageLen, which
// determines at which point messages should start to get truncated.
func TimeTracker(w io.Writer, c Clock, maxMessageLen int) TimerFunc {
	return func(start time.Time, msg string) {
		elapsed := c.Since(start)
		switch {
		case elapsed > time.Second:
			elapsed = elapsed.Round(time.Second)
		case elapsed > time.Millisecond:
			elapsed = elapsed.Round(time.Millisecond)
		default:
			elapsed = elapsed.Round(time.Microsecond)
		}

		if len(msg) > maxMessageLen {
			msg = msg[:maxMessageLen-3] + "..."
		}

		padding := strings.Repeat(" ", maxMessageLen-len(msg))
		fmt.Fprintf(w, "%s %stook: %s\n", msg, padding, elapsed)
	}
}
