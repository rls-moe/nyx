package snowflakes

import (
	"errors"
	"sync"
	"time"
)

const (
	counterLen  = 10
	counterMask = -1 ^ (-1 << counterLen)
)

var (
	errNoFuture = errors.New("Start Time cannot be set in the future")
)

// Generator is a fountain for new snowflakes. StartTime must be
// initialized to a past point in time and Instance ID can be any
// positive value or 0.
//
// If any value is not correctly set, new IDs cannot be produced.
type Generator struct {
	StartTime int64
	mutex     *sync.Mutex
	sequence  int32
	now       int64
}

// NewID generates a new, unique snowflake value
//
// Up to 8192 snowflakes per second can be requested
// If exhausted, it blocks and sleeps until a new second
// of unix time starts.
//
// The return value is signed but always positive.
//
// Additionally, the return value is monotonic for a single
// instance and weakly monotonic for many instances.
func (g *Generator) NewID() (int64, error) {
	if g.mutex == nil {
		g.mutex = new(sync.Mutex)
	}
	if g.StartTime > time.Now().Unix() {
		return 0, errNoFuture
	}
	g.mutex.Lock()
	defer g.mutex.Unlock()

	var (
		now   int64
		flake int64
	)
	now = int64(time.Now().Unix())

	if now == g.now {
		g.sequence = (g.sequence + 1) & counterMask
		if g.sequence == 0 {
			for now <= g.now {
				now = int64(time.Now().Unix())
				time.Sleep(time.Microsecond * 100)
			}
		}
	} else {
		g.sequence = 0
	}

	g.now = now

	flake = int64(
		((now - g.StartTime) << counterLen) |
			int64(g.sequence))

	return flake, nil
}
