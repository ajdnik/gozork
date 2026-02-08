package zork

import (
	"encoding/binary"
)

type RndSelect struct {
	Unselected []string
	Selected   []string
}

var (
	Moves      int
	Score      int
	BaseScore  int
	Lit        bool
	SuperBrief bool
	Verbose    bool
	Dead       bool
	Deaths     int
	// HelloSailor counts occurences of 'hello, sailor'
	HelloSailor int
	// IsSprayed is a flag indicating if the player is wearing Grue repellent
	IsSprayed bool
	Version   = [24]byte{0, 0, 0, 0, 119, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 56, 56, 48, 52, 50, 57}
	WonGame   bool
	// DescObj stores the last object which was described
	DescObj *Object
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
	Lucky        = true
	FumbleNumber = 7
	FumbleProb   = 8
)

func VVerbose(arg ActArg) bool {
	Verbose = true
	SuperBrief = false
	Print("Maximum verbosity.", Newline)
	return true
}

func VBrief(arg ActArg) bool {
	Verbose = false
	SuperBrief = false
	Print("Brief descriptions.", Newline)
	return true
}

func VSuperBrief(arg ActArg) bool {
	SuperBrief = true
	Print("Superbrief descriptions.", Newline)
	return true
}

func VInventory(arg ActArg) bool {
	if Winner.HasChildren() {
		return PrintCont(Winner, false, 0)
	}
	Print("You are empty-handed.", Newline)
	return true
}

func VRestart(arg ActArg) bool {
	VScore(ActUnk)
	Print("Do you wish to restart? (Y is affirmative): ", NoNewline)
	if IsYes() {
		Print("Restarting.", Newline)
		Restart()
		Print("Failed.", Newline)
		return true
	}
	return false
}

func VRestore(arg ActArg) bool {
	if Restore() {
		Print("Ok.", Newline)
		return VFirstLook(ActUnk)
	}
	Print("Failed.", Newline)
	return true
}

func VSave(arg ActArg) bool {
	if Save() {
		Print("Ok.", Newline)
		return true
	}
	Print("Failed.", Newline)
	return true
}

func VScript(arg ActArg) bool {
	// This code turns on the first bit in the 8th word from the beginning
	// Put(0, 8, Bor(Get(0, 8), 1))
	Script = true
	Print("Here begins a transcript of interaction with", Newline)
	VVersion(ActUnk)
	return true
}

func VUnscript(arg ActArg) bool {
	Print("Here ends a transcript of interaction with", Newline)
	VVersion(ActUnk)
	// This code turns off the first bit in the 8th word from the beginning
	// Put(0, 8, Band(Get(0, 8), -2))
	Script = false
	return true
}

func VVerify(arg ActArg) bool {
	Print("Verifying disk...", Newline)
	if Verify() {
		Print("The disk is correct.", Newline)
		return true
	}
	NewLine()
	Print("** Disk Failure **", Newline)
	return true
}

// TODO: VCommandFile, VRandom, VRecord, VUnrecord

func VVersion(arg ActArg) bool {
	Print("ZORK I: The Great Underground Empire", Newline)
	Print("Infocom interactive fiction - a fantasy story", Newline)
	Print("Copyright (c) 1981, 1982, 1983, 1984, 1985, 1986", NoNewline)
	Print(" Infocom, Inc. All rights reserved.", Newline)
	Print("ZORK is a registered trademark of Infocom, Inc.", Newline)
	Print("Release ", NoNewline)
	num := binary.LittleEndian.Uint16(Version[4:6])
	PrintNumber(int(num & 2047))
	Print(" / Serial number ", NoNewline)
	for offset := 18; offset <= 23; offset++ {
		Print(string(Version[offset]), NoNewline)
	}
	NewLine()
	return true
}

func VAdvent(arg ActArg) bool {
	Print("A hollow voice says \"Fool.\"", Newline)
	return true
}

func VAlarm(arg ActArg) bool {
	if !DirObj.Has(FlgPerson) {
		Print("The ", NoNewline)
		PrintObject(DirObj)
		Print(" isn't sleeping.", Newline)
		return true
	}
	if DirObj.Strength <= 0 {
		Print("He's wide awake, or haven't you noticed...", Newline)
		return true
	}
	Print("The ", NoNewline)
	PrintObject(DirObj)
	Print(" is rudely awakened.", Newline)
	return Awaken(DirObj)
}

func VAnswer(arg ActArg) bool {
	Print("Nobody seems to be awaiting your answer.", Newline)
	Params.Continue = NumUndef
	Params.InQuotes = false
	return true
}

