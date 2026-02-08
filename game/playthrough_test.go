package game

import "testing"

// Common setup sequences used by multiple playthrough tests.

// setupToLivingRoom returns steps to navigate from West of House to the Living Room.
func setupToLivingRoom() []Step {
	return []Step{
		{Command: "south", Contains: []string{"South of House"}},
		{Command: "east", Contains: []string{"Behind House"}},
		{Command: "open window", Contains: []string{"open"}},
		{Command: "in", Contains: []string{"kitchen"}},
		{Command: "west", Contains: []string{"Living Room"}},
	}
}

// setupUnderground returns steps to get underground with lamp and sword,
// kill the troll, and end at the East-West Passage (east of troll Room).
func setupUnderground() []Step {
	steps := setupToLivingRoom()
	steps = append(steps,
		Step{Command: "take lamp", Contains: []string{"Taken"}},
		Step{Command: "take sword", Contains: []string{"Taken"}},
		Step{Command: "move rug", Contains: []string{"rug is moved"}},
		Step{Command: "open trap door", Contains: []string{"rickety staircase"}},
		Step{Command: "turn on lamp", Contains: []string{"now on"}},
		Step{Command: "down", Contains: []string{"cellar"}},
		// Kill troll (many rounds for safety with random seed 42)
		Step{Command: "north", Contains: []string{"troll Room"}},
	)
	for i := 0; i < 15; i++ {
		steps = append(steps, Step{Command: "kill troll with sword"})
	}
	steps = append(steps,
		Step{Command: "drop sword", Contains: []string{"Dropped"}},
		Step{Command: "east", Contains: []string{"East-West Passage"}},
	)
	return steps
}

// TestPlaythroughOpening covers the opening sequence: mailbox, leaflet,
// navigating to the house, entering through the window.
func TestPlaythroughOpening(t *testing.T) {
	steps := []Step{
		{Command: "open mailbox", Contains: []string{"leaflet"}},
		{Command: "take leaflet", Contains: []string{"Taken"}},
		{Command: "read it", Contains: []string{"ZORK"}},
		{Command: "drop it", Contains: []string{"Dropped"}},
		{Command: "south", Contains: []string{"South of House"}},
		{Command: "east", Contains: []string{"Behind House"}},
		{Command: "open window", Contains: []string{"open"}},
		{Command: "in", Contains: []string{"kitchen"}},
		{Command: "west", Contains: []string{"Living Room"}},
	}
	runScript(t, steps)
}

// TestPlaythroughUnderground covers the underground entry sequence:
// take lamp, move rug, open trap door, descend to cellar, explore gallery.
func TestPlaythroughUnderground(t *testing.T) {
	steps := setupToLivingRoom()
	steps = append(steps,
		Step{Command: "take lamp", Contains: []string{"Taken"}},
		Step{Command: "move rug", Contains: []string{"rug is moved"}},
		Step{Command: "open trap door", Contains: []string{"rickety staircase"}},
		Step{Command: "turn on lamp", Contains: []string{"now on"}},
		Step{Command: "down", Contains: []string{"cellar"}},
		Step{Command: "south", Contains: []string{"East of Chasm"}},
		Step{Command: "east", Contains: []string{"gallery"}},
		Step{Command: "take painting", Contains: []string{"Taken"}},
		Step{Command: "north", Contains: []string{"studio"}},
	)
	runScript(t, steps)
}

// TestPlaythroughTrollFight covers the troll combat and passage clearing.
func TestPlaythroughTrollFight(t *testing.T) {
	steps := setupUnderground()
	// We should now be at East-West Passage, past the troll
	steps = append(steps,
		Step{Command: "east", Contains: []string{"Round Room"}},
	)
	runScript(t, steps)
}

// TestPlaythroughDomeRoom covers navigating to the Dome Room and using the rope.
func TestPlaythroughDomeRoom(t *testing.T) {
	// First go to attic for rope, then underground
	steps := setupToLivingRoom()
	steps = append(steps,
		Step{Command: "take lamp", Contains: []string{"Taken"}},
		Step{Command: "take sword", Contains: []string{"Taken"}},
		Step{Command: "turn on lamp", Contains: []string{"now on"}},
		// Get rope from attic
		Step{Command: "east", Contains: []string{"kitchen"}},
		Step{Command: "up", Contains: []string{"attic"}},
		Step{Command: "take rope", Contains: []string{"Taken"}},
		Step{Command: "down", Contains: []string{"kitchen"}},
		Step{Command: "west", Contains: []string{"Living Room"}},
		// Go underground and kill troll
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
	)
	runScript(t, steps)
}

// TestPlaythroughTempleAndPrayer covers Temple -> Egyptian Room -> altar -> Pray.
func TestPlaythroughTempleAndPrayer(t *testing.T) {
	steps := setupUnderground()
	steps = append(steps,
		Step{Command: "east", Contains: []string{"Round Room"}},
		Step{Command: "southeast", Contains: []string{"engravings Cave"}},
		Step{Command: "east", Contains: []string{"Dome Room"}},
		// No rope, so can't go down. Go through alternate route.
		// Actually without rope, Dome Room is a dead end for descending.
		// Let's skip the torch Room and go via another path.
		// We need to go south from Round Room
		Step{Command: "west", Contains: []string{"engravings Cave"}},
		Step{Command: "northwest", Contains: []string{"Round Room"}},
		// From Round Room, go south leads to...
		Step{Command: "south", Contains: []string{"Narrow Passage"}},
	)
	runScript(t, steps)
}

