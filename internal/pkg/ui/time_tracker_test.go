package ui

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockClock struct {
	elapsed time.Duration
}

func (c mockClock) Since(time.Time) time.Duration {
	return c.elapsed
}

func TestTimerFunc(t *testing.T) {
	cases := []struct {
		name      string
		maxMsgLen int
		msg       string
		elapsed   time.Duration
		exp       string
	}{
		{
			name:      "microsecond scale",
			maxMsgLen: 1,
			msg:       "a",
			elapsed:   time.Nanosecond * 123500,
			exp:       "a took: 124µs\n",
		},
		{
			name:      "millisecond scale",
			maxMsgLen: 1,
			msg:       "a",
			elapsed:   time.Microsecond * 123500,
			exp:       "a took: 124ms\n",
		},
		{
			name:      "second scale",
			maxMsgLen: 1,
			msg:       "a",
			elapsed:   time.Millisecond * 123500,
			exp:       "a took: 2m4s\n",
		},
		{
			name:      "minute scale",
			maxMsgLen: 1,
			msg:       "a",
			elapsed:   time.Second * 123500,
			exp:       "a took: 34h18m20s\n",
		},
		{
			name:      "message same as truncate size",
			maxMsgLen: 10,
			msg:       "aaaaaaaaaa",
			elapsed:   0,
			exp:       "aaaaaaaaaa took: 0s\n",
		},
		{
			name:      "message over truncate size",
			maxMsgLen: 10,
			msg:       "aaaaaaaaaaa",
			elapsed:   0,
			exp:       "aaaaaaa... took: 0s\n",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			clock := mockClock{
				elapsed: c.elapsed,
			}

			buf := new(bytes.Buffer)
			tt := TimeTracker(buf, clock, c.maxMsgLen)

			tt(time.Now(), c.msg)

			assert.Equal(t, c.exp, buf.String())
		})
	}
}
