package engine

import (
	"testing"
)

func TestHandleNumberIgnoresInvalid(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	G.LexRes = []LexItem{{Norm: "12x", Orig: "12x"}}
	HandleNumber(0)
	if G.LexRes[0].Norm == "intnum" {
		t.Fatalf("expected invalid number to remain non-intnum")
	}
	if G.Params.Number != 0 {
		t.Fatalf("expected Params.Number to remain 0")
	}
}

func TestHandleNumberSkipsLargeValues(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	G.LexRes = []LexItem{{Norm: "100000", Orig: "100000"}}
	HandleNumber(0)
	if G.LexRes[0].Norm != "intnum" {
		t.Fatalf("expected token to be normalized to intnum even for large values")
	}
}

func TestHandleNumberRejectsInvalidHour(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	G.LexRes = []LexItem{{Norm: "25:00", Orig: "25:00"}}
	HandleNumber(0)
	if G.LexRes[0].Norm != "intnum" {
		t.Fatalf("expected token to normalize to intnum")
	}
	if G.Params.Number != 0 {
		t.Fatalf("expected invalid hour to leave Params.Number at 0, got %d", G.Params.Number)
	}
}
