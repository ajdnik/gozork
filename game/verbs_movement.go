package game

import . "github.com/ajdnik/gozork/engine"

func preBoard(arg ActionArg) bool {
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
	if G.DirObj == &water || G.DirObj == &globalWater {
		Perform(ActionVerb{Norm: "swim", Orig: "swim"}, G.DirObj, nil)
		return true
	}
	Printf("You have a theory on how to board a %s, perhaps?\n", G.DirObj.Desc)
	return RFatal()
}

func vBoard(arg ActionArg) bool {
	Printf("You are now in the %s.\n", G.DirObj.Desc)
	G.Winner.MoveTo(G.DirObj)
	if G.DirObj.Action != nil {
		G.DirObj.Action(ActEnter)
	}
	return true
}

func vClimbDown(arg ActionArg) bool {
	return vClimbFcn(Down, G.DirObj)
}

func vClimbFoo(arg ActionArg) bool {
	return vClimbFcn(Up, G.DirObj)
}

func vClimbOn(arg ActionArg) bool {
	if !G.DirObj.Has(FlgVeh) {
		Printf("You can't climb onto the %s.\n", G.DirObj.Desc)
		return true
	}
	Perform(ActionVerb{Norm: "board", Orig: "board"}, G.DirObj, nil)
	return true
}

func vClimbUp(arg ActionArg) bool {
	return vClimbFcn(Up, nil)
}

