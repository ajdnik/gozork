package zork


type RndSelect struct {
	Unselected []string
	Selected   []string
}

var (
	// HelloSailor counts occurences of 'hello, sailor'
	// IsSprayed is a flag indicating if the player is wearing Grue repellent
	Version   = [24]byte{0, 0, 0, 0, 119, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 56, 56, 48, 52, 50, 57}
	// DescObj stores the last object which was described
	Indents = [6]string{
		"",
		"  ",
		"    ",
		"      ",
		"        ",
		"          ",
	}
	Yuks = RndSelect{
		Unselected: []string{
			"A valiant attempt.",
			"You can't be serious.",
			"An interesting idea...",
			"What a concept!",
		},
	}
	Hellos = RndSelect{
		Unselected: []string{
			"Hello.",
			"Good day.",
			"Nice weather we've been having lately.",
			"Goodbye.",
		},
	}
	JumpLoss = RndSelect{
		Unselected: []string{
			"You should have looked before you leaped.",
			"In the movies, your life would be passing before your eyes.",
			"Geronimo...",
		},
	}
	Wheeee = RndSelect{
		Unselected: []string{
			"Very good. Now you can go to the second grade.",
			"Are you enjoying yourself?",
			"Wheeeeeeeeee!!!!!",
			"Do you expect me to applaud?",
		},
	}
	Hohum = RndSelect{
		Unselected: []string{
			" doesn't seem to work.",
			" isn't notably helpful.",
			" has no effect.",
		},
	}
)

func VAlarm(arg ActArg) bool {
	if !G.DirObj.Has(FlgPerson) {
		Print("The ", NoNewline)
		PrintObject(G.DirObj)
		Print(" isn't sleeping.", Newline)
		return true
	}
	if G.DirObj.Strength <= 0 {
		Print("He's wide awake, or haven't you noticed...", Newline)
		return true
	}
	Print("The ", NoNewline)
	PrintObject(G.DirObj)
	Print(" is rudely awakened.", Newline)
	return Awaken(G.DirObj)
}

func VAnswer(arg ActArg) bool {
	Print("Nobody seems to be awaiting your answer.", Newline)
	G.Params.Continue = NumUndef
	G.Params.InQuotes = false
	return true
}

func VAttack(arg ActArg) bool {
	if !G.DirObj.Has(FlgPerson) {
		Print("I've known strange people, but fighting a ", NoNewline)
		PrintObject(G.DirObj)
		Print("?", Newline)
		return true
	}
	if G.IndirObj == nil || G.IndirObj == &Hands {
		Print("Trying to attack a ", NoNewline)
		PrintObject(G.DirObj)
		Print(" with your bare hands is suicidal.", Newline)
		return true
	}
	if !G.IndirObj.IsIn(G.Winner) {
		Print("You aren't even holding the ", NoNewline)
		PrintObject(G.IndirObj)
		Print(".", Newline)
		return true
	}
	if !G.IndirObj.Has(FlgWeapon) {
		Print("Trying to attack the ", NoNewline)
		PrintObject(G.DirObj)
		Print(" with a ", NoNewline)
		PrintObject(G.IndirObj)
		Print(" is suicidal.", Newline)
		return true
	}
	return HeroBlow()
}

func VBack(arg ActArg) bool {
	Print("Sorry, my memory is poor. Please give a direction.", Newline)
	return true
}

func VBlast(arg ActArg) bool {
	Print("You can't blast anything by using words.", Newline)
	return true
}

func VBreathe(arg ActArg) bool {
	ret := Perform(ActionVerb{Norm: "inflate", Orig: "inflate"}, G.DirObj, &Lungs)
	if ret == PerfFatal {
		return RFatal()
	}
	return ret == PerfHndld
}

func VBrush(arg ActArg) bool {
	Print("If you wish, but heaven only knows why.", Newline)
	return true
}

func TellNoPrsi() bool {
	Print("You didn't say with what!", Newline)
	return true
}

func PreBurn(arg ActArg) bool {
	if G.IndirObj == nil {
		Print("You didn't say with what!", Newline)
		return true
	}
	if IsFlaming(G.IndirObj) {
		return false
	}
	Print("With a ", NoNewline)
	PrintObject(G.IndirObj)
	Print("??!?", Newline)
	return true
}

func VBurn(arg ActArg) bool {
	if !G.DirObj.Has(FlgBurn) {
		Print("You can't burn a ", NoNewline)
		PrintObject(G.DirObj)
		Print(".", Newline)
		return true
	}
	if !G.DirObj.IsIn(G.Winner) && !G.Winner.IsIn(G.DirObj) {
		RemoveCarefully(G.DirObj)
		Print("The ", NoNewline)
		PrintObject(G.DirObj)
		Print(" catches fire and is consumed.", Newline)
		return true
	}
	RemoveCarefully(G.DirObj)
	Print("The ", NoNewline)
	PrintObject(G.DirObj)
	Print(" catches fire. Unfortunately, you were ", NoNewline)
	if G.Winner.IsIn(G.DirObj) {
		Print("in", NoNewline)
	} else {
		Print("holding", NoNewline)
	}
	return JigsUp(" it at the time.", false)
}

func VChomp(arg ActArg) bool {
	Print("Preposterous!", Newline)
	return true
}

func VClose(arg ActArg) bool {
	if !G.DirObj.Has(FlgCont) && !G.DirObj.Has(FlgDoor) {
		Print("You must tell me how to do that to a ", NoNewline)
		PrintObject(G.DirObj)
		Print(".", Newline)
		return true
	}
	if !G.DirObj.Has(FlgSurf) && G.DirObj.Capacity != 0 {
		if G.DirObj.Has(FlgOpen) {
			G.DirObj.Take(FlgOpen)
			Print("Closed.", Newline)
			if G.Lit {
				G.Lit = IsLit(G.Here, true)
				if !G.Lit {
					Print("It is now pitch black.", Newline)
				}
			}
			return true
		}
		Print("It is already closed.", Newline)
		return true
	}
	if G.DirObj.Has(FlgDoor) {
		if G.DirObj.Has(FlgOpen) {
			G.DirObj.Take(FlgOpen)
			Print("The ", NoNewline)
			PrintObject(G.DirObj)
			Print(" is now closed.", Newline)
			return true
		}
		Print("It is already closed.", Newline)
		return true
	}
	Print("You cannot close that.", Newline)
	return true
}

func VCommand(arg ActArg) bool {
	if G.DirObj.Has(FlgPerson) {
		Print("The ", NoNewline)
		PrintObject(G.DirObj)
		Print(" pays no attention.", Newline)
		return true
	}
	Print("You cannot talk to that!", Newline)
	return true
}

