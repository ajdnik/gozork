package game

import . "github.com/ajdnik/gozork/engine"

func initClockFuncs() {
	G.ClockFuncs = map[string]func() bool{
		"IFight":        IFight,
		"ISword":        ISword,
		"IThief":        IThief,
		"ICandles":      ICandles,
		"ILantern":      ILantern,
		"ICure":         ICure,
		"ICyclops":      ICyclops,
		"IForestRandom": IForestRandom,
		"IMaintRoom":    IMaintRoom,
		"IMatch":        IMatch,
		"IRempty":       IRempty,
		"IRfill":        IRfill,
		"IRiver":        IRiver,
		"IXb":           IXb,
		"IXbh":          IXbh,
		"IXc":           IXc,
	}
}

func registerWellKnownObjects() {
	G.AllObjects = Objects
	G.RoomsObj = &Rooms
	G.GlobalObj = &GlobalObjects
	G.LocalGlobalObj = &LocalGlobals
	G.NotHereObj = &NotHereObject
	G.PseudoObj = &PseudoObject
	G.ItPronounObj = &It
	G.MeObj = &Me
	G.HandsObj = &Hands
	// Set the NotHereObject's Action so the engine can call it
	NotHereObject.Action = NotHereObjectFcn
}

// InitGame sets up all game state for a fresh game. Call once before MainLoop.
func InitGame() {
	if G == nil {
		G = NewGameState()
	}
	// Set up game-specific data
	G.GameData = NewZorkData()
	registerWellKnownObjects()
	G.ITakeFunc = ITake

	ResetGameState()
	initClockFuncs()
	FinalizeGameObjects()
	BuildObjectTree()
	BuildVocabulary(GameCommands, BuzzWords, Synonyms)
	InitReader()
	initSaveSystem()

	Queue("IFight", -1).Run = true
	Queue("ISword", -1)
	Queue("IThief", -1).Run = true
	Queue("ICandles", 40)
	Queue("ILantern", 200)
	InflatedBoat.VehType = FlgNonLand
	Def1Res[1] = Def1[2]
	Def1Res[2] = Def1[4]
	Def2Res[2] = Def2B[2]
	Def2Res[3] = Def2B[4]
	Def3Res[1] = Def3A[2]
	Def3Res[3] = Def3B[2]
	G.Here = &WestOfHouse
	ThisIsIt(&Mailbox)
	G.Lit = true
	G.Winner = &Adventurer
	G.Player = G.Winner
	G.Winner.MoveTo(G.Here)
}

func Run() {
	InitGame()
	if !G.Here.Has(FlgTouch) {
		VVersion(ActUnk)
		Printf("\n")
	}
	VLook(ActUnk)
	MainLoop()
}
