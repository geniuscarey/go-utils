package utils

import (
	"fmt"
	"time"
)

const (
	PHI_FACTOR = 0.46
)

type arrivalStats struct {
	index     int
	mean      float64
	full      bool
	intervals []time.Duration
	sum       time.Duration
	lastPhi   float64
}

func NewArrivalStats(size int) *arrivalStats {
	return &arrivalStats{
		intervals: make([]time.Duration, size),
	}
}

func (this *arrivalStats) Add(d time.Duration) {
	if this.index == len(this.intervals) {
		this.index = 0
		this.full = true
	}

	if this.full {
		this.sum -= this.intervals[this.index]
	}

	this.intervals[this.index] = d
	this.sum += d
	this.index++
	this.mean = float64(this.sum) / float64(this.Size())
}

func (this *arrivalStats) Size() int {
	if this.full {
		return len(this.intervals)
	} else {
		return this.index
	}
}

func (this *arrivalStats) GetIntervals() []time.Duration {
	return this.intervals
}

func (this *arrivalStats) Mean() float64 {
	return this.mean
}

type ArrivalWindow struct {
	stats   *arrivalStats
	tLast   time.Time
	prevPhi float64
}

func NewArrivalWindow() *ArrivalWindow {
	return &ArrivalWindow{
		stats: NewArrivalStats(1000),
	}
}

func (this *ArrivalWindow) Mean() float64 {
	return this.stats.Mean()
}

func (this *ArrivalWindow) PrevPhi() float64 {
	return this.prevPhi
}

func (this *ArrivalWindow) Phi(now time.Time) float64 {
	i := now.Sub(this.tLast)
	this.prevPhi = float64(i) / this.Mean() * 0.46
	return this.prevPhi
}

func (this *ArrivalWindow) String() string {
	return fmt.Sprintf("%v", this.stats.GetIntervals())
}

func (this *ArrivalWindow) MaxInterval() time.Duration {
	return 10 * time.Second
}

func (this *ArrivalWindow) InitialInterval() time.Duration {
	return 50 * time.Millisecond
}

func (this *ArrivalWindow) Add(now time.Time) {
	if !this.tLast.IsZero() {
		i := now.Sub(this.tLast)

		if i <= this.MaxInterval() {
			this.stats.Add(i)
		}
	} else {
		this.stats.Add(this.InitialInterval())
	}

	this.tLast = now
}

func (this *ArrivalWindow) LastArrival() time.Time {
	return this.tLast
}
