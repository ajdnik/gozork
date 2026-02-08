package engine

import (
	"testing"
)

func TestSaveRestoreRestartHooks(t *testing.T) {
	g := NewGameState()

	saveCalled := false
	restoreCalled := false
	restartCalled := false

	g.Save = func() error {
		saveCalled = true
		return nil
	}
	g.Restore = func() error {
		restoreCalled = true
		return nil
	}
	g.Restart = func() error {
		restartCalled = true
		return nil
	}

	if err := g.Save(); err != nil {
		t.Fatalf("expected Save hook to return nil, got %v", err)
	}
	if err := g.Restore(); err != nil {
		t.Fatalf("expected Restore hook to return nil, got %v", err)
	}
	if err := g.Restart(); err != nil {
		t.Fatalf("expected Restart hook to return nil, got %v", err)
	}

	if !saveCalled || !restoreCalled || !restartCalled {
		t.Fatalf("expected all hooks to be invoked")
	}
}
