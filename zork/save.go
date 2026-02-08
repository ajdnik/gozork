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
	s.Dead = Dead
	s.Deaths = Deaths
	s.Score = Score
	s.BaseScore = BaseScore
	s.Moves = Moves
	s.Lit = Lit
	s.SuperBrief = SuperBrief
	s.Verbose = Verbose
	s.WonGame = WonGame
	s.HelloSailor = HelloSailor
	s.IsSprayed = IsSprayed
	s.Lucky = Lucky
	s.FumbleNumber = FumbleNumber
	s.FumbleProb = FumbleProb

	// Capture dungeon flags
	s.TrollFlag = TrollFlag
	s.CyclopsFlag = CyclopsFlag
	s.MagicFlag = MagicFlag
	s.LowTide = LowTide
	s.DomeFlag = DomeFlag
	s.EmptyHanded = EmptyHanded
	s.LLDFlag = LLDFlag
	s.RainbowFlag = RainbowFlag
	s.DeflateFlag = DeflateFlag
	s.CoffinCure = CoffinCure
	s.GrateRevealed = GrateRevealed
	s.KitchenWindowFlag = KitchenWindowFlag
	s.CageTop = CageTop
	s.RugMoved = RugMoved
	s.GrUnlock = GrUnlock
	s.CycloWrath = CycloWrath
	s.MirrorMung = MirrorMung
	s.GateFlag = GateFlag
	s.GatesOpen = GatesOpen
	s.WaterLevel = WaterLevel
	s.MatchCount = MatchCount
	s.EggSolve = EggSolve
	s.ThiefHere = ThiefHere
	s.ThiefEngrossed = ThiefEngrossed
	s.LoudFlag = LoudFlag
	s.SingSong = SingSong
	s.BuoyFlag = BuoyFlag
	s.BeachDig = BeachDig
	s.LightShaft = LightShaft
	s.XB = XB
	s.XC = XC
	s.Deflate = Deflate
	s.LampTableIdx = LampTableIdx
	s.CandleTableIdx = CandleTableIdx
	s.LoadAllowed = LoadAllowed
	s.LoadMax = LoadMax

	// Capture key references
	s.HereIdx = objToIdx(Here)
	s.WinnerIdx = objToIdx(Winner)
	s.PlayerIdx = objToIdx(Player)
	s.ItObjIdx = objToIdx(Params.ItObj)
	s.DescObjIdx = objToIdx(DescObj)

	// Capture clock state
	for i := range QueueItms {
		s.QueueItems[i] = savedQueueItem{
			Run:     QueueItms[i].Run,
			Tick:    QueueItms[i].Tick,
			RtnName: rtnName(QueueItms[i].Rtn),
		}
	}
	s.QueueInts = QueueInts
	s.QueueDmns = QueueDmns
	s.ClockWait = ClockWait

	// Capture villain probabilities
	s.VillainProbs = make([]savedVillain, len(Villains))
	for i, v := range Villains {
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
	Dead = s.Dead
	Deaths = s.Deaths
	Score = s.Score
	BaseScore = s.BaseScore
	Moves = s.Moves
	Lit = s.Lit
	SuperBrief = s.SuperBrief
	Verbose = s.Verbose
	WonGame = s.WonGame
	HelloSailor = s.HelloSailor
	IsSprayed = s.IsSprayed
	Lucky = s.Lucky
	FumbleNumber = s.FumbleNumber
	FumbleProb = s.FumbleProb

	// Restore dungeon flags
	TrollFlag = s.TrollFlag
	CyclopsFlag = s.CyclopsFlag
	MagicFlag = s.MagicFlag
	LowTide = s.LowTide
	DomeFlag = s.DomeFlag
	EmptyHanded = s.EmptyHanded
	LLDFlag = s.LLDFlag
	RainbowFlag = s.RainbowFlag
	DeflateFlag = s.DeflateFlag
	CoffinCure = s.CoffinCure
	GrateRevealed = s.GrateRevealed
	KitchenWindowFlag = s.KitchenWindowFlag
	CageTop = s.CageTop
	RugMoved = s.RugMoved
	GrUnlock = s.GrUnlock
	CycloWrath = s.CycloWrath
	MirrorMung = s.MirrorMung
	GateFlag = s.GateFlag
	GatesOpen = s.GatesOpen
	WaterLevel = s.WaterLevel
	MatchCount = s.MatchCount
	EggSolve = s.EggSolve
	ThiefHere = s.ThiefHere
	ThiefEngrossed = s.ThiefEngrossed
	LoudFlag = s.LoudFlag
	SingSong = s.SingSong
	BuoyFlag = s.BuoyFlag
	BeachDig = s.BeachDig
	LightShaft = s.LightShaft
	XB = s.XB
	XC = s.XC
	Deflate = s.Deflate
	LampTableIdx = s.LampTableIdx
	CandleTableIdx = s.CandleTableIdx
	LoadAllowed = s.LoadAllowed
	LoadMax = s.LoadMax

	// Restore key references
	Here = idxToObj(s.HereIdx)
	Winner = idxToObj(s.WinnerIdx)
	Player = idxToObj(s.PlayerIdx)
	Params.ItObj = idxToObj(s.ItObjIdx)
	DescObj = idxToObj(s.DescObjIdx)

	// Restore clock state
	for i := range QueueItms {
		QueueItms[i].Run = s.QueueItems[i].Run
		QueueItms[i].Tick = s.QueueItems[i].Tick
		if s.QueueItems[i].RtnName != "" {
			QueueItms[i].Rtn = nameToRtn[s.QueueItems[i].RtnName]
		} else {
			QueueItms[i].Rtn = nil
		}
	}
	QueueInts = s.QueueInts
	QueueDmns = s.QueueDmns
	ClockWait = s.ClockWait

	// Restore villain probabilities
	for i := range Villains {
		if i < len(s.VillainProbs) {
			Villains[i].Prob = s.VillainProbs[i].Prob
		}
	}
}

// promptFilename prompts the user for a save filename and returns it.
func promptFilename(action string) string {
	fmt.Fprintf(GameOutput, "Enter a file name (default is \"%s\"): ", SaveFile)
	if Reader == nil {
		InitReader()
	}
	txt, err := Reader.ReadString('\n')
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
	Save = doSave
	Restore = doRestore
	Restart = doRestart
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
	Lit = IsLit(Here, true)
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
	Here = &WestOfHouse
	ThisIsIt(&Mailbox)
	Lit = true
	Winner = &Adventurer
	Player = Winner
	Winner.MoveTo(Here)

	// Re-register save system for the fresh game
	buildObjIndex()
	Save = doSave
	Restore = doRestore
	Restart = doRestart
	return true
}
