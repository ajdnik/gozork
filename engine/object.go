package engine

// Flags represents a bitfield of properties that game objects can have.
type Flags uint64

const (
	// FlgUnk means the value is undefined (zero value, no bits set).
	FlgUnk Flags = 0

	// FlgTake means the object can be picked up by the player.
	FlgTake Flags = 1 << iota
	// FlgTryTake means the object shouldn't be implicitly taken.
	FlgTryTake
	// FlgCont means the object can contain other objects.
	FlgCont
	// FlgDoor means the object is a door.
	FlgDoor
	// FlgOpen means the object can be opened.
	FlgOpen
	// FlgSurf refers to objects such as table, desk, countertop.
	FlgSurf
	// FlgLock means the object is locked and can't be implicitly opened.
	FlgLock
	// FlgWear means the object can be worn by the player.
	FlgWear
	// FlgWorn means the object is currently being worn by the player.
	FlgWorn
	// FlgRead means the object can be read.
	FlgRead
	// FlgLight means the object can be turned on/off.
	FlgLight
	// FlgOn means the object's properties such as Light or Flame are turned on.
	FlgOn
	// FlgFlame means the object can be a source of fire.
	FlgFlame
	// FlgBurn means the object can be burnt.
	FlgBurn
	// FlgTrans means the objects inside can be seen even when it's closed.
	FlgTrans
	// FlgNoDesc means the object shouldn't be described.
	FlgNoDesc
	// FlgInvis means the object shouldn't be found.
	FlgInvis
	// FlgTouch means the object has interacted with the player.
	FlgTouch
	// FlgSearch means to find anything in the object the parser should look as deep as possible.
	FlgSearch
	// FlgVeh means the object can transport the player.
	FlgVeh
	// FlgPerson means the object is a character.
	FlgPerson
	// FlgFemale means the object is a female character.
	FlgFemale
	// FlgVowel means the object's description starts with a vowel.
	FlgVowel
	// FlgNoArt means the object's description doesn't work with articles.
	FlgNoArt
	// FlgPlural means the object's description is a plural noun.
	FlgPlural
	// FlgLand means the object is a dry land room.
	FlgLand
	// FlgWater means the object is a water room.
	FlgWater
	// FlgAir means the object is a room mid-air.
	FlgAir
	// FlgOut means the object is an outdoors room.
	FlgOut
	// FlgIntegral means the object is an integral part of another object.
	FlgIntegral
	// FlgBodyPart means the object is a body part.
	FlgBodyPart
	// FlgNotAll means the object shouldn't be taken when taking all.
	FlgNotAll
	// FlgDrop means if dropping objects in the vehicle the object should stay in the vehicle.
	FlgDrop
	// FlgIn means the player should stay in the vehicle rather than on.
	FlgIn
	// FlgKludge is a syntax flag which can be used to support a syntax VERB PREPOSITION without any object.
	FlgKludge
	// FlgFight means the character is actively engaged in combat.
	FlgFight
	// FlgStagg means the character is recovering from a blow and will skip their next attack.
	FlgStagg
	// FlgSacred means the object is protected and cannot be stolen by the thief.
	FlgSacred
	// FlgTool means the object can be used as an implement (e.g. for digging, inflating, locking).
	FlgTool
	// FlgNonLand means the room is a non-land area.
	FlgNonLand
	// FlgMaze means the room is a maze.
	FlgMaze
	// FlgClimb means the object can be climbed.
	FlgClimb
	// FlgWeapon means the object is a weapon.
	FlgWeapon
	// FlgDrink means the object can be drunk.
	FlgDrink
	// FlgFood means the object can be eaten.
	FlgFood
	// FlgTurn means the object can be turned.
	FlgTurn
	// FlgRMung means the room is munged/destroyed.
	FlgRMung
	// FlgRLand means the room becomes dry land conditionally.
	FlgRLand
	// FlgActor means the object is an actor.
	FlgActor
)

