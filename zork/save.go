package zork

import (
	"encoding/gob"
	"fmt"
	"os"
	"reflect"
	"strings"
)

// SaveFile is the default path used for save/restore.
var SaveFile = "zork.sav"

// ---- object index mapping ----

// objIndex maps each *Object in Objects to its index. Built lazily.
var objIndex map[*Object]int

func buildObjIndex() {
	if objIndex != nil {
		return
	}
	objIndex = make(map[*Object]int, len(Objects)+1)
	for i, o := range Objects {
		objIndex[o] = i
	}
}

func objToIdx(o *Object) int {
	if o == nil {
		return -1
	}
	if idx, ok := objIndex[o]; ok {
		return idx
	}
	return -1
}

func idxToObj(idx int) *Object {
	if idx < 0 || idx >= len(Objects) {
		return nil
	}
	return Objects[idx]
}

// ---- clock function registry ----

var (
	rtnToName map[uintptr]string
	nameToRtn map[string]RtnFunc
)

func buildRtnRegistry() {
	if rtnToName != nil {
		return
	}
	entries := []struct {
		name string
		fn   RtnFunc
	}{
		{"IFight", IFight},
		{"ISword", ISword},
		{"IThief", IThief},
		{"ICandles", ICandles},
		{"ILantern", ILantern},
		{"ICure", ICure},
		{"ICyclops", ICyclops},
		{"IForestRandom", IForestRandom},
		{"IMaintRoom", IMaintRoom},
		{"IMatch", IMatch},
		{"IRempty", IRempty},
		{"IRfill", IRfill},
		{"IRiver", IRiver},
		{"IXb", IXb},
		{"IXbh", IXbh},
		{"IXc", IXc},
	}
	rtnToName = make(map[uintptr]string, len(entries))
	nameToRtn = make(map[string]RtnFunc, len(entries))
	for _, e := range entries {
		ptr := reflect.ValueOf(e.fn).Pointer()
		rtnToName[ptr] = e.name
		nameToRtn[e.name] = e.fn
	}
}

func rtnName(fn RtnFunc) string {
	if fn == nil {
		return ""
	}
	return rtnToName[reflect.ValueOf(fn).Pointer()]
}

// ---- serializable state structs ----

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
	Run     bool
	Tick    int
	RtnName string
}

type savedVillain struct {
	Prob int
}

type gameState struct {
	// Object states (indexed by Objects slice position)
	ObjStates []savedObject

	// Global game flags
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

	// Key references (object indices)
	HereIdx    int
	WinnerIdx  int
	PlayerIdx  int
	ItObjIdx   int
	DescObjIdx int

	// Clock state
	QueueItems [30]savedQueueItem
	QueueInts  int
	QueueDmns  int
	ClockWait  bool

	// Villain wake probabilities
	VillainProbs []savedVillain
}

