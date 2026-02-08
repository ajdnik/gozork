package engine

import (
	"testing"
)

func TestIsInGlobal(t *testing.T) {
	obj := &Object{Desc: "obj"}
	container := &Object{}
	if IsInGlobal(obj, container) {
		t.Fatalf("expected false when Global is nil")
	}
	container.Global = []*Object{obj}
	if !IsInGlobal(obj, container) {
		t.Fatalf("expected true when obj is in Global list")
	}
}

func TestIsHeld(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	holder := &Object{Desc: "holder"}
	G.Winner = holder

	bag := &Object{}
	item := &Object{}
	item.In = bag
	bag.In = holder

	if !IsHeld(item) {
		t.Fatalf("expected item to be held by winner via container chain")
	}

	orphan := &Object{}
	if IsHeld(orphan) {
		t.Fatalf("expected orphan to not be held")
	}
}
