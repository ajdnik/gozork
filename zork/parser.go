package zork

import "math/rand"

type ParseTbl struct {
	Verb         LexItm
	Prep1        LexItm
	ObjOrClause1 []LexItm
	Obj1Start    int
	Obj1End      int
	Prep2        LexItm
	ObjOrClause2 []LexItm
}

func (pt *ParseTbl) Set(tbl ParseTbl) {
	pt.Verb.Set(tbl.Verb)
	pt.Prep1.Set(tbl.Prep1)
	pt.ObjOrClause1 = append([]LexItm{}, tbl.ObjOrClause1...)
	pt.Obj1Start = tbl.Obj1Start
	pt.Obj1End = tbl.Obj1End
	pt.Prep2.Set(tbl.Prep2)
	pt.ObjOrClause2 = append([]LexItm{}, tbl.ObjOrClause2...)
}

type NotHereProps struct {
	Syn LexItm
	Adj LexItm
}

type FindProps struct {
	ObjFlags []Flag
	LocFlags LocFlags
	Syn      LexItm
	Adj      LexItm
}

type ClauseTyp int

const (
	ClauseUnk ClauseTyp = iota
	Clause1
	Clause2
)

type ClauseProps struct {
	Type ClauseTyp
	Syn  LexItm
	Adj  LexItm
}

func (cp ClauseProps) IsSet() bool {
	return cp.Type != ClauseUnk
}

type GetObjTyp int

const (
	GetUndef GetObjTyp = iota
	GetAll
	GetOne
	GetInhibit
)

type ParseProps struct {
	ObjOrClauseCnt int
	EndOnPrep      bool
	Number         int
	ShldOrphan     bool
	HasMerged      bool
	HasAnd         bool
	ItObj          *Object
	GetType        GetObjTyp
	AdjClause      ClauseProps
	Buts           []*Object
	OneObj         LexItm
	Continue       int
	InQuotes       bool
	BufLen         int
	WalkDir        string
}

const (
	ContEmpty int = -1
)

type FindTyp int

const (
	FindAll FindTyp = iota
	FindTop
	FindBottom
)

type ReserveProps struct {
	Idx    int
	IdxSet bool
	Buf    []LexItm
}

type OopsProps struct {
	Unk    int
	UnkSet bool
	Idx    int
}

type AgainProps struct {
	Buf []LexItm
	Dir string
}

var (
	DirObj            *Object
	IndirObj          *Object
	ActVerb           string
	DirObjPossibles   []*Object
	IndirObjPossibles []*Object
	Winner            *Object
	Here              *Object
	AlwaysLit         bool
	Search            FindProps
	ParsedSyntx       ParseTbl
	OrphanedSyntx     ParseTbl
	Params            ParseProps
	DetectedSyntx     *Syntx
	NotHere           NotHereProps
	LexRes            []LexItm
	Reserv            ReserveProps
	Again             AgainProps
	Oops              OopsProps
)

