package zork


func PreBoard(arg ActArg) bool {
	if G.DirObj.Has(FlgVeh) {
		if !G.DirObj.IsIn(G.Here) {
			Printf("The %s must be on the ground to be boarded.\n", G.DirObj.Desc)
			return RFatal()
		}
		if av := G.Winner.Location(); av != nil && av.Has(FlgVeh) {
			Printf("You are already in the %s!\n", av.Desc)
			return RFatal()
		}
		return false
	}
	if G.DirObj == &Water || G.DirObj == &GlobalWater {
		Perform(ActionVerb{Norm: "swim", Orig: "swim"}, G.DirObj, nil)
		return true
	}
	Printf("You have a theory on how to board a %s, perhaps?\n", G.DirObj.Desc)
	return RFatal()
}

func VBoard(arg ActArg) bool {
	Printf("You are now in the %s.\n", G.DirObj.Desc)
	G.Winner.MoveTo(G.DirObj)
	if G.DirObj.Action != nil {
		G.DirObj.Action(ActEnter)
	}
	return true
}

func VClimbDown(arg ActArg) bool {
	return VClimbFcn(Down, G.DirObj)
}

func VClimbFoo(arg ActArg) bool {
	return VClimbFcn(Up, G.DirObj)
}

func VClimbOn(arg ActArg) bool {
	if !G.DirObj.Has(FlgVeh) {
		Printf("You can't climb onto the %s.\n", G.DirObj.Desc)
		return true
	}
	Perform(ActionVerb{Norm: "board", Orig: "board"}, G.DirObj, nil)
	return true
}

func VClimbUp(arg ActArg) bool {
	return VClimbFcn(Up, nil)
}

func VClimbFcn(dir Direction, obj *Object) bool {
	if obj != nil && G.DirObj != &Rooms {
		obj = G.DirObj
	}
	if tx := G.Here.GetExit(dir); tx != nil {
		if obj != nil {
			if len(tx.NExit) > 0 || ((tx.CExit != nil || tx.DExit != nil || tx.UExit) && !IsInGlobal(G.DirObj, tx.RExit)) {
				Printf("The %s do", obj.Desc)
				if obj != &Stairs {
					Printf("es")
				}
				Printf("n't lead ")
				if dir == Up {
					Printf("up")
				} else {
					Printf("down")
				}
				Printf("ward.\n")
				return true
			}
		}
		DoWalk(dir)
		return true
	}
	if obj != nil && G.DirObj.Is("wall") {
		Printf("Climbing the walls is to no avail.\n")
		return true
	}
	if G.Here != &Path && (obj == nil || obj == &Tree) && IsInGlobal(&Tree, G.Here) {
		Printf("There are no climbable trees here.\n")
		return true
	}
	if obj == nil || obj == &Rooms {
		Printf("You can't go that way.\n")
		return true
	}
	Printf("You can't do that!\n")
	return true
}

func VDisembark(arg ActArg) bool {
	loc := G.Winner.Location()
	if G.DirObj == &Rooms && loc.Has(FlgVeh) {
		Perform(ActionVerb{Norm: "disembark", Orig: "disembark"}, loc, nil)
		return true
	}
	if loc != G.DirObj {
		Printf("You're not in that!\n")
		return RFatal()
	}
	if G.Here.Has(FlgLand) {
		Printf("You are on your own feet again.\n")
		G.Winner.MoveTo(G.Here)
		return true
	}
	Printf("You realize that getting out here would be fatal.\n")
	return RFatal()
}

func VEnter(arg ActArg) bool {
	return DoWalk(In)
}

func VExit(arg ActArg) bool {
	if (G.DirObj == nil || G.DirObj == &Rooms) && G.Winner.Location().Has(FlgVeh) {
		Perform(ActionVerb{Norm: "disembark", Orig: "disembark"}, G.Winner.Location(), nil)
		return true
	}
	if G.DirObj != nil && G.Winner.IsIn(G.DirObj) {
		Perform(ActionVerb{Norm: "disembark", Orig: "disembark"}, G.DirObj, nil)
		return true
	}
	return DoWalk(Out)
}