func VCount(arg ActArg) bool {
	if G.DirObj == &Blessings {
		Print("Well, for one, you are playing Zork...", Newline)
		return true
	}
	Print("You have lost your mind.", Newline)
	return true
}

func VCurses(arg ActArg) bool {
	if G.DirObj == nil {
		Print("Such language in a high-class establishment like this!", Newline)
		return true
	}
	if G.DirObj.Has(FlgPerson) {
		Print("Insults of this nature won't help you.", Newline)
		return true
	}
	Print("What a loony!", Newline)
	return true
}

func VCut(arg ActArg) bool {
	if G.DirObj.Has(FlgPerson) {
		ret := Perform(ActionVerb{Norm: "attack", Orig: "attack"}, G.DirObj, G.IndirObj)
		if ret == PerfFatal {
			return RFatal()
		}
		return ret == PerfHndld
	}
	if G.DirObj.Has(FlgBurn) && G.IndirObj.Has(FlgWeapon) {
		if G.Winner.IsIn(G.DirObj) {
			Print("Not a bright idea, especially since you're in it.", Newline)
			return true
		}
		RemoveCarefully(G.DirObj)
		Print("Your skillful ", NoNewline)
		PrintObject(G.IndirObj)
		Print("smanship slices the ", NoNewline)
		PrintObject(G.DirObj)
		Print(" into innumerable slivers which blow away.", Newline)
		return true
	}
	if !G.IndirObj.Has(FlgWeapon) {
		Print("The \"cutting edge\" of a ", NoNewline)
		PrintObject(G.IndirObj)
		Print(" is hardly adequate.", Newline)
		return true
	}
	Print("Strange concept, cutting the ", NoNewline)
	PrintObject(G.DirObj)
	Print("....", Newline)
	return true
}

func VDeflate(arg ActArg) bool {
	Print("Come on, now!", Newline)
	return true
}

func VDig(arg ActArg) bool {
	if G.IndirObj == nil {
		G.IndirObj = &Hands
	}
	if G.IndirObj == &Shovel {
		Print("There's no reason to be digging here.", Newline)
		return true
	}
	if G.IndirObj.Has(FlgTool) {
		Print("Digging with the ", NoNewline)
		PrintObject(G.IndirObj)
		Print(" is slow and tedious.", Newline)
		return true
	}
	Print("Digging with a ", NoNewline)
	PrintObject(G.IndirObj)
	Print(" is silly.", Newline)
	return true
}

func VDisenchant(arg ActArg) bool {
	Print("Nothing happens.", Newline)
	return true
}

func VDrink(arg ActArg) bool {
	return VEat(ActUnk)
}

func VDrinkFrom(act ActArg) bool {
	Print("How peculiar!", Newline)
	return true
}

func PreDrop(arg ActArg) bool {
	if G.DirObj == G.Winner.Location() {
		Perform(ActionVerb{Norm: "disembark", Orig: "disembark"}, G.DirObj, nil)
		return true
	}
	return false
}

func VDrop(arg ActArg) bool {
	if IDrop() {
		Print("Dropped.", Newline)
		return true
	}
	return false
}

func VEat(arg ActArg) bool {
	isEat := G.DirObj.Has(FlgFood)
	if isEat {
		if !G.DirObj.IsIn(G.Winner) && !G.DirObj.Location().IsIn(G.Winner) {
			Print("You're not holding that.", Newline)
			return true
		}
		if G.ActVerb.Norm == "drink" {
			Print("How can you drink that?", Newline)
			return true
		}
		Print("Thank you very much. It really hit the spot.", Newline)
		RemoveCarefully(G.DirObj)
		return true
	}
	isDrink := G.DirObj.Has(FlgDrink)
	if isDrink {
		nobj := G.DirObj.Location()
		if G.DirObj.IsIn(&GlobalObjects) || IsInGlobal(&GlobalWater, G.Here) || G.DirObj == &PseudoObject {
			return HitSpot()
		}
		if nobj == nil || !IsAccessible(nobj) {
			Print("There isn't any water here.", Newline)
			return true
		}
		if IsAccessible(nobj) && !nobj.IsIn(G.Winner) {
			Print("You have to be holding the ", NoNewline)
			PrintObject(nobj)
			Print(" first.", Newline)
			return true
		}
		if !nobj.Has(FlgOpen) {
			Print("You'll have to open the ", NoNewline)
			PrintObject(nobj)
			Print(" first.", Newline)
			return true
		}
		return HitSpot()
	}
	if !isEat && !isDrink {
		Print("I don't think that the ", NoNewline)
		PrintObject(G.DirObj)
		Print(" would agree with you.", Newline)
		return true
	}
	return false
}

func HitSpot() bool {
	if G.DirObj == &Water && !IsInGlobal(&GlobalWater, G.Here) {
		RemoveCarefully(G.DirObj)
	}
	Print("Thank you very much. I was rather thirsty (from all this talking, probably).", Newline)
	return true
}

func VEcho(arg ActArg) bool {
	if len(G.LexRes) <= 0 {
		Print("echo echo ...", Newline)
		return true
	}
	wrd := G.LexRes[len(G.LexRes)-1]
	Print(wrd.Orig+" "+wrd.Orig+" ...", Newline)
	return true
}

func VEnchant(arg ActArg) bool {
	return VDisenchant(ActUnk)
}

func VExcorcise(arg ActArg) bool {
	Print("What a bizarre concept!", Newline)
	return true
}

func PreFill(arg ActArg) bool {
	if G.IndirObj == nil {
		if IsInGlobal(&GlobalWater, G.Here) {
			Perform(ActionVerb{Norm: "fill", Orig: "fill"}, G.DirObj, &GlobalWater)
			return true
		}
		if Water.IsIn(G.Winner.Location()) {
			Perform(ActionVerb{Norm: "fill", Orig: "fill"}, G.DirObj, &Water)
			return true
		}
		Print("There is nothing to fill it with.", Newline)
		return true
	}
	if G.IndirObj == &Water {
		return false
	}
	if G.IndirObj != &GlobalWater {
		Perform(ActionVerb{Norm: "put", Orig: "put"}, G.IndirObj, G.DirObj)
		return true
	}
	return false
}

