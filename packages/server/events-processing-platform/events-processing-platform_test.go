package main

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
)

func TestTemporalWorkerDoesNotKillServer(t *testing.T) {
	ctx := context.Background()
	_, cancel := context.WithCancel(ctx)
	cfg := &config.Config{ServiceName: "test-service"}
	logger := initLogger(cfg)
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	// Run Temporal worker
	go runTemporalWorker(cfg, logger, &waitGroup)

	// Start a heartbeat
	done := make(chan interface{})
	defer close(done)
	const timeout = 2 * time.Second
	heartbeat := Heartbeat(done, timeout/2)

	// test we get a heartbeat even if temporal worker fails
	intSlice := []int{0, 1, 2, 3, 5}

	<-heartbeat

	i := 0
	for i < len(intSlice) {
		select {
		case r, ok := <-heartbeat:
			if ok == false {
				t.Error("heartbeat channel closed")
			} else if r != 1 {
				t.Errorf("heartbeat sent wrong value: got %d, expected %d", r, 1)
			}
			i++
		case <-time.After(timeout):
			t.Fatal("test timed out")
		}

		logger.Debug("pulse")
		i++
	}

	// Stop the server
	cancel()

	// Wait for the server to stop
	waitGroup.Wait()
}

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
