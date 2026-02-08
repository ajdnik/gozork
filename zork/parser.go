package zork

const (
	NumUndef int = -1
)

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

func (pt *ParseTbl) Clear() {
	pt.Verb.Clear()
	pt.Prep1.Clear()
	pt.ObjOrClause1 = nil
	pt.Obj1Start = NumUndef
	pt.Obj1End = NumUndef
	pt.Prep2.Clear()
	pt.ObjOrClause2 = nil
}

type NotHereProps struct {
	Syn LexItm
	Adj LexItm
}

type FindProps struct {
	ObjFlags Flags
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
	WalkDir        Direction
	HasWalkDir     bool
}

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
	Buf    []LexItm
	Dir    Direction
	HasDir bool
}

type ActionVerb struct {
	Norm string
	Orig string
}


func Parse() bool {
	if G.Params.ShldOrphan {
		G.OrphanedSyntx.Set(G.ParsedSyntx)
	}
	G.ParsedSyntx.Clear()
	bakWin := G.Winner
	bakMerg := G.Params.HasMerged
	G.Params.HasMerged = false
	G.Params.EndOnPrep = false
	G.Params.Buts = []*Object{}
	G.DirObjPossibles = []*Object{}
	G.IndirObjPossibles = []*Object{}
	if !G.Params.InQuotes && G.Winner != G.Player {
		G.Winner = G.Player
		G.Here = MetaLoc(G.Player)
		G.Lit = IsLit(G.Here, true)
	}
	beg := 0
	if G.Reserv.IdxSet {
		beg = G.Reserv.Idx
		G.LexRes = append([]LexItm{}, G.Reserv.Buf...)
		if !G.SuperBrief && G.Player == G.Winner {
			NewLine()
		}
		G.Reserv.IdxSet = false
		G.Reserv.Buf = nil
		G.Params.Continue = NumUndef
	} else if G.Params.Continue != NumUndef {
		beg = G.Params.Continue
		if !G.SuperBrief && G.Player == G.Winner && G.ActVerb.Norm != "say" {
			NewLine()
		}
		G.Params.Continue = NumUndef
	} else {
		G.Winner = G.Player
		G.Params.InQuotes = false
		if l := G.Winner.Location(); !l.Has(FlgVeh) {
			G.Here = l
		}
		G.Lit = IsLit(G.Here, true)
		if !G.SuperBrief {
			NewLine()
		}
		Print(">", NoNewline)
		_, G.LexRes = Read()
	}
	G.Params.BufLen = len(G.LexRes)
	if G.Params.BufLen == 0 {
		Print("I beg your pardon?", Newline)
		return false
	}
	if G.LexRes[beg].Is("oops") {
		if G.Params.BufLen > beg+1 && G.LexRes[beg+1].IsAny(".", ",") {
			beg++
			G.Params.BufLen--
		}
		if G.Params.BufLen <= 1 {
			Print("I can't help your clumsiness.", Newline)
			return false
		}
		if G.Oops.UnkSet {
			if G.Params.BufLen > beg+1 && G.LexRes[beg+1].Is("\"") {
				Print("Sorry, you can't correct mistakes in quoted text.", Newline)
				return false
			}
			if G.Params.BufLen > beg+2 {
				Print("Warning: only the first word after OOPS is used.", Newline)
			}
			G.Again.Buf[G.Oops.Unk].Set(G.LexRes[beg+1])
			G.Winner = bakWin
			G.LexRes = append([]LexItm{}, G.Again.Buf...)
			G.Params.BufLen = len(G.LexRes)
			beg = G.Oops.Idx
		} else {
			Print("There was no word to replace!", Newline)
			return false
		}
	} else if !G.LexRes[beg].IsAny("again", "g") {
		G.Params.Number = 0
	}
	var dir Direction
	hasDir := false
	if G.LexRes[beg].IsAny("again", "g") {
		if len(G.Again.Buf) == 0 {
			Print("Beg pardon?", Newline)
			return false
		}
		if G.Params.ShldOrphan {
			Print("It's difficult to repeat fragments.", Newline)
			return false
		}
		if !G.ParserOk {
			Print("That would just repeat a mistake.", Newline)
			return false
		}
		tmpLen := len(G.LexRes)
		if G.Params.BufLen > beg+1 {
			if !G.LexRes[beg+1].IsAny(".", ",", "then", "and") {
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
			G.Reserv.Idx = beg
			G.Reserv.IdxSet = true
			G.Reserv.Buf = append([]LexItm{}, G.LexRes...)
		} else {
			G.Reserv.IdxSet = false
		}
		G.Winner = bakWin
		G.Params.HasMerged = bakMerg
		G.LexRes = append([]LexItm{}, G.Again.Buf...)
		dir = G.Again.Dir
		hasDir = G.Again.HasDir
		G.ParsedSyntx.Set(G.OrphanedSyntx)
	} else {
		G.Again.Buf = append([]LexItm{}, G.LexRes...)
		G.Oops.Idx = beg
		G.Reserv.IdxSet = false
		G.Params.ObjOrClauseCnt = 0
		G.Params.GetType = GetUndef
		ln := G.Params.BufLen
		var lw, nw LexItm
		var vrb string
		var isOf bool
		G.Params.BufLen--
		for i := beg; G.Params.BufLen > -1; i, G.Params.BufLen = i+1, G.Params.BufLen-1 {
			HandleNumber(i)
			wrd := G.LexRes[i]
			if wrd.Types == nil && !wrd.Is("intnum") {
				UnknownWord(i)
				return false
			}
			if G.Params.BufLen == 0 {
				nw.Clear()
			} else {
				nw.Set(G.LexRes[i+1])
			}
			if wrd.Is("to") && vrb == "tell" {
				wrd.Set(mkBuzz("\""))
			} else if wrd.Is("then") && G.Params.BufLen > 0 && len(vrb) == 0 && !G.Params.InQuotes {
				if !lw.IsSet() || lw.Is(".") {
					wrd.Set(mkBuzz("the"))
				} else {
					G.ParsedSyntx.Verb.Norm = "tell"
					G.ParsedSyntx.Verb.Orig = "tell"
					G.ParsedSyntx.Verb.Types = WordTypes{WordVerb}
					wrd.Set(mkBuzz("\""))
				}
			}
			if wrd.IsAny("then", ".", "\"") {
				if wrd.Is("\"") {
					G.Params.InQuotes = !G.Params.InQuotes
				}
				if G.Params.BufLen != 0 {
					G.Params.Continue = i + 1
				}
				break
			} else if wrd.Types.Has(WordDir) &&
				(len(vrb) == 0 || vrb == "walk") &&
				(ln == 1 ||
					(ln == 2 && vrb == "walk") ||
					(nw.IsAny("then", ".", "\"") && ln >= 2) ||
					(G.Params.InQuotes && ln == 2 && nw.Is("\"")) ||
					(ln > 2 && nw.IsAny(",", "and"))) {
				dir, _ = StringToDir(wrd.Norm)
			hasDir = true
				if nw.IsAny(",", "and") {
					G.LexRes[i+1].Set(mkBuzz("then"))
				}
				if ln <= 2 {
					G.Params.InQuotes = false
					break
				}
			} else if wrd.Types.Has(WordVerb) && len(vrb) == 0 {
				vrb = wrd.Norm
				G.ParsedSyntx.Verb.Set(wrd)
			} else if wrd.Types.Has(WordPrep) || wrd.IsAny("all", "one") || wrd.Types.Has(WordAdj) || wrd.Types.Has(WordObj) {
				if G.Params.BufLen > 1 && nw.Is("of") && !wrd.Types.Has(WordPrep) && !wrd.IsAny("all", "one", "a") {
					isOf = true
				} else if wrd.Types.Has(WordPrep) && (G.Params.BufLen == 0 || nw.IsAny("then", ".")) {
					G.Params.EndOnPrep = true
					if G.Params.ObjOrClauseCnt < 2 {
						G.ParsedSyntx.Prep1.Set(wrd)
					}
				} else if G.Params.ObjOrClauseCnt == 2 {
					Print("There were too many nouns in that sentence.", Newline)
					return false
				} else {
					G.Params.ObjOrClauseCnt++
					ok, i := Clause(i, wrd)
					if !ok {
						return false
					}
					if i < 0 {
						G.Params.InQuotes = false
						break
					}
				}
			} else if wrd.Is("of") {
				if !isOf || nw.IsAny(".", "then") {
					CantUse(i)
					return false
				}
				isOf = false
			} else if wrd.Types.Has(WordBuzz) {
				lw.Set(wrd)
				continue
			} else if vrb == "tell" && wrd.Types.Has(WordVerb) && G.Winner == G.Player {
				Print("Please consult your manual for the correct way to talk to other people or creatures.", Newline)
				return false
			} else {
				CantUse(i)
				return false
			}
			lw.Set(wrd)
		}
		if G.Params.BufLen < 1 {
			G.Params.InQuotes = false
		}
	}
	G.Oops.UnkSet = false
	if hasDir {
		G.ActVerb.Norm = "walk"
		G.ActVerb.Orig = "walk"
		G.DirObj = nil
		G.Params.ShldOrphan = false
		G.Params.WalkDir = dir
		G.Params.HasWalkDir = true
		G.Again.Dir = dir
		G.Again.HasDir = true
	} else {
		if G.Params.ShldOrphan {
			OrphanMerge()
		}
		G.Params.HasWalkDir = false
		G.Again.HasDir = false
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
		Norm:  wrd,
		Orig:  wrd,
		Types: WordTypes{WordBuzz},
	}
}

func Clause(idx int, wrd LexItm) (bool, int) {
	if wrd.Types.Has(WordPrep) {
		if G.Params.ObjOrClauseCnt == 1 {
			G.ParsedSyntx.Prep1.Set(wrd)
		} else if G.Params.ObjOrClauseCnt == 2 {
			G.ParsedSyntx.Prep2.Set(wrd)
		}
		idx++
	} else {
		G.Params.BufLen++
	}
	if G.Params.BufLen == 0 {
		G.Params.ObjOrClauseCnt--
		return true, -1
	}
	cpyStart := idx
	if G.LexRes[idx].IsAny("the", "a", "an") {
		cpyStart++
	}
	var lw, nw LexItm
	isFirst := true
	var isAnd bool
	var i int
	G.Params.BufLen--
	for i = idx; G.Params.BufLen > -1; i, G.Params.BufLen = i+1, G.Params.BufLen-1 {
		HandleNumber(i)
		cw := G.LexRes[i]
		if cw.Types == nil && !cw.Is("intnum") {
			UnknownWord(i)
			return false, -1
		}
		if G.Params.BufLen == 0 {
			nw.Clear()
		} else {
			nw.Set(G.LexRes[i+1])
		}
		if cw.IsAny("and", ",") {
			isAnd = true
		} else if cw.IsAny("and", ",") {
			if nw.Is("of") {
				G.Params.BufLen--
				i++
			}
		} else if cw.IsAny("then", ".") || (cw.Types.Has(WordPrep) && G.ParsedSyntx.Verb.IsSet() && !isFirst) {
			G.Params.BufLen++
			if G.Params.ObjOrClauseCnt == 1 {
				G.ParsedSyntx.ObjOrClause1 = append([]LexItm{}, G.LexRes[cpyStart:i]...)
			} else if G.Params.ObjOrClauseCnt == 2 {
				G.ParsedSyntx.ObjOrClause2 = append([]LexItm{}, G.LexRes[cpyStart:i]...)
			}
			return true, i - 1
		} else if cw.Types.Has(WordObj) {
			if G.Params.BufLen > 0 && nw.Is("of") && !cw.IsAny("all", "one") {
				lw.Set(cw)
				isFirst = false
				continue
			}
			if cw.Types.Has(WordAdj) && nw.IsSet() && nw.Types.Has(WordObj) {
				lw.Set(cw)
				isFirst = false
				continue
			}
			if !isAnd && !nw.IsAny("but", "except", "and", ",") {
				if G.Params.ObjOrClauseCnt == 1 {
					G.ParsedSyntx.ObjOrClause1 = append([]LexItm{}, G.LexRes[cpyStart:i+1]...)
				} else if G.Params.ObjOrClauseCnt == 2 {
					G.ParsedSyntx.ObjOrClause2 = append([]LexItm{}, G.LexRes[cpyStart:i+1]...)
				}
				return true, i
			}
			isAnd = true
		} else if (G.Params.HasMerged || G.Params.ShldOrphan || G.ParsedSyntx.Verb.IsSet()) &&
			(cw.Types.Has(WordAdj) || cw.Types.Has(WordBuzz)) {
			lw.Set(cw)
			isFirst = false
			continue
		} else if isAnd && (cw.Types.Has(WordDir) || cw.Types.Has(WordVerb)) {
			i -= 2
			G.Params.BufLen += 2
			G.LexRes[i].Set(mkBuzz("then"))
		} else if cw.Types.Has(WordPrep) {
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
	if G.Params.ObjOrClauseCnt == 1 {
		G.ParsedSyntx.ObjOrClause1 = append([]LexItm{}, G.LexRes[cpyStart:i]...)
	} else if G.Params.ObjOrClauseCnt == 2 {
		G.ParsedSyntx.ObjOrClause2 = append([]LexItm{}, G.LexRes[cpyStart:i]...)
	}
	return true, -1
}

func UnknownWord(idx int) {
	G.Oops.UnkSet = true
	G.Oops.Unk = idx
	if G.ActVerb.Norm == "say" {
		Print("Nothing happens.", Newline)
		return
	}
	Print("I don't know the word \"", NoNewline)
	Print(G.LexRes[idx].Orig, NoNewline)
	Print("\".", Newline)
	G.Params.InQuotes = false
	G.Params.ShldOrphan = false
}

func CantUse(idx int) {
	if G.ActVerb.Norm == "say" {
		Print("Nothing happens.", Newline)
		return
	}
	Print("You used the word \"", NoNewline)
	Print(G.LexRes[idx].Orig, NoNewline)
	Print("\" in a way that I don't understand.", Newline)
	G.Params.InQuotes = false
	G.Params.ShldOrphan = false
}

// HandleNumber converts the lex item pointed
// to by idx into a number if possible.
func HandleNumber(idx int) {
	wrd := G.LexRes[idx]
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
	G.LexRes[idx].Norm = "intnum"
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
	G.Params.Number = sum
}

func OrphanMerge() {
	G.Params.ShldOrphan = false
	isAdj := false
	if G.ParsedSyntx.Verb.Types.Equals(G.OrphanedSyntx.Verb.Types) || G.ParsedSyntx.Verb.Types.Has(WordAdj) {
		isAdj = true
	} else if G.ParsedSyntx.Verb.Types.Has(WordObj) && G.Params.ObjOrClauseCnt == 0 {
		G.ParsedSyntx.Verb.Clear()
		G.ParsedSyntx.ObjOrClause1 = []LexItm{G.LexRes[0], G.LexRes[1]}
		G.Params.ObjOrClauseCnt = 1
	}
	if G.ParsedSyntx.Verb.IsSet() && !isAdj && !G.ParsedSyntx.Verb.Matches(G.OrphanedSyntx.Verb) {
		return
	}
	if G.Params.ObjOrClauseCnt == 2 {
		return
	}
	if G.OrphanedSyntx.ObjOrClause1 != nil && len(G.OrphanedSyntx.ObjOrClause1) == 0 {
		if !G.ParsedSyntx.Prep1.Matches(G.OrphanedSyntx.Prep1) && G.ParsedSyntx.Prep1.IsSet() {
			return
		}
		if isAdj {
			if G.Params.ObjOrClauseCnt == 0 {
				G.Params.ObjOrClauseCnt = 1
			}
			if len(G.ParsedSyntx.ObjOrClause1) == 0 {
				G.OrphanedSyntx.ObjOrClause1 = []LexItm{G.LexRes[0], G.LexRes[1]}
			} else {
				G.OrphanedSyntx.ObjOrClause1 = G.LexRes[0:G.ParsedSyntx.Obj1End]
			}
		} else {
			G.OrphanedSyntx.ObjOrClause1 = append([]LexItm{}, G.ParsedSyntx.ObjOrClause1...)
		}
	} else if G.OrphanedSyntx.ObjOrClause2 != nil && len(G.OrphanedSyntx.ObjOrClause2) == 0 {
		if !G.ParsedSyntx.Prep1.Matches(G.OrphanedSyntx.Prep2) && G.ParsedSyntx.Prep1.IsSet() {
			return
		}
		if isAdj {
			if len(G.ParsedSyntx.ObjOrClause1) == 0 {
				G.OrphanedSyntx.ObjOrClause2 = []LexItm{G.LexRes[0], G.LexRes[1]}
			} else {
				G.OrphanedSyntx.ObjOrClause2 = G.LexRes[0:G.ParsedSyntx.Obj1End]
			}
		} else {
			G.OrphanedSyntx.ObjOrClause2 = append([]LexItm{}, G.ParsedSyntx.ObjOrClause1...)
		}
		G.Params.ObjOrClauseCnt = 2
	} else if G.Params.AdjClause.Type != ClauseUnk {
		if G.Params.ObjOrClauseCnt != 1 && !isAdj {
			G.Params.AdjClause.Type = ClauseUnk
			return
		}
		beg := G.ParsedSyntx.Obj1Start
		if isAdj {
			beg = 0
			isAdj = false
		}
		var adj LexItm
		if G.ParsedSyntx.Obj1End == 0 {
			G.ParsedSyntx.Obj1End = 1
		}
		broken := false
		for i := beg; i < G.ParsedSyntx.Obj1End; i++ {
			wrd := G.LexRes[i]
			if !isAdj && (wrd.Types.Has(WordAdj) || wrd.IsAny("all", "one")) {
				adj.Set(wrd)
			} else if wrd.Is("one") {
				AclauseWin(adj)
				broken = true
				break
			} else if wrd.Types.Has(WordObj) {
				if wrd.Matches(G.Params.AdjClause.Syn) {
					AclauseWin(adj)
				} else {
					NclauseWin()
				}
				broken = true
				break
			}
		}
		if !broken {
			if G.ParsedSyntx.Obj1End == 1 {
				G.ParsedSyntx.ObjOrClause1 = []LexItm{G.LexRes[0]}
				G.Params.ObjOrClauseCnt = 1
			}
			if !adj.IsSet() {
				G.Params.AdjClause.Type = ClauseUnk
				return
			}
			AclauseWin(adj)
		}
	}
	G.ParsedSyntx.Set(G.OrphanedSyntx)
	G.Params.HasMerged = true
}

func AclauseWin(adj LexItm) {
	G.ParsedSyntx.Verb.Set(G.OrphanedSyntx.Verb)
	var tbl *[]LexItm = &[]LexItm{}
	if G.Params.AdjClause.Type == Clause1 {
		tbl = &G.OrphanedSyntx.ObjOrClause1
	} else if G.Params.AdjClause.Type == Clause2 {
		tbl = &G.OrphanedSyntx.ObjOrClause2
	}
	for idx, obj := range *tbl {
		if obj.Matches(G.Params.AdjClause.Adj) {
			*tbl = append((*tbl)[0:idx], append([]LexItm{adj}, (*tbl)[idx:len(*tbl)]...)...)
			break
		}
	}
	if G.OrphanedSyntx.ObjOrClause2 != nil {
		G.Params.ObjOrClauseCnt = 2
	}
	G.Params.AdjClause.Type = ClauseUnk
}

func NclauseWin() {
	if G.Params.AdjClause.Type == Clause1 {
		G.OrphanedSyntx.ObjOrClause1 = append([]LexItm{}, G.ParsedSyntx.ObjOrClause1...)
		G.OrphanedSyntx.Obj1Start = G.ParsedSyntx.Obj1Start
		G.OrphanedSyntx.Obj1End = G.ParsedSyntx.Obj1End
	} else if G.Params.AdjClause.Type == Clause2 {
		G.OrphanedSyntx.ObjOrClause2 = append([]LexItm{}, G.ParsedSyntx.ObjOrClause1...)
	}
	if G.OrphanedSyntx.ObjOrClause2 != nil {
		G.Params.ObjOrClauseCnt = 2
	}
	G.Params.AdjClause.Type = ClauseUnk
}

func TakeCheck() bool {
	if !ITakeCheck(G.DirObjPossibles, G.DetectedSyntx.Obj1.LocFlags) {
		return false
	}
	return ITakeCheck(G.IndirObjPossibles, G.DetectedSyntx.Obj2.LocFlags)
}

func ITakeCheck(tbl []*Object, ibits LocFlags) bool {
	if len(tbl) == 0 || (!LocHave.In(ibits) && !LocTake.In(ibits)) {
		return true
	}
	for _, obj := range tbl {
		if obj == &It {
			if !IsAccessible(G.Params.ItObj) {
				Print("I don't see what you're referring to.", Newline)
				return false
			}
			obj = G.Params.ItObj
		}
		var taken bool
		if !IsHeld(obj) && obj != &Hands && obj != &Me {
			G.DirObjPossibles = []*Object{obj}
			if obj.Has(FlgTryTake) {
				taken = true
			} else if G.Winner != &Adventurer {
				taken = false
			} else if LocTake.In(ibits) && ITake(false) {
				taken = false
			} else {
				taken = true
			}
			if taken && LocHave.In(ibits) && G.Winner == &Adventurer {
				if obj == &NotHereObject {
					Print("You don't have that!", Newline)
					return false
				}
				Print("You don't have the ", NoNewline)
				PrintObject(obj)
				Print(".", Newline)
				return false
			}
			if !taken && G.Winner == &Adventurer {
				Print("(Taken)", Newline)
			}
		}
	}
	return true
}

func ManyCheck() bool {
	loss := 0
	if len(G.DirObjPossibles) > 1 && !LocMany.In(G.DetectedSyntx.Obj1.LocFlags) {
		loss = 1
	} else if len(G.IndirObjPossibles) > 1 && !LocMany.In(G.DetectedSyntx.Obj2.LocFlags) {
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
	if !G.ParsedSyntx.Verb.IsSet() {
		Print("tell", NoNewline)
	} else if G.Params.ShldOrphan || G.Params.HasMerged {
		Print(G.ParsedSyntx.Verb.Norm, NoNewline)
	} else {
		Print(G.ParsedSyntx.Verb.Orig, NoNewline)
	}
	Print("\".", Newline)
	return false
}

func SnarfObjects() bool {
	G.Params.Buts = []*Object{}
	if G.ParsedSyntx.ObjOrClause2 != nil {
		G.Search.LocFlags = G.DetectedSyntx.Obj2.LocFlags
		res := Snarfem(false, G.ParsedSyntx.ObjOrClause2)
		if res == nil {
			return false
		}
		G.IndirObjPossibles = append(G.IndirObjPossibles, res...)
	}
	if G.ParsedSyntx.ObjOrClause1 != nil {
		G.Search.LocFlags = G.DetectedSyntx.Obj1.LocFlags
		res := Snarfem(true, G.ParsedSyntx.ObjOrClause1)
		if res == nil {
			return false
		}
		G.DirObjPossibles = append(G.DirObjPossibles, res...)
	}
	if len(G.Params.Buts) != 0 {
		l := len(G.DirObjPossibles)
		if len(G.ParsedSyntx.ObjOrClause1) != 0 {
			G.DirObjPossibles = ButMerge(G.DirObjPossibles)
		}
		if len(G.ParsedSyntx.ObjOrClause2) != 0 && (len(G.ParsedSyntx.ObjOrClause1) == 0 || l == len(G.DirObjPossibles)) {
			G.IndirObjPossibles = ButMerge(G.IndirObjPossibles)
		}
	}
	return true
}

func ButMerge(tbl []*Object) []*Object {
	res := []*Object{}
	for _, bts := range G.Params.Buts {
		for _, obj := range tbl {
			if obj == bts {
				res = append(res, obj)
			}
		}
	}
	return res
}

func Snarfem(isDirect bool, wrds []LexItm) []*Object {
	G.Params.HasAnd = false
	wasall := false
	if G.Params.GetType == GetAll {
		wasall = true
	}
	G.Search.ObjFlags = 0
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
			G.Params.GetType = GetAll
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
			but = &G.Params.Buts
			*but = []*Object{}
		} else if wrd.IsAny("a", "one") {
			if !G.Search.Adj.IsSet() {
				G.Params.GetType = GetOne
				if nw.Is("of") {
					continue
				}
			} else {
				G.Search.Syn.Set(G.Params.OneObj)
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
			G.Params.HasAnd = true
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
		} else if wrd.Types.Has(WordBuzz) {
			continue
		} else if wrd.IsAny("and", ",") {
			continue
		} else if wrd.Is("of") {
			if G.Params.GetType == GetUndef {
				G.Params.GetType = GetInhibit
			}
		} else if wrd.Types.Has(WordAdj) && !G.Search.Adj.IsSet() {
			G.Search.Adj.Set(wrd)
		} else if wrd.Types.Has(WordObj) {
			G.Search.Syn.Set(wrd)
			G.Params.OneObj.Set(wrd)
		}
	}
	out := GetObject(isDirect, true)
	if wasall {
		G.Params.GetType = GetAll
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
	if !G.ParsedSyntx.Verb.IsSet() {
		Print("There was no verb in that sentence!", Newline)
		return false
	}
	var findFirst, findSecond *Syntx
	for _, syn := range Commands {
		if !G.ParsedSyntx.Verb.Is(syn.Verb) {
			continue
		}
		if G.Params.ObjOrClauseCnt > syn.NumObjects() {
			continue
		}
		if syn.NumObjects() >= 1 && G.Params.ObjOrClauseCnt == 0 && (!G.ParsedSyntx.Prep1.IsSet() || syn.IsVrbPrep(G.ParsedSyntx.Prep1.Norm)) {
			findFirst = &syn
		} else if syn.IsVrbPrep(G.ParsedSyntx.Prep1.Norm) {
			if syn.NumObjects() == 2 && G.Params.ObjOrClauseCnt == 1 {
				findSecond = &syn
			} else if syn.IsObjPrep(G.ParsedSyntx.Prep2.Norm) {
				G.DetectedSyntx = &syn
				G.ActVerb.Norm = syn.GetNormVerb()
				G.ActVerb.Orig = syn.GetActionVerb()
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
			G.DirObjPossibles = []*Object{obj}
			G.DetectedSyntx = findFirst
			G.ActVerb.Norm = findFirst.GetNormVerb()
			G.ActVerb.Orig = findFirst.GetActionVerb()
			found = true
		}
	}
	if findSecond != nil && !found {
		obj := FindWhatIMean(findSecond.Obj2.ObjFlags, findSecond.Obj2.LocFlags, findSecond.ObjPrep)
		if obj != nil {
			G.IndirObjPossibles = []*Object{obj}
			G.DetectedSyntx = findSecond
			G.ActVerb.Norm = findSecond.GetNormVerb()
			G.ActVerb.Orig = findSecond.GetActionVerb()
			found = true
		}
	}
	if G.ParsedSyntx.Verb.Is("find") && !found {
		Print("That question can't be answered.", Newline)
		return false
	}
	if G.Winner != G.Player && !found {
		CanNotOrphan()
		return false
	}
	if !found {
		Orphan(findFirst, findSecond)
		Print("What do you want to ", NoNewline)
		if !G.OrphanedSyntx.Verb.IsSet() {
			Print("tell", NoNewline)
		} else {
			Print(G.OrphanedSyntx.Verb.Orig, NoNewline)
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
	G.OrphanedSyntx.Set(G.ParsedSyntx)
	if G.Params.ObjOrClauseCnt < 2 {
		G.OrphanedSyntx.ObjOrClause2 = []LexItm{}
	}
	if G.Params.ObjOrClauseCnt < 1 {
		G.OrphanedSyntx.ObjOrClause1 = []LexItm{}
	}
	if first != nil {
		G.OrphanedSyntx.Prep1.Norm = first.VrbPrep
		G.OrphanedSyntx.Prep1.Orig = first.VrbPrep
		G.OrphanedSyntx.Prep1.Types = WordTypes{WordPrep}
		G.OrphanedSyntx.ObjOrClause1 = []LexItm{}
	} else if second != nil {
		G.OrphanedSyntx.Prep2.Norm = second.ObjPrep
		G.OrphanedSyntx.Prep2.Orig = second.ObjPrep
		G.OrphanedSyntx.Prep2.Types = WordTypes{WordPrep}
		G.OrphanedSyntx.ObjOrClause2 = []LexItm{}
	}
}

func CanNotOrphan() {
	Print("\"I don't understand! What are you referring to?\"", Newline)
}

func FindWhatIMean(objFlags Flags, locFlags LocFlags, prep string) *Object {
	if objFlags&FlgKludge != 0 {
		return &Rooms
	}
	G.Search.ObjFlags = objFlags
	G.Search.LocFlags = locFlags
	res := GetObject(false, false)
	G.Search.ObjFlags = 0
	if len(res) != 1 {
		return nil
	}
	Print("(", NoNewline)
	if len(prep) == 0 || G.Params.EndOnPrep {
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
	if G.Params.GetType == GetInhibit {
		return []*Object{}
	}
	if !G.Search.Syn.IsSet() && G.Search.Adj.IsSet() && G.Search.Adj.Types.Has(WordObj) {
		G.Search.Syn.Set(G.Search.Adj)
		G.Search.Adj.Clear()
	}
	if !G.Search.Syn.IsSet() && !G.Search.Adj.IsSet() && G.Params.GetType != GetAll && G.Search.ObjFlags == 0 {
		if vrb {
			Print("There seems to be a noun missing in that sentence!", Newline)
		}
		return nil
	}
	xbits := G.Search.LocFlags
	if G.Params.GetType != GetAll || len(G.Search.LocFlags) == 0 {
		G.Search.LocFlags.All()
	}
	res := []*Object{}
	gcheck := false
	olen := 0
	for {
		if gcheck {
			res = append(res, GlobalCheck()...)
		} else {
			if G.Lit {
				G.Player.Take(FlgTrans)
				res = append(res, DoSL(G.Here, LocOnGrnd, LocInRoom)...)
				G.Player.Give(FlgTrans)
			}
			res = append(res, DoSL(G.Player, LocHeld, LocCarried)...)
		}
		// Deduplicate results
		res = dedup(res)
		ln := len(res)
		if G.Params.GetType == GetAll {
			G.Search.LocFlags = xbits
			G.Search.Syn.Clear()
			G.Search.Adj.Clear()
			return res
		}
		if G.Params.GetType == GetOne && ln != 0 {
			if ln > 1 {
				res = []*Object{res[G.Rand.Intn(len(res))]}
				Print("(How about the ", NoNewline)
				PrintObject(res[0])
				Print("?)", Newline)
			}
		} else if ln > 1 || (ln == 0 && G.Search.LocFlags.HasAll()) {
			if G.Search.LocFlags.HasAll() {
				G.Search.LocFlags = xbits
				olen = ln
				res = []*Object{}
				continue
			}
			if ln == 0 {
				ln = olen
			}
			if G.Winner != G.Player {
				CanNotOrphan()
				return nil
			}
			if vrb && G.Search.Syn.IsSet() {
				WhichPrint(isDirect, res)
				if isDirect {
					G.Params.AdjClause.Type = Clause1
				} else {
					G.Params.AdjClause.Type = Clause2
				}
				G.Params.AdjClause.Syn.Set(G.Search.Syn)
				G.Params.AdjClause.Adj.Set(G.Search.Adj)
				Orphan(nil, nil)
				G.Params.ShldOrphan = true
			} else if vrb {
				Print("There seems to be a noun missing in that sentence!", Newline)
			}
			G.Search.Syn.Clear()
			G.Search.Adj.Clear()
			return nil
		}
		if ln == 0 && gcheck {
			if vrb {
				G.Search.LocFlags = xbits
				if G.Lit || G.ActVerb.Norm == "tell" {
					res = append(res, &NotHereObject)
					G.NotHere.Syn.Set(G.Search.Syn)
					G.NotHere.Adj.Set(G.Search.Adj)
					G.Search.Syn.Clear()
					G.Search.Adj.Clear()
					return res
				}
				Print("It's too dark to see!", Newline)
			}
			G.Search.Syn.Clear()
			G.Search.Adj.Clear()
			return nil
		}
		if ln == 0 {
			gcheck = true
			continue
		}
		G.Search.LocFlags = xbits
		G.Search.Syn.Clear()
		G.Search.Adj.Clear()
		return res
	}
}

// WhichPrint outputs all of the possible matches
// when the game parser matches multiple game objects.
func WhichPrint(isDirect bool, tbl []*Object) {
	Print("Which ", NoNewline)
	if G.Params.ShldOrphan || G.Params.HasMerged || G.Params.HasAnd {
		if G.Search.Syn.IsSet() {
			Print(G.Search.Syn.Norm, NoNewline)
		} else if G.Search.Adj.IsSet() {
			Print(G.Search.Adj.Norm, NoNewline)
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
	if len(G.Here.Global) != 0 {
		for _, obj := range G.Here.Global {
			if IsThisIt(obj) {
				res = append(res, obj)
			}
		}
	}
	if len(G.Here.Pseudo) != 0 {
		for _, obj := range G.Here.Pseudo {
			if G.Search.Syn.Is(obj.Synonym) {
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
		if len(res) == 0 && (G.ActVerb.Norm == "look inside" || G.ActVerb.Norm == "search" || G.ActVerb.Norm == "examine") {
			if LocHave.In(G.Search.LocFlags) {
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
	search := &G.ParsedSyntx.ObjOrClause1
	if !isDirect {
		search = &G.ParsedSyntx.ObjOrClause2
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
			PrintNumber(G.Params.Number)
			pn = true
		} else {
			if isFirst && !pn && isThe {
				Print("the ", NoNewline)
			}
			if G.Params.ShldOrphan || G.Params.HasMerged {
				Print(wrd.Norm, NoNewline)
			} else if wrd.Is("it") && IsAccessible(G.Params.ItObj) {
				PrintObject(G.Params.ItObj)
			} else {
				Print(wrd.Orig, NoNewline)
			}
			isFirst = false
		}
	}
}

// IsLit checks if the current game room is lit.
func IsLit(room *Object, rmChk bool) bool {
	if G.AlwaysLit && G.Winner == G.Player {
		return true
	}
	G.Search.ObjFlags = FlgOn
	bak := G.Here
	G.Here = room
	if rmChk && room.Has(FlgOn) {
		G.Here = bak
		G.Search.ObjFlags = 0
		return true
	}
	G.Search.LocFlags = nil
	res := []*Object{}
	if bak == room {
		if nr := SearchList(G.Winner, FindAll); nr != nil {
			res = append(res, nr...)
		}
		if G.Winner != G.Player && G.Player.IsIn(room) {
			if nr := SearchList(G.Player, FindAll); nr != nil {
				res = append(res, nr...)
			}
		}
	}
	if nr := SearchList(room, FindAll); nr != nil {
		res = append(res, nr...)
	}
	G.Here = bak
	G.Search.ObjFlags = 0
	if len(res) > 0 {
		return true
	}
	return false
}

// DoSL performs a specific game object search
// based on the provided location flags and
// parameters defined in the Search global variable.
func DoSL(obj *Object, f1, f2 LocFlag) []*Object {
	if f1.In(G.Search.LocFlags) && f2.In(G.Search.LocFlags) {
		return SearchList(obj, FindAll)
	} else if f1.In(G.Search.LocFlags) {
		return SearchList(obj, FindTop)
	} else if f2.In(G.Search.LocFlags) {
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
	if G.Search.Syn.IsSet() && !obj.Is(G.Search.Syn.Norm) {
		return false
	}
	if G.Search.Adj.IsSet() && (len(obj.Adjectives) == 0 || !obj.Is(G.Search.Adj.Norm)) {
		return false
	}
	if !AnyFlagIn(G.Search.ObjFlags, obj.Flags) {
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
	if l == &LocalGlobals && IsInGlobal(obj, G.Here) {
		return true
	}
	ml := MetaLoc(obj)
	if ml != G.Here && G.Winner != nil && ml != G.Winner.Location() {
		return false
	}
	if l == G.Winner || l == G.Here || (G.Winner != nil && l == G.Winner.Location()) {
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

// dedup removes duplicate *Object entries from a slice.
func dedup(objs []*Object) []*Object {
	seen := make(map[*Object]bool, len(objs))
	result := make([]*Object, 0, len(objs))
	for _, o := range objs {
		if !seen[o] {
			seen[o] = true
			result = append(result, o)
		}
	}
	return result
}