func VFill(arg ActArg) bool {
	if G.IndirObj != nil {
		Print("You may know how to do that, but I don't.", Newline)
		return true
	}
	if IsInGlobal(&GlobalWater, G.Here) {
		Perform(ActionVerb{Norm: "fill", Orig: "fill"}, G.DirObj, &GlobalWater)
		return true
	}
	if Water.IsIn(G.Winner.Location()) {
		Perform(ActionVerb{Norm: "fill", Orig: "fill"}, G.DirObj, &Water)
		return true
	}
	Print("There is nothing to fill it with.", Newline)
	return true
}

func VFirstLook(arg ActArg) bool {
	if DescribeRoom(false) {
		if !G.SuperBrief {
			return DescribeObjects(false)
		}
	}
	return false
}

func VFind(arg ActArg) bool {
	if G.DirObj == &Hands || G.DirObj == &Lungs {
		Print("Within six feet of your head, assuming you haven't left that somewhere.", Newline)
		return true
	}
	if G.DirObj == &Me {
		Print("You're around here somewhere...", Newline)
		return true
	}
	l := G.DirObj.Location()
	if l == &GlobalObjects {
		Print("You find it.", Newline)
		return true
	}
	if G.DirObj.IsIn(G.Winner) {
		Print("You have it.", Newline)
		return true
	}
	if G.DirObj.IsIn(G.Here) || IsInGlobal(G.DirObj, G.Here) || G.DirObj == &PseudoObject {
		Print("It's right here.", Newline)
		return true
	}
	if l.Has(FlgPerson) {
		Print("The ", NoNewline)
		PrintObject(l)
		Print(" has it.", Newline)
		return true
	}
	if l.Has(FlgSurf) {
		Print("It's on the ", NoNewline)
		PrintObject(l)
		Print(".", Newline)
		return true
	}
	if l.Has(FlgCont) {
		Print("It's in the ", NoNewline)
		PrintObject(l)
		Print(".", Newline)
		return true
	}
	Print("Beats me.", Newline)
	return true
}

func VFrobozz(arg ActArg) bool {
	Print("The FROBOZZ Corporation created, owns, and operates this dungeon.", Newline)
	return true
}

func PreGive(arg ActArg) bool {
	if !IsHeld(G.DirObj) {
		Print("That's easy for you to say since you don't even have the ", NoNewline)
		PrintObject(G.DirObj)
		Print(".", Newline)
		return true
	}
	return false
}

func VGive(arg ActArg) bool {
	if !G.DirObj.Has(FlgPerson) {
		Print("You can't give a ", NoNewline)
		PrintObject(G.DirObj)
		Print(" to a ", NoNewline)
		PrintObject(G.IndirObj)
		Print("!", Newline)
		return true
	}
	Print("The ", NoNewline)
	PrintObject(G.IndirObj)
	Print(" refuses it politely.", Newline)
	return true
}

func VHatch(arg ActArg) bool {
	Print("Bizarre!", Newline)
	return true
}

func VHello(arg ActArg) bool {
	if G.DirObj == nil {
		Print(PickOne(Hellos), Newline)
		return true
	}
	if G.DirObj.Has(FlgPerson) {
		Print("The ", NoNewline)
		PrintObject(G.DirObj)
		Print(" bows his head to you in greeting.", Newline)
		return true
	}
	Print("It's a well known fact that only schizophrenics say \"Hello\" to a ", NoNewline)
	PrintObject(G.DirObj)
	Print(".", Newline)
	return true
}

func VIncant(arg ActArg) bool {
	Print("The incantation echoes back faintly, but nothing else happens.", Newline)
	G.Params.InQuotes = false
	G.Params.Continue = NumUndef
	return true
}

func VInflate(arg ActArg) bool {
	Print("How can you inflate that?", Newline)
	return true
}

func VKick(arg ActArg) bool {
	return HackHack("Kicking the ")
}

func VKiss(arg ActArg) bool {
	Print("I'd sooner kiss a pig.", Newline)
	return true
}

func VKnock(arg ActArg) bool {
	if G.DirObj.Has(FlgDoor) {
		Print("Nobody's home.", Newline)
		return true
	}
	Print("Why knock on a ", NoNewline)
	PrintObject(G.DirObj)
	Print("?", Newline)
	return true
}

func VLampOff(arg ActArg) bool {
	if !G.DirObj.Has(FlgLight) {
		Print("You can't turn that off.", Newline)
		return true
	}
	if !G.DirObj.Has(FlgOn) {
		Print("It is already off.", Newline)
		return true
	}
	G.DirObj.Take(FlgOn)
	if G.Lit {
		G.Lit = IsLit(G.Here, true)
	}
	Print("The ", NoNewline)
	PrintObject(G.DirObj)
	Print(" is now off.", Newline)
	if !G.Lit {
		Print("It is now pitch black.", Newline)
	}
	return true
}

func VLampOn(arg ActArg) bool {
	if !G.DirObj.Has(FlgLight) {
		if G.DirObj.Has(FlgBurn) {
			Print("If you wish to burn the ", NoNewline)
			PrintObject(G.DirObj)
			Print(", you should say so.", Newline)
			return true
		}
		Print("You can't turn that on.", Newline)
		return true
	}
	if G.DirObj.Has(FlgOn) {
		Print("It is already on.", Newline)
		return true
	}
	G.DirObj.Give(FlgOn)
	Print("The ", NoNewline)
	PrintObject(G.DirObj)
	Print(" is now on.", Newline)
	if !G.Lit {
		G.Lit = IsLit(G.Here, true)
		NewLine()
		return VLook(ActUnk)
	}
	return true
}

func VLaunch(arg ActArg) bool {
	if G.DirObj.Has(FlgVeh) {
		Print("You can't launch that by saying \"launch\"!", Newline)
		return true
	}
	Print("That's pretty weird.", Newline)
	return true
}

func VLeanOn(arg ActArg) bool {
	Print("Getting tired?", Newline)
	return true
}

func VListen(arg ActArg) bool {
	Print("The ", NoNewline)
	PrintObject(G.DirObj)
	Print(" makes no sound.", Newline)
	return true
}

func VLock(arg ActArg) bool {
	Print("It doesn't seem to work.", Newline)
	return true
}

func VLower(arg ActArg) bool {
	return HackHack("Playing in this way with the ")
}

func VLook(arg ActArg) bool {
	if DescribeRoom(true) {
		return DescribeObjects(true)
	}
	return false
}

func VLookBehind(arg ActArg) bool {
	Print("There is nothing behind the ", NoNewline)
	PrintObject(G.DirObj)
	Print(".", Newline)
	return true
}

