package game

import . "github.com/ajdnik/gozork/engine"

var (
	// HelloSailor counts occurrences of 'hello, sailor'
	// IsSprayed is a flag indicating if the player is wearing grue repellent
	version = [24]byte{0, 0, 0, 0, 119, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 56, 56, 48, 52, 50, 57}
	// DescObj stores the last object which was described
	indents = [6]string{
		"",
		"  ",
		"    ",
		"      ",
		"        ",
		"          ",
	}
	yuks = RndSelect{
		Unselected: []string{
			"A valiant attempt.",
			"You can't be serious.",
			"An interesting idea...",
			"What a concept!",
		},
	}
	hellos = RndSelect{
		Unselected: []string{
			"Hello.",
			"Good day.",
			"Nice weather we've been having lately.",
			"Goodbye.",
		},
	}
	jumpLoss = RndSelect{
		Unselected: []string{
			"You should have looked before you leaped.",
			"In the movies, your life would be passing before your eyes.",
			"Geronimo...",
		},
	}
	wheeee = RndSelect{
		Unselected: []string{
			"Very good. Now you can go to the second grade.",
			"Are you enjoying yourself?",
			"Wheeeeeeeeee!!!!!",
			"Do you expect me to applaud?",
		},
	}
	hohum = RndSelect{
		Unselected: []string{
			" doesn't seem to work.",
			" isn't notably helpful.",
			" has no effect.",
		},
	}
)

func vAlarm(arg ActionArg) bool {
	if !G.DirObj.Has(FlgPerson) {
		Printf("The %s isn't sleeping.\n", G.DirObj.Desc)
		return true
	}
	if G.DirObj.GetStrength() <= 0 {
		Printf("He's wide awake, or haven't you noticed...\n")
		return true
	}
	Printf("The %s is rudely awakened.\n", G.DirObj.Desc)
	return awaken(G.DirObj)
}

func vAnswer(arg ActionArg) bool {
	Printf("Nobody seems to be awaiting your answer.\n")
	G.Params.Continue = NumUndef
	G.Params.InQuotes = false
	return true
}

func vAttack(arg ActionArg) bool {
	if !G.DirObj.Has(FlgPerson) {
		Printf("I've known strange people, but fighting a %s?\n", G.DirObj.Desc)
		return true
	}
	if G.IndirObj == nil || G.IndirObj == &hands {
		Printf("Trying to attack a %s with your bare hands is suicidal.\n", G.DirObj.Desc)
		return true
	}
	if !G.IndirObj.IsIn(G.Winner) {
		Printf("You aren't even holding the %s.\n", G.IndirObj.Desc)
		return true
	}
	if !G.IndirObj.Has(FlgWeapon) {
		Printf("Trying to attack the %s with a %s is suicidal.\n", G.DirObj.Desc, G.IndirObj.Desc)
		return true
	}
	return heroBlow()
}

func vBack(arg ActionArg) bool {
	Printf("Sorry, my memory is poor. Please give a direction.\n")
	return true
}

func vBlast(arg ActionArg) bool {
	Printf("You can't blast anything by using words.\n")
	return true
}

func vBreathe(arg ActionArg) bool {
	ret := Perform(ActionVerb{Norm: "inflate", Orig: "inflate"}, G.DirObj, &lungs)
	if ret == PerfFatal {
		return RFatal()
	}
	return ret == PerfHndld
}

func vBrush(arg ActionArg) bool {
	Printf("If you wish, but heaven only knows why.\n")
	return true
}

func preBurn(arg ActionArg) bool {
	if G.IndirObj == nil {
		Printf("You didn't say with what!\n")
		return true
	}
	if IsFlaming(G.IndirObj) {
		return false
	}
	Printf("With a %s??!?\n", G.IndirObj.Desc)
	return true
}

func vBurn(arg ActionArg) bool {
	if !G.DirObj.Has(FlgBurn) {
		Printf("You can't burn a %s.\n", G.DirObj.Desc)
		return true
	}
	if !G.DirObj.IsIn(G.Winner) && !G.Winner.IsIn(G.DirObj) {
		removeCarefully(G.DirObj)
		Printf("The %s catches fire and is consumed.\n", G.DirObj.Desc)
		return true
	}
	removeCarefully(G.DirObj)
	Printf("The %s catches fire. Unfortunately, you were ", G.DirObj.Desc)
	if G.Winner.IsIn(G.DirObj) {
		Printf("in")
	} else {
		Printf("holding")
	}
	return jigsUp(" it at the time.", false)
}

func vChomp(arg ActionArg) bool {
	Printf("Preposterous!\n")
	return true
}

func vClose(arg ActionArg) bool {
	if !G.DirObj.Has(FlgCont) && !G.DirObj.Has(FlgDoor) {
		Printf("You must tell me how to do that to a %s.\n", G.DirObj.Desc)
		return true
	}
	if !G.DirObj.Has(FlgSurf) && G.DirObj.GetCapacity() != 0 {
		if G.DirObj.Has(FlgOpen) {
			G.DirObj.Take(FlgOpen)
			Printf("Closed.\n")
			if G.Lit {
				G.Lit = IsLit(G.Here, true)
				if !G.Lit {
					Printf("it is now pitch black.\n")
				}
			}
			return true
		}
		Printf("it is already closed.\n")
		return true
	}
	if G.DirObj.Has(FlgDoor) {
		if G.DirObj.Has(FlgOpen) {
			G.DirObj.Take(FlgOpen)
			Printf("The %s is now closed.\n", G.DirObj.Desc)
			return true
		}
		Printf("it is already closed.\n")
		return true
	}
	Printf("You cannot close that.\n")
	return true
}

func vCommand(arg ActionArg) bool {
	if G.DirObj.Has(FlgPerson) {
		Printf("The %s pays no attention.\n", G.DirObj.Desc)
		return true
	}
	Printf("You cannot talk to that!\n")
	return true
}

func vCount(arg ActionArg) bool {
	if G.DirObj == &blessings {
		Printf("Well, for one, you are playing Zork...\n")
		return true
	}
	Printf("You have lost your mind.\n")
	return true
}