func Parse() bool {
	if Params.ShldOrphan {
		OrphanedSyntx.Set(ParsedSyntx)
	}
	bakWin := Winner
	bakMerg := Params.HasMerged
	Params.HasMerged = false
	Params.EndOnPrep = false
	Params.Buts = []*Object{}
	DirObjPossibles = []*Object{}
	IndirObjPossibles = []*Object{}
	if !Params.InQuotes && Winner != Player {
		Winner = Player
		Here = MetaLoc(Player)
		Lit = IsLit(Here, true)
	}
	beg := 0
	if Reserv.IdxSet {
		beg = Reserv.Idx
		LexRes = append([]LexItm{}, Reserv.Buf...)
		if !SuperBrief && Player == Winner {
			NewLine()
		}
		Reserv.IdxSet = false
		Reserv.Buf = nil
		Params.Continue = ContEmpty
	} else if Params.Continue != ContEmpty {
		beg = Params.Continue
		if !SuperBrief && Player == Winner && ActVerb != "say" {
			NewLine()
		}
		Params.Continue = ContEmpty
	} else {
		Winner = Player
		Params.InQuotes = false
		if l := Winner.Location(); !l.Has(FlgVeh) {
			Here = l
		}
		Lit = IsLit(Here, true)
		if !SuperBrief {
			NewLine()
		}
		Print(">", NoNewline)
		_, LexRes = Read()
	}
	Params.BufLen = len(LexRes)
	if Params.BufLen == 0 {
		Print("I beg your pardon?", Newline)
		return false
	}
	if LexRes[beg].Is("oops") {
		if Params.BufLen > beg+1 && LexRes[beg+1].IsAny(".", ",") {
			beg++
			Params.BufLen--
		}
		if Params.BufLen <= 1 {
			Print("I can't help your clumsiness.", Newline)
			return false
		}
		if Oops.UnkSet {
			if Params.BufLen > beg+1 && LexRes[beg+1].Is("\"") {
				Print("Sorry, you can't correct mistakes in quoted text.", Newline)
				return false
			}
			if Params.BufLen > beg+2 {
				Print("Warning: only the first word after OOPS is used.", Newline)
			}
			Again.Buf[Oops.Unk].Set(LexRes[beg+1])
			Winner = bakWin
			LexRes = append([]LexItm{}, Again.Buf...)
			Params.BufLen = len(LexRes)
			beg = Oops.Idx
		} else {
			Print("There was no word to replace!", Newline)
			return false
		}
	} else if !LexRes[beg].IsAny("again", "g") {
		Params.Number = 0
	}
	var dir string
	if LexRes[beg].IsAny("again", "g") {
		if len(Again.Buf) == 0 {
			Print("Beg pardon?", Newline)
			return false
		}
		if Params.ShldOrphan {
			Print("It's difficult to repeat fragments.", Newline)
			return false
		}
		if !ParserOk {
			Print("That would just repeat a mistake.", Newline)
			return false
		}
		tmpLen := len(LexRes)
		if Params.BufLen > beg+1 {
			if !LexRes[beg+1].IsAny(".", ",", "then", "and") {
				Print("I couldn't understand that sentence.", Newline)
				return false
			}
			beg += 2
			tmpLen -= 2
		} else {
			beg++
			tmpLen--
		}
		if tmpLen > 0 {
			Reserv.Idx = beg
			Reserv.IdxSet = true
			Reserv.Buf = append([]LexItm{}, LexRes...)
		} else {
			Reserv.IdxSet = false
		}
		Winner = bakWin
		Params.HasMerged = bakMerg
		LexRes = append([]LexItm{}, Again.Buf...)
		dir = Again.Dir
		ParsedSyntx.Set(OrphanedSyntx)
	} else {
		Again.Buf = append([]LexItm{}, LexRes...)
		Oops.Idx = beg
		Reserv.IdxSet = false
		Params.ObjOrClauseCnt = 0
		Params.GetType = GetUndef
		ln := Params.BufLen
		var lw, nw LexItm
		var vrb string
		var isOf bool
		Params.BufLen--
		for i := beg; Params.BufLen > -1; i, Params.BufLen = i+1, Params.BufLen-1 {
			HandleNumber(i)
			wrd := LexRes[i]
			if wrd.Type == WordUnk && !wrd.Is("intnum") {
				UnknownWord(i)
				return false
			}
			if Params.BufLen == 0 {
				nw.Clear()
			} else {
				nw.Set(LexRes[i+1])
			}
			if wrd.Is("to") && vrb == "tell" {
				wrd.Set(mkBuzz("\""))
			} else if wrd.Is("then") && Params.BufLen > 0 && len(vrb) == 0 && !Params.InQuotes {
				if !lw.IsSet() || lw.Is(".") {
					wrd.Set(mkBuzz("the"))
				} else {
					ParsedSyntx.Verb.Norm = "tell"
					ParsedSyntx.Verb.Orig = "tell"
					ParsedSyntx.Verb.Type = WordVerb
					wrd.Set(mkBuzz("\""))
				}
			}
			if wrd.IsAny("then", ".", "\"") {
				if wrd.Is("\"") {
					Params.InQuotes = !Params.InQuotes
				}
				if Params.BufLen != 0 {
					Params.Continue = i + 1
				}
				break
			} else if wrd.Type == WordDir &&
				(len(vrb) == 0 || vrb == "walk") &&
				(ln == 1 ||
					(ln == 2 && vrb == "walk") ||
					(nw.IsAny("then", ".", "\"") && ln >= 2) ||
					(Params.InQuotes && ln == 2 && nw.Is("\"")) ||
					(ln > 2 && nw.IsAny(",", "and"))) {
				dir = wrd.Norm
				if nw.IsAny(",", "and") {
					LexRes[i+1].Set(mkBuzz("then"))
				}
				if ln <= 2 {
					Params.InQuotes = false
					break
				}
			} else if wrd.Type == WordVerb && len(vrb) == 0 {
				vrb = wrd.Norm
				ParsedSyntx.Verb.Set(wrd)
			} else if wrd.Type == WordPrep || wrd.IsAny("all", "one") || wrd.Type == WordAdj || wrd.Type == WordObj {
				if Params.BufLen > 1 && nw.Is("of") && wrd.Type != WordPrep && !wrd.IsAny("all", "one", "a") {
					isOf = true
				} else if wrd.Type == WordPrep && (Params.BufLen == 0 || nw.IsAny("then", ".")) {
					Params.EndOnPrep = true
					if Params.ObjOrClauseCnt < 2 {
						ParsedSyntx.Prep1.Set(wrd)
					}
				} else if Params.ObjOrClauseCnt == 2 {
					Print("There were too many nouns in that sentence.", Newline)
					return false
				} else {
					Params.ObjOrClauseCnt++
					ok, i := Clause(i, wrd)
					if !ok {
						return false
					}
					if i < 0 {
						Params.InQuotes = false
						break
					}
				}
			} else if wrd.Is("of") {
				if !isOf || nw.IsAny(".", "then") {
					CantUse(i)
					return false
				}
				isOf = false
			} else if wrd.Type == WordBuzz {
				lw.Set(wrd)
				continue
			} else if vrb == "tell" && wrd.Type == WordVerb && Winner == Player {
				Print("Please consult your manual for the correct way to talk to other people or creatures.", Newline)
				return false
			} else {
				CantUse(i)
				return false
			}
			lw.Set(wrd)
		}
		if Params.BufLen < 1 {
			Params.InQuotes = false
		}
	}
	Oops.UnkSet = false
	if len(dir) != 0 {
		ActVerb = "walk"
		DirObj = ToDirObj(dir)
		Params.ShldOrphan = false
		Params.WalkDir = dir
		Again.Dir = dir
	} else {
		if Params.ShldOrphan {
			OrphanMerge()
		}
		Params.WalkDir = ""
		Again.Dir = ""
		if !SyntaxCheck() {
			return false
		}
		if !SnarfObjects() {
			return false
		}
		if !ManyCheck() {
			return false
		}
		if !TakeCheck() {
			return false
		}
	}
	return true
}

