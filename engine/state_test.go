package engine

import (
	"testing"
)

func TestNewGameStateDefaults(t *testing.T) {
	g := NewGameState()

	if !g.Lucky {
		t.Fatalf("expected Lucky default to be true")
	}
	if g.QueueInts != MaxQueueEvents || g.QueueDmns != MaxQueueEvents {
		t.Fatalf("expected QueueInts/QueueDmns to default to %d", MaxQueueEvents)
	}
	if g.GameOutput == nil || g.GameInput == nil {
		t.Fatalf("expected GameOutput/GameInput to be initialized")
	}
	if g.Actions == nil || g.PreActions == nil || g.NormVerbs == nil || g.ClockFuncs == nil {
		t.Fatalf("expected maps to be initialized")
	}
	if err := g.Save(); err == nil {
		t.Fatalf("expected Save to return default error")
	}
	if err := g.Restore(); err == nil {
		t.Fatalf("expected Restore to return default error")
	}
	if err := g.Restart(); err == nil {
		t.Fatalf("expected Restart to return default error")
	}
}