func vCurses(arg ActionArg) bool {
	if G.DirObj == nil {
		Printf("Such language in a high-class establishment like this!\n")
		return true
	}
	if G.DirObj.Has(FlgPerson) {
		Printf("Insults of this nature won't help you.\n")
		return true
	}
	Printf("What a loony!\n")
	return true
}

func vCut(arg ActionArg) bool {
	if G.DirObj.Has(FlgPerson) {
		ret := Perform(ActionVerb{Norm: "attack", Orig: "attack"}, G.DirObj, G.IndirObj)
		if ret == PerfFatal {
			return RFatal()
		}
		return ret == PerfHndld
	}
	if G.DirObj.Has(FlgBurn) && G.IndirObj.Has(FlgWeapon) {
		if G.Winner.IsIn(G.DirObj) {
			Printf("Not a bright idea, especially since you're in it.\n")
			return true
		}
		removeCarefully(G.DirObj)
		Printf("Your skillful %ssmanship slices the %s into innumerable slivers which blow away.\n", G.IndirObj.Desc, G.DirObj.Desc)
		return true
	}
	if !G.IndirObj.Has(FlgWeapon) {
		Printf("The \"cutting edge\" of a %s is hardly adequate.\n", G.IndirObj.Desc)
		return true
	}
	Printf("Strange concept, cutting the %s....\n", G.DirObj.Desc)
	return true
}

func vDeflate(arg ActionArg) bool {
	Printf("Come on, now!\n")
	return true
}

func vDig(arg ActionArg) bool {
	if G.IndirObj == nil {
		G.IndirObj = &hands
	}
	if G.IndirObj == &shovel {
		Printf("There's no reason to be digging here.\n")
		return true
	}
	if G.IndirObj.Has(FlgTool) {
		Printf("Digging with the %s is slow and tedious.\n", G.IndirObj.Desc)
		return true
	}
	Printf("Digging with a %s is silly.\n", G.IndirObj.Desc)
	return true
}

func vDisenchant(arg ActionArg) bool {
	Printf("Nothing happens.\n")
	return true
}

func vDrink(arg ActionArg) bool {
	return vEat(ActUnk)
}

func vDrinkFrom(act ActionArg) bool {
	Printf("How peculiar!\n")
	return true
}

func preDrop(arg ActionArg) bool {
	if G.DirObj == G.Winner.Location() {
		Perform(ActionVerb{Norm: "disembark", Orig: "disembark"}, G.DirObj, nil)
		return true
	}
	return false
}

func vDrop(arg ActionArg) bool {
	if iDrop() {
		Printf("Dropped.\n")
		return true
	}
	return false
}

func vEat(arg ActionArg) bool {
	isEat := G.DirObj.Has(FlgFood)
	if isEat {
		if !G.DirObj.IsIn(G.Winner) && !G.DirObj.Location().IsIn(G.Winner) {
			Printf("You're not holding that.\n")
			return true
		}
		if G.ActVerb.Norm == "drink" {
			Printf("How can you drink that?\n")
			return true
		}
		Printf("Thank you very much. it really hit the spot.\n")
		removeCarefully(G.DirObj)
		return true
	}
	isDrink := G.DirObj.Has(FlgDrink)
	if isDrink {
		nobj := G.DirObj.Location()
		if G.DirObj.IsIn(&globalObjects) || IsInGlobal(&globalWater, G.Here) || G.DirObj == &pseudoObject {
			return hitSpot()
		}
		if nobj == nil || !IsAccessible(nobj) {
			Printf("There isn't any water here.\n")
			return true
		}
		if IsAccessible(nobj) && !nobj.IsIn(G.Winner) {
			Printf("You have to be holding the %s first.\n", nobj.Desc)
			return true
		}
		if !nobj.Has(FlgOpen) {
			Printf("You'll have to open the %s first.\n", nobj.Desc)
			return true
		}
		return hitSpot()
	}
	if !isEat && !isDrink {
		Printf("I don't think that the %s would agree with you.\n", G.DirObj.Desc)
		return true
	}
	return false
}

func hitSpot() bool {
	if G.DirObj == &water && !IsInGlobal(&globalWater, G.Here) {
		removeCarefully(G.DirObj)
	}
	Printf("Thank you very much. I was rather thirsty (from all this talking, probably).\n")
	return true
}

func vEcho(arg ActionArg) bool {
	if len(G.LexRes) <= 0 {
		Printf("echo echo ...\n")
		return true
	}
	wrd := G.LexRes[len(G.LexRes)-1]
	Printf("%s\n", wrd.Orig+" "+wrd.Orig+" ...")
	return true
}

func vEnchant(arg ActionArg) bool {
	return vDisenchant(ActUnk)
}

func vExcorcise(arg ActionArg) bool {
	Printf("What a bizarre concept!\n")
	return true
}

func preFill(arg ActionArg) bool {
	if G.IndirObj == nil {
		if IsInGlobal(&globalWater, G.Here) {
			Perform(ActionVerb{Norm: "fill", Orig: "fill"}, G.DirObj, &globalWater)
			return true
		}
		if water.IsIn(G.Winner.Location()) {
			Perform(ActionVerb{Norm: "fill", Orig: "fill"}, G.DirObj, &water)
			return true
		}
		Printf("There is nothing to fill it with.\n")
		return true
	}
	if G.IndirObj == &water {
		return false
	}
	if G.IndirObj != &globalWater {
		Perform(ActionVerb{Norm: "put", Orig: "put"}, G.IndirObj, G.DirObj)
		return true
	}
	return false
}

func vFill(arg ActionArg) bool {
	if G.IndirObj != nil {
		Printf("You may know how to do that, but I don't.\n")
		return true
	}
	if IsInGlobal(&globalWater, G.Here) {
		Perform(ActionVerb{Norm: "fill", Orig: "fill"}, G.DirObj, &globalWater)
		return true
	}
	if water.IsIn(G.Winner.Location()) {
		Perform(ActionVerb{Norm: "fill", Orig: "fill"}, G.DirObj, &water)
		return true
	}
	Printf("There is nothing to fill it with.\n")
	return true
}

