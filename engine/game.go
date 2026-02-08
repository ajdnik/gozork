package engine

type PerfRet int

const (
	PerfNotHndld PerfRet = iota
	PerfHndld
	PerfFatal
	PerfQuit
)

// RFatal signals a fatal action result (equivalent to ZIL's <RFATAL>).
func RFatal() bool {
	G.PerformFatal = true
	return true
}

// Quit signals that the game should exit.
func Quit() {
	G.QuitRequested = true
}

// Verify checks the integrity of the game file.
func Verify() bool {
	return true
}

// ResetGameState creates a fresh GameState with all defaults and restores the
// object tree. Tests and restart use this to get a clean game.
func ResetGameState() {
	out, in_, rng := G.GameOutput, G.GameInput, G.Rand
	clockFuncs := G.ClockFuncs
	allObjs := G.AllObjects
	roomsObj := G.RoomsObj
	globalObj := G.GlobalObj
	localGlobalObj := G.LocalGlobalObj
	notHereObj := G.NotHereObj
	pseudoObj := G.PseudoObj
	itPronounObj := G.ItPronounObj
	meObj := G.MeObj
	handsObj := G.HandsObj
	gameData := G.GameData

	G = NewGameState()
	if out != nil {
		G.GameOutput = out
	}
	if in_ != nil {
		G.GameInput = in_
	}
	if rng != nil {
		G.Rand = rng
	}
	G.ClockFuncs = clockFuncs
	G.AllObjects = allObjs
	G.RoomsObj = roomsObj
	G.GlobalObj = globalObj
	G.LocalGlobalObj = localGlobalObj
	G.NotHereObj = notHereObj
	G.PseudoObj = pseudoObj
	G.ItPronounObj = itPronounObj
	G.MeObj = meObj
	G.HandsObj = handsObj
	G.GameData = gameData

	ResetObjectTree()
}

func MainLoop() {
	G.Params.Continue = NumUndef
	for {
		if G.InputExhausted || G.QuitRequested {
			return
		}
		G.ParserOk = Parse()
		if G.InputExhausted || G.QuitRequested {
			return
		}
		if !G.ParserOk {
			G.Params.Continue = NumUndef
			continue
		}
		if G.Params.ItObj != nil && IsAccessible(G.Params.ItObj) {
			found := false
			for idx, obj := range G.IndirObjPossibles {
				if obj == G.ItPronounObj {
					G.IndirObjPossibles[idx] = G.Params.ItObj
					found = true
				}
			}
			if !found {
				for idx, obj := range G.DirObjPossibles {
					if obj == G.ItPronounObj {
						G.DirObjPossibles[idx] = G.Params.ItObj
					}
				}
			}
		}
		numObj := 1
		var obj *Object
		isDir := true
		if len(G.DirObjPossibles) == 0 {
			numObj = 0
		} else if len(G.DirObjPossibles) > 1 {
			if len(G.IndirObjPossibles) == 0 {
				obj = nil
			} else {
				obj = G.IndirObjPossibles[0]
			}
			numObj = len(G.DirObjPossibles)
		} else if len(G.IndirObjPossibles) > 1 {
			isDir = false
			obj = G.DirObjPossibles[0]
			numObj = len(G.IndirObjPossibles)
		}
		if obj == nil && len(G.IndirObjPossibles) == 1 {
			obj = G.IndirObjPossibles[0]
		}
		var res PerfRet
		if G.ActVerb.Norm == "walk" && G.Params.HasWalkDir {
			res = Perform(G.ActVerb, G.DirObj, nil)
		} else if numObj == 0 {
			if G.DetectedSyntx.NumObjects() == 0 {
				res = Perform(G.ActVerb, nil, nil)
				G.DirObj = nil
			} else if !G.Lit {
				Printf("It's too dark to see.\n")
				res = PerfNotHndld
			} else {
				Printf("It's not clear what you're referring to.\n")
				res = PerfNotHndld
			}
		} else {
			tmp := false
			notHere := 0
			var obj1 *Object
			var indir, dir *Object
			for i := 0; i < numObj; i++ {
				if isDir {
					obj1 = G.DirObjPossibles[i]
					dir, indir = obj1, obj
				} else {
					obj1 = G.IndirObjPossibles[i]
					dir, indir = obj, obj1
				}
				if numObj > 1 || (len(G.ParsedSyntx.ObjOrClause1) > 0 && G.ParsedSyntx.ObjOrClause1[0].Is("all")) {
					if dir == G.NotHereObj {
						notHere++
						continue
					}
					if G.ActVerb.Norm == "take" &&
						indir != nil &&
						len(G.ParsedSyntx.ObjOrClause1) > 0 &&
						G.ParsedSyntx.ObjOrClause1[0].Is("all") &&
						dir != nil &&
						!dir.IsIn(indir) {
						continue
					}
					if l := dir.Location(); G.Params.GetType == GetAll &&
						G.ActVerb.Norm == "take" &&
						((l != G.Winner &&
							l != G.Here &&
							l != G.Winner.Location() &&
							l != indir &&
							!l.Has(FlgSurf)) ||
							!(dir.Has(FlgTake) || dir.Has(FlgTryTake))) {
						continue
					}
					if obj1 == G.ItPronounObj {
						Printf("%s", G.Params.ItObj.Desc)
					} else {
						Printf("%s", obj1.Desc)
					}
					Printf(": ")
				}
				G.DirObj, G.IndirObj = dir, indir
				tmp = true
				res = Perform(G.ActVerb, G.DirObj, G.IndirObj)
				if res == PerfFatal || res == PerfQuit {
					break
				}
			}
			if notHere > 0 {
				Printf("The ")
				if notHere != numObj {
					Printf("other ")
				}
				Printf("object")
				if notHere != 1 {
					Printf("s")
				}
				Printf(" that you mentioned ")
				if notHere != 1 {
					Printf("are")
				} else {
					Printf("is")
				}
				Printf("n't here.\n")
			} else if !tmp {
				Printf("There's nothing here you can take.\n")
			}
		}
		if G.QuitRequested {
			return
		}
		if l := G.Winner.Location(); res != PerfFatal && res != PerfQuit && l != nil && l.Action != nil {
			if l.Action(ActEnd) {
				res = PerfHndld
			} else {
				res = PerfNotHndld
			}
		}
		if G.QuitRequested {
			return
		}
		if res == PerfFatal {
			G.Params.Continue = -1
		}
		if G.ActVerb.Norm == "tell" || G.ActVerb.Norm == "brief" || G.ActVerb.Norm == "superbrief" || G.ActVerb.Norm == "verbose" || G.ActVerb.Norm == "save" || G.ActVerb.Norm == "version" || G.ActVerb.Norm == "quit" || G.ActVerb.Norm == "restart" || G.ActVerb.Norm == "score" || G.ActVerb.Norm == "script" || G.ActVerb.Norm == "unscript" || G.ActVerb.Norm == "restore" {
			continue
		} else {
			Clocker()
		}
	}
}

