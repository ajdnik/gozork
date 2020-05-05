package zork

type LocFlag int

const (
	LocHeld LocFlag = iota
	LocCarried
	LocInRoom
	LocOnGrnd
	LocTake
	LocMany
	LocHave
)

func (lf LocFlag) In(flgs LocFlags) bool {
	for _, flg := range flgs {
		if flg == lf {
			return true
		}
	}
	return false
}

type LocFlags []LocFlag

func (lfs *LocFlags) All() LocFlags {
	*lfs = LocFlags{LocHeld, LocCarried, LocInRoom, LocOnGrnd, LocTake, LocMany, LocHave}
	return *lfs
}

func (lfs LocFlags) HasAll() bool {
	return len(lfs) == 7
}

type VrbAction func(ActArg) bool

type ObjProp struct {
	ObjFlags []Flag
	LocFlags LocFlags
	HasObj   bool
}

type Syntx struct {
	Verb      string
	VrbPrep   string
	Obj1      ObjProp
	ObjPrep   string
	Obj2      ObjProp
	Action    VrbAction
	PreAction VrbAction
}

func (s *Syntx) NumObjects() int {
	if !s.Obj1.HasObj {
		return 0
	}
	if !s.Obj2.HasObj {
		return 1
	}
	return 2
}

func (s *Syntx) IsVrbPrep(prep string) bool {
	return s.VrbPrep == prep
}

func (s *Syntx) IsObjPrep(prep string) bool {
	return s.ObjPrep == prep
}

func (s *Syntx) GetActionVerb() string {
	av := s.Verb
	if len(s.VrbPrep) > 0 {
		av += " " + s.VrbPrep
	}
	return av
}

var (
	BuzzWords = []string{
		"again", "g", "oops", "a", "an", "the",
		"is", "and", "of", "then", "all", "one",
		"but", "except", ".", ",", "\"", "yes",
		"no", "y", "here",
	}
	Synonyms = map[string]string{
		"using":      "with",
		"through":    "with",
		"thru":       "with",
		"into":       "in",
		"inside":     "in",
		"onto":       "on",
		"underneath": "under",
		"beneath":    "under",
		"below":      "under",
		"n":          "north",
		"s":          "south",
		"e":          "east",
		"w":          "west",
		"d":          "down",
		"u":          "up",
		"nw":         "northwest",
		"ne":         "northeast",
		"sw":         "southwest",
		"se":         "southeast",
		"superbrief": "super",
		"i":          "inventory",
		"q":          "quit",
		"l":          "look",
		"stare":      "look",
		"gaze":       "look",
		"describe":   "examine",
		"what":       "examine",
		"whats":      "examine",
	}
	Directions = []string{
		"north", "east", "west", "south", "northeast",
		"northwest", "southeast", "southwest", "up",
		"down", "in", "out", "land",
	}
	Commands = []Syntx{
		{
			Verb:   "quit",
			Action: VQuit,
		},
		{
			Verb:   "score",
			Action: VScore,
		},
		{
			Verb:   "version",
			Action: VVersion,
		},
		{
			Verb:   "examine",
			Obj1:   ObjProp{HasObj: true, LocFlags: LocFlags{LocMany}},
			Action: VExamine,
		},
		{
			Verb:    "examine",
			VrbPrep: "in",
			Obj1:    ObjProp{HasObj: true, LocFlags: LocFlags{LocHeld, LocCarried, LocInRoom, LocOnGrnd, LocMany}},
			Action:  VLookInside,
		},
		{
			Verb:    "examine",
			VrbPrep: "on",
			Obj1:    ObjProp{HasObj: true, LocFlags: LocFlags{LocHeld, LocCarried, LocInRoom, LocOnGrnd, LocMany}},
			Action:  VLookInside,
		},
		{
			Verb:   "look",
			Action: VLook,
		},
		{
			Verb:    "look",
			VrbPrep: "around",
			Obj1:    ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgKludge}},
			Action:  VLook,
		},
		{
			Verb:    "look",
			VrbPrep: "up",
			Obj1:    ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgKludge}},
			Action:  VLook,
		},
		{
			Verb:    "look",
			VrbPrep: "down",
			Obj1:    ObjProp{HasObj: true, ObjFlags: []Flag{FlgSearch, FlgKludge}},
			Action:  VLook,
		},
		{
			Verb:    "look",
			VrbPrep: "at",
			Obj1:    ObjProp{HasObj: true, LocFlags: LocFlags{LocHeld, LocCarried, LocInRoom, LocOnGrnd, LocMany}},
			Action:  VExamine,
		},
		{
			Verb:    "look",
			VrbPrep: "on",
			Obj1:    ObjProp{HasObj: true},
			Action:  VLookOn,
		},
		{
			Verb:    "look",
			VrbPrep: "with",
			Obj1:    ObjProp{HasObj: true, LocFlags: LocFlags{LocHeld, LocCarried, LocInRoom, LocOnGrnd, LocMany}},
			Action:  VLookInside,
		},
		{
			Verb:    "look",
			VrbPrep: "under",
			Obj1:    ObjProp{HasObj: true},
			Action:  VLookUnder,
		},
		{
			Verb:    "look",
			VrbPrep: "behind",
			Obj1:    ObjProp{HasObj: true},
			Action:  VLookBehind,
		},
		{
			Verb:    "look",
			VrbPrep: "in",
			Obj1:    ObjProp{HasObj: true, LocFlags: LocFlags{LocHeld, LocCarried, LocInRoom, LocOnGrnd, LocMany}},
			Action:  VLookInside,
		},
	}
	Actions    = make(map[string]VrbAction)
	PreActions = make(map[string]VrbAction)
)

func addToVocab(wrd string, typ WordTyp) {
	if _, ok := Vocabulary[wrd]; !ok {
		Vocabulary[wrd] = WordItm{
			Norm: wrd,
			Type: typ,
		}
	}
}

func BuildVocabulary() {
	// Add buzz words
	for _, bw := range BuzzWords {
		addToVocab(bw, WordBuzz)
	}
	// Add verbs
	for _, cmd := range Commands {
		addToVocab(cmd.Verb, WordVerb)
		if len(cmd.VrbPrep) > 0 {
			addToVocab(cmd.VrbPrep, WordPrep)
		}
		if len(cmd.ObjPrep) > 0 {
			addToVocab(cmd.ObjPrep, WordPrep)
		}
		Actions[cmd.GetActionVerb()] = cmd.Action
		if cmd.PreAction != nil {
			PreActions[cmd.GetActionVerb()] = cmd.PreAction
		}
	}
	// Add directions
	for _, dir := range Directions {
		addToVocab(dir, WordDir)
	}
	// Add objects
	for _, obj := range Objects {
		if obj.Synonyms != nil {
			for _, s := range obj.Synonyms {
				addToVocab(s, WordObj)
			}
		}
		if obj.Adjectives != nil {
			for _, a := range obj.Adjectives {
				addToVocab(a, WordAdj)
			}
		}
	}
	// Add synonyms
	for key, val := range Synonyms {
		if _, ok := Vocabulary[key]; !ok {
			Vocabulary[key] = WordItm{
				Norm: val,
				Type: WordUnk,
			}
			if el, ok := Vocabulary[val]; ok {
				Vocabulary[key] = WordItm{
					Norm: val,
					Type: el.Type,
				}
			}
		}
	}
}