func VLookInside(arg ActArg) bool {
	if G.DirObj.Has(FlgDoor) {
		Print("The ", NoNewline)
		PrintObject(G.DirObj)
		if G.DirObj.Has(FlgOpen) {
			Print(" is open, but I can't tell what's beyond it.", NoNewline)
		} else {
			Print(" is closed.", NoNewline)
		}
		NewLine()
		return true
	}
	if G.DirObj.Has(FlgCont) {
		if G.DirObj.Has(FlgPerson) {
			Print("There is nothing special to be seen.", Newline)
			return true
		}
		if !CanSeeInside(G.DirObj) {
			Print("The ", NoNewline)
			PrintObject(G.DirObj)
			Print(" is closed.", Newline)
			return true
		}
		if G.DirObj.HasChildren() && PrintCont(G.DirObj, false, 0) {
			return true
		}
		Print("The ", NoNewline)
		PrintObject(G.DirObj)
		Print(" is empty.", Newline)
		return true
	}
	Print("You can't look inside a ", NoNewline)
	PrintObject(G.DirObj)
	Print(".", Newline)
	return true
}

func VLookOn(arg ActArg) bool {
	if G.DirObj.Has(FlgSurf) {
		Perform(ActionVerb{Norm: "look inside", Orig: "look inside"}, G.DirObj, nil)
		return true
	}
	Print("Look on a ", NoNewline)
	PrintObject(G.DirObj)
	Print("???", Newline)
	return true
}

func VLookUnder(arg ActArg) bool {
	Print("There is nothing but dust there.", Newline)
	return true
}

func VExamine(arg ActArg) bool {
	if len(G.DirObj.Text) > 0 {
		Print(G.DirObj.Text, Newline)
		return true
	}
	if G.DirObj.Has(FlgCont) || G.DirObj.Has(FlgDoor) {
		return VLookInside(ActUnk)
	}
	Print("There's nothing special about the ", NoNewline)
	PrintObject(G.DirObj)
	Print(".", Newline)
	return true
}

func VMake(arg ActArg) bool {
	Print("You can't do that.", Newline)
	return true
}

func VMelt(arg ActArg) bool {
	Print("It's not clear that a ", NoNewline)
	PrintObject(G.DirObj)
	Print(" can be melted.", Newline)
	return true
}

func PreMove(arg ActArg) bool {
	if IsHeld(G.DirObj) {
		Print("You aren't an accomplished enough juggler.", Newline)
		return true
	}
	return false
}

func VMove(arg ActArg) bool {
	if G.DirObj.Has(FlgTake) {
		Print("Moving the ", NoNewline)
		PrintObject(G.DirObj)
		Print(" reveals nothing.", Newline)
		return true
	}
	Print("You can't move the ", NoNewline)
	PrintObject(G.DirObj)
	Print(".", Newline)
	return true
}

func VMumble(arg ActArg) bool {
	Print("You'll have to speak up if you expect me to hear you!", Newline)
	return true
}

func PreMung(arg ActArg) bool {
	if G.IndirObj == nil || !G.IndirObj.Has(FlgWeapon) {
		Print("Trying to destroy the ", NoNewline)
		PrintObject(G.DirObj)
		Print(" with ", NoNewline)
		if G.IndirObj == nil {
			Print("your bare hands", NoNewline)
		} else {
			Print("a", NoNewline)
			PrintObject(G.IndirObj)
		}
		Print(" is futile.", Newline)
		return true
	}
	return false
}

func VMung(arg ActArg) bool {
	if !G.DirObj.Has(FlgPerson) {
		Print("Nice try.", Newline)
		return true
	}
	Perform(ActionVerb{Norm: "attack", Orig: "attack"}, G.DirObj, nil)
	return true
}

func MungRoom(rm *Object, str string) {
	rm.Give(FlgKludge)
	rm.LongDesc = str
}

func VOdysseus(arg ActArg) bool {
	if G.Here != &CyclopsRoom || !Cyclops.IsIn(G.Here) || G.CyclopsFlag {
		Print("Wasn't he a sailor?", Newline)
		return true
	}
	QueueInt(ICyclops, false).Run = false
	G.CyclopsFlag = true
	Print("The cyclops, hearing the name of his father's deadly nemesis, flees the room by knocking down the wall on the east of the room.", Newline)
	G.MagicFlag = true
	Cyclops.Take(FlgFight)
	return RemoveCarefully(&Cyclops)
}

func VOil(arg ActArg) bool {
	Print("You probably put spinach in your gas tank, too.", Newline)
	return true
}

func VOpen(arg ActArg) bool {
	if !G.DirObj.Has(FlgCont) || G.DirObj.Capacity == 0 {
		if G.DirObj.Has(FlgDoor) {
			if G.DirObj.Has(FlgOpen) {
				Print("It is already open.", Newline)
				return true
			}
			Print("The ", NoNewline)
			PrintObject(G.DirObj)
			Print(" opens.", Newline)
			G.DirObj.Give(FlgOpen)
			return true
		}
		Print("You must tell me how to do that to a ", NoNewline)
		PrintObject(G.DirObj)
		Print(".", Newline)
		return true
	}
	if G.DirObj.Has(FlgOpen) {
		Print("It is already open.", Newline)
		return true
	}
	G.DirObj.Give(FlgOpen)
	G.DirObj.Give(FlgTouch)
	if !G.DirObj.HasChildren() || G.DirObj.Has(FlgTrans) {
		Print("Opened.", Newline)
		return true
	}
	if len(G.DirObj.Children) == 1 && !G.DirObj.Children[0].Has(FlgTouch) && len(G.DirObj.Children[0].FirstDesc) > 0 {
		str := G.DirObj.Children[0].FirstDesc
		Print("The ", NoNewline)
		PrintObject(G.DirObj)
		Print(" opens.", Newline)
		Print(str, Newline)
		return true
	}
	Print("Opening the ", NoNewline)
	PrintObject(G.DirObj)
	Print(" reveals ", NoNewline)
	PrintContents(G.DirObj)
	Print(".", Newline)
	return true
}

func VOverboard(arg ActArg) bool {
	locn := G.Winner.Location()
	if G.IndirObj == &Teeth {
		if locn.Has(FlgVeh) {
			Print("Ahoy -- ", NoNewline)
			PrintObject(G.IndirObj)
			Print(" overboard!", Newline)
			return true
		}
		Print("You're not in anything!", Newline)
		return true
	}
	if locn.Has(FlgVeh) {
		Perform(ActionVerb{Norm: "throw", Orig: "throw"}, G.DirObj, nil)
		return true
	}
	Print("Huh?", Newline)
	return true
}

func VPick(arg ActArg) bool {
	Print("You can't pick that.", Newline)
	return true
}

