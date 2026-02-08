package zork


func PreBoard(arg ActArg) bool {
	if DirObj.Has(FlgVeh) {
		if !DirObj.IsIn(Here) {
			Print("The ", NoNewline)
			PrintObject(DirObj)
			Print(" must be on the ground to be boarded.", Newline)
			return RFatal()
		}
		if av := Winner.Location(); av != nil && av.Has(FlgVeh) {
			Print("You are already in the ", NoNewline)
			PrintObject(av)
			Print("!", Newline)
			return RFatal()
		}
		return false
	}
	if DirObj == &Water || DirObj == &GlobalWater {
		Perform(ActionVerb{Norm: "swim", Orig: "swim"}, DirObj, nil)
		return true
	}
	Print("You have a theory on how to board a ", NoNewline)
	PrintObject(DirObj)
	Print(", perhaps?", Newline)
	return RFatal()
}

func VBoard(arg ActArg) bool {
	Print("You are now in the ", NoNewline)
	PrintObject(DirObj)
	Print(".", Newline)
	Winner.MoveTo(DirObj)
	if DirObj.Action != nil {
		DirObj.Action(ActEnter)
	}
	return true
}

func VClimbDown(arg ActArg) bool {
	return VClimbFcn("down", DirObj)
}

func VClimbFoo(arg ActArg) bool {
	return VClimbFcn("up", DirObj)
}

func VClimbOn(arg ActArg) bool {
	if !DirObj.Has(FlgVeh) {
		Print("You can't climb onto the ", NoNewline)
		PrintObject(DirObj)
		Print(".", Newline)
		return true
	}
	Perform(ActionVerb{Norm: "board", Orig: "board"}, DirObj, nil)
	return true
}

func VClimbUp(arg ActArg) bool {
	return VClimbFcn("up", nil)
}

func VClimbFcn(dir string, obj *Object) bool {
	if obj != nil && DirObj != &Rooms {
		obj = DirObj
	}
	if tx := Here.GetDir(dir); tx != nil {
		if obj != nil {
			if len(tx.NExit) > 0 || ((tx.CExit != nil || tx.DExit != nil || tx.UExit) && !IsInGlobal(DirObj, tx.RExit)) {
				Print("The ", NoNewline)
				PrintObject(obj)
				Print(" do", NoNewline)
				if obj != &Stairs {
					Print("es", NoNewline)
				}
				Print("n't lead ", NoNewline)
				if dir == "up" {
					Print("up", NoNewline)
				} else {
					Print("down", NoNewline)
				}
				Print("ward.", Newline)
				return true
			}
		}
		DoWalk(dir)
		return true
	}
	if obj != nil && DirObj.Is("wall") {
		Print("Climbing the walls is to no avail.", Newline)
		return true
	}
	if Here != &Path && (obj == nil || obj == &Tree) && IsInGlobal(&Tree, Here) {
		Print("There are no climbable trees here.", Newline)
		return true
	}
	if obj == nil || obj == &Rooms {
		Print("You can't go that way.", Newline)
		return true
	}
	Print("You can't do that!", Newline)
	return true
}

func VDisembark(arg ActArg) bool {
	loc := Winner.Location()
	if DirObj == &Rooms && loc.Has(FlgVeh) {
		Perform(ActionVerb{Norm: "disembark", Orig: "disembark"}, loc, nil)
		return true
	}
	if loc != DirObj {
		Print("You're not in that!", Newline)
		return RFatal()
	}
	if Here.Has(FlgLand) {
		Print("You are on your own feet again.", Newline)
		Winner.MoveTo(Here)
		return true
	}
	Print("You realize that getting out here would be fatal.", Newline)
	return RFatal()
}

func VEnter(arg ActArg) bool {
	return DoWalk("in")
}

func VExit(arg ActArg) bool {
	if (DirObj == nil || DirObj == &Rooms) && Winner.Location().Has(FlgVeh) {
		Perform(ActionVerb{Norm: "disembark", Orig: "disembark"}, Winner.Location(), nil)
		return true
	}
	if DirObj != nil && Winner.IsIn(DirObj) {
		Perform(ActionVerb{Norm: "disembark", Orig: "disembark"}, DirObj, nil)
		return true
	}
	return DoWalk("out")
}

