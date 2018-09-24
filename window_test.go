package slidingwindow

import (
	"go.uber.org/atomic"
	"testing"
	"time"

	"gotest.tools/assert"
)

func fill(win *Window, i ...int64) {
	for _, v := range i {
		var atom atomic.Int64
		atom.Store(v)
		win.samples = append(win.samples, atom)
	}
}

func TestWindow_Simple(t *testing.T) {
	win := Window{
		window:      time.Second * 5,
		granularity: time.Second,
		samples:     []atomic.Int64{},
		pos:         4,
	}
	fill(&win, 1, 2, 3, 4, 5)

	i, samples, err := win.Last(5)
	assert.NilError(t, err)
	assert.Equal(t, i, int64(15))
	assert.Equal(t, samples, int(5))
}

func TestWindow_Wrapping(t *testing.T) {
	win := Window{
		window:      time.Second * 5,
		granularity: time.Second,
		samples:     []atomic.Int64{},
		pos:         0,
	}
	fill(&win, 1, 2, 3, 4, 5)

	i, samples, err := win.Last(5)
	assert.NilError(t, err)
	assert.Equal(t, i, int64(15))
	assert.Equal(t, samples, int(5))
}

func TestWindow_Zero(t *testing.T) {
	win := Window{
		window:      time.Second * 5,
		granularity: time.Second,
		samples:     []atomic.Int64{},
		pos:         0,
	}
	fill(&win, 7, 0, 0, 0, 0)

	i, samples, err := win.Last(1)
	assert.NilError(t, err)
	assert.Equal(t, i, int64(7))
	assert.Equal(t, samples, int(1))
}
