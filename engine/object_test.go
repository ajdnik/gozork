package engine

import (
	"testing"
)

func TestDirectionStringAndParse(t *testing.T) {
	if got := North.String(); got != "north" {
		t.Fatalf("expected north, got %q", got)
	}
	if got := Direction(99).String(); got != "unknown" {
		t.Fatalf("expected unknown for out of range, got %q", got)
	}
	if dir, ok := StringToDir("north"); !ok || dir != North {
		t.Fatalf("expected StringToDir to parse north")
	}
	if _, ok := StringToDir("nope"); ok {
		t.Fatalf("expected StringToDir to reject unknown")
	}
}

func TestAnyFlagIn(t *testing.T) {
	if !AnyFlagIn(0, FlgTake) {
		t.Fatalf("expected AnyFlagIn to allow zero required flags")
	}
	if AnyFlagIn(FlgTake, FlgOpen) {
		t.Fatalf("expected AnyFlagIn to be false when no bits overlap")
	}
	if !AnyFlagIn(FlgOpen|FlgTake, FlgTake) {
		t.Fatalf("expected AnyFlagIn to be true when bits overlap")
	}
}

func TestObjectChildrenAndMove(t *testing.T) {
	parent := &Object{Desc: "parent"}
	child := &Object{Desc: "child"}
	parent.AddChild(child)
	if !parent.HasChildren() || parent.Children[0] != child {
		t.Fatalf("expected child to be added to parent children")
	}
	parent.RemoveChild(child)
	if parent.HasChildren() {
		t.Fatalf("expected child to be removed")
	}

	dest := &Object{Desc: "dest"}
	child.MoveTo(dest)
	if child.Location() != dest {
		t.Fatalf("expected child to move to dest")
	}
	if len(dest.Children) != 1 || dest.Children[0] != child {
		t.Fatalf("expected dest to contain child")
	}

	child.Remove()
	if child.Location() != nil {
		t.Fatalf("expected child to have no parent after Remove")
	}
}

func TestExitPropsSetGet(t *testing.T) {
	room := &Object{}
	if _, ok := room.GetExit(North); ok {
		t.Fatalf("expected empty exits to return false")
	}
	target := &Object{Desc: "target"}
	room.SetExit(North, ExitProps{UExit: true, RExit: target})
	dp, ok := room.GetExit(North)
	if !ok || dp.RExit != target {
		t.Fatalf("expected GetExit to return the configured exit")
	}
	room.SetExit(South, ExitProps{})
	if _, ok := room.GetExit(South); ok {
		t.Fatalf("expected unset ExitProps to be considered empty")
	}
}

func TestObjectAccessors(t *testing.T) {
	obj := &Object{}
	if obj.GetSize() != 0 || obj.GetValue() != 0 || obj.GetTValue() != 0 || obj.GetCapacity() != 0 {
		t.Fatalf("expected zero item accessors for nil ItemData")
	}
	if obj.GetStrength() != 0 {
		t.Fatalf("expected zero strength for nil CombatData")
	}
	if obj.GetVehType() != 0 {
		t.Fatalf("expected zero vehicle type for nil VehicleData")
	}

	obj.SetSize(3)
	obj.SetValue(7)
	obj.SetTValue(9)
	obj.SetCapacity(11)
	if obj.GetSize() != 3 || obj.GetValue() != 7 || obj.GetTValue() != 9 || obj.GetCapacity() != 11 {
		t.Fatalf("expected item accessors to round-trip")
	}

	obj.SetStrength(5)
	if obj.GetStrength() != 5 {
		t.Fatalf("expected combat strength to round-trip")
	}

	obj.SetVehType(FlgWater)
	if obj.GetVehType() != FlgWater {
		t.Fatalf("expected vehicle type to round-trip")
	}
}

func TestObjectItemSettersWithExistingItem(t *testing.T) {
	obj := &Object{
		Item: &ItemData{TValue: 1, Capacity: 2},
	}

	obj.SetTValue(9)
	if obj.Item.TValue != 9 {
		t.Fatalf("expected TValue to be updated, got %d", obj.Item.TValue)
	}

	obj.SetCapacity(7)
	if obj.Item.Capacity != 7 {
		t.Fatalf("expected Capacity to be updated, got %d", obj.Item.Capacity)
	}
}

func TestExitPropsIsSet(t *testing.T) {
	target := &Object{Desc: "target"}
	if (ExitProps{}).IsSet() {
		t.Fatalf("expected empty ExitProps to be unset")
	}
	if !(ExitProps{NExit: "nope"}).IsSet() {
		t.Fatalf("expected NExit to mark set")
	}
	if !(ExitProps{UExit: true, RExit: target}).IsSet() {
		t.Fatalf("expected UExit with RExit to mark set")
	}
	if !(ExitProps{FExit: func() *Object { return target }}).IsSet() {
		t.Fatalf("expected FExit to mark set")
	}
	if !(ExitProps{CExit: func() bool { return true }, RExit: target}).IsSet() {
		t.Fatalf("expected CExit with RExit to mark set")
	}
	if !(ExitProps{DExit: target, RExit: target}).IsSet() {
		t.Fatalf("expected DExit with RExit to mark set")
	}
}