func vFirstLook(arg ActionArg) bool {
	if describeRoom(false) {
		if !G.SuperBrief {
			return describeObjects(false)
		}
	}
	return false
}

func vFind(arg ActionArg) bool {
	if G.DirObj == &hands || G.DirObj == &lungs {
		Printf("Within six feet of your head, assuming you haven't left that somewhere.\n")
		return true
	}
	if G.DirObj == &me {
		Printf("You're around here somewhere...\n")
		return true
	}
	l := G.DirObj.Location()
	if l == &globalObjects {
		Printf("You find it.\n")
		return true
	}
	if G.DirObj.IsIn(G.Winner) {
		Printf("You have it.\n")
		return true
	}
	if G.DirObj.IsIn(G.Here) || IsInGlobal(G.DirObj, G.Here) || G.DirObj == &pseudoObject {
		Printf("it's right here.\n")
		return true
	}
	if l.Has(FlgPerson) {
		Printf("The %s has it.\n", l.Desc)
		return true
	}
	if l.Has(FlgSurf) {
		Printf("it's on the %s.\n", l.Desc)
		return true
	}
	if l.Has(FlgCont) {
		Printf("it's in the %s.\n", l.Desc)
		return true
	}
	Printf("Beats me.\n")
	return true
}

func vFrobozz(arg ActionArg) bool {
	Printf("The FROBOZZ Corporation created, owns, and operates this dungeon.\n")
	return true
}

func preGive(arg ActionArg) bool {
	if !IsHeld(G.DirObj) {
		Printf("That's easy for you to say since you don't even have the %s.\n", G.DirObj.Desc)
		return true
	}
	return false
}

func vGive(arg ActionArg) bool {
	if !G.DirObj.Has(FlgPerson) {
		Printf("You can't give a %s to a %s!\n", G.DirObj.Desc, G.IndirObj.Desc)
		return true
	}
	Printf("The %s refuses it politely.\n", G.IndirObj.Desc)
	return true
}

func vHatch(arg ActionArg) bool {
	Printf("Bizarre!\n")
	return true
}

func vHello(arg ActionArg) bool {
	if G.DirObj == nil {
		Printf("%s\n", PickOne(hellos))
		return true
	}
	if G.DirObj.Has(FlgPerson) {
		Printf("The %s bows his head to you in greeting.\n", G.DirObj.Desc)
		return true
	}
	Printf("it's a well known fact that only schizophrenics say \"Hello\" to a %s.\n", G.DirObj.Desc)
	return true
}

func vIncant(arg ActionArg) bool {
	Printf("The incantation echoes back faintly, but nothing else happens.\n")
	G.Params.InQuotes = false
	G.Params.Continue = NumUndef
	return true
}

func vInflate(arg ActionArg) bool {
	Printf("How can you inflate that?\n")
	return true
}

func vKick(arg ActionArg) bool {
	return hackHack("Kicking the ")
}

func vKiss(arg ActionArg) bool {
	Printf("I'd sooner kiss a pig.\n")
	return true
}

func vKnock(arg ActionArg) bool {
	if G.DirObj.Has(FlgDoor) {
		Printf("Nobody's home.\n")
		return true
	}
	Printf("Why knock on a %s?\n", G.DirObj.Desc)
	return true
}

func vLampOff(arg ActionArg) bool {
	if !G.DirObj.Has(FlgLight) {
		Printf("You can't turn that off.\n")
		return true
	}
	if !G.DirObj.Has(FlgOn) {
		Printf("it is already off.\n")
		return true
	}
	G.DirObj.Take(FlgOn)
	if G.Lit {
		G.Lit = IsLit(G.Here, true)
	}
	Printf("The %s is now off.\n", G.DirObj.Desc)
	if !G.Lit {
		Printf("it is now pitch black.\n")
	}
	return true
}

func vLampOn(arg ActionArg) bool {
	if !G.DirObj.Has(FlgLight) {
		if G.DirObj.Has(FlgBurn) {
			Printf("If you wish to burn the %s, you should say so.\n", G.DirObj.Desc)
			return true
		}
		Printf("You can't turn that on.\n")
		return true
	}
	if G.DirObj.Has(FlgOn) {
		Printf("it is already on.\n")
		return true
	}
	G.DirObj.Give(FlgOn)
	Printf("The %s is now on.\n", G.DirObj.Desc)
	if !G.Lit {
		G.Lit = IsLit(G.Here, true)
		Printf("\n")
		return vLook(ActUnk)
	}
	return true
}

func vLaunch(arg ActionArg) bool {
	if G.DirObj.Has(FlgVeh) {
		Printf("You can't launch that by saying \"launch\"!\n")
		return true
	}
	Printf("That's pretty weird.\n")
	return true
}

func vLeanOn(arg ActionArg) bool {
	Printf("Getting tired?\n")
	return true
}

func vListen(arg ActionArg) bool {
	Printf("The %s makes no sound.\n", G.DirObj.Desc)
	return true
}

func vLock(arg ActionArg) bool {
	Printf("it doesn't seem to work.\n")
	return true
}

func vLower(arg ActionArg) bool {
	return hackHack("Playing in this way with the ")
}

func vLook(arg ActionArg) bool {
	if describeRoom(true) {
		return describeObjects(true)
	}
	return false
}

func vLookBehind(arg ActionArg) bool {
	Printf("There is nothing behind the %s.\n", G.DirObj.Desc)
	return true
}