// callHandler invokes an action handler and checks for quit/fatal signals.
func callHandler(fn func(ActArg) bool, arg ActArg) (PerfRet, bool) {
	result := fn(arg)
	if G.QuitRequested {
		return PerfQuit, true
	}
	if G.PerformFatal {
		G.PerformFatal = false
		return PerfFatal, true
	}
	if result {
		return PerfHndld, true
	}
	return PerfNotHndld, false
}

func Perform(a ActionVerb, o, i *Object) PerfRet {
	oldActVerb := G.ActVerb
	oldDirObj := G.DirObj
	oldIndirObj := G.IndirObj
	defer func() {
		G.ActVerb = oldActVerb
		G.DirObj = oldDirObj
		G.IndirObj = oldIndirObj
	}()

	G.PerformFatal = false
	if (o == G.ItPronounObj || i == G.ItPronounObj) && !IsAccessible(G.Params.ItObj) {
		Printf("I don't see what you are referring to.\n")
		return PerfFatal
	}
	if o == G.ItPronounObj {
		o = G.Params.ItObj
	}
	if i == G.ItPronounObj {
		i = G.Params.ItObj
	}
	G.ActVerb = a
	G.DirObj = o
	if o != nil && a.Norm != "walk" {
		G.Params.ItObj = o
	}
	G.IndirObj = i

	if o == G.NotHereObj || i == G.NotHereObj {
		if G.NotHereObj.Action != nil {
			if ret, done := callHandler(G.NotHereObj.Action, ActUnk); done {
				return ret
			}
		}
	}
	if G.Winner != nil && G.Winner.Action != nil {
		if ret, done := callHandler(G.Winner.Action, ActUnk); done {
			return ret
		}
	}
	if l := G.Winner.Location(); l != nil && l.Action != nil {
		if ret, done := callHandler(l.Action, ActBegin); done {
			return ret
		}
	}
	if act, ok := G.PreActions[a.Orig]; ok && act != nil {
		if ret, done := callHandler(act, ActUnk); done {
			return ret
		}
	} else if a.Norm != a.Orig {
		if act, ok := G.PreActions[a.Norm]; ok && act != nil {
			if ret, done := callHandler(act, ActUnk); done {
				return ret
			}
		}
	}
	if i != nil && i.Action != nil {
		if ret, done := callHandler(i.Action, ActUnk); done {
			return ret
		}
	}
	if o != nil && a.Norm != "walk" && o.Location() != nil && o.Location().ContFcn != nil {
		if ret, done := callHandler(o.Location().ContFcn, ActUnk); done {
			return ret
		}
	}
	if o != nil && a.Norm != "walk" && o.Action != nil {
		if ret, done := callHandler(o.Action, ActUnk); done {
			return ret
		}
	}
	if act, ok := G.Actions[a.Orig]; ok && act != nil {
		if ret, done := callHandler(act, ActUnk); done {
			return ret
		}
	} else if a.Norm != a.Orig {
		if act, ok := G.Actions[a.Norm]; ok && act != nil {
			if ret, done := callHandler(act, ActUnk); done {
				return ret
			}
		}
	}
	return PerfNotHndld
}