func VFollow(arg ActArg) bool {
	Print("You're nuts!", Newline)
	return true
}

func VLeap(arg ActArg) bool {
	if DirObj != nil {
		if !DirObj.IsIn(Here) {
			Print("That would be a good trick.", Newline)
			return true
		}
		if DirObj.Has(FlgPerson) {
			Print("The ", NoNewline)
			PrintObject(DirObj)
			Print(" is too big to jump over.", Newline)
			return true
		}
		return VSkip(ActUnk)
	}
	tx := Here.GetDir("down")
	if tx.IsSet() {
		if len(tx.NExit) > 0 || (tx.CExit != nil && !tx.CExit()) {
			Print("This was not a very safe place to try jumping.", Newline)
			return JigsUp(PickOne(JumpLoss), false)
		}
		if Here == &UpATree {
			Print("In a feat of unaccustomed daring, you manage to land on your feet without killing yourself.", Newline)
			NewLine()
			DoWalk("down")
			return true
		}
	}
	return VSkip(ActUnk)
}

func VLeave(arg ActArg) bool {
	return DoWalk("out")
}

func VStand(arg ActArg) bool {
	loc := Winner.Location()
	if !loc.Has(FlgVeh) {
		Print("You are already standing, I think.", Newline)
		return true
	}
	Perform(ActionVerb{Norm: "disembark", Orig: "disembark"}, loc, nil)
	return true
}

func VStay(arg ActArg) bool {
	Print("You will be lost without me!", Newline)
	return true
}

func VSwim(arg ActArg) bool {
	if !IsInGlobal(&GlobalWater, Here) {
		Print("Go jump in a lake!", Newline)
		return true
	}
	Print("Swimming isn't usually allowed in the ", NoNewline)
	if DirObj != &Water && DirObj != &GlobalWater {
		PrintObject(DirObj)
		Print(".", Newline)
		return true
	}
	Print("dungeon.", Newline)
	return true
}

func VThrough(arg ActArg) bool {
	return Through(nil)
}

func Through(obj *Object) bool {
	m := OtherSide(DirObj)
	if DirObj.Has(FlgDoor) && len(m) > 0 {
		DoWalk(m)
		return true
	}
	if obj != nil && DirObj.Has(FlgVeh) {
		Perform(ActionVerb{Norm: "board", Orig: "board"}, DirObj, nil)
		return true
	}
	if obj != nil || !DirObj.Has(FlgTake) {
		Print("You hit your head against the ", NoNewline)
		PrintObject(DirObj)
		Print(" as you attempt this feat.", Newline)
		return true
	}
	if DirObj.IsIn(Winner) {
		Print("That would involve quite a contortion!", Newline)
		return true
	}
	Print(PickOne(Yuks), Newline)
	return true
}

func OtherSide(dobj *Object) string {
	dirs := []string{"north", "east", "west", "south", "northeast", "northwest", "southeast", "southwest", "up", "down", "in", "out", "land"}
	for _, d := range dirs {
		dirObj := Here.GetDir(d)
		if dirObj == nil {
			continue
		}
		if dirObj.DExit == dobj {
			return d
		}
	}
	return ""
}

func VWalk(arg ActArg) bool {
	if len(Params.WalkDir) == 0 {
		Perform(ActionVerb{Norm: "walk to", Orig: "walk to"}, DirObj, nil)
		return true
	}
	props := Here.GetDir(Params.WalkDir)
	if props == nil {
		if !Lit && Prob(80, false) && Winner == &Adventurer && !Here.Has(FlgNonLand) {
			if IsSprayed {
				Print("There are odd noises in the darkness, and there is no exit in that direction.", Newline)
				return RFatal()
			}
			return JigsUp("Oh, no! You have walked into the slavering fangs of a lurking grue!", false)
		}
		Print("You can't go that way.", Newline)
		return RFatal()
	}
	// Unconditional exit
	if props.UExit {
		return Goto(props.RExit, true)
	}
	// Non-exit
	if len(props.NExit) > 0 {
		Print(props.NExit, Newline)
		return RFatal()
	}
	// Functional exit
	if props.FExit != nil {
		rm := props.FExit()
		if rm == nil {
			return RFatal()
		}
		return Goto(rm, true)
	}
	// Conditional exit
	if props.CExit != nil {
		if props.CExit() {
			return Goto(props.RExit, true)
		}
		if len(props.CExitStr) > 0 {
			Print(props.CExitStr, Newline)
			return RFatal()
		}
		Print("You can't go that way.", Newline)
		return RFatal()
	}
	if props.DExit != nil {
		if props.DExit.Has(FlgOpen) {
			return Goto(props.RExit, true)
		}
		if len(props.DExitStr) > 0 {
			Print(props.DExitStr, Newline)
			return RFatal()
		}
		Print("The ", NoNewline)
		PrintObject(props.DExit)
		Print(" is closed.", Newline)
		ThisIsIt(props.DExit)
		return RFatal()
	}
	return false
}

