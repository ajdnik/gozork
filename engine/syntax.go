package engine

// LocFlag represents a single location constraint for parser object resolution.
type LocFlag uint8

const (
	// LocHeld means the object must be directly held by the actor.
	LocHeld LocFlag = 1 << iota
	// LocCarried means the object must be carried (held or in a carried container).
	LocCarried
	// LocInRoom means the object must be in the current room.
	LocInRoom
	// LocOnGrnd means the object must be on the ground.
	LocOnGrnd
	// LocTake means the parser may attempt an implicit take.
	LocTake
	// LocMany means multiple direct objects are allowed.
	LocMany
	// LocHave means the object must already be possessed.
	LocHave
)

// In returns true if this flag appears in the given flag set.
func (lf LocFlag) In(flgs LocFlags) bool {
	return flgs&LocFlags(lf) != 0
}

// LocFlags is a set of LocFlag constraints.
type LocFlags uint8

const LocAll LocFlags = LocFlags(LocHeld | LocCarried | LocInRoom | LocOnGrnd | LocTake | LocMany | LocHave)

// LocSet builds a LocFlags bitset from a list of flags.
func LocSet(flags ...LocFlag) LocFlags {
	var res LocFlags
	for _, flg := range flags {
		res |= LocFlags(flg)
	}
	return res
}

// All replaces the set with every possible LocFlag value.
func (lfs *LocFlags) All() LocFlags {
	*lfs = LocAll
	return *lfs
}

// HasAll returns true if all LocFlag values are present.
func (lfs LocFlags) HasAll() bool {
	return lfs == LocAll
}

// VerbAction is a handler function invoked when a verb is performed.
type VerbAction func(ActionArg) bool

// ObjProp describes the expected properties of an object slot in a syntax definition.
type ObjProp struct {
	ObjFlags Flags
	LocFlags LocFlags
	HasObj   bool
}

// Syntax defines a single command syntax pattern (verb + prepositions + object slots).
type Syntax struct {
	NormVerb  string
	Verb      string
	VrbPrep   string
	Obj1      ObjProp
	ObjPrep   string
	Obj2      ObjProp
	Action    VerbAction
	PreAction VerbAction
}

// NumObjects returns how many object slots this syntax expects (0, 1, or 2).
func (s *Syntax) NumObjects() int {
	if !s.Obj1.HasObj {
		return 0
	}
	if !s.Obj2.HasObj {
		return 1
	}
	return 2
}

// IsVrbPrep returns true if the verb preposition matches.
func (s *Syntax) IsVrbPrep(prep string) bool { return s.VrbPrep == prep }

// IsObjPrep returns true if the object preposition matches.
func (s *Syntax) IsObjPrep(prep string) bool { return s.ObjPrep == prep }

// GetActionVerb returns the full verb string including its preposition.
func (s *Syntax) GetActionVerb() string {
	av := s.Verb
	if len(s.VrbPrep) > 0 {
		av += " " + s.VrbPrep
	}
	return av
}

// GetNormVerb returns the normalized verb key for action dispatch.
func (s *Syntax) GetNormVerb() string {
	if len(s.NormVerb) > 0 {
		return s.NormVerb
	}
	if vrb, ok := G.NormVerbs[s.GetActionVerb()]; ok {
		return vrb
	}
	return s.GetActionVerb()
}

// RndSelect supports non-repeating random selection from a pool of strings.
type RndSelect struct {
	Unselected []string
	Selected   []string
}

// ActionVerb stores both the normalized and original forms of the current verb.
type ActionVerb struct {
	Norm string
	Orig string
}

// Commands holds the game's syntax definitions. Set by the game package
// during initialization and used by the parser's SyntaxCheck.
var Commands []Syntax

// AddToVocab adds a word with its type to the global Vocabulary map.
func AddToVocab(wrd string, typ WordTyp) {
	v, ok := Vocabulary[wrd]
	if !ok {
		Vocabulary[wrd] = WordItem{
			Norm:  wrd,
			Types: WordTypes{typ},
		}
	} else {
		Vocabulary[wrd] = WordItem{
			Norm:  wrd,
			Types: append(v.Types, typ),
		}
	}
}

// BuildVocabulary builds the vocabulary and action maps from the provided
// game-specific data. The game package calls this during initialization.
func BuildVocabulary(commands []Syntax, buzzWords []string, synonyms map[string]string) {
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
			Vocabulary[key] = WordItem{
				Norm:  val,
				Types: nil,
			}
			if el, ok := Vocabulary[val]; ok {
				Vocabulary[key] = WordItem{
					Norm:  val,
					Types: append(WordTypes{}, el.Types...),
				}
			}
		}
	}
}