func VAttack(arg ActArg) bool {
	if !DirObj.Has(FlgPerson) {
		Print("I've known strange people, but fighting a ", NoNewline)
		PrintObject(DirObj)
		Print("?", Newline)
		return true
	}
	if IndirObj == nil || IndirObj == &Hands {
		Print("Trying to attack a ", NoNewline)
		PrintObject(DirObj)
		Print(" with your bare hands is suicidal.", Newline)
		return true
	}
	if !IndirObj.IsIn(Winner) {
		Print("You aren't even holding the ", NoNewline)
		PrintObject(IndirObj)
		Print(".", Newline)
		return true
	}
	if !IndirObj.Has(FlgWeapon) {
		Print("Trying to attack the ", NoNewline)
		PrintObject(DirObj)
		Print(" with a ", NoNewline)
		PrintObject(IndirObj)
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

func PreBoard(arg ActArg) bool {
	if DirObj.Has(FlgVeh) {
		if !DirObj.IsIn(Here) {
			Print("The ", NoNewline)
			PrintObject(DirObj)
			Print(" must be on the ground to be boarded.", Newline)
			// TODO: return fatal
			return false
		}
		if av := Winner.Location(); av != nil && av.Has(FlgVeh) {
			Print("You are already in the ", NoNewline)
			PrintObject(av)
			Print("!", Newline)
			// TODO. return fatal
			return false
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
	// TODO: return fatal
	return false
}

func VBoard(arg ActArg) bool {
	Print("You are now in the ", NoNewline)
	PrintObject(DirObj)
	Print(".", Newline)
	Winner.MoveTo(Here)
	// TODO: Maybe check if Action is defined?
	DirObj.Action(ActEnter)
	return true
}

func VBreathe(arg ActArg) bool {
	// TODO: In original code this returns perform
	if Perform(ActionVerb{Norm: "inflate", Orig: "inflate"}, DirObj, &Lungs) == PerfHndld {
		return true
	}
	return false
}

func VBrush(arg ActArg) bool {
	Print("If you wish, but heaven only knows why.", Newline)
	return true
}

func VBug(arg ActArg) bool {
	Print("Bug? Not in a flawless program like this! (Cough, cough).", Newline)
	return true
}

func TellNoPrsi() bool {
	Print("You didn't say with what!", Newline)
	return true
}

func PreBurn(arg ActArg) bool {
	if IndirObj == nil {
		Print("You didn't say with what!", Newline)
		return true
	}
	if IsFlaming(IndirObj) {
		return false
	}
	Print("With a ", NoNewline)
	PrintObject(IndirObj)
	Print("??!?", Newline)
	return true
}

func VBurn(arg ActArg) bool {
	if !DirObj.Has(FlgBurn) {
		Print("You can't burn a ", NoNewline)
		PrintObject(DirObj)
		Print(".", Newline)
		return true
	}
	if !DirObj.IsIn(Winner) && !Winner.IsIn(DirObj) {
		RemoveCarefully(DirObj)
		Print("The ", NoNewline)
		PrintObject(DirObj)
		Print(" catches fire and is consumed.", Newline)
		return true
	}
	RemoveCarefully(DirObj)
	Print("The ", NoNewline)
	PrintObject(DirObj)
	Print(" catches fire. Unfortunately, you were ", NoNewline)
	if Winner.IsIn(DirObj) {
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

func VClose(arg ActArg) bool {
	if !DirObj.Has(FlgCont) && !DirObj.Has(FlgDoor) {
		Print("You must tell me how to do that to a ", NoNewline)
		PrintObject(DirObj)
		Print(".", Newline)
		return true
	}
	if !DirObj.Has(FlgSurf) && DirObj.Capacity != 0 {
		if DirObj.Has(FlgOpen) {
			DirObj.Take(FlgOpen)
			Print("Closed.", Newline)
			if Lit {
				Lit = IsLit(Here, true)
				if !Lit {
					Print("It is now pitch black.", Newline)
				}
			}
			return true
		}
		Print("It is already closed.", Newline)
		return true
	}
	if DirObj.Has(FlgDoor) {
		if DirObj.Has(FlgOpen) {
			DirObj.Take(FlgOpen)
			Print("The ", NoNewline)
			PrintObject(DirObj)
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
	if DirObj.Has(FlgPerson) {
		Print("The ", NoNewline)
		PrintObject(DirObj)
		Print(" pays no attention.", Newline)
		return true
	}
	Print("You cannot talk to that!", Newline)
	return true
}

func VCount(arg ActArg) bool {
	if DirObj == &Blessings {
		Print("Well, for one, you are playing Zork...", Newline)
		return true
	}
	Print("You have lost your mind.", Newline)
	return true
}

func VCross(arg ActArg) bool {
	Print("You can't cross that!", Newline)
	return true
}

func VCurses(arg ActArg) bool {
	if DirObj == nil {
		Print("Such language in a high-class establishment like this!", Newline)
		return true
	}
	if DirObj.Has(FlgPerson) {
		Print("Insults of this nature won't help you.", Newline)
		return true
	}
	Print("What a loony!", Newline)
	return true
}

func VCut(arg ActArg) bool {
	if DirObj.Has(FlgPerson) {
		// TODO: returns perform
		if Perform(ActionVerb{Norm: "attack", Orig: "attack"}, DirObj, IndirObj) == PerfHndld {
			return true
		}
		return false
	}
	if DirObj.Has(FlgBurn) && IndirObj.Has(FlgWeapon) {
		if Winner.IsIn(DirObj) {
			Print("Not a bright idea, especially since you're in it.", Newline)
			return true
		}
		RemoveCarefully(DirObj)
		Print("Your skillful ", NoNewline)
		PrintObject(IndirObj)
		Print("smanship slices the ", NoNewline)
		PrintObject(DirObj)
		Print(" into innumerable slivers which blow away.", Newline)
		return true
	}
	if !IndirObj.Has(FlgWeapon) {
		Print("The \"cutting edge\" of a ", NoNewline)
		PrintObject(IndirObj)
		Print(" is hardly adequate.", Newline)
		return true
	}
	Print("Strange concept, cutting the ", NoNewline)
	PrintObject(DirObj)
	Print("....", Newline)
	return true
}

func VDeflate(arg ActArg) bool {
	Print("Come on, now!", Newline)
	return true
}

func VDig(arg ActArg) bool {
	if IndirObj == nil {
		IndirObj = &Hands
	}
	if IndirObj == &Shovel {
		Print("There's no reason to be digging here.", Newline)
		return true
	}
	if IndirObj.Has(FlgTool) {
		Print("Digging with the ", NoNewline)
		PrintObject(IndirObj)
		Print(" is slow and tedious.", Newline)
		return true
	}
	Print("Digging with a ", NoNewline)
	PrintObject(IndirObj)
	Print(" is silly.", Newline)
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
		// TODO: return fatal
		return false
	}
	if Here.Has(FlgLand) {
		Print("You are on your own feet again.", Newline)
		Winner.MoveTo(Here)
		return true
	}
	Print("You realize that getting out here would be fatal.", Newline)
	// TODO: return fatal
	return false
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
	if DirObj == Winner.Location() {
		Perform(ActionVerb{Norm: "disembark", Orig: "disembark"}, DirObj, nil)
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
	isEat := DirObj.Has(FlgFood)
	if isEat {
		if !DirObj.IsIn(Winner) && !DirObj.Location().IsIn(Winner) {
			Print("You're not holding that.", Newline)
			return true
		}
		if ActVerb.Norm == "drink" {
			Print("How can you drink that?", Newline)
			return true
		}
		Print("Thank you very much. It really hit the spot.", Newline)
		RemoveCarefully(DirObj)
		return true
	}
	isDrink := DirObj.Has(FlgDrink)
	if isDrink {
		nobj := DirObj.Location()
		if DirObj.IsIn(&GlobalObjects) || IsInGlobal(&GlobalWater, Here) || DirObj == &PseudoObject {
			return HitSpot()
		}
		if nobj == nil || !IsAccessible(nobj) {
			Print("There isn't any water here.", Newline)
			return true
		}
		if IsAccessible(nobj) && !nobj.IsIn(Winner) {
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
		PrintObject(DirObj)
		Print(" would agree with you.", Newline)
		return true
	}
	return false
}

func HitSpot() bool {
	if DirObj == &Water && !IsInGlobal(&GlobalWater, Here) {
		RemoveCarefully(DirObj)
	}
	Print("Thank you very much. I was rather thirsty (from all this talking, probably).", Newline)
	return true
}

func VEcho(arg ActArg) bool {
	if len(LexRes) <= 0 {
		Print("echo echo ...", Newline)
		return true
	}
	wrd := LexRes[len(LexRes)-1]
	Print(wrd.Orig+" "+wrd.Orig+" ...", Newline)
	return true
}

func VEnchant(arg ActArg) bool {
	return VDisenchant(ActUnk)
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

func VExcorcise(arg ActArg) bool {
	Print("What a bizarre concept!", Newline)
	return true
}

func PreFill(arg ActArg) bool {
	if IndirObj == nil {
		if IsInGlobal(&GlobalWater, Here) {
			Perform(ActionVerb{Norm: "fill", Orig: "fill"}, DirObj, &GlobalWater)
			return true
		}
		if Water.IsIn(Winner.Location()) {
			Perform(ActionVerb{Norm: "fill", Orig: "fill"}, DirObj, &Water)
			return true
		}
		Print("There is nothing to fill it with.", Newline)
		return true
	}
	if IndirObj == &Water {
		return false
	}
	if IndirObj != &GlobalWater {
		Perform(ActionVerb{Norm: "put", Orig: "put"}, IndirObj, DirObj)
		return true
	}
	return false
}

func VFill(arg ActArg) bool {
	if IndirObj != nil {
		Print("You may know how to do that, but I don't.", Newline)
		return true
	}
	if IsInGlobal(&GlobalWater, Here) {
		Perform(ActionVerb{Norm: "fill", Orig: "fill"}, DirObj, &GlobalWater)
		return true
	}
	if Water.IsIn(Winner.Location()) {
		Perform(ActionVerb{Norm: "fill", Orig: "fill"}, DirObj, &Water)
		return true
	}
	Print("There is nothing to fill it with.", Newline)
	return true
}

func VFirstLook(arg ActArg) bool {
	if DescribeRoom(false) {
		if !SuperBrief {
			return DescribeObjects(false)
		}
	}
	return false
}

func VFind(arg ActArg) bool {
	if DirObj == &Hands || DirObj == &Lungs {
		Print("Within six feet of your head, assuming you haven't left that somewhere.", Newline)
		return true
	}
	if DirObj == &Me {
		Print("You're around here somewhere...", Newline)
		return true
	}
	l := DirObj.Location()
	if l == &GlobalObjects {
		Print("You find it.", Newline)
		return true
	}
	if DirObj.IsIn(Winner) {
		Print("You have it.", Newline)
		return true
	}
	if DirObj.IsIn(Here) || IsInGlobal(DirObj, Here) || DirObj == &PseudoObject {
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

func VFollow(arg ActArg) bool {
	Print("You're nuts!", Newline)
	return true
}

func VFrobozz(arg ActArg) bool {
	Print("The FROBOZZ Corporation created, owns, and operates this dungeon.", Newline)
	return true
}

func PreGive(arg ActArg) bool {
	if !IsHeld(DirObj) {
		Print("That's easy for you to say since you don't even have the ", NoNewline)
		PrintObject(DirObj)
		Print(".", Newline)
		return true
	}
	return false
}

func VGive(arg ActArg) bool {
	if !DirObj.Has(FlgPerson) {
		Print("You can't give a ", NoNewline)
		PrintObject(DirObj)
		Print(" to a ", NoNewline)
		PrintObject(IndirObj)
		Print("!", Newline)
		return true
	}
	Print("The ", NoNewline)
	PrintObject(IndirObj)
	Print(" refuses it politely.", Newline)
	return true
}

func VHatch(arg ActArg) bool {
	Print("Bizarre!", Newline)
	return true
}

func VHello(arg ActArg) bool {
	if DirObj == nil {
		Print(PickOne(Hellos), Newline)
		return true
	}
	if DirObj.Has(FlgPerson) {
		Print("The ", NoNewline)
		PrintObject(DirObj)
		Print(" bows his head to you in greeting.", Newline)
		return true
	}
	Print("It's a well known fact that only schizophrenics say \"Hello\" to a ", NoNewline)
	PrintObject(DirObj)
	Print(".", Newline)
	return true
}

func VIncant(arg ActArg) bool {
	Print("The incantation echoes back faintly, but nothing else happens.", Newline)
	Params.InQuotes = false
	Params.Continue = NumUndef
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
	if DirObj.Has(FlgDoor) {
		Print("Nobody's home.", Newline)
		return true
	}
	Print("Why knock on a ", NoNewline)
	PrintObject(DirObj)
	Print("?", Newline)
	return true
}

func VLampOff(arg ActArg) bool {
	if !DirObj.Has(FlgLight) {
		Print("You can't turn that off.", Newline)
		return true
	}
	if !DirObj.Has(FlgOn) {
		Print("It is already off.", Newline)
		return true
	}
	DirObj.Take(FlgOn)
	if Lit {
		Lit = IsLit(Here, true)
	}
	Print("The ", NoNewline)
	PrintObject(DirObj)
	Print(" is now off.", Newline)
	if !Lit {
		Print("It is now pitch black.", Newline)
	}
	return true
}

func VLampOn(arg ActArg) bool {
	if !DirObj.Has(FlgLight) {
		if DirObj.Has(FlgBurn) {
			Print("If you wish to burn the ", NoNewline)
			PrintObject(DirObj)
			Print(", you should say so.", Newline)
			return true
		}
		Print("You can't turn that on.", Newline)
		return true
	}
	if DirObj.Has(FlgOn) {
		Print("It is already on.", Newline)
		return true
	}
	DirObj.Give(FlgOn)
	Print("The ", NoNewline)
	PrintObject(DirObj)
	Print(" is now on.", Newline)
	if !Lit {
		Lit = IsLit(Here, true)
		NewLine()
		return VLook(ActUnk)
	}
	return true
}

func VLaunch(arg ActArg) bool {
	if DirObj.Has(FlgVeh) {
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

func VListen(arg ActArg) bool {
	Print("The ", NoNewline)
	PrintObject(DirObj)
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
	PrintObject(DirObj)
	Print(".", Newline)
	return true
}

func VLookInside(arg ActArg) bool {
	if DirObj.Has(FlgDoor) {
		Print("The ", NoNewline)
		PrintObject(DirObj)
		if DirObj.Has(FlgOpen) {
			Print(" is open, but I can't tell what's beyond it.", NoNewline)
		} else {
			Print(" is closed.", NoNewline)
		}
		NewLine()
		return true
	}
	if DirObj.Has(FlgCont) {
		if DirObj.Has(FlgPerson) {
			Print("There is nothing special to be seen.", Newline)
			return true
		}
		if !CanSeeInside(DirObj) {
			Print("The ", NoNewline)
			PrintObject(DirObj)
			Print(" is closed.", Newline)
			return true
		}
		if DirObj.HasChildren() && PrintCont(DirObj, false, 0) {
			return true
		}
		Print("The ", NoNewline)
		PrintObject(DirObj)
		Print(" is empty.", Newline)
		return true
	}
	Print("You can't look inside a ", NoNewline)
	PrintObject(DirObj)
	Print(".", Newline)
	return true
}

func VLookOn(arg ActArg) bool {
	if DirObj.Has(FlgSurf) {
		Perform(ActionVerb{Norm: "look inside", Orig: "look inside"}, DirObj, nil)
		return true
	}
	Print("Look on a ", NoNewline)
	PrintObject(DirObj)
	Print("???", Newline)
	return true
}

func VLookUnder(arg ActArg) bool {
	Print("There is nothing but dust there.", Newline)
	return true
}

func VExamine(arg ActArg) bool {
	if len(DirObj.Text) > 0 {
		Print(DirObj.Text, Newline)
		return true
	}
	if DirObj.Has(FlgCont) || DirObj.Has(FlgDoor) {
		return VLookInside(ActUnk)
	}
	Print("There's nothing special about the ", NoNewline)
	PrintObject(DirObj)
	Print(".", Newline)
	return true
}

func VMake(arg ActArg) bool {
	Print("You can't do that.", Newline)
	return true
}

func VMelt(arg ActArg) bool {
	Print("It's not clear that a ", NoNewline)
	PrintObject(DirObj)
	Print(" can be melted.", Newline)
	return true
}

func PreMove(arg ActArg) bool {
	if IsHeld(DirObj) {
		Print("You aren't an accomplished enough juggler.", Newline)
		return true
	}
	return false
}

func VMove(arg ActArg) bool {
	if DirObj.Has(FlgTake) {
		Print("Moving the ", NoNewline)
		PrintObject(DirObj)
		Print(" reveals nothing.", Newline)
		return true
	}
	Print("You can't move the ", NoNewline)
	PrintObject(DirObj)
	Print(".", Newline)
	return true
}

func VMumble(arg ActArg) bool {
	Print("You'll have to speak up if you expect me to hear you!", Newline)
	return true
}

func PreMung(arg ActArg) bool {
	if IndirObj == nil || !IndirObj.Has(FlgWeapon) {
		Print("Trying to destroy the ", NoNewline)
		PrintObject(DirObj)
		Print(" with ", NoNewline)
		if IndirObj == nil {
			Print("your bare hands", NoNewline)
		} else {
			Print("a", NoNewline)
			PrintObject(IndirObj)
		}
		Print(" is futile.", Newline)
		return true
	}
	return false
}

func VMung(arg ActArg) bool {
	if !DirObj.Has(FlgPerson) {
		Print("Nice try.", Newline)
		return true
	}
	Perform(ActionVerb{Norm: "attack", Orig: "attack"}, DirObj, nil)
	return true
}

func MungRoom(rm *Object, str string) {
	rm.Give(FlgKludge)
	rm.LongDesc = str
}

func VOdysseus(arg ActArg) bool {
	if Here != &CyclopsRoom || !Cyclops.IsIn(Here) || CyclopsFlag {
		Print("Wasn't he a sailor?", Newline)
		return true
	}
	QueueInt(ICyclops, false).Run = false
	CyclopsFlag = true
	Print("The cyclops, hearing the name of his father's deadly nemesis, flees the room by knocking down the wall on the east of the room.", Newline)
	MagicFlag = true
	Cyclops.Take(FlgFight)
	return RemoveCarefully(&Cyclops)
}

func VOil(arg ActArg) bool {
	Print("You probably put spinach in your gas tank, too.", Newline)
	return true
}

func VOpen(arg ActArg) bool {
	if !DirObj.Has(FlgCont) || DirObj.Capacity == 0 {
		if DirObj.Has(FlgDoor) {
			if DirObj.Has(FlgOpen) {
				Print("It is already open.", Newline)
				return true
			}
			Print("The ", NoNewline)
			PrintObject(DirObj)
			Print(" opens.", Newline)
			DirObj.Give(FlgOpen)
			return true
		}
		Print("You must tell me how to do that to a ", NoNewline)
		PrintObject(DirObj)
		Print(".", Newline)
		return true
	}
	if DirObj.Has(FlgOpen) {
		Print("It is already open.", Newline)
		return true
	}
	DirObj.Give(FlgOpen)
	DirObj.Give(FlgTouch)
	if !DirObj.HasChildren() || DirObj.Has(FlgTrans) {
		Print("Opened.", Newline)
		return true
	}
	if len(DirObj.Children) == 1 && !DirObj.Children[0].Has(FlgTouch) && len(DirObj.Children[0].FirstDesc) > 0 {
		str := DirObj.Children[0].FirstDesc
		Print("The ", NoNewline)
		PrintObject(DirObj)
		Print(" opens.", Newline)
		Print(str, Newline)
		return true
	}
	Print("Opening the ", NoNewline)
	PrintObject(DirObj)
	Print(" reveals ", NoNewline)
	PrintContents(DirObj)
	Print(".", Newline)
	return true
}

func VOverboard(arg ActArg) bool {
	locn := Winner.Location()
	if IndirObj == &Teeth {
		if locn.Has(FlgVeh) {
			Print("Ahoy -- ", NoNewline)
			PrintObject(IndirObj)
			Print(" overboard!", Newline)
			return true
		}
		Print("You're not in anything!", Newline)
		return true
	}
	if locn.Has(FlgVeh) {
		Perform(ActionVerb{Norm: "throw", Orig: "throw"}, DirObj, nil)
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
	if !DirObj.Has(FlgPerson) {
		Print("That's silly!", Newline)
		return true
	}
	Print("You become so engrossed in the role of the ", NoNewline)
	PrintObject(DirObj)
	Print(" that you kill yourself, just as he might have done!", Newline)
	return JigsUp("", false)
}

func VPlug(arg ActArg) bool {
	Print("This has no effect.", Newline)
	return true
}

func VPourOn(arg ActArg) bool {
	if DirObj == &Water {
		RemoveCarefully(DirObj)
		if IsFlaming(IndirObj) {
			Print("The ", NoNewline)
			PrintObject(IndirObj)
			Print(" is extinguished.", Newline)
			IndirObj.Take(FlgOn)
			IndirObj.Take(FlgFlame)
			return true
		}
		Print("The water spills over the ", NoNewline)
		PrintObject(IndirObj)
		Print(", to the floor, and evaporates.", Newline)
		return true
	}
	if DirObj == &Putty {
		if Perform(ActionVerb{Norm: "put", Orig: "put"}, &Putty, IndirObj) == PerfHndld {
			return true
		}
		return false
	}
	Print("You can't pour that.", Newline)
	return true
}

func VPray(arg ActArg) bool {
	if Here != &SouthTemple {
		Print("If you pray enough, your prayers may be answered.", Newline)
		return true
	}
	return Goto(&Forest1, true)
}

func VPump(arg ActArg) bool {
	if IndirObj != nil && IndirObj != &Pump {
		Print("Pump it up with a ", NoNewline)
		PrintObject(IndirObj)
		Print("?", Newline)
		return true
	}
	if Pump.IsIn(Winner) {
		if Perform(ActionVerb{Norm: "inflate", Orig: "inflate"}, DirObj, &Pump) == PerfHndld {
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
	if !IndirObj.Has(FlgOpen) && !IsOpenable(IndirObj) && !IndirObj.Has(FlgVeh) {
		Print("You can't do that.", Newline)
		return true
	}
	if !IndirObj.Has(FlgOpen) {
		Print("The ", NoNewline)
		PrintObject(IndirObj)
		Print(" isn't open.", Newline)
		ThisIsIt(IndirObj)
		return true
	}
	if IndirObj == DirObj {
		Print("How can you do that?", Newline)
		return true
	}
	if DirObj.IsIn(IndirObj) {
		Print("The ", NoNewline)
		PrintObject(DirObj)
		Print(" is already in the ", NoNewline)
		PrintObject(IndirObj)
		Print(".", Newline)
		return true
	}
	if Weight(IndirObj)+Weight(DirObj)-IndirObj.Size > IndirObj.Capacity {
		Print("There's no room.", Newline)
		return true
	}
	if !IsHeld(DirObj) && !ITake(true) {
		return true
	}
	DirObj.MoveTo(IndirObj)
	DirObj.Give(FlgTouch)
	ScoreObj(DirObj)
	Print("Done.", Newline)
	return true
}

func VWalkAround(arg ActArg) bool {
	Print("Use compass directions for movement.", Newline)
	return true
}

func VPutBehind(arg ActArg) bool {
	Print("That hiding place is too obvious.", Newline)
	return true
}

func VPutOn(arg ActArg) bool {
	if IndirObj == nil || IndirObj == &Ground {
		return VDrop(ActUnk)
	}
	if IndirObj.Has(FlgSurf) {
		return VPut(ActUnk)
	}
	Print("There's no good surface on the ", NoNewline)
	PrintObject(IndirObj)
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
	if !Lit {
		Print("It is impossible to read in the dark.", Newline)
		return true
	}
	if IndirObj != nil && !IndirObj.Has(FlgTrans) {
		Print("How does one look through a ", NoNewline)
		PrintObject(IndirObj)
		Print("?", Newline)
		return true
	}
	return false
}

func VRead(arg ActArg) bool {
	if !DirObj.Has(FlgRead) {
		Print("How does one read a ", NoNewline)
		PrintObject(DirObj)
		Print("?", Newline)
		return true
	}
	Print(DirObj.Text, Newline)
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
	PrintObject(DirObj)
	Print(" is interested.", Newline)
	Params.Continue = NumUndef
	Params.InQuotes = false
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
	if Params.Continue == NumUndef {
		Print("Say what?", Newline)
		return true
	}
	Params.InQuotes = false
	v := FindIn(Here, FlgPerson)
	if v != nil {
		Print("You must address the ", NoNewline)
		PrintObject(v)
		Print(" directly.", Newline)
		Params.Continue = NumUndef
		return true
	}
	if LexRes[Params.Continue].Norm == "hello" {
		Params.Continue = NumUndef
		Print("Talking to yourself is a sign of impending mental collapse.", Newline)
		return true
	}
	return false
}

func FindIn(where *Object, what Flag) *Object {
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
	if DirObj.Has(FlgPerson) {
		Print("Why would you send for the ", NoNewline)
		PrintObject(DirObj)
		Print("?", Newline)
		return true
	}
	Print("That doesn't make sends.", Newline)
	return true
}

func PreSGive(arg ActArg) bool {
	Perform(ActionVerb{Norm: "give", Orig: "give"}, IndirObj, DirObj)
	return true
}

func VSGive(arg ActArg) bool {
	Print("Foo!", Newline)
	return true
}

func VShake(arg ActArg) bool {
	if DirObj.Has(FlgPerson) {
		Print("This seems to have no effect.", Newline)
		return true
	}
	if !DirObj.Has(FlgTake) {
		Print("You can't take it; thus, you can't shake it!", Newline)
		return true
	}
	if !DirObj.Has(FlgCont) {
		Print("Shaken.", Newline)
		return true
	}
	if DirObj.Has(FlgOpen) {
		if !DirObj.HasChildren() {
			Print("Shaken.", Newline)
			return true
		}
		ShakeLoop()
		Print("The contents of the ", NoNewline)
		PrintObject(DirObj)
		Print(" spill ", NoNewline)
		if !Here.Has(FlgLand) {
			Print("out and disappears", NoNewline)
		} else {
			Print("to the ground", NoNewline)
		}
		Print(".", Newline)
		return true
	}
	if DirObj.HasChildren() {
		Print("It sounds like there is something inside the ", NoNewline)
		PrintObject(DirObj)
		Print(".", Newline)
		return true
	}
	Print("The ", NoNewline)
	PrintObject(DirObj)
	Print(" sounds empty.", Newline)
	return true
}

func ShakeLoop() {
	if !DirObj.HasChildren() {
		return
	}
	x := DirObj.Children[0]
	x.Give(FlgTouch)
	mv := Here
	if Here == &UpATree {
		mv = &Path
	} else if !Here.Has(FlgLand) {
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
	PrintObject(DirObj)
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
	if !DirObj.Has(FlgPerson) {
		Print("How singularly useless.", Newline)
		return true
	}
	Print("The ", NoNewline)
	PrintObject(DirObj)
	Print(" does not understand this.", Newline)
	return true
}

func VSSpray(arg ActArg) bool {
	if Perform(ActionVerb{Norm: "spray", Orig: "spray"}, IndirObj, DirObj) == PerfHndld {
		return true
	}
	return false
}

func VStab(arg ActArg) bool {
	w := FindWeapon(Winner)
	if w == nil {
		Print("No doubt you propose to stab the ", NoNewline)
		PrintObject(DirObj)
		Print(" with your pinky?", Newline)
		return true
	}
	Perform(ActionVerb{Norm: "attack", Orig: "attack"}, DirObj, w)
	return true
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

func VStrike(arg ActArg) bool {
	if DirObj.Has(FlgPerson) {
		Print("Since you aren't versed in hand-to-hand combat, you'd better attack the ", NoNewline)
		PrintObject(DirObj)
		Print(" with a weapon.", Newline)
		return true
	}
	Perform(ActionVerb{Norm: "lamp on", Orig: "lamp on"}, DirObj, nil)
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

func VSwing(arg ActArg) bool {
	if IndirObj == nil {
		Print("Whoosh!", Newline)
		return true
	}
	if Perform(ActionVerb{Norm: "attack", Orig: "attack"}, IndirObj, DirObj) == PerfHndld {
		return true
	}
	return false
}

func PreTake(arg ActArg) bool {
	if DirObj == Winner {
		if DirObj.Has(FlgWear) {
			Print("You are already wearing it.", Newline)
			return true
		}
		Print("You already have that!", Newline)
		return true
	}
	lcn := DirObj.Location()
	if lcn.Has(FlgCont) && !lcn.Has(FlgOpen) {
		Print("You can't reach something that's inside a closed container.", Newline)
		return true
	}
	if IndirObj != nil {
		if IndirObj == &Ground {
			IndirObj = nil
			return false
		}
		if IndirObj != DirObj.Location() {
			Print("The ", NoNewline)
			PrintObject(DirObj)
			Print(" isn't in the ", NoNewline)
			PrintObject(IndirObj)
			Print(".", Newline)
			return true
		}
		IndirObj = nil
		return false
	}
	if DirObj == Winner.Location() {
		Print("You're inside of it!", Newline)
		return true
	}
	return false
}

func VTake(arg ActArg) bool {
	if ITake(true) {
		if DirObj.Has(FlgWear) {
			Print("You are now wearing the ", NoNewline)
			PrintObject(DirObj)
			Print(".", Newline)
			return true
		}
		Print("Taken.", Newline)
		return true
	}
	return false
}

func VTell(arg ActArg) bool {
	if !DirObj.Has(FlgPerson) {
		Print("You can't talk to the ", NoNewline)
		PrintObject(DirObj)
		Print("!", Newline)
		Params.InQuotes = false
		Params.Continue = NumUndef
		// TODO: return fatal
		return false
	}
	if Params.Continue != NumUndef {
		Winner = DirObj
		Here = Winner.Location()
		return true
	}
	Print("The ", NoNewline)
	PrintObject(DirObj)
	Print(" pauses for a moment, perhaps thinking that you should reread the manual.", Newline)
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

func VThrow(arg ActArg) bool {
	if !IDrop() {
		Print("Huh?", Newline)
		return true
	}
	if IndirObj == &Me {
		Print("A terrific throw! The ", NoNewline)
		Winner = Player
		return JigsUp(" hits you squarely in the head. Normally, this wouldn't do much damage, but by incredible mischance, you fall over backwards trying to duck, and break your neck, justice being swift and merciful in the Great Underground Empire.", false)
	}
	if IndirObj != nil && IndirObj.Has(FlgPerson) {
		Print("The ", NoNewline)
		PrintObject(IndirObj)
		Print(" ducks as the ", NoNewline)
		PrintObject(DirObj)
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
	if IndirObj == Winner {
		Print("You can't tie anything to yourself.", Newline)
		return true
	}
	Print("You can't tie the ", NoNewline)
	PrintObject(DirObj)
	Print(" to that.", Newline)
	return true
}

func VTieUp(arg ActArg) bool {
	Print("You could certainly never tie it with that!", Newline)
	return true
}

func VTreasure(arg ActArg) bool {
	if Here == &NorthTemple {
		return Goto(&TreasureRoom, true)
	}
	if Here == &TreasureRoom {
		return Goto(&NorthTemple, true)
	}
	Print("Nothing happens.", Newline)
	return true
}

func PreTurn(arg ActArg) bool {
	if (IndirObj == nil || IndirObj == &Rooms) && DirObj != &Book {
		Print("Your bare hands don't appear to be enough.", Newline)
		return true
	}
	if !DirObj.Has(FlgTurn) {
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

func VWait(arg ActArg) bool {
	return Wait(3)
}

func Wait(num int) bool {
	Print("Time passes...", Newline)
	for i := num; i > 0; i-- {
		Clocker()
	}
	ClockWait = true
	return true
}

func VWave(arg ActArg) bool {
	return HackHack("Waving the ")
}

func VWear(arg ActArg) bool {
	if !DirObj.Has(FlgWear) {
		Print("You can't wear the ", NoNewline)
		PrintObject(DirObj)
		Print(".", Newline)
		return true
	}
	Perform(ActionVerb{Norm: "take", Orig: "take"}, DirObj, nil)
	return true
}

func VWin(arg ActArg) bool {
	Print("Naturally!", Newline)
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
				// TODO: return fatal
				return false
			}
			return JigsUp("Oh, no! You have walked into the slavering fangs of a lurking grue!", false)
		}
		Print("You can't go that way.", Newline)
		// TODO: return fatal
		return false
	}
	// Unconditional exit
	if props.UExit {
		return Goto(props.RExit, true)
	}
	// Non-exit
	if len(props.NExit) > 0 {
		Print(props.NExit, Newline)
		// TODO: return fatal
		return false
	}
	// Functional exit
	if props.FExit != nil {
		rm := props.FExit()
		if rm == nil {
			// TODO: return fatal
			return false
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
			// TODO: return fatal
			return false
		}
		Print("You can't go that way.", Newline)
		// TODO: return fatal
		return false
	}
	if props.DExit != nil {
		if props.DExit.Has(FlgOpen) {
			return Goto(props.RExit, true)
		}
		if len(props.DExitStr) > 0 {
			Print(props.DExitStr, Newline)
			// TODO: return fatal
			return false
		}
		Print("The ", NoNewline)
		PrintObject(props.DExit)
		Print(" is closed.", Newline)
		ThisIsIt(props.DExit)
		// TODO: return fatal
		return false
	}
	return false
}

func VWind(arg ActArg) bool {
	Print("You cannot wind up a ", NoNewline)
	PrintObject(DirObj)
	Print(".", Newline)
	return true
}

func VWish(arg ActArg) bool {
	Print("With luck, your wish will come true.", Newline)
	return true
}

func VYell(arg ActArg) bool {
	Print("Aaaarrrrgggghhhh!", Newline)
	return true
}

func VZork(arg ActArg) bool {
	Print("At your service!", Newline)
	return true
}

func NoGoTell(av Flag, wloc *Object) {
	if av != FlgUnk {
		Print("You can't go there in a ", NoNewline)
		PrintObject(wloc)
		Print(".", Newline)
		return
	}
	Print("You can't go there without a vehicle.", Newline)
}

func Goto(rm *Object, isV bool) bool {
	lb := rm.Has(FlgLand)
	wloc := Winner.Location()
	var av Flag
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
	// TODO: This statement should never be true
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

func VQuit(arg ActArg) bool {
	VScore(arg)
	Print("Do you wish to leave the game? (Y is affirmative): ", NoNewline)
	if IsYes() {
		Quit()
	} else {
		Print("Ok.", Newline)
	}
	return true
}

func RemoveCarefully(obj *Object) bool {
	if obj == Params.ItObj {
		Params.ItObj = nil
	}
	oLit := Lit
	obj.Remove()
	Lit = IsLit(Here, true)
	if oLit && oLit != Lit {
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
	if DirObj.IsIn(&GlobalObjects) && (ActVerb.Norm == "wave" || ActVerb.Norm == "raise" || ActVerb.Norm == "lower") {
		Print("The ", NoNewline)
		PrintObject(DirObj)
		Print(" isn't here!", Newline)
		return true
	}
	Print(str, NoNewline)
	PrintObject(DirObj)
	Print(PickOne(Hohum), Newline)
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

func IDrop() bool {
	if !DirObj.IsIn(Winner) && !DirObj.Location().IsIn(Winner) {
		Print("You're not carrying the ", NoNewline)
		PrintObject(DirObj)
		Print(".", Newline)
		return false
	}
	if !DirObj.IsIn(Winner) && !DirObj.Location().Has(FlgOpen) {
		Print("The ", NoNewline)
		PrintObject(DirObj)
		Print(" is closed.", Newline)
		return false
	}
	DirObj.MoveTo(Winner.Location())
	return true
}

func IsYes() bool {
	Print(">", NoNewline)
	_, lex := Read()
	if len(lex) > 0 && lex[0].IsAny("yes", "y") {
		return true
	}
	return false
}

func DescribeRoom(isLook bool) bool {
	v := Verbose
	if isLook {
		v = isLook
	}
	if !Lit {
		Print("It is pitch black.", NoNewline)
		if !IsSprayed {
			Print(" You are likely to be eaten by a grue.", NoNewline)
		}
		NewLine()
		return false
	}
	if !Here.Has(FlgTouch) {
		Here.Give(FlgTouch)
		v = true
	}
	if Here.Has(FlgMaze) {
		Here.Take(FlgTouch)
	}
	if Here.IsIn(&Rooms) {
		PrintObject(Here)
		if av := Winner.Location(); av.Has(FlgVeh) {
			Print(", in the ", NoNewline)
			PrintObject(av)
		}
		NewLine()
	}
	if !isLook && SuperBrief {
		return true
	}
	av := Winner.Location()
	if v && Here.Action != nil && Here.Action(ActLook) {
		return true
	}
	if v && len(Here.LongDesc) != 0 {
		Print(Here.LongDesc, Newline)
	} else if Here.Action != nil {
		Here.Action(ActFlash)
	}
	if Here != av && av.Has(FlgVeh) && av.Action != nil {
		av.Action(ActLook)
	}
	return true
}

func DescribeObjects(v bool) bool {
	if !Lit {
		Print("Only bats can see in the dark. And you're not one.", Newline)
		return true
	}
	if !Here.HasChildren() {
		return false
	}
	if !v {
		v = Verbose
	}
	return PrintCont(Here, v, -1)
}

func PrintCont(obj *Object, v bool, lvl int) bool {
	if !obj.HasChildren() {
		return true
	}
	var av *Object
	if Winner.Location().Has(FlgVeh) {
		av = Winner.Location()
	}
	isInv := false
	isPv := false
	shit := true
	isFirst := true
	if Winner == obj || Winner == obj.Location() {
		isInv = true
	} else {
		for _, child := range obj.Children {
			if child == av {
				isPv = true
			} else if Winner == child {
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
	if obj == Winner {
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
	DescObj = obj
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
		} else if obj.Has(FlgWear) && obj.IsIn(Winner) {
			Print(" (being worn)", NoNewline)
		}
	}
	if av := Winner.Location(); lvl == 0 && av != nil && av.Has(FlgVeh) {
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
	Params.ItObj = obj
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
		if obj == Winner {
			return true
		}
	}
}

func ScoreUpd(num int) bool {
	BaseScore += num
	Score += num
	if Score == 350 && !WonGame {
		WonGame = true
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
		if obj == Player && child.Has(FlgWear) {
			wt++
		} else {
			wt += Weight(child)
		}
	}
	return wt + obj.Size
}

func ITake(vb bool) bool {
	if Dead {
		if vb {
			Print("Your hand passes through its object.", Newline)
		}
		return false
	}
	if DirObj == nil {
		return false
	}
	if !DirObj.Has(FlgTake) {
		if vb {
			Print(PickOne(Yuks), Newline)
		}
		return false
	}
	if DirObj.Has(FlgCont) && !DirObj.Has(FlgOpen) {
		return false
	}
	if !DirObj.Location().IsIn(Winner) && Weight(DirObj)+Weight(Winner) > LoadAllowed {
		if vb {
			Print("Your load is too heavy", NoNewline)
			if LoadAllowed < LoadMax {
				Print(", especially in light of your condition.", NoNewline)
			} else {
				Print(".", NoNewline)
			}
			NewLine()
		}
		// TODO: rfatal, not rfalse
		return false
	}
	cnt := CCount(Winner)
	if ActVerb.Norm == "tell" && cnt > FumbleNumber && Prob(cnt*FumbleProb, false) {
		Print("You're holding too many things already!", Newline)
		return false
	}
	DirObj.MoveTo(Winner)
	DirObj.Take(FlgNoDesc)
	DirObj.Give(FlgTouch)
	ScoreObj(DirObj)
	return true
}

func Finish() bool {
	VScore(ActUnk)
	for {
		NewLine()
		Print("Would you like to restart the game from the beginning, restore a saved game position, or end this session of the game?", Newline)
		Print("(Type RESTART, RESTORE, or QUIT):", Newline)
		Print(">", NoNewline)
		_, lex := Read()
		if len(lex) == 0 {
			continue
		}
		wrd := lex[0]
		if wrd.Norm == "restart" {
			if !Restart() {
				Print("Failed.", Newline)
			}
			return true
		}
		if wrd.Norm == "restore" {
			if Restore() {
				Print("Ok.", Newline)
				return true
			}
			Print("Failed.", Newline)
			return true
		}
		if wrd.Norm == "quit" {
			Quit()
		}
	}
}

func Lkp(itm *Object, tbl []*Object) *Object {
	for idx, obj := range tbl {
		if obj == itm {
			if idx+1 <= len(tbl)-1 {
				return tbl[idx+1]
			}
			break
		}
	}
	return nil
}

func ToDirObj(dir string) *Object {
	return nil
}