func VWalkAround(arg ActArg) bool {
	Print("Use compass directions for movement.", Newline)
	return true
}

func VWalkTo(arg ActArg) bool {
	if DirObj != nil && (DirObj.IsIn(Here) || IsInGlobal(DirObj, Here)) {
		Print("It's here!", Newline)
		return true
	}
	Print("You should supply a direction!", Newline)
	return true
}

func DoWalk(dir string) bool {
	Params.WalkDir = dir
	dirObj := ToDirObj(dir)
	if Perform(ActionVerb{Norm: "walk", Orig: "walk"}, dirObj, nil) == PerfHndld {
		return true
	}
	return false
}

func NoGoTell(av Flags, wloc *Object) {
	if av != FlgUnk {
		Print("You can't go there in a ", NoNewline)
		PrintObject(wloc)
		Print(".", Newline)
		return
	}
	Print("You can't go there without a vehicle.", Newline)
}

func Goto(rm *Object, isV bool) bool {
	lb := rm.Has(FlgLand) || rm.Has(FlgRLand)
	wloc := Winner.Location()
	var av Flags
	olit := Lit
	ohere := Here
	if wloc.Has(FlgVeh) {
		av = wloc.VehType
	}
	if !lb && av == FlgUnk {
		NoGoTell(av, wloc)
		return false
	}
	if !lb && av != FlgUnk && !rm.Has(av) {
		NoGoTell(av, wloc)
		return false
	}
	if Here.Has(FlgLand) && lb && av != FlgUnk && av != FlgLand && !rm.Has(av) {
		NoGoTell(av, wloc)
		return false
	}
	if rm.Has(FlgKludge) {
		Print(rm.LongDesc, Newline)
		return false
	}
	if lb && !Here.Has(FlgLand) && !Dead && wloc.Has(FlgVeh) {
		Print("The ", NoNewline)
		PrintObject(wloc)
		Print(" comes to a rest on the shore.", Newline)
		NewLine()
	}
	if av != FlgUnk {
		wloc.MoveTo(rm)
	} else {
		Winner.MoveTo(rm)
	}
	Here = rm
	Lit = IsLit(Here, true)
	if !olit && !Lit && Prob(80, false) {
		if !IsSprayed {
			Print("Oh, no! A lurking grue slithered into the ", NoNewline)
			if Winner.Location().Has(FlgVeh) {
				PrintObject(Winner.Location())
			} else {
				Print("room", NoNewline)
			}
			JigsUp(" and devoured you!", false)
			return true
		}
		Print("There are sinister gurgling noises in the darkness all around you!", Newline)
	}
	if !Lit && Winner == &Adventurer {
		Print("You have moved into a dark place.", Newline)
		Params.Continue = NumUndef
	}
	if Here.Action != nil {
		Here.Action(ActEnter)
	}
	ScoreObj(rm)
	// If the room's enter action teleported the player elsewhere, stop here.
	if Here != rm {
		return true
	}
	if Winner != &Adventurer && Adventurer.IsIn(ohere) {
		Print("The ", NoNewline)
		PrintObject(Winner)
		Print(" leaves the room.", Newline)
		return true
	}
	if Here == ohere && Here == &EnteranceToHades {
		return true
	}
	if isV && Winner == &Adventurer {
		VFirstLook(ActUnk)
	}
	return true
}

func VCross(arg ActArg) bool {
	Print("You can't cross that!", Newline)
	return true
}

func ToDirObj(dir string) *Object {
	return nil
}
