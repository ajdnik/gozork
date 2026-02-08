package zork

import "math/rand"

func OpenClose(obj *Object, strOpn, strCls string) {
	if G.ActVerb.Norm == "open" {
		if obj.Has(FlgOpen) {
			Print(PickOne(Dummy), Newline)
		} else {
			Print(strOpn, Newline)
			obj.Give(FlgOpen)
		}
	} else if G.ActVerb.Norm == "close" {
		if obj.Has(FlgOpen) {
			Print(strCls, Newline)
			obj.Take(FlgOpen)
		} else {
			Print(PickOne(Dummy), Newline)
		}
	}
}

func LeavesAppear() bool {
	if !Grate.Has(FlgOpen) && !G.GrateRevealed {
		if G.ActVerb.Norm == "move" || G.ActVerb.Norm == "take" {
			Print("In disturbing the pile of leaves, a grating is revealed.", Newline)
		} else {
			Print("With the leaves moved, a grating is revealed.", Newline)
		}
		Grate.Take(FlgInvis)
		G.GrateRevealed = true
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
	if G.Here != &EnteranceToHades {
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
	return G.Here == &Forest1 || G.Here == &Forest2 || G.Here == &Forest3 ||
		G.Here == &Path || G.Here == &UpATree
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
	PrintNumber(G.Score)
	Print(" (total of 350 points), in ", NoNewline)
	PrintNumber(G.Moves)
	if G.Moves == 1 {
		Print(" move.", NoNewline)
	} else {
		Print(" moves.", NoNewline)
	}
	NewLine()
	Print("This gives you the rank of ", NoNewline)
	switch {
	case G.Score == 350:
		Print("Master Adventurer", NoNewline)
	case G.Score > 330:
		Print("Wizard", NoNewline)
	case G.Score > 300:
		Print("Master", NoNewline)
	case G.Score > 200:
		Print("Adventurer", NoNewline)
	case G.Score > 100:
		Print("Junior Adventurer", NoNewline)
	case G.Score > 50:
		Print("Novice Adventurer", NoNewline)
	case G.Score > 25:
		Print("Amateur Adventurer", NoNewline)
	default:
		Print("Beginner", NoNewline)
	}
	Print(".", Newline)
	return true
}

func VDiagnose(arg ActArg) bool {
	ms := FightStrength(false)
	wd := G.Winner.Strength
	rs := ms + wd
	// Check if healing is active
	cureActive := false
	for i := len(G.QueueItms) - 1; i >= 0; i-- {
		if G.QueueItms[i].Key == "ICure" && G.QueueItms[i].Run {
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
	if G.Deaths > 0 {
		Print("You have been killed ", NoNewline)
		if G.Deaths == 1 {
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
	G.Winner = &Adventurer
	if G.Dead {
		NewLine()
		Print("It takes a talented person to be killed while already dead. YOU are such a talent. Unfortunately, it takes a talented person to deal with it. I am not such a talent. Sorry.", Newline)
		return Finish()
	}
	Print(desc, Newline)
	if !G.Lucky {
		Print("Bad luck, huh?", Newline)
	}
	ScoreUpd(-10)
	NewLine()
	Print("    ****  You have died  ****", Newline)
	NewLine()
	if G.Winner.Location().Has(FlgVeh) {
		G.Winner.MoveTo(G.Here)
	}
	if G.Deaths >= 2 {
		Print("You clearly are a suicidal maniac. We don't allow psychotics in the cave, since they may harm other adventurers. Your remains will be installed in the Land of the Living Dead, where your fellow adventurers may gloat over them.", Newline)
		return Finish()
	}
	G.Deaths++
	G.Winner.MoveTo(G.Here)
	if SouthTemple.Has(FlgTouch) {
		Print("As you take your last breath, you feel relieved of your burdens. The feeling passes as you find yourself before the gates of Hell, where the spirits jeer at you and deny you entry. Your senses are disturbed. The objects in the dungeon appear indistinct, bleached of color, even unreal.", Newline)
		NewLine()
		G.Dead = true
		G.TrollFlag = true
		G.AlwaysLit = true
		G.Winner.Action = DeadFunction
		Goto(&EnteranceToHades, true)
	} else {
		Print("Now, let's take a look here...\nWell, you probably deserve another chance. I can't quite fix you up completely, but you can't have everything.", Newline)
		NewLine()
		Goto(&Forest1, true)
	}
	TrapDoor.Take(FlgTouch)
	G.Params.Continue = NumUndef
	RandomizeObjects()
	KillInterrupts()
	return false
}

func RandomizeObjects() {
	if Lamp.IsIn(G.Winner) {
		Lamp.MoveTo(&LivingRoom)
	}
	if Coffin.IsIn(G.Winner) {
		Coffin.MoveTo(&EgyptRoom)
	}
	Sword.TValue = 0
	// Copy children before iterating since MoveTo modifies the slice.
	children := make([]*Object, len(G.Winner.Children))
	copy(children, G.Winner.Children)
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
	QueueInt("IXb", false).Run = false
	QueueInt("IXc", false).Run = false
	QueueInt("ICyclops", false).Run = false
	QueueInt("ILantern", false).Run = false
	QueueInt("ICandles", false).Run = false
	QueueInt("ISword", false).Run = false
	QueueInt("IForestRandom", false).Run = false
	QueueInt("IMatch", false).Run = false
	Match.Take(FlgOn)
	return true
}

// ================================================================
// THE WHITE HOUSE
// ================================================================

func WestHouseFcn(arg ActArg) bool {
	if arg == ActLook {
		Print("You are standing in an open field west of a white house, with a boarded front door.", NoNewline)
		if G.WonGame {
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
	if G.Here == &Kitchen || G.Here == &LivingRoom || G.Here == &Attic {
		if G.ActVerb.Norm == "find" {
			Print("Why not find your brains?", Newline)
			return true
		}
		if G.ActVerb.Norm == "walk around" {
			GoNext(InHouseAround)
			return true
		}
	} else if G.Here != &WestOfHouse && G.Here != &NorthOfHouse && G.Here != &EastOfHouse && G.Here != &SouthOfHouse {
		if G.ActVerb.Norm == "find" {
			if G.Here == &Clearing {
				Print("It seems to be to the west.", Newline)
				return true
			}
			Print("It was here just a minute ago....", Newline)
			return true
		}
		Print("You're not at the house.", Newline)
		return true
	} else if G.ActVerb.Norm == "find" {
		Print("It's right here! Are you blind or something?", Newline)
		return true
	} else if G.ActVerb.Norm == "walk around" {
		GoNext(HouseAround)
		return true
	} else if G.ActVerb.Norm == "examine" {
		Print("The house is a beautiful colonial house which is painted white. It is clear that the owners must have been extremely wealthy.", Newline)
		return true
	} else if G.ActVerb.Norm == "through" || G.ActVerb.Norm == "open" {
		if G.Here == &EastOfHouse {
			if KitchenWindow.Has(FlgOpen) {
				return Goto(&Kitchen, true)
			}
			Print("The window is closed.", Newline)
			ThisIsIt(&KitchenWindow)
			return true
		}
		Print("I can't see how to get in from here.", Newline)
		return true
	} else if G.ActVerb.Norm == "burn" {
		Print("You must be joking.", Newline)
		return true
	}
	return false
}

func GoNext(tbl map[*Object]*Object) int {
	val, ok := tbl[G.Here]
	if !ok {
		return NumUndef
	}
	if !Goto(val, true) {
		return 2
	}
	return 1
}

func BoardFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "take" || G.ActVerb.Norm == "examine" {
		Print("The boards are securely fastened.", Newline)
		return true
	}
	return false
}

func TeethFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "brush" && G.DirObj == &Teeth {
		if G.IndirObj == &Putty && Putty.IsIn(G.Winner) {
			JigsUp("Well, you seem to have been brushing your teeth with some sort of glue. As a result, your mouth gets glued together (with your nose) and you die of respiratory failure.", false)
			return true
		}
		if G.IndirObj == nil {
			Print("Dental hygiene is highly recommended, but I'm not sure what you want to brush them with.", Newline)
			return true
		}
		Print("A nice idea, but with a ", NoNewline)
		PrintObject(G.IndirObj)
		Print("?", Newline)
		return true
	}
	return false
}

func GraniteWallFcn(arg ActArg) bool {
	if G.Here == &NorthTemple {
		if G.ActVerb.Norm == "find" {
			Print("The west wall is solid granite here.", Newline)
			return true
		}
		if G.ActVerb.Norm == "take" || G.ActVerb.Norm == "raise" || G.ActVerb.Norm == "lower" {
			Print("It's solid granite.", Newline)
			return true
		}
	} else if G.Here == &TreasureRoom {
		if G.ActVerb.Norm == "find" {
			Print("The east wall is solid granite here.", Newline)
			return true
		}
		if G.ActVerb.Norm == "take" || G.ActVerb.Norm == "raise" || G.ActVerb.Norm == "lower" {
			Print("It's solid granite.", Newline)
			return true
		}
	} else if G.Here == &SlideRoom {
		if G.ActVerb.Norm == "find" || G.ActVerb.Norm == "read" {
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
	if G.ActVerb.Norm == "find" || G.ActVerb.Norm == "take" {
		Print("The songbird is not here but is probably nearby.", Newline)
		return true
	}
	if G.ActVerb.Norm == "listen" {
		Print("You can't hear the songbird now.", Newline)
		return true
	}
	if G.ActVerb.Norm == "follow" {
		Print("It can't be followed.", Newline)
		return true
	}
	Print("You can't see any songbird here.", Newline)
	return true
}

func MountainRangeFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "climb up" || G.ActVerb.Norm == "climb down" || G.ActVerb.Norm == "climb" {
		Print("Don't you believe me? The mountains are impassable!", Newline)
		return true
	}
	return false
}

func ForestFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "walk around" {
		if G.Here == &WestOfHouse || G.Here == &NorthOfHouse || G.Here == &SouthOfHouse || G.Here == &EastOfHouse {
			Print("You aren't even in the forest.", Newline)
			return true
		}
		GoNext(ForestAround)
		return true
	}
	if G.ActVerb.Norm == "disembark" {
		Print("You will have to specify a direction.", Newline)
		return true
	}
	if G.ActVerb.Norm == "find" {
		Print("You cannot see the forest for the trees.", Newline)
		return true
	}
	if G.ActVerb.Norm == "listen" {
		Print("The pines and the hemlocks seem to be murmuring.", Newline)
		return true
	}
	return false
}

func KitchenWindowFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "open" || G.ActVerb.Norm == "close" {
		G.KitchenWindowFlag = true
		OpenClose(&KitchenWindow,
			"With great effort, you open the window far enough to allow entry.",
			"The window closes (more easily than it opened).")
		return true
	}
	if G.ActVerb.Norm == "examine" && !G.KitchenWindowFlag {
		Print("The window is slightly ajar, but not enough to allow entry.", Newline)
		return true
	}
	if G.ActVerb.Norm == "walk" || G.ActVerb.Norm == "board" || G.ActVerb.Norm == "through" {
		if G.Here == &Kitchen {
			DoWalk(East)
		} else {
			DoWalk(West)
		}
		return true
	}
	if G.ActVerb.Norm == "look inside" {
		Print("You can see ", NoNewline)
		if G.Here == &Kitchen {
			Print("a clear area leading towards a forest.", Newline)
		} else {
			Print("what appears to be a kitchen.", Newline)
		}
		return true
	}
	return false
}

func ChimneyFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "examine" {
		Print("The chimney leads ", NoNewline)
		if G.Here == &Kitchen {
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
	if G.ActVerb.Norm == "tell" {
		Print("The spirits jeer loudly and ignore you.", Newline)
		G.Params.Continue = NumUndef
		return true
	}
	if G.ActVerb.Norm == "exorcise" {
		Print("Only the ceremony itself has any effect.", Newline)
		return true
	}
	if (G.ActVerb.Norm == "attack" || G.ActVerb.Norm == "mung") && G.DirObj == &Ghosts {
		Print("How can you attack a spirit with material objects?", Newline)
		return true
	}
	Print("You seem unable to interact with these spirits.", Newline)
	return true
}

func BasketFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "raise" {
		if G.CageTop {
			Print(PickOne(Dummy), Newline)
		} else {
			RaisedBasket.MoveTo(&ShaftRoom)
			LoweredBasket.MoveTo(&LowerShaft)
			G.CageTop = true
			ThisIsIt(&RaisedBasket)
			Print("The basket is raised to the top of the shaft.", Newline)
		}
		return true
	}
	if G.ActVerb.Norm == "lower" {
		if !G.CageTop {
			Print(PickOne(Dummy), Newline)
		} else {
			RaisedBasket.MoveTo(&LowerShaft)
			LoweredBasket.MoveTo(&ShaftRoom)
			ThisIsIt(&LoweredBasket)
			Print("The basket is lowered to the bottom of the shaft.", Newline)
			G.CageTop = false
			if G.Lit && !IsLit(G.Here, true) {
				G.Lit = false
				Print("It is now pitch black.", Newline)
			}
		}
		return true
	}
	if G.DirObj == &LoweredBasket || G.IndirObj == &LoweredBasket {
		Print("The basket is at the other end of the chain.", Newline)
		return true
	}
	if G.ActVerb.Norm == "take" && (G.DirObj == &RaisedBasket || G.DirObj == &LoweredBasket) {
		Print("The cage is securely fastened to the iron chain.", Newline)
		return true
	}
	return false
}

func BatFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "tell" {
		Fweep(6)
		G.Params.Continue = NumUndef
		return true
	}
	if G.ActVerb.Norm == "take" || G.ActVerb.Norm == "attack" || G.ActVerb.Norm == "mung" {
		if Garlic.Location() == G.Winner || Garlic.IsIn(G.Here) {
			Print("You can't reach him; he's on the ceiling.", Newline)
			return true
		}
		FlyMe()
		return true
	}
	return false
}

func BellFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "ring" {
		if G.Here == &EnteranceToHades && !G.LLDFlag {
			return false
		}
		Print("Ding, dong.", Newline)
		return true
	}
	return false
}

func HotBellFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "take" {
		Print("The bell is very hot and cannot be taken.", Newline)
		return true
	}
	if G.ActVerb.Norm == "rub" || (G.ActVerb.Norm == "ring" && G.IndirObj != nil) {
		if G.IndirObj != nil && G.IndirObj.Has(FlgBurn) {
			Print("The ", NoNewline)
			PrintObject(G.IndirObj)
			Print(" burns and is consumed.", Newline)
			RemoveCarefully(G.IndirObj)
			return true
		}
		if G.IndirObj == &Hands {
			Print("The bell is too hot to touch.", Newline)
			return true
		}
		Print("The heat from the bell is too intense.", Newline)
		return true
	}
	if G.ActVerb.Norm == "pour on" {
		RemoveCarefully(G.DirObj)
		Print("The water cools the bell and is evaporated.", Newline)
		QueueInt("IXbh", false).Run = false
		IXbh()
		return true
	}
	if G.ActVerb.Norm == "ring" {
		Print("The bell is too hot to reach.", Newline)
		return true
	}
	return false
}

func AxeFcn(arg ActArg) bool {
	if G.TrollFlag {
		return false
	}
	return WeaponFunction(&Axe, &Troll)
}

func TrapDoorFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "raise" {
		Perform(ActionVerb{Norm: "open", Orig: "open"}, &TrapDoor, nil)
		return true
	}
	if (G.ActVerb.Norm == "open" || G.ActVerb.Norm == "close") && G.Here == &LivingRoom {
		OpenClose(G.DirObj,
			"The door reluctantly opens to reveal a rickety staircase descending into darkness.",
			"The door swings shut and closes.")
		return true
	}
	if G.ActVerb.Norm == "look under" && G.Here == &LivingRoom {
		if TrapDoor.Has(FlgOpen) {
			Print("You see a rickety staircase descending into darkness.", Newline)
		} else {
			Print("It's closed.", Newline)
		}
		return true
	}
	if G.Here == &Cellar {
		if (G.ActVerb.Norm == "open" || G.ActVerb.Norm == "unlock") && !TrapDoor.Has(FlgOpen) {
			Print("The door is locked from above.", Newline)
			return true
		}
		if G.ActVerb.Norm == "close" && !TrapDoor.Has(FlgOpen) {
			TrapDoor.Take(FlgTouch)
			TrapDoor.Take(FlgOpen)
			Print("The door closes and locks.", Newline)
			return true
		}
		if G.ActVerb.Norm == "open" || G.ActVerb.Norm == "close" {
			Print(PickOne(Dummy), Newline)
			return true
		}
	}
	return false
}

func FrontDoorFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "open" {
		Print("The door cannot be opened.", Newline)
		return true
	}
	if G.ActVerb.Norm == "burn" {
		Print("You cannot burn this door.", Newline)
		return true
	}
	if G.ActVerb.Norm == "mung" {
		Print("You can't seem to damage the door.", Newline)
		return true
	}
	if G.ActVerb.Norm == "look behind" {
		Print("It won't open.", Newline)
		return true
	}
	return false
}

func BarrowDoorFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "open" || G.ActVerb.Norm == "close" {
		Print("The door is too heavy.", Newline)
		return true
	}
	return false
}

func BarrowFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "through" {
		DoWalk(West)
		return true
	}
	return false
}

func BottleFcn(arg ActArg) bool {
	empty := false
	if G.ActVerb.Norm == "throw" && G.DirObj == &Bottle {
		RemoveCarefully(G.DirObj)
		empty = true
		Print("The bottle hits the far wall and shatters.", Newline)
	} else if G.ActVerb.Norm == "mung" {
		empty = true
		RemoveCarefully(G.DirObj)
		Print("A brilliant maneuver destroys the bottle.", Newline)
	} else if G.ActVerb.Norm == "shake" {
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
	if G.ActVerb.Norm == "through" {
		Print("You can't fit through the crack.", Newline)
		return true
	}
	return false
}

func GrateFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "open" && G.IndirObj == &Keys {
		Perform(ActionVerb{Norm: "unlock", Orig: "unlock"}, &Grate, &Keys)
		return true
	}
	if G.ActVerb.Norm == "lock" {
		if G.Here == &GratingRoom {
			G.GrUnlock = false
			Print("The grate is locked.", Newline)
			return true
		}
		if G.Here == &Clearing {
			Print("You can't lock it from this side.", Newline)
			return true
		}
	}
	if G.ActVerb.Norm == "unlock" && G.DirObj == &Grate {
		if G.Here == &GratingRoom && G.IndirObj == &Keys {
			G.GrUnlock = true
			Print("The grate is unlocked.", Newline)
			return true
		}
		if G.Here == &Clearing && G.IndirObj == &Keys {
			Print("You can't reach the lock from here.", Newline)
			return true
		}
		Print("Can you unlock a grating with a ", NoNewline)
		PrintObject(G.IndirObj)
		Print("?", Newline)
		return true
	}
	if G.ActVerb.Norm == "pick" {
		Print("You can't pick the lock.", Newline)
		return true
	}
	if G.ActVerb.Norm == "open" || G.ActVerb.Norm == "close" {
		if G.GrUnlock {
			var openStr string
			if G.Here == &Clearing {
				openStr = "The grating opens."
			} else {
				openStr = "The grating opens to reveal trees above you."
			}
			OpenClose(&Grate, openStr, "The grating is closed.")
			if Grate.Has(FlgOpen) {
				if G.Here != &Clearing && !G.GrateRevealed {
					Print("A pile of leaves falls onto your head and to the ground.", Newline)
					G.GrateRevealed = true
					Leaves.MoveTo(G.Here)
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
	if G.ActVerb.Norm == "put" && G.IndirObj == &Grate {
		if G.DirObj.Size > 20 {
			Print("It won't fit through the grating.", Newline)
		} else {
			G.DirObj.MoveTo(&GratingRoom)
			Print("The ", NoNewline)
			PrintObject(G.DirObj)
			Print(" goes through the grating into the darkness below.", Newline)
		}
		return true
	}
	return false
}

func KnifeFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "take" {
		AtticTable.Take(FlgNoDesc)
		return false
	}
	return false
}

func SkeletonFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "take" || G.ActVerb.Norm == "rub" || G.ActVerb.Norm == "move" ||
		G.ActVerb.Norm == "push" || G.ActVerb.Norm == "raise" || G.ActVerb.Norm == "lower" ||
		G.ActVerb.Norm == "attack" || G.ActVerb.Norm == "kick" || G.ActVerb.Norm == "kiss" {
		Print("A ghost appears in the room and is appalled at your desecration of the remains of a fellow adventurer. He casts a curse on your valuables and banishes them to the Land of the Living Dead. The ghost leaves, muttering obscenities.", Newline)
		Rob(G.Here, &LandOfLivingDead, 100)
		Rob(&Adventurer, &LandOfLivingDead, 0)
		return true
	}
	return false
}

func TorchFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "examine" {
		Print("The torch is burning.", Newline)
		return true
	}
	if G.ActVerb.Norm == "pour on" && G.IndirObj == &Torch {
		Print("The water evaporates before it gets close.", Newline)
		return true
	}
	if G.ActVerb.Norm == "lamp off" && G.DirObj.Has(FlgOn) {
		Print("You nearly burn your hand trying to extinguish the flame.", Newline)
		return true
	}
	return false
}

func RustyKnifeFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "take" {
		if Sword.IsIn(G.Winner) {
			Print("As you touch the rusty knife, your sword gives a single pulse of blinding blue light.", Newline)
		}
		return false
	}
	if (G.IndirObj == &RustyKnife && G.ActVerb.Norm == "attack") ||
		(G.ActVerb.Norm == "swing" && G.DirObj == &RustyKnife && G.IndirObj != nil) {
		RemoveCarefully(&RustyKnife)
		JigsUp("As the knife approaches its victim, your mind is submerged by an overmastering will. Slowly, your hand turns, until the rusty blade is an inch from your neck. The knife seems to sing as it savagely slits your throat.", false)
		return true
	}
	return false
}

func LeafPileFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "count" {
		Print("There are 69,105 leaves here.", Newline)
		return true
	}
	if G.ActVerb.Norm == "burn" {
		LeavesAppear()
		RemoveCarefully(G.DirObj)
		if G.DirObj.IsIn(G.Here) {
			Print("The leaves burn.", Newline)
		} else {
			JigsUp("The leaves burn, and so do you.", false)
		}
		return true
	}
	if G.ActVerb.Norm == "cut" {
		Print("You rustle the leaves around, making quite a mess.", Newline)
		LeavesAppear()
		return true
	}
	if G.ActVerb.Norm == "move" || G.ActVerb.Norm == "take" {
		if G.ActVerb.Norm == "move" {
			Print("Done.", Newline)
		}
		if G.GrateRevealed {
			return false
		}
		LeavesAppear()
		if G.ActVerb.Norm == "take" {
			return false
		}
		return true
	}
	if G.ActVerb.Norm == "look under" && !G.GrateRevealed {
		Print("Underneath the pile of leaves is a grating. As you release the leaves, the grating is once again concealed from view.", Newline)
		return true
	}
	return false
}

func MatchFcn(arg ActArg) bool {
	if (G.ActVerb.Norm == "lamp on" || G.ActVerb.Norm == "burn") && G.DirObj == &Match {
		if G.MatchCount > 0 {
			G.MatchCount--
		}
		if G.MatchCount <= 0 {
			Print("I'm afraid that you have run out of matches.", Newline)
			return true
		}
		if G.Here == &LowerShaft || G.Here == &TimberRoom {
			Print("This room is drafty, and the match goes out instantly.", Newline)
			return true
		}
		Match.Give(FlgFlame)
		Match.Give(FlgOn)
		Queue("IMatch", 2).Run = true
		Print("One of the matches starts to burn.", Newline)
		if !G.Lit {
			G.Lit = true
			VLook(ActUnk)
		}
		return true
	}
	if G.ActVerb.Norm == "lamp off" && Match.Has(FlgFlame) {
		Print("The match is out.", Newline)
		Match.Take(FlgFlame)
		Match.Take(FlgOn)
		G.Lit = IsLit(G.Here, true)
		if !G.Lit {
			Print("It's pitch black in here!", Newline)
		}
		QueueInt("IMatch", false).Run = false
		return true
	}
	if G.ActVerb.Norm == "count" || G.ActVerb.Norm == "open" {
		Print("You have ", NoNewline)
		cnt := G.MatchCount - 1
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
	if G.ActVerb.Norm == "examine" {
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
	if !G.MirrorMung && G.ActVerb.Norm == "rub" {
		if G.IndirObj != nil && G.IndirObj != &Hands {
			Print("You feel a faint tingling transmitted through the ", NoNewline)
			PrintObject(G.IndirObj)
			Print(".", Newline)
			return true
		}
		if G.Here == rm2 {
			rm2 = &MirrorRoom1
		}
		// Swap room contents
		var l1, l2 []*Object
		for _, c := range G.Here.Children {
			l1 = append(l1, c)
		}
		for _, c := range rm2.Children {
			l2 = append(l2, c)
		}
		for _, c := range l1 {
			c.MoveTo(rm2)
		}
		for _, c := range l2 {
			c.MoveTo(G.Here)
		}
		Goto(rm2, false)
		Print("There is a rumble from deep within the earth and the room shakes.", Newline)
		return true
	}
	if G.ActVerb.Norm == "look inside" || G.ActVerb.Norm == "examine" {
		if G.MirrorMung {
			Print("The mirror is broken into many pieces.", Newline)
		} else {
			Print("There is an ugly person staring back at you.", Newline)
		}
		return true
	}
	if G.ActVerb.Norm == "take" {
		Print("The mirror is many times your size. Give up.", Newline)
		return true
	}
	if G.ActVerb.Norm == "mung" || G.ActVerb.Norm == "throw" || G.ActVerb.Norm == "attack" {
		if G.MirrorMung {
			Print("Haven't you done enough damage already?", Newline)
		} else {
			G.MirrorMung = true
			G.Lucky = false
			Print("You have broken the mirror. I hope you have a seven years' supply of good luck handy.", Newline)
		}
		return true
	}
	return false
}

func PaintingFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "mung" {
		G.DirObj.TValue = 0
		G.DirObj.LongDesc = "There is a worthless piece of canvas here."
		Print("Congratulations! Unlike the other vandals, who merely stole the artist's masterpieces, you have destroyed one.", Newline)
		return true
	}
	return false
}

func CandlesFcn(arg ActArg) bool {
	if !Candles.Has(FlgTouch) {
		Queue("ICandles", -1).Run = true
	}
	if G.IndirObj == &Candles {
		return false
	}
	if G.ActVerb.Norm == "lamp on" || G.ActVerb.Norm == "burn" {
		if Candles.Has(FlgRMung) {
			Print("Alas, there's not much left of the candles. Certainly not enough to burn.", Newline)
			return true
		}
		if G.IndirObj == nil {
			if Match.Has(FlgFlame) {
				Print("(with the match)", Newline)
				Perform(ActionVerb{Norm: "lamp on", Orig: "light"}, &Candles, &Match)
				return true
			}
			Print("You should say what to light them with.", Newline)
			return true
		}
		if G.IndirObj == &Match && Match.Has(FlgOn) {
			Print("The candles are ", NoNewline)
			if Candles.Has(FlgOn) {
				Print("already lit.", Newline)
			} else {
				Candles.Give(FlgOn)
				Print("lit.", Newline)
				Queue("ICandles", -1).Run = true
			}
			return true
		}
		if G.IndirObj == &Torch {
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
	if G.ActVerb.Norm == "count" {
		Print("Let's see, how many objects in a pair? Don't tell me, I'll get it.", Newline)
		return true
	}
	if G.ActVerb.Norm == "lamp off" {
		QueueInt("ICandles", false).Run = false
		if Candles.Has(FlgOn) {
			Print("The flame is extinguished.", NoNewline)
			Candles.Take(FlgOn)
			Candles.Give(FlgTouch)
			G.Lit = IsLit(G.Here, true)
			if !G.Lit {
				Print(" It's really dark in here....", NoNewline)
			}
			NewLine()
			return true
		}
		Print("The candles are not lighted.", Newline)
		return true
	}
	if G.ActVerb.Norm == "put" && G.IndirObj != nil && G.IndirObj.Has(FlgBurn) {
		Print("That wouldn't be smart.", Newline)
		return true
	}
	if G.ActVerb.Norm == "examine" {
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
	if G.ActVerb.Norm == "take" {
		Print("A force keeps you from taking the bodies.", Newline)
		return true
	}
	if G.ActVerb.Norm == "mung" || G.ActVerb.Norm == "burn" {
		JigsUp("The voice of the guardian of the dungeon booms out from the darkness, \"Your disrespect costs you your life!\" and places your head on a sharp pole.", false)
		return true
	}
	return false
}

func BlackBookFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "open" {
		Print("The book is already open to page 569.", Newline)
		return true
	}
	if G.ActVerb.Norm == "close" {
		Print("As hard as you try, the book cannot be closed.", Newline)
		return true
	}
	if G.ActVerb.Norm == "turn" {
		Print("Beside page 569, there is only one other page with any legible printing on it. Most of it is unreadable, but the subject seems to be the banishment of evil. Apparently, certain noises, lights, and prayers are efficacious in this regard.", Newline)
		return true
	}
	if G.ActVerb.Norm == "burn" {
		RemoveCarefully(G.DirObj)
		JigsUp("A booming voice says \"Wrong, cretin!\" and you notice that you have turned into a pile of dust. How, I can't imagine.", false)
		return true
	}
	return false
}

func SceptreFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "wave" || G.ActVerb.Norm == "raise" {
		if G.Here == &AragainFalls || G.Here == &EndOfRainbow {
			if !G.RainbowFlag {
				PotOfGold.Take(FlgInvis)
				Print("Suddenly, the rainbow appears to become solid and, I venture, walkable (I think the giveaway was the stairs and bannister).", Newline)
				if G.Here == &EndOfRainbow && PotOfGold.IsIn(&EndOfRainbow) {
					Print("A shimmering pot of gold appears at the end of the rainbow.", Newline)
				}
				G.RainbowFlag = true
			} else {
				Rob(&OnRainbow, &Wall, 0)
				Print("The rainbow seems to have become somewhat run-of-the-mill.", Newline)
				G.RainbowFlag = false
				return true
			}
			return true
		}
		if G.Here == &OnRainbow {
			G.RainbowFlag = false
			JigsUp("The structural integrity of the rainbow is severely compromised, leaving you hanging in midair, supported only by water vapor. Bye.", false)
			return true
		}
		Print("A dazzling display of color briefly emanates from the sceptre.", Newline)
		return true
	}
	return false
}

func SlideFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "through" || G.ActVerb.Norm == "climb up" || G.ActVerb.Norm == "climb down" || G.ActVerb.Norm == "climb" ||
		(G.ActVerb.Norm == "put" && G.DirObj == &Me) {
		if G.Here == &Cellar {
			DoWalk(West)
			return true
		}
		Print("You tumble down the slide....", Newline)
		Goto(&Cellar, true)
		return true
	}
	if G.ActVerb.Norm == "put" {
		Slider(G.DirObj)
		return true
	}
	return false
}

func SandwichBagFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "smell" && Lunch.IsIn(G.DirObj) {
		Print("It smells of hot peppers.", Newline)
		return true
	}
	return false
}

func ToolChestFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "examine" {
		Print("The chests are all empty.", Newline)
		return true
	}
	if G.ActVerb.Norm == "take" || G.ActVerb.Norm == "open" || G.ActVerb.Norm == "put" {
		RemoveCarefully(&ToolChest)
		Print("The chests are so rusty and corroded that they crumble when you touch them.", Newline)
		return true
	}
	return false
}

func ButtonFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "read" {
		Print("They're greek to you.", Newline)
		return true
	}
	if G.ActVerb.Norm == "push" {
		if G.DirObj == &BlueButton {
			if G.WaterLevel == 0 {
				Leak.Take(FlgInvis)
				Print("There is a rumbling sound and a stream of water appears to burst from the east wall of the room (apparently, a leak has occurred in a pipe).", Newline)
				G.WaterLevel = 1
				Queue("IMaintRoom", -1).Run = true
				return true
			}
			Print("The blue button appears to be jammed.", Newline)
			return true
		}
		if G.DirObj == &RedButton {
			Print("The lights within the room ", NoNewline)
			if G.Here.Has(FlgOn) {
				G.Here.Take(FlgOn)
				Print("shut off.", Newline)
			} else {
				G.Here.Give(FlgOn)
				Print("come on.", Newline)
			}
			return true
		}
		if G.DirObj == &BrownButton {
			DamRoom.Take(FlgTouch)
			G.GateFlag = false
			Print("Click.", Newline)
			return true
		}
		if G.DirObj == &YellowButton {
			DamRoom.Take(FlgTouch)
			G.GateFlag = true
			Print("Click.", Newline)
			return true
		}
		return true
	}
	return false
}

func LeakFcn(arg ActArg) bool {
	if G.WaterLevel > 0 {
		if (G.ActVerb.Norm == "put" || G.ActVerb.Norm == "put on") && G.DirObj == &Putty {
			FixMaintLeak()
			return true
		}
		if G.ActVerb.Norm == "plug" {
			if G.IndirObj == &Putty {
				FixMaintLeak()
				return true
			}
			WithTell(G.IndirObj)
			return true
		}
	}
	return false
}

func MachineFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "take" && G.DirObj == &Machine {
		Print("It is far too large to carry.", Newline)
		return true
	}
	if G.ActVerb.Norm == "open" {
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
	if G.ActVerb.Norm == "close" {
		if Machine.Has(FlgOpen) {
			Print("The lid closes.", Newline)
			Machine.Take(FlgOpen)
		} else {
			Print(PickOne(Dummy), Newline)
		}
		return true
	}
	if G.ActVerb.Norm == "lamp on" {
		if G.IndirObj == nil {
			Print("It's not clear how to turn it on with your bare hands.", Newline)
		} else {
			Perform(ActionVerb{Norm: "turn", Orig: "turn"}, &MachineSwitch, G.IndirObj)
			return true
		}
		return true
	}
	return false
}

func MachineSwitchFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "turn" {
		if G.IndirObj == &Screwdriver {
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
			PrintObject(G.IndirObj)
			Print(" won't do.", Newline)
		}
		return true
	}
	return false
}

func PuttyFcn(arg ActArg) bool {
	if (G.ActVerb.Norm == "oil" && G.IndirObj == &Putty) || (G.ActVerb.Norm == "put" && G.DirObj == &Putty) {
		Print("The all-purpose gunk isn't a lubricant.", Newline)
		return true
	}
	return false
}

func TubeFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "put" && G.IndirObj == &Tube {
		Print("The tube refuses to accept anything.", Newline)
		return true
	}
	if G.ActVerb.Norm == "squeeze" {
		if G.DirObj.Has(FlgOpen) && Putty.IsIn(G.DirObj) {
			Putty.MoveTo(G.Winner)
			Print("The viscous material oozes into your hand.", Newline)
			return true
		}
		if G.DirObj.Has(FlgOpen) {
			Print("The tube is apparently empty.", Newline)
			return true
		}
		Print("The tube is closed.", Newline)
		return true
	}
	return false
}

func SwordFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "take" && G.Winner == &Adventurer {
		Queue("ISword", -1).Run = true
		return false
	}
	if G.ActVerb.Norm == "examine" {
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
	if G.ActVerb.Norm == "throw" {
		Print("The lamp has smashed into the floor, and the light has gone out.", Newline)
		QueueInt("ILantern", false).Run = false
		RemoveCarefully(&Lamp)
		BrokenLamp.MoveTo(G.Here)
		return true
	}
	if G.ActVerb.Norm == "lamp on" {
		if Lamp.Has(FlgRMung) {
			Print("A burned-out lamp won't light.", Newline)
			return true
		}
		itm := QueueInt("ILantern", false)
		if itm.Tick <= 0 {
			// First activation or timer expired: initialize countdown
			itm.Tick = -1
		}
		// Otherwise resume from where we left off
		itm.Run = true
		return false
	}
	if G.ActVerb.Norm == "lamp off" {
		if Lamp.Has(FlgRMung) {
			Print("The lamp has already burned out.", Newline)
			return true
		}
		QueueInt("ILantern", false).Run = false
		return false
	}
	if G.ActVerb.Norm == "examine" {
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
	if G.ActVerb.Norm == "take" && G.DirObj == &Mailbox {
		Print("It is securely anchored.", Newline)
		return true
	}
	return false
}

// ================================================================
// TROLL
// ================================================================

func TrollFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "tell" {
		G.Params.Continue = NumUndef
		Print("The troll isn't much of a conversationalist.", Newline)
		return true
	}
	if arg == ActArg(FBusy) {
		if Axe.IsIn(&Troll) {
			return false
		}
		if Axe.IsIn(G.Here) && Prob(75, true) {
			Axe.Give(FlgNoDesc)
			Axe.Take(FlgWeapon)
			Axe.MoveTo(&Troll)
			Troll.LongDesc = "A nasty-looking troll, brandishing a bloody axe, blocks all passages out of the room."
			if Troll.IsIn(G.Here) {
				Print("The troll, angered and humiliated, recovers his weapon. He appears to have an axe to grind with you.", Newline)
			}
			return true
		}
		if Troll.IsIn(G.Here) {
			Troll.LongDesc = "A pathetically babbling troll is here."
			Print("The troll, disarmed, cowers in terror, pleading for his life in the guttural tongue of the trolls.", Newline)
			return true
		}
		return false
	}
	if arg == ActArg(FDead) {
		if Axe.IsIn(&Troll) {
			Axe.MoveTo(G.Here)
			Axe.Take(FlgNoDesc)
			Axe.Give(FlgWeapon)
		}
		G.TrollFlag = true
		return true
	}
	if arg == ActArg(FUnconscious) {
		Troll.Take(FlgFight)
		if Axe.IsIn(&Troll) {
			Axe.MoveTo(G.Here)
			Axe.Take(FlgNoDesc)
			Axe.Give(FlgWeapon)
		}
		Troll.LongDesc = "An unconscious troll is sprawled on the floor. All passages out of the room are open."
		G.TrollFlag = true
		return true
	}
	if arg == ActArg(FConscious) {
		if Troll.IsIn(G.Here) {
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
		G.TrollFlag = false
		return true
	}
	if arg == ActArg(FFirst) {
		if Prob(33, false) {
			Troll.Give(FlgFight)
			G.Params.Continue = NumUndef
			return true
		}
		return false
	}
	// Default (no mode - regular verbs)
	if G.ActVerb.Norm == "examine" {
		Print(Troll.LongDesc, Newline)
		return true
	}
	if (G.ActVerb.Norm == "throw" || G.ActVerb.Norm == "give") && G.DirObj != nil && G.IndirObj == &Troll {
		Awaken(&Troll)
		if G.ActVerb.Norm == "throw" || G.ActVerb.Norm == "give" {
			if G.DirObj == &Axe && Axe.IsIn(G.Winner) {
				Print("The troll scratches his head in confusion, then takes the axe.", Newline)
				Troll.Give(FlgFight)
				Axe.MoveTo(&Troll)
				return true
			}
			if G.DirObj == &Troll || G.DirObj == &Axe {
				Print("You would have to get the ", NoNewline)
				PrintObject(G.DirObj)
				Print(" first, and that seems unlikely.", Newline)
				return true
			}
			if G.ActVerb.Norm == "throw" {
				Print("The troll, who is remarkably coordinated, catches the ", NoNewline)
				PrintObject(G.DirObj)
			} else {
				Print("The troll, who is not overly proud, graciously accepts the gift", NoNewline)
			}
			if Prob(20, false) && (G.DirObj == &Knife || G.DirObj == &Sword || G.DirObj == &Axe) {
				RemoveCarefully(G.DirObj)
				Print(" and eats it hungrily. Poor troll, he dies from an internal hemorrhage and his carcass disappears in a sinister black fog.", Newline)
				RemoveCarefully(&Troll)
				TrollFcn(ActArg(FDead))
				G.TrollFlag = true
			} else if G.DirObj == &Knife || G.DirObj == &Sword || G.DirObj == &Axe {
				G.DirObj.MoveTo(G.Here)
				Print(" and, being for the moment sated, throws it back. Fortunately, the troll has poor control, and the ", NoNewline)
				PrintObject(G.DirObj)
				Print(" falls to the floor. He does not look pleased.", Newline)
				Troll.Give(FlgFight)
			} else {
				Print(" and not having the most discriminating tastes, gleefully eats it.", Newline)
				RemoveCarefully(G.DirObj)
			}
			return true
		}
	}
	if G.ActVerb.Norm == "take" || G.ActVerb.Norm == "move" {
		Awaken(&Troll)
		if G.ActVerb.Norm == "take" || G.ActVerb.Norm == "move" {
			Print("The troll spits in your face, grunting \"Better luck next time\" in a rather barbarous accent.", Newline)
			return true
		}
	}
	if G.ActVerb.Norm == "mung" {
		Awaken(&Troll)
		Print("The troll laughs at your puny gesture.", Newline)
		return true
	}
	if G.ActVerb.Norm == "listen" {
		Print("Every so often the troll says something, probably uncomplimentary, in his guttural tongue.", Newline)
		return true
	}
	if G.TrollFlag && G.ActVerb.Norm == "hello" {
		Print("Unfortunately, the troll can't hear you.", Newline)
		return true
	}
	return false
}

// ================================================================
// CYCLOPS
// ================================================================

func CyclopsFcn(arg ActArg) bool {
	count := G.CycloWrath
	if G.Winner == &Cyclops {
		if G.CyclopsFlag {
			Print("No use talking to him. He's fast asleep.", Newline)
			return true
		}
		if G.ActVerb.Norm == "odysseus" {
			G.Winner = &Adventurer
			Perform(ActionVerb{Norm: "odysseus", Orig: "odysseus"}, nil, nil)
			return true
		}
		Print("The cyclops prefers eating to making conversation.", Newline)
		return true
	}
	if G.CyclopsFlag {
		if G.ActVerb.Norm == "examine" {
			Print("The cyclops is sleeping like a baby, albeit a very ugly one.", Newline)
			return true
		}
		if G.ActVerb.Norm == "alarm" || G.ActVerb.Norm == "kick" || G.ActVerb.Norm == "attack" || G.ActVerb.Norm == "burn" || G.ActVerb.Norm == "mung" {
			Print("The cyclops yawns and stares at the thing that woke him up.", Newline)
			G.CyclopsFlag = false
			Cyclops.Give(FlgFight)
			if count < 0 {
				G.CycloWrath = -count
			} else {
				G.CycloWrath = count
			}
			return true
		}
		return false
	}
	if G.ActVerb.Norm == "examine" {
		Print("A hungry cyclops is standing at the foot of the stairs.", Newline)
		return true
	}
	if G.ActVerb.Norm == "give" && G.IndirObj == &Cyclops {
		if G.DirObj == &Lunch {
			if count >= 0 {
				RemoveCarefully(&Lunch)
				Print("The cyclops says \"Mmm Mmm. I love hot peppers! But oh, could I use a drink. Perhaps I could drink the blood of that thing.\"  From the gleam in his eye, it could be surmised that you are \"that thing\".", Newline)
				G.CycloWrath = MinInt(-1, -count)
			}
			Queue("ICyclops", -1).Run = true
			return true
		}
		if G.DirObj == &Water || (G.DirObj == &Bottle && Water.IsIn(&Bottle)) {
			if count < 0 {
				RemoveCarefully(&Water)
				Bottle.MoveTo(G.Here)
				Bottle.Give(FlgOpen)
				Cyclops.Take(FlgFight)
				Print("The cyclops takes the bottle, checks that it's open, and drinks the water. A moment later, he lets out a yawn that nearly blows you over, and then falls fast asleep (what did you put in that drink, anyway?).", Newline)
				G.CyclopsFlag = true
			} else {
				Print("The cyclops apparently is not thirsty and refuses your generous offer.", Newline)
			}
			return true
		}
		if G.DirObj == &Garlic {
			Print("The cyclops may be hungry, but there is a limit.", Newline)
			return true
		}
		Print("The cyclops is not so stupid as to eat THAT!", Newline)
		return true
	}
	if G.ActVerb.Norm == "throw" || G.ActVerb.Norm == "attack" || G.ActVerb.Norm == "mung" {
		Queue("ICyclops", -1).Run = true
		if G.ActVerb.Norm == "mung" {
			Print("\"Do you think I'm as stupid as my father was?\", he says, dodging.", Newline)
		} else {
			Print("The cyclops shrugs but otherwise ignores your pitiful attempt.", Newline)
			if G.ActVerb.Norm == "throw" {
				G.DirObj.MoveTo(G.Here)
			}
			return true
		}
		return true
	}
	if G.ActVerb.Norm == "take" {
		Print("The cyclops doesn't take kindly to being grabbed.", Newline)
		return true
	}
	if G.ActVerb.Norm == "tie" {
		Print("You cannot tie the cyclops, though he is fit to be tied.", Newline)
		return true
	}
	if G.ActVerb.Norm == "listen" {
		Print("You can hear his stomach rumbling.", Newline)
		return true
	}
	return false
}

// ================================================================
// THIEF / ROBBER
// ================================================================

func DumbContainerFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "open" || G.ActVerb.Norm == "close" || G.ActVerb.Norm == "look inside" {
		Print("You can't do that.", Newline)
		return true
	}
	if G.ActVerb.Norm == "examine" {
		Print("It looks pretty much like a ", NoNewline)
		PrintObject(G.DirObj)
		Print(".", Newline)
		return true
	}
	return false
}

func ChaliceFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "take" {
		if G.DirObj.IsIn(&TreasureRoom) && Thief.IsIn(&TreasureRoom) && Thief.Has(FlgFight) && !Thief.Has(FlgInvis) && Thief.LongDesc != RobberUDesc {
			Print("You'd be stabbed in the back first.", Newline)
			return true
		}
		return false
	}
	if G.ActVerb.Norm == "put" && G.IndirObj == &Chalice {
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
	if G.ActVerb.Norm == "open" || G.ActVerb.Norm == "close" {
		Print("The "+str+" are safely inside; there's no need to do that.", Newline)
		return true
	}
	if G.ActVerb.Norm == "look inside" || G.ActVerb.Norm == "examine" {
		Print("There are lots of "+str+" in there.", Newline)
		return true
	}
	if G.ActVerb.Norm == "put" && G.IndirObj == obj {
		Print("Don't be silly. It wouldn't be a ", NoNewline)
		PrintObject(obj)
		Print(" anymore.", Newline)
		return true
	}
	return false
}

func GarlicFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "eat" {
		RemoveCarefully(G.DirObj)
		Print("What the heck! You won't make friends this way, but nobody around here is too friendly anyhow. Gulp!", Newline)
		return true
	}
	return false
}

func BatDescFcn(arg ActArg) bool {
	if Garlic.Location() == G.Winner || Garlic.IsIn(G.Here) {
		Print("In the corner of the room on the ceiling is a large vampire bat who is obviously deranged and holding his nose.", Newline)
	} else {
		Print("A large vampire bat, hanging from the ceiling, swoops down at you!", Newline)
	}
	return true
}

func TrophyCaseFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "take" && G.DirObj == &TrophyCase {
		Print("The trophy case is securely fastened to the wall.", Newline)
		return true
	}
	return false
}

func BoardedWindowFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "open" {
		Print("The windows are boarded and can't be opened.", Newline)
		return true
	}
	if G.ActVerb.Norm == "mung" {
		Print("You can't break the windows open.", Newline)
		return true
	}
	return false
}

func NailsPseudo(arg ActArg) bool {
	if G.ActVerb.Norm == "take" {
		Print("The nails, deeply imbedded in the door, cannot be removed.", Newline)
		return true
	}
	return false
}

func CliffObjectFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "leap" || (G.ActVerb.Norm == "put" && G.DirObj == &Me) {
		Print("That would be very unwise. Perhaps even fatal.", Newline)
		return true
	}
	if G.IndirObj == &ClimbableCliff {
		if G.ActVerb.Norm == "put" || G.ActVerb.Norm == "throw off" {
			Print("The ", NoNewline)
			PrintObject(G.DirObj)
			Print(" tumbles into the river and is seen no more.", Newline)
			RemoveCarefully(G.DirObj)
			return true
		}
	}
	return false
}

func WhiteCliffFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "climb up" || G.ActVerb.Norm == "climb down" || G.ActVerb.Norm == "climb" {
		Print("The cliff is too steep for climbing.", Newline)
		return true
	}
	return false
}

func RainbowFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "cross" || G.ActVerb.Norm == "through" {
		if G.Here == &CanyonView {
			Print("From here?!?", Newline)
			return true
		}
		if G.RainbowFlag {
			if G.Here == &AragainFalls {
				Goto(&EndOfRainbow, true)
			} else if G.Here == &EndOfRainbow {
				Goto(&AragainFalls, true)
			} else {
				Print("You'll have to say which way...", Newline)
			}
		} else {
			Print("Can you walk on water vapor?", Newline)
		}
		return true
	}
	if G.ActVerb.Norm == "look under" {
		Print("The Frigid River flows under the rainbow.", Newline)
		return true
	}
	return false
}

