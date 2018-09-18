package slidingwindow

import (
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestWindow_Simple(t *testing.T) {
	win := slidingWindow{
		window:      time.Second * 5,
		granularity: time.Second,
		samples:     []int64{1, 2, 3, 4, 5},
		pos:         4,
	}

	i, err := win.Last(5)
	assert.NilError(t, err)
	assert.Equal(t, i, int64(15))
}

func TestWindow_Wrapping(t *testing.T) {
	win := slidingWindow{
		window:      time.Second * 5,
		granularity: time.Second,
		samples:     []int64{1, 2, 3, 4, 5},
		pos:         0,
	}

	i, err := win.Last(5)
	assert.NilError(t, err)
	assert.Equal(t, i, int64(15))
}
