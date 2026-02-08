package engine

import (
	"testing"
)

type seqRNG struct {
	vals []int
	idx  int
}

var _ RNG = (*seqRNG)(nil)

func (r *seqRNG) Intn(n int) int {
	if n <= 0 {
		return 0
	}
	if len(r.vals) == 0 {
		return 0
	}
	v := r.vals[r.idx%len(r.vals)]
	r.idx++
	if v < 0 {
		v = -v
	}
	return v % n
}

func TestPickOneRefillFromSelected(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	G.Rand = &seqRNG{vals: []int{1}}
	data := RndSelect{Selected: []string{"a", "b", "c"}}
	got := PickOne(data)
	if got != "b" {
		t.Fatalf("expected PickOne to select \"b\", got %q", got)
	}
}

func TestPickOneFromUnselected(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	G.Rand = &seqRNG{vals: []int{0}}
	data := RndSelect{Unselected: []string{"x", "y"}}
	got := PickOne(data)
	if got != "x" {
		t.Fatalf("expected PickOne to select \"x\", got %q", got)
	}
}

func TestRandom(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	G.Rand = &seqRNG{vals: []int{2}}
	o1 := &Object{Desc: "one"}
	o2 := &Object{Desc: "two"}
	o3 := &Object{Desc: "three"}
	got := Random([]*Object{o1, o2, o3})
	if got != o3 {
		t.Fatalf("expected Random to select third object")
	}
}

func TestProb(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	G.Rand = &seqRNG{vals: []int{0, 50}}
	if Prob(1, false) {
		t.Fatalf("expected Prob to be false for base=1 and roll=0")
	}
	if !Prob(60, false) {
		t.Fatalf("expected Prob to be true for base=60 and roll=50")
	}
}

func TestZprob(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	G.Rand = &seqRNG{vals: []int{0}}
	G.Lucky = true
	if Zprob(1) {
		t.Fatalf("expected Zprob to be false for base=1 and roll=0 when lucky")
	}

	G.Rand = &seqRNG{vals: []int{0}}
	G.Lucky = false
	if !Zprob(2) {
		t.Fatalf("expected Zprob to be true for base=2 and roll=0 when not lucky")
	}
}

func TestIsFlaming(t *testing.T) {
	obj := &Object{Flags: FlgFlame | FlgOn}
	if !IsFlaming(obj) {
		t.Fatalf("expected IsFlaming to be true when FlgFlame and FlgOn are set")
	}
	obj.Flags = FlgFlame
	if IsFlaming(obj) {
		t.Fatalf("expected IsFlaming to be false without FlgOn")
	}
}

func TestIsOpenable(t *testing.T) {
	obj := &Object{Flags: FlgDoor}
	if !IsOpenable(obj) {
		t.Fatalf("expected IsOpenable to be true for doors")
	}
	obj.Flags = FlgCont
	if !IsOpenable(obj) {
		t.Fatalf("expected IsOpenable to be true for containers")
	}
	obj.Flags = FlgTake
	if IsOpenable(obj) {
		t.Fatalf("expected IsOpenable to be false for non-openable objects")
	}
}