func mkBuzz(wrd string) LexItm {
	return LexItm{
		Norm: wrd,
		Orig: wrd,
		Type: WordBuzz,
	}
}

func Clause(idx int, wrd LexItm) (bool, int) {
	if wrd.Type == WordPrep {
		if Params.ObjOrClauseCnt == 1 {
			ParsedSyntx.Prep1.Set(wrd)
		} else if Params.ObjOrClauseCnt == 2 {
			ParsedSyntx.Prep2.Set(wrd)
		}
		idx++
	} else {
		Params.BufLen++
	}
	if Params.BufLen == 0 {
		Params.ObjOrClauseCnt--
		return true, -1
	}
	cpyStart := idx
	if LexRes[idx].IsAny("the", "a", "an") {
		cpyStart++
	}
	var lw, nw LexItm
	isFirst := true
	var isAnd bool
	var i int
	Params.BufLen--
	for i = idx; Params.BufLen > -1; i, Params.BufLen = i+1, Params.BufLen-1 {
		HandleNumber(i)
		cw := LexRes[i]
		if cw.Type == WordUnk && !cw.Is("intnum") {
			UnknownWord(i)
			return false, -1
		}
		if Params.BufLen == 0 {
			nw.Clear()
		} else {
			nw.Set(LexRes[i+1])
		}
		if cw.IsAny("and", ",") {
			isAnd = true
		} else if cw.IsAny("and", ",") {
			if nw.Is("of") {
				Params.BufLen--
				i++
			}
		} else if cw.IsAny("then", ".") || (cw.Type == WordPrep && ParsedSyntx.Verb.IsSet() && !isFirst) {
			Params.BufLen++
			if Params.ObjOrClauseCnt == 1 {
				ParsedSyntx.ObjOrClause1 = append([]LexItm{}, LexRes[cpyStart:i]...)
			} else if Params.ObjOrClauseCnt == 2 {
				ParsedSyntx.ObjOrClause2 = append([]LexItm{}, LexRes[cpyStart:i]...)
			}
			return true, i - 1
		} else if cw.Type == WordObj {
			if Params.BufLen > 0 && nw.Is("of") && !cw.IsAny("all", "one") {
				lw.Set(cw)
				isFirst = false
				continue
			}
			// TODO: I'm not sure how this can ever be true.
			if cw.Type == WordAdj && nw.IsSet() && nw.Type == WordObj {
				lw.Set(cw)
				isFirst = false
				continue
			}
			if !isAnd && !nw.IsAny("but", "except", "and", ",") {
				if Params.ObjOrClauseCnt == 1 {
					ParsedSyntx.ObjOrClause1 = append([]LexItm{}, LexRes[cpyStart:i+1]...)
				} else if Params.ObjOrClauseCnt == 2 {
					ParsedSyntx.ObjOrClause2 = append([]LexItm{}, LexRes[cpyStart:i+1]...)
				}
				return true, i
			}
			isAnd = true
		} else if (Params.HasMerged || Params.ShldOrphan || ParsedSyntx.Verb.IsSet()) &&
			(cw.Type == WordAdj || cw.Type == WordBuzz) {
			lw.Set(cw)
			isFirst = false
			continue
		} else if isAnd && (cw.Type == WordDir || cw.Type == WordVerb) {
			i -= 2
			Params.BufLen += 2
			LexRes[i].Set(mkBuzz("then"))
		} else if cw.Type == WordPrep {
			lw.Set(cw)
			isFirst = false
			continue
		} else {
			CantUse(i)
			return false, -1
		}
		lw.Set(cw)
		isFirst = false
	}
	if Params.ObjOrClauseCnt == 1 {
		ParsedSyntx.ObjOrClause1 = append([]LexItm{}, LexRes[cpyStart:i]...)
	} else if Params.ObjOrClauseCnt == 2 {
		ParsedSyntx.ObjOrClause2 = append([]LexItm{}, LexRes[cpyStart:i]...)
	}
	return true, -1
}

func UnknownWord(idx int) {
	Oops.UnkSet = true
	Oops.Unk = idx
	if ActVerb == "say" {
		Print("Nothing happens.", Newline)
		return
	}
	Print("I don't know the word \"", NoNewline)
	Print(LexRes[idx].Orig, NoNewline)
	Print("\".", Newline)
	Params.InQuotes = false
	Params.ShldOrphan = false
}

func CantUse(idx int) {
	if ActVerb == "say" {
		Print("Nothing happens.", Newline)
		return
	}
	Print("You used the word \"", NoNewline)
	Print(LexRes[idx].Orig, NoNewline)
	Print("\" in a way that I don't understand.", Newline)
	Params.InQuotes = false
	Params.ShldOrphan = false
}

// HandleNumber converts the lex item pointed
// to by idx into a number if possible.
func HandleNumber(idx int) {
	wrd := LexRes[idx]
	var tim, sum int
	for _, chr := range wrd.Orig {
		if chr == 58 {
			tim = sum
			sum = 0
		} else if sum > 10000 {
			return
		} else if chr < 58 && chr > 47 {
			sum = (sum * 10) + (int(chr) - 48)
		} else {
			return
		}
	}
	LexRes[idx].Norm = "intnum"
	if sum > 1000 {
		return
	} else if tim != 0 {
		if tim < 8 {
			tim += 12
		} else if tim > 23 {
			return
		}
		sum = tim*60 + sum
	}
	Params.Number = sum
}