func RopeFcn(arg ActArg) bool {
	if G.Here != &DomeRoom {
		G.DomeFlag = false
		if G.ActVerb.Norm == "tie" {
			Print("You can't tie the rope to that.", Newline)
			return true
		}
		return false
	}
	if G.ActVerb.Norm == "tie" {
		if G.IndirObj == &Railing {
			if G.DomeFlag {
				Print("The rope is already tied to it.", Newline)
			} else {
				Print("The rope drops over the side and comes within ten feet of the floor.", Newline)
				G.DomeFlag = true
				Rope.Give(FlgNoDesc)
				rloc := Rope.Location()
				if rloc == nil || !rloc.IsIn(&Rooms) {
					Rope.MoveTo(G.Here)
				}
			}
			return true
		}
		return false
	}
	if G.ActVerb.Norm == "climb down" && (G.DirObj == &Rope || G.DirObj == &Rooms) && G.DomeFlag {
		DoWalk(Down)
		return true
	}
	if G.ActVerb.Norm == "tie up" && G.IndirObj == &Rope {
		if G.DirObj.Has(FlgActor) {
			if G.DirObj.Strength < 0 {
				Print("Your attempt to tie up the ", NoNewline)
				PrintObject(G.DirObj)
				Print(" awakens him.", NoNewline)
				Awaken(G.DirObj)
			} else {
				Print("The ", NoNewline)
				PrintObject(G.DirObj)
				Print(" struggles and you cannot tie him up.", Newline)
			}
		} else {
			Print("Why would you tie up a ", NoNewline)
			PrintObject(G.DirObj)
			Print("?", Newline)
		}
		return true
	}
	if G.ActVerb.Norm == "untie" {
		if G.DomeFlag {
			G.DomeFlag = false
			Rope.Take(FlgNoDesc)
			Print("The rope is now untied.", Newline)
		} else {
			Print("It is not tied to anything.", Newline)
		}
		return true
	}
	if G.ActVerb.Norm == "drop" && G.Here == &DomeRoom && !G.DomeFlag {
		Rope.MoveTo(&TorchRoom)
		Print("The rope drops gently to the floor below.", Newline)
		return true
	}
	if G.ActVerb.Norm == "take" {
		if G.DomeFlag {
			Print("The rope is tied to the railing.", Newline)
			return true
		}
	}
	return false
}

func EggObjectFcn(arg ActArg) bool {
	if (G.ActVerb.Norm == "open" || G.ActVerb.Norm == "mung") && G.DirObj == &Egg {
		if G.DirObj.Has(FlgOpen) {
			Print("The egg is already open.", Newline)
			return true
		}
		if G.IndirObj == nil {
			Print("You have neither the tools nor the expertise.", Newline)
			return true
		}
		if G.IndirObj == &Hands {
			Print("I doubt you could do that without damaging it.", Newline)
			return true
		}
		if G.IndirObj.Has(FlgWeapon) || G.IndirObj.Has(FlgTool) || G.ActVerb.Norm == "mung" {
			Print("The egg is now open, but the clumsiness of your attempt has seriously compromised its esthetic appeal.", NoNewline)
			BadEgg()
			NewLine()
			return true
		}
		if G.DirObj.Has(FlgFight) {
			Print("Not to say that using the ", NoNewline)
			PrintObject(G.IndirObj)
			Print(" isn't original too...", Newline)
			return true
		}
		Print("The concept of using a ", NoNewline)
		PrintObject(G.IndirObj)
		Print(" is certainly original.", Newline)
		G.DirObj.Give(FlgFight)
		return true
	}
	if G.ActVerb.Norm == "climb on" || G.ActVerb.Norm == "hatch" {
		Print("There is a noticeable crunch from beneath you, and inspection reveals that the egg is lying open, badly damaged.", NoNewline)
		BadEgg()
		NewLine()
		return true
	}
	if G.ActVerb.Norm == "open" || G.ActVerb.Norm == "mung" || G.ActVerb.Norm == "throw" {
		if G.ActVerb.Norm == "throw" {
			G.DirObj.MoveTo(G.Here)
		}
		Print("Your rather indelicate handling of the egg has caused it some damage, although you have succeeded in opening it.", NoNewline)
		BadEgg()
		NewLine()
		return true
	}
	return false
}

func CanaryObjectFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "wind" {
		if G.DirObj == &Canary {
			if !G.SingSong && ForestRoomQ() {
				Print("The canary chirps, slightly off-key, an aria from a forgotten opera. From out of the greenery flies a lovely songbird. It perches on a limb just over your head and opens its beak to sing. As it does so a beautiful brass bauble drops from its mouth, bounces off the top of your head, and lands glimmering in the grass. As the canary winds down, the songbird flies away.", Newline)
				G.SingSong = true
				dest := G.Here
				if G.Here == &UpATree {
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
	if G.ActVerb.Norm == "raise" {
		Print("The rug is too heavy to lift", NoNewline)
		if G.RugMoved {
			Print(".", Newline)
		} else {
			Print(", but in trying to take it you have noticed an irregularity beneath it.", Newline)
		}
		return true
	}
	if G.ActVerb.Norm == "move" || G.ActVerb.Norm == "push" {
		if G.RugMoved {
			Print("Having moved the carpet previously, you find it impossible to move it again.", Newline)
		} else {
			Print("With a great effort, the rug is moved to one side of the room, revealing the dusty cover of a closed trap door.", Newline)
			TrapDoor.Take(FlgInvis)
			ThisIsIt(&TrapDoor)
			G.RugMoved = true
		}
		return true
	}
	if G.ActVerb.Norm == "take" {
		Print("The rug is extremely heavy and cannot be carried.", Newline)
		return true
	}
	if G.ActVerb.Norm == "look under" && !G.RugMoved && !TrapDoor.Has(FlgOpen) {
		Print("Underneath the rug is a closed trap door. As you drop the corner of the rug, the trap door is once again concealed from view.", Newline)
		return true
	}
	if G.ActVerb.Norm == "climb on" {
		if !G.RugMoved && !TrapDoor.Has(FlgOpen) {
			Print("As you sit, you notice an irregularity underneath it. Rather than be uncomfortable, you stand up again.", Newline)
		} else {
			Print("I suppose you think it's a magic carpet?", Newline)
		}
		return true
	}
	return false
}

func SandFunction(arg ActArg) bool {
	if G.ActVerb.Norm == "dig" && G.IndirObj == &Shovel {
		G.BeachDig++
		if G.BeachDig > 3 {
			G.BeachDig = -1
			if Scarab.IsIn(G.Here) {
				Scarab.Give(FlgInvis)
			}
			JigsUp("The hole collapses, smothering you.", false)
			return true
		}
		if G.BeachDig == 3 {
			if Scarab.Has(FlgInvis) {
				Print("You can see a scarab here in the sand.", Newline)
				ThisIsIt(&Scarab)
				Scarab.Take(FlgInvis)
			}
		} else {
			Print(BDigs[G.BeachDig], Newline)
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
		if G.ActVerb.Norm == "climb up" && G.DirObj == &Stairs {
			DoWalk(Up)
			return true
		}
	}
	return false
}

func LivingRoomFcn(arg ActArg) bool {
	if arg == ActLook {
		Print("You are in the living room. There is a doorway to the east", NoNewline)
		if G.MagicFlag {
			Print(". To the west is a cyclops-shaped opening in an old wooden door, above which is some strange gothic lettering, ", NoNewline)
		} else {
			Print(", a wooden door with strange gothic lettering to the west, which appears to be nailed shut, ", NoNewline)
		}
		Print("a trophy case, ", NoNewline)
		if G.RugMoved && TrapDoor.Has(FlgOpen) {
			Print("and a rug lying beside an open trap door.", NoNewline)
		} else if G.RugMoved {
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
		if G.ActVerb.Norm == "take" || (G.ActVerb.Norm == "put" && G.IndirObj == &TrophyCase) {
			if G.DirObj.IsIn(&TrophyCase) {
				TouchAll(G.DirObj)
			}
			G.Score = G.BaseScore + OtvalFrob(&TrophyCase)
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
		if G.ActVerb.Norm == "enter" || (G.ActVerb.Norm == "walk" && G.Params.HasWalkDir && (G.Params.WalkDir == West || G.Params.WalkDir == In)) || (G.ActVerb.Norm == "through" && G.DirObj == &Barrow) {
			Print("Inside the Barrow\nAs you enter the barrow, the door closes inexorably behind you. Around you it is dark, but ahead is an enormous cavern, brightly lit. Through its center runs a wide stream. Spanning the stream is a small wooden footbridge, and beyond a path leads into a dark tunnel. Above the bridge, floating in the air, is a large sign. It reads:  All ye who stand before this bridge have completed a great and perilous adventure which has tested your wit and courage. You have mastered the first part of the ZORK trilogy. Those who pass over this bridge must be prepared to undertake an even greater adventure that will severely test your skill and bravery!\n\nThe ZORK trilogy continues with \"ZORK II: The Wizard of Frobozz\" and is completed in \"ZORK III: The Dungeon Master.\"", Newline)
			Finish()
			return true
		}
	}
	return false
}

func TrollRoomFcn(arg ActArg) bool {
	if arg == ActEnter && Troll.IsIn(G.Here) {
		ThisIsIt(&Troll)
	}
	return false
}

func ClearingFcn(arg ActArg) bool {
	if arg == ActEnter {
		if !G.GrateRevealed {
			Grate.Give(FlgInvis)
		}
		return false
	}
	if arg == ActLook {
		Print("You are in a clearing, with a forest surrounding you on all sides. A path leads south.", NoNewline)
		if Grate.Has(FlgOpen) {
			NewLine()
			Print("There is an open grating, descending into darkness.", NoNewline)
		} else if G.GrateRevealed {
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
		} else if G.GrUnlock {
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
		if G.CyclopsFlag && !G.MagicFlag {
			Print("The cyclops is sleeping blissfully at the foot of the stairs.", Newline)
		} else if G.MagicFlag {
			Print("The east wall, previously solid, now has a cyclops-sized opening in it.", Newline)
		} else if G.CycloWrath == 0 {
			Print("A cyclops, who looks prepared to eat horses (much less mere adventurers), blocks the staircase. From his state of health, and the bloodstains on the walls, you gather that he is not very friendly, though he likes people.", Newline)
		} else if G.CycloWrath > 0 {
			Print("The cyclops is standing in the corner, eyeing you closely. I don't think he likes you very much. He looks extremely hungry, even for a cyclops.", Newline)
		} else {
			Print("The cyclops, having eaten the hot peppers, appears to be gasping. His enflamed tongue protrudes from his man-sized mouth.", Newline)
		}
		return true
	}
	if arg == ActEnter {
		if G.CycloWrath == 0 {
			return false
		}
		Queue("ICyclops", -1).Run = true
		return false
	}
	return false
}

func TreasureRoomFcn(arg ActArg) bool {
	if arg == ActEnter && !G.Dead {
		if !Thief.IsIn(G.Here) {
			Print("You hear a scream of anguish as you violate the robber's hideaway. Using passages unknown to you, he rushes to its defense.", Newline)
			Thief.MoveTo(G.Here)
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
		if G.LowTide && G.GatesOpen {
			Print("You are in a long room, to the north of which was formerly a lake. However, with the water level lowered, there is merely a wide stream running through the center of the room.", NoNewline)
		} else if G.GatesOpen {
			Print("You are in a long room. To the north is a large lake, too deep to cross. You notice, however, that the water level appears to be dropping at a rapid rate. Before long, it might be possible to cross to the other side from here.", NoNewline)
		} else if G.LowTide {
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
	if arg == ActEnd && !G.Winner.Location().Has(FlgVeh) && !G.GatesOpen && G.LowTide {
		Print("You notice that the water level here is rising rapidly. The currents are also becoming stronger. Staying here seems quite perilous!", Newline)
		return true
	}
	if arg == ActLook {
		if G.LowTide {
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
		if G.LowTide && G.GatesOpen {
			Print("You are in a large cavernous room, the south of which was formerly a lake. However, with the water level lowered, there is merely a wide stream running through there.", NoNewline)
		} else if G.GatesOpen {
			Print("You are in a large cavernous area. To the south is a wide lake, whose water level appears to be falling rapidly.", NoNewline)
		} else if G.LowTide {
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
		if G.MirrorMung {
			Print("Unfortunately, the mirror has been destroyed by your recklessness.", Newline)
		}
		return true
	}
	return false
}

func Cave2RoomFcn(arg ActArg) bool {
	if arg == ActEnd {
		if Candles.IsIn(G.Winner) && Prob(50, true) && Candles.Has(FlgOn) {
			QueueInt("ICandles", false).Run = false
			Candles.Take(FlgOn)
			Print("A gust of wind blows out your candles!", Newline)
			G.Lit = IsLit(G.Here, true)
			if !G.Lit {
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
		if !G.LLDFlag && !G.Dead {
			Print("The way through the gate is barred by evil spirits, who jeer at your attempts to pass.", Newline)
		}
		return true
	}
	if arg == ActBegin {
		if G.ActVerb.Norm == "exorcise" && !G.LLDFlag {
			if Bell.IsIn(G.Winner) && Book.IsIn(G.Winner) && Candles.IsIn(G.Winner) {
				Print("You must perform the ceremony.", Newline)
			} else {
				Print("You aren't equipped for an exorcism.", Newline)
			}
			return true
		}
		if !G.LLDFlag && G.ActVerb.Norm == "ring" && G.DirObj == &Bell {
			G.XB = true
			RemoveCarefully(&Bell)
			ThisIsIt(&HotBell)
			HotBell.MoveTo(G.Here)
			Print("The bell suddenly becomes red hot and falls to the ground. The wraiths, as if paralyzed, stop their jeering and slowly turn to face you. On their ashen faces, the expression of a long-forgotten terror takes shape.", Newline)
			if Candles.IsIn(G.Winner) {
				Print("In your confusion, the candles drop to the ground (and they are out).", Newline)
				Candles.MoveTo(G.Here)
				Candles.Take(FlgOn)
				QueueInt("ICandles", false).Run = false
			}
			Queue("IXb", 6).Run = true
			Queue("IXbh", 20).Run = true
			return true
		}
		if G.XC && G.ActVerb.Norm == "read" && G.DirObj == &Book && !G.LLDFlag {
			Print("Each word of the prayer reverberates through the hall in a deafening confusion. As the last word fades, a voice, loud and commanding, speaks: \"Begone, fiends!\" A heart-stopping scream fills the cavern, and the spirits, sensing a greater power, flee through the walls.", Newline)
			RemoveCarefully(&Ghosts)
			G.LLDFlag = true
			QueueInt("IXc", false).Run = false
			return true
		}
	}
	if arg == ActEnd {
		if G.XB && Candles.IsIn(G.Winner) && Candles.Has(FlgOn) && !G.XC {
			G.XC = true
			Print("The flames flicker wildly and appear to dance. The earth beneath your feet trembles, and your legs nearly buckle beneath you. The spirits cower at your unearthly power.", Newline)
			QueueInt("IXb", false).Run = false
			Queue("IXc", 3).Run = true
		}
	}
	return false
}

func DomeRoomFcn(arg ActArg) bool {
	if arg == ActLook {
		Print("You are at the periphery of a large dome, which forms the ceiling of another room below. Protecting you from a precipitous drop is a wooden railing which circles the dome.", Newline)
		if G.DomeFlag {
			Print("Hanging down from the railing is a rope which ends about ten feet from the floor below.", Newline)
		}
		return true
	}
	if arg == ActEnter {
		if G.Dead {
			Print("As you enter the dome you feel a strong pull as if from a wind drawing you over the railing and down.", Newline)
			G.Winner.MoveTo(&TorchRoom)
			G.Here = &TorchRoom
			return true
		}
		if G.ActVerb.Norm == "leap" {
			JigsUp("I'm afraid that the leap you attempted has done you in.", false)
			return true
		}
	}
	return false
}

func TorchRoomFcn(arg ActArg) bool {
	if arg == ActLook {
		Print("This is a large room with a prominent doorway leading to a down staircase. Above you is a large dome. Up around the edge of the dome (20 feet up) is a wooden railing. In the center of the room sits a white marble pedestal.", Newline)
		if G.DomeFlag {
			Print("A piece of rope descends from the railing above, ending some five feet above your head.", Newline)
		}
		return true
	}
	return false
}

func SouthTempleFcn(arg ActArg) bool {
	if arg == ActBegin {
		G.CoffinCure = !Coffin.IsIn(G.Winner)
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
		if G.LoudFlag || (!G.GatesOpen && G.LowTide) {
			Print(" The room is eerie in its quietness.", NoNewline)
		} else {
			Print(" The room is deafeningly loud with an undetermined rushing sound. The sound seems to reverberate from all of the walls, making it difficult even to think.", NoNewline)
		}
		NewLine()
		return true
	}
	if arg == ActEnd && G.GatesOpen && !G.LowTide {
		Print("It is unbearably loud here, with an ear-splitting roar seeming to come from all around you. There is a pounding in your head which won't stop. With a tremendous effort, you scramble out of the room.", Newline)
		NewLine()
		dest := LoudRuns[rand.Intn(len(LoudRuns))]
		Goto(dest, true)
		return false
	}
	if arg == ActEnter {
		if G.LoudFlag || (!G.GatesOpen && G.LowTide) {
			return false
		}
		if G.GatesOpen && !G.LowTide {
			return false
		}
		// Room is loud - special input handling
		VFirstLook(ActUnk)
		if G.Params.Continue != NumUndef {
			Print("The rest of your commands have been lost in the noise.", Newline)
			G.Params.Continue = NumUndef
		}
		// In the original, this has a special read loop. We simplify.
		return false
	}
	if G.ActVerb.Norm == "echo" {
		if G.LoudFlag || (!G.GatesOpen && G.LowTide) {
			// Room is already quiet
			Print("echo echo ...", Newline)
			return true
		}
		Print("The acoustics of the room change subtly.", Newline)
		G.LoudFlag = true
		return true
	}
	return false
}

func DeepCanyonFcn(arg ActArg) bool {
	if arg == ActLook {
		Print("You are on the south edge of a deep canyon. Passages lead off to the east, northwest and southwest. A stairway leads down.", NoNewline)
		if G.GatesOpen && !G.LowTide {
			Print(" You can hear a loud roaring sound, like that of rushing water, from below.", NoNewline)
		} else if !G.GatesOpen && G.LowTide {
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
		if G.ActVerb.Norm == "lamp on" || G.ActVerb.Norm == "burn" {
			if G.DirObj == &Candles || G.DirObj == &Torch || G.DirObj == &Match {
				dummy = true
			}
		}
		if (Candles.IsIn(G.Winner) && Candles.Has(FlgOn)) ||
			(Torch.IsIn(G.Winner) && Torch.Has(FlgOn)) ||
			(Match.IsIn(G.Winner) && Match.Has(FlgOn)) {
			if dummy {
				Print("How sad for an aspiring adventurer to light a ", NoNewline)
				PrintObject(G.DirObj)
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
	if arg == ActEnter && !G.Dead {
		if Garlic.Location() != G.Winner && !Garlic.IsIn(G.Here) {
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
		f := G.Winner.Children
		G.EmptyHanded = true
		for _, child := range f {
			if Weight(child) > 4 {
				G.EmptyHanded = false
				break
			}
		}
		if G.Here == &LowerShaft && G.Lit {
			ScoreUpd(G.LightShaft)
			G.LightShaft = 0
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
		Queue("IForestRandom", -1).Run = true
		return false
	}
	if arg == ActBegin {
		if (G.ActVerb.Norm == "climb" || G.ActVerb.Norm == "climb up") && G.DirObj == &Tree {
			DoWalk(Up)
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
		if (G.ActVerb.Norm == "climb down") && (G.DirObj == &Tree || G.DirObj == &Rooms) {
			DoWalk(Down)
			return true
		}
		if (G.ActVerb.Norm == "climb up" || G.ActVerb.Norm == "climb") && G.DirObj == &Tree {
			DoWalk(Up)
			return true
		}
		if G.ActVerb.Norm == "drop" {
			if !IDrop() {
				return true
			}
			if G.DirObj == &Nest && Egg.IsIn(&Nest) {
				Print("The nest falls to the ground, and the egg spills out of it, seriously damaged.", Newline)
				RemoveCarefully(&Egg)
				BrokenEgg.MoveTo(&Path)
				return true
			}
			if G.DirObj == &Egg {
				Print("The egg falls to the ground and springs open, seriously damaged.", NoNewline)
				Egg.MoveTo(&Path)
				BadEgg()
				NewLine()
				return true
			}
			if G.DirObj != G.Winner && G.DirObj != &Tree {
				G.DirObj.MoveTo(&Path)
				Print("The ", NoNewline)
				PrintObject(G.DirObj)
				Print(" falls to the ground.", Newline)
			}
			return true
		}
	}
	if arg == ActEnter {
		Queue("IForestRandom", -1).Run = true
	}
	return false
}

func DeadFunction(arg ActArg) bool {
	if G.ActVerb.Norm == "walk" {
		if G.Here == &TimberRoom && G.Params.HasWalkDir && G.Params.WalkDir == West {
			Print("You cannot enter in your condition.", Newline)
			return true
		}
		return false
	}
	if G.ActVerb.Norm == "brief" || G.ActVerb.Norm == "verbose" || G.ActVerb.Norm == "super-brief" || G.ActVerb.Norm == "version" {
		return false
	}
	if G.ActVerb.Norm == "attack" || G.ActVerb.Norm == "mung" || G.ActVerb.Norm == "alarm" || G.ActVerb.Norm == "swing" {
		Print("All such attacks are vain in your condition.", Newline)
		return true
	}
	if G.ActVerb.Norm == "open" || G.ActVerb.Norm == "close" || G.ActVerb.Norm == "eat" || G.ActVerb.Norm == "drink" ||
		G.ActVerb.Norm == "inflate" || G.ActVerb.Norm == "deflate" || G.ActVerb.Norm == "turn" || G.ActVerb.Norm == "burn" ||
		G.ActVerb.Norm == "tie" || G.ActVerb.Norm == "untie" || G.ActVerb.Norm == "rub" {
		Print("Even such an action is beyond your capabilities.", Newline)
		return true
	}
	if G.ActVerb.Norm == "wait" {
		Print("Might as well. You've got an eternity.", Newline)
		return true
	}
	if G.ActVerb.Norm == "lamp on" {
		Print("You need no light to guide you.", Newline)
		return true
	}
	if G.ActVerb.Norm == "score" {
		Print("You're dead! How can you think of your score?", Newline)
		return true
	}
	if G.ActVerb.Norm == "take" || G.ActVerb.Norm == "rub" {
		Print("Your hand passes through its object.", Newline)
		return true
	}
	if G.ActVerb.Norm == "drop" || G.ActVerb.Norm == "throw" || G.ActVerb.Norm == "inventory" {
		Print("You have no possessions.", Newline)
		return true
	}
	if G.ActVerb.Norm == "diagnose" {
		Print("You are dead.", Newline)
		return true
	}
	if G.ActVerb.Norm == "look" {
		Print("The room looks strange and unearthly", NoNewline)
		if !G.Here.HasChildren() {
			Print(".", NoNewline)
		} else {
			Print(" and objects appear indistinct.", NoNewline)
		}
		NewLine()
		if !G.Here.Has(FlgOn) {
			Print("Although there is no light, the room seems dimly illuminated.", Newline)
		}
		NewLine()
		return false
	}
	if G.ActVerb.Norm == "pray" {
		if G.Here == &SouthTemple {
			Lamp.Take(FlgInvis)
			G.Winner.Action = nil
			G.AlwaysLit = false
			G.Dead = false
			if Troll.IsIn(&TrollRoom) {
				G.TrollFlag = false
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
	G.Params.Continue = NumUndef
	return true
}

// ================================================================
// INTERRUPT ROUTINES
// ================================================================

func ICandles() bool {
	Candles.Give(FlgTouch)
	if G.CandleTableIdx >= len(CandleTable) {
		return true
	}
	tick := CandleTable[G.CandleTableIdx].(int)
	Queue("ICandles", tick).Run = true
	LightInt(&Candles, G.CandleTableIdx, tick)
	if tick != 0 {
		G.CandleTableIdx += 2
	}
	return true
}

func ILantern() bool {
	if G.LampTableIdx >= len(LampTable) {
		return true
	}
	tick := LampTable[G.LampTableIdx].(int)
	Queue("ILantern", tick).Run = true
	LightInt(&Lamp, G.LampTableIdx, tick)
	if tick != 0 {
		G.LampTableIdx += 2
	}
	return true
}

// LightInt handles light source countdown warnings and expiry
func LightInt(obj *Object, tblIdx, tick int) {
	if tick == 0 {
		obj.Take(FlgOn)
		obj.Give(FlgRMung)
	}
	if IsHeld(obj) || obj.IsIn(G.Here) {
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
	s := G.Winner.Strength
	if s > 0 {
		s = 0
		G.Winner.Strength = s
	} else if s < 0 {
		s++
		G.Winner.Strength = s
	}
	if s < 0 {
		if G.LoadAllowed < G.LoadMax {
			G.LoadAllowed += 10
		}
		Queue("ICure", CureWait).Run = true
	} else {
		G.LoadAllowed = G.LoadMax
		QueueInt("ICure", false).Run = false
	}
	return false
}

func IMatch() bool {
	Print("The match has gone out.", Newline)
	Match.Take(FlgFlame)
	Match.Take(FlgOn)
	G.Lit = IsLit(G.Here, true)
	return true
}

func IXb() bool {
	if !G.XC {
		if G.Here == &EnteranceToHades {
			Print("The tension of this ceremony is broken, and the wraiths, amused but shaken at your clumsy attempt, resume their hideous jeering.", Newline)
		}
	}
	G.XB = false
	return true
}

func IXbh() bool {
	RemoveCarefully(&HotBell)
	Bell.MoveTo(&EnteranceToHades)
	if G.Here == &EnteranceToHades {
		Print("The bell appears to have cooled down.", Newline)
	}
	return true
}

func IXc() bool {
	G.XC = false
	IXb()
	return true
}

func ICyclops() bool {
	if G.CyclopsFlag || G.Dead {
		return true
	}
	if G.Here != &CyclopsRoom {
		QueueInt("ICyclops", false).Run = false
		return false
	}
	if AbsInt(G.CycloWrath) > 5 {
		QueueInt("ICyclops", false).Run = false
		JigsUp("The cyclops, tired of all of your games and trickery, grabs you firmly. As he licks his chops, he says \"Mmm. Just like Mom used to make 'em.\" It's nice to be appreciated.", false)
		return true
	}
	if G.CycloWrath < 0 {
		G.CycloWrath--
	} else {
		G.CycloWrath++
	}
	if !G.CyclopsFlag {
		idx := AbsInt(G.CycloWrath) - 2
		if idx >= 0 && idx < len(Cyclomad) {
			Print(Cyclomad[idx], Newline)
		}
	}
	return true
}

func IForestRandom() bool {
	if !ForestRoomQ() {
		QueueInt("IForestRandom", false).Run = false
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
	if G.GrateRevealed {
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
	if G.RugMoved {
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
	f := G.Winner.Children
	if len(f) == 0 {
		Print("Going up empty-handed is a bad idea.", Newline)
		return nil
	}
	// Check if player is carrying at most 1-2 items including the lamp
	count := 0
	for range f {
		count++
	}
	if count <= 2 && Lamp.IsIn(G.Winner) {
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
	if G.Here == &Maze2 {
		return &Maze4
	}
	if G.Here == &Maze7 {
		return &DeadEnd1
	}
	if G.Here == &Maze9 {
		return &Maze11
	}
	if G.Here == &Maze12 {
		return &Maze5
	}
	return nil
}

// ================================================================
// PSEUDO FUNCTIONS
// ================================================================

func ChasmPseudo(arg ActArg) bool {
	if G.ActVerb.Norm == "leap" || (G.ActVerb.Norm == "put" && G.DirObj == &Me) {
		Print("You look before leaping, and realize that you would never survive.", Newline)
		return true
	}
	if G.ActVerb.Norm == "cross" {
		Print("It's too far to jump, and there's no bridge.", Newline)
		return true
	}
	if (G.ActVerb.Norm == "put" || G.ActVerb.Norm == "throw off") && G.IndirObj == &PseudoObject {
		Print("The ", NoNewline)
		PrintObject(G.DirObj)
		Print(" drops out of sight into the chasm.", Newline)
		RemoveCarefully(G.DirObj)
		return true
	}
	return false
}

func LakePseudo(arg ActArg) bool {
	if G.LowTide {
		Print("There's not much lake left....", Newline)
		return true
	}
	if G.ActVerb.Norm == "cross" {
		Print("It's too wide to cross.", Newline)
		return true
	}
	if G.ActVerb.Norm == "through" {
		Print("You can't swim in this lake.", Newline)
		return true
	}
	return false
}

func StreamPseudo(arg ActArg) bool {
	if G.ActVerb.Norm == "swim" || G.ActVerb.Norm == "through" {
		Print("You can't swim in the stream.", Newline)
		return true
	}
	if G.ActVerb.Norm == "cross" {
		Print("The other side is a sheer rock cliff.", Newline)
		return true
	}
	return false
}

func DomePseudo(arg ActArg) bool {
	if G.ActVerb.Norm == "kiss" {
		Print("No.", Newline)
		return true
	}
	return false
}

func GatePseudo(arg ActArg) bool {
	if G.ActVerb.Norm == "through" {
		DoWalk(In)
		return true
	}
	Print("The gate is protected by an invisible force. It makes your teeth ache to touch it.", Newline)
	return true
}

func DoorPseudo(arg ActArg) bool {
	if G.ActVerb.Norm == "open" || G.ActVerb.Norm == "close" {
		Print("The door won't budge.", Newline)
		return true
	}
	if G.ActVerb.Norm == "through" {
		DoWalk(South)
		return true
	}
	return false
}

func PaintPseudo(arg ActArg) bool {
	if G.ActVerb.Norm == "mung" {
		Print("Some paint chips away, revealing more paint.", Newline)
		return true
	}
	return false
}

func GasPseudo(arg ActArg) bool {
	if G.ActVerb.Norm == "breathe" {
		Print("There is too much gas to blow away.", Newline)
		return true
	}
	if G.ActVerb.Norm == "smell" {
		Print("It smells like coal gas in here.", Newline)
		return true
	}
	return false
}

func ChainPseudo(arg ActArg) bool {
	if G.ActVerb.Norm == "take" || G.ActVerb.Norm == "move" {
		Print("The chain is secure.", Newline)
		return true
	}
	if G.ActVerb.Norm == "raise" || G.ActVerb.Norm == "lower" {
		Print("Perhaps you should do that to the basket.", Newline)
		return true
	}
	if G.ActVerb.Norm == "examine" {
		Print("The chain secures a basket within the shaft.", Newline)
		return true
	}
	return false
}

func BarrowDoorFcn2(arg ActArg) bool {
	return false
}
