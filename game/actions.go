package game

import . "github.com/ajdnik/gozork/engine"

func openClose(obj *Object, strOpn, strCls string) {
	switch G.ActVerb.Norm {
	case "open":
		if obj.Has(FlgOpen) {
			Printf("%s\n", PickOne(dummy))
		} else {
			Printf("%s\n", strOpn)
			obj.Give(FlgOpen)
		}
	case "close":
		if obj.Has(FlgOpen) {
			Printf("%s\n", strCls)
			obj.Take(FlgOpen)
		} else {
			Printf("%s\n", PickOne(dummy))
		}
	}
}

func leavesAppear() bool {
	if !grate.Has(FlgOpen) && !gD().GrateRevealed {
		switch G.ActVerb.Norm {
		case "move", "take":
			Printf("In disturbing the pile of leaves, a grating is revealed.\n")
		default:
			Printf("With the leaves moved, a grating is revealed.\n")
		}
		grate.Take(FlgInvis)
		gD().GrateRevealed = true
	}
	return false
}

func fweep(n int) {
	for i := 0; i < n; i++ {
		Printf("    fweep!\n")
	}
	Printf("\n")
}

func flyMe() bool {
	fweep(4)
	Printf("The bat grabs you by the scruff of your neck and lifts you away....\n\n")
	dest := batDrops[G.Rand.Intn(len(batDrops))]
	moveToRoom(dest, false)
	if G.Here != &enteranceToHades {
		vFirstLook(ActUnk)
	}
	return true
}

func touchAll(obj *Object) {
	for _, child := range obj.Children {
		child.Give(FlgTouch)
		if child.HasChildren() {
			touchAll(child)
		}
	}
}

func otvalFrob(o *Object) int {
	score := 0
	for _, child := range o.Children {
		score += child.GetTValue()
		if child.HasChildren() {
			score += otvalFrob(child)
		}
	}
	return score
}

func integralPart() {
	Printf("it is an integral part of the control panel.\n")
}

func withTell(obj *Object) {
	Printf("With a %s?\n", obj.Desc)
}

func badEgg() {
	if canary.IsIn(&egg) {
		Printf(" %s", brokenCanary.FirstDesc)
	} else {
		removeCarefully(&brokenCanary)
	}
	brokenEgg.MoveTo(egg.Location())
	removeCarefully(&egg)
}

func slider(obj *Object) {
	if obj.Has(FlgTake) {
		Printf("The %s falls into the slide and is gone.\n", obj.Desc)
		if obj == &water {
			removeCarefully(obj)
		} else {
			obj.MoveTo(&cellar)
		}
	} else {
		Printf("%s\n", PickOne(yuks))
	}
}

func forestRoomQ() bool {
	return G.Here == &forest1 || G.Here == &forest2 || G.Here == &forest3 ||
		G.Here == &path || G.Here == &upATree
}

// stealJunk and dropJunk are defined in the iThief section below

// ================================================================
// SCORE
// ================================================================

func vScore(arg ActionArg) bool {
	Printf("Your score is %d (total of %d points), in %d", G.Score, ScoreMax, G.Moves)
	if G.Moves == 1 {
		Printf(" move.")
	} else {
		Printf(" moves.")
	}
	Printf("\nThis gives you the rank of ")
	switch {
	case G.Score == ScoreMax:
		Printf("Master adventurer")
	case G.Score > 330:
		Printf("Wizard")
	case G.Score > 300:
		Printf("Master")
	case G.Score > 200:
		Printf("adventurer")
	case G.Score > 100:
		Printf("Junior adventurer")
	case G.Score > 50:
		Printf("Novice adventurer")
	case G.Score > 25:
		Printf("Amateur adventurer")
	default:
		Printf("Beginner")
	}
	Printf(".\n")
	return true
}

func vDiagnose(arg ActionArg) bool {
	ms := fightStrength(false)
	wd := G.Winner.GetStrength()
	rs := ms + wd
	// Check if healing is active
	cureActive := false
	for i := len(G.QueueItms) - 1; i >= 0; i-- {
		if G.QueueItms[i].Key == "iCure" && G.QueueItms[i].Run {
			cureActive = true
			break
		}
	}
	if !cureActive {
		wd = 0
	} else {
		wd = -wd
	}
	if wd == 0 {
		Printf("You are in perfect health.")
	} else {
		Printf("You have ")
		switch {
		case wd == 1:
			Printf("a light wound,")
		case wd == 2:
			Printf("a serious wound,")
		case wd == 3:
			Printf("several wounds,")
		default:
			Printf("serious wounds,")
		}
	}
	if wd != 0 {
		Printf(" which will be cured after some moves.")
	}
	Printf("\nYou can ")
	switch {
	case rs == 0:
		Printf("expect death soon")
	case rs == 1:
		Printf("be killed by one more light wound")
	case rs == 2:
		Printf("be killed by a serious wound")
	case rs == 3:
		Printf("survive one serious wound")
	default:
		Printf("survive several wounds")
	}
	Printf(".\n")
	if gD().Deaths > 0 {
		Printf("You have been killed ")
		if gD().Deaths == 1 {
			Printf("once")
		} else {
			Printf("twice")
		}
		Printf(".\n")
	}
	return true
}

// ================================================================
// DEATH AND RESTART
// ================================================================

func jigsUp(desc string, isPlyr bool) bool {
	G.Winner = &adventurer
	if gD().Dead {
		Printf("\nIt takes a talented person to be killed while already dead. YOU are such a talent. Unfortunately, it takes a talented person to deal with it. I am not such a talent. Sorry.\n")
		return finish()
	}
	Printf("%s\n", desc)
	if !G.Lucky {
		Printf("Bad luck, huh?\n")
	}
	scoreUpd(-10)
	Printf("\n    ****  You have died  ****\n\n")
	if G.Winner.Location().Has(FlgVeh) {
		G.Winner.MoveTo(G.Here)
	}
	if gD().Deaths >= 2 {
		Printf("You clearly are a suicidal maniac. We don't allow psychotics in the cave, since they may harm other adventurers. Your remains will be installed in the Land of the Living Dead, where your fellow adventurers may gloat over them.\n")
		return finish()
	}
	gD().Deaths++
	G.Winner.MoveTo(G.Here)
	if southTemple.Has(FlgTouch) {
		Printf("As you take your last breath, you feel relieved of your burdens. The feeling passes as you find yourself before the gates of Hell, where the spirits jeer at you and deny you entry. Your senses are disturbed. The objects in the dungeon appear indistinct, bleached of color, even unreal.\n\n")
		gD().Dead = true
		gD().TrollFlag = true
		G.AlwaysLit = true
		G.Winner.Action = deadFunction
		moveToRoom(&enteranceToHades, true)
	} else {
		Printf("Now, let's take a look here...\nWell, you probably deserve another chance. I can't quite fix you up completely, but you can't have everything.\n\n")
		moveToRoom(&forest1, true)
	}
	trapDoor.Take(FlgTouch)
	G.Params.Continue = NumUndef
	randomizeObjects()
	killInterrupts()
	return false
}

func randomizeObjects() {
	if lamp.IsIn(G.Winner) {
		lamp.MoveTo(&livingRoom)
	}
	if coffin.IsIn(G.Winner) {
		coffin.MoveTo(&egyptRoom)
	}
	sword.SetTValue(0)
	// Copy children before iterating since MoveTo modifies the slice.
	children := make([]*Object, len(G.Winner.Children))
	copy(children, G.Winner.Children)
	for _, child := range children {
		if child.GetTValue() <= 0 {
			child.MoveTo(Random(aboveGround))
			continue
		}
		for _, r := range rooms.Children {
			if r.Has(FlgLand) && !r.Has(FlgOn) && Prob(50, false) {
				child.MoveTo(r)
				break
			}
		}
	}
}

func killInterrupts() bool {
	QueueInt("iXb", false).Run = false
	QueueInt("iXc", false).Run = false
	QueueInt("iCyclops", false).Run = false
	QueueInt("iLantern", false).Run = false
	QueueInt("iCandles", false).Run = false
	QueueInt("iSword", false).Run = false
	QueueInt("iForestRandom", false).Run = false
	QueueInt("iMatch", false).Run = false
	match.Take(FlgOn)
	return true
}

// ================================================================
// THE WHITE HOUSE
// ================================================================

func westHouseFcn(arg ActionArg) bool {
	if arg == ActLook {
		Printf("You are standing in an open field west of a white house, with a boarded front door.")
		if gD().WonGame {
			Printf(" A secret path leads southwest into the forest.")
		}
		Printf("\n")
		return true
	}
	return false
}

func eastHouseFcn(arg ActionArg) bool {
	if arg == ActLook {
		Printf("You are behind the white house. A path leads into the forest to the east. In one corner of the house there is a small window which is ")
		if kitchenWindow.Has(FlgOpen) {
			Printf("open.\n")
		} else {
			Printf("slightly ajar.\n")
		}
		return true
	}
	return false
}

func whiteHouseFcn(arg ActionArg) bool {
	if G.Here == &kitchen || G.Here == &livingRoom || G.Here == &attic {
		switch G.ActVerb.Norm {
		case "find":
			Printf("Why not find your brains?\n")
			return true
		case "walk around":
			goNext(inHouseAround)
			return true
		}
	} else if G.Here != &westOfHouse && G.Here != &northOfHouse && G.Here != &eastOfHouse && G.Here != &southOfHouse {
		switch G.ActVerb.Norm {
		case "find":
			if G.Here == &clearing {
				Printf("it seems to be to the west.\n")
				return true
			}
			Printf("it was here just a minute ago....\n")
			return true
		default:
			Printf("You're not at the house.\n")
			return true
		}
	} else {
		switch G.ActVerb.Norm {
		case "find":
			Printf("it's right here! Are you blind or something?\n")
			return true
		case "walk around":
			goNext(houseAround)
			return true
		case "examine":
			Printf("The house is a beautiful colonial house which is painted white. it is clear that the owners must have been extremely wealthy.\n")
			return true
		case "through", "open":
			if G.Here == &eastOfHouse {
				if kitchenWindow.Has(FlgOpen) {
					return moveToRoom(&kitchen, true)
				}
				Printf("The window is closed.\n")
				thisIsIt(&kitchenWindow)
				return true
			}
			Printf("I can't see how to get in from here.\n")
			return true
		case "burn":
			Printf("You must be joking.\n")
			return true
		}
	}
	return false
}

func goNext(tbl map[*Object]*Object) int {
	val, ok := tbl[G.Here]
	if !ok {
		return NumUndef
	}
	if !moveToRoom(val, true) {
		return 2
	}
	return 1
}

func boardFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "take", "examine":
		Printf("The boards are securely fastened.\n")
		return true
	}
	return false
}

func teethFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "brush":
		if G.DirObj == &teeth {
			if G.IndirObj == &putty && putty.IsIn(G.Winner) {
				jigsUp("Well, you seem to have been brushing your teeth with some sort of glue. As a result, your mouth gets glued together (with your nose) and you die of respiratory failure.", false)
				return true
			}
			if G.IndirObj == nil {
				Printf("Dental hygiene is highly recommended, but I'm not sure what you want to brush them with.\n")
				return true
			}
			Printf("A nice idea, but with a %s?\n", G.IndirObj.Desc)
			return true
		}
	}
	return false
}

