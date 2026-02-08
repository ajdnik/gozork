package engine

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
	ObjFlags Flags
	LocFlags LocFlags
	HasObj   bool
}

type Syntx struct {
	NormVerb  string
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

func (s *Syntx) IsVrbPrep(prep string) bool { return s.VrbPrep == prep }
func (s *Syntx) IsObjPrep(prep string) bool { return s.ObjPrep == prep }

func (s *Syntx) GetActionVerb() string {
	av := s.Verb
	if len(s.VrbPrep) > 0 {
		av += " " + s.VrbPrep
	}
	return av
}

func (s *Syntx) GetNormVerb() string {
	if len(s.NormVerb) > 0 {
		return s.NormVerb
	}
	if vrb, ok := G.NormVerbs[s.GetActionVerb()]; ok {
		return vrb
	}
	return s.GetActionVerb()
}

type RndSelect struct {
	Unselected []string
	Selected   []string
}

type ActionVerb struct {
	Norm string
	Orig string
}

// Commands holds the game's syntax definitions. Set by the game package
// during initialization and used by the parser's SyntaxCheck.
var Commands []Syntx

// AddToVocab adds a word with its type to the global Vocabulary map.
func AddToVocab(wrd string, typ WordTyp) {
	v, ok := Vocabulary[wrd]
	if !ok {
		Vocabulary[wrd] = WordItm{
			Norm:  wrd,
			Types: WordTypes{typ},
		}
	} else {
		Vocabulary[wrd] = WordItm{
			Norm:  wrd,
			Types: append(v.Types, typ),
		}
	}
}

// BuildVocabulary builds the vocabulary and action maps from the provided
// game-specific data. The game package calls this during initialization.
func BuildVocabulary(commands []Syntx, buzzWords []string, synonyms map[string]string) {
	Commands = commands
	// Add buzz words
	for _, bw := range buzzWords {
		AddToVocab(bw, WordBuzz)
	}
	// Add verbs
	for _, cmd := range commands {
		AddToVocab(cmd.Verb, WordVerb)
		if len(cmd.VrbPrep) > 0 {
			AddToVocab(cmd.VrbPrep, WordPrep)
		}
		if len(cmd.ObjPrep) > 0 {
			AddToVocab(cmd.ObjPrep, WordPrep)
		}
		actionKey := cmd.GetActionVerb()
		if len(cmd.NormVerb) > 0 && cmd.NormVerb != actionKey {
			G.Actions[cmd.NormVerb] = cmd.Action
			if cmd.PreAction != nil {
				G.PreActions[cmd.NormVerb] = cmd.PreAction
			}
		} else {
			G.Actions[actionKey] = cmd.Action
			if cmd.PreAction != nil {
				G.PreActions[actionKey] = cmd.PreAction
			}
		}
		if len(cmd.NormVerb) == 0 || cmd.NormVerb == cmd.GetActionVerb() {
			G.NormVerbs[cmd.GetActionVerb()] = cmd.GetActionVerb()
		}
	}
	// Add directions
	for _, d := range AllDirections {
		AddToVocab(d.String(), WordDir)
	}
	// Add objects
	for _, obj := range G.AllObjects {
		if obj.Synonyms != nil {
			for _, s := range obj.Synonyms {
				AddToVocab(s, WordObj)
			}
		}
		if obj.Adjectives != nil {
			for _, a := range obj.Adjectives {
				AddToVocab(a, WordAdj)
			}
		}
	}
	// Add synonyms
	for key, val := range synonyms {
		if _, ok := Vocabulary[key]; !ok {
			Vocabulary[key] = WordItm{
				Norm:  val,
				Types: nil,
			}
			if el, ok := Vocabulary[val]; ok {
				Vocabulary[key] = WordItm{
					Norm:  val,
					Types: append(WordTypes{}, el.Types...),
				}
			}
		}
	}
}