func vLookInside(arg ActionArg) bool {
	if G.DirObj.Has(FlgDoor) {
		Printf("The %s", G.DirObj.Desc)
		if G.DirObj.Has(FlgOpen) {
			Printf(" is open, but I can't tell what's beyond it.")
		} else {
			Printf(" is closed.")
		}
		Printf("\n")
		return true
	}
	if G.DirObj.Has(FlgCont) {
		if G.DirObj.Has(FlgPerson) {
			Printf("There is nothing special to be seen.\n")
			return true
		}
		if !canSeeInside(G.DirObj) {
			Printf("The %s is closed.\n", G.DirObj.Desc)
			return true
		}
		if G.DirObj.HasChildren() && printCont(G.DirObj, false, 0) {
			return true
		}
		Printf("The %s is empty.\n", G.DirObj.Desc)
		return true
	}
	Printf("You can't look inside a %s.\n", G.DirObj.Desc)
	return true
}

func vLookOn(arg ActionArg) bool {
	if G.DirObj.Has(FlgSurf) {
		Perform(ActionVerb{Norm: "look inside", Orig: "look inside"}, G.DirObj, nil)
		return true
	}
	Printf("Look on a %s???\n", G.DirObj.Desc)
	return true
}

func vLookUnder(arg ActionArg) bool {
	Printf("There is nothing but dust there.\n")
	return true
}

func vExamine(arg ActionArg) bool {
	if len(G.DirObj.Text) > 0 {
		Printf("%s\n", G.DirObj.Text)
		return true
	}
	if G.DirObj.Has(FlgCont) || G.DirObj.Has(FlgDoor) {
		return vLookInside(ActUnk)
	}
	Printf("There's nothing special about the %s.\n", G.DirObj.Desc)
	return true
}

func vMake(arg ActionArg) bool {
	Printf("You can't do that.\n")
	return true
}

func vMelt(arg ActionArg) bool {
	Printf("it's not clear that a %s can be melted.\n", G.DirObj.Desc)
	return true
}

func preMove(arg ActionArg) bool {
	if IsHeld(G.DirObj) {
		Printf("You aren't an accomplished enough juggler.\n")
		return true
	}
	return false
}

func vMove(arg ActionArg) bool {
	if G.DirObj.Has(FlgTake) {
		Printf("Moving the %s reveals nothing.\n", G.DirObj.Desc)
		return true
	}
	Printf("You can't move the %s.\n", G.DirObj.Desc)
	return true
}

func vMumble(arg ActionArg) bool {
	Printf("You'll have to speak up if you expect me to hear you!\n")
	return true
}

func preMung(arg ActionArg) bool {
	if G.IndirObj == nil || !G.IndirObj.Has(FlgWeapon) {
		Printf("Trying to destroy the %s with ", G.DirObj.Desc)
		if G.IndirObj == nil {
			Printf("your bare hands")
		} else {
			Printf("a%s", G.IndirObj.Desc)
		}
		Printf(" is futile.\n")
		return true
	}
	return false
}

func vMung(arg ActionArg) bool {
	if !G.DirObj.Has(FlgPerson) {
		Printf("Nice try.\n")
		return true
	}
	Perform(ActionVerb{Norm: "attack", Orig: "attack"}, G.DirObj, nil)
	return true
}

func mungRoom(rm *Object, str string) {
	rm.Give(FlgKludge)
	rm.LongDesc = str
}

func vOdysseus(arg ActionArg) bool {
	if G.Here != &cyclopsRoom || !cyclops.IsIn(G.Here) || gD().CyclopsFlag {
		Printf("Wasn't he a sailor?\n")
		return true
	}
	QueueInt("iCyclops", false).Run = false
	gD().CyclopsFlag = true
	Printf("The cyclops, hearing the name of his father's deadly nemesis, flees the room by knocking down the wall on the east of the room.\n")
	gD().MagicFlag = true
	cyclops.Take(FlgFight)
	return removeCarefully(&cyclops)
}

func vOil(arg ActionArg) bool {
	Printf("You probably put spinach in your gas tank, too.\n")
	return true
}

func vOpen(arg ActionArg) bool {
	if !G.DirObj.Has(FlgCont) || G.DirObj.GetCapacity() == 0 {
		if G.DirObj.Has(FlgDoor) {
			if G.DirObj.Has(FlgOpen) {
				Printf("it is already open.\n")
				return true
			}
			Printf("The %s opens.\n", G.DirObj.Desc)
			G.DirObj.Give(FlgOpen)
			return true
		}
		Printf("You must tell me how to do that to a %s.\n", G.DirObj.Desc)
		return true
	}
	if G.DirObj.Has(FlgOpen) {
		Printf("it is already open.\n")
		return true
	}
	G.DirObj.Give(FlgOpen)
	G.DirObj.Give(FlgTouch)
	if !G.DirObj.HasChildren() || G.DirObj.Has(FlgTrans) {
		Printf("Opened.\n")
		return true
	}
	if len(G.DirObj.Children) == 1 && !G.DirObj.Children[0].Has(FlgTouch) && len(G.DirObj.Children[0].FirstDesc) > 0 {
		str := G.DirObj.Children[0].FirstDesc
		Printf("The %s opens.\n%s\n", G.DirObj.Desc, str)
		return true
	}
	Printf("Opening the %s reveals ", G.DirObj.Desc)
	printContents(G.DirObj)
	Printf(".\n")
	return true
}

func vOverboard(arg ActionArg) bool {
	locn := G.Winner.Location()
	if G.IndirObj == &teeth {
		if locn.Has(FlgVeh) {
			Printf("Ahoy -- %s overboard!\n", G.IndirObj.Desc)
			return true
		}
		Printf("You're not in anything!\n")
		return true
	}
	if locn.Has(FlgVeh) {
		Perform(ActionVerb{Norm: "throw", Orig: "throw"}, G.DirObj, nil)
		return true
	}
	Printf("Huh?\n")
	return true
}

func vPick(arg ActionArg) bool {
	Printf("You can't pick that.\n")
	return true
}

func vPlay(arg ActionArg) bool {
	if !G.DirObj.Has(FlgPerson) {
		Printf("That's silly!\n")
		return true
	}
	Printf("You become so engrossed in the role of the %s that you kill yourself, just as he might have done!\n", G.DirObj.Desc)
	return jigsUp("", false)
}

func vPlug(arg ActionArg) bool {
	Printf("This has no effect.\n")
	return true
}

