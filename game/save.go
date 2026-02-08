package game

import (
	. "github.com/ajdnik/gozork/engine"
	"encoding/gob"
	"fmt"
	"os"
	"strings"
)

var SaveFile = "zork.sav"

type savedObject struct {
	InIdx     int
	Flags     Flags
	Strength  int
	Value     int
	TValue    int
	Text      string
	LongDesc  string
	FirstDesc string
}

type savedQueueItem struct {
	Key  string
	Run  bool
	Tick int
}

type savedVillain struct {
	Prob int
}

type gameState struct {
	ObjStates []savedObject

	// Engine state
	Dead       bool
	Deaths     int
	Score      int
	BaseScore  int
	Moves      int
	Lit        bool
	SuperBrief bool
	Verbose    bool
	WonGame    bool
	HelloSailor int
	IsSprayed  bool
	Lucky      bool
	FumbleNumber int
	FumbleProb int

	// Dungeon flags
	TrollFlag         bool
	CyclopsFlag       bool
	MagicFlag         bool
	LowTide           bool
	DomeFlag          bool
	EmptyHanded       bool
	LLDFlag           bool
	RainbowFlag       bool
	DeflateFlag       bool
	CoffinCure        bool
	GrateRevealed     bool
	KitchenWindowFlag bool
	CageTop           bool
	RugMoved          bool
	GrUnlock          bool
	CycloWrath        int
	MirrorMung        bool
	GateFlag          bool
	GatesOpen         bool
	WaterLevel        int
	MatchCount        int
	EggSolve          bool
	ThiefHere         bool
	ThiefEngrossed    bool
	LoudFlag          bool
	SingSong          bool
	BuoyFlag          bool
	BeachDig          int
	LightShaft        int
	XB                bool
	XC                bool
	Deflate           bool
	LampTableIdx      int
	CandleTableIdx    int
	LoadAllowed       int
	LoadMax           int

	HereIdx    int
	WinnerIdx  int
	PlayerIdx  int
	ItObjIdx   int
	DescObjIdx int

	QueueItems [30]savedQueueItem
	QueueInts  int
	QueueDmns  int
	ClockWait  bool

	VillainProbs []savedVillain
}