func OrphanMerge() {
	Params.ShldOrphan = false
	isAdj := false
	if ParsedSyntx.Verb.Type == OrphanedSyntx.Verb.Type || ParsedSyntx.Verb.Type == WordAdj {
		isAdj = true
	} else if ParsedSyntx.Verb.Type == WordObj && Params.ObjOrClauseCnt == 0 {
		ParsedSyntx.Verb.Clear()
		ParsedSyntx.ObjOrClause1 = []LexItm{LexRes[0], LexRes[1]}
		Params.ObjOrClauseCnt = 1
	}
	if ParsedSyntx.Verb.IsSet() && !isAdj && !ParsedSyntx.Verb.Matches(OrphanedSyntx.Verb) {
		return
	}
	if Params.ObjOrClauseCnt == 2 {
		return
	}
	if OrphanedSyntx.ObjOrClause1 != nil && len(OrphanedSyntx.ObjOrClause1) == 0 {
		if !ParsedSyntx.Prep1.Matches(OrphanedSyntx.Prep1) && ParsedSyntx.Prep1.IsSet() {
			return
		}
		if isAdj {
			if Params.ObjOrClauseCnt == 0 {
				Params.ObjOrClauseCnt = 1
			}
			if len(ParsedSyntx.ObjOrClause1) == 0 {
				OrphanedSyntx.ObjOrClause1 = []LexItm{LexRes[0], LexRes[1]}
			} else {
				OrphanedSyntx.ObjOrClause1 = LexRes[0:ParsedSyntx.Obj1End]
			}
		} else {
			OrphanedSyntx.ObjOrClause1 = append([]LexItm{}, ParsedSyntx.ObjOrClause1...)
		}
	} else if OrphanedSyntx.ObjOrClause2 != nil && len(OrphanedSyntx.ObjOrClause2) == 0 {
		if !ParsedSyntx.Prep1.Matches(OrphanedSyntx.Prep2) && ParsedSyntx.Prep1.IsSet() {
			return
		}
		if isAdj {
			if len(ParsedSyntx.ObjOrClause1) == 0 {
				OrphanedSyntx.ObjOrClause2 = []LexItm{LexRes[0], LexRes[1]}
			} else {
				OrphanedSyntx.ObjOrClause2 = LexRes[0:ParsedSyntx.Obj1End]
			}
		} else {
			OrphanedSyntx.ObjOrClause2 = append([]LexItm{}, ParsedSyntx.ObjOrClause1...)
		}
		Params.ObjOrClauseCnt = 2
	} else if Params.AdjClause.Type != ClauseUnk {
		if Params.ObjOrClauseCnt != 1 && !isAdj {
			Params.AdjClause.Type = ClauseUnk
			return
		}
		beg := ParsedSyntx.Obj1Start
		if isAdj {
			beg = 0
			isAdj = false
		}
		var adj LexItm
		if ParsedSyntx.Obj1End == 0 {
			ParsedSyntx.Obj1End = 1
		}
		broken := false
		for i := beg; i < ParsedSyntx.Obj1End; i++ {
			wrd := LexRes[i]
			if !isAdj && (wrd.Type == WordAdj || wrd.IsAny("all", "one")) {
				adj.Set(wrd)
			} else if wrd.Is("one") {
				AclauseWin(adj)
				broken = true
				break
			} else if wrd.Type == WordObj {
				if wrd.Matches(Params.AdjClause.Syn) {
					AclauseWin(adj)
				} else {
					NclauseWin()
				}
				broken = true
				break
			}
		}
		if !broken {
			if ParsedSyntx.Obj1End == 1 {
				ParsedSyntx.ObjOrClause1 = []LexItm{LexRes[0]}
				Params.ObjOrClauseCnt = 1
			}
			if !adj.IsSet() {
				Params.AdjClause.Type = ClauseUnk
				return
			}
			AclauseWin(adj)
		}
	}
	ParsedSyntx.Set(OrphanedSyntx)
	Params.HasMerged = true
}

func AclauseWin(adj LexItm) {
	ParsedSyntx.Verb.Set(OrphanedSyntx.Verb)
	var tbl *[]LexItm = &[]LexItm{}
	if Params.AdjClause.Type == Clause1 {
		tbl = &OrphanedSyntx.ObjOrClause1
	} else if Params.AdjClause.Type == Clause2 {
		tbl = &OrphanedSyntx.ObjOrClause2
	}
	for idx, obj := range *tbl {
		if obj.Matches(Params.AdjClause.Adj) {
			*tbl = append((*tbl)[0:idx], append([]LexItm{adj}, (*tbl)[idx:len(*tbl)]...)...)
			break
		}
	}
	if OrphanedSyntx.ObjOrClause2 != nil {
		Params.ObjOrClauseCnt = 2
	}
	Params.AdjClause.Type = ClauseUnk
}

func NclauseWin() {
	if Params.AdjClause.Type == Clause1 {
		OrphanedSyntx.ObjOrClause1 = append([]LexItm{}, ParsedSyntx.ObjOrClause1...)
		OrphanedSyntx.Obj1Start = ParsedSyntx.Obj1Start
		OrphanedSyntx.Obj1End = ParsedSyntx.Obj1End
	} else if Params.AdjClause.Type == Clause2 {
		OrphanedSyntx.ObjOrClause2 = append([]LexItm{}, ParsedSyntx.ObjOrClause1...)
	}
	if OrphanedSyntx.ObjOrClause2 != nil {
		Params.ObjOrClauseCnt = 2
	}
	Params.AdjClause.Type = ClauseUnk
}

