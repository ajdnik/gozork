package engine

import (
	"testing"
)

func TestObjIndexHelpers(t *testing.T) {
	oldG := G
	oldIndex := ObjIndex
	G = NewGameState()
	ObjIndex = nil
	t.Cleanup(func() {
		G = oldG
		ObjIndex = oldIndex
	})

	o1 := &Object{Desc: "o1"}
	o2 := &Object{Desc: "o2"}
	G.AllObjects = []*Object{o1, o2}

	BuildObjIndex()

	if ObjToIdx(nil) != -1 {
		t.Fatalf("expected ObjToIdx(nil) to return -1")
	}
	if ObjToIdx(o1) != 0 || ObjToIdx(o2) != 1 {
		t.Fatalf("unexpected ObjToIdx results: %d, %d", ObjToIdx(o1), ObjToIdx(o2))
	}
	if IdxToObj(0) != o1 || IdxToObj(1) != o2 {
		t.Fatalf("unexpected IdxToObj results")
	}
	if IdxToObj(2) != nil {
		t.Fatalf("expected IdxToObj out of range to return nil")
	}
	if ObjToIdx(&Object{Desc: "unknown"}) != -1 {
		t.Fatalf("expected ObjToIdx unknown to return -1")
	}

	o3 := &Object{Desc: "o3"}
	G.AllObjects = append(G.AllObjects, o3)
	BuildObjIndex()
	if ObjToIdx(o3) != -1 {
		t.Fatalf("expected BuildObjIndex to be idempotent and not rebuild index")
	}
}