func VPlay(arg ActArg) bool {
	if !G.DirObj.Has(FlgPerson) {
		Print("That's silly!", Newline)
		return true
	}
	Print("You become so engrossed in the role of the ", NoNewline)
	PrintObject(G.DirObj)
	Print(" that you kill yourself, just as he might have done!", Newline)
	return JigsUp("", false)
}

func VPlug(arg ActArg) bool {
	Print("This has no effect.", Newline)
	return true
}

func VPourOn(arg ActArg) bool {
	if G.DirObj == &Water {
		RemoveCarefully(G.DirObj)
		if IsFlaming(G.IndirObj) {
			Print("The ", NoNewline)
			PrintObject(G.IndirObj)
			Print(" is extinguished.", Newline)
			G.IndirObj.Take(FlgOn)
			G.IndirObj.Take(FlgFlame)
			return true
		}
		Print("The water spills over the ", NoNewline)
		PrintObject(G.IndirObj)
		Print(", to the floor, and evaporates.", Newline)
		return true
	}
	if G.DirObj == &Putty {
		if Perform(ActionVerb{Norm: "put", Orig: "put"}, &Putty, G.IndirObj) == PerfHndld {
			return true
		}
		return false
	}
	Print("You can't pour that.", Newline)
	return true
}

func VPray(arg ActArg) bool {
	if G.Here != &SouthTemple {
		Print("If you pray enough, your prayers may be answered.", Newline)
		return true
	}
	return Goto(&Forest1, true)
}

func VPump(arg ActArg) bool {
	if G.IndirObj != nil && G.IndirObj != &Pump {
		Print("Pump it up with a ", NoNewline)
		PrintObject(G.IndirObj)
		Print("?", Newline)
		return true
	}
	if Pump.IsIn(G.Winner) {
		if Perform(ActionVerb{Norm: "inflate", Orig: "inflate"}, G.DirObj, &Pump) == PerfHndld {
			return true
		}
		return false
	}
	Print("It's really not clear how.", Newline)
	return true
}

func VPush(arg ActArg) bool {
	return HackHack("Pushing the ")
}

func VPushTo(arg ActArg) bool {
	Print("You can't push things to that.", Newline)
	return true
}

func PrePut(arg ActArg) bool {
	return PreGive(arg)
}

func VPut(arg ActArg) bool {
	if !G.IndirObj.Has(FlgOpen) && !IsOpenable(G.IndirObj) && !G.IndirObj.Has(FlgVeh) {
		Print("You can't do that.", Newline)
		return true
	}
	if !G.IndirObj.Has(FlgOpen) {
		Print("The ", NoNewline)
		PrintObject(G.IndirObj)
		Print(" isn't open.", Newline)
		ThisIsIt(G.IndirObj)
		return true
	}
	if G.IndirObj == G.DirObj {
		Print("How can you do that?", Newline)
		return true
	}
	if G.DirObj.IsIn(G.IndirObj) {
		Print("The ", NoNewline)
		PrintObject(G.DirObj)
		Print(" is already in the ", NoNewline)
		PrintObject(G.IndirObj)
		Print(".", Newline)
		return true
	}
	if Weight(G.IndirObj)+Weight(G.DirObj)-G.IndirObj.Size > G.IndirObj.Capacity {
		Print("There's no room.", Newline)
		return true
	}
	if !IsHeld(G.DirObj) && !ITake(true) {
		return true
	}
	G.DirObj.MoveTo(G.IndirObj)
	G.DirObj.Give(FlgTouch)
	ScoreObj(G.DirObj)
	Print("Done.", Newline)
	return true
}

func VPutBehind(arg ActArg) bool {
	Print("That hiding place is too obvious.", Newline)
	return true
}

func VPutOn(arg ActArg) bool {
	if G.IndirObj == nil || G.IndirObj == &Ground {
		return VDrop(ActUnk)
	}
	if G.IndirObj.Has(FlgSurf) {
		return VPut(ActUnk)
	}
	Print("There's no good surface on the ", NoNewline)
	PrintObject(G.IndirObj)
	Print(".", Newline)
	return true
}

func VPutUnder(arg ActArg) bool {
	Print("You can't do that.", Newline)
	return true
}

func VRaise(arg ActArg) bool {
	return VLower(arg)
}

func VRape(arg ActArg) bool {
	Print("What a (ahem!) strange idea.", Newline)
	return true
}

func PreRead(arg ActArg) bool {
	if !G.Lit {
		Print("It is impossible to read in the dark.", Newline)
		return true
	}
	if G.IndirObj != nil && !G.IndirObj.Has(FlgTrans) {
		Print("How does one look through a ", NoNewline)
		PrintObject(G.IndirObj)
		Print("?", Newline)
		return true
	}
	return false
}

func VRead(arg ActArg) bool {
	if !G.DirObj.Has(FlgRead) {
		Print("How does one read a ", NoNewline)
		PrintObject(G.DirObj)
		Print("?", Newline)
		return true
	}
	Print(G.DirObj.Text, Newline)
	return true
}

func VReadPage(arg ActArg) bool {
	return VRead(ActUnk)
}

func VRepent(arg ActArg) bool {
	Print("It could very well be too late!", Newline)
	return true
}

func VReply(arg ActArg) bool {
	Print("It is hardly likely that the ", NoNewline)
	PrintObject(G.DirObj)
	Print(" is interested.", Newline)
	G.Params.Continue = NumUndef
	G.Params.InQuotes = false
	return true
}

func VRing(arg ActArg) bool {
	Print("How, exactly, can you ring that?", Newline)
	return true
}

func VRub(arg ActArg) bool {
	return HackHack("Fiddling with the ")
}

func VSay(arg ActArg) bool {
	if G.Params.Continue == NumUndef {
		Print("Say what?", Newline)
		return true
	}
	G.Params.InQuotes = false
	v := FindIn(G.Here, FlgPerson)
	if v != nil {
		Print("You must address the ", NoNewline)
		PrintObject(v)
		Print(" directly.", Newline)
		G.Params.Continue = NumUndef
		return true
	}
	if G.LexRes[G.Params.Continue].Norm == "hello" {
		G.Params.Continue = NumUndef
		Print("Talking to yourself is a sign of impending mental collapse.", Newline)
		return true
	}
	return false
}

func FindIn(where *Object, what Flags) *Object {
	if !where.HasChildren() {
		return nil
	}
	for w := 0; w < len(where.Children); w++ {
		if where.Children[w].Has(what) && where.Children[w] != &Adventurer {
			return where.Children[w]
		}
	}
	return nil
}