// AnyFlagIn returns true if required is zero (no requirement) or if
// any bit in required is also set in actual.
func AnyFlagIn(required, actual Flags) bool {
	if required == 0 {
		return true
	}
	return required&actual != 0
}

// ActArg represents an argument enum that is passed to Action functions.
type ActArg int

const (
	ActUnk ActArg = iota
	ActBegin
	ActEnter
	ActLook
	ActFlash
	ActObjDesc
	ActEnd
)

// Action is a handler function attached to a game object, invoked with a context argument.
type Action func(ActArg) bool

// PseudoObj are special game objects which only have a single synonym and an action.
type PseudoObj struct {
	Synonym string
	Action  Action
}

// FDir is a function-based room exit that computes the destination dynamically.
type FDir func() *Object

// CDir is a conditional exit check that returns true if passage is allowed.
type CDir func() bool

// Direction represents a compass direction or vertical movement.
type Direction int

const (
	North Direction = iota
	South
	East
	West
	NorthEast
	NorthWest
	SouthEast
	SouthWest
	Up
	Down
	In
	Out
	Land
	NumDirections // sentinel — must be last
)

// dirNames maps Direction values to their canonical string representation.
var dirNames = [NumDirections]string{
	North:     "north",
	South:     "south",
	East:      "east",
	West:      "west",
	NorthEast: "northeast",
	NorthWest: "northwest",
	SouthEast: "southeast",
	SouthWest: "southwest",
	Up:        "up",
	Down:      "down",
	In:        "in",
	Out:       "out",
	Land:      "land",
}

// String returns the canonical name of the direction.
func (d Direction) String() string {
	if d >= 0 && d < NumDirections {
		return dirNames[d]
	}
	return "unknown"
}

// StringToDir maps a direction name to a Direction value.
// Returns the direction and true if found, or -1 and false if not.
func StringToDir(s string) (Direction, bool) {
	d, ok := strToDir[s]
	return d, ok
}

var strToDir = map[string]Direction{
	"north":     North,
	"south":     South,
	"east":      East,
	"west":      West,
	"northeast": NorthEast,
	"northwest": NorthWest,
	"southeast": SouthEast,
	"southwest": SouthWest,
	"up":        Up,
	"down":      Down,
	"in":        In,
	"out":       Out,
	"land":      Land,
}

// AllDirections is the list of all valid directions.
var AllDirections = []Direction{
	North, South, East, West,
	NorthEast, NorthWest, SouthEast, SouthWest,
	Up, Down, In, Out, Land,
}

// DirProps describes how a room exit works (unconditional, conditional, etc.).
type DirProps struct {
	NExit    string
	UExit    bool
	RExit    *Object
	FExit    FDir
	CExit    CDir
	CExitStr string
	DExit    *Object
	DExitStr string
}

// IsSet returns true if any exit data has been configured.
func (dp DirProps) IsSet() bool {
	return len(dp.NExit) > 0 ||
		(dp.UExit && dp.RExit != nil) ||
		dp.FExit != nil ||
		(dp.CExit != nil && dp.RExit != nil) ||
		(dp.DExit != nil && dp.RExit != nil)
}

// ItemData holds properties specific to takeable items and containers.
// Only objects that can be picked up, scored, or hold other objects need this.
type ItemData struct {
	Size     int // weight / bulk of the object
	Value    int // base point value (awarded on first take)
	TValue   int // remaining point value (zeroed after scoring)
	Capacity int // how many child objects this container can hold
}

// CombatData holds properties specific to NPCs that participate in combat.
type CombatData struct {
	Strength int // combat strength (higher = tougher)
}

// VehicleData holds properties specific to vehicles the player can ride.
type VehicleData struct {
	Type Flags // e.g. FlgNonLand for boats
}