func vPourOn(arg ActionArg) bool {
	if G.DirObj == &water {
		removeCarefully(G.DirObj)
		if IsFlaming(G.IndirObj) {
			Printf("The %s is extinguished.\n", G.IndirObj.Desc)
			G.IndirObj.Take(FlgOn)
			G.IndirObj.Take(FlgFlame)
			return true
		}
		Printf("The water spills over the %s, to the floor, and evaporates.\n", G.IndirObj.Desc)
		return true
	}
	if G.DirObj == &putty {
		return Perform(ActionVerb{Norm: "put", Orig: "put"}, &putty, G.IndirObj) == PerfHndld
	}
	Printf("You can't pour that.\n")
	return true
}

func vPray(arg ActionArg) bool {
	if G.Here != &southTemple {
		Printf("If you pray enough, your prayers may be answered.\n")
		return true
	}
	return gotoRoom(&forest1, true)
}

func vPump(arg ActionArg) bool {
	if G.IndirObj != nil && G.IndirObj != &pump {
		Printf("pump it up with a %s?\n", G.IndirObj.Desc)
		return true
	}
	if pump.IsIn(G.Winner) {
		return Perform(ActionVerb{Norm: "inflate", Orig: "inflate"}, G.DirObj, &pump) == PerfHndld
	}
	Printf("it's really not clear how.\n")
	return true
}

func vPush(arg ActionArg) bool {
	return hackHack("Pushing the ")
}

func vPushTo(arg ActionArg) bool {
	Printf("You can't push things to that.\n")
	return true
}

func prePut(arg ActionArg) bool {
	return preGive(arg)
}

func vPut(arg ActionArg) bool {
	if !G.IndirObj.Has(FlgOpen) && !IsOpenable(G.IndirObj) && !G.IndirObj.Has(FlgVeh) {
		Printf("You can't do that.\n")
		return true
	}
	if !G.IndirObj.Has(FlgOpen) {
		Printf("The %s isn't open.\n", G.IndirObj.Desc)
		thisIsIt(G.IndirObj)
		return true
	}
	if G.IndirObj == G.DirObj {
		Printf("How can you do that?\n")
		return true
	}
	if G.DirObj.IsIn(G.IndirObj) {
		Printf("The %s is already in the %s.\n", G.DirObj.Desc, G.IndirObj.Desc)
		return true
	}
	if weight(G.IndirObj)+weight(G.DirObj)-G.IndirObj.GetSize() > G.IndirObj.GetCapacity() {
		Printf("There's no room.\n")
		return true
	}
	if !IsHeld(G.DirObj) && !iTake(true) {
		return true
	}
	G.DirObj.MoveTo(G.IndirObj)
	G.DirObj.Give(FlgTouch)
	scoreObj(G.DirObj)
	Printf("Done.\n")
	return true
}

func vPutBehind(arg ActionArg) bool {
	Printf("That hiding place is too obvious.\n")
	return true
}

func vPutOn(arg ActionArg) bool {
	if G.IndirObj == nil || G.IndirObj == &ground {
		return vDrop(ActUnk)
	}
	if G.IndirObj.Has(FlgSurf) {
		return vPut(ActUnk)
	}
	Printf("There's no good surface on the %s.\n", G.IndirObj.Desc)
	return true
}

func vPutUnder(arg ActionArg) bool {
	Printf("You can't do that.\n")
	return true
}

func vRaise(arg ActionArg) bool {
	return vLower(arg)
}

func vRape(arg ActionArg) bool {
	Printf("What a (ahem!) strange idea.\n")
	return true
}

func preRead(arg ActionArg) bool {
	if !G.Lit {
		Printf("it is impossible to read in the dark.\n")
		return true
	}
	if G.IndirObj != nil && !G.IndirObj.Has(FlgTrans) {
		Printf("How does one look through a %s?\n", G.IndirObj.Desc)
		return true
	}
	return false
}

func vRead(arg ActionArg) bool {
	if !G.DirObj.Has(FlgRead) {
		Printf("How does one read a %s?\n", G.DirObj.Desc)
		return true
	}
	Printf("%s\n", G.DirObj.Text)
	return true
}

func vReadPage(arg ActionArg) bool {
	return vRead(ActUnk)
}

func vRepent(arg ActionArg) bool {
	Printf("it could very well be too late!\n")
	return true
}

func vReply(arg ActionArg) bool {
	Printf("it is hardly likely that the %s is interested.\n", G.DirObj.Desc)
	G.Params.Continue = NumUndef
	G.Params.InQuotes = false
	return true
}

func vRing(arg ActionArg) bool {
	Printf("How, exactly, can you ring that?\n")
	return true
}

func vRub(arg ActionArg) bool {
	return hackHack("Fiddling with the ")
}

func vSay(arg ActionArg) bool {
	if G.Params.Continue == NumUndef {
		Printf("Say what?\n")
		return true
	}
	G.Params.InQuotes = false
	v := findIn(G.Here, FlgPerson)
	if v != nil {
		Printf("You must address the %s directly.\n", v.Desc)
		G.Params.Continue = NumUndef
		return true
	}
	if G.LexRes[G.Params.Continue].Norm == "hello" {
		G.Params.Continue = NumUndef
		Printf("Talking to yourself is a sign of impending mental collapse.\n")
		return true
	}
	return false
}

func findIn(where *Object, what Flags) *Object {
	if !where.HasChildren() {
		return nil
	}
	for w := 0; w < len(where.Children); w++ {
		if where.Children[w].Has(what) && where.Children[w] != &adventurer {
			return where.Children[w]
		}
	}
	return nil
}

func vSearch(arg ActionArg) bool {
	Printf("You find nothing unusual.\n")
	return true
}

func vSend(arg ActionArg) bool {
	if G.DirObj.Has(FlgPerson) {
		Printf("Why would you send for the %s?\n", G.DirObj.Desc)
		return true
	}
	Printf("That doesn't make sends.\n")
	return true
}

