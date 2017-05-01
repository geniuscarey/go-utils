package utils

import (
	"fmt"
	"time"
)

type FailureDetector interface {
}

type NodeID uint32

type PhiFailureDetector struct {
	nodesWindow map[NodeID]*ArrivalWindow
	lastCheck   time.Time
	lastPause   time.Time
}

func NewPhiFD() *PhiFailureDetector {
	return &PhiFailureDetector{
		nodesWindow: make(map[NodeID]*ArrivalWindow),
		lastCheck:   time.Now(),
	}
}

func (this *PhiFailureDetector) Report(id NodeID, now time.Time) {
	nodeHwd, ok := this.nodesWindow[id]
	if !ok {
		nodeHwd = NewArrivalWindow()
		nodeHwd.Add(now)
		this.nodesWindow[id] = nodeHwd
	} else {
		nodeHwd.Add(now)
	}

	fmt.Printf("node %v heatbeat window's mean:%v\n", id, nodeHwd.Mean()/float64(time.Millisecond))
}

func (this *PhiFailureDetector) MaxPauseTime() time.Duration {
	return 200 * time.Millisecond
}

func (this *PhiFailureDetector) Check(id NodeID) {
	nodeHwd, ok := this.nodesWindow[id]
	if !ok {
		//fmt.Printf("node %v has no data\n", id)
		fmt.Printf("node %v has no data\n", id)
		return
	}

	now := time.Now()
	diff := now.Sub(this.lastCheck)
	pauseDiff := now.Sub(this.lastPause)
	this.lastCheck = now

	if diff > this.MaxPauseTime() {
		fmt.Printf("skip this check due to large shedule delay %v\n", diff)
		this.lastPause = now
		return
	}

	if pauseDiff < this.MaxPauseTime() {
		fmt.Printf("protect not over, skip this check: %v\n", pauseDiff)
		return
	}

	phi := nodeHwd.Phi(now)

	if phi > 8 {
		fmt.Printf("Node %v is fault at %v, phi value: %v\n", id, now, phi)
	}
}