func captureState() *gameState {
	buildObjIndex()
	buildRtnRegistry()

	s := &gameState{}

	// Capture object states
	s.ObjStates = make([]savedObject, len(Objects))
	for i, obj := range Objects {
		s.ObjStates[i] = savedObject{
			InIdx:     objToIdx(obj.In),
			Flags:     obj.Flags,
			Strength:  obj.Strength,
			Value:     obj.Value,
			TValue:    obj.TValue,
			Text:      obj.Text,
			LongDesc:  obj.LongDesc,
			FirstDesc: obj.FirstDesc,
		}
	}

	// Capture global game flags
	s.Dead = G.Dead
	s.Deaths = G.Deaths
	s.Score = G.Score
	s.BaseScore = G.BaseScore
	s.Moves = G.Moves
	s.Lit = G.Lit
	s.SuperBrief = G.SuperBrief
	s.Verbose = G.Verbose
	s.WonGame = G.WonGame
	s.HelloSailor = G.HelloSailor
	s.IsSprayed = G.IsSprayed
	s.Lucky = G.Lucky
	s.FumbleNumber = G.FumbleNumber
	s.FumbleProb = G.FumbleProb

	// Capture dungeon flags
	s.TrollFlag = G.TrollFlag
	s.CyclopsFlag = G.CyclopsFlag
	s.MagicFlag = G.MagicFlag
	s.LowTide = G.LowTide
	s.DomeFlag = G.DomeFlag
	s.EmptyHanded = G.EmptyHanded
	s.LLDFlag = G.LLDFlag
	s.RainbowFlag = G.RainbowFlag
	s.DeflateFlag = G.DeflateFlag
	s.CoffinCure = G.CoffinCure
	s.GrateRevealed = G.GrateRevealed
	s.KitchenWindowFlag = G.KitchenWindowFlag
	s.CageTop = G.CageTop
	s.RugMoved = G.RugMoved
	s.GrUnlock = G.GrUnlock
	s.CycloWrath = G.CycloWrath
	s.MirrorMung = G.MirrorMung
	s.GateFlag = G.GateFlag
	s.GatesOpen = G.GatesOpen
	s.WaterLevel = G.WaterLevel
	s.MatchCount = G.MatchCount
	s.EggSolve = G.EggSolve
	s.ThiefHere = G.ThiefHere
	s.ThiefEngrossed = G.ThiefEngrossed
	s.LoudFlag = G.LoudFlag
	s.SingSong = G.SingSong
	s.BuoyFlag = G.BuoyFlag
	s.BeachDig = G.BeachDig
	s.LightShaft = G.LightShaft
	s.XB = G.XB
	s.XC = G.XC
	s.Deflate = G.Deflate
	s.LampTableIdx = G.LampTableIdx
	s.CandleTableIdx = G.CandleTableIdx
	s.LoadAllowed = G.LoadAllowed
	s.LoadMax = G.LoadMax

	// Capture key references
	s.HereIdx = objToIdx(G.Here)
	s.WinnerIdx = objToIdx(G.Winner)
	s.PlayerIdx = objToIdx(G.Player)
	s.ItObjIdx = objToIdx(G.Params.ItObj)
	s.DescObjIdx = objToIdx(G.DescObj)

	// Capture clock state
	for i := range G.QueueItms {
		s.QueueItems[i] = savedQueueItem{
			Run:     G.QueueItms[i].Run,
			Tick:    G.QueueItms[i].Tick,
			RtnName: rtnName(G.QueueItms[i].Rtn),
		}
	}
	s.QueueInts = G.QueueInts
	s.QueueDmns = G.QueueDmns
	s.ClockWait = G.ClockWait

	// Capture villain probabilities
	s.VillainProbs = make([]savedVillain, len(G.Villains))
	for i, v := range G.Villains {
		s.VillainProbs[i] = savedVillain{Prob: v.Prob}
	}

	return s
}