func VSearch(arg ActArg) bool {
	Print("You find nothing unusual.", Newline)
	return true
}

func VSend(arg ActArg) bool {
	if G.DirObj.Has(FlgPerson) {
		Print("Why would you send for the ", NoNewline)
		PrintObject(G.DirObj)
		Print("?", Newline)
		return true
	}
	Print("That doesn't make sends.", Newline)
	return true
}

func PreSGive(arg ActArg) bool {
	Perform(ActionVerb{Norm: "give", Orig: "give"}, G.IndirObj, G.DirObj)
	return true
}

func VSGive(arg ActArg) bool {
	Print("Foo!", Newline)
	return true
}

func VShake(arg ActArg) bool {
	if G.DirObj.Has(FlgPerson) {
		Print("This seems to have no effect.", Newline)
		return true
	}
	if !G.DirObj.Has(FlgTake) {
		Print("You can't take it; thus, you can't shake it!", Newline)
		return true
	}
	if !G.DirObj.Has(FlgCont) {
		Print("Shaken.", Newline)
		return true
	}
	if G.DirObj.Has(FlgOpen) {
		if !G.DirObj.HasChildren() {
			Print("Shaken.", Newline)
			return true
		}
		ShakeLoop()
		Print("The contents of the ", NoNewline)
		PrintObject(G.DirObj)
		Print(" spill ", NoNewline)
		if !G.Here.Has(FlgLand) {
			Print("out and disappears", NoNewline)
		} else {
			Print("to the ground", NoNewline)
		}
		Print(".", Newline)
		return true
	}
	if G.DirObj.HasChildren() {
		Print("It sounds like there is something inside the ", NoNewline)
		PrintObject(G.DirObj)
		Print(".", Newline)
		return true
	}
	Print("The ", NoNewline)
	PrintObject(G.DirObj)
	Print(" sounds empty.", Newline)
	return true
}

func ShakeLoop() {
	if !G.DirObj.HasChildren() {
		return
	}
	x := G.DirObj.Children[0]
	x.Give(FlgTouch)
	mv := G.Here
	if G.Here == &UpATree {
		mv = &Path
	} else if !G.Here.Has(FlgLand) {
		mv = &PseudoObject
	}
	x.MoveTo(mv)
}

func VSkip(arg ActArg) bool {
	Print(PickOne(Wheeee), Newline)
	return true
}

func VSmell(arg ActArg) bool {
	Print("It smells like a ", NoNewline)
	PrintObject(G.DirObj)
	Print(".", Newline)
	return true
}

func VSpin(arg ActArg) bool {
	Print("You can't spin that!", Newline)
	return true
}

func VSpray(arg ActArg) bool {
	return VSqueeze(arg)
}

func VSqueeze(arg ActArg) bool {
	if !G.DirObj.Has(FlgPerson) {
		Print("How singularly useless.", Newline)
		return true
	}
	Print("The ", NoNewline)
	PrintObject(G.DirObj)
	Print(" does not understand this.", Newline)
	return true
}

func VSSpray(arg ActArg) bool {
	if Perform(ActionVerb{Norm: "spray", Orig: "spray"}, G.IndirObj, G.DirObj) == PerfHndld {
		return true
	}
	return false
}

func VStab(arg ActArg) bool {
	w := FindWeapon(G.Winner)
	if w == nil {
		Print("No doubt you propose to stab the ", NoNewline)
		PrintObject(G.DirObj)
		Print(" with your pinky?", Newline)
		return true
	}
	Perform(ActionVerb{Norm: "attack", Orig: "attack"}, G.DirObj, w)
	return true
}

func VStrike(arg ActArg) bool {
	if G.DirObj.Has(FlgPerson) {
		Print("Since you aren't versed in hand-to-hand combat, you'd better attack the ", NoNewline)
		PrintObject(G.DirObj)
		Print(" with a weapon.", Newline)
		return true
	}
	Perform(ActionVerb{Norm: "lamp on", Orig: "lamp on"}, G.DirObj, nil)
	return true
}

func VSwing(arg ActArg) bool {
	if G.IndirObj == nil {
		Print("Whoosh!", Newline)
		return true
	}
	if Perform(ActionVerb{Norm: "attack", Orig: "attack"}, G.IndirObj, G.DirObj) == PerfHndld {
		return true
	}
	return false
}

func PreTake(arg ActArg) bool {
	if G.DirObj == G.Winner {
		if G.DirObj.Has(FlgWear) {
			Print("You are already wearing it.", Newline)
			return true
		}
		Print("You already have that!", Newline)
		return true
	}
	lcn := G.DirObj.Location()
	if lcn.Has(FlgCont) && !lcn.Has(FlgOpen) {
		Print("You can't reach something that's inside a closed container.", Newline)
		return true
	}
	if G.IndirObj != nil {
		if G.IndirObj == &Ground {
			G.IndirObj = nil
			return false
		}
		if G.IndirObj != G.DirObj.Location() {
			Print("The ", NoNewline)
			PrintObject(G.DirObj)
			Print(" isn't in the ", NoNewline)
			PrintObject(G.IndirObj)
			Print(".", Newline)
			return true
		}
		G.IndirObj = nil
		return false
	}
	if G.DirObj == G.Winner.Location() {
		Print("You're inside of it!", Newline)
		return true
	}
	return false
}

func VTake(arg ActArg) bool {
	if ITake(true) {
		if G.DirObj.Has(FlgWear) {
			Print("You are now wearing the ", NoNewline)
			PrintObject(G.DirObj)
			Print(".", Newline)
			return true
		}
		Print("Taken.", Newline)
		return true
	}
	return false
}

func VTell(arg ActArg) bool {
	if !G.DirObj.Has(FlgPerson) {
		Print("You can't talk to the ", NoNewline)
		PrintObject(G.DirObj)
		Print("!", Newline)
		G.Params.InQuotes = false
		G.Params.Continue = NumUndef
		return RFatal()
	}
	if G.Params.Continue != NumUndef {
		G.Winner = G.DirObj
		G.Here = G.Winner.Location()
		return true
	}
	Print("The ", NoNewline)
	PrintObject(G.DirObj)
	Print(" pauses for a moment, perhaps thinking that you should reread the manual.", Newline)
	return true
}

