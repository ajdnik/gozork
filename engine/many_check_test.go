package engine

import (
	"bytes"
	"strings"
	"testing"
)

func TestManyCheckRejectsMultiple(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	var out bytes.Buffer
	G.GameOutput = &out

	G.ParsedSyntx.Verb = LexItem{Norm: "take", Orig: "take", Types: WordTypes{WordVerb}}
	G.DirObjPossibles = []*Object{{Desc: "a"}, {Desc: "b"}}
	G.DetectedSyntx = &Syntax{
		Verb: "take",
		Obj1: ObjProp{LocFlags: 0},
	}

	if ManyCheck() {
		t.Fatalf("expected ManyCheck to fail for multiple direct objects")
	}
	if !strings.Contains(out.String(), "You can't use multiple") {
		t.Fatalf("expected multiple objects message, got %q", out.String())
	}
}

func TestManyCheckRejectsMultipleIndirect(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	var out bytes.Buffer
	G.GameOutput = &out

	G.ParsedSyntx.Verb = LexItem{Norm: "put", Orig: "put", Types: WordTypes{WordVerb}}
	G.IndirObjPossibles = []*Object{{Desc: "a"}, {Desc: "b"}}
	G.DetectedSyntx = &Syntax{
		Verb: "put",
		Obj2: ObjProp{LocFlags: 0},
	}

	if ManyCheck() {
		t.Fatalf("expected ManyCheck to fail for multiple indirect objects")
	}
	if !strings.Contains(out.String(), "indirect objects") {
		t.Fatalf("expected indirect objects message, got %q", out.String())
	}
}
