package zork

import (
	"bytes"
	"math/rand"
	"strings"
	"testing"
)

// Step represents one player command and the expected response properties.
type Step struct {
	Command  string   // command to type
	Contains []string // substrings that MUST appear in the response
	Excludes []string // substrings that must NOT appear in the response
}

// runScript initializes a fresh game, feeds the commands, captures all output,
// then verifies each step's expectations against the corresponding output segment.
//
// Output is segmented by scanning for the ">" prompt that precedes each command read.
// The output between two prompts is the game's response to the preceding command.
func runScript(t *testing.T, steps []Step) {
	t.Helper()

	// Seed rand for deterministic results. Seed 1 produces survivable
	// combat outcomes (seed 42 causes instant death from the troll).
	rand.Seed(1)

	// Build the command stream: all commands joined by newlines
	var input strings.Builder
	for _, s := range steps {
		input.WriteString(s.Command + "\n")
	}

	// Plug in fake I/O
	var output bytes.Buffer
	GameInput = strings.NewReader(input.String())
	GameOutput = &output
	Reader = nil // force re-init

	// Run the game, catching panics from Quit()
	func() {
		defer func() {
			if r := recover(); r != nil {
				if r != ErrQuit {
					t.Fatalf("unexpected panic: %v", r)
				}
			}
		}()
		InitGame()
		VLook(ActUnk)
		MainLoop()
	}()

	// Split output into segments aligned with commands.
	// The game prints ">" before each Read(). We use that as a delimiter.
	// The first segment (before any ">") is the opening game text.
	// Segment i+1 corresponds to steps[i].
	raw := output.String()
	segments := splitByPrompt(raw)

	// segments[0] = initial room description (before first command)
	// segments[i+1] = response to steps[i]
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

// splitByPrompt splits the raw output on the ">" prompt that the parser prints
// before each Read() call. The pattern in the output is "\n>" (newline then ">").
// We split on that boundary. segments[0] is the initial text before the first
// prompt, segments[1] is the response to the first command, etc.
func splitByPrompt(raw string) []string {
	// The parser prints "\n>" before reading. Split on that.
	parts := strings.Split(raw, "\n>")
	// The first part is the initial output (before any prompt).
	// Subsequent parts start right after a ">" prompt.
	return parts
}

// ================================================================
// GAMEPLAY SCENARIO TESTS
// ================================================================

func TestOpeningRoomDescription(t *testing.T) {
	runScript(t, []Step{
		{
			Command: "look",
			Contains: []string{
				"West of House",
				"white house",
				"small mailbox",
			},
		},
	})
}

func TestOpenMailbox(t *testing.T) {
	runScript(t, []Step{
		{
			Command:  "open mailbox",
			Contains: []string{"leaflet"},
		},
		{
			Command:  "take leaflet",
			Contains: []string{"Taken"},
		},
	})
}

func TestNavigation(t *testing.T) {
	runScript(t, []Step{
		{
			Command:  "north",
			Contains: []string{"North of House"},
		},
		{
			Command:  "east",
			Contains: []string{"Behind House"},
		},
		{
			Command:  "south",
			Contains: []string{"South of House"},
		},
		{
			Command:  "west",
			Contains: []string{"West of House"},
		},
	})
}

func TestScoreCommand(t *testing.T) {
	runScript(t, []Step{
		{
			Command:  "score",
			Contains: []string{"Your score is", "0", "Beginner"},
		},
	})
}

func TestInventoryEmpty(t *testing.T) {
	runScript(t, []Step{
		{
			Command:  "inventory",
			Contains: []string{"empty-handed"},
		},
	})
}

func TestTakeAndInventory(t *testing.T) {
	runScript(t, []Step{
		{
			Command:  "open mailbox",
			Contains: []string{"leaflet"},
		},
		{
			Command:  "take leaflet",
			Contains: []string{"Taken"},
		},
		{
			Command:  "inventory",
			Contains: []string{"leaflet"},
			Excludes: []string{"empty-handed"},
		},
	})
}

func TestEnterHouse(t *testing.T) {
	runScript(t, []Step{
		{
			Command:  "north",
			Contains: []string{"North of House"},
		},
		{
			Command:  "east",
			Contains: []string{"Behind House"},
		},
		{
			Command:  "open window",
			Contains: nil, // just should not crash
		},
		{
			Command:  "in",
			Contains: []string{"Kitchen"},
		},
	})
}

func TestKitchenContents(t *testing.T) {
	runScript(t, []Step{
		// Navigate to kitchen
		{"north", nil, nil},
		{"east", nil, nil},
		{"open window", nil, nil},
		{"in", nil, nil},
		{
			Command:  "look",
			Contains: []string{"Kitchen", "table"},
		},
	})
}

func TestLivingRoomTrophyCase(t *testing.T) {
	runScript(t, []Step{
		// Navigate to kitchen then living room
		{"north", nil, nil},
		{"east", nil, nil},
		{"open window", nil, nil},
		{"in", nil, nil},
		{
			Command:  "west",
			Contains: []string{"Living Room"},
		},
		{
			Command:  "look",
			Contains: []string{"trophy case", "sword", "lantern"},
		},
	})
}

func TestTakeLamp(t *testing.T) {
	runScript(t, []Step{
		{"north", nil, nil},
		{"east", nil, nil},
		{"open window", nil, nil},
		{"in", nil, nil},
		{"west", nil, nil},
		{
			Command:  "take lamp",
			Contains: []string{"Taken"},
		},
		{
			Command:  "inventory",
			Contains: []string{"lantern"},
		},
	})
}

func TestTakeSword(t *testing.T) {
	runScript(t, []Step{
		{"north", nil, nil},
		{"east", nil, nil},
		{"open window", nil, nil},
		{"in", nil, nil},
		{"west", nil, nil},
		{
			Command:  "take sword",
			Contains: []string{"Taken"},
		},
		{
			Command:  "inventory",
			Contains: []string{"sword"},
		},
	})
}

func TestVerboseMode(t *testing.T) {
	runScript(t, []Step{
		{
			Command:  "verbose",
			Contains: []string{"verbosity"},
		},
	})
}

func TestDiagnose(t *testing.T) {
	runScript(t, []Step{
		{
			Command:  "diagnose",
			Contains: []string{"health"},
		},
	})
}

func TestExamineMailbox(t *testing.T) {
	runScript(t, []Step{
		{
			Command:  "examine mailbox",
			Contains: []string{"mailbox"},
		},
	})
}

func TestCantGoThatWay(t *testing.T) {
	runScript(t, []Step{
		{
			Command: "east",
			// Can't go east from West of House (boarded door)
			Contains: []string{"door"},
		},
	})
}

func TestDropItem(t *testing.T) {
	runScript(t, []Step{
		{"open mailbox", nil, nil},
		{
			Command:  "take leaflet",
			Contains: []string{"Taken"},
		},
		{
			Command:  "drop leaflet",
			Contains: []string{"Dropped"},
		},
	})
}

func TestReadLeaflet(t *testing.T) {
	runScript(t, []Step{
		{"open mailbox", nil, nil},
		{"take leaflet", nil, nil},
		{
			Command:  "read it",
			Contains: []string{"ZORK"},
		},
	})
}
