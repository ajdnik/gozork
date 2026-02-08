package engine

import (
	"bytes"
	"testing"
)

func TestPrintfWritesToGameOutput(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	var buf bytes.Buffer
	G.GameOutput = &buf

	Printf("hello %s", "world")
	if got := buf.String(); got != "hello world" {
		t.Fatalf("expected Printf to write to GameOutput, got %q", got)
	}
}