func VFollow(arg ActArg) bool {
	Printf("You're nuts!\n")
	return true
}

func VLeap(arg ActArg) bool {
	if G.DirObj != nil {
		if !G.DirObj.IsIn(G.Here) {
			Printf("That would be a good trick.\n")
			return true
		}
		if G.DirObj.Has(FlgPerson) {
			Printf("The %s is too big to jump over.\n", G.DirObj.Desc)
			return true
		}
		return VSkip(ActUnk)
	}
	tx := G.Here.GetExit(Down)
	if tx != nil && tx.IsSet() {
		if len(tx.NExit) > 0 || (tx.CExit != nil && !tx.CExit()) {
			Printf("This was not a very safe place to try jumping.\n")
			return JigsUp(PickOne(JumpLoss), false)
		}
		if G.Here == &UpATree {
			Printf("In a feat of unaccustomed daring, you manage to land on your feet without killing yourself.\n\n")
			DoWalk(Down)
			return true
		}
	}
	return VSkip(ActUnk)
}

func VLeave(arg ActArg) bool {
	return DoWalk(Out)
}

func VStand(arg ActArg) bool {
	loc := G.Winner.Location()
	if !loc.Has(FlgVeh) {
		Printf("You are already standing, I think.\n")
		return true
	}
	Perform(ActionVerb{Norm: "disembark", Orig: "disembark"}, loc, nil)
	return true
}

func VStay(arg ActArg) bool {
	Printf("You will be lost without me!\n")
	return true
}

func VSwim(arg ActArg) bool {
	if !IsInGlobal(&GlobalWater, G.Here) {
		Printf("Go jump in a lake!\n")
		return true
	}
	Printf("Swimming isn't usually allowed in the ")
	if G.DirObj != &Water && G.DirObj != &GlobalWater {
		Printf("%s.\n", G.DirObj.Desc)
		return true
	}
	Printf("dungeon.\n")
	return true
}

func VThrough(arg ActArg) bool {
	return Through(nil)
}

func Through(obj *Object) bool {
	m, ok := OtherSide(G.DirObj)
	if G.DirObj.Has(FlgDoor) && ok {
		DoWalk(m)
		return true
	}
	if obj != nil && G.DirObj.Has(FlgVeh) {
		Perform(ActionVerb{Norm: "board", Orig: "board"}, G.DirObj, nil)
		return true
	}
	if obj != nil || !G.DirObj.Has(FlgTake) {
		Printf("You hit your head against the %s as you attempt this feat.\n", G.DirObj.Desc)
		return true
	}
	if G.DirObj.IsIn(G.Winner) {
		Printf("That would involve quite a contortion!\n")
		return true
	}
	Printf("%s\n", PickOne(Yuks))
	return true
}

func OtherSide(dobj *Object) (Direction, bool) {
	for _, d := range AllDirections {
		dp := G.Here.GetExit(d)
		if dp == nil {
			continue
		}
		if dp.DExit == dobj {
			return d, true
		}
	}
	return 0, false
}