func preSGive(arg ActionArg) bool {
	Perform(ActionVerb{Norm: "give", Orig: "give"}, G.IndirObj, G.DirObj)
	return true
}

func vSGive(arg ActionArg) bool {
	Printf("Foo!\n")
	return true
}

func vShake(arg ActionArg) bool {
	if G.DirObj.Has(FlgPerson) {
		Printf("This seems to have no effect.\n")
		return true
	}
	if !G.DirObj.Has(FlgTake) {
		Printf("You can't take it; thus, you can't shake it!\n")
		return true
	}
	if !G.DirObj.Has(FlgCont) {
		Printf("Shaken.\n")
		return true
	}
	if G.DirObj.Has(FlgOpen) {
		if !G.DirObj.HasChildren() {
			Printf("Shaken.\n")
			return true
		}
		shakeLoop()
		Printf("The contents of the %s spill ", G.DirObj.Desc)
		if !G.Here.Has(FlgLand) {
			Printf("out and disappears")
		} else {
			Printf("to the ground")
		}
		Printf(".\n")
		return true
	}
	if G.DirObj.HasChildren() {
		Printf("it sounds like there is something inside the %s.\n", G.DirObj.Desc)
		return true
	}
	Printf("The %s sounds empty.\n", G.DirObj.Desc)
	return true
}

func shakeLoop() {
	if !G.DirObj.HasChildren() {
		return
	}
	x := G.DirObj.Children[0]
	x.Give(FlgTouch)
	mv := G.Here
	if G.Here == &upATree {
		mv = &path
	} else if !G.Here.Has(FlgLand) {
		mv = &pseudoObject
	}
	x.MoveTo(mv)
}

func vSkip(arg ActionArg) bool {
	Printf("%s\n", PickOne(wheeee))
	return true
}

func vSmell(arg ActionArg) bool {
	Printf("it smells like a %s.\n", G.DirObj.Desc)
	return true
}

func vSpin(arg ActionArg) bool {
	Printf("You can't spin that!\n")
	return true
}

func vSpray(arg ActionArg) bool {
	return vSqueeze(arg)
}

func vSqueeze(arg ActionArg) bool {
	if !G.DirObj.Has(FlgPerson) {
		Printf("How singularly useless.\n")
		return true
	}
	Printf("The %s does not understand this.\n", G.DirObj.Desc)
	return true
}

func vSSpray(arg ActionArg) bool {
	return Perform(ActionVerb{Norm: "spray", Orig: "spray"}, G.IndirObj, G.DirObj) == PerfHndld
}

func vStab(arg ActionArg) bool {
	w := findWeapon(G.Winner)
	if w == nil {
		Printf("No doubt you propose to stab the %s with your pinky?\n", G.DirObj.Desc)
		return true
	}
	Perform(ActionVerb{Norm: "attack", Orig: "attack"}, G.DirObj, w)
	return true
}

func vStrike(arg ActionArg) bool {
	if G.DirObj.Has(FlgPerson) {
		Printf("Since you aren't versed in hand-to-hand combat, you'd better attack the %s with a weapon.\n", G.DirObj.Desc)
		return true
	}
	Perform(ActionVerb{Norm: "lamp on", Orig: "lamp on"}, G.DirObj, nil)
	return true
}

func vSwing(arg ActionArg) bool {
	if G.IndirObj == nil {
		Printf("Whoosh!\n")
		return true
	}
	if Perform(ActionVerb{Norm: "attack", Orig: "attack"}, G.IndirObj, G.DirObj) == PerfHndld {
		return true
	}
	return false
}

func preTake(arg ActionArg) bool {
	if G.DirObj == G.Winner {
		if G.DirObj.Has(FlgWear) {
			Printf("You are already wearing it.\n")
			return true
		}
		Printf("You already have that!\n")
		return true
	}
	lcn := G.DirObj.Location()
	if lcn.Has(FlgCont) && !lcn.Has(FlgOpen) {
		Printf("You can't reach something that's inside a closed container.\n")
		return true
	}
	if G.IndirObj != nil {
		if G.IndirObj == &ground {
			G.IndirObj = nil
			return false
		}
		if G.IndirObj != G.DirObj.Location() {
			Printf("The %s isn't in the %s.\n", G.DirObj.Desc, G.IndirObj.Desc)
			return true
		}
		G.IndirObj = nil
		return false
	}
	if G.DirObj == G.Winner.Location() {
		Printf("You're inside of it!\n")
		return true
	}
	return false
}

func vTake(arg ActionArg) bool {
	if iTake(true) {
		if G.DirObj.Has(FlgWear) {
			Printf("You are now wearing the %s.\n", G.DirObj.Desc)
			return true
		}
		Printf("Taken.\n")
		return true
	}
	return false
}

func vTell(arg ActionArg) bool {
	if !G.DirObj.Has(FlgPerson) {
		Printf("You can't talk to the %s!\n", G.DirObj.Desc)
		G.Params.InQuotes = false
		G.Params.Continue = NumUndef
		return RFatal()
	}
	if G.Params.Continue != NumUndef {
		G.Winner = G.DirObj
		G.Here = G.Winner.Location()
		return true
	}
	Printf("The %s pauses for a moment, perhaps thinking that you should reread the manual.\n", G.DirObj.Desc)
	return true
}

func vThrow(arg ActionArg) bool {
	if !iDrop() {
		Printf("Huh?\n")
		return true
	}
	if G.IndirObj == &me {
		Printf("A terrific throw! The ")
		G.Winner = G.Player
		return jigsUp(" hits you squarely in the head. Normally, this wouldn't do much damage, but by incredible mischance, you fall over backwards trying to duck, and break your neck, justice being swift and merciful in the Great Underground Empire.", false)
	}
	if G.IndirObj != nil && G.IndirObj.Has(FlgPerson) {
		Printf("The %s ducks as the %s flies by and crashes to the ground.\n", G.IndirObj.Desc, G.DirObj.Desc)
		return true
	}
	Printf("Thrown.\n")
	return true
}

func vThrowOff(arg ActionArg) bool {
	Printf("You can't throw anything off of that!\n")
	return true
}