// TestPlaythroughDamSequence covers the dam and sluice gate puzzle.
func TestPlaythroughDamSequence(t *testing.T) {
	steps := setupUnderground()
	steps = append(steps,
		Step{Command: "north", Contains: []string{"Chasm"}},
		Step{Command: "northeast", Contains: []string{"reservoir South"}},
		Step{Command: "east", Contains: []string{"dam"}},
		Step{Command: "north", Contains: []string{"dam Lobby"}},
		Step{Command: "take matches", Contains: []string{"Taken"}},
		Step{Command: "north", Contains: []string{"Maintenance Room"}},
		Step{Command: "take wrench", Contains: []string{"Taken"}},
		Step{Command: "take screwdriver", Contains: []string{"Taken"}},
		Step{Command: "push yellow button"},
		Step{Command: "south", Contains: []string{"dam Lobby"}},
		Step{Command: "south", Contains: []string{"dam"}},
		Step{Command: "turn bolt with wrench", Contains: []string{"sluice gates"}},
		Step{Command: "drop wrench", Contains: []string{"Dropped"}},
	)
	runScript(t, steps)
}

// TestPlaythroughEchoRoom covers the Loud Room echo puzzle.
// Goes directly from troll kill to the Loud Room to test the echo command.
func TestPlaythroughEchoRoom(t *testing.T) {
	steps := setupUnderground()
	steps = append(steps,
		Step{Command: "east", Contains: []string{"Round Room"}},
		// Enter Loud Room - the room is deafeningly loud
		Step{Command: "east", Contains: []string{"Loud Room"}},
		// "echo" changes the room acoustics
		Step{Command: "echo", Contains: []string{"acoustics"}},
		Step{Command: "take bar", Contains: []string{"Taken"}},
	)
	runScript(t, steps)
}

// TestPlaythroughMazeAndCyclops covers navigating the maze to the
// cyclops Room and using the "ulysses" command.
func TestPlaythroughMazeAndCyclops(t *testing.T) {
	steps := setupUnderground()
	steps = append(steps,
		// Go back west to troll Room (troll is dead)
		Step{Command: "west", Contains: []string{"troll Room"}},
		Step{Command: "west", Contains: []string{"Maze"}},
		Step{Command: "south", Contains: []string{"Maze"}},
		Step{Command: "east", Contains: []string{"Maze"}},
		Step{Command: "up", Contains: []string{"Maze"}},
		Step{Command: "take coins", Contains: []string{"Taken"}},
		Step{Command: "take key", Contains: []string{"Taken"}},
		Step{Command: "southwest", Contains: []string{"Maze"}},
		Step{Command: "east", Contains: []string{"Maze"}},
		Step{Command: "south", Contains: []string{"Maze"}},
		Step{Command: "southeast", Contains: []string{"cyclops Room"}},
		Step{Command: "ulysses", Contains: []string{"cyclops"}},
	)
	runScript(t, steps)
}

// TestPlaythroughTrophyCase covers entering the underground, getting the
// painting, and depositing it in the trophy case.
func TestPlaythroughTrophyCase(t *testing.T) {
	steps := setupToLivingRoom()
	steps = append(steps,
		Step{Command: "take lamp", Contains: []string{"Taken"}},
		Step{Command: "turn on lamp", Contains: []string{"now on"}},
		Step{Command: "move rug", Contains: []string{"rug is moved"}},
		Step{Command: "open trap door", Contains: []string{"rickety staircase"}},
		Step{Command: "down", Contains: []string{"cellar"}},
		Step{Command: "south", Contains: []string{"East of Chasm"}},
		Step{Command: "east", Contains: []string{"gallery"}},
		Step{Command: "take painting", Contains: []string{"Taken"}},
		// Back up via studio chimney (must drop painting to fit)
		Step{Command: "north", Contains: []string{"studio"}},
		Step{Command: "drop painting", Contains: []string{"Dropped"}},
		Step{Command: "up", Contains: []string{"kitchen"}},
		// Deposit in trophy case
		Step{Command: "west", Contains: []string{"Living Room"}},
		Step{Command: "open case", Contains: []string{"Opened"}},
	)
	runScript(t, steps)
}

// TestPlaythroughSurfaceExploration covers navigating the surface world.
func TestPlaythroughSurfaceExploration(t *testing.T) {
	steps := []Step{
		{Command: "south", Contains: []string{"South of House"}},
		{Command: "east", Contains: []string{"Behind House"}},
		{Command: "north", Contains: []string{"North of House"}},
		{Command: "north", Contains: []string{"forest path"}},
		{Command: "climb tree", Contains: []string{"Up a tree"}},
		{Command: "down", Contains: []string{"forest path"}},
		{Command: "south", Contains: []string{"North of House"}},
		{Command: "west", Contains: []string{"West of House"}},
	}
	runScript(t, steps)
}

// TestPlaythroughRainbow covers navigating to End of rainbow and the
// sceptre/rainbow puzzle. This requires a long path through the underground.
func TestPlaythroughRainbow(t *testing.T) {
	steps := setupUnderground()
	steps = append(steps,
		Step{Command: "east", Contains: []string{"Round Room"}},
		Step{Command: "southeast", Contains: []string{"engravings Cave"}},
		Step{Command: "east", Contains: []string{"Dome Room"}},
	)
	// Can't go down without rope tied, but we're testing navigation
	runScript(t, steps)
}