func VThrow(arg ActArg) bool {
	if !IDrop() {
		Print("Huh?", Newline)
		return true
	}
	if G.IndirObj == &Me {
		Print("A terrific throw! The ", NoNewline)
		G.Winner = G.Player
		return JigsUp(" hits you squarely in the head. Normally, this wouldn't do much damage, but by incredible mischance, you fall over backwards trying to duck, and break your neck, justice being swift and merciful in the Great Underground Empire.", false)
	}
	if G.IndirObj != nil && G.IndirObj.Has(FlgPerson) {
		Print("The ", NoNewline)
		PrintObject(G.IndirObj)
		Print(" ducks as the ", NoNewline)
		PrintObject(G.DirObj)
		Print(" flies by and crashes to the ground.", Newline)
		return true
	}
	Print("Thrown.", Newline)
	return true
}

func VThrowOff(arg ActArg) bool {
	Print("You can't throw anything off of that!", Newline)
	return true
}

func VTie(arg ActArg) bool {
	if G.IndirObj == G.Winner {
		Print("You can't tie anything to yourself.", Newline)
		return true
	}
	Print("You can't tie the ", NoNewline)
	PrintObject(G.DirObj)
	Print(" to that.", Newline)
	return true
}

func VTieUp(arg ActArg) bool {
	Print("You could certainly never tie it with that!", Newline)
	return true
}

func VTreasure(arg ActArg) bool {
	if G.Here == &NorthTemple {
		return Goto(&TreasureRoom, true)
	}
	if G.Here == &TreasureRoom {
		return Goto(&NorthTemple, true)
	}
	Print("Nothing happens.", Newline)
	return true
}

func PreTurn(arg ActArg) bool {
	if (G.IndirObj == nil || G.IndirObj == &Rooms) && G.DirObj != &Book {
		Print("Your bare hands don't appear to be enough.", Newline)
		return true
	}
	if !G.DirObj.Has(FlgTurn) {
		Print("You can't turn that!", Newline)
		return true
	}
	return false
}

func VTurn(arg ActArg) bool {
	Print("This has no effect.", Newline)
	return true
}

func VUnlock(arg ActArg) bool {
	return VLock(arg)
}

func VUntie(arg ActArg) bool {
	Print("This cannot be tied, so it cannot be untied!", Newline)
	return true
}

func VWave(arg ActArg) bool {
	return HackHack("Waving the ")
}

func VWear(arg ActArg) bool {
	if !G.DirObj.Has(FlgWear) {
		Print("You can't wear the ", NoNewline)
		PrintObject(G.DirObj)
		Print(".", Newline)
		return true
	}
	Perform(ActionVerb{Norm: "take", Orig: "take"}, G.DirObj, nil)
	return true
}

func VWind(arg ActArg) bool {
	Print("You cannot wind up a ", NoNewline)
	PrintObject(G.DirObj)
	Print(".", Newline)
	return true
}

func VYell(arg ActArg) bool {
	Print("Aaaarrrrgggghhhh!", Newline)
	return true
}

func RemoveCarefully(obj *Object) bool {
	if obj == G.Params.ItObj {
		G.Params.ItObj = nil
	}
	oLit := G.Lit
	obj.Remove()
	G.Lit = IsLit(G.Here, true)
	if oLit && oLit != G.Lit {
		Print("You are left in the dark...", Newline)
	}
	return true
}

func PrintContents(obj *Object) bool {
	if !obj.HasChildren() {
		return false
	}
	var itObj *Object
	twoIs := false
	for n := 0; n < len(obj.Children); n++ {
		if n != 0 {
			Print(", ", NoNewline)
			if n+1 >= len(obj.Children) {
				Print("and ", NoNewline)
			}
		}
		Print("a ", NoNewline)
		PrintObject(obj.Children[n])
		if itObj == nil && !twoIs {
			itObj = obj.Children[n]
		} else {
			twoIs = true
			itObj = nil
		}
	}
	if itObj != nil && !twoIs {
		ThisIsIt(itObj)
	}
	return true
}

func HackHack(str string) bool {
	if G.DirObj.IsIn(&GlobalObjects) && (G.ActVerb.Norm == "wave" || G.ActVerb.Norm == "raise" || G.ActVerb.Norm == "lower") {
		Print("The ", NoNewline)
		PrintObject(G.DirObj)
		Print(" isn't here!", Newline)
		return true
	}
	Print(str, NoNewline)
	PrintObject(G.DirObj)
	Print(PickOne(Hohum), Newline)
	return true
}

func IDrop() bool {
	if !G.DirObj.IsIn(G.Winner) && !G.DirObj.Location().IsIn(G.Winner) {
		Print("You're not carrying the ", NoNewline)
		PrintObject(G.DirObj)
		Print(".", Newline)
		return false
	}
	if !G.DirObj.IsIn(G.Winner) && !G.DirObj.Location().Has(FlgOpen) {
		Print("The ", NoNewline)
		PrintObject(G.DirObj)
		Print(" is closed.", Newline)
		return false
	}
	G.DirObj.MoveTo(G.Winner.Location())
	return true
}

func DescribeRoom(isLook bool) bool {
	v := G.Verbose
	if isLook {
		v = isLook
	}
	if !G.Lit {
		Print("It is pitch black.", NoNewline)
		if !G.IsSprayed {
			Print(" You are likely to be eaten by a grue.", NoNewline)
		}
		NewLine()
		return false
	}
	if !G.Here.Has(FlgTouch) {
		G.Here.Give(FlgTouch)
		v = true
	}
	if G.Here.Has(FlgMaze) {
		G.Here.Take(FlgTouch)
	}
	if G.Here.IsIn(&Rooms) {
		PrintObject(G.Here)
		if av := G.Winner.Location(); av.Has(FlgVeh) {
			Print(", in the ", NoNewline)
			PrintObject(av)
		}
		NewLine()
	}
	if !isLook && G.SuperBrief {
		return true
	}
	av := G.Winner.Location()
	if v && G.Here.Action != nil && G.Here.Action(ActLook) {
		return true
	}
	if v && len(G.Here.LongDesc) != 0 {
		Print(G.Here.LongDesc, Newline)
	} else if G.Here.Action != nil {
		G.Here.Action(ActFlash)
	}
	if G.Here != av && av.Has(FlgVeh) && av.Action != nil {
		av.Action(ActLook)
	}
	return true
}

func DescribeObjects(v bool) bool {
	if !G.Lit {
		Print("Only bats can see in the dark. And you're not one.", Newline)
		return true
	}
	if !G.Here.HasChildren() {
		return false
	}
	if !v {
		v = G.Verbose
	}
	return PrintCont(G.Here, v, -1)
}

