package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestStats(t *testing.T) {
	stats := NewArrivalStats(100)
	for i := 0; i < 10; i++ {
		stats.Add(time.Millisecond)
	}

	if stats.Size() != 10 {
		t.Errorf("get size failed %v", stats.Size())
	}

	for i := 10; i < 200; i++ {
		stats.Add(time.Millisecond)
	}

	if stats.Size() != 100 {
		t.Errorf("get size failed %v", stats.Size())
	}

	if stats.Mean() != float64(time.Millisecond) {
		t.Errorf("calc mean failed %v, %v", stats.Mean(), float64(time.Millisecond))
	}

	//fmt.Println(stats.GetIntervals())
}

func TestWindow(t *testing.T) {
	w := NewArrivalWindow()
	now := time.Now()
	for i := 0; i < 10; i++ {
		w.Add(now)
		now = now.Add(50 * time.Millisecond)
	}

	fmt.Println(w.Mean())
	fmt.Println(w.Phi(now.Add(50 * time.Millisecond)))
	fmt.Println(w.Phi(now.Add(500 * time.Millisecond)))
	fmt.Println(w.String())
}
func TestPhiFD(t *testing.T) {
	fd := NewPhiFD()

	now := time.Now()
	for i := 0; i < 10; i++ {
		fd.Report(NodeID(0), now)
		now = now.Add(50 * time.Millisecond)
	}

	fd.Check(NodeID(1))
	time.Sleep(time.Second)
	fd.Check(NodeID(0))
	time.Sleep(50 * time.Millisecond)
	fd.Check(NodeID(0))
	time.Sleep(50 * time.Millisecond)
	fd.Check(NodeID(0))
	time.Sleep(50 * time.Millisecond)
	fd.Check(NodeID(0))
	time.Sleep(50 * time.Millisecond)
	fd.Check(NodeID(0))
	time.Sleep(50 * time.Millisecond)
	fd.Check(NodeID(0))
	time.Sleep(50 * time.Millisecond)
	fd.Check(NodeID(0))
	fd.Report(NodeID(0), time.Now())
	time.Sleep(50 * time.Millisecond)
	fd.Check(NodeID(0))
}
