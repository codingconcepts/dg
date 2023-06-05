package random

import "time"

var (
	r = newSplitMix64(time.Now().UnixNano())
)

type splitMix64 struct {
	s uint64
}

func newSplitMix64(seed int64) *splitMix64 {
	return &splitMix64{
		s: uint64(seed),
	}
}

// Intn returns a non-negative pseudo-random int.
func Intn(n int) int {
	return int(r.uint64()&(1<<63-1)) % n

}

func (x *splitMix64) uint64() uint64 {
	x.s = x.s + uint64(0x9E3779B97F4A7C15)
	z := x.s
	z = (z ^ (z >> 30)) * uint64(0xBF58476D1CE4E5B9)
	z = (z ^ (z >> 27)) * uint64(0x94D049BB133111EB)
	return z ^ (z >> 31)
}
