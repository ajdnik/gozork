package game

import . "github.com/ajdnik/gozork/engine"

var (

	// ================================================================
	// ROOMS - Maze
	// ================================================================

	Maze1 = Object{
		In:       &Rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    FlgLand | FlgMaze,
	}
	Maze2 = Object{
		In:       &Rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    FlgLand | FlgMaze,
	}
	Maze3 = Object{
		In:       &Rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    FlgLand | FlgMaze,
	}
	Maze4 = Object{
		In:       &Rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    FlgLand | FlgMaze,
	}
	DeadEnd1 = Object{
		In:       &Rooms,
		LongDesc: "You have come to a dead end in the maze.",
		Desc:     "Dead End",
		Flags:    FlgLand | FlgMaze,
	}
	Maze5 = Object{
		In:       &Rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike. A skeleton, probably the remains of a luckless adventurer, lies here.",
		Desc:     "Maze",
		Flags:    FlgLand | FlgMaze,
	}
	DeadEnd2 = Object{
		In:       &Rooms,
		LongDesc: "You have come to a dead end in the maze.",
		Desc:     "Dead End",
		Flags:    FlgLand | FlgMaze,
	}
	Maze6 = Object{
		In:       &Rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    FlgLand | FlgMaze,
	}
	Maze7 = Object{
		In:       &Rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    FlgLand | FlgMaze,
	}
	Maze8 = Object{
		In:       &Rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    FlgLand | FlgMaze,
	}
	DeadEnd3 = Object{
		In:       &Rooms,
		LongDesc: "You have come to a dead end in the maze.",
		Desc:     "Dead End",
		Flags:    FlgLand | FlgMaze,
	}
	Maze9 = Object{
		In:       &Rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    FlgLand | FlgMaze,
	}
	Maze10 = Object{
		In:       &Rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    FlgLand | FlgMaze,
	}
	Maze11 = Object{
		In:       &Rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    FlgLand | FlgMaze,
	}
	GratingRoom = Object{
		In:     &Rooms,
		Desc:   "Grating Room",
		// Action set in FinalizeGameObjects to avoid init cycle
		Flags:  FlgLand,
		Global: []*Object{&Grate},
	}
	Maze12 = Object{
		In:       &Rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    FlgLand | FlgMaze,
	}
	DeadEnd4 = Object{
		In:       &Rooms,
		LongDesc: "You have come to a dead end in the maze.",
		Desc:     "Dead End",
		Flags:    FlgLand | FlgMaze,
	}
	Maze13 = Object{
		In:       &Rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    FlgLand | FlgMaze,
	}
	Maze14 = Object{
		In:       &Rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    FlgLand | FlgMaze,
	}
	Maze15 = Object{
		In:       &Rooms,
		LongDesc: "This is part of a maze of twisty little passages, all alike.",
		Desc:     "Maze",
		Flags:    FlgLand | FlgMaze,
	}

)