func vTie(arg ActionArg) bool {
	if G.IndirObj == G.Winner {
		Printf("You can't tie anything to yourself.\n")
		return true
	}
	Printf("You can't tie the %s to that.\n", G.DirObj.Desc)
	return true
}

func vTieUp(arg ActionArg) bool {
	Printf("You could certainly never tie it with that!\n")
	return true
}

func vTreasure(arg ActionArg) bool {
	if G.Here == &northTemple {
		return gotoRoom(&treasureRoom, true)
	}
	if G.Here == &treasureRoom {
		return gotoRoom(&northTemple, true)
	}
	Printf("Nothing happens.\n")
	return true
}

func preTurn(arg ActionArg) bool {
	if (G.IndirObj == nil || G.IndirObj == &rooms) && G.DirObj != &book {
		Printf("Your bare hands don't appear to be enough.\n")
		return true
	}
	if !G.DirObj.Has(FlgTurn) {
		Printf("You can't turn that!\n")
		return true
	}
	return false
}

func vTurn(arg ActionArg) bool {
	Printf("This has no effect.\n")
	return true
}

func vUnlock(arg ActionArg) bool {
	return vLock(arg)
}

func vUntie(arg ActionArg) bool {
	Printf("This cannot be tied, so it cannot be untied!\n")
	return true
}

func vWave(arg ActionArg) bool {
	return hackHack("Waving the ")
}

func vWear(arg ActionArg) bool {
	if !G.DirObj.Has(FlgWear) {
		Printf("You can't wear the %s.\n", G.DirObj.Desc)
		return true
	}
	Perform(ActionVerb{Norm: "take", Orig: "take"}, G.DirObj, nil)
	return true
}

func vWind(arg ActionArg) bool {
	Printf("You cannot wind up a %s.\n", G.DirObj.Desc)
	return true
}

func vYell(arg ActionArg) bool {
	Printf("Aaaarrrrgggghhhh!\n")
	return true
}

func removeCarefully(obj *Object) bool {
	if obj == G.Params.ItObj {
		G.Params.ItObj = nil
	}
	oLit := G.Lit
	obj.Remove()
	G.Lit = IsLit(G.Here, true)
	if oLit && oLit != G.Lit {
		Printf("You are left in the dark...\n")
	}
	return true
}

func printContents(obj *Object) bool {
	if !obj.HasChildren() {
		return false
	}
	var itObj *Object
	twoIs := false
	for n := 0; n < len(obj.Children); n++ {
		if n != 0 {
			Printf(", ")
			if n+1 >= len(obj.Children) {
				Printf("and ")
			}
		}
		Printf("a %s", obj.Children[n].Desc)
		if itObj == nil && !twoIs {
			itObj = obj.Children[n]
		} else {
			twoIs = true
			itObj = nil
		}
	}
	if itObj != nil && !twoIs {
		thisIsIt(itObj)
	}
	return true
}

func hackHack(str string) bool {
	if G.DirObj.IsIn(&globalObjects) && (G.ActVerb.Norm == "wave" || G.ActVerb.Norm == "raise" || G.ActVerb.Norm == "lower") {
		Printf("The %s isn't here!\n", G.DirObj.Desc)
		return true
	}
	Printf("%s%s%s\n", str, G.DirObj.Desc, PickOne(hohum))
	return true
}

func iDrop() bool {
	if !G.DirObj.IsIn(G.Winner) && !G.DirObj.Location().IsIn(G.Winner) {
		Printf("You're not carrying the %s.\n", G.DirObj.Desc)
		return false
	}
	if !G.DirObj.IsIn(G.Winner) && !G.DirObj.Location().Has(FlgOpen) {
		Printf("The %s is closed.\n", G.DirObj.Desc)
		return false
	}
	G.DirObj.MoveTo(G.Winner.Location())
	return true
}

func describeRoom(isLook bool) bool {
	v := G.Verbose
	if isLook {
		v = isLook
	}
	if !G.Lit {
		Printf("it is pitch black.")
		if !gD().IsSprayed {
			Printf(" You are likely to be eaten by a grue.")
		}
		Printf("\n")
		return false
	}
	if !G.Here.Has(FlgTouch) {
		G.Here.Give(FlgTouch)
		v = true
	}
	if G.Here.Has(FlgMaze) {
		G.Here.Take(FlgTouch)
	}
	if G.Here.IsIn(&rooms) {
		Printf("%s", G.Here.Desc)
		if av := G.Winner.Location(); av.Has(FlgVeh) {
			Printf(", in the %s", av.Desc)
		}
		Printf("\n")
	}
	if !isLook && G.SuperBrief {
		return true
	}
	av := G.Winner.Location()
	if v && G.Here.Action != nil && G.Here.Action(ActLook) {
		return true
	}
	if v && len(G.Here.LongDesc) != 0 {
		Printf("%s\n", G.Here.LongDesc)
	} else if G.Here.Action != nil {
		G.Here.Action(ActFlash)
	}
	if G.Here != av && av.Has(FlgVeh) && av.Action != nil {
		av.Action(ActLook)
	}
	return true
}

func describeObjects(v bool) bool {
	if !G.Lit {
		Printf("Only bats can see in the dark. And you're not one.\n")
		return true
	}
	if !G.Here.HasChildren() {
		return false
	}
	if !v {
		v = G.Verbose
	}
	return printCont(G.Here, v, -1)
}