// Object represents a game object which can be a character, room, vehicle etc.
// Role-specific data is stored in optional facet pointers that are nil when
// not applicable, keeping the core struct lean.
type Object struct {
	// ---- Core identity (all objects) ----
	Flags      Flags
	In         *Object
	Children   []*Object
	Synonyms   []string
	Adjectives []string
	Desc       string
	LongDesc   string
	FirstDesc  string
	Text       string

	// ---- Behavior ----
	Action  Action
	ContFcn Action
	DescFcn Action

	// ---- Scope / parser hints ----
	Global []*Object
	Pseudo []PseudoObj

	// ---- Room data ----
	Exits map[Direction]DirProps

	// ---- Optional facets (nil when not applicable) ----
	Item    *ItemData    // non-nil for takeable items / containers
	Combat  *CombatData  // non-nil for NPCs that fight
	Vehicle *VehicleData // non-nil for vehicles
}

// HasChildren checks if the game object has any children
func (o *Object) HasChildren() bool {
	return len(o.Children) > 0
}

// GetExit returns the direction properties for the given direction, or nil.
func (o *Object) GetExit(d Direction) *DirProps {
	if o.Exits == nil {
		return nil
	}
	dp, ok := o.Exits[d]
	if !ok || !dp.IsSet() {
		return nil
	}
	return &dp
}

// SetExit sets exit properties for a direction, initializing the map if needed.
func (o *Object) SetExit(d Direction, dp DirProps) {
	if o.Exits == nil {
		o.Exits = make(map[Direction]DirProps)
	}
	o.Exits[d] = dp
}

// Remove detaches the object from its parent.
func (o *Object) Remove() {
	if o.In != nil {
		o.In.RemoveChild(o)
	}
	o.In = nil
}

// RemoveChild removes a direct child from this object's children list.
func (o *Object) RemoveChild(obj *Object) {
	if o.Children == nil {
		return
	}
	if obj == nil {
		return
	}
	found := -1
	for idx, chld := range o.Children {
		if obj == chld {
			found = idx
			break
		}
	}
	if found == -1 {
		return
	}
	o.Children[found] = o.Children[len(o.Children)-1]
	o.Children[len(o.Children)-1] = nil
	o.Children = o.Children[:len(o.Children)-1]
}

// AddChild adds a child to this object's children list (no-op if already present).
func (o *Object) AddChild(child *Object) {
	if o.Children == nil {
		o.Children = []*Object{}
	}
	for _, ch := range o.Children {
		if ch == child {
			return
		}
	}
	o.Children = append(o.Children, child)
}

// MoveTo removes the object from its current parent and places it inside dest.
func (o *Object) MoveTo(dest *Object) {
	if o.In != nil {
		o.In.RemoveChild(o)
	}
	o.In = dest
	dest.AddChild(o)
}

// Has returns true if any of the given flag bits are set.
func (o *Object) Has(f Flags) bool {
	return o.Flags&f != 0
}

// Give sets the given flag bits on the object.
func (o *Object) Give(f Flags) {
	o.Flags |= f
}

// Take clears the given flag bits on the object.
func (o *Object) Take(f Flags) {
	o.Flags &^= f
}

// Is returns true if wrd appears in the object's synonyms or adjectives.
func (o *Object) Is(wrd string) bool {
	for _, syn := range o.Synonyms {
		if syn == wrd {
			return true
		}
	}
	for _, adj := range o.Adjectives {
		if adj == wrd {
			return true
		}
	}
	return false
}

// Location returns the object's parent (container or room).
func (o *Object) Location() *Object {
	return o.In
}

// IsIn returns true if the object's direct parent is loc.
func (o *Object) IsIn(loc *Object) bool {
	return o.In == loc
}

// objectSnapshot holds the original state of an object for test resets.
type objectSnapshot struct {
	In       *Object
	Flags    Flags
	Strength int
	Value    int
	TValue   int
	Text     string
}

// ---- Nil-safe accessors for facet fields ----

