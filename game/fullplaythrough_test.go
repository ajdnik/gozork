package game

import "testing"

// TestFullPlaythrough runs a near-complete game based on playthrough.txt.
// Combat sections use extra rounds to handle RNG variation with seed 1.
func TestFullPlaythrough(t *testing.T) {
	var steps []Step

	// === Opening: mailbox, leaflet, enter house ===
	steps = append(steps,
		Step{Command: "open mailbox", Contains: []string{"leaflet"}},
		Step{Command: "take leaflet", Contains: []string{"Taken"}},
		Step{Command: "read leaflet", Contains: []string{"WELCOME TO ZORK"}},
		Step{Command: "drop leaflet", Contains: []string{"Dropped"}},
		Step{Command: "go south", Contains: []string{"South of House"}},
		Step{Command: "go east", Contains: []string{"Behind House"}},
		Step{Command: "open window", Contains: []string{"open"}},
		Step{Command: "enter house", Contains: []string{"Kitchen"}},
	)

	// === Living Room: take lamp, open underground ===
	steps = append(steps,
		Step{Command: "go west", Contains: []string{"Living Room"}},
		Step{Command: "take lamp", Contains: []string{"Taken"}},
		Step{Command: "move rug", Contains: []string{"rug is moved"}},
		Step{Command: "open trap door", Contains: []string{"rickety staircase"}},
		Step{Command: "turn on lamp", Contains: []string{"now on"}},
	)

	// === Underground: Gallery painting ===
	steps = append(steps,
		Step{Command: "go down", Contains: []string{"Cellar"}},
		Step{Command: "go south", Contains: []string{"East of Chasm"}},
		Step{Command: "go east", Contains: []string{"Gallery"}},
		Step{Command: "take painting", Contains: []string{"Taken"}},
		Step{Command: "go north", Contains: []string{"Studio"}},
		Step{Command: "go up", Contains: []string{"Kitchen"}}, // chimney
	)

	// === Attic: knife and rope ===
	steps = append(steps,
		Step{Command: "go up", Contains: []string{"Attic"}},
		Step{Command: "take knife", Contains: []string{"Taken"}},
		Step{Command: "take rope", Contains: []string{"Taken"}},
		Step{Command: "go down", Contains: []string{"Kitchen"}},
	)

	// === Living Room: store painting, get sword ===
	steps = append(steps,
		Step{Command: "go west", Contains: []string{"Living Room"}},
		Step{Command: "open case", Contains: []string{"Opened"}},
		Step{Command: "put painting in case", Contains: []string{"Done"}},
		Step{Command: "drop knife", Contains: []string{"Dropped"}},
		Step{Command: "take sword", Contains: []string{"Taken"}},
	)

	// === Troll fight (troll dies quickly with seed 1) ===
	steps = append(steps,
		Step{Command: "open trap door", Contains: []string{"rickety staircase"}},
		Step{Command: "go down", Contains: []string{"Cellar"}},
		Step{Command: "go north", Contains: []string{"Troll Room"}},
	)
	for i := 0; i < 5; i++ {
		steps = append(steps, Step{Command: "kill troll with sword"})
	}
	steps = append(steps,
		Step{Command: "drop sword", Contains: []string{"Dropped"}},
	)

	// === Dome Room and rope ===
	steps = append(steps,
		Step{Command: "go east", Contains: []string{"East-West Passage"}},
		Step{Command: "go east", Contains: []string{"Round Room"}},
		Step{Command: "go southeast", Contains: []string{"Engravings Cave"}},
		Step{Command: "go east", Contains: []string{"Dome Room"}},
		Step{Command: "tie rope to railing", Contains: []string{"rope"}},
		Step{Command: "go down", Contains: []string{"Torch Room"}},
	)

	// === Temple, Egyptian Room, Altar, Pray ===
	steps = append(steps,
		Step{Command: "go south", Contains: []string{"Temple"}},
		Step{Command: "go east", Contains: []string{"Egyptian Room"}},
		Step{Command: "take coffin", Contains: []string{"Taken"}},
		Step{Command: "go west", Contains: []string{"Temple"}},
		Step{Command: "go south", Contains: []string{"Altar"}},
		Step{Command: "pray", Contains: []string{"Forest"}},
	)

	// === Surface: Canyon to End of Rainbow ===
	steps = append(steps,
		Step{Command: "turn off lamp", Contains: []string{"now off"}},
		Step{Command: "go south", Contains: []string{"Forest"}},
		Step{Command: "go north", Contains: []string{"Clearing"}},
		Step{Command: "go east", Contains: []string{"Canyon View"}},
		Step{Command: "go down", Contains: []string{"Rocky Ledge"}},
		Step{Command: "go down", Contains: []string{"Canyon Bottom"}},
		Step{Command: "go north", Contains: []string{"End of Rainbow"}},
	)

	// === Rainbow puzzle ===
	steps = append(steps,
		Step{Command: "drop coffin", Contains: []string{"Dropped"}},
		Step{Command: "open coffin", Contains: []string{"sceptre"}},
		Step{Command: "take sceptre", Contains: []string{"Taken"}},
		Step{Command: "wave sceptre", Contains: []string{"rainbow"}},
		Step{Command: "take pot", Contains: []string{"Taken"}},
		Step{Command: "take coffin", Contains: []string{"Taken"}},
	)

	// === Back to house, store treasures ===
	steps = append(steps,
		Step{Command: "go southwest", Contains: []string{"Canyon Bottom"}},
		Step{Command: "go up", Contains: []string{"Rocky Ledge"}},
		Step{Command: "go up", Contains: []string{"Canyon View"}},
		Step{Command: "go northwest", Contains: []string{"Clearing"}},
		Step{Command: "go west", Contains: []string{"Behind House"}},
		Step{Command: "enter house", Contains: []string{"Kitchen"}},
		Step{Command: "open bag", Contains: []string{"lunch"}},
		Step{Command: "take garlic", Contains: []string{"Taken"}},
		Step{Command: "go west", Contains: []string{"Living Room"}},
		Step{Command: "put coffin in case", Contains: []string{"Done"}},
		Step{Command: "put gold in case", Contains: []string{"Done"}},
		Step{Command: "put sceptre in case", Contains: []string{"Done"}},
	)

	// === Dam sequence ===
	steps = append(steps,
		Step{Command: "open trap door", Contains: []string{"rickety staircase"}},
		Step{Command: "turn on lamp", Contains: []string{"now on"}},
		Step{Command: "go down", Contains: []string{"Cellar"}},
		Step{Command: "go north", Contains: []string{"Troll Room"}},
		Step{Command: "go east", Contains: []string{"East-West Passage"}},
		Step{Command: "go north", Contains: []string{"Chasm"}},
		Step{Command: "go northeast", Contains: []string{"Reservoir South"}},
		Step{Command: "go east", Contains: []string{"Dam"}},
		Step{Command: "go north", Contains: []string{"Dam Lobby"}},
		Step{Command: "take matches", Contains: []string{"Taken"}},
		Step{Command: "go north", Contains: []string{"Maintenance Room"}},
		Step{Command: "take wrench", Contains: []string{"Taken"}},
		Step{Command: "take screwdriver", Contains: []string{"Taken"}},
		Step{Command: "push yellow button"},
		Step{Command: "go south", Contains: []string{"Dam Lobby"}},
		Step{Command: "go south", Contains: []string{"Dam"}},
		Step{Command: "turn bolt with wrench", Contains: []string{"sluice gates"}},
		Step{Command: "drop wrench", Contains: []string{"Dropped"}},
	)

	// === Deep Canyon to Torch Room (avoid loud room bounce) ===
	steps = append(steps,
		Step{Command: "go south", Contains: []string{"Deep Canyon"}},
		Step{Command: "go southwest", Contains: []string{"Passage"}},
		Step{Command: "go south", Contains: []string{"Round Room"}},
		Step{Command: "go southeast", Contains: []string{"Engravings Cave"}},
		Step{Command: "go east", Contains: []string{"Dome Room"}},
		Step{Command: "go down", Contains: []string{"Torch Room"}},
		Step{Command: "take torch", Contains: []string{"Taken"}},
		Step{Command: "turn off lamp", Contains: []string{"now off"}},
	)

	// === Temple: bell, Altar: candles+book, Hades ===
	steps = append(steps,
		Step{Command: "go south", Contains: []string{"Temple"}},
		Step{Command: "take bell", Contains: []string{"Taken"}},
		Step{Command: "go south", Contains: []string{"Altar"}},
		Step{Command: "take candles", Contains: []string{"Taken"}},
		Step{Command: "take book", Contains: []string{"Taken"}},
		Step{Command: "go down", Contains: []string{"Cave"}},
		Step{Command: "go down", Contains: []string{"Entrance to Hades"}},
	)

	// === Hades puzzle ===
	steps = append(steps,
		Step{Command: "ring bell", Contains: []string{"red hot"}},
		Step{Command: "take candles", Contains: []string{"Taken"}},
		Step{Command: "light match", Contains: []string{"starts to burn"}},
		Step{Command: "light candles with match", Contains: []string{"candles are lit"}},
		Step{Command: "read book", Contains: []string{"Begone"}},
		Step{Command: "drop book", Contains: []string{"Dropped"}},
	)

	// === Land of the Dead ===
	steps = append(steps,
		Step{Command: "go south", Contains: []string{"Land of the Dead"}},
		Step{Command: "take skull", Contains: []string{"Taken"}},
		Step{Command: "go north", Contains: []string{"Entrance to Hades"}},
		Step{Command: "go up", Contains: []string{"Cave"}},
	)

	// === Mirror Room to Mine ===
	steps = append(steps,
		Step{Command: "go north", Contains: []string{"Mirror Room"}},
		Step{Command: "rub mirror", Contains: []string{"rumble"}},
		Step{Command: "go north", Contains: []string{"Cold Passage"}},
		Step{Command: "go west", Contains: []string{"Slide Room"}},
		Step{Command: "go north", Contains: []string{"Mine Entrance"}},
		Step{Command: "go west", Contains: []string{"Squeaky Room"}},
	)

	// === Bat Room, Shaft Room: deposit torch+screwdriver in basket ===
	steps = append(steps,
		Step{Command: "go north", Contains: []string{"Bat Room"}},
		Step{Command: "go east", Contains: []string{"Shaft Room"}},
		Step{Command: "put torch in basket", Contains: []string{"Done"}},
		Step{Command: "put screwdriver in basket", Contains: []string{"Done"}},
		Step{Command: "turn on lamp", Contains: []string{"now on"}},
	)

	// === Coal Mine: get coal ===
	steps = append(steps,
		Step{Command: "go north", Contains: []string{"Smelly Room"}},
		Step{Command: "go down", Contains: []string{"Gas Room"}},
		Step{Command: "go east", Contains: []string{"Coal Mine"}},
		Step{Command: "go northeast", Contains: []string{"Coal Mine"}},
		Step{Command: "go southeast", Contains: []string{"Coal Mine"}},
		Step{Command: "go southwest", Contains: []string{"Coal Mine"}},
		Step{Command: "go down", Contains: []string{"Ladder Top"}},
		Step{Command: "go down", Contains: []string{"Ladder Bottom"}},
		Step{Command: "go south", Contains: []string{"Dead End"}},
		Step{Command: "take coal", Contains: []string{"Taken"}},
	)

	// === Back to Shaft Room with coal ===
	steps = append(steps,
		Step{Command: "go north", Contains: []string{"Ladder Bottom"}},
		Step{Command: "go up", Contains: []string{"Ladder Top"}},
		Step{Command: "go up", Contains: []string{"Coal Mine"}},
		Step{Command: "go north", Contains: []string{"Coal Mine"}},
		Step{Command: "go east", Contains: []string{"Coal Mine"}},
		Step{Command: "go south", Contains: []string{"Coal Mine"}},
		Step{Command: "go north", Contains: []string{"Gas Room"}},
		Step{Command: "go up", Contains: []string{"Smelly Room"}},
		Step{Command: "go south", Contains: []string{"Shaft Room"}},
	)

	// === Lower basket, navigate to Drafty Room ===
	steps = append(steps,
		Step{Command: "put coal in basket", Contains: []string{"Done"}},
		Step{Command: "lower basket"},
		Step{Command: "go north", Contains: []string{"Smelly Room"}},
		Step{Command: "go down", Contains: []string{"Gas Room"}},
		Step{Command: "go east", Contains: []string{"Coal Mine"}},
		Step{Command: "go northeast", Contains: []string{"Coal Mine"}},
		Step{Command: "go southeast", Contains: []string{"Coal Mine"}},
		Step{Command: "go southwest", Contains: []string{"Coal Mine"}},
		Step{Command: "go down", Contains: []string{"Ladder Top"}},
		Step{Command: "go down", Contains: []string{"Ladder Bottom"}},
		Step{Command: "go west", Contains: []string{"Timber Room"}},
		Step{Command: "drop all"},
	)

	// === Drafty Room: Machine Room diamond ===
	steps = append(steps,
		Step{Command: "go west", Contains: []string{"Drafty Room"}},
		Step{Command: "take coal", Contains: []string{"Taken"}},
		Step{Command: "take screwdriver", Contains: []string{"Taken"}},
		Step{Command: "take torch", Contains: []string{"Taken"}},
		Step{Command: "go south", Contains: []string{"Machine Room"}},
		Step{Command: "open lid", Contains: []string{"lid opens"}},
		Step{Command: "put coal in machine", Contains: []string{"Done"}},
		Step{Command: "close lid", Contains: []string{"lid closes"}},
		Step{Command: "turn switch with screwdriver", Contains: []string{"dazzling display"}},
		Step{Command: "drop screwdriver", Contains: []string{"Dropped"}},
		Step{Command: "open lid", Contains: []string{"diamond"}},
		Step{Command: "take diamond", Contains: []string{"Taken"}},
	)

	// === Put diamond+torch in basket, collect items from Timber Room ===
	steps = append(steps,
		Step{Command: "go north", Contains: []string{"Drafty Room"}},
		Step{Command: "put torch in basket", Contains: []string{"Done"}},
		Step{Command: "put diamond in basket", Contains: []string{"Done"}},
		Step{Command: "go east", Contains: []string{"Timber Room"}},
		Step{Command: "take skull", Contains: []string{"Taken"}},
		Step{Command: "take lamp", Contains: []string{"Taken"}},
		Step{Command: "take garlic", Contains: []string{"Taken"}},
	)

	// === Back up through mine to Shaft Room ===
	steps = append(steps,
		Step{Command: "go east", Contains: []string{"Ladder Bottom"}},
		Step{Command: "go up", Contains: []string{"Ladder Top"}},
		Step{Command: "go up", Contains: []string{"Coal Mine"}},
		Step{Command: "go north", Contains: []string{"Coal Mine"}},
		Step{Command: "go east", Contains: []string{"Coal Mine"}},
		Step{Command: "go south", Contains: []string{"Coal Mine"}},
		Step{Command: "go north", Contains: []string{"Gas Room"}},
		Step{Command: "take bracelet", Contains: []string{"Taken"}},
		Step{Command: "go up", Contains: []string{"Smelly Room"}},
		Step{Command: "go south", Contains: []string{"Shaft Room"}},
	)

	// === Raise basket, collect diamond+torch ===
	steps = append(steps,
		Step{Command: "raise basket"},
		Step{Command: "take diamond", Contains: []string{"Taken"}},
		Step{Command: "take torch", Contains: []string{"Taken"}},
		Step{Command: "turn off lamp", Contains: []string{"now off"}},
	)

	// === Bat Room to Living Room via slide ===
	steps = append(steps,
		Step{Command: "go west", Contains: []string{"Bat Room"}},
		Step{Command: "take jade", Contains: []string{"Taken"}},
		Step{Command: "go south", Contains: []string{"Squeaky Room"}},
		Step{Command: "go east", Contains: []string{"Mine Entrance"}},
		Step{Command: "go south", Contains: []string{"Slide Room"}},
		Step{Command: "go down", Contains: []string{"Cellar"}},
		Step{Command: "go up", Contains: []string{"Living Room"}},
	)

	// === Store jade and diamond ===
	steps = append(steps,
		Step{Command: "put jade in case", Contains: []string{"Done"}},
		Step{Command: "put diamond in case", Contains: []string{"Done"}},
	)

	// === Reservoir: trunk and pump ===
	steps = append(steps,
		Step{Command: "turn on lamp", Contains: []string{"now on"}},
		Step{Command: "go down", Contains: []string{"Cellar"}},
		Step{Command: "go north", Contains: []string{"Troll Room"}},
		Step{Command: "go east", Contains: []string{"East-West Passage"}},
		Step{Command: "go north", Contains: []string{"Chasm"}},
		Step{Command: "go northeast", Contains: []string{"Reservoir South"}},
		Step{Command: "go north", Contains: []string{"Reservoir"}},
		Step{Command: "take trunk", Contains: []string{"Taken"}},
		Step{Command: "go north", Contains: []string{"Reservoir North"}},
		Step{Command: "take pump", Contains: []string{"Taken"}},
	)

	// === Atlantis Room: trident ===
	steps = append(steps,
		Step{Command: "go north", Contains: []string{"Atlantis Room"}},
		Step{Command: "drop torch", Contains: []string{"Dropped"}},
		Step{Command: "take trident", Contains: []string{"Taken"}},
	)

	// === Dam Base: inflate boat, river journey ===
	steps = append(steps,
		Step{Command: "go south", Contains: []string{"Reservoir North"}},
		Step{Command: "go south", Contains: []string{"Reservoir"}},
		Step{Command: "go south", Contains: []string{"Reservoir South"}},
		Step{Command: "go east", Contains: []string{"Dam"}},
		Step{Command: "go east", Contains: []string{"Dam Base"}},
		Step{Command: "inflate plastic with pump", Contains: []string{"boat inflates"}},
		Step{Command: "drop pump", Contains: []string{"Dropped"}},
	)

	// === River (ZIL speeds: River1=4, River2=4, River3=3, River4=2, River5=1)
	// VWait runs Clocker 3x per call.
	// Launch: River1, tick=4. EOT: 4→3.
	// wait1 (3 ticks): 3→2→1→0 fire! → River2. Queue(4).
	// wait2 (3 ticks): 4→3→2→1.
	// wait3 (3 ticks): 1→0 fire! → River3. Queue(3). 3→2→1.
	// look (1 tick): 1→0 fire! → River4. Queue(2). Player sees River4!
	// take buoy (1 tick): 2→1. Still at River4!
	// go east (1 tick): lands at Sandy Beach.
	steps = append(steps,
		Step{Command: "board boat"},
		Step{Command: "launch boat"},
		Step{Command: "wait"},
		Step{Command: "wait"},
		Step{Command: "wait"},
		Step{Command: "look", Contains: []string{"buoy"}},
		Step{Command: "take buoy", Contains: []string{"Taken"}},
		Step{Command: "go east", Contains: []string{"Sandy Beach"}},
	)

	// === Sandy Beach: dig for scarab ===
	steps = append(steps,
		Step{Command: "exit"},
		Step{Command: "drop garlic", Contains: []string{"Dropped"}},
		Step{Command: "drop buoy", Contains: []string{"Dropped"}},
		Step{Command: "take shovel", Contains: []string{"Taken"}},
		Step{Command: "go northeast", Contains: []string{"Sandy Cave"}},
		Step{Command: "dig sand with shovel"},
		Step{Command: "dig sand with shovel"},
		Step{Command: "dig sand with shovel"},
		Step{Command: "dig sand with shovel", Contains: []string{"scarab"}},
		Step{Command: "drop shovel", Contains: []string{"Dropped"}},
		Step{Command: "take scarab", Contains: []string{"Taken"}},
	)

	// === Buoy emerald, cross rainbow ===
	steps = append(steps,
		Step{Command: "go southwest", Contains: []string{"Sandy Beach"}},
		Step{Command: "open buoy", Contains: []string{"emerald"}},
		Step{Command: "take emerald", Contains: []string{"Taken"}},
		Step{Command: "take garlic", Contains: []string{"Taken"}},
		Step{Command: "go south", Contains: []string{"Shore"}},
		Step{Command: "go south", Contains: []string{"Aragain Falls"}},
		Step{Command: "cross rainbow", Contains: []string{"End of Rainbow"}},
		Step{Command: "turn off lamp", Contains: []string{"now off"}},
	)

	// === Back to house, store treasures ===
	steps = append(steps,
		Step{Command: "go southwest", Contains: []string{"Canyon Bottom"}},
		Step{Command: "go up", Contains: []string{"Rocky Ledge"}},
		Step{Command: "go up", Contains: []string{"Canyon View"}},
		Step{Command: "go northwest", Contains: []string{"Clearing"}},
		Step{Command: "go west", Contains: []string{"Behind House"}},
		Step{Command: "enter house", Contains: []string{"Kitchen"}},
		Step{Command: "go west", Contains: []string{"Living Room"}},
		Step{Command: "put emerald in case", Contains: []string{"Done"}},
		Step{Command: "put scarab in case", Contains: []string{"Done"}},
		Step{Command: "put trident in case", Contains: []string{"Done"}},
		Step{Command: "put jewels in case", Contains: []string{"Done"}},
	)

	// === Get egg from tree ===
	steps = append(steps,
		Step{Command: "go east", Contains: []string{"Kitchen"}},
		Step{Command: "go east", Contains: []string{"Behind House"}},
		Step{Command: "go north", Contains: []string{"North of House"}},
		Step{Command: "go north", Contains: []string{"Forest Path"}},
		Step{Command: "climb tree", Contains: []string{"Up a Tree"}},
		Step{Command: "take egg", Contains: []string{"Taken"}},
		Step{Command: "go down", Contains: []string{"Forest Path"}},
		Step{Command: "go south", Contains: []string{"North of House"}},
		Step{Command: "go east", Contains: []string{"Behind House"}},
		Step{Command: "enter house", Contains: []string{"Kitchen"}},
		Step{Command: "go west", Contains: []string{"Living Room"}},
	)

	// === Maze to Cyclops Room ===
	steps = append(steps,
		Step{Command: "turn on lamp", Contains: []string{"now on"}},
		Step{Command: "go down", Contains: []string{"Cellar"}},
		Step{Command: "go north", Contains: []string{"Troll Room"}},
		Step{Command: "go west", Contains: []string{"Maze"}},
		Step{Command: "go south", Contains: []string{"Maze"}},
		Step{Command: "go east", Contains: []string{"Maze"}},
		Step{Command: "go up", Contains: []string{"Maze"}},
		Step{Command: "take coins", Contains: []string{"Taken"}},
		Step{Command: "take key", Contains: []string{"Taken"}},
		Step{Command: "go southwest", Contains: []string{"Maze"}},
		Step{Command: "go east", Contains: []string{"Maze"}},
		Step{Command: "go south", Contains: []string{"Maze"}},
		Step{Command: "go southeast", Contains: []string{"Cyclops Room"}},
		Step{Command: "ulysses", Contains: []string{"cyclops"}},
	)

	// === Thief: give egg, navigate, fight ===
	// After "ulysses", the cyclops flees. Go up to Treasure Room, give egg.
	steps = append(steps,
		Step{Command: "go up", Contains: []string{"Treasure Room"}},
		Step{Command: "give egg to thief", Contains: []string{"thief"}},
		Step{Command: "go down", Contains: []string{"Cyclops Room"}},
		Step{Command: "go east", Contains: []string{"Strange Passage"}},
		Step{Command: "go east", Contains: []string{"Living Room"}},
		Step{Command: "put coins in case", Contains: []string{"Done"}},
		Step{Command: "take knife", Contains: []string{"Taken"}},
		Step{Command: "go west", Contains: []string{"Strange Passage"}},
		Step{Command: "go west", Contains: []string{"Cyclops Room"}},
		Step{Command: "go up", Contains: []string{"Treasure Room"}},
	)
	for i := 0; i < 5; i++ {
		steps = append(steps, Step{Command: "kill thief with knife"})
	}

	// === After thief: collect treasures ===
	steps = append(steps,
		Step{Command: "take all"},
		Step{Command: "drop knife"},
		Step{Command: "drop stiletto"},
		Step{Command: "drop torch"},
	)

	// === Maze to Grating Room ===
	steps = append(steps,
		Step{Command: "go down", Contains: []string{"Cyclops Room"}},
		Step{Command: "go northwest", Contains: []string{"Maze"}},
		Step{Command: "go south", Contains: []string{"Maze"}},
		Step{Command: "go west", Contains: []string{"Maze"}},
		Step{Command: "go up", Contains: []string{"Maze"}},
		Step{Command: "go down", Contains: []string{"Maze"}},
		Step{Command: "go northeast", Contains: []string{"Grating Room"}},
		Step{Command: "unlock grate with key"},
		Step{Command: "open grate", Contains: []string{"grating opens"}},
		Step{Command: "go up", Contains: []string{"Clearing"}},
	)

	// === Canary bauble from tree ===
	steps = append(steps,
		Step{Command: "go south", Contains: []string{"Forest Path"}},
		Step{Command: "climb tree", Contains: []string{"Up a Tree"}},
		Step{Command: "wind up canary", Contains: []string{"bauble"}},
		Step{Command: "go down", Contains: []string{"Forest Path"}},
		Step{Command: "drop knife"},
		Step{Command: "take bauble", Contains: []string{"Taken"}},
	)

	// === Final treasure deposits ===
	steps = append(steps,
		Step{Command: "go south", Contains: []string{"North of House"}},
		Step{Command: "go east", Contains: []string{"Behind House"}},
		Step{Command: "enter house", Contains: []string{"Kitchen"}},
		Step{Command: "go west", Contains: []string{"Living Room"}},
		Step{Command: "put bauble in case", Contains: []string{"Done"}},
		Step{Command: "put chalice in case", Contains: []string{"Done"}},
		Step{Command: "take canary from egg", Contains: []string{"Taken"}},
		Step{Command: "put canary in case", Contains: []string{"Done"}},
		Step{Command: "put egg in case", Contains: []string{"Done"}},
		Step{Command: "put bracelet in case", Contains: []string{"Done"}},
		Step{Command: "put skull in case", Contains: []string{"Done"}},
	)

	// === Echo Room: platinum bar ===
	steps = append(steps,
		Step{Command: "go down", Contains: []string{"Cellar"}},
		Step{Command: "go north", Contains: []string{"Troll Room"}},
		Step{Command: "go east", Contains: []string{"East-West Passage"}},
		Step{Command: "go east", Contains: []string{"Round Room"}},
		Step{Command: "go east", Contains: []string{"Loud Room"}},
		Step{Command: "echo", Contains: []string{"acoustics"}},
		Step{Command: "take bar", Contains: []string{"Taken"}},
	)

	// === Store bar, end ===
	steps = append(steps,
		Step{Command: "go west", Contains: []string{"Round Room"}},
		Step{Command: "go west", Contains: []string{"East-West Passage"}},
		Step{Command: "go west", Contains: []string{"Troll Room"}},
		Step{Command: "go south", Contains: []string{"Cellar"}},
		Step{Command: "go up", Contains: []string{"Living Room"}},
		Step{Command: "put bar in case", Contains: []string{"Done"}},
	)

	runScript(t, steps)
}
