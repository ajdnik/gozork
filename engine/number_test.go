package engine

import (
	"testing"
)

func TestHandleNumberParsesTime(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	G.LexRes = []LexItem{{Norm: "8:30", Orig: "8:30"}}
	HandleNumber(0)
	if G.LexRes[0].Norm != "intnum" {
		t.Fatalf("expected token to be normalized to intnum")
	}
	if G.Params.Number != 8*60+30 {
		t.Fatalf("expected time to parse to minutes, got %d", G.Params.Number)
	}
}

func TestHandleNumberParsesPlainNumber(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	G.LexRes = []LexItem{{Norm: "42", Orig: "42"}}
	HandleNumber(0)
	if G.LexRes[0].Norm != "intnum" {
		t.Fatalf("expected token to be normalized to intnum")
	}
	if G.Params.Number != 42 {
		t.Fatalf("expected number to parse to 42, got %d", G.Params.Number)
	}
}

func TestHandleNumberRejectsInvalidTime(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	G.LexRes = []LexItem{{Norm: "24:00", Orig: "24:00"}}
	HandleNumber(0)
	if G.LexRes[0].Norm != "intnum" {
		t.Fatalf("expected token to normalize to intnum")
	}
	if G.Params.Number != 0 {
		t.Fatalf("expected invalid time to leave Params.Number at 0, got %d", G.Params.Number)
	}
}

func TestHandleNumberParsesShortMinutes(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	G.LexRes = []LexItem{{Norm: "7:5", Orig: "7:5"}}
	HandleNumber(0)
	if G.LexRes[0].Norm != "intnum" {
		t.Fatalf("expected token to normalize to intnum")
	}
	if G.Params.Number != 19*60+5 {
		t.Fatalf("expected time to parse to minutes with PM adjustment, got %d", G.Params.Number)
	}
}

func TestHandleNumberParsesLateTime(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	G.LexRes = []LexItem{{Norm: "23:59", Orig: "23:59"}}
	HandleNumber(0)
	if G.LexRes[0].Norm != "intnum" {
		t.Fatalf("expected token to normalize to intnum")
	}
	if G.Params.Number != 23*60+59 {
		t.Fatalf("expected time to parse to minutes, got %d", G.Params.Number)
	}
}

func TestHandleNumberParsesMidnight(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	G.LexRes = []LexItem{{Norm: "0:00", Orig: "0:00"}}
	HandleNumber(0)
	if G.LexRes[0].Norm != "intnum" {
		t.Fatalf("expected token to normalize to intnum")
	}
	if G.Params.Number != 0 {
		t.Fatalf("expected midnight to parse to 0 minutes, got %d", G.Params.Number)
	}
}

func TestHandleNumberParsesNoon(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	G.LexRes = []LexItem{{Norm: "12:00", Orig: "12:00"}}
	HandleNumber(0)
	if G.LexRes[0].Norm != "intnum" {
		t.Fatalf("expected token to normalize to intnum")
	}
	if G.Params.Number != 12*60 {
		t.Fatalf("expected noon to parse to minutes, got %d", G.Params.Number)
	}
}

func TestHandleNumberParsesAfterNoon(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	G.LexRes = []LexItem{{Norm: "12:01", Orig: "12:01"}}
	HandleNumber(0)
	if G.LexRes[0].Norm != "intnum" {
		t.Fatalf("expected token to normalize to intnum")
	}
	if G.Params.Number != 12*60+1 {
		t.Fatalf("expected time to parse to minutes, got %d", G.Params.Number)
	}
}
