package slidingwindow

import (
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestWindow_Simple(t *testing.T) {
	win := Window{
		window:      time.Second * 5,
		granularity: time.Second,
		samples:     []int64{1, 2, 3, 4, 5},
		pos:         4,
	}

	i, samples, err := win.Last(5)
	assert.NilError(t, err)
	assert.Equal(t, i, int64(15))
	assert.Equal(t, samples, int(5))
}

func TestWindow_Wrapping(t *testing.T) {
	win := Window{
		window:      time.Second * 5,
		granularity: time.Second,
		samples:     []int64{1, 2, 3, 4, 5},
		pos:         0,
	}

	i, samples, err := win.Last(5)
	assert.NilError(t, err)
	assert.Equal(t, i, int64(15))
	assert.Equal(t, samples, int(5))
}

func TestWindow_Zero(t *testing.T) {
	win := Window{
		window:      time.Second * 5,
		granularity: time.Second,
		samples:     []int64{7, 0, 0, 0, 0},
		pos:         0,
	}

	i, samples, err := win.Last(1)
	assert.NilError(t, err)
	assert.Equal(t, i, int64(7))
	assert.Equal(t, samples, int(1))
}