func TakeCheck() bool {
	if !ITakeCheck(DirObjPossibles, DetectedSyntx.Obj1.LocFlags) {
		return false
	}
	return ITakeCheck(IndirObjPossibles, DetectedSyntx.Obj2.LocFlags)
}

func ITakeCheck(tbl []*Object, ibits LocFlags) bool {
	if len(tbl) == 0 || (!LocHave.In(ibits) && !LocTake.In(ibits)) {
		return true
	}
	for _, obj := range tbl {
		if obj == &It {
			if !IsAccessible(Params.ItObj) {
				Print("I don't see what you're referring to.", Newline)
				return false
			}
			obj = Params.ItObj
		}
		var taken bool
		if !IsHeld(obj) && obj != &Hands && obj != &Me {
			DirObjPossibles = []*Object{obj}
			if obj.Has(FlgTryTake) {
				taken = true
			} else if Winner != &Adventurer {
				taken = false
			} else if LocTake.In(ibits) && ITake(false) {
				taken = false
			} else {
				taken = true
			}
			if taken && LocHave.In(ibits) && Winner == &Adventurer {
				if obj == &NotHereObject {
					Print("You don't have that!", Newline)
					return false
				}
				Print("You don't have the ", NoNewline)
				PrintObject(obj)
				Print(".", Newline)
				return false
			}
			if !taken && Winner == &Adventurer {
				Print("(Taken)", Newline)
			}
		}
	}
	return true
}

func ManyCheck() bool {
	loss := 0
	if len(DirObjPossibles) > 1 && !LocMany.In(DetectedSyntx.Obj1.LocFlags) {
		loss = 1
	} else if len(IndirObjPossibles) > 1 && !LocMany.In(DetectedSyntx.Obj2.LocFlags) {
		loss = 2
	}
	if loss == 0 {
		return true
	}
	Print("You can't use multiple ", NoNewline)
	if loss == 2 {
		Print("in", NoNewline)
	}
	Print("direct objects with \"", NoNewline)
	if !ParsedSyntx.Verb.IsSet() {
		Print("tell", NoNewline)
	} else if Params.ShldOrphan || Params.HasMerged {
		Print(ParsedSyntx.Verb.Norm, NoNewline)
	} else {
		Print(ParsedSyntx.Verb.Orig, NoNewline)
	}
	Print("\".", Newline)
	return false
}

func SnarfObjects() bool {
	Params.Buts = []*Object{}
	if ParsedSyntx.ObjOrClause2 != nil {
		Search.LocFlags = DetectedSyntx.Obj2.LocFlags
		res := Snarfem(false, ParsedSyntx.ObjOrClause2)
		if res == nil {
			return false
		}
		IndirObjPossibles = append(IndirObjPossibles, res...)
	}
	if ParsedSyntx.ObjOrClause1 != nil {
		Search.LocFlags = DetectedSyntx.Obj1.LocFlags
		res := Snarfem(true, ParsedSyntx.ObjOrClause1)
		if res == nil {
			return false
		}
		DirObjPossibles = append(DirObjPossibles, res...)
	}
	if len(Params.Buts) != 0 {
		l := len(DirObjPossibles)
		if len(ParsedSyntx.ObjOrClause1) != 0 {
			DirObjPossibles = ButMerge(DirObjPossibles)
		}
		if len(ParsedSyntx.ObjOrClause2) != 0 && (len(ParsedSyntx.ObjOrClause1) == 0 || l == len(DirObjPossibles)) {
			IndirObjPossibles = ButMerge(IndirObjPossibles)
		}
	}
	return true
}

func ButMerge(tbl []*Object) []*Object {
	res := []*Object{}
	for _, bts := range Params.Buts {
		for _, obj := range tbl {
			if obj == bts {
				res = append(res, obj)
			}
		}
	}
	return res
}

func Snarfem(isDirect bool, wrds []LexItm) []*Object {
	Params.HasAnd = false
	wasall := false
	if Params.GetType == GetAll {
		wasall = true
	}
	Search.ObjFlags = nil
	res := []*Object{}
	var but *[]*Object
	var nw LexItm
	for idx, wrd := range wrds {
		if idx != len(wrds)-1 {
			nw.Set(wrds[idx+1])
		} else {
			nw.Clear()
		}
		if wrd.Is("all") {
			Params.GetType = GetAll
			if nw.Is("of") {
				continue
			}
		} else if wrd.IsAny("but", "except") {
			out := GetObject(isDirect, true)
			if out == nil {
				return nil
			}
			if but != nil {
				*but = append(*but, out...)
			} else {
				res = append(res, out...)
			}
			but = &Params.Buts
			*but = []*Object{}
		} else if wrd.IsAny("a", "one") {
			if !Search.Adj.IsSet() {
				Params.GetType = GetOne
				if nw.Is("of") {
					continue
				}
			} else {
				Search.Syn.Set(Params.OneObj)
				out := GetObject(isDirect, true)
				if out == nil {
					return nil
				}
				if but != nil {
					*but = append(*but, out...)
				} else {
					res = append(res, out...)
				}
				if idx == len(wrds)-1 {
					return res
				}
			}
		} else if wrd.IsAny("and", ",") && !nw.IsAny("and", ",") {
			Params.HasAnd = true
			out := GetObject(isDirect, true)
			if out == nil {
				return nil
			}
			if but != nil {
				*but = append(*but, out...)
			} else {
				res = append(res, out...)
			}
			continue
		} else if wrd.Type == WordBuzz {
			continue
		} else if wrd.IsAny("and", ",") {
			continue
		} else if wrd.Is("of") {
			if Params.GetType == GetUndef {
				Params.GetType = GetInhibit
			}
		} else if wrd.Type == WordAdj && !Search.Adj.IsSet() {
			Search.Adj.Set(wrd)
		} else if wrd.Type == WordObj {
			Search.Syn.Set(wrd)
			Params.OneObj.Set(wrd)
		}
	}
	out := GetObject(isDirect, true)
	if wasall {
		Params.GetType = GetAll
	}
	if out == nil {
		return nil
	}
	if but != nil {
		*but = append(*but, out...)
	} else {
		res = append(res, out...)
	}
	return res
}

