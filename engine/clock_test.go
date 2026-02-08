package engine

import (
	"testing"
)

func TestQueueAndClocker(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	called := 0
	G.ClockFuncs = map[string]func() bool{
		"tick": func() bool {
			called++
			return true
		},
	}

	G.QueueInts = MaxQueueEvents
	G.QueueDmns = MaxQueueEvents
	G.ParserOk = true

	evt := Queue("tick", 2)
	if evt == nil || evt.Tick != 2 {
		t.Fatalf("expected queued event with tick 2")
	}
	evt.Run = true

	Clocker()
	if called != 0 {
		t.Fatalf("expected not called yet")
	}

	Clocker()
	if called != 1 {
		t.Fatalf("expected called once, got %d", called)
	}

	if evt.Tick != 0 {
		t.Fatalf("expected event tick to reach 0")
	}
}

// ---- clock_extra_test.go ----

func TestQueueIntReturnsExisting(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	G.QueueInts = MaxQueueEvents - 1
	G.QueueDmns = MaxQueueEvents - 1
	G.QueueItms[MaxQueueEvents-1] = ClockEvent{Key: "tick"}

	ev := QueueInt("tick", false)
	if ev == nil || ev.Key != "tick" {
		t.Fatalf("expected existing clock event")
	}
}

func TestQueueIntWhenFull(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	G.QueueInts = 0
	G.QueueDmns = 0

	ev := QueueInt("full", true)
	if ev == nil || ev.Key != "" {
		t.Fatalf("expected QueueInt to return first slot when full")
	}
}

func TestClockerWaitSkips(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	G.ClockWait = true
	if Clocker() {
		t.Fatalf("expected Clocker to return false when waiting")
	}
	if G.ClockWait {
		t.Fatalf("expected ClockWait to be cleared")
	}
}
