package zork

// Flag represents various properties game objects can have
type Flag int

const (
	// Take means the object can be picked up by the player
	FlgTake Flag = iota
	// TryTake means the object shouldn't be implicitly taken
	FlgTryTake
	// Container means the object can contain other objects
	FlgCont
	// Door means the object is a door
	FlgDoor
	// Open means the object can be opened
	FlgOpen
	// Surface refres to objects such as table, desk, countertop...
	FlgSurf
	// Locked means the object is locked and can't be implicitly opened
	FlgLock
	// Wear means the object can be worn by the player
	FlgWear
	// Worn means the object is currently being worn by the player
	FlgWorn
	// Read means the object can be read
	FlgRead
	// Light means the object can be turned on/off
	FlgLight
	// On means the object's properties such as Light or Flame are turned on
	FlgOn
	// Flame means the object can be a source of fire
	FlgFlame
	// Burn means the object can be burnt
	FlgBurn
	// Transparent means the objects inside can be seen even when it's closed
	FlgTrans
	// NoDescription means the object shouldn't be described
	FlgNoDesc
	// Invisible means the object shouldn't be found
	FlgInvis
	// Touch means the object has interacted with the player
	FlgTouch
	// Search means to find anything in the object the parser should look as deep as possible
	FlgSearch
	// Vehicle means the object can transport the player
	FlgVeh
	// Person means the object is a character
	FlgPerson
	// Female means the object is a female character
	FlgFemale
	// Vowel means the object's description starts with a vowel
	FlgVowel
	// NoArticle means the object's description doesn't work with articles
	FlgNoArt
	// Plural means the object's description is a plural noun
	FlgPlural
	// RoomLand means the object is a dry land room
	FlgLand
	// RoomWater means the object is a water room
	FlgWater
	// RoomAir means the object is a room mid-air
	FlgAir
	// Outside means the object is an outdoors room
	FlgOut
	// Integral means the object is an integral part of another object and can't be taken or dropped independently
	FlgIntegral
	// BodyPart means the object is a body part
	FlgBodyPart
	// NotAll means the object shouldn't be taken when taking all
	FlgNotAll
	// Drop means if dropping objects in the vehicle the object should stay in the vehicle
	FlgDrop
	// In means the player should stay in the vehicle rather than on
	FlgIn
	// Kludge is a syntax flag which can be used to support a syntax VERB PREPOSITION without any object
	FlgKludge
	// TODO: describe what the flag is used for
	FlgFight
	// TODO: describe what the flag is used for
	FlgStagg
	// TODO: describe
	FlgSacred
	// TODO: describe
	FlgTool
	FlgNonLand
	// FlgMaze means the room is a maze
	FlgMaze
	FlgClimb
)

// In function checks if the current flag is in the flag slice.
// Returns true if it's in otherwise it returns false.
func (f Flag) In(flgs []Flag) bool {
	for _, fl := range flgs {
		if fl == f {
			return true
		}
	}
	return false
}

// AnyFlagIn compares two flag slices and returns true
// if any of the flags in the first slice are present in the
// second slice.
func AnyFlagIn(any, flgs []Flag) bool {
	if len(any) == 0 {
		return true
	}
	for _, af := range any {
		for _, fl := range flgs {
			if af == fl {
				return true
			}
		}
	}
	return false
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

// Object represents a game object which can be a character, room, vehicle etc.
type Object struct {
	Flags      []Flag
	In         *Object
	Children   []*Object
	Synonyms   []string
	Adjectives []string
	Action     Action
	Global     []*Object
	Pseudo     []PseudoObj
	ContFcn    Action
	DescFcn    Action
	VehType    Flag
	Capacity   int
	Size       int
	Value      int
	Strength   int
	Text       string
	Desc       string
	LongDesc   string
	FirstDesc  string
}

func (o *Object) HasChildren() bool {
	return len(o.Children) > 0
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
// its child.
func (o *Object) MoveTo(dest *Object) {
	o.In = dest
	dest.AddChild(o)
}

// Has checks if the current game object has a certain
// flag set.
func (o *Object) Has(f Flag) bool {
	return f.In(o.Flags)
}

// Give gives the game object a flag.
// If the flag is alread set nothing happend.
func (o *Object) Give(f Flag) {
	if o.Flags == nil {
		o.Flags = []Flag{}
	}
	if f.In(o.Flags) {
		return
	}
	o.Flags = append(o.Flags, f)
}

// Take removes a flag from the game object.
// If the flag isn't set nothing happens.
func (o *Object) Take(f Flag) {
	found := -1
	for idx, flg := range o.Flags {
		if flg == f {
			found = idx
			break
		}
	}
	if found == -1 {
		return
	}
	o.Flags[found] = o.Flags[len(o.Flags)-1]
	o.Flags[len(o.Flags)-1] = -1
	o.Flags = o.Flags[:len(o.Flags)-1]
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

// BuildObjectTree populates each objects
// childrens present in the Objects global variable.
func BuildObjectTree() {
	for _, obj := range Objects {
		if obj.Location() != nil {
			obj.Location().AddChild(obj)
		}
	}
}