// GetSize returns the item size, or 0 if the object has no ItemData.
func (o *Object) GetSize() int {
	if o.Item == nil {
		return 0
	}
	return o.Item.Size
}

// SetSize sets the item size, allocating ItemData if needed.
func (o *Object) SetSize(v int) {
	if o.Item == nil {
		o.Item = &ItemData{}
	}
	o.Item.Size = v
}

// GetValue returns the scoring value, or 0 if the object has no ItemData.
func (o *Object) GetValue() int {
	if o.Item == nil {
		return 0
	}
	return o.Item.Value
}

// SetValue sets the scoring value, allocating ItemData if needed.
func (o *Object) SetValue(v int) {
	if o.Item == nil {
		o.Item = &ItemData{}
	}
	o.Item.Value = v
}

// GetTValue returns the treasure value, or 0 if the object has no ItemData.
func (o *Object) GetTValue() int {
	if o.Item == nil {
		return 0
	}
	return o.Item.TValue
}

// SetTValue sets the treasure value, allocating ItemData if needed.
func (o *Object) SetTValue(v int) {
	if o.Item == nil {
		o.Item = &ItemData{}
	}
	o.Item.TValue = v
}

// GetCapacity returns the container capacity, or 0 if the object has no ItemData.
func (o *Object) GetCapacity() int {
	if o.Item == nil {
		return 0
	}
	return o.Item.Capacity
}

// SetCapacity sets the container capacity, allocating ItemData if needed.
func (o *Object) SetCapacity(v int) {
	if o.Item == nil {
		o.Item = &ItemData{}
	}
	o.Item.Capacity = v
}

// GetStrength returns the combat strength, or 0 if the object has no CombatData.
func (o *Object) GetStrength() int {
	if o.Combat == nil {
		return 0
	}
	return o.Combat.Strength
}

// SetStrength sets the combat strength, allocating CombatData if needed.
func (o *Object) SetStrength(v int) {
	if o.Combat == nil {
		o.Combat = &CombatData{}
	}
	o.Combat.Strength = v
}

// GetVehType returns the vehicle type flags, or 0 if the object has no VehicleData.
func (o *Object) GetVehType() Flags {
	if o.Vehicle == nil {
		return 0
	}
	return o.Vehicle.Type
}

// SetVehType sets the vehicle type flags, allocating VehicleData if needed.
func (o *Object) SetVehType(v Flags) {
	if o.Vehicle == nil {
		o.Vehicle = &VehicleData{}
	}
	o.Vehicle.Type = v
}

var originalState map[*Object]objectSnapshot

// BuildObjectTree populates each object's children from G.AllObjects.
func BuildObjectTree() {
	saveOriginals := originalState == nil
	if saveOriginals {
		originalState = make(map[*Object]objectSnapshot, len(G.AllObjects))
	}
	for _, obj := range G.AllObjects {
		if saveOriginals {
			originalState[obj] = objectSnapshot{
				In:       obj.In,
				Flags:    obj.Flags,
				Strength: obj.GetStrength(),
				Value:    obj.GetValue(),
				TValue:   obj.GetTValue(),
				Text:     obj.Text,
			}
		}
		if obj.Location() != nil {
			obj.Location().AddChild(obj)
		}
	}
}

// ResetObjectTree restores every object to its original state and rebuilds
// the children tree. Used by tests to get a fresh game state.
func ResetObjectTree() {
	if originalState == nil {
		return
	}
	for _, obj := range G.AllObjects {
		obj.Children = nil
	}
	for obj, snap := range originalState {
		obj.In = snap.In
		obj.Flags = snap.Flags
		obj.Text = snap.Text
		obj.SetStrength(snap.Strength)
		obj.SetValue(snap.Value)
		obj.SetTValue(snap.TValue)
	}
	for _, obj := range G.AllObjects {
		if obj.Location() != nil {
			obj.Location().AddChild(obj)
		}
	}
}
