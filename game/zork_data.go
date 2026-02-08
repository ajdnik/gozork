package game

import . "github.com/ajdnik/gozork/engine"

// ZorkData holds all Zork I specific mutable state that is not part
// of the generic engine. Stored in G.GameData during play.
type ZorkData struct {
	// ---- Game progress ----
	Dead    bool
	Deaths  int
	WonGame bool

	// ---- Misc counters ----
	HelloSailor  int
	IsSprayed    bool
	FumbleNumber int
	FumbleProb   int
	DescObj      *Object

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
}

// newZorkData returns a fresh ZorkData with non-zero defaults.
func newZorkData() *ZorkData {
	return &ZorkData{
		FumbleNumber: 7,
		FumbleProb:   8,
		CageTop:      true,
		MatchCount:   6,
		BuoyFlag:     true,
		BeachDig:     -1,
		LightShaft:   13,
		LoadAllowed:  100,
		LoadMax:      100,
	}
}

// gD returns the game-specific ZorkData from the engine's GameState.
func gD() *ZorkData {
	return G.GameData.(*ZorkData)
}
