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
		Step{Command: "enter house", Contains: []string{"kitchen"}},
	)

	// === Living Room: take lamp, open underground ===
	steps = append(steps,
		Step{Command: "go west", Contains: []string{"Living Room"}},
		Step{Command: "take lamp", Contains: []string{"Taken"}},
		Step{Command: "move rug", Contains: []string{"rug is moved"}},
		Step{Command: "open trap door", Contains: []string{"rickety staircase"}},
		Step{Command: "turn on lamp", Contains: []string{"now on"}},
	)

	// === Underground: gallery painting ===
	steps = append(steps,
		Step{Command: "go down", Contains: []string{"cellar"}},
		Step{Command: "go south", Contains: []string{"East of Chasm"}},
		Step{Command: "go east", Contains: []string{"gallery"}},
		Step{Command: "take painting", Contains: []string{"Taken"}},
		Step{Command: "go north", Contains: []string{"studio"}},
		Step{Command: "go up", Contains: []string{"kitchen"}}, // chimney
	)

	// === attic: knife and rope ===
	steps = append(steps,
		Step{Command: "go up", Contains: []string{"attic"}},
		Step{Command: "take knife", Contains: []string{"Taken"}},
		Step{Command: "take rope", Contains: []string{"Taken"}},
		Step{Command: "go down", Contains: []string{"kitchen"}},
	)

	// === Living Room: store painting, get sword ===
	steps = append(steps,
		Step{Command: "go west", Contains: []string{"Living Room"}},
		Step{Command: "open case", Contains: []string{"Opened"}},
		Step{Command: "put painting in case", Contains: []string{"Done"}},
		Step{Command: "drop knife", Contains: []string{"Dropped"}},
		Step{Command: "take sword", Contains: []string{"Taken"}},
	)

	// === troll fight (troll dies quickly with seed 1) ===
	steps = append(steps,
		Step{Command: "open trap door", Contains: []string{"rickety staircase"}},
		Step{Command: "go down", Contains: []string{"cellar"}},
		Step{Command: "go north", Contains: []string{"troll Room"}},
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
		Step{Command: "go southeast", Contains: []string{"engravings Cave"}},
		Step{Command: "go east", Contains: []string{"Dome Room"}},
		Step{Command: "tie rope to railing", Contains: []string{"rope"}},
		Step{Command: "go down", Contains: []string{"torch Room"}},
	)

	// === Temple, Egyptian Room, altar, Pray ===
	steps = append(steps,
		Step{Command: "go south", Contains: []string{"Temple"}},
		Step{Command: "go east", Contains: []string{"Egyptian Room"}},
		Step{Command: "take coffin", Contains: []string{"Taken"}},
		Step{Command: "go west", Contains: []string{"Temple"}},
		Step{Command: "go south", Contains: []string{"altar"}},
		Step{Command: "pray", Contains: []string{"forest"}},
	)

	// === Surface: Canyon to End of rainbow ===
	steps = append(steps,
		Step{Command: "turn off lamp", Contains: []string{"now off"}},
		Step{Command: "go south", Contains: []string{"forest"}},
		Step{Command: "go north", Contains: []string{"clearing"}},
		Step{Command: "go east", Contains: []string{"Canyon View"}},
		Step{Command: "go down", Contains: []string{"Rocky Ledge"}},
		Step{Command: "go down", Contains: []string{"Canyon Bottom"}},
		Step{Command: "go north", Contains: []string{"End of rainbow"}},
	)

	// === rainbow puzzle ===
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
		Step{Command: "go northwest", Contains: []string{"clearing"}},
		Step{Command: "go west", Contains: []string{"Behind House"}},
		Step{Command: "enter house", Contains: []string{"kitchen"}},
		Step{Command: "open bag", Contains: []string{"lunch"}},
		Step{Command: "take garlic", Contains: []string{"Taken"}},
		Step{Command: "go west", Contains: []string{"Living Room"}},
		Step{Command: "put coffin in case", Contains: []string{"Done"}},
		Step{Command: "put gold in case", Contains: []string{"Done"}},
		Step{Command: "put sceptre in case", Contains: []string{"Done"}},
	)

	// === dam sequence ===
	steps = append(steps,
		Step{Command: "open trap door", Contains: []string{"rickety staircase"}},
		Step{Command: "turn on lamp", Contains: []string{"now on"}},
		Step{Command: "go down", Contains: []string{"cellar"}},
		Step{Command: "go north", Contains: []string{"troll Room"}},
		Step{Command: "go east", Contains: []string{"East-West Passage"}},
		Step{Command: "go north", Contains: []string{"Chasm"}},
		Step{Command: "go northeast", Contains: []string{"reservoir South"}},
		Step{Command: "go east", Contains: []string{"dam"}},
		Step{Command: "go north", Contains: []string{"dam Lobby"}},
		Step{Command: "take matches", Contains: []string{"Taken"}},
		Step{Command: "go north", Contains: []string{"Maintenance Room"}},
		Step{Command: "take wrench", Contains: []string{"Taken"}},
		Step{Command: "take screwdriver", Contains: []string{"Taken"}},
		Step{Command: "push yellow button"},
		Step{Command: "go south", Contains: []string{"dam Lobby"}},
		Step{Command: "go south", Contains: []string{"dam"}},
		Step{Command: "turn bolt with wrench", Contains: []string{"sluice gates"}},
		Step{Command: "drop wrench", Contains: []string{"Dropped"}},
	)

	// === Deep Canyon to torch Room (avoid loud room bounce) ===
	steps = append(steps,
		Step{Command: "go south", Contains: []string{"Deep Canyon"}},
		Step{Command: "go southwest", Contains: []string{"Passage"}},
		Step{Command: "go south", Contains: []string{"Round Room"}},
		Step{Command: "go southeast", Contains: []string{"engravings Cave"}},
		Step{Command: "go east", Contains: []string{"Dome Room"}},
		Step{Command: "go down", Contains: []string{"torch Room"}},
		Step{Command: "take torch", Contains: []string{"Taken"}},
		Step{Command: "turn off lamp", Contains: []string{"now off"}},
	)

	// === Temple: bell, altar: candles+book, Hades ===
	steps = append(steps,
		Step{Command: "go south", Contains: []string{"Temple"}},
		Step{Command: "take bell", Contains: []string{"Taken"}},
		Step{Command: "go south", Contains: []string{"altar"}},
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
		Step{Command: "go west", Contains: []string{"slide Room"}},
		Step{Command: "go north", Contains: []string{"Mine Entrance"}},
		Step{Command: "go west", Contains: []string{"Squeaky Room"}},
	)

	// === bat Room, Shaft Room: deposit torch+screwdriver in basket ===
	steps = append(steps,
		Step{Command: "go north", Contains: []string{"bat Room"}},
		Step{Command: "go east", Contains: []string{"Shaft Room"}},
		Step{Command: "put torch in basket", Contains: []string{"Done"}},
		Step{Command: "put screwdriver in basket", Contains: []string{"Done"}},
		Step{Command: "turn on lamp", Contains: []string{"now on"}},
	)

	// === coal Mine: get coal ===
	steps = append(steps,
		Step{Command: "go north", Contains: []string{"Smelly Room"}},
		Step{Command: "go down", Contains: []string{"Gas Room"}},
		Step{Command: "go east", Contains: []string{"coal Mine"}},
		Step{Command: "go northeast", Contains: []string{"coal Mine"}},
		Step{Command: "go southeast", Contains: []string{"coal Mine"}},
		Step{Command: "go southwest", Contains: []string{"coal Mine"}},
		Step{Command: "go down", Contains: []string{"ladder Top"}},
		Step{Command: "go down", Contains: []string{"ladder Bottom"}},
		Step{Command: "go south", Contains: []string{"Dead End"}},
		Step{Command: "take coal", Contains: []string{"Taken"}},
	)

	// === Back to Shaft Room with coal ===
	steps = append(steps,
		Step{Command: "go north", Contains: []string{"ladder Bottom"}},
		Step{Command: "go up", Contains: []string{"ladder Top"}},
		Step{Command: "go up", Contains: []string{"coal Mine"}},
		Step{Command: "go north", Contains: []string{"coal Mine"}},
		Step{Command: "go east", Contains: []string{"coal Mine"}},
		Step{Command: "go south", Contains: []string{"coal Mine"}},
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
		Step{Command: "go east", Contains: []string{"coal Mine"}},
		Step{Command: "go northeast", Contains: []string{"coal Mine"}},
		Step{Command: "go southeast", Contains: []string{"coal Mine"}},
		Step{Command: "go southwest", Contains: []string{"coal Mine"}},
		Step{Command: "go down", Contains: []string{"ladder Top"}},
		Step{Command: "go down", Contains: []string{"ladder Bottom"}},
		Step{Command: "go west", Contains: []string{"Timber Room"}},
		Step{Command: "drop all"},
	)

	// === Drafty Room: machine Room diamond ===
	steps = append(steps,
		Step{Command: "go west", Contains: []string{"Drafty Room"}},
		Step{Command: "take coal", Contains: []string{"Taken"}},
		Step{Command: "take screwdriver", Contains: []string{"Taken"}},
		Step{Command: "take torch", Contains: []string{"Taken"}},
		Step{Command: "go south", Contains: []string{"machine Room"}},
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
		Step{Command: "go east", Contains: []string{"ladder Bottom"}},
		Step{Command: "go up", Contains: []string{"ladder Top"}},
		Step{Command: "go up", Contains: []string{"coal Mine"}},
		Step{Command: "go north", Contains: []string{"coal Mine"}},
		Step{Command: "go east", Contains: []string{"coal Mine"}},
		Step{Command: "go south", Contains: []string{"coal Mine"}},
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

	// === bat Room to Living Room via slide ===
	steps = append(steps,
		Step{Command: "go west", Contains: []string{"bat Room"}},
		Step{Command: "take jade", Contains: []string{"Taken"}},
		Step{Command: "go south", Contains: []string{"Squeaky Room"}},
		Step{Command: "go east", Contains: []string{"Mine Entrance"}},
		Step{Command: "go south", Contains: []string{"slide Room"}},
		Step{Command: "go down", Contains: []string{"cellar"}},
		Step{Command: "go up", Contains: []string{"Living Room"}},
	)

	// === Store jade and diamond ===
	steps = append(steps,
		Step{Command: "put jade in case", Contains: []string{"Done"}},
		Step{Command: "put diamond in case", Contains: []string{"Done"}},
	)

	// === reservoir: trunk and pump ===
	steps = append(steps,
		Step{Command: "turn on lamp", Contains: []string{"now on"}},
		Step{Command: "go down", Contains: []string{"cellar"}},
		Step{Command: "go north", Contains: []string{"troll Room"}},
		Step{Command: "go east", Contains: []string{"East-West Passage"}},
		Step{Command: "go north", Contains: []string{"Chasm"}},
		Step{Command: "go northeast", Contains: []string{"reservoir South"}},
		Step{Command: "go north", Contains: []string{"reservoir"}},
		Step{Command: "take trunk", Contains: []string{"Taken"}},
		Step{Command: "go north", Contains: []string{"reservoir North"}},
		Step{Command: "take pump", Contains: []string{"Taken"}},
	)

	// === Atlantis Room: trident ===
	steps = append(steps,
		Step{Command: "go north", Contains: []string{"Atlantis Room"}},
		Step{Command: "drop torch", Contains: []string{"Dropped"}},
		Step{Command: "take trident", Contains: []string{"Taken"}},
	)

	// === dam Base: inflate boat, river journey ===
	steps = append(steps,
		Step{Command: "go south", Contains: []string{"reservoir North"}},
		Step{Command: "go south", Contains: []string{"reservoir"}},
		Step{Command: "go south", Contains: []string{"reservoir South"}},
		Step{Command: "go east", Contains: []string{"dam"}},
		Step{Command: "go east", Contains: []string{"dam Base"}},
		Step{Command: "inflate plastic with pump", Contains: []string{"boat inflates"}},
		Step{Command: "drop pump", Contains: []string{"Dropped"}},
	)

	// === river (ZIL speeds: river1=4, river2=4, river3=3, river4=2, river5=1)
	// vWait runs Clocker 3x per call.
	// Launch: river1, tick=4. EOT: 4→3.
	// wait1 (3 ticks): 3→2→1→0 fire! → river2. Queue(4).
	// wait2 (3 ticks): 4→3→2→1.
	// wait3 (3 ticks): 1→0 fire! → river3. Queue(3). 3→2→1.
	// look (1 tick): 1→0 fire! → river4. Queue(2). Player sees river4!
	// take buoy (1 tick): 2→1. Still at river4!
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

	// === buoy emerald, cross rainbow ===
	steps = append(steps,
		Step{Command: "go southwest", Contains: []string{"Sandy Beach"}},
		Step{Command: "open buoy", Contains: []string{"emerald"}},
		Step{Command: "take emerald", Contains: []string{"Taken"}},
		Step{Command: "take garlic", Contains: []string{"Taken"}},
		Step{Command: "go south", Contains: []string{"shore"}},
		Step{Command: "go south", Contains: []string{"Aragain Falls"}},
		Step{Command: "cross rainbow", Contains: []string{"End of rainbow"}},
		Step{Command: "turn off lamp", Contains: []string{"now off"}},
	)

	// === Back to house, store treasures ===
	steps = append(steps,
		Step{Command: "go southwest", Contains: []string{"Canyon Bottom"}},
		Step{Command: "go up", Contains: []string{"Rocky Ledge"}},
		Step{Command: "go up", Contains: []string{"Canyon View"}},
		Step{Command: "go northwest", Contains: []string{"clearing"}},
		Step{Command: "go west", Contains: []string{"Behind House"}},
		Step{Command: "enter house", Contains: []string{"kitchen"}},
		Step{Command: "go west", Contains: []string{"Living Room"}},
		Step{Command: "put emerald in case", Contains: []string{"Done"}},
		Step{Command: "put scarab in case", Contains: []string{"Done"}},
		Step{Command: "put trident in case", Contains: []string{"Done"}},
		Step{Command: "put jewels in case", Contains: []string{"Done"}},
	)

	// === Get egg from tree ===
	steps = append(steps,
		Step{Command: "go east", Contains: []string{"kitchen"}},
		Step{Command: "go east", Contains: []string{"Behind House"}},
		Step{Command: "go north", Contains: []string{"North of House"}},
		Step{Command: "go north", Contains: []string{"forest path"}},
		Step{Command: "climb tree", Contains: []string{"Up a tree"}},
		Step{Command: "take egg", Contains: []string{"Taken"}},
		Step{Command: "go down", Contains: []string{"forest path"}},
		Step{Command: "go south", Contains: []string{"North of House"}},
		Step{Command: "go east", Contains: []string{"Behind House"}},
		Step{Command: "enter house", Contains: []string{"kitchen"}},
		Step{Command: "go west", Contains: []string{"Living Room"}},
	)

	// === Maze to cyclops Room ===
	steps = append(steps,
		Step{Command: "turn on lamp", Contains: []string{"now on"}},
		Step{Command: "go down", Contains: []string{"cellar"}},
		Step{Command: "go north", Contains: []string{"troll Room"}},
		Step{Command: "go west", Contains: []string{"Maze"}},
		Step{Command: "go south", Contains: []string{"Maze"}},
		Step{Command: "go east", Contains: []string{"Maze"}},
		Step{Command: "go up", Contains: []string{"Maze"}},
		Step{Command: "take coins", Contains: []string{"Taken"}},
		Step{Command: "take key", Contains: []string{"Taken"}},
		Step{Command: "go southwest", Contains: []string{"Maze"}},
		Step{Command: "go east", Contains: []string{"Maze"}},
		Step{Command: "go south", Contains: []string{"Maze"}},
		Step{Command: "go southeast", Contains: []string{"cyclops Room"}},
		Step{Command: "ulysses", Contains: []string{"cyclops"}},
	)

	// === thief: give egg, navigate, fight ===
	// After "ulysses", the cyclops flees. Go up to Treasure Room, give egg.
	steps = append(steps,
		Step{Command: "go up", Contains: []string{"Treasure Room"}},
		Step{Command: "give egg to thief", Contains: []string{"thief"}},
		Step{Command: "go down", Contains: []string{"cyclops Room"}},
		Step{Command: "go east", Contains: []string{"Strange Passage"}},
		Step{Command: "go east", Contains: []string{"Living Room"}},
		Step{Command: "put coins in case", Contains: []string{"Done"}},
		Step{Command: "take knife", Contains: []string{"Taken"}},
		Step{Command: "go west", Contains: []string{"Strange Passage"}},
		Step{Command: "go west", Contains: []string{"cyclops Room"}},
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
		Step{Command: "go down", Contains: []string{"cyclops Room"}},
		Step{Command: "go northwest", Contains: []string{"Maze"}},
		Step{Command: "go south", Contains: []string{"Maze"}},
		Step{Command: "go west", Contains: []string{"Maze"}},
		Step{Command: "go up", Contains: []string{"Maze"}},
		Step{Command: "go down", Contains: []string{"Maze"}},
		Step{Command: "go northeast", Contains: []string{"Grating Room"}},
		Step{Command: "unlock grate with key"},
		Step{Command: "open grate", Contains: []string{"grating opens"}},
		Step{Command: "go up", Contains: []string{"clearing"}},
	)

	// === canary bauble from tree ===
	steps = append(steps,
		Step{Command: "go south", Contains: []string{"forest path"}},
		Step{Command: "climb tree", Contains: []string{"Up a tree"}},
		Step{Command: "wind up canary", Contains: []string{"bauble"}},
		Step{Command: "go down", Contains: []string{"forest path"}},
		Step{Command: "drop knife"},
		Step{Command: "take bauble", Contains: []string{"Taken"}},
	)

	// === Final treasure deposits ===
	steps = append(steps,
		Step{Command: "go south", Contains: []string{"North of House"}},
		Step{Command: "go east", Contains: []string{"Behind House"}},
		Step{Command: "enter house", Contains: []string{"kitchen"}},
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
		Step{Command: "go down", Contains: []string{"cellar"}},
		Step{Command: "go north", Contains: []string{"troll Room"}},
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
		Step{Command: "go west", Contains: []string{"troll Room"}},
		Step{Command: "go south", Contains: []string{"cellar"}},
		Step{Command: "go up", Contains: []string{"Living Room"}},
		Step{Command: "put bar in case", Contains: []string{"Done"}},
	)

	runScript(t, steps)
}
