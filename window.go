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
	size        int
	stopping    chan struct{}
}

func newWindow(window, granularity time.Duration) (*Window, error) {
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

	return sw, nil
}

func New(window, granularity time.Duration) (*Window, error) {
	w, err := newWindow(window, granularity)
	if err != nil {
		return nil, err
	}
	go w.shifter()
	return w, nil
}

func MustNew(window, granularity time.Duration) *Window {
	w, err := New(window, granularity)
	if err != nil {
		panic(err)
	}
	return w
}

func MustNewFromSamples(window, granularity time.Duration, samples []int64) *Window {
	w, err := newWindow(window, granularity)
	if err != nil {
		panic(err)
	}

	fillSamples(w, samples...)
	go w.shifter()
	return w
}

func (sw *Window) shifter() {
	ticker := time.NewTicker(sw.granularity)

	for {
		select {
		case <-ticker.C:
			sw.nextPosition()
		case <-sw.stopping:
			return
		}
	}
}

func (sw *Window) nextPosition() {
	if sw.pos = sw.pos + 1; sw.pos >= len(sw.samples) {
		sw.pos = 0
	}
	sw.samples[sw.pos].Swap(0)
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

	// if position - (n - 1) is higher than or equal to zero, then
	lastIdx := sw.pos - (n - 1)
	if lastIdx >= 0 {
		// We have enough samples to process this request, therefore we iterate till the last index
		for i := sw.pos; i >= lastIdx; i-- {
			val := sw.samples[i].Load()
			if val != 0 {
				result += val
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
			val := sw.samples[i].Load()
			if val != 0 {
				result += val
				samples++
			}
		}
	}
	return result, samples, nil
}

func fillSamples(win *Window, i ...int64) {
	for j, v := range i {
		win.Add(v)
		if j != len(i)-1 {
			win.nextPosition()
		}
	}
}