func applyState(s *gameState) {
	buildObjIndex()
	buildRtnRegistry()

	// Restore object states
	// First clear all children
	for _, obj := range Objects {
		obj.Children = nil
	}
	for i, obj := range Objects {
		so := s.ObjStates[i]
		obj.In = idxToObj(so.InIdx)
		obj.Flags = so.Flags
		obj.Strength = so.Strength
		obj.Value = so.Value
		obj.TValue = so.TValue
		obj.Text = so.Text
		obj.LongDesc = so.LongDesc
		obj.FirstDesc = so.FirstDesc
	}
	// Rebuild children from restored In pointers
	for _, obj := range Objects {
		if obj.In != nil {
			obj.In.AddChild(obj)
		}
	}

	// Restore global game flags
	G.Dead = s.Dead
	G.Deaths = s.Deaths
	G.Score = s.Score
	G.BaseScore = s.BaseScore
	G.Moves = s.Moves
	G.Lit = s.Lit
	G.SuperBrief = s.SuperBrief
	G.Verbose = s.Verbose
	G.WonGame = s.WonGame
	G.HelloSailor = s.HelloSailor
	G.IsSprayed = s.IsSprayed
	G.Lucky = s.Lucky
	G.FumbleNumber = s.FumbleNumber
	G.FumbleProb = s.FumbleProb

	// Restore dungeon flags
	G.TrollFlag = s.TrollFlag
	G.CyclopsFlag = s.CyclopsFlag
	G.MagicFlag = s.MagicFlag
	G.LowTide = s.LowTide
	G.DomeFlag = s.DomeFlag
	G.EmptyHanded = s.EmptyHanded
	G.LLDFlag = s.LLDFlag
	G.RainbowFlag = s.RainbowFlag
	G.DeflateFlag = s.DeflateFlag
	G.CoffinCure = s.CoffinCure
	G.GrateRevealed = s.GrateRevealed
	G.KitchenWindowFlag = s.KitchenWindowFlag
	G.CageTop = s.CageTop
	G.RugMoved = s.RugMoved
	G.GrUnlock = s.GrUnlock
	G.CycloWrath = s.CycloWrath
	G.MirrorMung = s.MirrorMung
	G.GateFlag = s.GateFlag
	G.GatesOpen = s.GatesOpen
	G.WaterLevel = s.WaterLevel
	G.MatchCount = s.MatchCount
	G.EggSolve = s.EggSolve
	G.ThiefHere = s.ThiefHere
	G.ThiefEngrossed = s.ThiefEngrossed
	G.LoudFlag = s.LoudFlag
	G.SingSong = s.SingSong
	G.BuoyFlag = s.BuoyFlag
	G.BeachDig = s.BeachDig
	G.LightShaft = s.LightShaft
	G.XB = s.XB
	G.XC = s.XC
	G.Deflate = s.Deflate
	G.LampTableIdx = s.LampTableIdx
	G.CandleTableIdx = s.CandleTableIdx
	G.LoadAllowed = s.LoadAllowed
	G.LoadMax = s.LoadMax

	// Restore key references
	G.Here = idxToObj(s.HereIdx)
	G.Winner = idxToObj(s.WinnerIdx)
	G.Player = idxToObj(s.PlayerIdx)
	G.Params.ItObj = idxToObj(s.ItObjIdx)
	G.DescObj = idxToObj(s.DescObjIdx)

	// Restore clock state
	for i := range G.QueueItms {
		G.QueueItms[i].Run = s.QueueItems[i].Run
		G.QueueItms[i].Tick = s.QueueItems[i].Tick
		if s.QueueItems[i].RtnName != "" {
			G.QueueItms[i].Rtn = nameToRtn[s.QueueItems[i].RtnName]
		} else {
			G.QueueItms[i].Rtn = nil
		}
	}
	G.QueueInts = s.QueueInts
	G.QueueDmns = s.QueueDmns
	G.ClockWait = s.ClockWait

	// Restore villain probabilities
	for i := range G.Villains {
		if i < len(s.VillainProbs) {
			G.Villains[i].Prob = s.VillainProbs[i].Prob
		}
	}
}

// promptFilename prompts the user for a save filename and returns it.
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

// initSaveSystem wires up the real Save/Restore/Restart implementations.
// Called from InitGame after all package-level variables are initialized,
// breaking the init cycle between Objects and function references.
func initSaveSystem() {
	buildObjIndex()
	buildRtnRegistry()
	G.Save = doSave
	G.Restore = doRestore
	G.Restart = doRestart
}

// doSave serializes the current game state to a file.
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

// doRestore loads game state from a file, replacing all current state.
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

	if len(s.ObjStates) != len(Objects) {
		return false
	}

	applyState(&s)
	// Recalculate lighting after restoring state
	G.Lit = IsLit(G.Here, true)
	return true
}

// doRestart reinitializes the game to its starting state.
func doRestart() bool {
	ResetGameState()
	FinalizeGameObjects()
	BuildObjectTree()

	Queue(IFight, -1).Run = true
	Queue(ISword, -1)
	Queue(IThief, -1).Run = true
	Queue(ICandles, 40)
	Queue(ILantern, 200)
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

	// Re-register save system for the fresh game
	buildObjIndex()
	G.Save = doSave
	G.Restore = doRestore
	G.Restart = doRestart
	return true
}
