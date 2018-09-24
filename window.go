package slidingwindow

import (
	"sync"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/atomic"
)

// Window defines a sliding window for use to count averages over time.
type Window struct {
	sync.RWMutex
	window      time.Duration
	granularity time.Duration
	samples     []atomic.Int64
	pos         int
	stopping    chan struct{}
}

func New(window, granularity time.Duration) (*Window, error) {
	if window == 0 {
		return nil, errors.New("sliding window cannot be zero")
	}
	if granularity == 0 {
		return nil, errors.New("granularity cannot be zero")
	}
	if window <= granularity || window%granularity != 0 {
		return nil, errors.New("window size has to be a multiplier of granularity size")
	}

	sw := &Window{
		window:      window,
		granularity: granularity,
		samples:     make([]atomic.Int64, int(window/granularity)),
		stopping:    make(chan struct{}, 1),
	}
	go sw.shifter()

	return sw, nil
}

func MustNew(window, granularity time.Duration) *Window {
	w, err := New(window, granularity)
	if err != nil {
		panic(err)
	}
	return w
}

func (sw *Window) shifter() {
	ticker := time.NewTicker(sw.granularity)

	for {
		select {
		case <-ticker.C:
			if sw.pos = sw.pos + 1; sw.pos >= len(sw.samples) {
				sw.pos = 0
			}
			sw.samples[sw.pos].Swap(0)
		case <-sw.stopping:
			return
		}
	}
}

func (sw *Window) Add(v int64) {
	sw.samples[sw.pos].Add(v)
}

// Last retrieves the last N granularity samples and returns the total and number of samples
func (sw *Window) Last(n int) (total int64, samples int, err error) {
	if n <= 0 {
		return 0, 0, errors.New("cannot retrieve negative number of samples")
	}
	if n > len(sw.samples) {
		return 0, 0, errors.Errorf("cannot retrieve %d samples: only %d samples available", n, len(sw.samples))
	}

	var result int64
	if sw.pos >= (n - 1) {
		// (n - 1) >= sw.pos; in this case, we just count down and add
		lastIdx := sw.pos - (n - 1)
		for i := n - 1; i >= lastIdx; i-- {
			result += sw.samples[i].Load()
			if result != 0 {
				samples++
			}
		}
	} else {
		// We are somewhere in the middle; in this case, we subtract the index from position and then wrap around
		for i := 0; i < n; i++ {
			idx := sw.pos - i
			if idx < 0 {
				idx = len(sw.samples) - idx
			}
			result += sw.samples[i].Load()
			if result != 0 {
				samples++
			}
		}
	}
	return result, samples, nil
}