// SyntaxCheck tries to find a matching syntax based on
// the parsed syntax.
func SyntaxCheck() bool {
	if !ParsedSyntx.Verb.IsSet() {
		Print("There was no verb in that sentence!", Newline)
		return false
	}
	var findFirst, findSecond *Syntx
	for _, syn := range Commands {
		if !ParsedSyntx.Verb.Is(syn.Verb) {
			continue
		}
		if Params.ObjOrClauseCnt > syn.NumObjects() {
			continue
		}
		if syn.NumObjects() >= 1 && Params.ObjOrClauseCnt == 0 && (!ParsedSyntx.Prep1.IsSet() || syn.IsVrbPrep(ParsedSyntx.Prep1.Norm)) {
			findFirst = &syn
		} else if syn.IsVrbPrep(ParsedSyntx.Prep1.Norm) {
			if syn.NumObjects() == 2 && Params.ObjOrClauseCnt == 1 {
				findSecond = &syn
			} else if syn.IsObjPrep(ParsedSyntx.Prep2.Norm) {
				DetectedSyntx = &syn
				ActVerb = syn.GetActionVerb()
				return true
			}
		}
	}
	if findFirst == nil && findSecond == nil {
		Print("That sentence isn't one I recognize.", Newline)
		return false
	}
	found := false
	if findFirst != nil {
		obj := FindWhatIMean(findFirst.Obj1.ObjFlags, findFirst.Obj1.LocFlags, findFirst.VrbPrep)
		if obj != nil {
			DirObjPossibles = []*Object{obj}
			DetectedSyntx = findFirst
			ActVerb = findFirst.GetActionVerb()
			found = true
		}
	}
	if findSecond != nil && !found {
		obj := FindWhatIMean(findSecond.Obj2.ObjFlags, findSecond.Obj2.LocFlags, findSecond.ObjPrep)
		if obj != nil {
			IndirObjPossibles = []*Object{obj}
			DetectedSyntx = findSecond
			ActVerb = findSecond.GetActionVerb()
			found = true
		}
	}
	if ParsedSyntx.Verb.Is("find") && !found {
		Print("That question can't be answered.", Newline)
		return false
	}
	if Winner != Player && !found {
		CanNotOrphan()
		return false
	}
	if !found {
		Orphan(findFirst, findSecond)
		Print("What do you want to ", NoNewline)
		if !OrphanedSyntx.Verb.IsSet() {
			Print("tell", NoNewline)
		} else {
			Print(OrphanedSyntx.Verb.Orig, NoNewline)
		}
		if findSecond != nil {
			Print(" ", NoNewline)
			ThingPrint(true, true)
		}
		if findFirst != nil {
			Print(" "+findFirst.VrbPrep, NoNewline)
		} else if findSecond != nil {
			Print(" "+findSecond.ObjPrep, NoNewline)
		}
		Print("?", Newline)
		return false
	}
	return true
}

func Orphan(first, second *Syntx) {
	OrphanedSyntx.Set(ParsedSyntx)
	if Params.ObjOrClauseCnt < 2 {
		OrphanedSyntx.ObjOrClause2 = []LexItm{}
	}
	if Params.ObjOrClauseCnt < 1 {
		OrphanedSyntx.ObjOrClause1 = []LexItm{}
	}
	if first != nil {
		OrphanedSyntx.Prep1.Norm = first.VrbPrep
		OrphanedSyntx.Prep1.Orig = first.VrbPrep
		OrphanedSyntx.Prep1.Type = WordPrep
		OrphanedSyntx.ObjOrClause1 = []LexItm{}
	} else if second != nil {
		OrphanedSyntx.Prep2.Norm = second.ObjPrep
		OrphanedSyntx.Prep2.Orig = second.ObjPrep
		OrphanedSyntx.Prep2.Type = WordPrep
		OrphanedSyntx.ObjOrClause2 = []LexItm{}
	}
}

func CanNotOrphan() {
	Print("\"I don't understand! What are you referring to?\"", Newline)
}

func FindWhatIMean(objFlags []Flag, locFlags LocFlags, prep string) *Object {
	if FlgKludge.In(objFlags) {
		return &Rooms
	}
	Search.ObjFlags = objFlags
	Search.LocFlags = locFlags
	res := GetObject(false, false)
	Search.ObjFlags = nil
	if len(res) != 1 {
		return nil
	}
	Print("(", NoNewline)
	if len(prep) == 0 || Params.EndOnPrep {
		PrintObject(res[0])
		Print(")", Newline)
		return res[0]
	}
	Print(prep, NoNewline)
	if prep == "out" {
		Print(" of", NoNewline)
	}
	if res[0] == &Hands {
		Print(" your hands", NoNewline)
	} else {
		Print(" the ", NoNewline)
		PrintObject(res[0])
	}
	Print(")", Newline)
	return res[0]
}