func printCont(obj *Object, v bool, lvl int) bool {
	if !obj.HasChildren() {
		return true
	}
	var av *Object
	if G.Winner.Location().Has(FlgVeh) {
		av = G.Winner.Location()
	}
	isInv := false
	isPv := false
	shit := true
	isFirst := true
	if G.Winner == obj || G.Winner == obj.Location() {
		isInv = true
	} else {
		for _, child := range obj.Children {
			if child == av {
				isPv = true
			} else if G.Winner == child {
				continue
			} else if !child.Has(FlgInvis) && !child.Has(FlgTouch) && len(child.FirstDesc) > 0 {
				if !child.Has(FlgNoDesc) {
					Printf("%s\n", child.FirstDesc)
					shit = false
				}
				if canSeeInside(child) && child.Location().DescFcn == nil && child.HasChildren() {
					if printCont(child, v, 0) {
						isFirst = false
					}
				}
			}
		}
	}
	for _, child := range obj.Children {
		if child == av || child == &adventurer {
			continue
		}
		if child.Has(FlgInvis) || (!isInv && !child.Has(FlgTouch) && len(child.FirstDesc) > 0) {
			continue
		}
		if !child.Has(FlgNoDesc) {
			if isFirst {
				if firster(obj, lvl) && lvl < 0 {
					lvl = 0
				}
				lvl++
				isFirst = false
			}
			if lvl < 0 {
				lvl = 0
			}
			describeObject(child, v, lvl)
		} else if child.HasChildren() && canSeeInside(child) {
			lvl++
			printCont(child, v, lvl)
			lvl--
		}
	}
	if isPv && av != nil && av.HasChildren() {
		lvl++
		printCont(av, v, lvl)
	}
	if isFirst && shit {
		return false
	}
	return true
}

func firster(obj *Object, lvl int) bool {
	if obj == &trophyCase {
		Printf("Your collection of treasures consists of:\n")
		return true
	}
	if obj == G.Winner {
		Printf("You are carrying:\n")
		return true
	}
	if obj.IsIn(&rooms) {
		return false
	}
	if lvl > 0 {
		Printf("%s", indents[lvl])
	}
	if obj.Has(FlgSurf) {
		Printf("Sitting on the %s is: \n", obj.Desc)
		return true
	}
	if obj.Has(FlgPerson) {
		Printf("The %s is holding: \n", obj.Desc)
		return true
	}
	Printf("The %s contains:\n", obj.Desc)
	return true
}

func describeObject(obj *Object, v bool, lvl int) bool {
	gD().DescObj = obj
	if lvl == 0 && obj.DescFcn != nil && obj.DescFcn(ActObjDesc) {
		return true
	}
	if lvl == 0 && ((!obj.Has(FlgTouch) && len(obj.FirstDesc) > 0) || len(obj.LongDesc) > 0) {
		if !obj.Has(FlgTouch) && len(obj.FirstDesc) > 0 {
			Printf("%s", obj.FirstDesc)
		} else {
			Printf("%s", obj.LongDesc)
		}
	} else if lvl == 0 {
		Printf("There is a %s here", obj.Desc)
		if obj.Has(FlgOn) {
			Printf(" (providing light)")
		}
		Printf(".")
	} else {
		Printf("%sA %s", indents[lvl], obj.Desc)
		if obj.Has(FlgOn) {
			Printf(" (providing light)")
		} else if obj.Has(FlgWear) && obj.IsIn(G.Winner) {
			Printf(" (being worn)")
		}
	}
	if av := G.Winner.Location(); lvl == 0 && av != nil && av.Has(FlgVeh) {
		Printf(" (outside the %s)", av.Desc)
	}
	Printf("\n")
	if canSeeInside(obj) && obj.HasChildren() {
		return printCont(obj, v, lvl)
	}
	return true
}

func canSeeInside(obj *Object) bool {
	if obj.Has(FlgInvis) {
		return false
	}
	if obj.Has(FlgTrans) || obj.Has(FlgOpen) {
		return true
	}
	return false
}

func thisIsIt(obj *Object) {
	G.Params.ItObj = obj
}

func scoreUpd(num int) bool {
	G.BaseScore += num
	G.Score += num
	if G.Score == 350 && !gD().WonGame {
		gD().WonGame = true
		mapObj.Take(FlgInvis)
		westOfHouse.Take(FlgTouch)
		Printf("An almost inaudible voice whispers in your ear, \"Look to your treasures for the final secret.\"\n")
	}
	return true
}

func scoreObj(obj *Object) {
	if obj.GetValue() <= 0 {
		return
	}
	scoreUpd(obj.GetValue())
	obj.SetValue(0)
}

func cCount(obj *Object) int {
	cnt := 0
	for _, child := range obj.Children {
		if !child.Has(FlgWear) {
			cnt++
		}
	}
	return cnt
}

func weight(obj *Object) int {
	wt := 0
	for _, child := range obj.Children {
		if obj == G.Player && child.Has(FlgWear) {
			wt++
		} else {
			wt += weight(child)
		}
	}
	return wt + obj.GetSize()
}

func iTake(vb bool) bool {
	if gD().Dead {
		if vb {
			Printf("Your hand passes through its object.\n")
		}
		return false
	}
	if G.DirObj == nil {
		return false
	}
	if !G.DirObj.Has(FlgTake) {
		if vb {
			Printf("%s\n", PickOne(yuks))
		}
		return false
	}
	// ZIL: <FSET? <LOC ,PRSO> ,CONTBIT> / <NOT <FSET? <LOC ,PRSO> ,OPENBIT>>
	// Prevent taking objects from inside a closed container.
	loc := G.DirObj.Location()
	if loc != nil && loc.Has(FlgCont) && !loc.Has(FlgOpen) {
		return false
	}
	if !G.DirObj.Location().IsIn(G.Winner) && weight(G.DirObj)+weight(G.Winner) > gD().LoadAllowed {
		if vb {
			Printf("Your load is too heavy")
			if gD().LoadAllowed < gD().LoadMax {
				Printf(", especially in light of your condition.")
			} else {
				Printf(".")
			}
			Printf("\n")
		}
		G.PerformFatal = true
		return false
	}
	cnt := cCount(G.Winner)
	if G.ActVerb.Norm == "tell" && cnt > gD().FumbleNumber && Prob(cnt*gD().FumbleProb, false) {
		Printf("You're holding too many things already!\n")
		return false
	}
	G.DirObj.MoveTo(G.Winner)
	G.DirObj.Take(FlgNoDesc)
	G.DirObj.Give(FlgTouch)
	scoreObj(G.DirObj)
	return true
}
