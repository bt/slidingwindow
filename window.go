package slidingwindow

import (
	"sync"
	"time"

	"github.com/pkg/errors"
)

// slidingWindow defines a sliding window for use to count averages over time.
type slidingWindow struct {
	sync.RWMutex
	window      time.Duration
	granularity time.Duration
	samples     []int64
	pos         int
	stopping    chan struct{}
}

func New(window, granularity time.Duration) (*slidingWindow, error) {
	if window == 0 {
		return nil, errors.New("sliding window cannot be zero")
	}
	if granularity == 0 {
		return nil, errors.New("granularity cannot be zero")
	}
	if window <= granularity || window%granularity != 0 {
		return nil, errors.New("window size has to be a multiplier of granularity size")
	}

	sw := &slidingWindow{
		window:      window,
		granularity: granularity,
		samples:     make([]int64, int(window/granularity)),
		stopping:    make(chan struct{}, 1),
	}
	go sw.shifter()

	return sw, nil
}

func (sw *slidingWindow) shifter() {
	ticker := time.NewTicker(sw.granularity)

	for {
		select {
		case <-ticker.C:
			sw.Lock()
			if sw.pos = sw.pos + 1; sw.pos >= len(sw.samples) {
				sw.pos = 0
			}
			sw.samples[sw.pos] = 0
			sw.Unlock()
		case <-sw.stopping:
			return
		}
	}
}

func (sw *slidingWindow) Add(v int64) {
	sw.Lock()
	defer sw.Unlock()
	sw.samples[sw.pos] += v
}

// Last retrieves the last N granularity samples and returns the total
func (sw *slidingWindow) Last(n int) (int64, error) {
	if n <= 0 {
		return 0, errors.New("cannot retrieve negative number of samples")
	}
	if n > len(sw.samples) {
		return 0, errors.Errorf("cannot retrieve %d samples: only %d samples available", n, len(sw.samples))
	}

	var result int64
	if sw.pos >= (n - 1) {
		// (n - 1) >= sw.pos; in this case, we just count down and add
		lastIdx := sw.pos - (n - 1)
		for i := n - 1; i >= lastIdx; i-- {
			result += sw.samples[i]
		}
	} else {
		// We are somewhere in the middle; in this case, we subtract the index from position and then wrap around
		for i := 0; i < n; i++ {
			idx := sw.pos - i
			if idx < 0 {
				idx = len(sw.samples) - idx
			}
			result += sw.samples[i]
		}
	}
	return result, nil
}
