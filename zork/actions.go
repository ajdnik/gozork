package zork

import "math/rand"

func OpenClose(obj *Object, strOpn, strCls string) {
	if ActVerb.Norm == "open" {
		if obj.Has(FlgOpen) {
			Print(PickOne(Dummy), Newline)
		} else {
			Print(strOpn, Newline)
			obj.Give(FlgOpen)
		}
	} else if ActVerb.Norm == "close" {
		if obj.Has(FlgOpen) {
			Print(strCls, Newline)
			obj.Take(FlgOpen)
		} else {
			Print(PickOne(Dummy), Newline)
		}
	}
}

func LeavesAppear() bool {
	if !Grate.Has(FlgOpen) && !GrateRevealed {
		if ActVerb.Norm == "move" || ActVerb.Norm == "take" {
			Print("In disturbing the pile of leaves, a grating is revealed.", Newline)
		} else {
			Print("With the leaves moved, a grating is revealed.", Newline)
		}
		Grate.Take(FlgInvis)
		GrateRevealed = true
	}
	return false
}

func Fweep(n int) {
	for i := 0; i < n; i++ {
		Print("    Fweep!", Newline)
	}
	NewLine()
}

func FlyMe() bool {
	Fweep(4)
	Print("The bat grabs you by the scruff of your neck and lifts you away....", Newline)
	NewLine()
	dest := BatDrops[rand.Intn(len(BatDrops))]
	Goto(dest, false)
	if Here != &EnteranceToHades {
		VFirstLook(ActUnk)
	}
	return true
}

func TouchAll(obj *Object) {
	for _, child := range obj.Children {
		child.Give(FlgTouch)
		if child.HasChildren() {
			TouchAll(child)
		}
	}
}

func OtvalFrob(o *Object) int {
	score := 0
	for _, child := range o.Children {
		score += child.TValue
		if child.HasChildren() {
			score += OtvalFrob(child)
		}
	}
	return score
}

func IntegralPart() {
	Print("It is an integral part of the control panel.", Newline)
}

func WithTell(obj *Object) {
	Print("With a ", NoNewline)
	PrintObject(obj)
	Print("?", Newline)
}


func BadEgg() {
	if Canary.IsIn(&Egg) {
		Print(" ", NoNewline)
		Print(BrokenCanary.FirstDesc, NoNewline)
	} else {
		RemoveCarefully(&BrokenCanary)
	}
	BrokenEgg.MoveTo(Egg.Location())
	RemoveCarefully(&Egg)
}

func Slider(obj *Object) {
	if obj.Has(FlgTake) {
		Print("The ", NoNewline)
		PrintObject(obj)
		Print(" falls into the slide and is gone.", Newline)
		if obj == &Water {
			RemoveCarefully(obj)
		} else {
			obj.MoveTo(&Cellar)
		}
	} else {
		Print(PickOne(Yuks), Newline)
	}
}

func ForestRoomQ() bool {
	return Here == &Forest1 || Here == &Forest2 || Here == &Forest3 ||
		Here == &Path || Here == &UpATree
}

