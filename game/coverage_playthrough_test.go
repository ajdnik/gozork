package game

import (
	"bytes"
	"strings"
	"testing"

	. "github.com/ajdnik/gozork/engine"
)

// Additional playthroughs to exercise broader game behavior.

// runScriptNoThief is a variant of runScript that disables the thief daemon.
func runScriptNoThief(t *testing.T, steps []Step) {
	t.Helper()

	var input strings.Builder
	for _, s := range steps {
		input.WriteString(s.Command + "\n")
	}
	if G == nil {
		G = NewGameState()
	}
	var output bytes.Buffer
	G.GameInput = strings.NewReader(input.String())
	G.GameOutput = &output
	G.Rand = newSeededRNG(1)
	G.Reader = nil

	InitGame()
	QueueInt("iThief", false).Run = false
	vLook(ActUnk)
	MainLoop()

	raw := output.String()
	segments := splitByPrompt(raw)

	for i, step := range steps {
		segIdx := i + 1
		var seg string
		if segIdx < len(segments) {
			seg = segments[segIdx]
		}
		for _, want := range step.Contains {
			if !strings.Contains(seg, want) {
				t.Errorf("step %d %q: expected output to contain %q\ngot:\n%s",
					i, step.Command, want, seg)
			}
		}
		for _, reject := range step.Excludes {
			if strings.Contains(seg, reject) {
				t.Errorf("step %d %q: expected output NOT to contain %q\ngot:\n%s",
					i, step.Command, reject, seg)
			}
		}
	}
}

func TestPlaythroughMaintenanceLeakAndPutty(t *testing.T) {
	steps := setupUnderground()
	steps = append(steps,
		Step{Command: "north", Contains: []string{"Chasm"}},
		Step{Command: "northeast", Contains: []string{"reservoir South"}},
		Step{Command: "east", Contains: []string{"dam"}},
		Step{Command: "north", Contains: []string{"dam Lobby"}},
		Step{Command: "north", Contains: []string{"Maintenance Room"}},
		Step{Command: "push blue button", Contains: []string{"leak"}},
		Step{Command: "take tube", Contains: []string{"Taken"}},
		Step{Command: "open tube"},
		Step{Command: "squeeze tube", Contains: []string{"oozes"}},
		Step{Command: "put gunk on leak", Contains: []string{"managed to stop the leak"}},
	)
	runScriptNoThief(t, steps)
}

func TestPlaythroughBoatAndBuoy(t *testing.T) {
	steps := setupUnderground()
	steps = append(steps,
		Step{Command: "north", Contains: []string{"Chasm"}},
		Step{Command: "northeast", Contains: []string{"reservoir South"}},
		Step{Command: "east", Contains: []string{"dam"}},
		Step{Command: "north", Contains: []string{"dam Lobby"}},
		Step{Command: "north", Contains: []string{"Maintenance Room"}},
		Step{Command: "take wrench", Contains: []string{"Taken"}},
		Step{Command: "push yellow button", Contains: []string{"Click"}},
		Step{Command: "south", Contains: []string{"dam Lobby"}},
		Step{Command: "south", Contains: []string{"dam"}},
		Step{Command: "turn bolt with wrench", Contains: []string{"sluice gates"}},
		Step{Command: "wait"},
		Step{Command: "wait"},
		Step{Command: "wait"},
		Step{Command: "west", Contains: []string{"reservoir South"}},
		Step{Command: "north", Contains: []string{"reservoir"}},
		Step{Command: "north", Contains: []string{"reservoir North"}},
		Step{Command: "take pump", Contains: []string{"Taken"}},
		Step{Command: "south", Contains: []string{"reservoir"}},
		Step{Command: "south", Contains: []string{"reservoir South"}},
		Step{Command: "east", Contains: []string{"dam"}},
		Step{Command: "east", Contains: []string{"dam Base"}},
		Step{Command: "inflate plastic with pump", Contains: []string{"boat inflates"}},
		Step{Command: "board boat"},
		Step{Command: "launch boat"},
		Step{Command: "wait"},
		Step{Command: "wait"},
		Step{Command: "wait"},
		Step{Command: "look", Contains: []string{"buoy"}},
		Step{Command: "take buoy", Contains: []string{"Taken"}},
		Step{Command: "go east", Contains: []string{"Sandy Beach"}},
		Step{Command: "exit"},
		Step{Command: "open buoy", Contains: []string{"emerald"}},
		Step{Command: "take emerald", Contains: []string{"Taken"}},
	)
	runScriptNoThief(t, steps)
}
