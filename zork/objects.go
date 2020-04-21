package zork

type ActionArg int

const (
	MEnd ActionArg = iota
	MHandled
	MEnter
	MLook
	MWake
	MFight
	MBefore
)

type Action func(ActionArg) bool

type Object struct {
	Properties  []Property
	Name        string
	Description string
	Action      Action
	Parent      *Object
	Contains    []*Object
	Initial     string
	Initial2    Action
	Adjectives  []string
	Nouns       []string
	Capacity    int
	VType       int
	Strength    int
}

const (
	NoDescription = ""
	NoInitial     = ""
)

func (o *Object) Has(prop Property) bool {
	if o.Properties == nil {
		return false
	}
	for _, p := range o.Properties {
		if p == prop {
			return true
		}
	}
	return false
}

func (o *Object) Give(prop Property) {
	if o.Properties == nil {
		o.Properties = []Property{}
	}
	if o.Has(prop) {
		return
	}
	o.Properties = append(o.Properties, prop)
}

func (o *Object) Take(prop Property) {
	if o.Properties == nil {
		return
	}
	i := -1
	for idx, p := range o.Properties {
		if p == prop {
			i = idx
			break
		}
	}
	if i == -1 {
		return
	}
	o.Properties[i] = o.Properties[len(o.Properties)-1]
	o.Properties[len(o.Properties)-1] = Property("")
	o.Properties = o.Properties[:len(o.Properties)-1]
}

func (o *Object) MoveTo(obj *Object) {
	if obj.Contains == nil {
		obj.Contains = []*Object{}
	}
	o.Parent = obj
	for _, itm := range obj.Contains {
		if itm.Name == o.Name {
			return
		}
	}
	obj.Contains = append(obj.Contains, o)
}

var (
	WestOfHouse = &Object{
		Name:       "West of House",
		Properties: []Property{DryLand, Light, Sacred},
		Action: func(rarg ActionArg) bool {
			if rarg != MLook {
				return false
			}
			Print("You are standing in an open field west of a white house, with a boarded front door.")
			if WonFlag {
				Print(" A secret path leads southwest into the forest.")
			}
			NewLine()
			return true
		},
	}
	SmallMailbox = &Object{
		Name: "small mailbox",
		Action: func(_ ActionArg) bool {
			if Verb != "take" {
				return false
			}
			if Noun == nil {
				return false
			}
			if Noun.Name != "small mailbox" {
				return false
			}
			Print("It is securely anchored.")
			return true
		},
		Adjectives: []string{"small"},
		Capacity:   10,
		Nouns:      []string{"mailbox", "box"},
		Properties: []Property{Container, TryTakeBit},
	}
	Cretin = &Object{
		Name:     "cretin",
		Nouns:    []string{"adventurer"},
		Strength: 0,
		Action: func(_ ActionArg) bool {
			return false
		},
		Properties: []Property{Animate, Concealed, Sacred, Scenery},
	}
	TrophyCase = &Object{}
	MagicBoat  = &Object{}
)