func graniteWallFcn(arg ActionArg) bool {
	if G.Here == &northTemple {
		switch G.ActVerb.Norm {
		case "find":
			Printf("The west wall is solid granite here.\n")
			return true
		case "take", "raise", "lower":
			Printf("it's solid granite.\n")
			return true
		}
	} else if G.Here == &treasureRoom {
		switch G.ActVerb.Norm {
		case "find":
			Printf("The east wall is solid granite here.\n")
			return true
		case "take", "raise", "lower":
			Printf("it's solid granite.\n")
			return true
		}
	} else if G.Here == &slideRoom {
		switch G.ActVerb.Norm {
		case "find", "read":
			Printf("it only SAYS \"Granite wall\".\n")
			return true
		default:
			Printf("The wall isn't granite.\n")
			return true
		}
	} else {
		Printf("There is no granite wall here.\n")
		return true
	}
	return false
}

func songbirdFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "find", "take":
		Printf("The songbird is not here but is probably nearby.\n")
		return true
	case "listen":
		Printf("You can't hear the songbird now.\n")
		return true
	case "follow":
		Printf("it can't be followed.\n")
		return true
	default:
		Printf("You can't see any songbird here.\n")
		return true
	}
}

func mountainRangeFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "climb up", "climb down", "climb":
		Printf("Don't you believe me? The mountains are impassable!\n")
		return true
	}
	return false
}

func forestFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "walk around":
		if G.Here == &westOfHouse || G.Here == &northOfHouse || G.Here == &southOfHouse || G.Here == &eastOfHouse {
			Printf("You aren't even in the forest.\n")
			return true
		}
		goNext(forestAround)
		return true
	case "disembark":
		Printf("You will have to specify a direction.\n")
		return true
	case "find":
		Printf("You cannot see the forest for the trees.\n")
		return true
	case "listen":
		Printf("The pines and the hemlocks seem to be murmuring.\n")
		return true
	}
	return false
}

func kitchenWindowFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "open", "close":
		gD().KitchenWindowFlag = true
		openClose(&kitchenWindow,
			"With great effort, you open the window far enough to allow entry.",
			"The window closes (more easily than it opened).")
		return true
	case "examine":
		if !gD().KitchenWindowFlag {
			Printf("The window is slightly ajar, but not enough to allow entry.\n")
			return true
		}
	case "walk", "board", "through":
		if G.Here == &kitchen {
			doWalk(East)
		} else {
			doWalk(West)
		}
		return true
	case "look inside":
		Printf("You can see ")
		if G.Here == &kitchen {
			Printf("a clear area leading towards a forest.\n")
		} else {
			Printf("what appears to be a kitchen.\n")
		}
		return true
	}
	return false
}

func chimneyFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "examine":
		Printf("The chimney leads ")
		if G.Here == &kitchen {
			Printf("down")
		} else {
			Printf("up")
		}
		Printf("ward, and looks climbable.\n")
		return true
	}
	return false
}

func ghostsFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "tell":
		Printf("The spirits jeer loudly and ignore you.\n")
		G.Params.Continue = NumUndef
		return true
	case "exorcise":
		Printf("Only the ceremony itself has any effect.\n")
		return true
	case "attack", "mung":
		if G.DirObj == &ghosts {
			Printf("How can you attack a spirit with material objects?\n")
			return true
		}
	}
	Printf("You seem unable to interact with these spirits.\n")
	return true
}

func basketFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "raise":
		if gD().CageTop {
			Printf("%s\n", PickOne(dummy))
		} else {
			raisedBasket.MoveTo(&shaftRoom)
			loweredBasket.MoveTo(&lowerShaft)
			gD().CageTop = true
			thisIsIt(&raisedBasket)
			Printf("The basket is raised to the top of the shaft.\n")
		}
		return true
	case "lower":
		if !gD().CageTop {
			Printf("%s\n", PickOne(dummy))
		} else {
			raisedBasket.MoveTo(&lowerShaft)
			loweredBasket.MoveTo(&shaftRoom)
			thisIsIt(&loweredBasket)
			Printf("The basket is lowered to the bottom of the shaft.\n")
			gD().CageTop = false
			if G.Lit && !IsLit(G.Here, true) {
				G.Lit = false
				Printf("it is now pitch black.\n")
			}
		}
		return true
	}
	if G.DirObj == &loweredBasket || G.IndirObj == &loweredBasket {
		Printf("The basket is at the other end of the chain.\n")
		return true
	}
	switch G.ActVerb.Norm {
	case "take":
		if G.DirObj == &raisedBasket || G.DirObj == &loweredBasket {
			Printf("The cage is securely fastened to the iron chain.\n")
			return true
		}
	}
	return false
}

func batFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "tell":
		fweep(6)
		G.Params.Continue = NumUndef
		return true
	case "take", "attack", "mung":
		if garlic.Location() == G.Winner || garlic.IsIn(G.Here) {
			Printf("You can't reach him; he's on the ceiling.\n")
			return true
		}
		flyMe()
		return true
	}
	return false
}

func bellFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "ring":
		if G.Here == &enteranceToHades && !gD().LLDFlag {
			return false
		}
		Printf("Ding, dong.\n")
		return true
	}
	return false
}

func hotBellFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "take":
		Printf("The bell is very hot and cannot be taken.\n")
		return true
	case "rub":
		if G.IndirObj != nil && G.IndirObj.Has(FlgBurn) {
			Printf("The %s burns and is consumed.\n", G.IndirObj.Desc)
			removeCarefully(G.IndirObj)
			return true
		}
		if G.IndirObj == &hands {
			Printf("The bell is too hot to touch.\n")
			return true
		}
		Printf("The heat from the bell is too intense.\n")
		return true
	case "ring":
		if G.IndirObj != nil {
			if G.IndirObj.Has(FlgBurn) {
				Printf("The %s burns and is consumed.\n", G.IndirObj.Desc)
				removeCarefully(G.IndirObj)
				return true
			}
			if G.IndirObj == &hands {
				Printf("The bell is too hot to touch.\n")
				return true
			}
			Printf("The heat from the bell is too intense.\n")
			return true
		}
		Printf("The bell is too hot to reach.\n")
		return true
	case "pour on":
		removeCarefully(G.DirObj)
		Printf("The water cools the bell and is evaporated.\n")
		QueueInt("iXbh", false).Run = false
		iXbh()
		return true
	}
	return false
}

func axeFcn(arg ActionArg) bool {
	if gD().TrollFlag {
		return false
	}
	return weaponFunction(&axe, &troll)
}

func trapDoorFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "raise":
		Perform(ActionVerb{Norm: "open", Orig: "open"}, &trapDoor, nil)
		return true
	case "open", "close":
		if G.Here == &livingRoom {
			openClose(G.DirObj,
				"The door reluctantly opens to reveal a rickety staircase descending into darkness.",
				"The door swings shut and closes.")
			return true
		}
	case "look under":
		if G.Here == &livingRoom {
			if trapDoor.Has(FlgOpen) {
				Printf("You see a rickety staircase descending into darkness.\n")
			} else {
				Printf("it's closed.\n")
			}
			return true
		}
	}
	if G.Here == &cellar {
		switch G.ActVerb.Norm {
		case "open":
			if !trapDoor.Has(FlgOpen) {
				Printf("The door is locked from above.\n")
				return true
			}
			Printf("%s\n", PickOne(dummy))
			return true
		case "unlock":
			if !trapDoor.Has(FlgOpen) {
				Printf("The door is locked from above.\n")
				return true
			}
			return false
		case "close":
			if !trapDoor.Has(FlgOpen) {
				trapDoor.Take(FlgTouch)
				trapDoor.Take(FlgOpen)
				Printf("The door closes and locks.\n")
				return true
			}
			Printf("%s\n", PickOne(dummy))
			return true
		}
	}
	return false
}

func frontDoorFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "open":
		Printf("The door cannot be opened.\n")
		return true
	case "burn":
		Printf("You cannot burn this door.\n")
		return true
	case "mung":
		Printf("You can't seem to damage the door.\n")
		return true
	case "look behind":
		Printf("it won't open.\n")
		return true
	}
	return false
}

func barrowDoorFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "open", "close":
		Printf("The door is too heavy.\n")
		return true
	}
	return false
}

func barrowFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "through":
		doWalk(West)
		return true
	}
	return false
}

func bottleFcn(arg ActionArg) bool {
	empty := false
	switch G.ActVerb.Norm {
	case "throw":
		if G.DirObj == &bottle {
			removeCarefully(G.DirObj)
			empty = true
			Printf("The bottle hits the far wall and shatters.\n")
		}
	case "mung":
		empty = true
		removeCarefully(G.DirObj)
		Printf("A brilliant maneuver destroys the bottle.\n")
	case "shake":
		if bottle.Has(FlgOpen) && water.IsIn(&bottle) {
			empty = true
		}
	}
	if empty && water.IsIn(&bottle) {
		Printf("The water spills to the floor and evaporates.\n")
		removeCarefully(&water)
		return true
	}
	if empty {
		return true
	}
	return false
}

func crackFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "through":
		Printf("You can't fit through the crack.\n")
		return true
	}
	return false
}

func grateFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "lock":
		if G.Here == &gratingRoom {
			gD().GrUnlock = false
			Printf("The grate is locked.\n")
			return true
		}
		if G.Here == &clearing {
			Printf("You can't lock it from this side.\n")
			return true
		}
	case "unlock":
		if G.DirObj != &grate {
			break
		}
		if G.Here == &gratingRoom && G.IndirObj == &keys {
			gD().GrUnlock = true
			Printf("The grate is unlocked.\n")
			return true
		}
		if G.Here == &clearing && G.IndirObj == &keys {
			Printf("You can't reach the lock from here.\n")
			return true
		}
		Printf("Can you unlock a grating with a %s?\n", G.IndirObj.Desc)
		return true
	case "pick":
		Printf("You can't pick the lock.\n")
		return true
	case "open", "close":
		if G.IndirObj == &keys {
			Perform(ActionVerb{Norm: "unlock", Orig: "unlock"}, &grate, &keys)
			return true
		}
		if gD().GrUnlock {
			var openStr string
			if G.Here == &clearing {
				openStr = "The grating opens."
			} else {
				openStr = "The grating opens to reveal trees above you."
			}
			openClose(&grate, openStr, "The grating is closed.")
			if grate.Has(FlgOpen) {
				if G.Here != &clearing && !gD().GrateRevealed {
					Printf("A pile of leaves falls onto your head and to the ground.\n")
					gD().GrateRevealed = true
					leaves.MoveTo(G.Here)
				}
				gratingRoom.Give(FlgOn)
			} else {
				gratingRoom.Take(FlgOn)
			}
		} else {
			Printf("The grating is locked.\n")
		}
		return true
	case "put":
		if G.IndirObj == &grate {
			if G.DirObj.GetSize() > 20 {
				Printf("it won't fit through the grating.\n")
			} else {
				G.DirObj.MoveTo(&gratingRoom)
				Printf("The %s goes through the grating into the darkness below.\n", G.DirObj.Desc)
			}
			return true
		}
	}
	return false
}

func knifeFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "take":
		atticTable.Take(FlgNoDesc)
		return false
	}
	return false
}

func skeletonFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "take", "rub", "move", "push", "raise", "lower", "attack", "kick", "kiss":
		Printf("A ghost appears in the room and is appalled at your desecration of the remains of a fellow adventurer. He casts a curse on your valuables and banishes them to the Land of the Living Dead. The ghost leaves, muttering obscenities.\n")
		rob(G.Here, &landOfLivingDead, 100)
		rob(&adventurer, &landOfLivingDead, 0)
		return true
	}
	return false
}

func torchFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "examine":
		Printf("The torch is burning.\n")
		return true
	case "pour on":
		if G.IndirObj == &torch {
			Printf("The water evaporates before it gets close.\n")
			return true
		}
	case "lamp off":
		if G.DirObj.Has(FlgOn) {
			Printf("You nearly burn your hand trying to extinguish the flame.\n")
			return true
		}
	}
	return false
}

func rustyKnifeFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "take":
		if sword.IsIn(G.Winner) {
			Printf("As you touch the rusty knife, your sword gives a single pulse of blinding blue light.\n")
		}
		return false
	case "attack":
		if G.IndirObj == &rustyKnife {
			removeCarefully(&rustyKnife)
			jigsUp("As the knife approaches its victim, your mind is submerged by an overmastering will. Slowly, your hand turns, until the rusty blade is an inch from your neck. The knife seems to sing as it savagely slits your throat.", false)
			return true
		}
	case "swing":
		if G.DirObj == &rustyKnife && G.IndirObj != nil {
			removeCarefully(&rustyKnife)
			jigsUp("As the knife approaches its victim, your mind is submerged by an overmastering will. Slowly, your hand turns, until the rusty blade is an inch from your neck. The knife seems to sing as it savagely slits your throat.", false)
			return true
		}
	}
	return false
}

func leafPileFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "count":
		Printf("There are 69,105 leaves here.\n")
		return true
	case "burn":
		leavesAppear()
		removeCarefully(G.DirObj)
		if G.DirObj.IsIn(G.Here) {
			Printf("The leaves burn.\n")
		} else {
			jigsUp("The leaves burn, and so do you.", false)
		}
		return true
	case "cut":
		Printf("You rustle the leaves around, making quite a mess.\n")
		leavesAppear()
		return true
	case "move", "take":
		if G.ActVerb.Norm == "move" {
			Printf("Done.\n")
		}
		if gD().GrateRevealed {
			return false
		}
		leavesAppear()
		return G.ActVerb.Norm != "take"
	case "look under":
		if !gD().GrateRevealed {
			Printf("Underneath the pile of leaves is a grating. As you release the leaves, the grating is once again concealed from view.\n")
			return true
		}
	}
	return false
}

func matchFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "lamp on", "burn":
		if G.DirObj == &match {
			if gD().MatchCount > 0 {
				gD().MatchCount--
			}
			if gD().MatchCount <= 0 {
				Printf("I'm afraid that you have run out of matches.\n")
				return true
			}
			if G.Here == &lowerShaft || G.Here == &timberRoom {
				Printf("This room is drafty, and the match goes out instantly.\n")
				return true
			}
			match.Give(FlgFlame)
			match.Give(FlgOn)
			Queue("iMatch", 2).Run = true
			Printf("One of the matches starts to burn.\n")
			if !G.Lit {
				G.Lit = true
				vLook(ActUnk)
			}
			return true
		}
	case "lamp off":
		if match.Has(FlgFlame) {
			Printf("The match is out.\n")
			match.Take(FlgFlame)
			match.Take(FlgOn)
			G.Lit = IsLit(G.Here, true)
			if !G.Lit {
				Printf("it's pitch black in here!\n")
			}
			QueueInt("iMatch", false).Run = false
			return true
		}
	case "count", "open":
		Printf("You have ")
		cnt := gD().MatchCount - 1
		if cnt <= 0 {
			Printf("no")
		} else {
			Printf("%d", cnt)
		}
		Printf(" match")
		if cnt != 1 {
			Printf("es.")
		} else {
			Printf(".")
		}
		Printf("\n")
		return true
	case "examine":
		if match.Has(FlgOn) {
			Printf("The match is burning.\n")
		} else {
			Printf("The matchbook isn't very interesting, except for what's written on it.\n")
		}
		return true
	}
	return false
}

func mirrorMirrorFcn(arg ActionArg) bool {
	rm2 := &mirrorRoom2
	switch G.ActVerb.Norm {
	case "rub":
		if !gD().MirrorMung {
			if G.IndirObj != nil && G.IndirObj != &hands {
				Printf("You feel a faint tingling transmitted through the %s.\n", G.IndirObj.Desc)
				return true
			}
			if G.Here == rm2 {
				rm2 = &mirrorRoom1
			}
			// Swap room contents
			l1 := append([]*Object{}, G.Here.Children...)
			l2 := append([]*Object{}, rm2.Children...)
			for _, c := range l1 {
				c.MoveTo(rm2)
			}
			for _, c := range l2 {
				c.MoveTo(G.Here)
			}
			moveToRoom(rm2, false)
			Printf("There is a rumble from deep within the earth and the room shakes.\n")
			return true
		}
		return false
	case "look inside", "examine":
		if gD().MirrorMung {
			Printf("The mirror is broken into many pieces.\n")
		} else {
			Printf("There is an ugly person staring back at you.\n")
		}
		return true
	case "take":
		Printf("The mirror is many times your size. Give up.\n")
		return true
	case "mung", "throw", "attack":
		if gD().MirrorMung {
			Printf("Haven't you done enough damage already?\n")
		} else {
			gD().MirrorMung = true
			G.Lucky = false
			Printf("You have broken the mirror. I hope you have a seven years' supply of good luck handy.\n")
		}
		return true
	}
	return false
}

func paintingFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "mung":
		G.DirObj.SetTValue(0)
		G.DirObj.LongDesc = "There is a worthless piece of canvas here."
		Printf("Congratulations! Unlike the other vandals, who merely stole the artist's masterpieces, you have destroyed one.\n")
		return true
	}
	return false
}

func candlesFcn(arg ActionArg) bool {
	if !candles.Has(FlgTouch) {
		Queue("iCandles", -1).Run = true
	}
	if G.IndirObj == &candles {
		return false
	}
	switch G.ActVerb.Norm {
	case "lamp on", "burn":
		if candles.Has(FlgDestroyed) {
			Printf("Alas, there's not much left of the candles. Certainly not enough to burn.\n")
			return true
		}
		if G.IndirObj == nil {
			if match.Has(FlgFlame) {
				Printf("(with the match)\n")
				Perform(ActionVerb{Norm: "lamp on", Orig: "light"}, &candles, &match)
				return true
			}
			Printf("You should say what to light them with.\n")
			return true
		}
		if G.IndirObj == &match && match.Has(FlgOn) {
			Printf("The candles are ")
			if candles.Has(FlgOn) {
				Printf("already lit.\n")
			} else {
				candles.Give(FlgOn)
				Printf("lit.\n")
				Queue("iCandles", -1).Run = true
			}
			return true
		}
		if G.IndirObj == &torch {
			if candles.Has(FlgOn) {
				Printf("You realize, just in time, that the candles are already lighted.\n")
			} else {
				Printf("The heat from the torch is so intense that the candles are vaporized.\n")
				removeCarefully(&candles)
			}
			return true
		}
		Printf("You have to light them with something that's burning, you know.\n")
		return true
	case "count":
		Printf("Let's see, how many objects in a pair? Don't tell me, I'll get it.\n")
		return true
	case "lamp off":
		QueueInt("iCandles", false).Run = false
		if candles.Has(FlgOn) {
			Printf("The flame is extinguished.")
			candles.Take(FlgOn)
			candles.Give(FlgTouch)
			G.Lit = IsLit(G.Here, true)
			if !G.Lit {
				Printf(" it's really dark in here....")
			}
			Printf("\n")
			return true
		}
		Printf("The candles are not lighted.\n")
		return true
	case "put":
		if G.IndirObj != nil && G.IndirObj.Has(FlgBurn) {
			Printf("That wouldn't be smart.\n")
			return true
		}
	case "examine":
		Printf("The candles are ")
		if candles.Has(FlgOn) {
			Printf("burning.\n")
		} else {
			Printf("out.\n")
		}
		return true
	}
	return false
}

func gunkFcn(arg ActionArg) bool {
	removeCarefully(&gunk)
	Printf("The slag was rather insubstantial, and crumbles into dust at your touch.\n")
	return true
}

func bodyFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "take":
		Printf("A force keeps you from taking the bodies.\n")
		return true
	case "mung", "burn":
		jigsUp("The voice of the guardian of the dungeon booms out from the darkness, \"Your disrespect costs you your life!\" and places your head on a sharp pole.", false)
		return true
	}
	return false
}

func blackBookFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "open":
		Printf("The book is already open to page 569.\n")
		return true
	case "close":
		Printf("As hard as you try, the book cannot be closed.\n")
		return true
	case "turn":
		Printf("Beside page 569, there is only one other page with any legible printing on it. Most of it is unreadable, but the subject seems to be the banishment of evil. Apparently, certain noises, lights, and prayers are efficacious in this regard.\n")
		return true
	case "burn":
		removeCarefully(G.DirObj)
		jigsUp("A booming voice says \"Wrong, cretin!\" and you notice that you have turned into a pile of dust. How, I can't imagine.", false)
		return true
	}
	return false
}

func sceptreFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "wave", "raise":
		if G.Here == &aragainFalls || G.Here == &endOfRainbow {
			if !gD().RainbowFlag {
				potOfGold.Take(FlgInvis)
				Printf("Suddenly, the rainbow appears to become solid and, I venture, walkable (I think the giveaway was the stairs and bannister).\n")
				if G.Here == &endOfRainbow && potOfGold.IsIn(&endOfRainbow) {
					Printf("A shimmering pot of gold appears at the end of the rainbow.\n")
				}
				gD().RainbowFlag = true
			} else {
				rob(&onRainbow, &wall, 0)
				Printf("The rainbow seems to have become somewhat run-of-the-mill.\n")
				gD().RainbowFlag = false
				return true
			}
			return true
		}
		if G.Here == &onRainbow {
			gD().RainbowFlag = false
			jigsUp("The structural integrity of the rainbow is severely compromised, leaving you hanging in midair, supported only by water vapor. Bye.", false)
			return true
		}
		Printf("A dazzling display of color briefly emanates from the sceptre.\n")
		return true
	}
	return false
}

func slideFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "through", "climb up", "climb down", "climb":
		if G.Here == &cellar {
			doWalk(West)
			return true
		}
		Printf("You tumble down the slide....\n")
		moveToRoom(&cellar, true)
		return true
	case "put":
		if G.DirObj == &me {
			if G.Here == &cellar {
				doWalk(West)
				return true
			}
			Printf("You tumble down the slide....\n")
			moveToRoom(&cellar, true)
			return true
		}
		slider(G.DirObj)
		return true
	}
	return false
}

func sandwichBagFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "smell":
		if lunch.IsIn(G.DirObj) {
			Printf("it smells of hot peppers.\n")
			return true
		}
	}
	return false
}

func toolChestFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "examine":
		Printf("The chests are all empty.\n")
		return true
	case "take", "open", "put":
		removeCarefully(&toolChest)
		Printf("The chests are so rusty and corroded that they crumble when you touch them.\n")
		return true
	}
	return false
}

func buttonFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "read":
		Printf("They're greek to you.\n")
		return true
	case "push":
		if G.DirObj == &blueButton {
			if gD().WaterLevel == 0 {
				leak.Take(FlgInvis)
				Printf("There is a rumbling sound and a stream of water appears to burst from the east wall of the room (apparently, a leak has occurred in a pipe).\n")
				gD().WaterLevel = 1
				Queue("iMaintRoom", -1).Run = true
				return true
			}
			Printf("The blue button appears to be jammed.\n")
			return true
		}
		if G.DirObj == &redButton {
			Printf("The lights within the room ")
			if G.Here.Has(FlgOn) {
				G.Here.Take(FlgOn)
				Printf("shut off.\n")
			} else {
				G.Here.Give(FlgOn)
				Printf("come on.\n")
			}
			return true
		}
		if G.DirObj == &brownButton {
			damRoom.Take(FlgTouch)
			gD().GateFlag = false
			Printf("Click.\n")
			return true
		}
		if G.DirObj == &yellowButton {
			damRoom.Take(FlgTouch)
			gD().GateFlag = true
			Printf("Click.\n")
			return true
		}
		return true
	}
	return false
}

func leakFcn(arg ActionArg) bool {
	if gD().WaterLevel > 0 {
		switch G.ActVerb.Norm {
		case "put", "put on":
			if G.DirObj == &putty {
				fixMaintLeak()
				return true
			}
		case "plug":
			if G.IndirObj == &putty {
				fixMaintLeak()
				return true
			}
			withTell(G.IndirObj)
			return true
		}
	}
	return false
}

func machineFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "take":
		if G.DirObj == &machine {
			Printf("it is far too large to carry.\n")
			return true
		}
	case "open":
		if machine.Has(FlgOpen) {
			Printf("%s\n", PickOne(dummy))
		} else if machine.HasChildren() {
			Printf("The lid opens, revealing ")
			printContents(&machine)
			Printf(".\n")
			machine.Give(FlgOpen)
		} else {
			Printf("The lid opens.\n")
			machine.Give(FlgOpen)
		}
		return true
	case "close":
		if machine.Has(FlgOpen) {
			Printf("The lid closes.\n")
			machine.Take(FlgOpen)
		} else {
			Printf("%s\n", PickOne(dummy))
		}
		return true
	case "lamp on":
		if G.IndirObj == nil {
			Printf("it's not clear how to turn it on with your bare hands.\n")
		} else {
			Perform(ActionVerb{Norm: "turn", Orig: "turn"}, &machineSwitch, G.IndirObj)
			return true
		}
		return true
	}
	return false
}

func machineSwitchFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "turn":
		if G.IndirObj == &screwdriver {
			if machine.Has(FlgOpen) {
				Printf("The machine doesn't seem to want to do anything.\n")
			} else {
				Printf("The machine comes to life (figuratively) with a dazzling display of colored lights and bizarre noises. After a few moments, the excitement abates.\n")
				if coal.IsIn(&machine) {
					removeCarefully(&coal)
					diamond.MoveTo(&machine)
				} else {
					// Remove everything and put gunk in
					toRemove := append([]*Object{}, machine.Children...)
					for _, o := range toRemove {
						removeCarefully(o)
					}
					gunk.MoveTo(&machine)
				}
			}
		} else {
			Printf("it seems that a %s won't do.\n", G.IndirObj.Desc)
		}
		return true
	}
	return false
}

func puttyFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "oil":
		if G.IndirObj == &putty {
			Printf("The all-purpose gunk isn't a lubricant.\n")
			return true
		}
	case "put":
		if G.DirObj == &putty {
			Printf("The all-purpose gunk isn't a lubricant.\n")
			return true
		}
	}
	return false
}

func tubeFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "put":
		if G.IndirObj == &tube {
			Printf("The tube refuses to accept anything.\n")
			return true
		}
	case "squeeze":
		if G.DirObj.Has(FlgOpen) && putty.IsIn(G.DirObj) {
			putty.MoveTo(G.Winner)
			Printf("The viscous material oozes into your hand.\n")
			return true
		}
		if G.DirObj.Has(FlgOpen) {
			Printf("The tube is apparently empty.\n")
			return true
		}
		Printf("The tube is closed.\n")
		return true
	}
	return false
}

func swordFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "take":
		if G.Winner == &adventurer {
			Queue("iSword", -1).Run = true
			return false
		}
	case "examine":
		g := sword.GetTValue()
		if g == 1 {
			Printf("Your sword is glowing with a faint blue glow.\n")
			return true
		}
		if g == 2 {
			Printf("Your sword is glowing very brightly.\n")
			return true
		}
	}
	return false
}

func lanternFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "throw":
		Printf("The lamp has smashed into the floor, and the light has gone out.\n")
		QueueInt("iLantern", false).Run = false
		removeCarefully(&lamp)
		brokenLamp.MoveTo(G.Here)
		return true
	case "lamp on":
		if lamp.Has(FlgDestroyed) {
			Printf("A burned-out lamp won't light.\n")
			return true
		}
		itm := QueueInt("iLantern", false)
		if itm.Tick <= 0 {
			// First activation or timer expired: initialize countdown
			itm.Tick = -1
		}
		// Otherwise resume from where we left off
		itm.Run = true
		return false
	case "lamp off":
		if lamp.Has(FlgDestroyed) {
			Printf("The lamp has already burned out.\n")
			return true
		}
		QueueInt("iLantern", false).Run = false
		return false
	case "examine":
		Printf("The lamp ")
		if lamp.Has(FlgDestroyed) {
			Printf("has burned out.\n")
		} else if lamp.Has(FlgOn) {
			Printf("is on.\n")
		} else {
			Printf("is turned off.\n")
		}
		return true
	}
	return false
}

func mailboxFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "take":
		if G.DirObj == &mailbox {
			Printf("it is securely anchored.\n")
			return true
		}
	}
	return false
}

// ================================================================
// TROLL
// ================================================================

func trollFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "tell":
		G.Params.Continue = NumUndef
		Printf("The troll isn't much of a conversationalist.\n")
		return true
	}
	if arg == ActBusy {
		if axe.IsIn(&troll) {
			return false
		}
		if axe.IsIn(G.Here) && Prob(75, true) {
			axe.Give(FlgNoDesc)
			axe.Take(FlgWeapon)
			axe.MoveTo(&troll)
			troll.LongDesc = "A nasty-looking troll, brandishing a bloody axe, blocks all passages out of the room."
			if troll.IsIn(G.Here) {
				Printf("The troll, angered and humiliated, recovers his weapon. He appears to have an axe to grind with you.\n")
			}
			return true
		}
		if troll.IsIn(G.Here) {
			troll.LongDesc = "A pathetically babbling troll is here."
			Printf("The troll, disarmed, cowers in terror, pleading for his life in the guttural tongue of the trolls.\n")
			return true
		}
		return false
	}
	if arg == ActDead {
		if axe.IsIn(&troll) {
			axe.MoveTo(G.Here)
			axe.Take(FlgNoDesc)
			axe.Give(FlgWeapon)
		}
		gD().TrollFlag = true
		return true
	}
	if arg == ActUnconscious {
		troll.Take(FlgFight)
		if axe.IsIn(&troll) {
			axe.MoveTo(G.Here)
			axe.Take(FlgNoDesc)
			axe.Give(FlgWeapon)
		}
		troll.LongDesc = "An unconscious troll is sprawled on the floor. All passages out of the room are open."
		gD().TrollFlag = true
		return true
	}
	if arg == ActConscious {
		if troll.IsIn(G.Here) {
			troll.Give(FlgFight)
			Printf("The troll stirs, quickly resuming a fighting stance.\n")
		}
		if axe.IsIn(&troll) {
			troll.LongDesc = "A nasty-looking troll, brandishing a bloody axe, blocks all passages out of the room."
		} else if axe.IsIn(&trollRoom) {
			axe.Give(FlgNoDesc)
			axe.Take(FlgWeapon)
			axe.MoveTo(&troll)
			troll.LongDesc = "A nasty-looking troll, brandishing a bloody axe, blocks all passages out of the room."
		} else {
			troll.LongDesc = "A troll is here."
		}
		gD().TrollFlag = false
		return true
	}
	if arg == ActFirst {
		if Prob(33, false) {
			troll.Give(FlgFight)
			G.Params.Continue = NumUndef
			return true
		}
		return false
	}
	// Default (no mode - regular verbs)
	switch G.ActVerb.Norm {
	case "examine":
		Printf("%s\n", troll.LongDesc)
		return true
	case "throw", "give":
		if G.DirObj != nil && G.IndirObj == &troll {
			awaken(&troll)
			if G.DirObj == &axe && axe.IsIn(G.Winner) {
				Printf("The troll scratches his head in confusion, then takes the axe.\n")
				troll.Give(FlgFight)
				axe.MoveTo(&troll)
				return true
			}
			if G.DirObj == &troll || G.DirObj == &axe {
				Printf("You would have to get the %s first, and that seems unlikely.\n", G.DirObj.Desc)
				return true
			}
			if G.ActVerb.Norm == "throw" {
				Printf("The troll, who is remarkably coordinated, catches the %s", G.DirObj.Desc)
			} else {
				Printf("The troll, who is not overly proud, graciously accepts the gift")
			}
			if Prob(20, false) && (G.DirObj == &knife || G.DirObj == &sword || G.DirObj == &axe) {
				removeCarefully(G.DirObj)
				Printf(" and eats it hungrily. Poor troll, he dies from an internal hemorrhage and his carcass disappears in a sinister black fog.\n")
				removeCarefully(&troll)
				trollFcn(ActDead)
				gD().TrollFlag = true
			} else if G.DirObj == &knife || G.DirObj == &sword || G.DirObj == &axe {
				G.DirObj.MoveTo(G.Here)
				Printf(" and, being for the moment sated, throws it back. Fortunately, the troll has poor control, and the %s falls to the floor. He does not look pleased.\n", G.DirObj.Desc)
				troll.Give(FlgFight)
			} else {
				Printf(" and not having the most discriminating tastes, gleefully eats it.\n")
				removeCarefully(G.DirObj)
			}
			return true
		}
	case "take", "move":
		awaken(&troll)
		Printf("The troll spits in your face, grunting \"Better luck next time\" in a rather barbarous accent.\n")
		return true
	case "mung":
		awaken(&troll)
		Printf("The troll laughs at your puny gesture.\n")
		return true
	case "listen":
		Printf("Every so often the troll says something, probably uncomplimentary, in his guttural tongue.\n")
		return true
	case "hello":
		if gD().TrollFlag {
			Printf("Unfortunately, the troll can't hear you.\n")
			return true
		}
	}
	return false
}

// ================================================================
// CYCLOPS
// ================================================================

