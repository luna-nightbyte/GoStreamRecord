package utils

import (
	"time"
)

// FPSTimer is a simple timer to calculate frames per second (FPS).
// It keeps track of the start time and the number of frames that have passed.
type FPSTimer struct {
	startTime  time.Time
	frameCount int
}

// NewFPSTimer creates and returns a new FPSTimer.
// The timer starts as soon as it is created.
func NewFPSTimer() *FPSTimer {
	return &FPSTimer{
		startTime: time.Now(),
	}
}

func (t *FPSTimer) Count() int {
	return t.frameCount
}

// Tick increments the frame counter.
// This function should be called once per frame or loop iteration.
func (t *FPSTimer) Tick() {
	t.frameCount++
}

// FPS returns the current frames per second as a floating-point number.
// It calculates the elapsed time and divides the frame count by it.
func (t *FPSTimer) FPS() float64 {
	elapsed := time.Since(t.startTime).Seconds()
	if elapsed == 0 {
		return 0
	}
	return float64(t.frameCount) / elapsed
}

// Reset resets the timer, setting the start time to the current time and the
// frame count back to zero.
func (t *FPSTimer) Reset() {
	t.startTime = time.Now()
	t.frameCount = 0
}
