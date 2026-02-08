package game

import . "github.com/ajdnik/gozork/engine"

func initClockFuncs() {
	G.ClockFuncs = map[string]func() bool{
		"iFight":        iFight,
		"iSword":        iSword,
		"iThief":        iThief,
		"iCandles":      iCandles,
		"iLantern":      iLantern,
		"iCure":         iCure,
		"iCyclops":      iCyclops,
		"iForestRandom": iForestRandom,
		"iMaintRoom":    iMaintRoom,
		"iMatch":        iMatch,
		"iRempty":       iRempty,
		"iRfill":        iRfill,
		"iRiver":        iRiver,
		"iXb":           iXb,
		"iXbh":          iXbh,
		"iXc":           iXc,
	}
}

func registerWellKnownObjects() {
	G.AllObjects = objects
	G.RoomsObj = &rooms
	G.GlobalObj = &globalObjects
	G.LocalGlobalObj = &localGlobals
	G.NotHereObj = &notHereObject
	G.PseudoObj = &pseudoObject
	G.ItPronounObj = &it
	G.MeObj = &me
	G.HandsObj = &hands
	// Set the notHereObject's Action so the engine can call it
	notHereObject.Action = notHereObjectFcn
}

// InitGame sets up all game state for a fresh game. Call once before MainLoop.
func InitGame() {
	if G == nil {
		G = NewGameState()
	}
	// Set up game-specific data
	G.GameData = newZorkData()
	registerWellKnownObjects()
	G.ITakeFunc = iTake

	ResetGameState()
	initClockFuncs()
	finalizeGameObjects()
	BuildObjectTree()
	BuildVocabulary(gameCommands, buzzWords, synonyms)
	InitReader()
	initSaveSystem()

	Queue("iFight", -1).Run = true
	Queue("iSword", -1)
	Queue("iThief", -1).Run = true
	Queue("iCandles", 40)
	Queue("iLantern", 200)
	inflatedBoat.SetVehType(FlgNonLand)
	def1Res[1] = def1[2]
	def1Res[2] = def1[4]
	def2Res[2] = def2B[2]
	def2Res[3] = def2B[4]
	def3Res[1] = def3A[2]
	def3Res[3] = def3B[2]
	G.Here = &westOfHouse
	thisIsIt(&mailbox)
	G.Lit = true
	G.Winner = &adventurer
	G.Player = G.Winner
	G.Winner.MoveTo(G.Here)
}

func Run() {
	InitGame()
	if !G.Here.Has(FlgTouch) {
		vVersion(ActUnk)
		Printf("\n")
	}
	vLook(ActUnk)
	MainLoop()
}
