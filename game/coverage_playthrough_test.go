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

// runScriptNoThiefNoBats disables both the thief and bat room shenanigans.
func runScriptNoThiefNoBats(t *testing.T, steps []Step) {
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
	batRoom.Action = nil
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
	runScriptNoThiefNoBats(t, steps)
}

func TestPlaythroughBoatAndBuoy(t *testing.T) {
	steps := setupToLivingRoom()
	steps = append(steps,
		Step{Command: "take lamp", Contains: []string{"Taken"}},
		Step{Command: "take sword", Contains: []string{"Taken"}},
		Step{Command: "east", Contains: []string{"kitchen"}},
		Step{Command: "take bottle", Contains: []string{"Taken"}},
		Step{Command: "west", Contains: []string{"Living Room"}},
		Step{Command: "move rug", Contains: []string{"rug is moved"}},
		Step{Command: "open trap door", Contains: []string{"rickety staircase"}},
		Step{Command: "turn on lamp", Contains: []string{"now on"}},
		Step{Command: "down", Contains: []string{"cellar"}},
		Step{Command: "north", Contains: []string{"troll Room"}},
	)
	for i := 0; i < 15; i++ {
		steps = append(steps, Step{Command: "kill troll with sword"})
	}
	steps = append(steps,
		Step{Command: "drop sword", Contains: []string{"Dropped"}},
		Step{Command: "east", Contains: []string{"East-West Passage"}},
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
		Step{Command: "swim in water", Contains: []string{"Swimming isn't usually allowed"}},
		Step{Command: "fill bottle", Contains: []string{}},
		Step{Command: "north", Contains: []string{"reservoir North"}},
		Step{Command: "take pump", Contains: []string{"Taken"}},
		Step{Command: "south", Contains: []string{"reservoir"}},
		Step{Command: "south", Contains: []string{"reservoir South"}},
		Step{Command: "east", Contains: []string{"dam"}},
		Step{Command: "east", Contains: []string{"dam Base"}},
		Step{Command: "inflate plastic with pump", Contains: []string{"boat inflates"}},
		Step{Command: "board boat"},
		Step{Command: "launch boat"},
		Step{Command: "overboard lamp"},
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

func TestPlaythroughTempleHadesRitual(t *testing.T) {
	steps := []Step{
		{Command: "open mailbox", Contains: []string{"leaflet"}},
		{Command: "take leaflet", Contains: []string{"Taken"}},
	}
	steps = append(steps, setupToLivingRoom()...)
	steps = append(steps,
		Step{Command: "take lamp", Contains: []string{"Taken"}},
		Step{Command: "take sword", Contains: []string{"Taken"}},
		Step{Command: "turn on lamp", Contains: []string{"now on"}},
		Step{Command: "east", Contains: []string{"kitchen"}},
		Step{Command: "up", Contains: []string{"attic"}},
		Step{Command: "take rope", Contains: []string{"Taken"}},
		Step{Command: "down", Contains: []string{"kitchen"}},
		Step{Command: "west", Contains: []string{"Living Room"}},
		Step{Command: "move rug", Contains: []string{"rug is moved"}},
		Step{Command: "open trap door", Contains: []string{"rickety staircase"}},
		Step{Command: "down", Contains: []string{"cellar"}},
		Step{Command: "north", Contains: []string{"troll Room"}},
	)
	for i := 0; i < 15; i++ {
		steps = append(steps, Step{Command: "kill troll with sword"})
	}
	steps = append(steps,
		Step{Command: "drop sword", Contains: []string{"Dropped"}},
		Step{Command: "east", Contains: []string{"East-West Passage"}},
		// Detour for matches.
		Step{Command: "north", Contains: []string{"Chasm"}},
		Step{Command: "northeast", Contains: []string{"reservoir South"}},
		Step{Command: "east", Contains: []string{"dam"}},
		Step{Command: "north", Contains: []string{"dam Lobby"}},
		Step{Command: "take matches", Contains: []string{"Taken"}},
		Step{Command: "south", Contains: []string{"dam"}},
		Step{Command: "south", Contains: []string{"Deep Canyon"}},
		Step{Command: "southwest", Contains: []string{"Passage"}},
		Step{Command: "south", Contains: []string{"Round Room"}},
		Step{Command: "southeast", Contains: []string{"engravings Cave"}},
		Step{Command: "east", Contains: []string{"Dome Room"}},
		Step{Command: "tie rope to railing", Contains: []string{"rope"}},
		Step{Command: "down", Contains: []string{"torch Room"}},
		Step{Command: "take torch", Contains: []string{"Taken"}},
		Step{Command: "drop leaflet", Contains: []string{"Dropped"}},
		Step{Command: "burn leaflet with torch", Contains: []string{"catches fire"}},
		Step{Command: "south", Contains: []string{"Temple"}},
		Step{Command: "take bell", Contains: []string{"Taken"}},
		Step{Command: "south", Contains: []string{"altar"}},
		Step{Command: "take candles", Contains: []string{"Taken"}},
		Step{Command: "take book", Contains: []string{"Taken"}},
		Step{Command: "down", Contains: []string{"Cave"}},
		Step{Command: "down", Contains: []string{"Entrance to Hades"}},
		Step{Command: "ring bell", Contains: []string{"red hot"}},
		Step{Command: "light match", Contains: []string{"burn"}},
		Step{Command: "light candles with match", Contains: []string{"candles"}},
		Step{Command: "read book", Contains: []string{"Commandment"}},
	)
	runScriptNoThief(t, steps)
}

func TestPlaythroughMirrorRoom(t *testing.T) {
	steps := setupToLivingRoom()
	steps = append(steps,
		Step{Command: "take lamp", Contains: []string{"Taken"}},
		Step{Command: "take sword", Contains: []string{"Taken"}},
		Step{Command: "turn on lamp", Contains: []string{"now on"}},
		Step{Command: "east", Contains: []string{"kitchen"}},
		Step{Command: "up", Contains: []string{"attic"}},
		Step{Command: "take rope", Contains: []string{"Taken"}},
		Step{Command: "down", Contains: []string{"kitchen"}},
		Step{Command: "west", Contains: []string{"Living Room"}},
		Step{Command: "move rug", Contains: []string{"rug is moved"}},
		Step{Command: "open trap door", Contains: []string{"rickety staircase"}},
		Step{Command: "down", Contains: []string{"cellar"}},
		Step{Command: "north", Contains: []string{"troll Room"}},
	)
	for i := 0; i < 15; i++ {
		steps = append(steps, Step{Command: "kill troll with sword"})
	}
	steps = append(steps,
		Step{Command: "drop sword", Contains: []string{"Dropped"}},
		Step{Command: "east", Contains: []string{"East-West Passage"}},
		Step{Command: "east", Contains: []string{"Round Room"}},
		Step{Command: "southeast", Contains: []string{"engravings Cave"}},
		Step{Command: "east", Contains: []string{"Dome Room"}},
		Step{Command: "tie rope to railing", Contains: []string{"rope"}},
		Step{Command: "down", Contains: []string{"torch Room"}},
		Step{Command: "south", Contains: []string{"Temple"}},
		Step{Command: "south", Contains: []string{"altar"}},
		Step{Command: "down", Contains: []string{"Cave"}},
		Step{Command: "north", Contains: []string{"Mirror Room"}},
		Step{Command: "rub mirror", Contains: []string{"rumble"}},
	)
	runScriptNoThief(t, steps)
}

func TestPlaythroughSurfaceActions(t *testing.T) {
	steps := []Step{
		{Command: "open mailbox", Contains: []string{"leaflet"}},
		{Command: "take leaflet", Contains: []string{"Taken"}},
		{Command: "read leaflet", Contains: []string{"ZORK"}},
		{Command: "south", Contains: []string{"South of House"}},
		{Command: "east", Contains: []string{"Behind House"}},
		{Command: "open window", Contains: []string{"open"}},
		{Command: "in", Contains: []string{"kitchen"}},
		{Command: "take bottle", Contains: []string{"Taken"}},
		{Command: "west", Contains: []string{"Living Room"}},
		{Command: "take lamp", Contains: []string{"Taken"}},
		{Command: "take sword", Contains: []string{"Taken"}},
		{Command: "find lamp", Contains: []string{"You have it"}},
		{Command: "put lamp on sword"},
		{Command: "east", Contains: []string{"kitchen"}},
		{Command: "out", Contains: []string{"Behind House"}},
		{Command: "north", Contains: []string{"North of House"}},
		{Command: "north", Contains: []string{"forest path"}},
		{Command: "cut tree with sword", Contains: []string{"Strange concept"}},
		{Command: "climb tree", Contains: []string{"Up a tree"}},
		{Command: "take nest", Contains: []string{"Taken"}},
		{Command: "drop nest", Contains: []string{"nest falls"}},
		{Command: "leap", Contains: []string{"unaccustomed daring"}},
		{Command: "throw leaflet at tree"},
		{Command: "south", Contains: []string{"North of House"}},
		{Command: "west", Contains: []string{"West of House"}},
	}
	runScriptNoThief(t, steps)
}

func TestPlaythroughSlideRoom(t *testing.T) {
	steps := setupToLivingRoom()
	steps = append(steps,
		Step{Command: "take lamp", Contains: []string{"Taken"}},
		Step{Command: "take sword", Contains: []string{"Taken"}},
		Step{Command: "turn on lamp", Contains: []string{"now on"}},
		Step{Command: "east", Contains: []string{"kitchen"}},
		Step{Command: "up", Contains: []string{"attic"}},
		Step{Command: "take rope", Contains: []string{"Taken"}},
		Step{Command: "down", Contains: []string{"kitchen"}},
		Step{Command: "west", Contains: []string{"Living Room"}},
		Step{Command: "move rug", Contains: []string{"rug is moved"}},
		Step{Command: "open trap door", Contains: []string{"rickety staircase"}},
		Step{Command: "down", Contains: []string{"cellar"}},
		Step{Command: "north", Contains: []string{"troll Room"}},
	)
	for i := 0; i < 15; i++ {
		steps = append(steps, Step{Command: "kill troll with sword"})
	}
	steps = append(steps,
		Step{Command: "drop sword", Contains: []string{"Dropped"}},
		Step{Command: "east", Contains: []string{"East-West Passage"}},
		Step{Command: "east", Contains: []string{"Round Room"}},
		Step{Command: "southeast", Contains: []string{"engravings Cave"}},
		Step{Command: "east", Contains: []string{"Dome Room"}},
		Step{Command: "tie rope to railing", Contains: []string{"rope"}},
		Step{Command: "down", Contains: []string{"torch Room"}},
		Step{Command: "south", Contains: []string{"Temple"}},
		Step{Command: "south", Contains: []string{"altar"}},
		Step{Command: "down", Contains: []string{"Cave"}},
		Step{Command: "north", Contains: []string{"Mirror Room"}},
		Step{Command: "rub mirror", Contains: []string{"rumble"}},
		Step{Command: "north", Contains: []string{"Cold Passage"}},
		Step{Command: "west", Contains: []string{"slide Room"}},
		Step{Command: "climb down slide", Contains: []string{"tumble down the slide"}},
	)
	runScriptNoThief(t, steps)
}