func captureState() *gameState {
	BuildObjIndex()
	gd := GD()
	s := &gameState{}

	s.ObjStates = make([]savedObject, len(G.AllObjects))
	for i, obj := range G.AllObjects {
		s.ObjStates[i] = savedObject{
			InIdx:     ObjToIdx(obj.In),
			Flags:     obj.Flags,
			Strength:  obj.Strength,
			Value:     obj.Value,
			TValue:    obj.TValue,
			Text:      obj.Text,
			LongDesc:  obj.LongDesc,
			FirstDesc: obj.FirstDesc,
		}
	}

	s.Dead = gd.Dead
	s.Deaths = gd.Deaths
	s.Score = G.Score
	s.BaseScore = G.BaseScore
	s.Moves = G.Moves
	s.Lit = G.Lit
	s.SuperBrief = G.SuperBrief
	s.Verbose = G.Verbose
	s.WonGame = gd.WonGame
	s.HelloSailor = gd.HelloSailor
	s.IsSprayed = gd.IsSprayed
	s.Lucky = G.Lucky
	s.FumbleNumber = gd.FumbleNumber
	s.FumbleProb = gd.FumbleProb

	s.TrollFlag = gd.TrollFlag
	s.CyclopsFlag = gd.CyclopsFlag
	s.MagicFlag = gd.MagicFlag
	s.LowTide = gd.LowTide
	s.DomeFlag = gd.DomeFlag
	s.EmptyHanded = gd.EmptyHanded
	s.LLDFlag = gd.LLDFlag
	s.RainbowFlag = gd.RainbowFlag
	s.DeflateFlag = gd.DeflateFlag
	s.CoffinCure = gd.CoffinCure
	s.GrateRevealed = gd.GrateRevealed
	s.KitchenWindowFlag = gd.KitchenWindowFlag
	s.CageTop = gd.CageTop
	s.RugMoved = gd.RugMoved
	s.GrUnlock = gd.GrUnlock
	s.CycloWrath = gd.CycloWrath
	s.MirrorMung = gd.MirrorMung
	s.GateFlag = gd.GateFlag
	s.GatesOpen = gd.GatesOpen
	s.WaterLevel = gd.WaterLevel
	s.MatchCount = gd.MatchCount
	s.EggSolve = gd.EggSolve
	s.ThiefHere = gd.ThiefHere
	s.ThiefEngrossed = gd.ThiefEngrossed
	s.LoudFlag = gd.LoudFlag
	s.SingSong = gd.SingSong
	s.BuoyFlag = gd.BuoyFlag
	s.BeachDig = gd.BeachDig
	s.LightShaft = gd.LightShaft
	s.XB = gd.XB
	s.XC = gd.XC
	s.Deflate = gd.Deflate
	s.LampTableIdx = gd.LampTableIdx
	s.CandleTableIdx = gd.CandleTableIdx
	s.LoadAllowed = gd.LoadAllowed
	s.LoadMax = gd.LoadMax

	s.HereIdx = ObjToIdx(G.Here)
	s.WinnerIdx = ObjToIdx(G.Winner)
	s.PlayerIdx = ObjToIdx(G.Player)
	s.ItObjIdx = ObjToIdx(G.Params.ItObj)
	s.DescObjIdx = ObjToIdx(gd.DescObj)

	for i := range G.QueueItms {
		s.QueueItems[i] = savedQueueItem{
			Key:  G.QueueItms[i].Key,
			Run:  G.QueueItms[i].Run,
			Tick: G.QueueItms[i].Tick,
		}
	}
	s.QueueInts = G.QueueInts
	s.QueueDmns = G.QueueDmns
	s.ClockWait = G.ClockWait

	s.VillainProbs = make([]savedVillain, len(gd.Villains))
	for i, v := range gd.Villains {
		s.VillainProbs[i] = savedVillain{Prob: v.Prob}
	}

	return s
}

