package engine

import (
	"testing"
)

func TestBuildAndResetObjectTree(t *testing.T) {
	oldG := G
	oldOriginal := originalState
	t.Cleanup(func() {
		G = oldG
		originalState = oldOriginal
	})

	G = NewGameState()
	parent := &Object{Desc: "parent"}
	child := &Object{Desc: "child", In: parent}
	child.SetStrength(2)
	child.SetValue(5)
	child.SetTValue(7)
	G.AllObjects = []*Object{parent, child}

	originalState = nil
	BuildObjectTree()
	if len(parent.Children) != 1 || parent.Children[0] != child {
		t.Fatalf("expected child to be attached after BuildObjectTree")
	}

	// Mutate state and ensure ResetObjectTree restores snapshot.
	child.Flags = FlgTake
	child.SetStrength(9)
	child.SetValue(11)
	child.SetTValue(13)
	child.Text = "changed"
	child.In = nil
	ResetObjectTree()
	if child.In != parent {
		t.Fatalf("expected child location restored")
	}
	if child.Flags != 0 {
		t.Fatalf("expected flags restored to original")
	}
	if child.GetStrength() != 2 || child.GetValue() != 5 || child.GetTValue() != 7 {
		t.Fatalf("expected combat/value restored")
	}
	if child.Text != "" {
		t.Fatalf("expected text restored")
	}
}