func cyclopsFcn(arg ActionArg) bool {
	count := gD().CycloWrath
	if G.Winner == &cyclops {
		if gD().CyclopsFlag {
			Printf("No use talking to him. He's fast asleep.\n")
			return true
		}
		switch G.ActVerb.Norm {
		case "odysseus":
			G.Winner = &adventurer
			Perform(ActionVerb{Norm: "odysseus", Orig: "odysseus"}, nil, nil)
			return true
		}
		Printf("The cyclops prefers eating to making conversation.\n")
		return true
	}
	if gD().CyclopsFlag {
		switch G.ActVerb.Norm {
		case "examine":
			Printf("The cyclops is sleeping like a baby, albeit a very ugly one.\n")
			return true
		case "alarm", "kick", "attack", "burn", "mung":
			Printf("The cyclops yawns and stares at the thing that woke him up.\n")
			gD().CyclopsFlag = false
			cyclops.Give(FlgFight)
			if count < 0 {
				gD().CycloWrath = -count
			} else {
				gD().CycloWrath = count
			}
			return true
		default:
			return false
		}
	}
	switch G.ActVerb.Norm {
	case "examine":
		Printf("A hungry cyclops is standing at the foot of the stairs.\n")
		return true
	case "give":
		if G.IndirObj == &cyclops {
			if G.DirObj == &lunch {
				if count >= 0 {
					removeCarefully(&lunch)
					Printf("The cyclops says \"Mmm Mmm. I love hot peppers! But oh, could I use a drink. Perhaps I could drink the blood of that thing.\"  From the gleam in his eye, it could be surmised that you are \"that thing\".\n")
					gD().CycloWrath = min(-1, -count)
				}
				Queue("iCyclops", -1).Run = true
				return true
			}
			if G.DirObj == &water || (G.DirObj == &bottle && water.IsIn(&bottle)) {
				if count < 0 {
					removeCarefully(&water)
					bottle.MoveTo(G.Here)
					bottle.Give(FlgOpen)
					cyclops.Take(FlgFight)
					Printf("The cyclops takes the bottle, checks that it's open, and drinks the water. A moment later, he lets out a yawn that nearly blows you over, and then falls fast asleep (what did you put in that drink, anyway?).\n")
					gD().CyclopsFlag = true
				} else {
					Printf("The cyclops apparently is not thirsty and refuses your generous offer.\n")
				}
				return true
			}
			if G.DirObj == &garlic {
				Printf("The cyclops may be hungry, but there is a limit.\n")
				return true
			}
			Printf("The cyclops is not so stupid as to eat THAT!\n")
			return true
		}
	case "throw", "attack", "mung":
		Queue("iCyclops", -1).Run = true
		if G.ActVerb.Norm == "mung" {
			Printf("\"Do you think I'm as stupid as my father was?\", he says, dodging.\n")
		} else {
			Printf("The cyclops shrugs but otherwise ignores your pitiful attempt.\n")
			if G.ActVerb.Norm == "throw" {
				G.DirObj.MoveTo(G.Here)
			}
			return true
		}
		return true
	case "take":
		Printf("The cyclops doesn't take kindly to being grabbed.\n")
		return true
	case "tie":
		Printf("You cannot tie the cyclops, though he is fit to be tied.\n")
		return true
	case "listen":
		Printf("You can hear his stomach rumbling.\n")
		return true
	}
	return false
}

// ================================================================
// THIEF / ROBBER
// ================================================================

func dumbContainerFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "open", "close", "look inside":
		Printf("You can't do that.\n")
		return true
	case "examine":
		Printf("it looks pretty much like a %s.\n", G.DirObj.Desc)
		return true
	}
	return false
}

func chaliceFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "take":
		if G.DirObj.IsIn(&treasureRoom) && thief.IsIn(&treasureRoom) && thief.Has(FlgFight) && !thief.Has(FlgInvis) && thief.LongDesc != robberUDesc {
			Printf("You'd be stabbed in the back first.\n")
			return true
		}
		return false
	case "put":
		if G.IndirObj == &chalice {
			Printf("You can't. it's not a very good chalice, is it?\n")
			return true
		}
	}
	return dumbContainerFcn(arg)
}

func trunkFcn(arg ActionArg) bool {
	return stupidContainer(&trunk, "jewels")
}

func bagOfCoinsFcn(arg ActionArg) bool {
	return stupidContainer(&bagOfCoins, "coins")
}

func stupidContainer(obj *Object, str string) bool {
	switch G.ActVerb.Norm {
	case "open", "close":
		Printf("The %s are safely inside; there's no need to do that.\n", str)
		return true
	case "look inside", "examine":
		Printf("There are lots of %s in there.\n", str)
		return true
	case "put":
		if G.IndirObj == obj {
			Printf("Don't be silly. it wouldn't be a %s anymore.\n", obj.Desc)
			return true
		}
	}
	return false
}

func garlicFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "eat":
		removeCarefully(G.DirObj)
		Printf("What the heck! You won't make friends this way, but nobody around here is too friendly anyhow. Gulp!\n")
		return true
	}
	return false
}

func batDescFcn(arg ActionArg) bool {
	if garlic.Location() == G.Winner || garlic.IsIn(G.Here) {
		Printf("In the corner of the room on the ceiling is a large vampire bat who is obviously deranged and holding his nose.\n")
	} else {
		Printf("A large vampire bat, hanging from the ceiling, swoops down at you!\n")
	}
	return true
}

func trophyCaseFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "take":
		if G.DirObj == &trophyCase {
			Printf("The trophy case is securely fastened to the wall.\n")
			return true
		}
	}
	return false
}

func boardedWindowFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "open":
		Printf("The windows are boarded and can't be opened.\n")
		return true
	case "mung":
		Printf("You can't break the windows open.\n")
		return true
	}
	return false
}

func nailsPseudo(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "take":
		Printf("The nails, deeply imbedded in the door, cannot be removed.\n")
		return true
	}
	return false
}

func cliffObjectFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "leap":
		Printf("That would be very unwise. Perhaps even fatal.\n")
		return true
	case "put":
		if G.DirObj == &me {
			Printf("That would be very unwise. Perhaps even fatal.\n")
			return true
		}
		if G.IndirObj == &climbableCliff {
			Printf("The %s tumbles into the river and is seen no more.\n", G.DirObj.Desc)
			removeCarefully(G.DirObj)
			return true
		}
	case "throw off":
		if G.IndirObj == &climbableCliff {
			Printf("The %s tumbles into the river and is seen no more.\n", G.DirObj.Desc)
			removeCarefully(G.DirObj)
			return true
		}
	}
	return false
}

func whiteCliffFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "climb up", "climb down", "climb":
		Printf("The cliff is too steep for climbing.\n")
		return true
	}
	return false
}

func rainbowFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "cross", "through":
		if G.Here == &canyonView {
			Printf("From here?!?\n")
			return true
		}
		if gD().RainbowFlag {
			if G.Here == &aragainFalls {
				moveToRoom(&endOfRainbow, true)
			} else if G.Here == &endOfRainbow {
				moveToRoom(&aragainFalls, true)
			} else {
				Printf("You'll have to say which way...\n")
			}
		} else {
			Printf("Can you walk on water vapor?\n")
		}
		return true
	case "look under":
		Printf("The Frigid river flows under the rainbow.\n")
		return true
	}
	return false
}

func ropeFcn(arg ActionArg) bool {
	if G.Here != &domeRoom {
		gD().DomeFlag = false
		switch G.ActVerb.Norm {
		case "tie":
			Printf("You can't tie the rope to that.\n")
			return true
		default:
			return false
		}
	}
	switch G.ActVerb.Norm {
	case "tie":
		if G.IndirObj == &railing {
			if gD().DomeFlag {
				Printf("The rope is already tied to it.\n")
			} else {
				Printf("The rope drops over the side and comes within ten feet of the floor.\n")
				gD().DomeFlag = true
				rope.Give(FlgNoDesc)
				rloc := rope.Location()
				if rloc == nil || !rloc.IsIn(&rooms) {
					rope.MoveTo(G.Here)
				}
			}
			return true
		}
		return false
	case "climb down":
		if (G.DirObj == &rope || G.DirObj == &rooms) && gD().DomeFlag {
			doWalk(Down)
			return true
		}
	case "tie up":
		if G.IndirObj == &rope {
			if G.DirObj.Has(FlgActor) {
				if G.DirObj.GetStrength() < 0 {
					Printf("Your attempt to tie up the %s awakens him.", G.DirObj.Desc)
					awaken(G.DirObj)
				} else {
					Printf("The %s struggles and you cannot tie him up.\n", G.DirObj.Desc)
				}
			} else {
				Printf("Why would you tie up a %s?\n", G.DirObj.Desc)
			}
			return true
		}
	case "untie":
		if gD().DomeFlag {
			gD().DomeFlag = false
			rope.Take(FlgNoDesc)
			Printf("The rope is now untied.\n")
		} else {
			Printf("it is not tied to anything.\n")
		}
		return true
	case "drop":
		if !gD().DomeFlag {
			rope.MoveTo(&torchRoom)
			Printf("The rope drops gently to the floor below.\n")
			return true
		}
	case "take":
		if gD().DomeFlag {
			Printf("The rope is tied to the railing.\n")
			return true
		}
	}
	return false
}

func eggObjectFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "open", "mung":
		if G.DirObj == &egg {
			if G.DirObj.Has(FlgOpen) {
				Printf("The egg is already open.\n")
				return true
			}
			if G.IndirObj == nil {
				Printf("You have neither the tools nor the expertise.\n")
				return true
			}
			if G.IndirObj == &hands {
				Printf("I doubt you could do that without damaging it.\n")
				return true
			}
			if G.IndirObj.Has(FlgWeapon) || G.IndirObj.Has(FlgTool) || G.ActVerb.Norm == "mung" {
				Printf("The egg is now open, but the clumsiness of your attempt has seriously compromised its esthetic appeal.")
				badEgg()
				Printf("\n")
				return true
			}
			if G.DirObj.Has(FlgFight) {
				Printf("Not to say that using the %s isn't original too...\n", G.IndirObj.Desc)
				return true
			}
			Printf("The concept of using a %s is certainly original.\n", G.IndirObj.Desc)
			G.DirObj.Give(FlgFight)
			return true
		}
		Printf("Your rather indelicate handling of the egg has caused it some damage, although you have succeeded in opening it.")
		badEgg()
		Printf("\n")
		return true
	case "climb on", "hatch":
		Printf("There is a noticeable crunch from beneath you, and inspection reveals that the egg is lying open, badly damaged.")
		badEgg()
		Printf("\n")
		return true
	case "throw":
		G.DirObj.MoveTo(G.Here)
		Printf("Your rather indelicate handling of the egg has caused it some damage, although you have succeeded in opening it.")
		badEgg()
		Printf("\n")
		return true
	}
	return false
}

func canaryObjectFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "wind":
		if G.DirObj == &canary {
			if !gD().SingSong && forestRoomQ() {
				Printf("The canary chirps, slightly off-key, an aria from a forgotten opera. From out of the greenery flies a lovely songbird. it perches on a limb just over your head and opens its beak to sing. As it does so a beautiful brass bauble drops from its mouth, bounces off the top of your head, and lands glimmering in the grass. As the canary winds down, the songbird flies away.\n")
				gD().SingSong = true
				dest := G.Here
				if G.Here == &upATree {
					dest = &path
				}
				bauble.MoveTo(dest)
			} else {
				Printf("The canary chirps blithely, if somewhat tinnily, for a short time.\n")
			}
		} else {
			Printf("There is an unpleasant grinding noise from inside the canary.\n")
		}
		return true
	}
	return false
}

func rugFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "raise":
		Printf("The rug is too heavy to lift")
		if gD().RugMoved {
			Printf(".\n")
		} else {
			Printf(", but in trying to take it you have noticed an irregularity beneath it.\n")
		}
		return true
	case "move", "push":
		if gD().RugMoved {
			Printf("Having moved the carpet previously, you find it impossible to move it again.\n")
		} else {
			Printf("With a great effort, the rug is moved to one side of the room, revealing the dusty cover of a closed trap door.\n")
			trapDoor.Take(FlgInvis)
			thisIsIt(&trapDoor)
			gD().RugMoved = true
		}
		return true
	case "take":
		Printf("The rug is extremely heavy and cannot be carried.\n")
		return true
	case "look under":
		if !gD().RugMoved && !trapDoor.Has(FlgOpen) {
			Printf("Underneath the rug is a closed trap door. As you drop the corner of the rug, the trap door is once again concealed from view.\n")
			return true
		}
	case "climb on":
		if !gD().RugMoved && !trapDoor.Has(FlgOpen) {
			Printf("As you sit, you notice an irregularity underneath it. Rather than be uncomfortable, you stand up again.\n")
		} else {
			Printf("I suppose you think it's a magic carpet?\n")
		}
		return true
	}
	return false
}

func sandFunction(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "dig":
		if G.IndirObj == &shovel {
			gD().BeachDig++
			if gD().BeachDig > 3 {
				gD().BeachDig = -1
				if scarab.IsIn(G.Here) {
					scarab.Give(FlgInvis)
				}
				jigsUp("The hole collapses, smothering you.", false)
				return true
			}
			if gD().BeachDig == 3 {
				if scarab.Has(FlgInvis) {
					Printf("You can see a scarab here in the sand.\n")
					thisIsIt(&scarab)
					scarab.Take(FlgInvis)
				}
			} else {
				Printf("%s\n", bDigs[gD().BeachDig])
			}
			return true
		}
	}
	return false
}

// ================================================================
// ROOM ACTION FUNCTIONS
// ================================================================

func kitchenFcn(arg ActionArg) bool {
	if arg == ActLook {
		Printf("You are in the kitchen of the white house. A table seems to have been used recently for the preparation of food. A passage leads to the west and a dark staircase can be seen leading upward. A dark chimney leads down and to the east is a small window which is ")
		if kitchenWindow.Has(FlgOpen) {
			Printf("open.\n")
		} else {
			Printf("slightly ajar.\n")
		}
		return true
	}
	if arg == ActBegin {
		switch G.ActVerb.Norm {
		case "climb up":
			if G.DirObj == &stairs {
				doWalk(Up)
				return true
			}
		}
	}
	return false
}

func livingRoomFcn(arg ActionArg) bool {
	if arg == ActLook {
		Printf("You are in the living room. There is a doorway to the east")
		if gD().MagicFlag {
			Printf(". To the west is a cyclops-shaped opening in an old wooden door, above which is some strange gothic lettering, ")
		} else {
			Printf(", a wooden door with strange gothic lettering to the west, which appears to be nailed shut, ")
		}
		Printf("a trophy case, ")
		if gD().RugMoved && trapDoor.Has(FlgOpen) {
			Printf("and a rug lying beside an open trap door.")
		} else if gD().RugMoved {
			Printf("and a closed trap door at your feet.")
		} else if trapDoor.Has(FlgOpen) {
			Printf("and an open trap door at your feet.")
		} else {
			Printf("and a large oriental rug in the center of the room.")
		}
		Printf("\n")
		return true
	}
	if arg == ActEnd {
		switch G.ActVerb.Norm {
		case "take":
			if G.DirObj.IsIn(&trophyCase) {
				touchAll(G.DirObj)
			}
			G.Score = G.BaseScore + otvalFrob(&trophyCase)
			scoreUpd(0)
			return false
		case "put":
			if G.IndirObj == &trophyCase {
				if G.DirObj.IsIn(&trophyCase) {
					touchAll(G.DirObj)
				}
				G.Score = G.BaseScore + otvalFrob(&trophyCase)
				scoreUpd(0)
				return false
			}
		}
	}
	return false
}

func cellarFcn(arg ActionArg) bool {
	if arg == ActLook {
		Printf("You are in a dark and damp cellar with a narrow passageway leading north, and a crawlway to the south. On the west is the bottom of a steep metal ramp which is unclimbable.\n")
		return true
	}
	if arg == ActEnter {
		if trapDoor.Has(FlgOpen) && !trapDoor.Has(FlgTouch) {
			trapDoor.Take(FlgOpen)
			trapDoor.Give(FlgTouch)
			Printf("The trap door crashes shut, and you hear someone barring it.\n\n")
		}
		return false
	}
	return false
}

func stoneBarrowFcn(arg ActionArg) bool {
	if arg == ActBegin {
		switch G.ActVerb.Norm {
		case "enter":
			Printf("Inside the barrow\nAs you enter the barrow, the door closes inexorably behind you. Around you it is dark, but ahead is an enormous cavern, brightly lit. through its center runs a wide stream. Spanning the stream is a small wooden footbridge, and beyond a path leads into a dark tunnel. Above the bridge, floating in the air, is a large sign. it reads:  All ye who stand before this bridge have completed a great and perilous adventure which has tested your wit and courage. You have mastered the first part of the ZORK trilogy. Those who pass over this bridge must be prepared to undertake an even greater adventure that will severely test your skill and bravery!\n\nThe ZORK trilogy continues with \"ZORK II: The Wizard of Frobozz\" and is completed in \"ZORK III: The Dungeon Master.\"\n")
			finish()
			return true
		case "walk":
			if G.Params.HasWalkDir && (G.Params.WalkDir == West || G.Params.WalkDir == In) {
				Printf("Inside the barrow\nAs you enter the barrow, the door closes inexorably behind you. Around you it is dark, but ahead is an enormous cavern, brightly lit. through its center runs a wide stream. Spanning the stream is a small wooden footbridge, and beyond a path leads into a dark tunnel. Above the bridge, floating in the air, is a large sign. it reads:  All ye who stand before this bridge have completed a great and perilous adventure which has tested your wit and courage. You have mastered the first part of the ZORK trilogy. Those who pass over this bridge must be prepared to undertake an even greater adventure that will severely test your skill and bravery!\n\nThe ZORK trilogy continues with \"ZORK II: The Wizard of Frobozz\" and is completed in \"ZORK III: The Dungeon Master.\"\n")
				finish()
				return true
			}
		case "through":
			if G.DirObj == &barrow {
				Printf("Inside the barrow\nAs you enter the barrow, the door closes inexorably behind you. Around you it is dark, but ahead is an enormous cavern, brightly lit. through its center runs a wide stream. Spanning the stream is a small wooden footbridge, and beyond a path leads into a dark tunnel. Above the bridge, floating in the air, is a large sign. it reads:  All ye who stand before this bridge have completed a great and perilous adventure which has tested your wit and courage. You have mastered the first part of the ZORK trilogy. Those who pass over this bridge must be prepared to undertake an even greater adventure that will severely test your skill and bravery!\n\nThe ZORK trilogy continues with \"ZORK II: The Wizard of Frobozz\" and is completed in \"ZORK III: The Dungeon Master.\"\n")
				finish()
				return true
			}
		}
	}
	return false
}

func trollRoomFcn(arg ActionArg) bool {
	if arg == ActEnter && troll.IsIn(G.Here) {
		thisIsIt(&troll)
	}
	return false
}

func clearingFcn(arg ActionArg) bool {
	if arg == ActEnter {
		if !gD().GrateRevealed {
			grate.Give(FlgInvis)
		}
		return false
	}
	if arg == ActLook {
		Printf("You are in a clearing, with a forest surrounding you on all sides. A path leads south.")
		if grate.Has(FlgOpen) {
			Printf("\nThere is an open grating, descending into darkness.")
		} else if gD().GrateRevealed {
			Printf("\nThere is a grating securely fastened into the ground.")
		}
		Printf("\n")
		return true
	}
	return false
}

func maze11Fcn(arg ActionArg) bool {
	if arg == ActEnter {
		grate.Take(FlgInvis)
		return false
	}
	if arg == ActLook {
		Printf("You are in a small room near the maze. There are twisty passages in the immediate vicinity.\n")
		if grate.Has(FlgOpen) {
			Printf("Above you is an open grating with sunlight pouring in.")
		} else if gD().GrUnlock {
			Printf("Above you is a grating.")
		} else {
			Printf("Above you is a grating locked with a skull-and-crossbones lock.")
		}
		Printf("\n")
		return true
	}
	return false
}

func cyclopsRoomFcn(arg ActionArg) bool {
	if arg == ActLook {
		Printf("This room has an exit on the northwest, and a staircase leading up.\n")
		if gD().CyclopsFlag && !gD().MagicFlag {
			Printf("The cyclops is sleeping blissfully at the foot of the stairs.\n")
		} else if gD().MagicFlag {
			Printf("The east wall, previously solid, now has a cyclops-sized opening in it.\n")
		} else if gD().CycloWrath == 0 {
			Printf("A cyclops, who looks prepared to eat horses (much less mere adventurers), blocks the staircase. From his state of health, and the bloodstains on the walls, you gather that he is not very friendly, though he likes people.\n")
		} else if gD().CycloWrath > 0 {
			Printf("The cyclops is standing in the corner, eyeing you closely. I don't think he likes you very much. He looks extremely hungry, even for a cyclops.\n")
		} else {
			Printf("The cyclops, having eaten the hot peppers, appears to be gasping. His enflamed tongue protrudes from his man-sized mouth.\n")
		}
		return true
	}
	if arg == ActEnter {
		if gD().CycloWrath == 0 {
			return false
		}
		Queue("iCyclops", -1).Run = true
		return false
	}
	return false
}

func treasureRoomFcn(arg ActionArg) bool {
	if arg == ActEnter && !gD().Dead {
		if !thief.IsIn(G.Here) {
			Printf("You hear a scream of anguish as you violate the robber's hideaway. Using passages unknown to you, he rushes to its defense.\n")
			thief.MoveTo(G.Here)
		}
		thief.Give(FlgFight)
		thief.Take(FlgInvis)
		thiefInTreasure()
		return true
	}
	return false
}

func reservoirSouthFcn(arg ActionArg) bool {
	if arg == ActLook {
		if gD().LowTide && gD().GatesOpen {
			Printf("You are in a long room, to the north of which was formerly a lake. However, with the water level lowered, there is merely a wide stream running through the center of the room.")
		} else if gD().GatesOpen {
			Printf("You are in a long room. To the north is a large lake, too deep to cross. You notice, however, that the water level appears to be dropping at a rapid rate. Before long, it might be possible to cross to the other side from here.")
		} else if gD().LowTide {
			Printf("You are in a long room, to the north of which is a wide area which was formerly a reservoir, but now is merely a stream. You notice, however, that the level of the stream is rising quickly and that before long it will be impossible to cross here.")
		} else {
			Printf("You are in a long room on the south shore of a large lake, far too deep and wide for crossing.")
		}
		Printf("\nThere is a path along the stream to the east or west, a steep pathway climbing southwest along the edge of a chasm, and a path leading into a canyon to the southeast.\n")
		return true
	}
	return false
}

func reservoirFcn(arg ActionArg) bool {
	if arg == ActEnd && !G.Winner.Location().Has(FlgVeh) && !gD().GatesOpen && gD().LowTide {
		Printf("You notice that the water level here is rising rapidly. The currents are also becoming stronger. Staying here seems quite perilous!\n")
		return true
	}
	if arg == ActLook {
		if gD().LowTide {
			Printf("You are on what used to be a large lake, but which is now a large mud pile. There are \"shores\" to the north and south.")
		} else {
			Printf("You are on the lake. Beaches can be seen north and south. Upstream a small stream enters the lake through a narrow cleft in the rocks. The dam can be seen downstream.")
		}
		Printf("\n")
		return true
	}
	return false
}