func GetObject(isDirect, vrb bool) []*Object {
	if Params.GetType == GetInhibit {
		return []*Object{}
	}
	if !Search.Syn.IsSet() && Search.Adj.IsSet() && Search.Adj.Type == WordObj {
		Search.Syn.Set(Search.Adj)
		Search.Adj.Clear()
	}
	if !Search.Syn.IsSet() && !Search.Adj.IsSet() && Params.GetType != GetAll && len(Search.ObjFlags) == 0 {
		if vrb {
			Print("There seems to be a noun missing in that sentence!", Newline)
		}
		return nil
	}
	xbits := Search.LocFlags
	if Params.GetType != GetAll || len(Search.LocFlags) == 0 {
		Search.LocFlags.All()
	}
	res := []*Object{}
	gcheck := false
	olen := 0
	for {
		if gcheck {
			res = append(res, GlobalCheck()...)
		} else {
			if Lit {
				Player.Take(FlgTrans)
				res = append(res, DoSL(Here, LocOnGrnd, LocInRoom)...)
				Player.Give(FlgTrans)
			}
			res = append(res, DoSL(Player, LocHeld, LocCarried)...)
		}
		ln := len(res)
		if Params.GetType == GetAll {
			Search.LocFlags = xbits
			Search.Syn.Clear()
			Search.Adj.Clear()
			return res
		}
		if Params.GetType == GetOne && ln != 0 {
			if ln > 1 {
				res = []*Object{res[rand.Intn(len(res))]}
				Print("(How about the ", NoNewline)
				PrintObject(res[0])
				Print("?)", Newline)
			}
		} else if ln > 1 || (ln == 0 && Search.LocFlags.HasAll()) {
			if Search.LocFlags.HasAll() {
				Search.LocFlags = xbits
				olen = ln
				res = []*Object{}
				continue
			}
			if ln == 0 {
				ln = olen
			}
			if Winner != Player {
				CanNotOrphan()
				return nil
			}
			if vrb && Search.Syn.IsSet() {
				WhichPrint(isDirect, res)
				if isDirect {
					Params.AdjClause.Type = Clause1
				} else {
					Params.AdjClause.Type = Clause2
				}
				Params.AdjClause.Syn.Set(Search.Syn)
				Params.AdjClause.Adj.Set(Search.Adj)
				Orphan(nil, nil)
				Params.ShldOrphan = true
			} else if vrb {
				Print("There seems to be a noun missing in that sentence!", Newline)
			}
			Search.Syn.Clear()
			Search.Adj.Clear()
			return nil
		}
		if ln == 0 && gcheck {
			if vrb {
				Search.LocFlags = xbits
				if Lit || ActVerb == "tell" {
					res = append(res, &NotHereObject)
					NotHere.Syn.Set(Search.Syn)
					NotHere.Adj.Set(Search.Adj)
					Search.Syn.Clear()
					Search.Adj.Clear()
					return res
				}
				Print("It's too dark to see!", Newline)
			}
			Search.Syn.Clear()
			Search.Adj.Clear()
			return nil
		}
		if ln == 0 {
			gcheck = true
			continue
		}
		Search.LocFlags = xbits
		Search.Syn.Clear()
		Search.Adj.Clear()
		return res
	}
}

// WhichPrint outputs all of the possible matches
// when the game parser matches multiple game objects.
func WhichPrint(isDirect bool, tbl []*Object) {
	Print("Which ", NoNewline)
	if Params.ShldOrphan || Params.HasMerged || Params.HasAnd {
		if Search.Syn.IsSet() {
			Print(Search.Syn.Norm, NoNewline)
		} else if Search.Adj.IsSet() {
			Print(Search.Adj.Norm, NoNewline)
		} else {
			Print("one", NoNewline)
		}
	} else {
		ThingPrint(isDirect, false)
	}
	Print(" do you mean, ", NoNewline)
	ln := len(tbl)
	for _, obj := range tbl {
		Print("the ", NoNewline)
		PrintObject(obj)
		if ln == 2 {
			if len(tbl) != 2 {
				Print(",", NoNewline)
			}
			Print(" or ", NoNewline)
		} else if ln > 2 {
			Print(", ", NoNewline)
		}
		ln--
		if ln < 1 {
			Print("?", Newline)
		}
	}
}

// GlobalCheck looks through global objects if any match
// the parameters defined in the Search global variable
func GlobalCheck() []*Object {
	res := []*Object{}
	if len(Here.Global) != 0 {
		for _, obj := range Here.Global {
			if IsThisIt(obj) {
				res = append(res, obj)
			}
		}
	}
	if len(Here.Pseudo) != 0 {
		for _, obj := range Here.Pseudo {
			if Search.Syn.Is(obj.Synonym) {
				PseudoObject.Action = obj.Action
				res = append(res, &PseudoObject)
				break
			}
		}
	}
	if len(res) == 0 {
		if g := SearchList(&GlobalObjects, FindAll); g != nil {
			res = append(res, g...)
		}
		if len(res) == 0 && (ActVerb == "look inside" || ActVerb == "search" || ActVerb == "examine") {
			if LocHave.In(Search.LocFlags) {
				if r := SearchList(&Rooms, FindAll); r != nil {
					res = append(res, r...)
				}
			}
		}
	}
	return res
}

