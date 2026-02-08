package zork

import (
	"bufio"
	"io"
	"math/rand"
	"os"
	"time"
)

// RNG abstracts random number generation so tests can inject
// a deterministic source. The single method matches math/rand.Rand.Intn.
type RNG interface {
	Intn(n int) int
}

// G is the current game state. All mutable game state is accessed through
// this pointer, making it easy to create fresh instances for tests or
// save/restore without relying on scattered global variables.
var G *GameState

// GameState holds all mutable state for a single game session.
type GameState struct {
	// ---- Core references ----
	Here   *Object
	Winner *Object
	Player *Object

	// ---- Current command (set by parser / Perform) ----
	DirObj            *Object
	IndirObj          *Object
	ActVerb           ActionVerb
	DirObjPossibles   []*Object
	IndirObjPossibles []*Object
	DetectedSyntx     *Syntx

	// ---- Game progress ----
	Moves     int
	Score     int
	BaseScore int
	Dead      bool
	Deaths    int
	Lit       bool
	WonGame   bool

	// ---- Settings ----
	SuperBrief bool
	Verbose    bool
	Lucky      bool

	// ---- Misc counters ----
	HelloSailor  int
	IsSprayed    bool
	FumbleNumber int
	FumbleProb   int
	DescObj      *Object

	// ---- Parser internals ----
	ParserOk      bool
	Script        bool
	PerformFatal  bool
	QuitRequested bool
	AlwaysLit     bool
	Search        FindProps
	ParsedSyntx   ParseTbl
	OrphanedSyntx ParseTbl
	Params        ParseProps
	NotHere       NotHereProps
	LexRes        []LexItm
	Reserv        ReserveProps
	Again         AgainProps
	Oops          OopsProps

	// ---- Clock / interrupt queue ----
	QueueItms [30]ClockEvent
	QueueInts int
	QueueDmns int
	ClockWait bool

	// ---- Dungeon flags ----
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
	LampTableIdx      int
	CandleTableIdx    int
	XB                bool
	XC                bool
	Deflate           bool
	LoadAllowed       int
	LoadMax           int

	// ---- Combat ----
	Villains []*VillainEntry

	// ---- I/O ----
	GameOutput     io.Writer
	GameInput      io.Reader
	Reader         *bufio.Reader
	InputExhausted bool

	// ---- RNG ----
	Rand RNG

	// ---- Save/Restore/Restart function hooks ----
	Save    func() bool
	Restore func() bool
	Restart func() bool
}

// NewGameState creates a fresh GameState with all default values set.
func NewGameState() *GameState {
	return &GameState{
		// Non-zero defaults
		Lucky:        true,
		FumbleNumber: 7,
		FumbleProb:   8,
		CageTop:      true,
		MatchCount:   6,
		BuoyFlag:     true,
		BeachDig:     -1,
		LightShaft:   13,
		LoadAllowed:  100,
		LoadMax:      100,
		QueueInts:    30,
		QueueDmns:    30,

		// I/O defaults
		GameOutput: os.Stdout,
		GameInput:  os.Stdin,

		// RNG default: time-seeded source
		Rand: rand.New(rand.NewSource(time.Now().UnixNano())),

		// Stub function hooks (replaced by initSaveSystem)
		Save:    func() bool { return false },
		Restore: func() bool { return false },
		Restart: func() bool { return false },
	}
}
