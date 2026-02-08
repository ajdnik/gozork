package game

import . "github.com/ajdnik/gozork/engine"

var (
	// ================================================================
	// ROOMS - forest and Outside
	// ================================================================

	westOfHouse = Object{
		In:   &rooms,
		Desc: "West of House",
		// Action set in finalizeGameObjects to avoid init cycle
		Flags:  FlgLand | FlgOn | FlgSacred,
		Global: []*Object{&whiteHouse, &board, &forest},
	}
	stoneBarrow = Object{
		In:       &rooms,
		LongDesc: "You are standing in front of a massive barrow of stone. In the east face is a huge stone door which is open. You cannot see into the dark of the tomb.",
		Desc:     "Stone barrow",
		Flags:    FlgLand | FlgOn | FlgSacred,
		// Action set in finalizeGameObjects to avoid init cycle
	}
	northOfHouse = Object{
		In:       &rooms,
		LongDesc: "You are facing the north side of a white house. There is no door here, and all the windows are boarded up. To the north a narrow path winds through the trees.",
		Desc:     "North of House",
		Flags:    FlgLand | FlgOn | FlgSacred,
		Global:   []*Object{&boardedWindow, &board, &whiteHouse, &forest},
	}
	southOfHouse = Object{
		In:       &rooms,
		LongDesc: "You are facing the south side of a white house. There is no door here, and all the windows are boarded.",
		Desc:     "South of House",
		Flags:    FlgLand | FlgOn | FlgSacred,
		Global:   []*Object{&boardedWindow, &board, &whiteHouse, &forest},
	}
	eastOfHouse = Object{
		In:   &rooms,
		Desc: "Behind House",
		// Action set in finalizeGameObjects to avoid init cycle
		Flags:  FlgLand | FlgOn | FlgSacred,
		Global: []*Object{&whiteHouse, &kitchenWindow, &forest},
	}
	forest1 = Object{
		In:       &rooms,
		LongDesc: "This is a forest, with trees in all directions. To the east, there appears to be sunlight.",
		Desc:     "forest",
		// Action set in finalizeGameObjects to avoid init cycle
		Flags:  FlgLand | FlgOn | FlgSacred,
		Global: []*Object{&tree, &songbird, &whiteHouse, &forest},
	}
	forest2 = Object{
		In:       &rooms,
		LongDesc: "This is a dimly lit forest, with large trees all around.",
		Desc:     "forest",
		// Action set in finalizeGameObjects to avoid init cycle
		Flags:  FlgLand | FlgOn | FlgSacred,
		Global: []*Object{&tree, &songbird, &whiteHouse, &forest},
	}
	mountains = Object{
		In:       &rooms,
		LongDesc: "The forest thins out, revealing impassable mountains.",
		Desc:     "forest",
		Flags:    FlgLand | FlgOn | FlgSacred,
		Global:   []*Object{&tree, &whiteHouse},
	}
	forest3 = Object{
		In:       &rooms,
		LongDesc: "This is a dimly lit forest, with large trees all around.",
		Desc:     "forest",
		// Action set in finalizeGameObjects to avoid init cycle
		Flags:  FlgLand | FlgOn | FlgSacred,
		Global: []*Object{&tree, &songbird, &whiteHouse, &forest},
	}
	path = Object{
		In:       &rooms,
		LongDesc: "This is a path winding through a dimly lit forest. The path heads north-south here. One particularly large tree with some low branches stands at the edge of the path.",
		Desc:     "forest path",
		// Action set in finalizeGameObjects to avoid init cycle
		Flags:  FlgLand | FlgOn | FlgSacred,
		Global: []*Object{&tree, &songbird, &whiteHouse, &forest},
	}
	upATree = Object{
		In:   &rooms,
		Desc: "Up a tree",
		// Action set in finalizeGameObjects to avoid init cycle
		Flags:  FlgLand | FlgOn | FlgSacred,
		Global: []*Object{&tree, &forest, &songbird, &whiteHouse},
	}
	gratingClearing = Object{
		In:     &rooms,
		Desc:   "clearing",
		Flags:  FlgLand | FlgOn | FlgSacred,
		Global: []*Object{&whiteHouse, &grate},
		// Action set in finalizeGameObjects to avoid init cycle
	}
	clearing = Object{
		In:       &rooms,
		LongDesc: "You are in a small clearing in a well marked forest path that extends to the east and west.",
		Desc:     "clearing",
		Flags:    FlgLand | FlgOn | FlgSacred,
		Global:   []*Object{&tree, &songbird, &whiteHouse, &forest},
		// Action set in finalizeGameObjects to avoid init cycle
	}

	// ================================================================
	// ROOMS - House
	// ================================================================

	kitchen = Object{
		In:     &rooms,
		Desc:   "kitchen",
		Flags:  FlgLand | FlgOn | FlgSacred,
		Global: []*Object{&kitchenWindow, &chimney, &stairs},
		// Action set in finalizeGameObjects to avoid init cycle
		Item: &ItemData{Value: 10},
	}
	attic = Object{
		In:       &rooms,
		LongDesc: "This is the attic. The only exit is a stairway leading down.",
		Desc:     "attic",
		Flags:    FlgLand | FlgSacred,
		Global:   []*Object{&stairs},
	}
	livingRoom = Object{
		In:     &rooms,
		Desc:   "Living Room",
		Flags:  FlgLand | FlgOn | FlgSacred,
		Global: []*Object{&stairs},
		Pseudo: []PseudoObj{
			{Synonym: "nails", Action: nailsPseudo},
			{Synonym: "nail", Action: nailsPseudo},
		},
		// Action set in finalizeGameObjects to avoid init cycle
	}

	// ================================================================
	// ROOMS - cellar and Vicinity
	// ================================================================

	cellar = Object{
		In:   &rooms,
		Desc: "cellar",
		// Action set in finalizeGameObjects to avoid init cycle
		Flags:  FlgLand,
		Global: []*Object{&slide, &stairs},
		Item:   &ItemData{Value: 25},
	}
	trollRoom = Object{
		In:       &rooms,
		LongDesc: "This is a small room with passages to the east and south and a forbidding hole leading west. Bloodstains and deep scratches (perhaps made by an axe) mar the walls.",
		Desc:     "The troll Room",
		Flags:    FlgLand,
		// Action set in finalizeGameObjects to avoid init cycle
	}
	eastOfChasm = Object{
		In:       &rooms,
		LongDesc: "You are on the east edge of a chasm, the bottom of which cannot be seen. A narrow passage goes north, and the path you are on continues to the east.",
		Desc:     "East of Chasm",
		Flags:    FlgLand,
		Pseudo:   []PseudoObj{{Synonym: "chasm", Action: chasmPseudo}},
	}
	gallery = Object{
		In:       &rooms,
		LongDesc: "This is an art gallery. Most of the paintings have been stolen by vandals with exceptional taste. The vandals left through either the north or west exits.",
		Desc:     "gallery",
		Flags:    FlgLand | FlgOn,
	}
	studio = Object{
		In:       &rooms,
		LongDesc: "This appears to have been an artist's studio. The walls and floors are splattered with paints of 69 different colors. Strangely enough, nothing of value is hanging here. At the south end of the room is an open door (also covered with paint). A dark and narrow chimney leads up from a fireplace; although you might be able to get up it, it seems unlikely you could get back down.",
		Desc:     "studio",
		Flags:    FlgLand,
		Global:   []*Object{&chimney},
		Pseudo: []PseudoObj{
			{Synonym: "door", Action: doorPseudo},
			{Synonym: "paint", Action: paintPseudo},
		},
	}
)