func VWalk(arg ActArg) bool {
	if !G.Params.HasWalkDir {
		Perform(ActionVerb{Norm: "walk to", Orig: "walk to"}, G.DirObj, nil)
		return true
	}
	props := G.Here.GetExit(G.Params.WalkDir)
	if props == nil {
		if !G.Lit && Prob(80, false) && G.Winner == &Adventurer && !G.Here.Has(FlgNonLand) {
			if G.IsSprayed {
				Printf("There are odd noises in the darkness, and there is no exit in that direction.\n")
				return RFatal()
			}
			return JigsUp("Oh, no! You have walked into the slavering fangs of a lurking grue!", false)
		}
		Printf("You can't go that way.\n")
		return RFatal()
	}
	// Unconditional exit
	if props.UExit {
		return Goto(props.RExit, true)
	}
	// Non-exit
	if len(props.NExit) > 0 {
		Printf("%s\n", props.NExit)
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
			Printf("%s\n", props.CExitStr)
			return RFatal()
		}
		Printf("You can't go that way.\n")
		return RFatal()
	}
	if props.DExit != nil {
		if props.DExit.Has(FlgOpen) {
			return Goto(props.RExit, true)
		}
		if len(props.DExitStr) > 0 {
			Printf("%s\n", props.DExitStr)
			return RFatal()
		}
		Printf("The %s is closed.\n", props.DExit.Desc)
		ThisIsIt(props.DExit)
		return RFatal()
	}
	return false
}

func VWalkAround(arg ActArg) bool {
	Printf("Use compass directions for movement.\n")
	return true
}

func VWalkTo(arg ActArg) bool {
	if G.DirObj != nil && (G.DirObj.IsIn(G.Here) || IsInGlobal(G.DirObj, G.Here)) {
		Printf("It's here!\n")
		return true
	}
	Printf("You should supply a direction!\n")
	return true
}

func DoWalk(dir Direction) bool {
	G.Params.WalkDir = dir
	G.Params.HasWalkDir = true
	if Perform(ActionVerb{Norm: "walk", Orig: "walk"}, nil, nil) == PerfHndld {
		return true
	}
	return false
}

func NoGoTell(av Flags, wloc *Object) {
	if av != FlgUnk {
		Printf("You can't go there in a %s.\n", wloc.Desc)
		return
	}
	Printf("You can't go there without a vehicle.\n")
}

func Goto(rm *Object, isV bool) bool {
	lb := rm.Has(FlgLand) || rm.Has(FlgRLand)
	wloc := G.Winner.Location()
	var av Flags
	olit := G.Lit
	ohere := G.Here
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
	if G.Here.Has(FlgLand) && lb && av != FlgUnk && av != FlgLand && !rm.Has(av) {
		NoGoTell(av, wloc)
		return false
	}
	if rm.Has(FlgKludge) {
		Printf("%s\n", rm.LongDesc)
		return false
	}
	if lb && !G.Here.Has(FlgLand) && !G.Dead && wloc.Has(FlgVeh) {
		Printf("The %s comes to a rest on the shore.\n\n", wloc.Desc)
	}
	if av != FlgUnk {
		wloc.MoveTo(rm)
	} else {
		G.Winner.MoveTo(rm)
	}
	G.Here = rm
	G.Lit = IsLit(G.Here, true)
	if !olit && !G.Lit && Prob(80, false) {
		if !G.IsSprayed {
			Printf("Oh, no! A lurking grue slithered into the ")
			if G.Winner.Location().Has(FlgVeh) {
				Printf("%s", G.Winner.Location().Desc)
			} else {
				Printf("room")
			}
			JigsUp(" and devoured you!", false)
			return true
		}
		Printf("There are sinister gurgling noises in the darkness all around you!\n")
	}
	if !G.Lit && G.Winner == &Adventurer {
		Printf("You have moved into a dark place.\n")
		G.Params.Continue = NumUndef
	}
	if G.Here.Action != nil {
		G.Here.Action(ActEnter)
	}
	ScoreObj(rm)
	// If the room's enter action teleported the player elsewhere, stop here.
	if G.Here != rm {
		return true
	}
	if G.Winner != &Adventurer && Adventurer.IsIn(ohere) {
		Printf("The %s leaves the room.\n", G.Winner.Desc)
		return true
	}
	if G.Here == ohere && G.Here == &EnteranceToHades {
		return true
	}
	if isV && G.Winner == &Adventurer {
		VFirstLook(ActUnk)
	}
	return true
}

func VCross(arg ActArg) bool {
	Printf("You can't cross that!\n")
	return true
}

