package engine

import (
	"testing"
)

func TestDedupPreservesOrder(t *testing.T) {
	a := &Object{Desc: "a"}
	b := &Object{Desc: "b"}
	c := &Object{Desc: "c"}

	input := []*Object{a, b, a, c, b}
	out := dedup(input)
	if len(out) != 3 {
		t.Fatalf("expected 3 unique objects, got %d", len(out))
	}
	if out[0] != a || out[1] != b || out[2] != c {
		t.Fatalf("expected dedup to preserve first-seen order")
	}
}
