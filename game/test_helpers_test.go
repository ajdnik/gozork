package game

import (
	"bytes"
	. "github.com/ajdnik/gozork/engine"
	"strings"
	"testing"
)

func setupTestGame(t *testing.T, input string) *bytes.Buffer {
	t.Helper()

	oldG := G
	G = NewGameState()
	InitGame()

	var out bytes.Buffer
	G.GameOutput = &out
	if input != "" {
		G.GameInput = strings.NewReader(input)
		G.Reader = nil
	}
	G.InputExhausted = false
	G.QuitRequested = false

	t.Cleanup(func() { G = oldG })
	return &out
}