func reservoirNorthFcn(arg ActionArg) bool {
	if arg == ActLook {
		if gD().LowTide && gD().GatesOpen {
			Printf("You are in a large cavernous room, the south of which was formerly a lake. However, with the water level lowered, there is merely a wide stream running through there.")
		} else if gD().GatesOpen {
			Printf("You are in a large cavernous area. To the south is a wide lake, whose water level appears to be falling rapidly.")
		} else if gD().LowTide {
			Printf("You are in a cavernous area, to the south of which is a very wide stream. The level of the stream is rising rapidly, and it appears that before long it will be impossible to cross to the other side.")
		} else {
			Printf("You are in a large cavernous room, north of a large lake.")
		}
		Printf("\nThere is a slimy stairway leaving the room to the north.\n")
		return true
	}
	return false
}

func mirrorRoomFcn(arg ActionArg) bool {
	if arg == ActLook {
		Printf("You are in a large square room with tall ceilings. On the south wall is an enormous mirror which fills the entire wall. There are exits on the other three sides of the room.\n")
		if gD().MirrorMung {
			Printf("Unfortunately, the mirror has been destroyed by your recklessness.\n")
		}
		return true
	}
	return false
}

func cave2RoomFcn(arg ActionArg) bool {
	if arg == ActEnd {
		if candles.IsIn(G.Winner) && Prob(50, true) && candles.Has(FlgOn) {
			QueueInt("iCandles", false).Run = false
			candles.Take(FlgOn)
			Printf("A gust of wind blows out your candles!\n")
			G.Lit = IsLit(G.Here, true)
			if !G.Lit {
				Printf("it is now completely dark.\n")
			}
			return true
		}
	}
	return false
}

func lLDRoomFcn(arg ActionArg) bool {
	if arg == ActLook {
		Printf("You are outside a large gateway, on which is inscribed\n\n  Abandon every hope\nall ye who enter here!\n\nThe gate is open; through it you can see a desolation, with a pile of mangled bodies in one corner. Thousands of voices, lamenting some hideous fate, can be heard.\n")
		if !gD().LLDFlag && !gD().Dead {
			Printf("The way through the gate is barred by evil spirits, who jeer at your attempts to pass.\n")
		}
		return true
	}
	if arg == ActBegin {
		switch G.ActVerb.Norm {
		case "exorcise":
			if !gD().LLDFlag {
				if bell.IsIn(G.Winner) && book.IsIn(G.Winner) && candles.IsIn(G.Winner) {
					Printf("You must perform the ceremony.\n")
				} else {
					Printf("You aren't equipped for an exorcism.\n")
				}
				return true
			}
		case "ring":
			if !gD().LLDFlag && G.DirObj == &bell {
				gD().XB = true
				removeCarefully(&bell)
				thisIsIt(&hotBell)
				hotBell.MoveTo(G.Here)
				Printf("The bell suddenly becomes red hot and falls to the ground. The wraiths, as if paralyzed, stop their jeering and slowly turn to face you. On their ashen faces, the expression of a long-forgotten terror takes shape.\n")
				if candles.IsIn(G.Winner) {
					Printf("In your confusion, the candles drop to the ground (and they are out).\n")
					candles.MoveTo(G.Here)
					candles.Take(FlgOn)
					QueueInt("iCandles", false).Run = false
				}
				Queue("iXb", 6).Run = true
				Queue("iXbh", 20).Run = true
				return true
			}
		case "read":
			if gD().XC && G.DirObj == &book && !gD().LLDFlag {
				Printf("Each word of the prayer reverberates through the hall in a deafening confusion. As the last word fades, a voice, loud and commanding, speaks: \"Begone, fiends!\" A heart-stopping scream fills the cavern, and the spirits, sensing a greater power, flee through the walls.\n")
				removeCarefully(&ghosts)
				gD().LLDFlag = true
				QueueInt("iXc", false).Run = false
				return true
			}
		}
	}
	if arg == ActEnd {
		if gD().XB && candles.IsIn(G.Winner) && candles.Has(FlgOn) && !gD().XC {
			gD().XC = true
			Printf("The flames flicker wildly and appear to dance. The earth beneath your feet trembles, and your legs nearly buckle beneath you. The spirits cower at your unearthly power.\n")
			QueueInt("iXb", false).Run = false
			Queue("iXc", 3).Run = true
		}
	}
	return false
}

func domeRoomFcn(arg ActionArg) bool {
	if arg == ActLook {
		Printf("You are at the periphery of a large dome, which forms the ceiling of another room below. Protecting you from a precipitous drop is a wooden railing which circles the dome.\n")
		if gD().DomeFlag {
			Printf("Hanging down from the railing is a rope which ends about ten feet from the floor below.\n")
		}
		return true
	}
	if arg == ActEnter {
		if gD().Dead {
			Printf("As you enter the dome you feel a strong pull as if from a wind drawing you over the railing and down.\n")
			G.Winner.MoveTo(&torchRoom)
			G.Here = &torchRoom
			return true
		}
		switch G.ActVerb.Norm {
		case "leap":
			jigsUp("I'm afraid that the leap you attempted has done you in.", false)
			return true
		}
	}
	return false
}

func torchRoomFcn(arg ActionArg) bool {
	if arg == ActLook {
		Printf("This is a large room with a prominent doorway leading to a down staircase. Above you is a large dome. Up around the edge of the dome (20 feet up) is a wooden railing. In the center of the room sits a white marble pedestal.\n")
		if gD().DomeFlag {
			Printf("A piece of rope descends from the railing above, ending some five feet above your head.\n")
		}
		return true
	}
	return false
}

func southTempleFcn(arg ActionArg) bool {
	if arg == ActBegin {
		gD().CoffinCure = !coffin.IsIn(G.Winner)
		return false
	}
	return false
}

func machineRoomFcn(arg ActionArg) bool {
	if arg == ActLook {
		Printf("This is a large, cold room whose sole exit is to the north. In one corner there is a machine which is reminiscent of a clothes dryer. On its face is a switch which is labelled \"START\". The switch does not appear to be manipulable by any human hand (unless the fingers are about 1/16 by 1/4 inch). On the front of the machine is a large lid, which is ")
		if machine.Has(FlgOpen) {
			Printf("open.\n")
		} else {
			Printf("closed.\n")
		}
		return true
	}
	return false
}

func loudRoomFcn(arg ActionArg) bool {
	if arg == ActLook {
		Printf("This is a large room with a ceiling which cannot be detected from the ground. There is a narrow passage from east to west and a stone stairway leading upward.")
		if gD().LoudFlag || (!gD().GatesOpen && gD().LowTide) {
			Printf(" The room is eerie in its quietness.")
		} else {
			Printf(" The room is deafeningly loud with an undetermined rushing sound. The sound seems to reverberate from all of the walls, making it difficult even to think.")
		}
		Printf("\n")
		return true
	}
	if arg == ActEnd && gD().GatesOpen && !gD().LowTide {
		Printf("it is unbearably loud here, with an ear-splitting roar seeming to come from all around you. There is a pounding in your head which won't stop. With a tremendous effort, you scramble out of the room.\n\n")
		dest := loudRuns[G.Rand.Intn(len(loudRuns))]
		moveToRoom(dest, true)
		return false
	}
	if arg == ActEnter {
		if gD().LoudFlag || (!gD().GatesOpen && gD().LowTide) {
			return false
		}
		if gD().GatesOpen && !gD().LowTide {
			return false
		}
		// Room is loud - special input handling
		vFirstLook(ActUnk)
		if G.Params.Continue != NumUndef {
			Printf("The rest of your commands have been lost in the noise.\n")
			G.Params.Continue = NumUndef
		}
		// In the original, this has a special read loop. We simplify.
		return false
	}
	switch G.ActVerb.Norm {
	case "echo":
		if gD().LoudFlag || (!gD().GatesOpen && gD().LowTide) {
			// Room is already quiet
			Printf("echo echo ...\n")
			return true
		}
		Printf("The acoustics of the room change subtly.\n")
		gD().LoudFlag = true
		return true
	}
	return false
}

func deepCanyonFcn(arg ActionArg) bool {
	if arg == ActLook {
		Printf("You are on the south edge of a deep canyon. Passages lead off to the east, northwest and southwest. A stairway leads down.")
		if gD().GatesOpen && !gD().LowTide {
			Printf(" You can hear a loud roaring sound, like that of rushing water, from below.")
		} else if !gD().GatesOpen && gD().LowTide {
			Printf("\n")
			return true
		} else {
			Printf(" You can hear the sound of flowing water from below.")
		}
		Printf("\n")
		return true
	}
	return false
}

func boomRoomFcn(arg ActionArg) bool {
	if arg == ActEnd {
		dummy := false
		switch G.ActVerb.Norm {
		case "lamp on", "burn":
			if G.DirObj == &candles || G.DirObj == &torch || G.DirObj == &match {
				dummy = true
			}
		}
		if (candles.IsIn(G.Winner) && candles.Has(FlgOn)) ||
			(torch.IsIn(G.Winner) && torch.Has(FlgOn)) ||
			(match.IsIn(G.Winner) && match.Has(FlgOn)) {
			if dummy {
				Printf("How sad for an aspiring adventurer to light a %s in a room which reeks of gas. Fortunately, there is justice in the world.\n", G.DirObj.Desc)
			} else {
				Printf("Oh dear. it appears that the smell coming from this room was coal gas. I would have thought twice about carrying flaming objects in here.\n")
			}
			jigsUp("\n      ** BOOOOOOOOOOOM **", false)
			return true
		}
	}
	return false
}

func batsRoomFcn(arg ActionArg) bool {
	if arg == ActLook {
		Printf("You are in a small room which has doors only to the east and south.\n")
		return true
	}
	if arg == ActEnter && !gD().Dead {
		if garlic.Location() != G.Winner && !garlic.IsIn(G.Here) {
			vLook(ActUnk)
			Printf("\n")
			flyMe()
			return true
		}
	}
	return false
}

func noObjsFcn(arg ActionArg) bool {
	if arg == ActBegin {
		f := G.Winner.Children
		gD().EmptyHanded = true
		for _, child := range f {
			if weight(child) > 4 {
				gD().EmptyHanded = false
				break
			}
		}
		if G.Here == &lowerShaft && G.Lit {
			scoreUpd(gD().LightShaft)
			gD().LightShaft = 0
		}
		return false
	}
	return false
}

func canyonViewFcn(arg ActionArg) bool {
	return false
}

func forestRoomFcn(arg ActionArg) bool {
	if arg == ActEnter {
		Queue("iForestRandom", -1).Run = true
		return false
	}
	if arg == ActBegin {
		switch G.ActVerb.Norm {
		case "climb", "climb up":
			if G.DirObj == &tree {
				doWalk(Up)
				return true
			}
		}
	}
	return false
}