func PrintCont(obj *Object, v bool, lvl int) bool {
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
					Print(child.FirstDesc, Newline)
					shit = false
				}
				if CanSeeInside(child) && child.Location().DescFcn == nil && child.HasChildren() {
					if PrintCont(child, v, 0) {
						isFirst = false
					}
				}
			}
		}
	}
	for _, child := range obj.Children {
		if child == av || child == &Adventurer {
			continue
		}
		if child.Has(FlgInvis) || (!isInv && !child.Has(FlgTouch) && len(child.FirstDesc) > 0) {
			continue
		}
		if !child.Has(FlgNoDesc) {
			if isFirst {
				if Firster(obj, lvl) && lvl < 0 {
					lvl = 0
				}
				lvl++
				isFirst = false
			}
			if lvl < 0 {
				lvl = 0
			}
			DescribeObject(child, v, lvl)
		} else if child.HasChildren() && CanSeeInside(child) {
			lvl++
			PrintCont(child, v, lvl)
			lvl--
		}
	}
	if isPv && av != nil && av.HasChildren() {
		lvl++
		PrintCont(av, v, lvl)
	}
	if isFirst && shit {
		return false
	}
	return true
}

func Firster(obj *Object, lvl int) bool {
	if obj == &TrophyCase {
		Print("Your collection of treasures consists of:", Newline)
		return true
	}
	if obj == G.Winner {
		Print("You are carrying:", Newline)
		return true
	}
	if obj.IsIn(&Rooms) {
		return false
	}
	if lvl > 0 {
		Print(Indents[lvl], NoNewline)
	}
	if obj.Has(FlgSurf) {
		Print("Sitting on the ", NoNewline)
		PrintObject(obj)
		Print(" is: ", Newline)
		return true
	}
	if obj.Has(FlgPerson) {
		Print("The ", NoNewline)
		PrintObject(obj)
		Print(" is holding: ", Newline)
		return true
	}
	Print("The ", NoNewline)
	PrintObject(obj)
	Print(" contains:", Newline)
	return true
}

func DescribeObject(obj *Object, v bool, lvl int) bool {
	G.DescObj = obj
	if lvl == 0 && obj.DescFcn != nil && obj.DescFcn(ActObjDesc) {
		return true
	}
	if lvl == 0 && ((!obj.Has(FlgTouch) && len(obj.FirstDesc) > 0) || len(obj.LongDesc) > 0) {
		if !obj.Has(FlgTouch) && len(obj.FirstDesc) > 0 {
			Print(obj.FirstDesc, NoNewline)
		} else {
			Print(obj.LongDesc, NoNewline)
		}
	} else if lvl == 0 {
		Print("There is a ", NoNewline)
		PrintObject(obj)
		Print(" here", NoNewline)
		if obj.Has(FlgOn) {
			Print(" (providing light)", NoNewline)
		}
		Print(".", NoNewline)
	} else {
		Print(Indents[lvl], NoNewline)
		Print("A ", NoNewline)
		PrintObject(obj)
		if obj.Has(FlgOn) {
			Print(" (providing light)", NoNewline)
		} else if obj.Has(FlgWear) && obj.IsIn(G.Winner) {
			Print(" (being worn)", NoNewline)
		}
	}
	if av := G.Winner.Location(); lvl == 0 && av != nil && av.Has(FlgVeh) {
		Print(" (outside the ", NoNewline)
		PrintObject(av)
		Print(")", NoNewline)
	}
	NewLine()
	if CanSeeInside(obj) && obj.HasChildren() {
		return PrintCont(obj, v, lvl)
	}
	return true
}

func CanSeeInside(obj *Object) bool {
	if obj.Has(FlgInvis) {
		return false
	}
	if obj.Has(FlgTrans) || obj.Has(FlgOpen) {
		return true
	}
	return false
}

func ThisIsIt(obj *Object) {
	G.Params.ItObj = obj
}

func IsInGlobal(obj1, obj2 *Object) bool {
	if obj2.Global == nil {
		return false
	}
	for _, o := range obj2.Global {
		if o == obj1 {
			return true
		}
	}
	return false
}

func IsHeld(obj *Object) bool {
	for {
		obj = obj.Location()
		if obj == nil {
			return false
		}
		if obj == G.Winner {
			return true
		}
	}
}

func ScoreUpd(num int) bool {
	G.BaseScore += num
	G.Score += num
	if G.Score == 350 && !G.WonGame {
		G.WonGame = true
		Map.Take(FlgInvis)
		WestOfHouse.Take(FlgTouch)
		Print("An almost inaudible voice whispers in your ear, \"Look to your treasures for the final secret.\"", Newline)
	}
	return true
}

func ScoreObj(obj *Object) {
	if obj.Value <= 0 {
		return
	}
	ScoreUpd(obj.Value)
	obj.Value = 0
}

func CCount(obj *Object) int {
	cnt := 0
	for _, child := range obj.Children {
		if !child.Has(FlgWear) {
			cnt++
		}
	}
	return cnt
}

func Weight(obj *Object) int {
	wt := 0
	for _, child := range obj.Children {
		if obj == G.Player && child.Has(FlgWear) {
			wt++
		} else {
			wt += Weight(child)
		}
	}
	return wt + obj.Size
}

func ITake(vb bool) bool {
	if G.Dead {
		if vb {
			Print("Your hand passes through its object.", Newline)
		}
		return false
	}
	if G.DirObj == nil {
		return false
	}
	if !G.DirObj.Has(FlgTake) {
		if vb {
			Print(PickOne(Yuks), Newline)
		}
		return false
	}
	// ZIL: <FSET? <LOC ,PRSO> ,CONTBIT> / <NOT <FSET? <LOC ,PRSO> ,OPENBIT>>
	// Prevent taking objects from inside a closed container.
	loc := G.DirObj.Location()
	if loc != nil && loc.Has(FlgCont) && !loc.Has(FlgOpen) {
		return false
	}
	if !G.DirObj.Location().IsIn(G.Winner) && Weight(G.DirObj)+Weight(G.Winner) > G.LoadAllowed {
		if vb {
			Print("Your load is too heavy", NoNewline)
			if G.LoadAllowed < G.LoadMax {
				Print(", especially in light of your condition.", NoNewline)
			} else {
				Print(".", NoNewline)
			}
			NewLine()
		}
		G.PerformFatal = true
		return false
	}
	cnt := CCount(G.Winner)
	if G.ActVerb.Norm == "tell" && cnt > G.FumbleNumber && Prob(cnt*G.FumbleProb, false) {
		Print("You're holding too many things already!", Newline)
		return false
	}
	G.DirObj.MoveTo(G.Winner)
	G.DirObj.Take(FlgNoDesc)
	G.DirObj.Give(FlgTouch)
	ScoreObj(G.DirObj)
	return true
}