func applyState(s *gameState) {
	BuildObjIndex()
	gd := GD()

	for _, obj := range G.AllObjects {
		obj.Children = nil
	}
	for i, obj := range G.AllObjects {
		so := s.ObjStates[i]
		obj.In = IdxToObj(so.InIdx)
		obj.Flags = so.Flags
		obj.Strength = so.Strength
		obj.Value = so.Value
		obj.TValue = so.TValue
		obj.Text = so.Text
		obj.LongDesc = so.LongDesc
		obj.FirstDesc = so.FirstDesc
	}
	for _, obj := range G.AllObjects {
		if obj.In != nil {
			obj.In.AddChild(obj)
		}
	}

	gd.Dead = s.Dead
	gd.Deaths = s.Deaths
	G.Score = s.Score
	G.BaseScore = s.BaseScore
	G.Moves = s.Moves
	G.Lit = s.Lit
	G.SuperBrief = s.SuperBrief
	G.Verbose = s.Verbose
	gd.WonGame = s.WonGame
	gd.HelloSailor = s.HelloSailor
	gd.IsSprayed = s.IsSprayed
	G.Lucky = s.Lucky
	gd.FumbleNumber = s.FumbleNumber
	gd.FumbleProb = s.FumbleProb

	gd.TrollFlag = s.TrollFlag
	gd.CyclopsFlag = s.CyclopsFlag
	gd.MagicFlag = s.MagicFlag
	gd.LowTide = s.LowTide
	gd.DomeFlag = s.DomeFlag
	gd.EmptyHanded = s.EmptyHanded
	gd.LLDFlag = s.LLDFlag
	gd.RainbowFlag = s.RainbowFlag
	gd.DeflateFlag = s.DeflateFlag
	gd.CoffinCure = s.CoffinCure
	gd.GrateRevealed = s.GrateRevealed
	gd.KitchenWindowFlag = s.KitchenWindowFlag
	gd.CageTop = s.CageTop
	gd.RugMoved = s.RugMoved
	gd.GrUnlock = s.GrUnlock
	gd.CycloWrath = s.CycloWrath
	gd.MirrorMung = s.MirrorMung
	gd.GateFlag = s.GateFlag
	gd.GatesOpen = s.GatesOpen
	gd.WaterLevel = s.WaterLevel
	gd.MatchCount = s.MatchCount
	gd.EggSolve = s.EggSolve
	gd.ThiefHere = s.ThiefHere
	gd.ThiefEngrossed = s.ThiefEngrossed
	gd.LoudFlag = s.LoudFlag
	gd.SingSong = s.SingSong
	gd.BuoyFlag = s.BuoyFlag
	gd.BeachDig = s.BeachDig
	gd.LightShaft = s.LightShaft
	gd.XB = s.XB
	gd.XC = s.XC
	gd.Deflate = s.Deflate
	gd.LampTableIdx = s.LampTableIdx
	gd.CandleTableIdx = s.CandleTableIdx
	gd.LoadAllowed = s.LoadAllowed
	gd.LoadMax = s.LoadMax

	G.Here = IdxToObj(s.HereIdx)
	G.Winner = IdxToObj(s.WinnerIdx)
	G.Player = IdxToObj(s.PlayerIdx)
	G.Params.ItObj = IdxToObj(s.ItObjIdx)
	gd.DescObj = IdxToObj(s.DescObjIdx)

	for i := range G.QueueItms {
		G.QueueItms[i].Key = s.QueueItems[i].Key
		G.QueueItms[i].Run = s.QueueItems[i].Run
		G.QueueItms[i].Tick = s.QueueItems[i].Tick
		if s.QueueItems[i].Key != "" {
			G.QueueItms[i].Fn = G.ClockFuncs[s.QueueItems[i].Key]
		} else {
			G.QueueItms[i].Fn = nil
		}
	}
	G.QueueInts = s.QueueInts
	G.QueueDmns = s.QueueDmns
	G.ClockWait = s.ClockWait

	for i := range gd.Villains {
		if i < len(s.VillainProbs) {
			gd.Villains[i].Prob = s.VillainProbs[i].Prob
		}
	}
}

func promptFilename(action string) string {
	fmt.Fprintf(G.GameOutput, "Enter a file name (default is \"%s\"): ", SaveFile)
	if G.Reader == nil {
		InitReader()
	}
	txt, err := G.Reader.ReadString('\n')
	if err != nil && len(txt) == 0 {
		return SaveFile
	}
	txt = strings.TrimSpace(txt)
	if txt == "" {
		return SaveFile
	}
	return txt
}

func initSaveSystem() {
	BuildObjIndex()
	G.Save = doSave
	G.Restore = doRestore
	G.Restart = doRestart
}

func doSave() bool {
	fname := promptFilename("save")
	s := captureState()

	f, err := os.Create(fname)
	if err != nil {
		return false
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	if err := enc.Encode(s); err != nil {
		return false
	}
	return true
}

func doRestore() bool {
	fname := promptFilename("restore")

	f, err := os.Open(fname)
	if err != nil {
		return false
	}
	defer f.Close()

	var s gameState
	dec := gob.NewDecoder(f)
	if err := dec.Decode(&s); err != nil {
		return false
	}

	if len(s.ObjStates) != len(G.AllObjects) {
		return false
	}

	applyState(&s)
	G.Lit = IsLit(G.Here, true)
	return true
}

func doRestart() bool {
	G.GameData = NewZorkData()
	registerWellKnownObjects()
	G.ITakeFunc = ITake
	ResetGameState()
	FinalizeGameObjects()
	BuildObjectTree()

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

	BuildObjIndex()
	G.Save = doSave
	G.Restore = doRestore
	G.Restart = doRestart
	return true
}
