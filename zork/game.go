package zork

import "os"

var (
	ParserOk bool
	Player   *Object
)

type PerfRet int

const (
	PerfNotHndld PerfRet = iota
	PerfHndld
	PerfFatal
)

func Restart() {
	// TODO: Implement restart
}

func Restore() bool {
	// TODO: Implement restore
	return false
}

func Quit() {
	// TODO: Implement quit
	os.Exit(0)
}

func Run() {
	FinalizeGameObjects()
	BuildObjectTree()
	BuildVocabulary()
	for {
		Queue(IFight, -1).Run = true
		Queue(ISword, -1)
		Queue(IThief, -1).Run = true
		Queue(ICandles, 40)
		Queue(ILantern, 200)
		InflatedBoat.VehType = FlgNonLand
		Def1Res[1] = Def1[2]
		Def1Res[2] = Def1[4]
		Def2Res[2] = Def2B[2]
		Def2Res[3] = Def2B[4]
		Def3Res[1] = Def3A[2]
		Def3Res[3] = Def3B[2]
		Here = &WestOfHouse
		ThisIsIt(&Mailbox)
		if !Here.Has(FlgTouch) {
			VVersion(ActUnk)
			NewLine()
		}
		Lit = true
		Winner = &Adventurer
		Player = Winner
		Winner.MoveTo(Here)
		VLook(ActUnk)
		MainLoop()
	}
}

func MainLoop() {
	Params.Continue = NumUndef
	for {
		ParserOk = Parse()
		if !ParserOk {
			Params.Continue = NumUndef
			continue
		}
		if Params.ItObj != nil && IsAccessible(Params.ItObj) {
			found := false
			for idx, obj := range IndirObjPossibles {
				if obj == &It {
					IndirObjPossibles[idx] = Params.ItObj
					found = true
				}
			}
			if !found {
				for idx, obj := range DirObjPossibles {
					if obj == &It {
						DirObjPossibles[idx] = Params.ItObj
					}
				}
			}
		}
		numObj := 1
		var obj *Object
		isDir := true
		if len(DirObjPossibles) == 0 {
			numObj = 0
		} else if len(DirObjPossibles) > 1 {
			if len(IndirObjPossibles) == 0 {
				obj = nil
			} else {
				obj = IndirObjPossibles[0]
			}
			numObj = len(DirObjPossibles)
		} else if len(IndirObjPossibles) > 1 {
			isDir = false
			obj = DirObjPossibles[0]
			numObj = len(IndirObjPossibles)
		}
		if obj == nil && len(IndirObjPossibles) == 1 {
			obj = IndirObjPossibles[0]
		}
		var res PerfRet
		if ActVerb == "walk" && len(Params.WalkDir) != 0 {
			res = Perform(ActVerb, DirObj, nil)
		} else if numObj == 0 {
			if DetectedSyntx.NumObjects() == 0 {
				res = Perform(ActVerb, nil, nil)
				DirObj = nil
			} else if !Lit {
				Print("It's not clear what you're referring to.", Newline)
				res = PerfNotHndld
			}
		} else {
			tmp := false
			notHere := 0
			var obj1 *Object
			var indir, dir *Object
			for i := 0; i < numObj; i++ {
				if isDir {
					obj1 = DirObjPossibles[i]
					dir, indir = obj1, obj
				} else {
					obj1 = IndirObjPossibles[i]
					dir, indir = obj, obj1
				}
				if numObj > 1 || (len(ParsedSyntx.ObjOrClause1) > 0 && ParsedSyntx.ObjOrClause1[0].Is("all")) {
					if dir == &NotHereObject {
						notHere++
						continue
					}
					if ActVerb == "take" &&
						indir != nil &&
						len(ParsedSyntx.ObjOrClause1) > 0 &&
						ParsedSyntx.ObjOrClause1[0].Is("all") &&
						dir != nil &&
						!dir.IsIn(indir) {
						continue
					}
					if l := dir.Location(); Params.GetType == GetAll &&
						ActVerb == "take" &&
						((l != Winner &&
							l != Here &&
							l != Winner.Location() &&
							l != indir &&
							!l.Has(FlgSurf)) ||
							!(dir.Has(FlgTake) || dir.Has(FlgTryTake))) {
						continue
					}
					if obj1 == &It {
						PrintObject(Params.ItObj)
					} else {
						PrintObject(obj1)
					}
					Print(": ", NoNewline)
				}
				DirObj, IndirObj = dir, indir
				tmp = true
				res = Perform(ActVerb, DirObj, IndirObj)
				if res == PerfFatal {
					break
				}
			}
			if notHere > 0 {
				Print("The ", NoNewline)
				if notHere != numObj {
					Print("other ", NoNewline)
				}
				Print("object", NoNewline)
				if notHere != 1 {
					Print("s", NoNewline)
				}
				Print(" that you mentioned ", NoNewline)
				if notHere != 1 {
					Print("are", NoNewline)
				} else {
					Print("is", NoNewline)
				}
				Print("n't here.", Newline)
			} else if !tmp {
				Print("There's nothing here you can take.", Newline)
			}
		}
		if l := Winner.Location(); res != PerfFatal && l != nil && l.Action != nil {
			if l.Action(ActEnd) {
				res = PerfHndld
			} else {
				res = PerfNotHndld
			}
		}
		if res == PerfFatal {
			Params.Continue = -1
		}
		if ActVerb == "tell" || ActVerb == "brief" || ActVerb == "superbrief" || ActVerb == "verbose" || ActVerb == "save" || ActVerb == "version" || ActVerb == "quit" || ActVerb == "restart" || ActVerb == "score" || ActVerb == "script" || ActVerb == "unscript" || ActVerb == "restore" {
			continue
		} else {
			Clocker()
		}
	}
}

func Perform(a string, o, i *Object) PerfRet {
	if (o == &It || i == &It) && IsAccessible(Params.ItObj) {
		Print("I don't see what you are referring to.", Newline)
		return PerfFatal
	}
	if o == &It {
		o = Params.ItObj
	}
	if i == &It {
		i = Params.ItObj
	}
	if o != nil && IndirObj != &It && a == "walk" {
		Params.ItObj = o
	}
	if o == &NotHereObject || i == &NotHereObject {
		if NotHereObjectFcn(ActUnk) {
			return PerfHndld
		}
	}
	if Winner != nil && Winner.Action != nil {
		if Winner.Action(ActUnk) {
			return PerfHndld
		}
	}
	if l := Winner.Location(); l != nil && l.Action != nil {
		if l.Action(ActBegin) {
			return PerfHndld
		}
	}
	if act, ok := PreActions[a]; ok && act != nil {
		if act(ActUnk) {
			return PerfHndld
		}
	}
	if i != nil && i.Action != nil {
		if i.Action(ActUnk) {
			return PerfHndld
		}
	}
	if o != nil && a != "walk" && o.Location() != nil && o.Location().ContFcn != nil {
		if o.Location().ContFcn(ActUnk) {
			return PerfHndld
		}
	}
	if o != nil && a != "walk" && o.Action != nil {
		if o.Action(ActUnk) {
			return PerfHndld
		}
	}
	if act, ok := Actions[a]; ok && act != nil {
		if act(ActUnk) {
			return PerfHndld
		}
	}
	return PerfNotHndld
}
