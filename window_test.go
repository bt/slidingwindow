package slidingwindow

import (
	"testing"
	"time"

	"gotest.tools/assert"
)


func TestWindow_Simple(t *testing.T) {
	win, err := newWindow(time.Second * 5, time.Second)
	assert.NilError(t, err)
	fillSamples(win, 1,2,3,4,5)
	win.pos = 4

	i, samples, err := win.Last(5)
	assert.NilError(t, err)
	assert.Equal(t, i, int64(15))
	assert.Equal(t, samples, int(5))
}

func TestWindow_LastTime(t *testing.T) {
	win, err := New(time.Hour * 24, time.Minute * 10)
	assert.NilError(t, err)
	win.Add(123)
	win.Add(456)

	val, _, err := win.Last(1)
	assert.NilError(t, err)
	assert.Equal(t, val, int64(579))
}

func TestWindow_Position(t *testing.T) {
	win, err := New(time.Hour * 24, time.Minute * 10)
	assert.NilError(t, err)

	win.Add(123)
	win.nextPosition()
	win.Add(456)

	val, _, err := win.Last(1)
	assert.NilError(t, err)
	assert.Equal(t, val, int64(456))
}

func TestWindow_LoadedSamplesAppend(t *testing.T) {
	win := MustNewFromSamples(time.Hour * 24, time.Minute * 10, []int64{0, 0, 0, 0, 123})
	win.Add(456)
	val, _, err := win.Last(1)
	assert.NilError(t, err)
	assert.Equal(t, val, int64(579))
}

func TestWindow_Wrapping(t *testing.T) {
	win, err := newWindow(time.Second * 5, time.Second)
	assert.NilError(t, err)
	fillSamples(win, 1,2,3,4,5)
	win.pos = 0

	i, samples, err := win.Last(5)
	assert.NilError(t, err)
	assert.Equal(t, i, int64(15))
	assert.Equal(t, samples, int(5))
}

func TestWindow_Zero(t *testing.T) {
	win, err := newWindow(time.Second * 5, time.Second)
	assert.NilError(t, err)
	fillSamples(win, 7, 0, 0, 0, 0)
	win.pos = 0

	i, samples, err := win.Last(1)
	assert.NilError(t, err)
	assert.Equal(t, i, int64(7))
	assert.Equal(t, samples, int(1))
}

func TestWindow_LoadSamples(t *testing.T) {
	win := MustNewFromSamples(time.Minute, time.Second, []int64{10, 20, 30, 40})
	i, samples, err := win.Last(10)
	assert.NilError(t, err)
	assert.Equal(t, i, int64(100))
	assert.Equal(t, samples, int(4))

	win.Add(50)
	i, samples, err = win.Last(10)
	assert.NilError(t, err)
	assert.Equal(t, i, int64(150))
	assert.Equal(t, samples, int(4))
}