func ThingPrint(isDirect, isThe bool) {
	nsp, isFirst := true, true
	pn := false
	search := &ParsedSyntx.ObjOrClause1
	if !isDirect {
		search = &ParsedSyntx.ObjOrClause2
	}
	for _, wrd := range *search {
		if wrd.Is(",") {
			Print(", ", NoNewline)
		} else if nsp {
			nsp = false
		} else {
			Print(" ", NoNewline)
		}
		if wrd.IsAny(".", ",") {
			nsp = true
		} else if wrd.Is("me") {
			PrintObject(&Me)
			pn = true
		} else if wrd.Is("intnum") {
			PrintNumber(Params.Number)
			pn = true
		} else {
			if isFirst && !pn && isThe {
				Print("the ", NoNewline)
			}
			if Params.ShldOrphan || Params.HasMerged {
				Print(wrd.Norm, NoNewline)
			} else if wrd.Is("it") && IsAccessible(Params.ItObj) {
				PrintObject(Params.ItObj)
			} else {
				Print(wrd.Orig, NoNewline)
			}
			isFirst = false
		}
	}
}

// IsLit checks if the current game room is lit.
func IsLit(room *Object, rmChk bool) bool {
	if AlwaysLit && Winner == Player {
		return true
	}
	Search.ObjFlags = []Flag{FlgOn}
	bak := Here
	Here = room
	if rmChk && room.Has(FlgOn) {
		Here = bak
		Search.ObjFlags = nil
		return true
	}
	Search.LocFlags = nil
	res := []*Object{}
	if bak == room {
		if nr := SearchList(Winner, FindAll); nr != nil {
			res = append(res, nr...)
		}
		if Winner != Player && Player.IsIn(room) {
			if nr := SearchList(Player, FindAll); nr != nil {
				res = append(res, nr...)
			}
		}
	}
	if nr := SearchList(room, FindAll); nr != nil {
		res = append(res, nr...)
	}
	Here = bak
	Search.ObjFlags = nil
	if len(res) > 0 {
		return true
	}
	return false
}

// DoSL performs a specific game object search
// based on the provided location flags and
// parameters defined in the Search global variable.
func DoSL(obj *Object, f1, f2 LocFlag) []*Object {
	if f1.In(Search.LocFlags) && f2.In(Search.LocFlags) {
		return SearchList(obj, FindAll)
	} else if f1.In(Search.LocFlags) {
		return SearchList(obj, FindTop)
	} else if f2.In(Search.LocFlags) {
		return SearchList(obj, FindBottom)
	}
	return []*Object{}
}

// SearchList traverses the game object tree down and looks
// for game objects that match search parameters defined
// in the Search global variable.
// It returns a list of all found game objects.
func SearchList(obj *Object, typ FindTyp) []*Object {
	if obj == nil || !obj.HasChildren() {
		return nil
	}
	found := []*Object{}
	for _, child := range obj.Children {
		if typ != FindBottom && child.Synonyms != nil && IsThisIt(child) {
			found = append(found, child)
		}
		if (typ != FindTop || child.Has(FlgSearch) || child.Has(FlgSurf)) &&
			child.HasChildren() &&
			(child.Has(FlgOpen) || child.Has(FlgTrans)) {
			var res []*Object
			if child.Has(FlgSurf) || child.Has(FlgSearch) {
				res = SearchList(child, FindAll)
			} else {
				res = SearchList(child, FindTop)
			}
			found = append(found, res...)
		}
	}
	return found
}

// IsThisIt checks if the provided game object matches search
// parameters defined the the Search global variable.
func IsThisIt(obj *Object) bool {
	if obj.Has(FlgInvis) {
		return false
	}
	if Search.Syn.IsSet() && !obj.Is(Search.Syn.Norm) {
		return false
	}
	if Search.Adj.IsSet() && (len(obj.Adjectives) == 0 || !obj.Is(Search.Adj.Norm)) {
		return false
	}
	if !AnyFlagIn(Search.ObjFlags, obj.Flags) {
		return false
	}
	return true
}

// IsAccessible returns if the game object can be touched by the current (winner) character.
// In most cases the current (winner) character is the player.
func IsAccessible(obj *Object) bool {
	if obj.Has(FlgInvis) {
		return false
	}
	l := obj.Location()
	if l == nil {
		return false
	}
	if l == &GlobalObjects {
		return true
	}
	if l == &LocalGlobals && IsInGlobal(obj, Here) {
		return true
	}
	ml := MetaLoc(obj)
	if ml != Here && Winner != nil && ml != Winner.Location() {
		return false
	}
	if l == Winner || l == Here || (Winner != nil && l == Winner.Location()) {
		return true
	}
	if obj.Has(FlgOpen) && IsAccessible(l) {
		return true
	}
	return false
}

// MetaLoc returns the game objects top most location in the
// object's hierarchy tree which can either be a room or
// the GlobalObjects.
func MetaLoc(obj *Object) *Object {
	for {
		if obj == nil {
			return nil
		}
		if obj.IsIn(&GlobalObjects) {
			return &GlobalObjects
		}
		if obj.IsIn(&Rooms) {
			return obj
		}
		obj = obj.Location()
	}
}