// StealJunk and DropJunk are defined in the IThief section below

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func AbsInt(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// ================================================================
// SCORE
// ================================================================

func VScore(arg ActArg) bool {
	Print("Your score is ", NoNewline)
	PrintNumber(Score)
	Print(" (total of 350 points), in ", NoNewline)
	PrintNumber(Moves)
	if Moves == 1 {
		Print(" move.", NoNewline)
	} else {
		Print(" moves.", NoNewline)
	}
	NewLine()
	Print("This gives you the rank of ", NoNewline)
	switch {
	case Score == 350:
		Print("Master Adventurer", NoNewline)
	case Score > 330:
		Print("Wizard", NoNewline)
	case Score > 300:
		Print("Master", NoNewline)
	case Score > 200:
		Print("Adventurer", NoNewline)
	case Score > 100:
		Print("Junior Adventurer", NoNewline)
	case Score > 50:
		Print("Novice Adventurer", NoNewline)
	case Score > 25:
		Print("Amateur Adventurer", NoNewline)
	default:
		Print("Beginner", NoNewline)
	}
	Print(".", Newline)
	return true
}

func VDiagnose(arg ActArg) bool {
	ms := FightStrength(false)
	wd := Winner.Strength
	rs := ms + wd
	// Check if healing is active
	cureActive := false
	for i := len(QueueItms) - 1; i >= 0; i-- {
		if QueueItms[i].Rtn != nil && QueueItms[i].Run {
			// Can't directly compare function pointers in Go, so we use a flag
			// The cure interrupt is identified by its tick pattern
			break
		}
	}
	if !cureActive {
		wd = 0
	} else {
		wd = -wd
	}
	if wd == 0 {
		Print("You are in perfect health.", NoNewline)
	} else {
		Print("You have ", NoNewline)
		switch {
		case wd == 1:
			Print("a light wound,", NoNewline)
		case wd == 2:
			Print("a serious wound,", NoNewline)
		case wd == 3:
			Print("several wounds,", NoNewline)
		default:
			Print("serious wounds,", NoNewline)
		}
	}
	if wd != 0 {
		Print(" which will be cured after some moves.", NoNewline)
	}
	NewLine()
	Print("You can ", NoNewline)
	switch {
	case rs == 0:
		Print("expect death soon", NoNewline)
	case rs == 1:
		Print("be killed by one more light wound", NoNewline)
	case rs == 2:
		Print("be killed by a serious wound", NoNewline)
	case rs == 3:
		Print("survive one serious wound", NoNewline)
	default:
		Print("survive several wounds", NoNewline)
	}
	Print(".", Newline)
	if Deaths > 0 {
		Print("You have been killed ", NoNewline)
		if Deaths == 1 {
			Print("once", NoNewline)
		} else {
			Print("twice", NoNewline)
		}
		Print(".", Newline)
	}
	return true
}

// ================================================================
// DEATH AND RESTART
// ================================================================

func JigsUp(desc string, isPlyr bool) bool {
	Winner = &Adventurer
	if Dead {
		NewLine()
		Print("It takes a talented person to be killed while already dead. YOU are such a talent. Unfortunately, it takes a talented person to deal with it. I am not such a talent. Sorry.", Newline)
		return Finish()
	}
	Print(desc, Newline)
	if !Lucky {
		Print("Bad luck, huh?", Newline)
	}
	ScoreUpd(-10)
	NewLine()
	Print("    ****  You have died  ****", Newline)
	NewLine()
	if Winner.Location().Has(FlgVeh) {
		Winner.MoveTo(Here)
	}
	if Deaths >= 2 {
		Print("You clearly are a suicidal maniac. We don't allow psychotics in the cave, since they may harm other adventurers. Your remains will be installed in the Land of the Living Dead, where your fellow adventurers may gloat over them.", Newline)
		return Finish()
	}
	Deaths++
	Winner.MoveTo(Here)
	if SouthTemple.Has(FlgTouch) {
		Print("As you take your last breath, you feel relieved of your burdens. The feeling passes as you find yourself before the gates of Hell, where the spirits jeer at you and deny you entry. Your senses are disturbed. The objects in the dungeon appear indistinct, bleached of color, even unreal.", Newline)
		NewLine()
		Dead = true
		TrollFlag = true
		AlwaysLit = true
		Winner.Action = DeadFunction
		Goto(&EnteranceToHades, true)
	} else {
		Print("Now, let's take a look here...\nWell, you probably deserve another chance. I can't quite fix you up completely, but you can't have everything.", Newline)
		NewLine()
		Goto(&Forest1, true)
	}
	TrapDoor.Take(FlgTouch)
	Params.Continue = NumUndef
	RandomizeObjects()
	KillInterrupts()
	return false
}

func RandomizeObjects() {
	if Lamp.IsIn(Winner) {
		Lamp.MoveTo(&LivingRoom)
	}
	if Coffin.IsIn(Winner) {
		Coffin.MoveTo(&EgyptRoom)
	}
	Sword.TValue = 0
	// Copy children before iterating since MoveTo modifies the slice.
	children := make([]*Object, len(Winner.Children))
	copy(children, Winner.Children)
	for _, child := range children {
		if child.TValue <= 0 {
			child.MoveTo(Random(AboveGround))
			continue
		}
		for _, r := range Rooms.Children {
			if r.Has(FlgLand) && !r.Has(FlgOn) && Prob(50, false) {
				child.MoveTo(r)
				break
			}
		}
	}
}

func KillInterrupts() bool {
	QueueInt(IXb, false).Run = false
	QueueInt(IXc, false).Run = false
	QueueInt(ICyclops, false).Run = false
	QueueInt(ILantern, false).Run = false
	QueueInt(ICandles, false).Run = false
	QueueInt(ISword, false).Run = false
	QueueInt(IForestRandom, false).Run = false
	QueueInt(IMatch, false).Run = false
	Match.Take(FlgOn)
	return true
}

// ================================================================
// THE WHITE HOUSE
// ================================================================

func WestHouseFcn(arg ActArg) bool {
	if arg == ActLook {
		Print("You are standing in an open field west of a white house, with a boarded front door.", NoNewline)
		if WonGame {
			Print(" A secret path leads southwest into the forest.", NoNewline)
		}
		NewLine()
		return true
	}
	return false
}

func EastHouseFcn(arg ActArg) bool {
	if arg == ActLook {
		Print("You are behind the white house. A path leads into the forest to the east. In one corner of the house there is a small window which is ", NoNewline)
		if KitchenWindow.Has(FlgOpen) {
			Print("open.", Newline)
		} else {
			Print("slightly ajar.", Newline)
		}
		return true
	}
	return false
}

func WhiteHouseFcn(arg ActArg) bool {
	if Here == &Kitchen || Here == &LivingRoom || Here == &Attic {
		if ActVerb.Norm == "find" {
			Print("Why not find your brains?", Newline)
			return true
		}
		if ActVerb.Norm == "walk around" {
			GoNext(InHouseAround)
			return true
		}
	} else if Here != &WestOfHouse && Here != &NorthOfHouse && Here != &EastOfHouse && Here != &SouthOfHouse {
		if ActVerb.Norm == "find" {
			if Here == &Clearing {
				Print("It seems to be to the west.", Newline)
				return true
			}
			Print("It was here just a minute ago....", Newline)
			return true
		}
		Print("You're not at the house.", Newline)
		return true
	} else if ActVerb.Norm == "find" {
		Print("It's right here! Are you blind or something?", Newline)
		return true
	} else if ActVerb.Norm == "walk around" {
		GoNext(HouseAround)
		return true
	} else if ActVerb.Norm == "examine" {
		Print("The house is a beautiful colonial house which is painted white. It is clear that the owners must have been extremely wealthy.", Newline)
		return true
	} else if ActVerb.Norm == "through" || ActVerb.Norm == "open" {
		if Here == &EastOfHouse {
			if KitchenWindow.Has(FlgOpen) {
				return Goto(&Kitchen, true)
			}
			Print("The window is closed.", Newline)
			ThisIsIt(&KitchenWindow)
			return true
		}
		Print("I can't see how to get in from here.", Newline)
		return true
	} else if ActVerb.Norm == "burn" {
		Print("You must be joking.", Newline)
		return true
	}
	return false
}

func GoNext(tbl []*Object) int {
	val := Lkp(Here, tbl)
	if val == nil {
		return NumUndef
	}
	if !Goto(val, true) {
		return 2
	}
	return 1
}

func BoardFcn(arg ActArg) bool {
	if ActVerb.Norm == "take" || ActVerb.Norm == "examine" {
		Print("The boards are securely fastened.", Newline)
		return true
	}
	return false
}

func TeethFcn(arg ActArg) bool {
	if ActVerb.Norm == "brush" && DirObj == &Teeth {
		if IndirObj == &Putty && Putty.IsIn(Winner) {
			JigsUp("Well, you seem to have been brushing your teeth with some sort of glue. As a result, your mouth gets glued together (with your nose) and you die of respiratory failure.", false)
			return true
		}
		if IndirObj == nil {
			Print("Dental hygiene is highly recommended, but I'm not sure what you want to brush them with.", Newline)
			return true
		}
		Print("A nice idea, but with a ", NoNewline)
		PrintObject(IndirObj)
		Print("?", Newline)
		return true
	}
	return false
}

func GraniteWallFcn(arg ActArg) bool {
	if Here == &NorthTemple {
		if ActVerb.Norm == "find" {
			Print("The west wall is solid granite here.", Newline)
			return true
		}
		if ActVerb.Norm == "take" || ActVerb.Norm == "raise" || ActVerb.Norm == "lower" {
			Print("It's solid granite.", Newline)
			return true
		}
	} else if Here == &TreasureRoom {
		if ActVerb.Norm == "find" {
			Print("The east wall is solid granite here.", Newline)
			return true
		}
		if ActVerb.Norm == "take" || ActVerb.Norm == "raise" || ActVerb.Norm == "lower" {
			Print("It's solid granite.", Newline)
			return true
		}
	} else if Here == &SlideRoom {
		if ActVerb.Norm == "find" || ActVerb.Norm == "read" {
			Print("It only SAYS \"Granite Wall\".", Newline)
			return true
		}
		Print("The wall isn't granite.", Newline)
		return true
	} else {
		Print("There is no granite wall here.", Newline)
		return true
	}
	return false
}

func SongbirdFcn(arg ActArg) bool {
	if ActVerb.Norm == "find" || ActVerb.Norm == "take" {
		Print("The songbird is not here but is probably nearby.", Newline)
		return true
	}
	if ActVerb.Norm == "listen" {
		Print("You can't hear the songbird now.", Newline)
		return true
	}
	if ActVerb.Norm == "follow" {
		Print("It can't be followed.", Newline)
		return true
	}
	Print("You can't see any songbird here.", Newline)
	return true
}

func MountainRangeFcn(arg ActArg) bool {
	if ActVerb.Norm == "climb up" || ActVerb.Norm == "climb down" || ActVerb.Norm == "climb" {
		Print("Don't you believe me? The mountains are impassable!", Newline)
		return true
	}
	return false
}

func ForestFcn(arg ActArg) bool {
	if ActVerb.Norm == "walk around" {
		if Here == &WestOfHouse || Here == &NorthOfHouse || Here == &SouthOfHouse || Here == &EastOfHouse {
			Print("You aren't even in the forest.", Newline)
			return true
		}
		GoNext(ForestAround)
		return true
	}
	if ActVerb.Norm == "disembark" {
		Print("You will have to specify a direction.", Newline)
		return true
	}
	if ActVerb.Norm == "find" {
		Print("You cannot see the forest for the trees.", Newline)
		return true
	}
	if ActVerb.Norm == "listen" {
		Print("The pines and the hemlocks seem to be murmuring.", Newline)
		return true
	}
	return false
}

func KitchenWindowFcn(arg ActArg) bool {
	if ActVerb.Norm == "open" || ActVerb.Norm == "close" {
		KitchenWindowFlag = true
		OpenClose(&KitchenWindow,
			"With great effort, you open the window far enough to allow entry.",
			"The window closes (more easily than it opened).")
		return true
	}
	if ActVerb.Norm == "examine" && !KitchenWindowFlag {
		Print("The window is slightly ajar, but not enough to allow entry.", Newline)
		return true
	}
	if ActVerb.Norm == "walk" || ActVerb.Norm == "board" || ActVerb.Norm == "through" {
		if Here == &Kitchen {
			DoWalk("east")
		} else {
			DoWalk("west")
		}
		return true
	}
	if ActVerb.Norm == "look inside" {
		Print("You can see ", NoNewline)
		if Here == &Kitchen {
			Print("a clear area leading towards a forest.", Newline)
		} else {
			Print("what appears to be a kitchen.", Newline)
		}
		return true
	}
	return false
}

func ChimneyFcn(arg ActArg) bool {
	if ActVerb.Norm == "examine" {
		Print("The chimney leads ", NoNewline)
		if Here == &Kitchen {
			Print("down", NoNewline)
		} else {
			Print("up", NoNewline)
		}
		Print("ward, and looks climbable.", Newline)
		return true
	}
	return false
}

func GhostsFcn(arg ActArg) bool {
	if ActVerb.Norm == "tell" {
		Print("The spirits jeer loudly and ignore you.", Newline)
		Params.Continue = NumUndef
		return true
	}
	if ActVerb.Norm == "exorcise" {
		Print("Only the ceremony itself has any effect.", Newline)
		return true
	}
	if (ActVerb.Norm == "attack" || ActVerb.Norm == "mung") && DirObj == &Ghosts {
		Print("How can you attack a spirit with material objects?", Newline)
		return true
	}
	Print("You seem unable to interact with these spirits.", Newline)
	return true
}

func BasketFcn(arg ActArg) bool {
	if ActVerb.Norm == "raise" {
		if CageTop {
			Print(PickOne(Dummy), Newline)
		} else {
			RaisedBasket.MoveTo(&ShaftRoom)
			LoweredBasket.MoveTo(&LowerShaft)
			CageTop = true
			ThisIsIt(&RaisedBasket)
			Print("The basket is raised to the top of the shaft.", Newline)
		}
		return true
	}
	if ActVerb.Norm == "lower" {
		if !CageTop {
			Print(PickOne(Dummy), Newline)
		} else {
			RaisedBasket.MoveTo(&LowerShaft)
			LoweredBasket.MoveTo(&ShaftRoom)
			ThisIsIt(&LoweredBasket)
			Print("The basket is lowered to the bottom of the shaft.", Newline)
			CageTop = false
			if Lit && !IsLit(Here, true) {
				Lit = false
				Print("It is now pitch black.", Newline)
			}
		}
		return true
	}
	if DirObj == &LoweredBasket || IndirObj == &LoweredBasket {
		Print("The basket is at the other end of the chain.", Newline)
		return true
	}
	if ActVerb.Norm == "take" && (DirObj == &RaisedBasket || DirObj == &LoweredBasket) {
		Print("The cage is securely fastened to the iron chain.", Newline)
		return true
	}
	return false
}

func BatFcn(arg ActArg) bool {
	if ActVerb.Norm == "tell" {
		Fweep(6)
		Params.Continue = NumUndef
		return true
	}
	if ActVerb.Norm == "take" || ActVerb.Norm == "attack" || ActVerb.Norm == "mung" {
		if Garlic.Location() == Winner || Garlic.IsIn(Here) {
			Print("You can't reach him; he's on the ceiling.", Newline)
			return true
		}
		FlyMe()
		return true
	}
	return false
}

func BellFcn(arg ActArg) bool {
	if ActVerb.Norm == "ring" {
		if Here == &EnteranceToHades && !LLDFlag {
			return false
		}
		Print("Ding, dong.", Newline)
		return true
	}
	return false
}

func HotBellFcn(arg ActArg) bool {
	if ActVerb.Norm == "take" {
		Print("The bell is very hot and cannot be taken.", Newline)
		return true
	}
	if ActVerb.Norm == "rub" || (ActVerb.Norm == "ring" && IndirObj != nil) {
		if IndirObj != nil && IndirObj.Has(FlgBurn) {
			Print("The ", NoNewline)
			PrintObject(IndirObj)
			Print(" burns and is consumed.", Newline)
			RemoveCarefully(IndirObj)
			return true
		}
		if IndirObj == &Hands {
			Print("The bell is too hot to touch.", Newline)
			return true
		}
		Print("The heat from the bell is too intense.", Newline)
		return true
	}
	if ActVerb.Norm == "pour on" {
		RemoveCarefully(DirObj)
		Print("The water cools the bell and is evaporated.", Newline)
		QueueInt(IXbh, false).Run = false
		IXbh()
		return true
	}
	if ActVerb.Norm == "ring" {
		Print("The bell is too hot to reach.", Newline)
		return true
	}
	return false
}

func AxeFcn(arg ActArg) bool {
	if TrollFlag {
		return false
	}
	return WeaponFunction(&Axe, &Troll)
}

func TrapDoorFcn(arg ActArg) bool {
	if ActVerb.Norm == "raise" {
		Perform(ActionVerb{Norm: "open", Orig: "open"}, &TrapDoor, nil)
		return true
	}
	if (ActVerb.Norm == "open" || ActVerb.Norm == "close") && Here == &LivingRoom {
		OpenClose(DirObj,
			"The door reluctantly opens to reveal a rickety staircase descending into darkness.",
			"The door swings shut and closes.")
		return true
	}
	if ActVerb.Norm == "look under" && Here == &LivingRoom {
		if TrapDoor.Has(FlgOpen) {
			Print("You see a rickety staircase descending into darkness.", Newline)
		} else {
			Print("It's closed.", Newline)
		}
		return true
	}
	if Here == &Cellar {
		if (ActVerb.Norm == "open" || ActVerb.Norm == "unlock") && !TrapDoor.Has(FlgOpen) {
			Print("The door is locked from above.", Newline)
			return true
		}
		if ActVerb.Norm == "close" && !TrapDoor.Has(FlgOpen) {
			TrapDoor.Take(FlgTouch)
			TrapDoor.Take(FlgOpen)
			Print("The door closes and locks.", Newline)
			return true
		}
		if ActVerb.Norm == "open" || ActVerb.Norm == "close" {
			Print(PickOne(Dummy), Newline)
			return true
		}
	}
	return false
}

func FrontDoorFcn(arg ActArg) bool {
	if ActVerb.Norm == "open" {
		Print("The door cannot be opened.", Newline)
		return true
	}
	if ActVerb.Norm == "burn" {
		Print("You cannot burn this door.", Newline)
		return true
	}
	if ActVerb.Norm == "mung" {
		Print("You can't seem to damage the door.", Newline)
		return true
	}
	if ActVerb.Norm == "look behind" {
		Print("It won't open.", Newline)
		return true
	}
	return false
}

func BarrowDoorFcn(arg ActArg) bool {
	if ActVerb.Norm == "open" || ActVerb.Norm == "close" {
		Print("The door is too heavy.", Newline)
		return true
	}
	return false
}

func BarrowFcn(arg ActArg) bool {
	if ActVerb.Norm == "through" {
		DoWalk("west")
		return true
	}
	return false
}

func BottleFcn(arg ActArg) bool {
	empty := false
	if ActVerb.Norm == "throw" && DirObj == &Bottle {
		RemoveCarefully(DirObj)
		empty = true
		Print("The bottle hits the far wall and shatters.", Newline)
	} else if ActVerb.Norm == "mung" {
		empty = true
		RemoveCarefully(DirObj)
		Print("A brilliant maneuver destroys the bottle.", Newline)
	} else if ActVerb.Norm == "shake" {
		if Bottle.Has(FlgOpen) && Water.IsIn(&Bottle) {
			empty = true
		}
	}
	if empty && Water.IsIn(&Bottle) {
		Print("The water spills to the floor and evaporates.", Newline)
		RemoveCarefully(&Water)
		return true
	}
	if empty {
		return true
	}
	return false
}

func CrackFcn(arg ActArg) bool {
	if ActVerb.Norm == "through" {
		Print("You can't fit through the crack.", Newline)
		return true
	}
	return false
}

func GrateFcn(arg ActArg) bool {
	if ActVerb.Norm == "open" && IndirObj == &Keys {
		Perform(ActionVerb{Norm: "unlock", Orig: "unlock"}, &Grate, &Keys)
		return true
	}
	if ActVerb.Norm == "lock" {
		if Here == &GratingRoom {
			GrUnlock = false
			Print("The grate is locked.", Newline)
			return true
		}
		if Here == &Clearing {
			Print("You can't lock it from this side.", Newline)
			return true
		}
	}
	if ActVerb.Norm == "unlock" && DirObj == &Grate {
		if Here == &GratingRoom && IndirObj == &Keys {
			GrUnlock = true
			Print("The grate is unlocked.", Newline)
			return true
		}
		if Here == &Clearing && IndirObj == &Keys {
			Print("You can't reach the lock from here.", Newline)
			return true
		}
		Print("Can you unlock a grating with a ", NoNewline)
		PrintObject(IndirObj)
		Print("?", Newline)
		return true
	}
	if ActVerb.Norm == "pick" {
		Print("You can't pick the lock.", Newline)
		return true
	}
	if ActVerb.Norm == "open" || ActVerb.Norm == "close" {
		if GrUnlock {
			var openStr string
			if Here == &Clearing {
				openStr = "The grating opens."
			} else {
				openStr = "The grating opens to reveal trees above you."
			}
			OpenClose(&Grate, openStr, "The grating is closed.")
			if Grate.Has(FlgOpen) {
				if Here != &Clearing && !GrateRevealed {
					Print("A pile of leaves falls onto your head and to the ground.", Newline)
					GrateRevealed = true
					Leaves.MoveTo(Here)
				}
				GratingRoom.Give(FlgOn)
			} else {
				GratingRoom.Take(FlgOn)
			}
		} else {
			Print("The grating is locked.", Newline)
		}
		return true
	}
	if ActVerb.Norm == "put" && IndirObj == &Grate {
		if DirObj.Size > 20 {
			Print("It won't fit through the grating.", Newline)
		} else {
			DirObj.MoveTo(&GratingRoom)
			Print("The ", NoNewline)
			PrintObject(DirObj)
			Print(" goes through the grating into the darkness below.", Newline)
		}
		return true
	}
	return false
}

func KnifeFcn(arg ActArg) bool {
	if ActVerb.Norm == "take" {
		AtticTable.Take(FlgNoDesc)
		return false
	}
	return false
}

func SkeletonFcn(arg ActArg) bool {
	if ActVerb.Norm == "take" || ActVerb.Norm == "rub" || ActVerb.Norm == "move" ||
		ActVerb.Norm == "push" || ActVerb.Norm == "raise" || ActVerb.Norm == "lower" ||
		ActVerb.Norm == "attack" || ActVerb.Norm == "kick" || ActVerb.Norm == "kiss" {
		Print("A ghost appears in the room and is appalled at your desecration of the remains of a fellow adventurer. He casts a curse on your valuables and banishes them to the Land of the Living Dead. The ghost leaves, muttering obscenities.", Newline)
		Rob(Here, &LandOfLivingDead, 100)
		Rob(&Adventurer, &LandOfLivingDead, 0)
		return true
	}
	return false
}

func TorchFcn(arg ActArg) bool {
	if ActVerb.Norm == "examine" {
		Print("The torch is burning.", Newline)
		return true
	}
	if ActVerb.Norm == "pour on" && IndirObj == &Torch {
		Print("The water evaporates before it gets close.", Newline)
		return true
	}
	if ActVerb.Norm == "lamp off" && DirObj.Has(FlgOn) {
		Print("You nearly burn your hand trying to extinguish the flame.", Newline)
		return true
	}
	return false
}

func RustyKnifeFcn(arg ActArg) bool {
	if ActVerb.Norm == "take" {
		if Sword.IsIn(Winner) {
			Print("As you touch the rusty knife, your sword gives a single pulse of blinding blue light.", Newline)
		}
		return false
	}
	if (IndirObj == &RustyKnife && ActVerb.Norm == "attack") ||
		(ActVerb.Norm == "swing" && DirObj == &RustyKnife && IndirObj != nil) {
		RemoveCarefully(&RustyKnife)
		JigsUp("As the knife approaches its victim, your mind is submerged by an overmastering will. Slowly, your hand turns, until the rusty blade is an inch from your neck. The knife seems to sing as it savagely slits your throat.", false)
		return true
	}
	return false
}

func LeafPileFcn(arg ActArg) bool {
	if ActVerb.Norm == "count" {
		Print("There are 69,105 leaves here.", Newline)
		return true
	}
	if ActVerb.Norm == "burn" {
		LeavesAppear()
		RemoveCarefully(DirObj)
		if DirObj.IsIn(Here) {
			Print("The leaves burn.", Newline)
		} else {
			JigsUp("The leaves burn, and so do you.", false)
		}
		return true
	}
	if ActVerb.Norm == "cut" {
		Print("You rustle the leaves around, making quite a mess.", Newline)
		LeavesAppear()
		return true
	}
	if ActVerb.Norm == "move" || ActVerb.Norm == "take" {
		if ActVerb.Norm == "move" {
			Print("Done.", Newline)
		}
		if GrateRevealed {
			return false
		}
		LeavesAppear()
		if ActVerb.Norm == "take" {
			return false
		}
		return true
	}
	if ActVerb.Norm == "look under" && !GrateRevealed {
		Print("Underneath the pile of leaves is a grating. As you release the leaves, the grating is once again concealed from view.", Newline)
		return true
	}
	return false
}

func MatchFcn(arg ActArg) bool {
	if (ActVerb.Norm == "lamp on" || ActVerb.Norm == "burn") && DirObj == &Match {
		if MatchCount > 0 {
			MatchCount--
		}
		if MatchCount <= 0 {
			Print("I'm afraid that you have run out of matches.", Newline)
			return true
		}
		if Here == &LowerShaft || Here == &TimberRoom {
			Print("This room is drafty, and the match goes out instantly.", Newline)
			return true
		}
		Match.Give(FlgFlame)
		Match.Give(FlgOn)
		Queue(IMatch, 2).Run = true
		Print("One of the matches starts to burn.", Newline)
		if !Lit {
			Lit = true
			VLook(ActUnk)
		}
		return true
	}
	if ActVerb.Norm == "lamp off" && Match.Has(FlgFlame) {
		Print("The match is out.", Newline)
		Match.Take(FlgFlame)
		Match.Take(FlgOn)
		Lit = IsLit(Here, true)
		if !Lit {
			Print("It's pitch black in here!", Newline)
		}
		QueueInt(IMatch, false).Run = false
		return true
	}
	if ActVerb.Norm == "count" || ActVerb.Norm == "open" {
		Print("You have ", NoNewline)
		cnt := MatchCount - 1
		if cnt <= 0 {
			Print("no", NoNewline)
		} else {
			PrintNumber(cnt)
		}
		Print(" match", NoNewline)
		if cnt != 1 {
			Print("es.", NoNewline)
		} else {
			Print(".", NoNewline)
		}
		NewLine()
		return true
	}
	if ActVerb.Norm == "examine" {
		if Match.Has(FlgOn) {
			Print("The match is burning.", Newline)
		} else {
			Print("The matchbook isn't very interesting, except for what's written on it.", Newline)
		}
		return true
	}
	return false
}

func MirrorMirrorFcn(arg ActArg) bool {
	rm2 := &MirrorRoom2
	if !MirrorMung && ActVerb.Norm == "rub" {
		if IndirObj != nil && IndirObj != &Hands {
			Print("You feel a faint tingling transmitted through the ", NoNewline)
			PrintObject(IndirObj)
			Print(".", Newline)
			return true
		}
		if Here == rm2 {
			rm2 = &MirrorRoom1
		}
		// Swap room contents
		var l1, l2 []*Object
		for _, c := range Here.Children {
			l1 = append(l1, c)
		}
		for _, c := range rm2.Children {
			l2 = append(l2, c)
		}
		for _, c := range l1 {
			c.MoveTo(rm2)
		}
		for _, c := range l2 {
			c.MoveTo(Here)
		}
		Goto(rm2, false)
		Print("There is a rumble from deep within the earth and the room shakes.", Newline)
		return true
	}
	if ActVerb.Norm == "look inside" || ActVerb.Norm == "examine" {
		if MirrorMung {
			Print("The mirror is broken into many pieces.", Newline)
		} else {
			Print("There is an ugly person staring back at you.", Newline)
		}
		return true
	}
	if ActVerb.Norm == "take" {
		Print("The mirror is many times your size. Give up.", Newline)
		return true
	}
	if ActVerb.Norm == "mung" || ActVerb.Norm == "throw" || ActVerb.Norm == "attack" {
		if MirrorMung {
			Print("Haven't you done enough damage already?", Newline)
		} else {
			MirrorMung = true
			Lucky = false
			Print("You have broken the mirror. I hope you have a seven years' supply of good luck handy.", Newline)
		}
		return true
	}
	return false
}

func PaintingFcn(arg ActArg) bool {
	if ActVerb.Norm == "mung" {
		DirObj.TValue = 0
		DirObj.LongDesc = "There is a worthless piece of canvas here."
		Print("Congratulations! Unlike the other vandals, who merely stole the artist's masterpieces, you have destroyed one.", Newline)
		return true
	}
	return false
}

func CandlesFcn(arg ActArg) bool {
	if !Candles.Has(FlgTouch) {
		Queue(ICandles, -1).Run = true
	}
	if IndirObj == &Candles {
		return false
	}
	if ActVerb.Norm == "lamp on" || ActVerb.Norm == "burn" {
		if Candles.Has(FlgRMung) {
			Print("Alas, there's not much left of the candles. Certainly not enough to burn.", Newline)
			return true
		}
		if IndirObj == nil {
			if Match.Has(FlgFlame) {
				Print("(with the match)", Newline)
				Perform(ActionVerb{Norm: "lamp on", Orig: "light"}, &Candles, &Match)
				return true
			}
			Print("You should say what to light them with.", Newline)
			return true
		}
		if IndirObj == &Match && Match.Has(FlgOn) {
			Print("The candles are ", NoNewline)
			if Candles.Has(FlgOn) {
				Print("already lit.", Newline)
			} else {
				Candles.Give(FlgOn)
				Print("lit.", Newline)
				Queue(ICandles, -1).Run = true
			}
			return true
		}
		if IndirObj == &Torch {
			if Candles.Has(FlgOn) {
				Print("You realize, just in time, that the candles are already lighted.", Newline)
			} else {
				Print("The heat from the torch is so intense that the candles are vaporized.", Newline)
				RemoveCarefully(&Candles)
			}
			return true
		}
		Print("You have to light them with something that's burning, you know.", Newline)
		return true
	}
	if ActVerb.Norm == "count" {
		Print("Let's see, how many objects in a pair? Don't tell me, I'll get it.", Newline)
		return true
	}
	if ActVerb.Norm == "lamp off" {
		QueueInt(ICandles, false).Run = false
		if Candles.Has(FlgOn) {
			Print("The flame is extinguished.", NoNewline)
			Candles.Take(FlgOn)
			Candles.Give(FlgTouch)
			Lit = IsLit(Here, true)
			if !Lit {
				Print(" It's really dark in here....", NoNewline)
			}
			NewLine()
			return true
		}
		Print("The candles are not lighted.", Newline)
		return true
	}
	if ActVerb.Norm == "put" && IndirObj != nil && IndirObj.Has(FlgBurn) {
		Print("That wouldn't be smart.", Newline)
		return true
	}
	if ActVerb.Norm == "examine" {
		Print("The candles are ", NoNewline)
		if Candles.Has(FlgOn) {
			Print("burning.", Newline)
		} else {
			Print("out.", Newline)
		}
		return true
	}
	return false
}

func GunkFcn(arg ActArg) bool {
	RemoveCarefully(&Gunk)
	Print("The slag was rather insubstantial, and crumbles into dust at your touch.", Newline)
	return true
}

func BodyFcn(arg ActArg) bool {
	if ActVerb.Norm == "take" {
		Print("A force keeps you from taking the bodies.", Newline)
		return true
	}
	if ActVerb.Norm == "mung" || ActVerb.Norm == "burn" {
		JigsUp("The voice of the guardian of the dungeon booms out from the darkness, \"Your disrespect costs you your life!\" and places your head on a sharp pole.", false)
		return true
	}
	return false
}

func BlackBookFcn(arg ActArg) bool {
	if ActVerb.Norm == "open" {
		Print("The book is already open to page 569.", Newline)
		return true
	}
	if ActVerb.Norm == "close" {
		Print("As hard as you try, the book cannot be closed.", Newline)
		return true
	}
	if ActVerb.Norm == "turn" {
		Print("Beside page 569, there is only one other page with any legible printing on it. Most of it is unreadable, but the subject seems to be the banishment of evil. Apparently, certain noises, lights, and prayers are efficacious in this regard.", Newline)
		return true
	}
	if ActVerb.Norm == "burn" {
		RemoveCarefully(DirObj)
		JigsUp("A booming voice says \"Wrong, cretin!\" and you notice that you have turned into a pile of dust. How, I can't imagine.", false)
		return true
	}
	return false
}

func SceptreFcn(arg ActArg) bool {
	if ActVerb.Norm == "wave" || ActVerb.Norm == "raise" {
		if Here == &AragainFalls || Here == &EndOfRainbow {
			if !RainbowFlag {
				PotOfGold.Take(FlgInvis)
				Print("Suddenly, the rainbow appears to become solid and, I venture, walkable (I think the giveaway was the stairs and bannister).", Newline)
				if Here == &EndOfRainbow && PotOfGold.IsIn(&EndOfRainbow) {
					Print("A shimmering pot of gold appears at the end of the rainbow.", Newline)
				}
				RainbowFlag = true
			} else {
				Rob(&OnRainbow, &Wall, 0)
				Print("The rainbow seems to have become somewhat run-of-the-mill.", Newline)
				RainbowFlag = false
				return true
			}
			return true
		}
		if Here == &OnRainbow {
			RainbowFlag = false
			JigsUp("The structural integrity of the rainbow is severely compromised, leaving you hanging in midair, supported only by water vapor. Bye.", false)
			return true
		}
		Print("A dazzling display of color briefly emanates from the sceptre.", Newline)
		return true
	}
	return false
}

func SlideFcn(arg ActArg) bool {
	if ActVerb.Norm == "through" || ActVerb.Norm == "climb up" || ActVerb.Norm == "climb down" || ActVerb.Norm == "climb" ||
		(ActVerb.Norm == "put" && DirObj == &Me) {
		if Here == &Cellar {
			DoWalk("west")
			return true
		}
		Print("You tumble down the slide....", Newline)
		Goto(&Cellar, true)
		return true
	}
	if ActVerb.Norm == "put" {
		Slider(DirObj)
		return true
	}
	return false
}

func SandwichBagFcn(arg ActArg) bool {
	if ActVerb.Norm == "smell" && Lunch.IsIn(DirObj) {
		Print("It smells of hot peppers.", Newline)
		return true
	}
	return false
}

func ToolChestFcn(arg ActArg) bool {
	if ActVerb.Norm == "examine" {
		Print("The chests are all empty.", Newline)
		return true
	}
	if ActVerb.Norm == "take" || ActVerb.Norm == "open" || ActVerb.Norm == "put" {
		RemoveCarefully(&ToolChest)
		Print("The chests are so rusty and corroded that they crumble when you touch them.", Newline)
		return true
	}
	return false
}

func ButtonFcn(arg ActArg) bool {
	if ActVerb.Norm == "read" {
		Print("They're greek to you.", Newline)
		return true
	}
	if ActVerb.Norm == "push" {
		if DirObj == &BlueButton {
			if WaterLevel == 0 {
				Leak.Take(FlgInvis)
				Print("There is a rumbling sound and a stream of water appears to burst from the east wall of the room (apparently, a leak has occurred in a pipe).", Newline)
				WaterLevel = 1
				Queue(IMaintRoom, -1).Run = true
				return true
			}
			Print("The blue button appears to be jammed.", Newline)
			return true
		}
		if DirObj == &RedButton {
			Print("The lights within the room ", NoNewline)
			if Here.Has(FlgOn) {
				Here.Take(FlgOn)
				Print("shut off.", Newline)
			} else {
				Here.Give(FlgOn)
				Print("come on.", Newline)
			}
			return true
		}
		if DirObj == &BrownButton {
			DamRoom.Take(FlgTouch)
			GateFlag = false
			Print("Click.", Newline)
			return true
		}
		if DirObj == &YellowButton {
			DamRoom.Take(FlgTouch)
			GateFlag = true
			Print("Click.", Newline)
			return true
		}
		return true
	}
	return false
}

func LeakFcn(arg ActArg) bool {
	if WaterLevel > 0 {
		if (ActVerb.Norm == "put" || ActVerb.Norm == "put on") && DirObj == &Putty {
			FixMaintLeak()
			return true
		}
		if ActVerb.Norm == "plug" {
			if IndirObj == &Putty {
				FixMaintLeak()
				return true
			}
			WithTell(IndirObj)
			return true
		}
	}
	return false
}

func MachineFcn(arg ActArg) bool {
	if ActVerb.Norm == "take" && DirObj == &Machine {
		Print("It is far too large to carry.", Newline)
		return true
	}
	if ActVerb.Norm == "open" {
		if Machine.Has(FlgOpen) {
			Print(PickOne(Dummy), Newline)
		} else if Machine.HasChildren() {
			Print("The lid opens, revealing ", NoNewline)
			PrintContents(&Machine)
			Print(".", Newline)
			Machine.Give(FlgOpen)
		} else {
			Print("The lid opens.", Newline)
			Machine.Give(FlgOpen)
		}
		return true
	}
	if ActVerb.Norm == "close" {
		if Machine.Has(FlgOpen) {
			Print("The lid closes.", Newline)
			Machine.Take(FlgOpen)
		} else {
			Print(PickOne(Dummy), Newline)
		}
		return true
	}
	if ActVerb.Norm == "lamp on" {
		if IndirObj == nil {
			Print("It's not clear how to turn it on with your bare hands.", Newline)
		} else {
			Perform(ActionVerb{Norm: "turn", Orig: "turn"}, &MachineSwitch, IndirObj)
			return true
		}
		return true
	}
	return false
}

func MachineSwitchFcn(arg ActArg) bool {
	if ActVerb.Norm == "turn" {
		if IndirObj == &Screwdriver {
			if Machine.Has(FlgOpen) {
				Print("The machine doesn't seem to want to do anything.", Newline)
			} else {
				Print("The machine comes to life (figuratively) with a dazzling display of colored lights and bizarre noises. After a few moments, the excitement abates.", Newline)
				if Coal.IsIn(&Machine) {
					RemoveCarefully(&Coal)
					Diamond.MoveTo(&Machine)
				} else {
					// Remove everything and put gunk in
					var toRemove []*Object
					for _, o := range Machine.Children {
						toRemove = append(toRemove, o)
					}
					for _, o := range toRemove {
						RemoveCarefully(o)
					}
					Gunk.MoveTo(&Machine)
				}
			}
		} else {
			Print("It seems that a ", NoNewline)
			PrintObject(IndirObj)
			Print(" won't do.", Newline)
		}
		return true
	}
	return false
}

func PuttyFcn(arg ActArg) bool {
	if (ActVerb.Norm == "oil" && IndirObj == &Putty) || (ActVerb.Norm == "put" && DirObj == &Putty) {
		Print("The all-purpose gunk isn't a lubricant.", Newline)
		return true
	}
	return false
}

func TubeFcn(arg ActArg) bool {
	if ActVerb.Norm == "put" && IndirObj == &Tube {
		Print("The tube refuses to accept anything.", Newline)
		return true
	}
	if ActVerb.Norm == "squeeze" {
		if DirObj.Has(FlgOpen) && Putty.IsIn(DirObj) {
			Putty.MoveTo(Winner)
			Print("The viscous material oozes into your hand.", Newline)
			return true
		}
		if DirObj.Has(FlgOpen) {
			Print("The tube is apparently empty.", Newline)
			return true
		}
		Print("The tube is closed.", Newline)
		return true
	}
	return false
}

func SwordFcn(arg ActArg) bool {
	if ActVerb.Norm == "take" && Winner == &Adventurer {
		Queue(ISword, -1).Run = true
		return false
	}
	if ActVerb.Norm == "examine" {
		g := Sword.TValue
		if g == 1 {
			Print("Your sword is glowing with a faint blue glow.", Newline)
			return true
		}
		if g == 2 {
			Print("Your sword is glowing very brightly.", Newline)
			return true
		}
	}
	return false
}

func LanternFcn(arg ActArg) bool {
	if ActVerb.Norm == "throw" {
		Print("The lamp has smashed into the floor, and the light has gone out.", Newline)
		QueueInt(ILantern, false).Run = false
		RemoveCarefully(&Lamp)
		BrokenLamp.MoveTo(Here)
		return true
	}
	if ActVerb.Norm == "lamp on" {
		if Lamp.Has(FlgRMung) {
			Print("A burned-out lamp won't light.", Newline)
			return true
		}
		itm := QueueInt(ILantern, false)
		if itm.Tick <= 0 {
			// First activation or timer expired: initialize countdown
			itm.Tick = -1
		}
		// Otherwise resume from where we left off
		itm.Run = true
		return false
	}
	if ActVerb.Norm == "lamp off" {
		if Lamp.Has(FlgRMung) {
			Print("The lamp has already burned out.", Newline)
			return true
		}
		QueueInt(ILantern, false).Run = false
		return false
	}
	if ActVerb.Norm == "examine" {
		Print("The lamp ", NoNewline)
		if Lamp.Has(FlgRMung) {
			Print("has burned out.", Newline)
		} else if Lamp.Has(FlgOn) {
			Print("is on.", Newline)
		} else {
			Print("is turned off.", Newline)
		}
		return true
	}
	return false
}

func MailboxFcn(arg ActArg) bool {
	if ActVerb.Norm == "take" && DirObj == &Mailbox {
		Print("It is securely anchored.", Newline)
		return true
	}
	return false
}

// ================================================================
// TROLL
// ================================================================

func TrollFcn(arg ActArg) bool {
	if ActVerb.Norm == "tell" {
		Params.Continue = NumUndef
		Print("The troll isn't much of a conversationalist.", Newline)
		return true
	}
	if arg == ActArg(FBusy) {
		if Axe.IsIn(&Troll) {
			return false
		}
		if Axe.IsIn(Here) && Prob(75, true) {
			Axe.Give(FlgNoDesc)
			Axe.Take(FlgWeapon)
			Axe.MoveTo(&Troll)
			Troll.LongDesc = "A nasty-looking troll, brandishing a bloody axe, blocks all passages out of the room."
			if Troll.IsIn(Here) {
				Print("The troll, angered and humiliated, recovers his weapon. He appears to have an axe to grind with you.", Newline)
			}
			return true
		}
		if Troll.IsIn(Here) {
			Troll.LongDesc = "A pathetically babbling troll is here."
			Print("The troll, disarmed, cowers in terror, pleading for his life in the guttural tongue of the trolls.", Newline)
			return true
		}
		return false
	}
	if arg == ActArg(FDead) {
		if Axe.IsIn(&Troll) {
			Axe.MoveTo(Here)
			Axe.Take(FlgNoDesc)
			Axe.Give(FlgWeapon)
		}
		TrollFlag = true
		return true
	}
	if arg == ActArg(FUnconscious) {
		Troll.Take(FlgFight)
		if Axe.IsIn(&Troll) {
			Axe.MoveTo(Here)
			Axe.Take(FlgNoDesc)
			Axe.Give(FlgWeapon)
		}
		Troll.LongDesc = "An unconscious troll is sprawled on the floor. All passages out of the room are open."
		TrollFlag = true
		return true
	}
	if arg == ActArg(FConscious) {
		if Troll.IsIn(Here) {
			Troll.Give(FlgFight)
			Print("The troll stirs, quickly resuming a fighting stance.", Newline)
		}
		if Axe.IsIn(&Troll) {
			Troll.LongDesc = "A nasty-looking troll, brandishing a bloody axe, blocks all passages out of the room."
		} else if Axe.IsIn(&TrollRoom) {
			Axe.Give(FlgNoDesc)
			Axe.Take(FlgWeapon)
			Axe.MoveTo(&Troll)
			Troll.LongDesc = "A nasty-looking troll, brandishing a bloody axe, blocks all passages out of the room."
		} else {
			Troll.LongDesc = "A troll is here."
		}
		TrollFlag = false
		return true
	}
	if arg == ActArg(FFirst) {
		if Prob(33, false) {
			Troll.Give(FlgFight)
			Params.Continue = NumUndef
			return true
		}
		return false
	}
	// Default (no mode - regular verbs)
	if ActVerb.Norm == "examine" {
		Print(Troll.LongDesc, Newline)
		return true
	}
	if (ActVerb.Norm == "throw" || ActVerb.Norm == "give") && DirObj != nil && IndirObj == &Troll {
		Awaken(&Troll)
		if ActVerb.Norm == "throw" || ActVerb.Norm == "give" {
			if DirObj == &Axe && Axe.IsIn(Winner) {
				Print("The troll scratches his head in confusion, then takes the axe.", Newline)
				Troll.Give(FlgFight)
				Axe.MoveTo(&Troll)
				return true
			}
			if DirObj == &Troll || DirObj == &Axe {
				Print("You would have to get the ", NoNewline)
				PrintObject(DirObj)
				Print(" first, and that seems unlikely.", Newline)
				return true
			}
			if ActVerb.Norm == "throw" {
				Print("The troll, who is remarkably coordinated, catches the ", NoNewline)
				PrintObject(DirObj)
			} else {
				Print("The troll, who is not overly proud, graciously accepts the gift", NoNewline)
			}
			if Prob(20, false) && (DirObj == &Knife || DirObj == &Sword || DirObj == &Axe) {
				RemoveCarefully(DirObj)
				Print(" and eats it hungrily. Poor troll, he dies from an internal hemorrhage and his carcass disappears in a sinister black fog.", Newline)
				RemoveCarefully(&Troll)
				TrollFcn(ActArg(FDead))
				TrollFlag = true
			} else if DirObj == &Knife || DirObj == &Sword || DirObj == &Axe {
				DirObj.MoveTo(Here)
				Print(" and, being for the moment sated, throws it back. Fortunately, the troll has poor control, and the ", NoNewline)
				PrintObject(DirObj)
				Print(" falls to the floor. He does not look pleased.", Newline)
				Troll.Give(FlgFight)
			} else {
				Print(" and not having the most discriminating tastes, gleefully eats it.", Newline)
				RemoveCarefully(DirObj)
			}
			return true
		}
	}
	if ActVerb.Norm == "take" || ActVerb.Norm == "move" {
		Awaken(&Troll)
		if ActVerb.Norm == "take" || ActVerb.Norm == "move" {
			Print("The troll spits in your face, grunting \"Better luck next time\" in a rather barbarous accent.", Newline)
			return true
		}
	}
	if ActVerb.Norm == "mung" {
		Awaken(&Troll)
		Print("The troll laughs at your puny gesture.", Newline)
		return true
	}
	if ActVerb.Norm == "listen" {
		Print("Every so often the troll says something, probably uncomplimentary, in his guttural tongue.", Newline)
		return true
	}
	if TrollFlag && ActVerb.Norm == "hello" {
		Print("Unfortunately, the troll can't hear you.", Newline)
		return true
	}
	return false
}

// ================================================================
// CYCLOPS
// ================================================================

func CyclopsFcn(arg ActArg) bool {
	count := CycloWrath
	if Winner == &Cyclops {
		if CyclopsFlag {
			Print("No use talking to him. He's fast asleep.", Newline)
			return true
		}
		if ActVerb.Norm == "odysseus" {
			Winner = &Adventurer
			Perform(ActionVerb{Norm: "odysseus", Orig: "odysseus"}, nil, nil)
			return true
		}
		Print("The cyclops prefers eating to making conversation.", Newline)
		return true
	}
	if CyclopsFlag {
		if ActVerb.Norm == "examine" {
			Print("The cyclops is sleeping like a baby, albeit a very ugly one.", Newline)
			return true
		}
		if ActVerb.Norm == "alarm" || ActVerb.Norm == "kick" || ActVerb.Norm == "attack" || ActVerb.Norm == "burn" || ActVerb.Norm == "mung" {
			Print("The cyclops yawns and stares at the thing that woke him up.", Newline)
			CyclopsFlag = false
			Cyclops.Give(FlgFight)
			if count < 0 {
				CycloWrath = -count
			} else {
				CycloWrath = count
			}
			return true
		}
		return false
	}
	if ActVerb.Norm == "examine" {
		Print("A hungry cyclops is standing at the foot of the stairs.", Newline)
		return true
	}
	if ActVerb.Norm == "give" && IndirObj == &Cyclops {
		if DirObj == &Lunch {
			if count >= 0 {
				RemoveCarefully(&Lunch)
				Print("The cyclops says \"Mmm Mmm. I love hot peppers! But oh, could I use a drink. Perhaps I could drink the blood of that thing.\"  From the gleam in his eye, it could be surmised that you are \"that thing\".", Newline)
				CycloWrath = MinInt(-1, -count)
			}
			Queue(ICyclops, -1).Run = true
			return true
		}
		if DirObj == &Water || (DirObj == &Bottle && Water.IsIn(&Bottle)) {
			if count < 0 {
				RemoveCarefully(&Water)
				Bottle.MoveTo(Here)
				Bottle.Give(FlgOpen)
				Cyclops.Take(FlgFight)
				Print("The cyclops takes the bottle, checks that it's open, and drinks the water. A moment later, he lets out a yawn that nearly blows you over, and then falls fast asleep (what did you put in that drink, anyway?).", Newline)
				CyclopsFlag = true
			} else {
				Print("The cyclops apparently is not thirsty and refuses your generous offer.", Newline)
			}
			return true
		}
		if DirObj == &Garlic {
			Print("The cyclops may be hungry, but there is a limit.", Newline)
			return true
		}
		Print("The cyclops is not so stupid as to eat THAT!", Newline)
		return true
	}
	if ActVerb.Norm == "throw" || ActVerb.Norm == "attack" || ActVerb.Norm == "mung" {
		Queue(ICyclops, -1).Run = true
		if ActVerb.Norm == "mung" {
			Print("\"Do you think I'm as stupid as my father was?\", he says, dodging.", Newline)
		} else {
			Print("The cyclops shrugs but otherwise ignores your pitiful attempt.", Newline)
			if ActVerb.Norm == "throw" {
				DirObj.MoveTo(Here)
			}
			return true
		}
		return true
	}
	if ActVerb.Norm == "take" {
		Print("The cyclops doesn't take kindly to being grabbed.", Newline)
		return true
	}
	if ActVerb.Norm == "tie" {
		Print("You cannot tie the cyclops, though he is fit to be tied.", Newline)
		return true
	}
	if ActVerb.Norm == "listen" {
		Print("You can hear his stomach rumbling.", Newline)
		return true
	}
	return false
}

// ================================================================
// THIEF / ROBBER
// ================================================================

func DumbContainerFcn(arg ActArg) bool {
	if ActVerb.Norm == "open" || ActVerb.Norm == "close" || ActVerb.Norm == "look inside" {
		Print("You can't do that.", Newline)
		return true
	}
	if ActVerb.Norm == "examine" {
		Print("It looks pretty much like a ", NoNewline)
		PrintObject(DirObj)
		Print(".", Newline)
		return true
	}
	return false
}

func ChaliceFcn(arg ActArg) bool {
	if ActVerb.Norm == "take" {
		if DirObj.IsIn(&TreasureRoom) && Thief.IsIn(&TreasureRoom) && Thief.Has(FlgFight) && !Thief.Has(FlgInvis) && Thief.LongDesc != RobberUDesc {
			Print("You'd be stabbed in the back first.", Newline)
			return true
		}
		return false
	}
	if ActVerb.Norm == "put" && IndirObj == &Chalice {
		Print("You can't. It's not a very good chalice, is it?", Newline)
		return true
	}
	return DumbContainerFcn(arg)
}

func TrunkFcn(arg ActArg) bool {
	return StupidContainer(&Trunk, "jewels")
}

func BagOfCoinsFcn(arg ActArg) bool {
	return StupidContainer(&BagOfCoins, "coins")
}

func StupidContainer(obj *Object, str string) bool {
	if ActVerb.Norm == "open" || ActVerb.Norm == "close" {
		Print("The "+str+" are safely inside; there's no need to do that.", Newline)
		return true
	}
	if ActVerb.Norm == "look inside" || ActVerb.Norm == "examine" {
		Print("There are lots of "+str+" in there.", Newline)
		return true
	}
	if ActVerb.Norm == "put" && IndirObj == obj {
		Print("Don't be silly. It wouldn't be a ", NoNewline)
		PrintObject(obj)
		Print(" anymore.", Newline)
		return true
	}
	return false
}

func GarlicFcn(arg ActArg) bool {
	if ActVerb.Norm == "eat" {
		RemoveCarefully(DirObj)
		Print("What the heck! You won't make friends this way, but nobody around here is too friendly anyhow. Gulp!", Newline)
		return true
	}
	return false
}

func BatDescFcn(arg ActArg) bool {
	if Garlic.Location() == Winner || Garlic.IsIn(Here) {
		Print("In the corner of the room on the ceiling is a large vampire bat who is obviously deranged and holding his nose.", Newline)
	} else {
		Print("A large vampire bat, hanging from the ceiling, swoops down at you!", Newline)
	}
	return true
}

func TrophyCaseFcn(arg ActArg) bool {
	if ActVerb.Norm == "take" && DirObj == &TrophyCase {
		Print("The trophy case is securely fastened to the wall.", Newline)
		return true
	}
	return false
}

func BoardedWindowFcn(arg ActArg) bool {
	if ActVerb.Norm == "open" {
		Print("The windows are boarded and can't be opened.", Newline)
		return true
	}
	if ActVerb.Norm == "mung" {
		Print("You can't break the windows open.", Newline)
		return true
	}
	return false
}

func NailsPseudo(arg ActArg) bool {
	if ActVerb.Norm == "take" {
		Print("The nails, deeply imbedded in the door, cannot be removed.", Newline)
		return true
	}
	return false
}

func CliffObjectFcn(arg ActArg) bool {
	if ActVerb.Norm == "leap" || (ActVerb.Norm == "put" && DirObj == &Me) {
		Print("That would be very unwise. Perhaps even fatal.", Newline)
		return true
	}
	if IndirObj == &ClimbableCliff {
		if ActVerb.Norm == "put" || ActVerb.Norm == "throw off" {
			Print("The ", NoNewline)
			PrintObject(DirObj)
			Print(" tumbles into the river and is seen no more.", Newline)
			RemoveCarefully(DirObj)
			return true
		}
	}
	return false
}

func WhiteCliffFcn(arg ActArg) bool {
	if ActVerb.Norm == "climb up" || ActVerb.Norm == "climb down" || ActVerb.Norm == "climb" {
		Print("The cliff is too steep for climbing.", Newline)
		return true
	}
	return false
}

func RainbowFcn(arg ActArg) bool {
	if ActVerb.Norm == "cross" || ActVerb.Norm == "through" {
		if Here == &CanyonView {
			Print("From here?!?", Newline)
			return true
		}
		if RainbowFlag {
			if Here == &AragainFalls {
				Goto(&EndOfRainbow, true)
			} else if Here == &EndOfRainbow {
				Goto(&AragainFalls, true)
			} else {
				Print("You'll have to say which way...", Newline)
			}
		} else {
			Print("Can you walk on water vapor?", Newline)
		}
		return true
	}
	if ActVerb.Norm == "look under" {
		Print("The Frigid River flows under the rainbow.", Newline)
		return true
	}
	return false
}

func RopeFcn(arg ActArg) bool {
	if Here != &DomeRoom {
		DomeFlag = false
		if ActVerb.Norm == "tie" {
			Print("You can't tie the rope to that.", Newline)
			return true
		}
		return false
	}
	if ActVerb.Norm == "tie" {
		if IndirObj == &Railing {
			if DomeFlag {
				Print("The rope is already tied to it.", Newline)
			} else {
				Print("The rope drops over the side and comes within ten feet of the floor.", Newline)
				DomeFlag = true
				Rope.Give(FlgNoDesc)
				rloc := Rope.Location()
				if rloc == nil || !rloc.IsIn(&Rooms) {
					Rope.MoveTo(Here)
				}
			}
			return true
		}
		return false
	}
	if ActVerb.Norm == "climb down" && (DirObj == &Rope || DirObj == &Rooms) && DomeFlag {
		DoWalk("down")
		return true
	}
	if ActVerb.Norm == "tie up" && IndirObj == &Rope {
		if DirObj.Has(FlgActor) {
			if DirObj.Strength < 0 {
				Print("Your attempt to tie up the ", NoNewline)
				PrintObject(DirObj)
				Print(" awakens him.", NoNewline)
				Awaken(DirObj)
			} else {
				Print("The ", NoNewline)
				PrintObject(DirObj)
				Print(" struggles and you cannot tie him up.", Newline)
			}
		} else {
			Print("Why would you tie up a ", NoNewline)
			PrintObject(DirObj)
			Print("?", Newline)
		}
		return true
	}
	if ActVerb.Norm == "untie" {
		if DomeFlag {
			DomeFlag = false
			Rope.Take(FlgNoDesc)
			Print("The rope is now untied.", Newline)
		} else {
			Print("It is not tied to anything.", Newline)
		}
		return true
	}
	if ActVerb.Norm == "drop" && Here == &DomeRoom && !DomeFlag {
		Rope.MoveTo(&TorchRoom)
		Print("The rope drops gently to the floor below.", Newline)
		return true
	}
	if ActVerb.Norm == "take" {
		if DomeFlag {
			Print("The rope is tied to the railing.", Newline)
			return true
		}
	}
	return false
}

func EggObjectFcn(arg ActArg) bool {
	if (ActVerb.Norm == "open" || ActVerb.Norm == "mung") && DirObj == &Egg {
		if DirObj.Has(FlgOpen) {
			Print("The egg is already open.", Newline)
			return true
		}
		if IndirObj == nil {
			Print("You have neither the tools nor the expertise.", Newline)
			return true
		}
		if IndirObj == &Hands {
			Print("I doubt you could do that without damaging it.", Newline)
			return true
		}
		if IndirObj.Has(FlgWeapon) || IndirObj.Has(FlgTool) || ActVerb.Norm == "mung" {
			Print("The egg is now open, but the clumsiness of your attempt has seriously compromised its esthetic appeal.", NoNewline)
			BadEgg()
			NewLine()
			return true
		}
		if DirObj.Has(FlgFight) {
			Print("Not to say that using the ", NoNewline)
			PrintObject(IndirObj)
			Print(" isn't original too...", Newline)
			return true
		}
		Print("The concept of using a ", NoNewline)
		PrintObject(IndirObj)
		Print(" is certainly original.", Newline)
		DirObj.Give(FlgFight)
		return true
	}
	if ActVerb.Norm == "climb on" || ActVerb.Norm == "hatch" {
		Print("There is a noticeable crunch from beneath you, and inspection reveals that the egg is lying open, badly damaged.", NoNewline)
		BadEgg()
		NewLine()
		return true
	}
	if ActVerb.Norm == "open" || ActVerb.Norm == "mung" || ActVerb.Norm == "throw" {
		if ActVerb.Norm == "throw" {
			DirObj.MoveTo(Here)
		}
		Print("Your rather indelicate handling of the egg has caused it some damage, although you have succeeded in opening it.", NoNewline)
		BadEgg()
		NewLine()
		return true
	}
	return false
}

func CanaryObjectFcn(arg ActArg) bool {
	if ActVerb.Norm == "wind" {
		if DirObj == &Canary {
			if !SingSong && ForestRoomQ() {
				Print("The canary chirps, slightly off-key, an aria from a forgotten opera. From out of the greenery flies a lovely songbird. It perches on a limb just over your head and opens its beak to sing. As it does so a beautiful brass bauble drops from its mouth, bounces off the top of your head, and lands glimmering in the grass. As the canary winds down, the songbird flies away.", Newline)
				SingSong = true
				dest := Here
				if Here == &UpATree {
					dest = &Path
				}
				Bauble.MoveTo(dest)
			} else {
				Print("The canary chirps blithely, if somewhat tinnily, for a short time.", Newline)
			}
		} else {
			Print("There is an unpleasant grinding noise from inside the canary.", Newline)
		}
		return true
	}
	return false
}

func RugFcn(arg ActArg) bool {
	if ActVerb.Norm == "raise" {
		Print("The rug is too heavy to lift", NoNewline)
		if RugMoved {
			Print(".", Newline)
		} else {
			Print(", but in trying to take it you have noticed an irregularity beneath it.", Newline)
		}
		return true
	}
	if ActVerb.Norm == "move" || ActVerb.Norm == "push" {
		if RugMoved {
			Print("Having moved the carpet previously, you find it impossible to move it again.", Newline)
		} else {
			Print("With a great effort, the rug is moved to one side of the room, revealing the dusty cover of a closed trap door.", Newline)
			TrapDoor.Take(FlgInvis)
			ThisIsIt(&TrapDoor)
			RugMoved = true
		}
		return true
	}
	if ActVerb.Norm == "take" {
		Print("The rug is extremely heavy and cannot be carried.", Newline)
		return true
	}
	if ActVerb.Norm == "look under" && !RugMoved && !TrapDoor.Has(FlgOpen) {
		Print("Underneath the rug is a closed trap door. As you drop the corner of the rug, the trap door is once again concealed from view.", Newline)
		return true
	}
	if ActVerb.Norm == "climb on" {
		if !RugMoved && !TrapDoor.Has(FlgOpen) {
			Print("As you sit, you notice an irregularity underneath it. Rather than be uncomfortable, you stand up again.", Newline)
		} else {
			Print("I suppose you think it's a magic carpet?", Newline)
		}
		return true
	}
	return false
}

func SandFunction(arg ActArg) bool {
	if ActVerb.Norm == "dig" && IndirObj == &Shovel {
		BeachDig++
		if BeachDig > 3 {
			BeachDig = -1
			if Scarab.IsIn(Here) {
				Scarab.Give(FlgInvis)
			}
			JigsUp("The hole collapses, smothering you.", false)
			return true
		}
		if BeachDig == 3 {
			if Scarab.Has(FlgInvis) {
				Print("You can see a scarab here in the sand.", Newline)
				ThisIsIt(&Scarab)
				Scarab.Take(FlgInvis)
			}
		} else {
			Print(BDigs[BeachDig], Newline)
		}
		return true
	}
	return false
}

// ================================================================
// ROOM ACTION FUNCTIONS
// ================================================================

func KitchenFcn(arg ActArg) bool {
	if arg == ActLook {
		Print("You are in the kitchen of the white house. A table seems to have been used recently for the preparation of food. A passage leads to the west and a dark staircase can be seen leading upward. A dark chimney leads down and to the east is a small window which is ", NoNewline)
		if KitchenWindow.Has(FlgOpen) {
			Print("open.", Newline)
		} else {
			Print("slightly ajar.", Newline)
		}
		return true
	}
	if arg == ActBegin {
		if ActVerb.Norm == "climb up" && DirObj == &Stairs {
			DoWalk("up")
			return true
		}
	}
	return false
}

func LivingRoomFcn(arg ActArg) bool {
	if arg == ActLook {
		Print("You are in the living room. There is a doorway to the east", NoNewline)
		if MagicFlag {
			Print(". To the west is a cyclops-shaped opening in an old wooden door, above which is some strange gothic lettering, ", NoNewline)
		} else {
			Print(", a wooden door with strange gothic lettering to the west, which appears to be nailed shut, ", NoNewline)
		}
		Print("a trophy case, ", NoNewline)
		if RugMoved && TrapDoor.Has(FlgOpen) {
			Print("and a rug lying beside an open trap door.", NoNewline)
		} else if RugMoved {
			Print("and a closed trap door at your feet.", NoNewline)
		} else if TrapDoor.Has(FlgOpen) {
			Print("and an open trap door at your feet.", NoNewline)
		} else {
			Print("and a large oriental rug in the center of the room.", NoNewline)
		}
		NewLine()
		return true
	}
	if arg == ActEnd {
		if ActVerb.Norm == "take" || (ActVerb.Norm == "put" && IndirObj == &TrophyCase) {
			if DirObj.IsIn(&TrophyCase) {
				TouchAll(DirObj)
			}
			Score = BaseScore + OtvalFrob(&TrophyCase)
			ScoreUpd(0)
			return false
		}
	}
	return false
}

func CellarFcn(arg ActArg) bool {
	if arg == ActLook {
		Print("You are in a dark and damp cellar with a narrow passageway leading north, and a crawlway to the south. On the west is the bottom of a steep metal ramp which is unclimbable.", Newline)
		return true
	}
	if arg == ActEnter {
		if TrapDoor.Has(FlgOpen) && !TrapDoor.Has(FlgTouch) {
			TrapDoor.Take(FlgOpen)
			TrapDoor.Give(FlgTouch)
			Print("The trap door crashes shut, and you hear someone barring it.", Newline)
			NewLine()
		}
		return false
	}
	return false
}

func StoneBarrowFcn(arg ActArg) bool {
	if arg == ActBegin {
		if ActVerb.Norm == "enter" || (ActVerb.Norm == "walk" && (DirObj == ToDirObj("west") || DirObj == ToDirObj("in"))) || (ActVerb.Norm == "through" && DirObj == &Barrow) {
			Print("Inside the Barrow\nAs you enter the barrow, the door closes inexorably behind you. Around you it is dark, but ahead is an enormous cavern, brightly lit. Through its center runs a wide stream. Spanning the stream is a small wooden footbridge, and beyond a path leads into a dark tunnel. Above the bridge, floating in the air, is a large sign. It reads:  All ye who stand before this bridge have completed a great and perilous adventure which has tested your wit and courage. You have mastered the first part of the ZORK trilogy. Those who pass over this bridge must be prepared to undertake an even greater adventure that will severely test your skill and bravery!\n\nThe ZORK trilogy continues with \"ZORK II: The Wizard of Frobozz\" and is completed in \"ZORK III: The Dungeon Master.\"", Newline)
			Finish()
			return true
		}
	}
	return false
}

func TrollRoomFcn(arg ActArg) bool {
	if arg == ActEnter && Troll.IsIn(Here) {
		ThisIsIt(&Troll)
	}
	return false
}

func ClearingFcn(arg ActArg) bool {
	if arg == ActEnter {
		if !GrateRevealed {
			Grate.Give(FlgInvis)
		}
		return false
	}
	if arg == ActLook {
		Print("You are in a clearing, with a forest surrounding you on all sides. A path leads south.", NoNewline)
		if Grate.Has(FlgOpen) {
			NewLine()
			Print("There is an open grating, descending into darkness.", NoNewline)
		} else if GrateRevealed {
			NewLine()
			Print("There is a grating securely fastened into the ground.", NoNewline)
		}
		NewLine()
		return true
	}
	return false
}

func Maze11Fcn(arg ActArg) bool {
	if arg == ActEnter {
		Grate.Take(FlgInvis)
		return false
	}
	if arg == ActLook {
		Print("You are in a small room near the maze. There are twisty passages in the immediate vicinity.", Newline)
		if Grate.Has(FlgOpen) {
			Print("Above you is an open grating with sunlight pouring in.", NoNewline)
		} else if GrUnlock {
			Print("Above you is a grating.", NoNewline)
		} else {
			Print("Above you is a grating locked with a skull-and-crossbones lock.", NoNewline)
		}
		NewLine()
		return true
	}
	return false
}

func CyclopsRoomFcn(arg ActArg) bool {
	if arg == ActLook {
		Print("This room has an exit on the northwest, and a staircase leading up.", Newline)
		if CyclopsFlag && !MagicFlag {
			Print("The cyclops is sleeping blissfully at the foot of the stairs.", Newline)
		} else if MagicFlag {
			Print("The east wall, previously solid, now has a cyclops-sized opening in it.", Newline)
		} else if CycloWrath == 0 {
			Print("A cyclops, who looks prepared to eat horses (much less mere adventurers), blocks the staircase. From his state of health, and the bloodstains on the walls, you gather that he is not very friendly, though he likes people.", Newline)
		} else if CycloWrath > 0 {
			Print("The cyclops is standing in the corner, eyeing you closely. I don't think he likes you very much. He looks extremely hungry, even for a cyclops.", Newline)
		} else {
			Print("The cyclops, having eaten the hot peppers, appears to be gasping. His enflamed tongue protrudes from his man-sized mouth.", Newline)
		}
		return true
	}
	if arg == ActEnter {
		if CycloWrath == 0 {
			return false
		}
		Queue(ICyclops, -1).Run = true
		return false
	}
	return false
}

func TreasureRoomFcn(arg ActArg) bool {
	if arg == ActEnter && !Dead {
		if !Thief.IsIn(Here) {
			Print("You hear a scream of anguish as you violate the robber's hideaway. Using passages unknown to you, he rushes to its defense.", Newline)
			Thief.MoveTo(Here)
		}
		Thief.Give(FlgFight)
		Thief.Take(FlgInvis)
		ThiefInTreasure()
		return true
	}
	return false
}

func ReservoirSouthFcn(arg ActArg) bool {
	if arg == ActLook {
		if LowTide && GatesOpen {
			Print("You are in a long room, to the north of which was formerly a lake. However, with the water level lowered, there is merely a wide stream running through the center of the room.", NoNewline)
		} else if GatesOpen {
			Print("You are in a long room. To the north is a large lake, too deep to cross. You notice, however, that the water level appears to be dropping at a rapid rate. Before long, it might be possible to cross to the other side from here.", NoNewline)
		} else if LowTide {
			Print("You are in a long room, to the north of which is a wide area which was formerly a reservoir, but now is merely a stream. You notice, however, that the level of the stream is rising quickly and that before long it will be impossible to cross here.", NoNewline)
		} else {
			Print("You are in a long room on the south shore of a large lake, far too deep and wide for crossing.", NoNewline)
		}
		NewLine()
		Print("There is a path along the stream to the east or west, a steep pathway climbing southwest along the edge of a chasm, and a path leading into a canyon to the southeast.", Newline)
		return true
	}
	return false
}

func ReservoirFcn(arg ActArg) bool {
	if arg == ActEnd && !Winner.Location().Has(FlgVeh) && !GatesOpen && LowTide {
		Print("You notice that the water level here is rising rapidly. The currents are also becoming stronger. Staying here seems quite perilous!", Newline)
		return true
	}
	if arg == ActLook {
		if LowTide {
			Print("You are on what used to be a large lake, but which is now a large mud pile. There are \"shores\" to the north and south.", NoNewline)
		} else {
			Print("You are on the lake. Beaches can be seen north and south. Upstream a small stream enters the lake through a narrow cleft in the rocks. The dam can be seen downstream.", NoNewline)
		}
		NewLine()
		return true
	}
	return false
}

func ReservoirNorthFcn(arg ActArg) bool {
	if arg == ActLook {
		if LowTide && GatesOpen {
			Print("You are in a large cavernous room, the south of which was formerly a lake. However, with the water level lowered, there is merely a wide stream running through there.", NoNewline)
		} else if GatesOpen {
			Print("You are in a large cavernous area. To the south is a wide lake, whose water level appears to be falling rapidly.", NoNewline)
		} else if LowTide {
			Print("You are in a cavernous area, to the south of which is a very wide stream. The level of the stream is rising rapidly, and it appears that before long it will be impossible to cross to the other side.", NoNewline)
		} else {
			Print("You are in a large cavernous room, north of a large lake.", NoNewline)
		}
		NewLine()
		Print("There is a slimy stairway leaving the room to the north.", Newline)
		return true
	}
	return false
}

func MirrorRoomFcn(arg ActArg) bool {
	if arg == ActLook {
		Print("You are in a large square room with tall ceilings. On the south wall is an enormous mirror which fills the entire wall. There are exits on the other three sides of the room.", Newline)
		if MirrorMung {
			Print("Unfortunately, the mirror has been destroyed by your recklessness.", Newline)
		}
		return true
	}
	return false
}

func Cave2RoomFcn(arg ActArg) bool {
	if arg == ActEnd {
		if Candles.IsIn(Winner) && Prob(50, true) && Candles.Has(FlgOn) {
			QueueInt(ICandles, false).Run = false
			Candles.Take(FlgOn)
			Print("A gust of wind blows out your candles!", Newline)
			Lit = IsLit(Here, true)
			if !Lit {
				Print("It is now completely dark.", Newline)
			}
			return true
		}
	}
	return false
}

func LLDRoomFcn(arg ActArg) bool {
	if arg == ActLook {
		Print("You are outside a large gateway, on which is inscribed\n\n  Abandon every hope\nall ye who enter here!\n\nThe gate is open; through it you can see a desolation, with a pile of mangled bodies in one corner. Thousands of voices, lamenting some hideous fate, can be heard.", Newline)
		if !LLDFlag && !Dead {
			Print("The way through the gate is barred by evil spirits, who jeer at your attempts to pass.", Newline)
		}
		return true
	}
	if arg == ActBegin {
		if ActVerb.Norm == "exorcise" && !LLDFlag {
			if Bell.IsIn(Winner) && Book.IsIn(Winner) && Candles.IsIn(Winner) {
				Print("You must perform the ceremony.", Newline)
			} else {
				Print("You aren't equipped for an exorcism.", Newline)
			}
			return true
		}
		if !LLDFlag && ActVerb.Norm == "ring" && DirObj == &Bell {
			XB = true
			RemoveCarefully(&Bell)
			ThisIsIt(&HotBell)
			HotBell.MoveTo(Here)
			Print("The bell suddenly becomes red hot and falls to the ground. The wraiths, as if paralyzed, stop their jeering and slowly turn to face you. On their ashen faces, the expression of a long-forgotten terror takes shape.", Newline)
			if Candles.IsIn(Winner) {
				Print("In your confusion, the candles drop to the ground (and they are out).", Newline)
				Candles.MoveTo(Here)
				Candles.Take(FlgOn)
				QueueInt(ICandles, false).Run = false
			}
			Queue(IXb, 6).Run = true
			Queue(IXbh, 20).Run = true
			return true
		}
		if XC && ActVerb.Norm == "read" && DirObj == &Book && !LLDFlag {
			Print("Each word of the prayer reverberates through the hall in a deafening confusion. As the last word fades, a voice, loud and commanding, speaks: \"Begone, fiends!\" A heart-stopping scream fills the cavern, and the spirits, sensing a greater power, flee through the walls.", Newline)
			RemoveCarefully(&Ghosts)
			LLDFlag = true
			QueueInt(IXc, false).Run = false
			return true
		}
	}
	if arg == ActEnd {
		if XB && Candles.IsIn(Winner) && Candles.Has(FlgOn) && !XC {
			XC = true
			Print("The flames flicker wildly and appear to dance. The earth beneath your feet trembles, and your legs nearly buckle beneath you. The spirits cower at your unearthly power.", Newline)
			QueueInt(IXb, false).Run = false
			Queue(IXc, 3).Run = true
		}
	}
	return false
}

func DomeRoomFcn(arg ActArg) bool {
	if arg == ActLook {
		Print("You are at the periphery of a large dome, which forms the ceiling of another room below. Protecting you from a precipitous drop is a wooden railing which circles the dome.", Newline)
		if DomeFlag {
			Print("Hanging down from the railing is a rope which ends about ten feet from the floor below.", Newline)
		}
		return true
	}
	if arg == ActEnter {
		if Dead {
			Print("As you enter the dome you feel a strong pull as if from a wind drawing you over the railing and down.", Newline)
			Winner.MoveTo(&TorchRoom)
			Here = &TorchRoom
			return true
		}
		if ActVerb.Norm == "leap" {
			JigsUp("I'm afraid that the leap you attempted has done you in.", false)
			return true
		}
	}
	return false
}

func TorchRoomFcn(arg ActArg) bool {
	if arg == ActLook {
		Print("This is a large room with a prominent doorway leading to a down staircase. Above you is a large dome. Up around the edge of the dome (20 feet up) is a wooden railing. In the center of the room sits a white marble pedestal.", Newline)
		if DomeFlag {
			Print("A piece of rope descends from the railing above, ending some five feet above your head.", Newline)
		}
		return true
	}
	return false
}

func SouthTempleFcn(arg ActArg) bool {
	if arg == ActBegin {
		CoffinCure = !Coffin.IsIn(Winner)
		return false
	}
	return false
}

func MachineRoomFcn(arg ActArg) bool {
	if arg == ActLook {
		Print("This is a large, cold room whose sole exit is to the north. In one corner there is a machine which is reminiscent of a clothes dryer. On its face is a switch which is labelled \"START\". The switch does not appear to be manipulable by any human hand (unless the fingers are about 1/16 by 1/4 inch). On the front of the machine is a large lid, which is ", NoNewline)
		if Machine.Has(FlgOpen) {
			Print("open.", Newline)
		} else {
			Print("closed.", Newline)
		}
		return true
	}
	return false
}

func LoudRoomFcn(arg ActArg) bool {
	if arg == ActLook {
		Print("This is a large room with a ceiling which cannot be detected from the ground. There is a narrow passage from east to west and a stone stairway leading upward.", NoNewline)
		if LoudFlag || (!GatesOpen && LowTide) {
			Print(" The room is eerie in its quietness.", NoNewline)
		} else {
			Print(" The room is deafeningly loud with an undetermined rushing sound. The sound seems to reverberate from all of the walls, making it difficult even to think.", NoNewline)
		}
		NewLine()
		return true
	}
	if arg == ActEnd && GatesOpen && !LowTide {
		Print("It is unbearably loud here, with an ear-splitting roar seeming to come from all around you. There is a pounding in your head which won't stop. With a tremendous effort, you scramble out of the room.", Newline)
		NewLine()
		dest := LoudRuns[rand.Intn(len(LoudRuns))]
		Goto(dest, true)
		return false
	}
	if arg == ActEnter {
		if LoudFlag || (!GatesOpen && LowTide) {
			return false
		}
		if GatesOpen && !LowTide {
			return false
		}
		// Room is loud - special input handling
		VFirstLook(ActUnk)
		if Params.Continue != NumUndef {
			Print("The rest of your commands have been lost in the noise.", Newline)
			Params.Continue = NumUndef
		}
		// In the original, this has a special read loop. We simplify.
		return false
	}
	if ActVerb.Norm == "echo" {
		if LoudFlag || (!GatesOpen && LowTide) {
			// Room is already quiet
			Print("echo echo ...", Newline)
			return true
		}
		Print("The acoustics of the room change subtly.", Newline)
		LoudFlag = true
		return true
	}
	return false
}

func DeepCanyonFcn(arg ActArg) bool {
	if arg == ActLook {
		Print("You are on the south edge of a deep canyon. Passages lead off to the east, northwest and southwest. A stairway leads down.", NoNewline)
		if GatesOpen && !LowTide {
			Print(" You can hear a loud roaring sound, like that of rushing water, from below.", NoNewline)
		} else if !GatesOpen && LowTide {
			NewLine()
			return true
		} else {
			Print(" You can hear the sound of flowing water from below.", NoNewline)
		}
		NewLine()
		return true
	}
	return false
}

func BoomRoomFcn(arg ActArg) bool {
	if arg == ActEnd {
		dummy := false
		if ActVerb.Norm == "lamp on" || ActVerb.Norm == "burn" {
			if DirObj == &Candles || DirObj == &Torch || DirObj == &Match {
				dummy = true
			}
		}
		if (Candles.IsIn(Winner) && Candles.Has(FlgOn)) ||
			(Torch.IsIn(Winner) && Torch.Has(FlgOn)) ||
			(Match.IsIn(Winner) && Match.Has(FlgOn)) {
			if dummy {
				Print("How sad for an aspiring adventurer to light a ", NoNewline)
				PrintObject(DirObj)
				Print(" in a room which reeks of gas. Fortunately, there is justice in the world.", Newline)
			} else {
				Print("Oh dear. It appears that the smell coming from this room was coal gas. I would have thought twice about carrying flaming objects in here.", Newline)
			}
			JigsUp("\n      ** BOOOOOOOOOOOM **", false)
			return true
		}
	}
	return false
}

func BatsRoomFcn(arg ActArg) bool {
	if arg == ActLook {
		Print("You are in a small room which has doors only to the east and south.", Newline)
		return true
	}
	if arg == ActEnter && !Dead {
		if Garlic.Location() != Winner && !Garlic.IsIn(Here) {
			VLook(ActUnk)
			NewLine()
			FlyMe()
			return true
		}
	}
	return false
}

func NoObjsFcn(arg ActArg) bool {
	if arg == ActBegin {
		f := Winner.Children
		EmptyHanded = true
		for _, child := range f {
			if Weight(child) > 4 {
				EmptyHanded = false
				break
			}
		}
		if Here == &LowerShaft && Lit {
			ScoreUpd(LightShaft)
			LightShaft = 0
		}
		return false
	}
	return false
}

func CanyonViewFcn(arg ActArg) bool {
	return false
}

func ForestRoomFcn(arg ActArg) bool {
	if arg == ActEnter {
		Queue(IForestRandom, -1).Run = true
		return false
	}
	if arg == ActBegin {
		if (ActVerb.Norm == "climb" || ActVerb.Norm == "climb up") && DirObj == &Tree {
			DoWalk("up")
			return true
		}
	}
	return false
}

func TreeRoomFcn(arg ActArg) bool {
	if arg == ActLook {
		Print("You are about 10 feet above the ground nestled among some large branches. The nearest branch above you is above your reach.", Newline)
		if Path.HasChildren() && len(Path.Children) > 0 {
			Print("On the ground below you can see:  ", NoNewline)
			PrintContents(&Path)
			Print(".", Newline)
		}
		return true
	}
	if arg == ActBegin {
		if (ActVerb.Norm == "climb down") && (DirObj == &Tree || DirObj == &Rooms) {
			DoWalk("down")
			return true
		}
		if (ActVerb.Norm == "climb up" || ActVerb.Norm == "climb") && DirObj == &Tree {
			DoWalk("up")
			return true
		}
		if ActVerb.Norm == "drop" {
			if !IDrop() {
				return true
			}
			if DirObj == &Nest && Egg.IsIn(&Nest) {
				Print("The nest falls to the ground, and the egg spills out of it, seriously damaged.", Newline)
				RemoveCarefully(&Egg)
				BrokenEgg.MoveTo(&Path)
				return true
			}
			if DirObj == &Egg {
				Print("The egg falls to the ground and springs open, seriously damaged.", NoNewline)
				Egg.MoveTo(&Path)
				BadEgg()
				NewLine()
				return true
			}
			if DirObj != Winner && DirObj != &Tree {
				DirObj.MoveTo(&Path)
				Print("The ", NoNewline)
				PrintObject(DirObj)
				Print(" falls to the ground.", Newline)
			}
			return true
		}
	}
	if arg == ActEnter {
		Queue(IForestRandom, -1).Run = true
	}
	return false
}

func DeadFunction(arg ActArg) bool {
	if ActVerb.Norm == "walk" {
		if Here == &TimberRoom && DirObj == ToDirObj("west") {
			Print("You cannot enter in your condition.", Newline)
			return true
		}
		return false
	}
	if ActVerb.Norm == "brief" || ActVerb.Norm == "verbose" || ActVerb.Norm == "super-brief" || ActVerb.Norm == "version" {
		return false
	}
	if ActVerb.Norm == "attack" || ActVerb.Norm == "mung" || ActVerb.Norm == "alarm" || ActVerb.Norm == "swing" {
		Print("All such attacks are vain in your condition.", Newline)
		return true
	}
	if ActVerb.Norm == "open" || ActVerb.Norm == "close" || ActVerb.Norm == "eat" || ActVerb.Norm == "drink" ||
		ActVerb.Norm == "inflate" || ActVerb.Norm == "deflate" || ActVerb.Norm == "turn" || ActVerb.Norm == "burn" ||
		ActVerb.Norm == "tie" || ActVerb.Norm == "untie" || ActVerb.Norm == "rub" {
		Print("Even such an action is beyond your capabilities.", Newline)
		return true
	}
	if ActVerb.Norm == "wait" {
		Print("Might as well. You've got an eternity.", Newline)
		return true
	}
	if ActVerb.Norm == "lamp on" {
		Print("You need no light to guide you.", Newline)
		return true
	}
	if ActVerb.Norm == "score" {
		Print("You're dead! How can you think of your score?", Newline)
		return true
	}
	if ActVerb.Norm == "take" || ActVerb.Norm == "rub" {
		Print("Your hand passes through its object.", Newline)
		return true
	}
	if ActVerb.Norm == "drop" || ActVerb.Norm == "throw" || ActVerb.Norm == "inventory" {
		Print("You have no possessions.", Newline)
		return true
	}
	if ActVerb.Norm == "diagnose" {
		Print("You are dead.", Newline)
		return true
	}
	if ActVerb.Norm == "look" {
		Print("The room looks strange and unearthly", NoNewline)
		if !Here.HasChildren() {
			Print(".", NoNewline)
		} else {
			Print(" and objects appear indistinct.", NoNewline)
		}
		NewLine()
		if !Here.Has(FlgOn) {
			Print("Although there is no light, the room seems dimly illuminated.", Newline)
		}
		NewLine()
		return false
	}
	if ActVerb.Norm == "pray" {
		if Here == &SouthTemple {
			Lamp.Take(FlgInvis)
			Winner.Action = nil
			AlwaysLit = false
			Dead = false
			if Troll.IsIn(&TrollRoom) {
				TrollFlag = false
			}
			Print("From the distance the sound of a lone trumpet is heard. The room becomes very bright and you feel disembodied. In a moment, the brightness fades and you find yourself rising as if from a long sleep, deep in the woods. In the distance you can faintly hear a songbird and the sounds of the forest.", Newline)
			NewLine()
			Goto(&Forest1, true)
			return true
		}
		Print("Your prayers are not heard.", Newline)
		return true
	}
	Print("You can't even do that.", Newline)
	Params.Continue = NumUndef
	return true
}

// ================================================================
// INTERRUPT ROUTINES
// ================================================================

func ICandles() bool {
	Candles.Give(FlgTouch)
	if CandleTableIdx >= len(CandleTable) {
		return true
	}
	tick := CandleTable[CandleTableIdx].(int)
	Queue(ICandles, tick).Run = true
	LightInt(&Candles, CandleTableIdx, tick)
	if tick != 0 {
		CandleTableIdx += 2
	}
	return true
}

func ILantern() bool {
	if LampTableIdx >= len(LampTable) {
		return true
	}
	tick := LampTable[LampTableIdx].(int)
	Queue(ILantern, tick).Run = true
	LightInt(&Lamp, LampTableIdx, tick)
	if tick != 0 {
		LampTableIdx += 2
	}
	return true
}

// LightInt handles light source countdown warnings and expiry
func LightInt(obj *Object, tblIdx, tick int) {
	if tick == 0 {
		obj.Take(FlgOn)
		obj.Give(FlgRMung)
	}
	if IsHeld(obj) || obj.IsIn(Here) {
		if tick == 0 {
			Print("You'd better have more light than from the ", NoNewline)
			PrintObject(obj)
			Print(".", Newline)
		} else {
			// Print the warning message from the table
			var tbl []interface{}
			if obj == &Candles {
				tbl = CandleTable
			} else {
				tbl = LampTable
			}
			if tblIdx+1 < len(tbl) {
				if msg, ok := tbl[tblIdx+1].(string); ok {
					Print(msg, Newline)
				}
			}
		}
	}
}

// ICure heals the player gradually
func ICure() bool {
	s := Winner.Strength
	if s > 0 {
		s = 0
		Winner.Strength = s
	} else if s < 0 {
		s++
		Winner.Strength = s
	}
	if s < 0 {
		if LoadAllowed < LoadMax {
			LoadAllowed += 10
		}
		Queue(ICure, CureWait).Run = true
	} else {
		LoadAllowed = LoadMax
		QueueInt(ICure, false).Run = false
	}
	return false
}

func IMatch() bool {
	Print("The match has gone out.", Newline)
	Match.Take(FlgFlame)
	Match.Take(FlgOn)
	Lit = IsLit(Here, true)
	return true
}

func IXb() bool {
	if !XC {
		if Here == &EnteranceToHades {
			Print("The tension of this ceremony is broken, and the wraiths, amused but shaken at your clumsy attempt, resume their hideous jeering.", Newline)
		}
	}
	XB = false
	return true
}

func IXbh() bool {
	RemoveCarefully(&HotBell)
	Bell.MoveTo(&EnteranceToHades)
	if Here == &EnteranceToHades {
		Print("The bell appears to have cooled down.", Newline)
	}
	return true
}

func IXc() bool {
	XC = false
	IXb()
	return true
}

func ICyclops() bool {
	if CyclopsFlag || Dead {
		return true
	}
	if Here != &CyclopsRoom {
		QueueInt(ICyclops, false).Run = false
		return false
	}
	if AbsInt(CycloWrath) > 5 {
		QueueInt(ICyclops, false).Run = false
		JigsUp("The cyclops, tired of all of your games and trickery, grabs you firmly. As he licks his chops, he says \"Mmm. Just like Mom used to make 'em.\" It's nice to be appreciated.", false)
		return true
	}
	if CycloWrath < 0 {
		CycloWrath--
	} else {
		CycloWrath++
	}
	if !CyclopsFlag {
		idx := AbsInt(CycloWrath) - 2
		if idx >= 0 && idx < len(Cyclomad) {
			Print(Cyclomad[idx], Newline)
		}
	}
	return true
}

func IForestRandom() bool {
	if !ForestRoomQ() {
		QueueInt(IForestRandom, false).Run = false
		return false
	}
	if Prob(15, false) {
		Print("You hear in the distance the chirping of a song bird.", Newline)
	}
	return true
}

// ================================================================
// EXIT FUNCTIONS
// ================================================================

func GratingExitFcn() *Object {
	if GrateRevealed {
		if Grate.Has(FlgOpen) {
			return &GratingRoom
		}
		Print("The grating is closed!", Newline)
		ThisIsIt(&Grate)
		return nil
	}
	Print("You can't go that way.", Newline)
	return nil
}

func TrapDoorExitFcn() *Object {
	if RugMoved {
		if TrapDoor.Has(FlgOpen) {
			return &Cellar
		}
		Print("The trap door is closed.", Newline)
		ThisIsIt(&TrapDoor)
		return nil
	}
	Print("You can't go that way.", Newline)
	return nil
}

func UpChimneyFcn() *Object {
	f := Winner.Children
	if len(f) == 0 {
		Print("Going up empty-handed is a bad idea.", Newline)
		return nil
	}
	// Check if player is carrying at most 1-2 items including the lamp
	count := 0
	for range f {
		count++
	}
	if count <= 2 && Lamp.IsIn(Winner) {
		if !TrapDoor.Has(FlgOpen) {
			TrapDoor.Take(FlgTouch)
		}
		return &Kitchen
	}
	Print("You can't get up there with what you're carrying.", Newline)
	return nil
}

func MazeDiodesFcn() *Object {
	Print("You won't be able to get back up to the tunnel you are going through when it gets to the next room.", Newline)
	NewLine()
	if Here == &Maze2 {
		return &Maze4
	}
	if Here == &Maze7 {
		return &DeadEnd1
	}
	if Here == &Maze9 {
		return &Maze11
	}
	if Here == &Maze12 {
		return &Maze5
	}
	return nil
}

// ================================================================
// PSEUDO FUNCTIONS
// ================================================================

func ChasmPseudo(arg ActArg) bool {
	if ActVerb.Norm == "leap" || (ActVerb.Norm == "put" && DirObj == &Me) {
		Print("You look before leaping, and realize that you would never survive.", Newline)
		return true
	}
	if ActVerb.Norm == "cross" {
		Print("It's too far to jump, and there's no bridge.", Newline)
		return true
	}
	if (ActVerb.Norm == "put" || ActVerb.Norm == "throw off") && IndirObj == &PseudoObject {
		Print("The ", NoNewline)
		PrintObject(DirObj)
		Print(" drops out of sight into the chasm.", Newline)
		RemoveCarefully(DirObj)
		return true
	}
	return false
}

func LakePseudo(arg ActArg) bool {
	if LowTide {
		Print("There's not much lake left....", Newline)
		return true
	}
	if ActVerb.Norm == "cross" {
		Print("It's too wide to cross.", Newline)
		return true
	}
	if ActVerb.Norm == "through" {
		Print("You can't swim in this lake.", Newline)
		return true
	}
	return false
}

func StreamPseudo(arg ActArg) bool {
	if ActVerb.Norm == "swim" || ActVerb.Norm == "through" {
		Print("You can't swim in the stream.", Newline)
		return true
	}
	if ActVerb.Norm == "cross" {
		Print("The other side is a sheer rock cliff.", Newline)
		return true
	}
	return false
}

func DomePseudo(arg ActArg) bool {
	if ActVerb.Norm == "kiss" {
		Print("No.", Newline)
		return true
	}
	return false
}

func GatePseudo(arg ActArg) bool {
	if ActVerb.Norm == "through" {
		DoWalk("in")
		return true
	}
	Print("The gate is protected by an invisible force. It makes your teeth ache to touch it.", Newline)
	return true
}

func DoorPseudo(arg ActArg) bool {
	if ActVerb.Norm == "open" || ActVerb.Norm == "close" {
		Print("The door won't budge.", Newline)
		return true
	}
	if ActVerb.Norm == "through" {
		DoWalk("south")
		return true
	}
	return false
}

func PaintPseudo(arg ActArg) bool {
	if ActVerb.Norm == "mung" {
		Print("Some paint chips away, revealing more paint.", Newline)
		return true
	}
	return false
}

func GasPseudo(arg ActArg) bool {
	if ActVerb.Norm == "breathe" {
		Print("There is too much gas to blow away.", Newline)
		return true
	}
	if ActVerb.Norm == "smell" {
		Print("It smells like coal gas in here.", Newline)
		return true
	}
	return false
}

func ChainPseudo(arg ActArg) bool {
	if ActVerb.Norm == "take" || ActVerb.Norm == "move" {
		Print("The chain is secure.", Newline)
		return true
	}
	if ActVerb.Norm == "raise" || ActVerb.Norm == "lower" {
		Print("Perhaps you should do that to the basket.", Newline)
		return true
	}
	if ActVerb.Norm == "examine" {
		Print("The chain secures a basket within the shaft.", Newline)
		return true
	}
	return false
}

func BarrowDoorFcn2(arg ActArg) bool {
	return false
}