func treeRoomFcn(arg ActionArg) bool {
	if arg == ActLook {
		Printf("You are about 10 feet above the ground nestled among some large branches. The nearest branch above you is above your reach.\n")
		if path.HasChildren() && len(path.Children) > 0 {
			Printf("On the ground below you can see:  ")
			printContents(&path)
			Printf(".\n")
		}
		return true
	}
	if arg == ActBegin {
		switch G.ActVerb.Norm {
		case "climb down":
			if G.DirObj == &tree || G.DirObj == &rooms {
				doWalk(Down)
				return true
			}
		case "climb up", "climb":
			if G.DirObj == &tree {
				doWalk(Up)
				return true
			}
		case "drop":
			if !iDrop() {
				return true
			}
			if G.DirObj == &nest && egg.IsIn(&nest) {
				Printf("The nest falls to the ground, and the egg spills out of it, seriously damaged.\n")
				removeCarefully(&egg)
				brokenEgg.MoveTo(&path)
				return true
			}
			if G.DirObj == &egg {
				Printf("The egg falls to the ground and springs open, seriously damaged.")
				egg.MoveTo(&path)
				badEgg()
				Printf("\n")
				return true
			}
			if G.DirObj != G.Winner && G.DirObj != &tree {
				G.DirObj.MoveTo(&path)
				Printf("The %s falls to the ground.\n", G.DirObj.Desc)
			}
			return true
		}
	}
	if arg == ActEnter {
		Queue("iForestRandom", -1).Run = true
	}
	return false
}

func deadFunction(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "walk":
		if G.Here == &timberRoom && G.Params.HasWalkDir && G.Params.WalkDir == West {
			Printf("You cannot enter in your condition.\n")
			return true
		}
		return false
	case "brief", "verbose", "super-brief", "version":
		return false
	case "attack", "mung", "alarm", "swing":
		Printf("All such attacks are vain in your condition.\n")
		return true
	case "open", "close", "eat", "drink", "inflate", "deflate", "turn", "burn", "tie", "untie", "rub":
		Printf("Even such an action is beyond your capabilities.\n")
		return true
	case "wait":
		Printf("Might as well. You've got an eternity.\n")
		return true
	case "lamp on":
		Printf("You need no light to guide you.\n")
		return true
	case "score":
		Printf("You're dead! How can you think of your score?\n")
		return true
	case "take":
		Printf("Your hand passes through its object.\n")
		return true
	case "drop", "throw", "inventory":
		Printf("You have no possessions.\n")
		return true
	case "diagnose":
		Printf("You are dead.\n")
		return true
	case "look":
		Printf("The room looks strange and unearthly")
		if !G.Here.HasChildren() {
			Printf(".")
		} else {
			Printf(" and objects appear indistinct.")
		}
		Printf("\n")
		if !G.Here.Has(FlgOn) {
			Printf("Although there is no light, the room seems dimly illuminated.\n")
		}
		Printf("\n")
		return false
	case "pray":
		if G.Here == &southTemple {
			lamp.Take(FlgInvis)
			G.Winner.Action = nil
			G.AlwaysLit = false
			gD().Dead = false
			if troll.IsIn(&trollRoom) {
				gD().TrollFlag = false
			}
			Printf("From the distance the sound of a lone trumpet is heard. The room becomes very bright and you feel disembodied. In a moment, the brightness fades and you find yourself rising as if from a long sleep, deep in the woods. In the distance you can faintly hear a songbird and the sounds of the forest.\n\n")
			moveToRoom(&forest1, true)
			return true
		}
		Printf("Your prayers are not heard.\n")
		return true
	}
	Printf("You can't even do that.\n")
	G.Params.Continue = NumUndef
	return true
}

// ================================================================
// INTERRUPT ROUTINES
// ================================================================

func iCandles() bool {
	candles.Give(FlgTouch)
	if gD().CandleTableIdx >= len(candleTable) {
		return true
	}
	tick := candleTable[gD().CandleTableIdx].(int)
	Queue("iCandles", tick).Run = true
	lightInt(&candles, gD().CandleTableIdx, tick)
	if tick != 0 {
		gD().CandleTableIdx += 2
	}
	return true
}

func iLantern() bool {
	if gD().LampTableIdx >= len(lampTable) {
		return true
	}
	tick := lampTable[gD().LampTableIdx].(int)
	Queue("iLantern", tick).Run = true
	lightInt(&lamp, gD().LampTableIdx, tick)
	if tick != 0 {
		gD().LampTableIdx += 2
	}
	return true
}

// lightInt handles light source countdown warnings and expiry
func lightInt(obj *Object, tblIdx, tick int) {
	if tick == 0 {
		obj.Take(FlgOn)
		obj.Give(FlgDestroyed)
	}
	if IsHeld(obj) || obj.IsIn(G.Here) {
		if tick == 0 {
			Printf("You'd better have more light than from the %s.\n", obj.Desc)
		} else {
			// Print the warning message from the table
			var tbl []interface{}
			if obj == &candles {
				tbl = candleTable
			} else {
				tbl = lampTable
			}
			if tblIdx+1 < len(tbl) {
				if msg, ok := tbl[tblIdx+1].(string); ok {
					Printf("%s\n", msg)
				}
			}
		}
	}
}

// iCure heals the player gradually
func iCure() bool {
	s := G.Winner.GetStrength()
	if s > 0 {
		s = 0
		G.Winner.SetStrength(s)
	} else if s < 0 {
		s++
		G.Winner.SetStrength(s)
	}
	if s < 0 {
		if gD().LoadAllowed < gD().LoadMax {
			gD().LoadAllowed += 10
		}
		Queue("iCure", cureWait).Run = true
	} else {
		gD().LoadAllowed = gD().LoadMax
		QueueInt("iCure", false).Run = false
	}
	return false
}

func iMatch() bool {
	Printf("The match has gone out.\n")
	match.Take(FlgFlame)
	match.Take(FlgOn)
	G.Lit = IsLit(G.Here, true)
	return true
}

func iXb() bool {
	if !gD().XC {
		if G.Here == &enteranceToHades {
			Printf("The tension of this ceremony is broken, and the wraiths, amused but shaken at your clumsy attempt, resume their hideous jeering.\n")
		}
	}
	gD().XB = false
	return true
}

func iXbh() bool {
	removeCarefully(&hotBell)
	bell.MoveTo(&enteranceToHades)
	if G.Here == &enteranceToHades {
		Printf("The bell appears to have cooled down.\n")
	}
	return true
}

func iXc() bool {
	gD().XC = false
	iXb()
	return true
}

func iCyclops() bool {
	if gD().CyclopsFlag || gD().Dead {
		return true
	}
	if G.Here != &cyclopsRoom {
		QueueInt("iCyclops", false).Run = false
		return false
	}
	if gD().CycloWrath < -5 || gD().CycloWrath > 5 {
		QueueInt("iCyclops", false).Run = false
		jigsUp("The cyclops, tired of all of your games and trickery, grabs you firmly. As he licks his chops, he says \"Mmm. Just like Mom used to make 'em.\" it's nice to be appreciated.", false)
		return true
	}
	if gD().CycloWrath < 0 {
		gD().CycloWrath--
	} else {
		gD().CycloWrath++
	}
	if !gD().CyclopsFlag {
		idx := gD().CycloWrath
		if idx < 0 {
			idx = -idx
		}
		idx -= 2
		if idx >= 0 && idx < len(cyclomad) {
			Printf("%s\n", cyclomad[idx])
		}
	}
	return true
}

func iForestRandom() bool {
	if !forestRoomQ() {
		QueueInt("iForestRandom", false).Run = false
		return false
	}
	if Prob(15, false) {
		Printf("You hear in the distance the chirping of a song bird.\n")
	}
	return true
}

// ================================================================
// EXIT FUNCTIONS
// ================================================================

func gratingExitFcn() *Object {
	if gD().GrateRevealed {
		if grate.Has(FlgOpen) {
			return &gratingRoom
		}
		Printf("The grating is closed!\n")
		thisIsIt(&grate)
		return nil
	}
	Printf("You can't go that way.\n")
	return nil
}

func trapDoorExitFcn() *Object {
	if gD().RugMoved {
		if trapDoor.Has(FlgOpen) {
			return &cellar
		}
		Printf("The trap door is closed.\n")
		thisIsIt(&trapDoor)
		return nil
	}
	Printf("You can't go that way.\n")
	return nil
}

func upChimneyFcn() *Object {
	f := G.Winner.Children
	if len(f) == 0 {
		Printf("Going up empty-handed is a bad idea.\n")
		return nil
	}
	// Check if player is carrying at most 1-2 items including the lamp
	count := 0
	for range f {
		count++
	}
	if count <= 2 && lamp.IsIn(G.Winner) {
		if !trapDoor.Has(FlgOpen) {
			trapDoor.Take(FlgTouch)
		}
		return &kitchen
	}
	Printf("You can't get up there with what you're carrying.\n")
	return nil
}

func mazeDiodesFcn() *Object {
	Printf("You won't be able to get back up to the tunnel you are going through when it gets to the next room.\n\n")
	if G.Here == &maze2 {
		return &maze4
	}
	if G.Here == &maze7 {
		return &deadEnd1
	}
	if G.Here == &maze9 {
		return &maze11
	}
	if G.Here == &maze12 {
		return &maze5
	}
	return nil
}

// ================================================================
// PSEUDO FUNCTIONS
// ================================================================

func chasmPseudo(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "leap":
		Printf("You look before leaping, and realize that you would never survive.\n")
		return true
	case "put":
		if G.DirObj == &me {
			Printf("You look before leaping, and realize that you would never survive.\n")
			return true
		}
		if G.IndirObj == &pseudoObject {
			Printf("The %s drops out of sight into the chasm.\n", G.DirObj.Desc)
			removeCarefully(G.DirObj)
			return true
		}
	case "cross":
		Printf("it's too far to jump, and there's no bridge.\n")
		return true
	case "throw off":
		if G.IndirObj == &pseudoObject {
			Printf("The %s drops out of sight into the chasm.\n", G.DirObj.Desc)
			removeCarefully(G.DirObj)
			return true
		}
	}
	return false
}

func lakePseudo(arg ActionArg) bool {
	if gD().LowTide {
		Printf("There's not much lake left....\n")
		return true
	}
	switch G.ActVerb.Norm {
	case "cross":
		Printf("it's too wide to cross.\n")
		return true
	case "through":
		Printf("You can't swim in this lake.\n")
		return true
	}
	return false
}

func streamPseudo(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "swim", "through":
		Printf("You can't swim in the stream.\n")
		return true
	case "cross":
		Printf("The other side is a sheer rock cliff.\n")
		return true
	}
	return false
}

func domePseudo(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "kiss":
		Printf("No.\n")
		return true
	}
	return false
}

func gatePseudo(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "through":
		doWalk(In)
		return true
	}
	Printf("The gate is protected by an invisible force. it makes your teeth ache to touch it.\n")
	return true
}

func doorPseudo(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "open", "close":
		Printf("The door won't budge.\n")
		return true
	case "through":
		doWalk(South)
		return true
	}
	return false
}

func paintPseudo(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "mung":
		Printf("Some paint chips away, revealing more paint.\n")
		return true
	}
	return false
}

func gasPseudo(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "breathe":
		Printf("There is too much gas to blow away.\n")
		return true
	case "smell":
		Printf("it smells like coal gas in here.\n")
		return true
	}
	return false
}

func chainPseudo(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "take", "move":
		Printf("The chain is secure.\n")
		return true
	case "raise", "lower":
		Printf("Perhaps you should do that to the basket.\n")
		return true
	case "examine":
		Printf("The chain secures a basket within the shaft.\n")
		return true
	}
	return false
}
