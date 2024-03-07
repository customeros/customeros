package main

import (
	"testing"
	"time"
)

func TestHeartbeat(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	intSlice := []int{0, 1, 2, 3, 5}
	const timeout = 2 * time.Second
	heartbeat := Heartbeat(done, timeout/2, intSlice...)

	<-heartbeat

	i := 0
	for i < len(intSlice) {
		select {
		case r, ok := <-heartbeat:
			if ok == false {
				return
			} else if r != 1 {
				t.Errorf(
					"heartbeat sent wrong value: got %d, expected %d",
					r,
					1,
				)
			}
			i++
		case <-time.After(timeout):
			t.Fatal("test timed out")
		}
	}
}