func vClimbFcn(dir Direction, obj *Object) bool {
	if obj != nil && G.DirObj != &rooms {
		obj = G.DirObj
	}
	if tx := G.Here.GetExit(dir); tx != nil {
		if obj != nil {
			if len(tx.NExit) > 0 || ((tx.CExit != nil || tx.DExit != nil || tx.UExit) && !IsInGlobal(G.DirObj, tx.RExit)) {
				Printf("The %s do", obj.Desc)
				if obj != &stairs {
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
		doWalk(dir)
		return true
	}
	if obj != nil && G.DirObj.Is("wall") {
		Printf("Climbing the walls is to no avail.\n")
		return true
	}
	if G.Here != &path && (obj == nil || obj == &tree) && IsInGlobal(&tree, G.Here) {
		Printf("There are no climbable trees here.\n")
		return true
	}
	if obj == nil || obj == &rooms {
		Printf("You can't go that way.\n")
		return true
	}
	Printf("You can't do that!\n")
	return true
}

func vDisembark(arg ActionArg) bool {
	loc := G.Winner.Location()
	if G.DirObj == &rooms && loc.Has(FlgVeh) {
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

func vEnter(arg ActionArg) bool {
	return doWalk(In)
}

func vExit(arg ActionArg) bool {
	if (G.DirObj == nil || G.DirObj == &rooms) && G.Winner.Location().Has(FlgVeh) {
		Perform(ActionVerb{Norm: "disembark", Orig: "disembark"}, G.Winner.Location(), nil)
		return true
	}
	if G.DirObj != nil && G.Winner.IsIn(G.DirObj) {
		Perform(ActionVerb{Norm: "disembark", Orig: "disembark"}, G.DirObj, nil)
		return true
	}
	return doWalk(Out)
}

func vFollow(arg ActionArg) bool {
	Printf("You're nuts!\n")
	return true
}

func vLeap(arg ActionArg) bool {
	if G.DirObj != nil {
		if !G.DirObj.IsIn(G.Here) {
			Printf("That would be a good trick.\n")
			return true
		}
		if G.DirObj.Has(FlgPerson) {
			Printf("The %s is too big to jump over.\n", G.DirObj.Desc)
			return true
		}
		return vSkip(ActUnk)
	}
	tx := G.Here.GetExit(Down)
	if tx != nil && tx.IsSet() {
		if len(tx.NExit) > 0 || (tx.CExit != nil && !tx.CExit()) {
			Printf("This was not a very safe place to try jumping.\n")
			return jigsUp(PickOne(jumpLoss), false)
		}
		if G.Here == &upATree {
			Printf("In a feat of unaccustomed daring, you manage to land on your feet without killing yourself.\n\n")
			doWalk(Down)
			return true
		}
	}
	return vSkip(ActUnk)
}

func vLeave(arg ActionArg) bool {
	return doWalk(Out)
}

func vStand(arg ActionArg) bool {
	loc := G.Winner.Location()
	if !loc.Has(FlgVeh) {
		Printf("You are already standing, I think.\n")
		return true
	}
	Perform(ActionVerb{Norm: "disembark", Orig: "disembark"}, loc, nil)
	return true
}

func vStay(arg ActionArg) bool {
	Printf("You will be lost without me!\n")
	return true
}

func vSwim(arg ActionArg) bool {
	if !IsInGlobal(&globalWater, G.Here) {
		Printf("Go jump in a lake!\n")
		return true
	}
	Printf("Swimming isn't usually allowed in the ")
	if G.DirObj != &water && G.DirObj != &globalWater {
		Printf("%s.\n", G.DirObj.Desc)
		return true
	}
	Printf("dungeon.\n")
	return true
}

func vThrough(arg ActionArg) bool {
	return through(nil)
}

func through(obj *Object) bool {
	m, ok := otherSide(G.DirObj)
	if G.DirObj.Has(FlgDoor) && ok {
		doWalk(m)
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
	Printf("%s\n", PickOne(yuks))
	return true
}

func otherSide(dobj *Object) (Direction, bool) {
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

func vWalk(arg ActionArg) bool {
	if !G.Params.HasWalkDir {
		Perform(ActionVerb{Norm: "walk to", Orig: "walk to"}, G.DirObj, nil)
		return true
	}
	props := G.Here.GetExit(G.Params.WalkDir)
	if props == nil {
		if !G.Lit && Prob(80, false) && G.Winner == &adventurer && !G.Here.Has(FlgNonLand) {
			if gD().IsSprayed {
				Printf("There are odd noises in the darkness, and there is no exit in that direction.\n")
				return RFatal()
			}
			return jigsUp("Oh, no! You have walked into the slavering fangs of a lurking grue!", false)
		}
		Printf("You can't go that way.\n")
		return RFatal()
	}
	// Unconditional exit
	if props.UExit {
		return gotoRoom(props.RExit, true)
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
		return gotoRoom(rm, true)
	}
	// Conditional exit
	if props.CExit != nil {
		if props.CExit() {
			return gotoRoom(props.RExit, true)
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
			return gotoRoom(props.RExit, true)
		}
		if len(props.DExitStr) > 0 {
			Printf("%s\n", props.DExitStr)
			return RFatal()
		}
		Printf("The %s is closed.\n", props.DExit.Desc)
		thisIsIt(props.DExit)
		return RFatal()
	}
	return false
}

func vWalkAround(arg ActionArg) bool {
	Printf("Use compass directions for movement.\n")
	return true
}

func vWalkTo(arg ActionArg) bool {
	if G.DirObj != nil && (G.DirObj.IsIn(G.Here) || IsInGlobal(G.DirObj, G.Here)) {
		Printf("it's here!\n")
		return true
	}
	Printf("You should supply a direction!\n")
	return true
}

func doWalk(dir Direction) bool {
	G.Params.WalkDir = dir
	G.Params.HasWalkDir = true
	return Perform(ActionVerb{Norm: "walk", Orig: "walk"}, nil, nil) == PerfHndld
}

func noGoTell(av Flags, wloc *Object) {
	if av != FlgUnk {
		Printf("You can't go there in a %s.\n", wloc.Desc)
		return
	}
	Printf("You can't go there without a vehicle.\n")
}

func gotoRoom(rm *Object, isV bool) bool {
	lb := rm.Has(FlgLand) || rm.Has(FlgRLand)
	wloc := G.Winner.Location()
	var av Flags
	olit := G.Lit
	ohere := G.Here
	if wloc.Has(FlgVeh) {
		av = wloc.GetVehType()
	}
	if !lb && av == FlgUnk {
		noGoTell(av, wloc)
		return false
	}
	if !lb && av != FlgUnk && !rm.Has(av) {
		noGoTell(av, wloc)
		return false
	}
	if G.Here.Has(FlgLand) && lb && av != FlgUnk && av != FlgLand && !rm.Has(av) {
		noGoTell(av, wloc)
		return false
	}
	if rm.Has(FlgKludge) {
		Printf("%s\n", rm.LongDesc)
		return false
	}
	if lb && !G.Here.Has(FlgLand) && !gD().Dead && wloc.Has(FlgVeh) {
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
		if !gD().IsSprayed {
			Printf("Oh, no! A lurking grue slithered into the ")
			if G.Winner.Location().Has(FlgVeh) {
				Printf("%s", G.Winner.Location().Desc)
			} else {
				Printf("room")
			}
			jigsUp(" and devoured you!", false)
			return true
		}
		Printf("There are sinister gurgling noises in the darkness all around you!\n")
	}
	if !G.Lit && G.Winner == &adventurer {
		Printf("You have moved into a dark place.\n")
		G.Params.Continue = NumUndef
	}
	if G.Here.Action != nil {
		G.Here.Action(ActEnter)
	}
	scoreObj(rm)
	// If the room's enter action teleported the player elsewhere, stop here.
	if G.Here != rm {
		return true
	}
	if G.Winner != &adventurer && adventurer.IsIn(ohere) {
		Printf("The %s leaves the room.\n", G.Winner.Desc)
		return true
	}
	if G.Here == ohere && G.Here == &enteranceToHades {
		return true
	}
	if isV && G.Winner == &adventurer {
		vFirstLook(ActUnk)
	}
	return true
}

func vCross(arg ActionArg) bool {
	Printf("You can't cross that!\n")
	return true
}
