package zork

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

type Action func(ActArg) bool

// PseudoObject are special game objects which only have a single synonym and an action.
type PseudoObj struct {
	Synonym string
	Action  Action
}

type FDir func() *Object
type CDir func() bool

type DirProps struct {
	// NExit represents a non-exit direction where
	// the game outputs NExit string and doesn't
	// move in direction
	NExit string
	// UExit represents a flag which means the user moves to
	// the room in RExit unconditionally
	UExit bool
	// RExit contains a room object where the user should move
	RExit *Object
	// FExit represents a function which executes and if it
	// evaluates to a room object the player moves to that room
	FExit FDir
	// CExit represents a conditional direction where the player moves
	// if it evaluates to true the player is moved to RExit
	CExit CDir
	// CExitStr represents a message printed if the CExit evaluates to false
	CExitStr string
	// DExit represents a game object which if it's open moves the player to RExit
	DExit *Object
	// DExitStr represents a message printed if the DExit is closed
	DExitStr string
}

func (dp DirProps) IsSet() bool {
	return len(dp.NExit) > 0 ||
		(dp.UExit && dp.RExit != nil) ||
		dp.FExit != nil ||
		(dp.CExit != nil && dp.RExit != nil) ||
		(dp.DExit != nil && dp.RExit != nil)
}

// Object represents a game object which can be a character, room, vehicle etc.
type Object struct {
	Flags      Flags
	In         *Object
	Children   []*Object
	Synonyms   []string
	Adjectives []string
	Action     Action
	Global     []*Object
	Pseudo     []PseudoObj
	ContFcn    Action
	DescFcn    Action
	VehType    Flags
	Capacity   int
	Size       int
	Value      int
	TValue     int
	Strength   int
	Text       string
	Desc       string
	LongDesc   string
	FirstDesc  string
	North      DirProps
	South      DirProps
	West       DirProps
	East       DirProps
	NorthWest  DirProps
	NorthEast  DirProps
	SouthWest  DirProps
	SouthEast  DirProps
	Up         DirProps
	Down       DirProps
	Into       DirProps
	Out        DirProps
	Land       DirProps
}

// HasChildren checks if the game object has any children
func (o *Object) HasChildren() bool {
	return len(o.Children) > 0
}

// GetDir returns game object's direction properties if they are set
func (o *Object) GetDir(dir string) *DirProps {
	switch {
	case dir == "north" && o.North.IsSet():
		return &o.North
	case dir == "east" && o.East.IsSet():
		return &o.East
	case dir == "west" && o.West.IsSet():
		return &o.West
	case dir == "south" && o.South.IsSet():
		return &o.South
	case dir == "northeast" && o.NorthEast.IsSet():
		return &o.NorthEast
	case dir == "northwest" && o.NorthWest.IsSet():
		return &o.NorthWest
	case dir == "southeast" && o.SouthEast.IsSet():
		return &o.SouthEast
	case dir == "southwest" && o.SouthWest.IsSet():
		return &o.SouthWest
	case dir == "up" && o.Up.IsSet():
		return &o.Up
	case dir == "down" && o.Down.IsSet():
		return &o.Down
	case dir == "in" && o.Into.IsSet():
		return &o.Into
	case dir == "out" && o.Out.IsSet():
		return &o.Out
	case dir == "land" && o.Land.IsSet():
		return &o.Land
	}
	return nil
}

func (o *Object) Remove() {
	if o.In != nil {
		o.In.RemoveChild(o)
	}
	o.In = nil
}

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

// AddChild adds the game object as a child of the current
// game object. If the child is already present the function
// doesn't add it again.
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

// MoveTo moves the game object to the destination
// supplied. This means the current game object will become
// its child. It also removes the object from its previous
// parent's children list.
func (o *Object) MoveTo(dest *Object) {
	if o.In != nil {
		o.In.RemoveChild(o)
	}
	o.In = dest
	dest.AddChild(o)
}

// Has checks if the current game object has a certain flag set.
func (o *Object) Has(f Flags) bool {
	return o.Flags&f != 0
}

// Give sets a flag on the game object.
func (o *Object) Give(f Flags) {
	o.Flags |= f
}

// Take clears a flag from the game object.
func (o *Object) Take(f Flags) {
	o.Flags &^= f
}

// Is checks if the word is present amongst object's
// synonyms or adjectives. If present it returns true.
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

// Location returns the objects parent.
func (o *Object) Location() *Object {
	return o.In
}

// In returns true if the object's parent
// is the same as the object supplied.
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

// originalState stores the initial state for each object, set once during
// the first BuildObjectTree call. Used by ResetObjectTree to restore the
// object tree to its initial state for tests.
var originalState map[*Object]objectSnapshot

// BuildObjectTree populates each object's
// children present in the Objects global variable.
func BuildObjectTree() {
	saveOriginals := originalState == nil
	if saveOriginals {
		originalState = make(map[*Object]objectSnapshot, len(Objects))
	}
	for _, obj := range Objects {
		if saveOriginals {
			originalState[obj] = objectSnapshot{
				In:       obj.In,
				Flags:    obj.Flags,
				Strength: obj.Strength,
				Value:    obj.Value,
				TValue:   obj.TValue,
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
	// Clear all children first
	for _, obj := range Objects {
		obj.Children = nil
	}
	// Restore original state
	for obj, snap := range originalState {
		obj.In = snap.In
		obj.Flags = snap.Flags
		obj.Strength = snap.Strength
		obj.Value = snap.Value
		obj.TValue = snap.TValue
		obj.Text = snap.Text
	}
	// Rebuild the tree from restored In pointers
	for _, obj := range Objects {
		if obj.Location() != nil {
			obj.Location().AddChild(obj)
		}
	}
}
