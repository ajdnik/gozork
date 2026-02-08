package game

import . "github.com/ajdnik/gozork/engine"

var (

	// ================================================================
	// ROOMS - Maze
	// ================================================================

	maze1 = Object{
		In:       &rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    FlgLand | FlgMaze,
	}
	maze2 = Object{
		In:       &rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    FlgLand | FlgMaze,
	}
	maze3 = Object{
		In:       &rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    FlgLand | FlgMaze,
	}
	maze4 = Object{
		In:       &rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    FlgLand | FlgMaze,
	}
	deadEnd1 = Object{
		In:       &rooms,
		LongDesc: "You have come to a dead end in the maze.",
		Desc:     "Dead End",
		Flags:    FlgLand | FlgMaze,
	}
	maze5 = Object{
		In:       &rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike. A skeleton, probably the remains of a luckless adventurer, lies here.",
		Desc:     "Maze",
		Flags:    FlgLand | FlgMaze,
	}
	deadEnd2 = Object{
		In:       &rooms,
		LongDesc: "You have come to a dead end in the maze.",
		Desc:     "Dead End",
		Flags:    FlgLand | FlgMaze,
	}
	maze6 = Object{
		In:       &rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    FlgLand | FlgMaze,
	}
	maze7 = Object{
		In:       &rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    FlgLand | FlgMaze,
	}
	maze8 = Object{
		In:       &rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    FlgLand | FlgMaze,
	}
	deadEnd3 = Object{
		In:       &rooms,
		LongDesc: "You have come to a dead end in the maze.",
		Desc:     "Dead End",
		Flags:    FlgLand | FlgMaze,
	}
	maze9 = Object{
		In:       &rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    FlgLand | FlgMaze,
	}
	maze10 = Object{
		In:       &rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    FlgLand | FlgMaze,
	}
	maze11 = Object{
		In:       &rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    FlgLand | FlgMaze,
	}
	gratingRoom = Object{
		In:   &rooms,
		Desc: "Grating Room",
		// Action set in finalizeGameObjects to avoid init cycle
		Flags:  FlgLand,
		Global: []*Object{&grate},
	}
	maze12 = Object{
		In:       &rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    FlgLand | FlgMaze,
	}
	deadEnd4 = Object{
		In:       &rooms,
		LongDesc: "You have come to a dead end in the maze.",
		Desc:     "Dead End",
		Flags:    FlgLand | FlgMaze,
	}
	maze13 = Object{
		In:       &rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    FlgLand | FlgMaze,
	}
	maze14 = Object{
		In:       &rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    FlgLand | FlgMaze,
	}
	maze15 = Object{
		In:       &rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    FlgLand | FlgMaze,
	}
)
