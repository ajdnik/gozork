package engine

import (
	"bufio"
	"fmt"
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
// Engine-generic fields live here directly; game-specific extensions
// are stored in the GameData field.
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
	Lit       bool

	// ---- Settings ----
	SuperBrief bool
	Verbose    bool
	Lucky      bool

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

	// ---- I/O ----
	GameOutput     io.Writer
	GameInput      io.Reader
	Reader         *bufio.Reader
	InputExhausted bool

	// ---- RNG ----
	Rand RNG

	// ---- Save/Restore/Restart function hooks ----
	Save    func() error
	Restore func() error
	Restart func() error

	// ---- Well-known objects (set by game during init) ----
	AllObjects     []*Object // complete list of all game objects
	RoomsObj       *Object   // container for all rooms
	GlobalObj      *Object   // global objects container
	LocalGlobalObj *Object   // local globals container
	NotHereObj     *Object   // sentinel for "not here" objects
	PseudoObj      *Object   // pseudo object for pseudo-object actions
	ItPronounObj   *Object   // the "it" pronoun object
	MeObj          *Object   // the "me" player-reference object
	HandsObj       *Object   // the player's hands object

	// ---- Vocabulary registries (populated by BuildVocabulary) ----
	Actions    map[string]VrbAction
	PreActions map[string]VrbAction
	NormVerbs  map[string]string

	// ---- Clock function registry (populated by game) ----
	ClockFuncs map[string]func() bool

	// ---- Game-specific callbacks (set by game during init) ----
	// ITakeFunc is the implicit-take handler called by the parser.
	ITakeFunc func(vb bool) bool

	// ---- Game-specific extension data ----
	// The game package stores its own state struct here and accesses
	// it via a type assertion (e.g. G.GameData.(*ZorkData)).
	GameData interface{}
}

// NewGameState creates a fresh GameState with all default values set.
func NewGameState() *GameState {
	return &GameState{
		// Non-zero defaults
		Lucky:     true,
		QueueInts: 30,
		QueueDmns: 30,

		// I/O defaults
		GameOutput: os.Stdout,
		GameInput:  os.Stdin,

		// RNG default: time-seeded source
		Rand: rand.New(rand.NewSource(time.Now().UnixNano())),

		// Stub function hooks (replaced by game's initSaveSystem)
		Save:    func() error { return fmt.Errorf("save system not initialized") },
		Restore: func() error { return fmt.Errorf("restore system not initialized") },
		Restart: func() error { return fmt.Errorf("restart system not initialized") },

		// Initialize maps
		Actions:    make(map[string]VrbAction),
		PreActions: make(map[string]VrbAction),
		NormVerbs:  make(map[string]string),
		ClockFuncs: make(map[string]func() bool),
	}
}
