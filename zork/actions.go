package zork

func OpenClose(obj *Object, strOpn, strCls string) {
	if G.ActVerb.Norm == "open" {
		if obj.Has(FlgOpen) {
			Printf("%s\n", PickOne(Dummy))
		} else {
			Printf("%s\n", strOpn)
			obj.Give(FlgOpen)
		}
	} else if G.ActVerb.Norm == "close" {
		if obj.Has(FlgOpen) {
			Printf("%s\n", strCls)
			obj.Take(FlgOpen)
		} else {
			Printf("%s\n", PickOne(Dummy))
		}
	}
}

func LeavesAppear() bool {
	if !Grate.Has(FlgOpen) && !G.GrateRevealed {
		if G.ActVerb.Norm == "move" || G.ActVerb.Norm == "take" {
			Printf("In disturbing the pile of leaves, a grating is revealed.\n")
		} else {
			Printf("With the leaves moved, a grating is revealed.\n")
		}
		Grate.Take(FlgInvis)
		G.GrateRevealed = true
	}
	return false
}

func Fweep(n int) {
	for i := 0; i < n; i++ {
		Printf("    Fweep!\n")
	}
	Printf("\n")
}

func FlyMe() bool {
	Fweep(4)
	Printf("The bat grabs you by the scruff of your neck and lifts you away....\n\n")
	dest := BatDrops[G.Rand.Intn(len(BatDrops))]
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
	Printf("It is an integral part of the control panel.\n")
}

func WithTell(obj *Object) {
	Printf("With a %s?\n", obj.Desc)
}


func BadEgg() {
	if Canary.IsIn(&Egg) {
		Printf(" %s", BrokenCanary.FirstDesc)
	} else {
		RemoveCarefully(&BrokenCanary)
	}
	BrokenEgg.MoveTo(Egg.Location())
	RemoveCarefully(&Egg)
}

func Slider(obj *Object) {
	if obj.Has(FlgTake) {
		Printf("The %s falls into the slide and is gone.\n", obj.Desc)
		if obj == &Water {
			RemoveCarefully(obj)
		} else {
			obj.MoveTo(&Cellar)
		}
	} else {
		Printf("%s\n", PickOne(Yuks))
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
	Printf("Your score is %d (total of 350 points), in %d", G.Score, G.Moves)
	if G.Moves == 1 {
		Printf(" move.")
	} else {
		Printf(" moves.")
	}
	Printf("\nThis gives you the rank of ")
	switch {
	case G.Score == 350:
		Printf("Master Adventurer")
	case G.Score > 330:
		Printf("Wizard")
	case G.Score > 300:
		Printf("Master")
	case G.Score > 200:
		Printf("Adventurer")
	case G.Score > 100:
		Printf("Junior Adventurer")
	case G.Score > 50:
		Printf("Novice Adventurer")
	case G.Score > 25:
		Printf("Amateur Adventurer")
	default:
		Printf("Beginner")
	}
	Printf(".\n")
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
	if G.Deaths > 0 {
		Printf("You have been killed ")
		if G.Deaths == 1 {
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

func JigsUp(desc string, isPlyr bool) bool {
	G.Winner = &Adventurer
	if G.Dead {
		Printf("\nIt takes a talented person to be killed while already dead. YOU are such a talent. Unfortunately, it takes a talented person to deal with it. I am not such a talent. Sorry.\n")
		return Finish()
	}
	Printf("%s\n", desc)
	if !G.Lucky {
		Printf("Bad luck, huh?\n")
	}
	ScoreUpd(-10)
	Printf("\n    ****  You have died  ****\n\n")
	if G.Winner.Location().Has(FlgVeh) {
		G.Winner.MoveTo(G.Here)
	}
	if G.Deaths >= 2 {
		Printf("You clearly are a suicidal maniac. We don't allow psychotics in the cave, since they may harm other adventurers. Your remains will be installed in the Land of the Living Dead, where your fellow adventurers may gloat over them.\n")
		return Finish()
	}
	G.Deaths++
	G.Winner.MoveTo(G.Here)
	if SouthTemple.Has(FlgTouch) {
		Printf("As you take your last breath, you feel relieved of your burdens. The feeling passes as you find yourself before the gates of Hell, where the spirits jeer at you and deny you entry. Your senses are disturbed. The objects in the dungeon appear indistinct, bleached of color, even unreal.\n\n")
		G.Dead = true
		G.TrollFlag = true
		G.AlwaysLit = true
		G.Winner.Action = DeadFunction
		Goto(&EnteranceToHades, true)
	} else {
		Printf("Now, let's take a look here...\nWell, you probably deserve another chance. I can't quite fix you up completely, but you can't have everything.\n\n")
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
		Printf("You are standing in an open field west of a white house, with a boarded front door.")
		if G.WonGame {
			Printf(" A secret path leads southwest into the forest.")
		}
		Printf("\n")
		return true
	}
	return false
}

func EastHouseFcn(arg ActArg) bool {
	if arg == ActLook {
		Printf("You are behind the white house. A path leads into the forest to the east. In one corner of the house there is a small window which is ")
		if KitchenWindow.Has(FlgOpen) {
			Printf("open.\n")
		} else {
			Printf("slightly ajar.\n")
		}
		return true
	}
	return false
}

func WhiteHouseFcn(arg ActArg) bool {
	if G.Here == &Kitchen || G.Here == &LivingRoom || G.Here == &Attic {
		if G.ActVerb.Norm == "find" {
			Printf("Why not find your brains?\n")
			return true
		}
		if G.ActVerb.Norm == "walk around" {
			GoNext(InHouseAround)
			return true
		}
	} else if G.Here != &WestOfHouse && G.Here != &NorthOfHouse && G.Here != &EastOfHouse && G.Here != &SouthOfHouse {
		if G.ActVerb.Norm == "find" {
			if G.Here == &Clearing {
				Printf("It seems to be to the west.\n")
				return true
			}
			Printf("It was here just a minute ago....\n")
			return true
		}
		Printf("You're not at the house.\n")
		return true
	} else if G.ActVerb.Norm == "find" {
		Printf("It's right here! Are you blind or something?\n")
		return true
	} else if G.ActVerb.Norm == "walk around" {
		GoNext(HouseAround)
		return true
	} else if G.ActVerb.Norm == "examine" {
		Printf("The house is a beautiful colonial house which is painted white. It is clear that the owners must have been extremely wealthy.\n")
		return true
	} else if G.ActVerb.Norm == "through" || G.ActVerb.Norm == "open" {
		if G.Here == &EastOfHouse {
			if KitchenWindow.Has(FlgOpen) {
				return Goto(&Kitchen, true)
			}
			Printf("The window is closed.\n")
			ThisIsIt(&KitchenWindow)
			return true
		}
		Printf("I can't see how to get in from here.\n")
		return true
	} else if G.ActVerb.Norm == "burn" {
		Printf("You must be joking.\n")
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
		Printf("The boards are securely fastened.\n")
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
			Printf("Dental hygiene is highly recommended, but I'm not sure what you want to brush them with.\n")
			return true
		}
		Printf("A nice idea, but with a %s?\n", G.IndirObj.Desc)
		return true
	}
	return false
}

func GraniteWallFcn(arg ActArg) bool {
	if G.Here == &NorthTemple {
		if G.ActVerb.Norm == "find" {
			Printf("The west wall is solid granite here.\n")
			return true
		}
		if G.ActVerb.Norm == "take" || G.ActVerb.Norm == "raise" || G.ActVerb.Norm == "lower" {
			Printf("It's solid granite.\n")
			return true
		}
	} else if G.Here == &TreasureRoom {
		if G.ActVerb.Norm == "find" {
			Printf("The east wall is solid granite here.\n")
			return true
		}
		if G.ActVerb.Norm == "take" || G.ActVerb.Norm == "raise" || G.ActVerb.Norm == "lower" {
			Printf("It's solid granite.\n")
			return true
		}
	} else if G.Here == &SlideRoom {
		if G.ActVerb.Norm == "find" || G.ActVerb.Norm == "read" {
			Printf("It only SAYS \"Granite Wall\".\n")
			return true
		}
		Printf("The wall isn't granite.\n")
		return true
	} else {
		Printf("There is no granite wall here.\n")
		return true
	}
	return false
}

func SongbirdFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "find" || G.ActVerb.Norm == "take" {
		Printf("The songbird is not here but is probably nearby.\n")
		return true
	}
	if G.ActVerb.Norm == "listen" {
		Printf("You can't hear the songbird now.\n")
		return true
	}
	if G.ActVerb.Norm == "follow" {
		Printf("It can't be followed.\n")
		return true
	}
	Printf("You can't see any songbird here.\n")
	return true
}

func MountainRangeFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "climb up" || G.ActVerb.Norm == "climb down" || G.ActVerb.Norm == "climb" {
		Printf("Don't you believe me? The mountains are impassable!\n")
		return true
	}
	return false
}

func ForestFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "walk around" {
		if G.Here == &WestOfHouse || G.Here == &NorthOfHouse || G.Here == &SouthOfHouse || G.Here == &EastOfHouse {
			Printf("You aren't even in the forest.\n")
			return true
		}
		GoNext(ForestAround)
		return true
	}
	if G.ActVerb.Norm == "disembark" {
		Printf("You will have to specify a direction.\n")
		return true
	}
	if G.ActVerb.Norm == "find" {
		Printf("You cannot see the forest for the trees.\n")
		return true
	}
	if G.ActVerb.Norm == "listen" {
		Printf("The pines and the hemlocks seem to be murmuring.\n")
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
		Printf("The window is slightly ajar, but not enough to allow entry.\n")
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
		Printf("You can see ")
		if G.Here == &Kitchen {
			Printf("a clear area leading towards a forest.\n")
		} else {
			Printf("what appears to be a kitchen.\n")
		}
		return true
	}
	return false
}

func ChimneyFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "examine" {
		Printf("The chimney leads ")
		if G.Here == &Kitchen {
			Printf("down")
		} else {
			Printf("up")
		}
		Printf("ward, and looks climbable.\n")
		return true
	}
	return false
}

func GhostsFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "tell" {
		Printf("The spirits jeer loudly and ignore you.\n")
		G.Params.Continue = NumUndef
		return true
	}
	if G.ActVerb.Norm == "exorcise" {
		Printf("Only the ceremony itself has any effect.\n")
		return true
	}
	if (G.ActVerb.Norm == "attack" || G.ActVerb.Norm == "mung") && G.DirObj == &Ghosts {
		Printf("How can you attack a spirit with material objects?\n")
		return true
	}
	Printf("You seem unable to interact with these spirits.\n")
	return true
}

func BasketFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "raise" {
		if G.CageTop {
			Printf("%s\n", PickOne(Dummy))
		} else {
			RaisedBasket.MoveTo(&ShaftRoom)
			LoweredBasket.MoveTo(&LowerShaft)
			G.CageTop = true
			ThisIsIt(&RaisedBasket)
			Printf("The basket is raised to the top of the shaft.\n")
		}
		return true
	}
	if G.ActVerb.Norm == "lower" {
		if !G.CageTop {
			Printf("%s\n", PickOne(Dummy))
		} else {
			RaisedBasket.MoveTo(&LowerShaft)
			LoweredBasket.MoveTo(&ShaftRoom)
			ThisIsIt(&LoweredBasket)
			Printf("The basket is lowered to the bottom of the shaft.\n")
			G.CageTop = false
			if G.Lit && !IsLit(G.Here, true) {
				G.Lit = false
				Printf("It is now pitch black.\n")
			}
		}
		return true
	}
	if G.DirObj == &LoweredBasket || G.IndirObj == &LoweredBasket {
		Printf("The basket is at the other end of the chain.\n")
		return true
	}
	if G.ActVerb.Norm == "take" && (G.DirObj == &RaisedBasket || G.DirObj == &LoweredBasket) {
		Printf("The cage is securely fastened to the iron chain.\n")
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
			Printf("You can't reach him; he's on the ceiling.\n")
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
		Printf("Ding, dong.\n")
		return true
	}
	return false
}

func HotBellFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "take" {
		Printf("The bell is very hot and cannot be taken.\n")
		return true
	}
	if G.ActVerb.Norm == "rub" || (G.ActVerb.Norm == "ring" && G.IndirObj != nil) {
		if G.IndirObj != nil && G.IndirObj.Has(FlgBurn) {
			Printf("The %s burns and is consumed.\n", G.IndirObj.Desc)
			RemoveCarefully(G.IndirObj)
			return true
		}
		if G.IndirObj == &Hands {
			Printf("The bell is too hot to touch.\n")
			return true
		}
		Printf("The heat from the bell is too intense.\n")
		return true
	}
	if G.ActVerb.Norm == "pour on" {
		RemoveCarefully(G.DirObj)
		Printf("The water cools the bell and is evaporated.\n")
		QueueInt("IXbh", false).Run = false
		IXbh()
		return true
	}
	if G.ActVerb.Norm == "ring" {
		Printf("The bell is too hot to reach.\n")
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
			Printf("You see a rickety staircase descending into darkness.\n")
		} else {
			Printf("It's closed.\n")
		}
		return true
	}
	if G.Here == &Cellar {
		if (G.ActVerb.Norm == "open" || G.ActVerb.Norm == "unlock") && !TrapDoor.Has(FlgOpen) {
			Printf("The door is locked from above.\n")
			return true
		}
		if G.ActVerb.Norm == "close" && !TrapDoor.Has(FlgOpen) {
			TrapDoor.Take(FlgTouch)
			TrapDoor.Take(FlgOpen)
			Printf("The door closes and locks.\n")
			return true
		}
		if G.ActVerb.Norm == "open" || G.ActVerb.Norm == "close" {
			Printf("%s\n", PickOne(Dummy))
			return true
		}
	}
	return false
}

func FrontDoorFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "open" {
		Printf("The door cannot be opened.\n")
		return true
	}
	if G.ActVerb.Norm == "burn" {
		Printf("You cannot burn this door.\n")
		return true
	}
	if G.ActVerb.Norm == "mung" {
		Printf("You can't seem to damage the door.\n")
		return true
	}
	if G.ActVerb.Norm == "look behind" {
		Printf("It won't open.\n")
		return true
	}
	return false
}

func BarrowDoorFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "open" || G.ActVerb.Norm == "close" {
		Printf("The door is too heavy.\n")
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
		Printf("The bottle hits the far wall and shatters.\n")
	} else if G.ActVerb.Norm == "mung" {
		empty = true
		RemoveCarefully(G.DirObj)
		Printf("A brilliant maneuver destroys the bottle.\n")
	} else if G.ActVerb.Norm == "shake" {
		if Bottle.Has(FlgOpen) && Water.IsIn(&Bottle) {
			empty = true
		}
	}
	if empty && Water.IsIn(&Bottle) {
		Printf("The water spills to the floor and evaporates.\n")
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
		Printf("You can't fit through the crack.\n")
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
			Printf("The grate is locked.\n")
			return true
		}
		if G.Here == &Clearing {
			Printf("You can't lock it from this side.\n")
			return true
		}
	}
	if G.ActVerb.Norm == "unlock" && G.DirObj == &Grate {
		if G.Here == &GratingRoom && G.IndirObj == &Keys {
			G.GrUnlock = true
			Printf("The grate is unlocked.\n")
			return true
		}
		if G.Here == &Clearing && G.IndirObj == &Keys {
			Printf("You can't reach the lock from here.\n")
			return true
		}
		Printf("Can you unlock a grating with a %s?\n", G.IndirObj.Desc)
		return true
	}
	if G.ActVerb.Norm == "pick" {
		Printf("You can't pick the lock.\n")
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
					Printf("A pile of leaves falls onto your head and to the ground.\n")
					G.GrateRevealed = true
					Leaves.MoveTo(G.Here)
				}
				GratingRoom.Give(FlgOn)
			} else {
				GratingRoom.Take(FlgOn)
			}
		} else {
			Printf("The grating is locked.\n")
		}
		return true
	}
	if G.ActVerb.Norm == "put" && G.IndirObj == &Grate {
		if G.DirObj.Size > 20 {
			Printf("It won't fit through the grating.\n")
		} else {
			G.DirObj.MoveTo(&GratingRoom)
			Printf("The %s goes through the grating into the darkness below.\n", G.DirObj.Desc)
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
		Printf("A ghost appears in the room and is appalled at your desecration of the remains of a fellow adventurer. He casts a curse on your valuables and banishes them to the Land of the Living Dead. The ghost leaves, muttering obscenities.\n")
		Rob(G.Here, &LandOfLivingDead, 100)
		Rob(&Adventurer, &LandOfLivingDead, 0)
		return true
	}
	return false
}

func TorchFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "examine" {
		Printf("The torch is burning.\n")
		return true
	}
	if G.ActVerb.Norm == "pour on" && G.IndirObj == &Torch {
		Printf("The water evaporates before it gets close.\n")
		return true
	}
	if G.ActVerb.Norm == "lamp off" && G.DirObj.Has(FlgOn) {
		Printf("You nearly burn your hand trying to extinguish the flame.\n")
		return true
	}
	return false
}

func RustyKnifeFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "take" {
		if Sword.IsIn(G.Winner) {
			Printf("As you touch the rusty knife, your sword gives a single pulse of blinding blue light.\n")
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
		Printf("There are 69,105 leaves here.\n")
		return true
	}
	if G.ActVerb.Norm == "burn" {
		LeavesAppear()
		RemoveCarefully(G.DirObj)
		if G.DirObj.IsIn(G.Here) {
			Printf("The leaves burn.\n")
		} else {
			JigsUp("The leaves burn, and so do you.", false)
		}
		return true
	}
	if G.ActVerb.Norm == "cut" {
		Printf("You rustle the leaves around, making quite a mess.\n")
		LeavesAppear()
		return true
	}
	if G.ActVerb.Norm == "move" || G.ActVerb.Norm == "take" {
		if G.ActVerb.Norm == "move" {
			Printf("Done.\n")
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
		Printf("Underneath the pile of leaves is a grating. As you release the leaves, the grating is once again concealed from view.\n")
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
			Printf("I'm afraid that you have run out of matches.\n")
			return true
		}
		if G.Here == &LowerShaft || G.Here == &TimberRoom {
			Printf("This room is drafty, and the match goes out instantly.\n")
			return true
		}
		Match.Give(FlgFlame)
		Match.Give(FlgOn)
		Queue("IMatch", 2).Run = true
		Printf("One of the matches starts to burn.\n")
		if !G.Lit {
			G.Lit = true
			VLook(ActUnk)
		}
		return true
	}
	if G.ActVerb.Norm == "lamp off" && Match.Has(FlgFlame) {
		Printf("The match is out.\n")
		Match.Take(FlgFlame)
		Match.Take(FlgOn)
		G.Lit = IsLit(G.Here, true)
		if !G.Lit {
			Printf("It's pitch black in here!\n")
		}
		QueueInt("IMatch", false).Run = false
		return true
	}
	if G.ActVerb.Norm == "count" || G.ActVerb.Norm == "open" {
		Printf("You have ")
		cnt := G.MatchCount - 1
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
	}
	if G.ActVerb.Norm == "examine" {
		if Match.Has(FlgOn) {
			Printf("The match is burning.\n")
		} else {
			Printf("The matchbook isn't very interesting, except for what's written on it.\n")
		}
		return true
	}
	return false
}

func MirrorMirrorFcn(arg ActArg) bool {
	rm2 := &MirrorRoom2
	if !G.MirrorMung && G.ActVerb.Norm == "rub" {
		if G.IndirObj != nil && G.IndirObj != &Hands {
			Printf("You feel a faint tingling transmitted through the %s.\n", G.IndirObj.Desc)
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
		Printf("There is a rumble from deep within the earth and the room shakes.\n")
		return true
	}
	if G.ActVerb.Norm == "look inside" || G.ActVerb.Norm == "examine" {
		if G.MirrorMung {
			Printf("The mirror is broken into many pieces.\n")
		} else {
			Printf("There is an ugly person staring back at you.\n")
		}
		return true
	}
	if G.ActVerb.Norm == "take" {
		Printf("The mirror is many times your size. Give up.\n")
		return true
	}
	if G.ActVerb.Norm == "mung" || G.ActVerb.Norm == "throw" || G.ActVerb.Norm == "attack" {
		if G.MirrorMung {
			Printf("Haven't you done enough damage already?\n")
		} else {
			G.MirrorMung = true
			G.Lucky = false
			Printf("You have broken the mirror. I hope you have a seven years' supply of good luck handy.\n")
		}
		return true
	}
	return false
}

func PaintingFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "mung" {
		G.DirObj.TValue = 0
		G.DirObj.LongDesc = "There is a worthless piece of canvas here."
		Printf("Congratulations! Unlike the other vandals, who merely stole the artist's masterpieces, you have destroyed one.\n")
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
			Printf("Alas, there's not much left of the candles. Certainly not enough to burn.\n")
			return true
		}
		if G.IndirObj == nil {
			if Match.Has(FlgFlame) {
				Printf("(with the match)\n")
				Perform(ActionVerb{Norm: "lamp on", Orig: "light"}, &Candles, &Match)
				return true
			}
			Printf("You should say what to light them with.\n")
			return true
		}
		if G.IndirObj == &Match && Match.Has(FlgOn) {
			Printf("The candles are ")
			if Candles.Has(FlgOn) {
				Printf("already lit.\n")
			} else {
				Candles.Give(FlgOn)
				Printf("lit.\n")
				Queue("ICandles", -1).Run = true
			}
			return true
		}
		if G.IndirObj == &Torch {
			if Candles.Has(FlgOn) {
				Printf("You realize, just in time, that the candles are already lighted.\n")
			} else {
				Printf("The heat from the torch is so intense that the candles are vaporized.\n")
				RemoveCarefully(&Candles)
			}
			return true
		}
		Printf("You have to light them with something that's burning, you know.\n")
		return true
	}
	if G.ActVerb.Norm == "count" {
		Printf("Let's see, how many objects in a pair? Don't tell me, I'll get it.\n")
		return true
	}
	if G.ActVerb.Norm == "lamp off" {
		QueueInt("ICandles", false).Run = false
		if Candles.Has(FlgOn) {
			Printf("The flame is extinguished.")
			Candles.Take(FlgOn)
			Candles.Give(FlgTouch)
			G.Lit = IsLit(G.Here, true)
			if !G.Lit {
				Printf(" It's really dark in here....")
			}
			Printf("\n")
			return true
		}
		Printf("The candles are not lighted.\n")
		return true
	}
	if G.ActVerb.Norm == "put" && G.IndirObj != nil && G.IndirObj.Has(FlgBurn) {
		Printf("That wouldn't be smart.\n")
		return true
	}
	if G.ActVerb.Norm == "examine" {
		Printf("The candles are ")
		if Candles.Has(FlgOn) {
			Printf("burning.\n")
		} else {
			Printf("out.\n")
		}
		return true
	}
	return false
}

func GunkFcn(arg ActArg) bool {
	RemoveCarefully(&Gunk)
	Printf("The slag was rather insubstantial, and crumbles into dust at your touch.\n")
	return true
}

func BodyFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "take" {
		Printf("A force keeps you from taking the bodies.\n")
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
		Printf("The book is already open to page 569.\n")
		return true
	}
	if G.ActVerb.Norm == "close" {
		Printf("As hard as you try, the book cannot be closed.\n")
		return true
	}
	if G.ActVerb.Norm == "turn" {
		Printf("Beside page 569, there is only one other page with any legible printing on it. Most of it is unreadable, but the subject seems to be the banishment of evil. Apparently, certain noises, lights, and prayers are efficacious in this regard.\n")
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
				Printf("Suddenly, the rainbow appears to become solid and, I venture, walkable (I think the giveaway was the stairs and bannister).\n")
				if G.Here == &EndOfRainbow && PotOfGold.IsIn(&EndOfRainbow) {
					Printf("A shimmering pot of gold appears at the end of the rainbow.\n")
				}
				G.RainbowFlag = true
			} else {
				Rob(&OnRainbow, &Wall, 0)
				Printf("The rainbow seems to have become somewhat run-of-the-mill.\n")
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
		Printf("A dazzling display of color briefly emanates from the sceptre.\n")
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
		Printf("You tumble down the slide....\n")
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
		Printf("It smells of hot peppers.\n")
		return true
	}
	return false
}

func ToolChestFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "examine" {
		Printf("The chests are all empty.\n")
		return true
	}
	if G.ActVerb.Norm == "take" || G.ActVerb.Norm == "open" || G.ActVerb.Norm == "put" {
		RemoveCarefully(&ToolChest)
		Printf("The chests are so rusty and corroded that they crumble when you touch them.\n")
		return true
	}
	return false
}

func ButtonFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "read" {
		Printf("They're greek to you.\n")
		return true
	}
	if G.ActVerb.Norm == "push" {
		if G.DirObj == &BlueButton {
			if G.WaterLevel == 0 {
				Leak.Take(FlgInvis)
				Printf("There is a rumbling sound and a stream of water appears to burst from the east wall of the room (apparently, a leak has occurred in a pipe).\n")
				G.WaterLevel = 1
				Queue("IMaintRoom", -1).Run = true
				return true
			}
			Printf("The blue button appears to be jammed.\n")
			return true
		}
		if G.DirObj == &RedButton {
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
		if G.DirObj == &BrownButton {
			DamRoom.Take(FlgTouch)
			G.GateFlag = false
			Printf("Click.\n")
			return true
		}
		if G.DirObj == &YellowButton {
			DamRoom.Take(FlgTouch)
			G.GateFlag = true
			Printf("Click.\n")
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
		Printf("It is far too large to carry.\n")
		return true
	}
	if G.ActVerb.Norm == "open" {
		if Machine.Has(FlgOpen) {
			Printf("%s\n", PickOne(Dummy))
		} else if Machine.HasChildren() {
			Printf("The lid opens, revealing ")
			PrintContents(&Machine)
			Printf(".\n")
			Machine.Give(FlgOpen)
		} else {
			Printf("The lid opens.\n")
			Machine.Give(FlgOpen)
		}
		return true
	}
	if G.ActVerb.Norm == "close" {
		if Machine.Has(FlgOpen) {
			Printf("The lid closes.\n")
			Machine.Take(FlgOpen)
		} else {
			Printf("%s\n", PickOne(Dummy))
		}
		return true
	}
	if G.ActVerb.Norm == "lamp on" {
		if G.IndirObj == nil {
			Printf("It's not clear how to turn it on with your bare hands.\n")
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
				Printf("The machine doesn't seem to want to do anything.\n")
			} else {
				Printf("The machine comes to life (figuratively) with a dazzling display of colored lights and bizarre noises. After a few moments, the excitement abates.\n")
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
			Printf("It seems that a %s won't do.\n", G.IndirObj.Desc)
		}
		return true
	}
	return false
}

func PuttyFcn(arg ActArg) bool {
	if (G.ActVerb.Norm == "oil" && G.IndirObj == &Putty) || (G.ActVerb.Norm == "put" && G.DirObj == &Putty) {
		Printf("The all-purpose gunk isn't a lubricant.\n")
		return true
	}
	return false
}

func TubeFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "put" && G.IndirObj == &Tube {
		Printf("The tube refuses to accept anything.\n")
		return true
	}
	if G.ActVerb.Norm == "squeeze" {
		if G.DirObj.Has(FlgOpen) && Putty.IsIn(G.DirObj) {
			Putty.MoveTo(G.Winner)
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

func SwordFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "take" && G.Winner == &Adventurer {
		Queue("ISword", -1).Run = true
		return false
	}
	if G.ActVerb.Norm == "examine" {
		g := Sword.TValue
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

func LanternFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "throw" {
		Printf("The lamp has smashed into the floor, and the light has gone out.\n")
		QueueInt("ILantern", false).Run = false
		RemoveCarefully(&Lamp)
		BrokenLamp.MoveTo(G.Here)
		return true
	}
	if G.ActVerb.Norm == "lamp on" {
		if Lamp.Has(FlgRMung) {
			Printf("A burned-out lamp won't light.\n")
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
			Printf("The lamp has already burned out.\n")
			return true
		}
		QueueInt("ILantern", false).Run = false
		return false
	}
	if G.ActVerb.Norm == "examine" {
		Printf("The lamp ")
		if Lamp.Has(FlgRMung) {
			Printf("has burned out.\n")
		} else if Lamp.Has(FlgOn) {
			Printf("is on.\n")
		} else {
			Printf("is turned off.\n")
		}
		return true
	}
	return false
}

func MailboxFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "take" && G.DirObj == &Mailbox {
		Printf("It is securely anchored.\n")
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
		Printf("The troll isn't much of a conversationalist.\n")
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
				Printf("The troll, angered and humiliated, recovers his weapon. He appears to have an axe to grind with you.\n")
			}
			return true
		}
		if Troll.IsIn(G.Here) {
			Troll.LongDesc = "A pathetically babbling troll is here."
			Printf("The troll, disarmed, cowers in terror, pleading for his life in the guttural tongue of the trolls.\n")
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
			Printf("The troll stirs, quickly resuming a fighting stance.\n")
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
		Printf("%s\n", Troll.LongDesc)
		return true
	}
	if (G.ActVerb.Norm == "throw" || G.ActVerb.Norm == "give") && G.DirObj != nil && G.IndirObj == &Troll {
		Awaken(&Troll)
		if G.ActVerb.Norm == "throw" || G.ActVerb.Norm == "give" {
			if G.DirObj == &Axe && Axe.IsIn(G.Winner) {
				Printf("The troll scratches his head in confusion, then takes the axe.\n")
				Troll.Give(FlgFight)
				Axe.MoveTo(&Troll)
				return true
			}
			if G.DirObj == &Troll || G.DirObj == &Axe {
				Printf("You would have to get the %s first, and that seems unlikely.\n", G.DirObj.Desc)
				return true
			}
			if G.ActVerb.Norm == "throw" {
				Printf("The troll, who is remarkably coordinated, catches the %s", G.DirObj.Desc)
			} else {
				Printf("The troll, who is not overly proud, graciously accepts the gift")
			}
			if Prob(20, false) && (G.DirObj == &Knife || G.DirObj == &Sword || G.DirObj == &Axe) {
				RemoveCarefully(G.DirObj)
				Printf(" and eats it hungrily. Poor troll, he dies from an internal hemorrhage and his carcass disappears in a sinister black fog.\n")
				RemoveCarefully(&Troll)
				TrollFcn(ActArg(FDead))
				G.TrollFlag = true
			} else if G.DirObj == &Knife || G.DirObj == &Sword || G.DirObj == &Axe {
				G.DirObj.MoveTo(G.Here)
				Printf(" and, being for the moment sated, throws it back. Fortunately, the troll has poor control, and the %s falls to the floor. He does not look pleased.\n", G.DirObj.Desc)
				Troll.Give(FlgFight)
			} else {
				Printf(" and not having the most discriminating tastes, gleefully eats it.\n")
				RemoveCarefully(G.DirObj)
			}
			return true
		}
	}
	if G.ActVerb.Norm == "take" || G.ActVerb.Norm == "move" {
		Awaken(&Troll)
		if G.ActVerb.Norm == "take" || G.ActVerb.Norm == "move" {
			Printf("The troll spits in your face, grunting \"Better luck next time\" in a rather barbarous accent.\n")
			return true
		}
	}
	if G.ActVerb.Norm == "mung" {
		Awaken(&Troll)
		Printf("The troll laughs at your puny gesture.\n")
		return true
	}
	if G.ActVerb.Norm == "listen" {
		Printf("Every so often the troll says something, probably uncomplimentary, in his guttural tongue.\n")
		return true
	}
	if G.TrollFlag && G.ActVerb.Norm == "hello" {
		Printf("Unfortunately, the troll can't hear you.\n")
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
			Printf("No use talking to him. He's fast asleep.\n")
			return true
		}
		if G.ActVerb.Norm == "odysseus" {
			G.Winner = &Adventurer
			Perform(ActionVerb{Norm: "odysseus", Orig: "odysseus"}, nil, nil)
			return true
		}
		Printf("The cyclops prefers eating to making conversation.\n")
		return true
	}
	if G.CyclopsFlag {
		if G.ActVerb.Norm == "examine" {
			Printf("The cyclops is sleeping like a baby, albeit a very ugly one.\n")
			return true
		}
		if G.ActVerb.Norm == "alarm" || G.ActVerb.Norm == "kick" || G.ActVerb.Norm == "attack" || G.ActVerb.Norm == "burn" || G.ActVerb.Norm == "mung" {
			Printf("The cyclops yawns and stares at the thing that woke him up.\n")
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
		Printf("A hungry cyclops is standing at the foot of the stairs.\n")
		return true
	}
	if G.ActVerb.Norm == "give" && G.IndirObj == &Cyclops {
		if G.DirObj == &Lunch {
			if count >= 0 {
				RemoveCarefully(&Lunch)
				Printf("The cyclops says \"Mmm Mmm. I love hot peppers! But oh, could I use a drink. Perhaps I could drink the blood of that thing.\"  From the gleam in his eye, it could be surmised that you are \"that thing\".\n")
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
				Printf("The cyclops takes the bottle, checks that it's open, and drinks the water. A moment later, he lets out a yawn that nearly blows you over, and then falls fast asleep (what did you put in that drink, anyway?).\n")
				G.CyclopsFlag = true
			} else {
				Printf("The cyclops apparently is not thirsty and refuses your generous offer.\n")
			}
			return true
		}
		if G.DirObj == &Garlic {
			Printf("The cyclops may be hungry, but there is a limit.\n")
			return true
		}
		Printf("The cyclops is not so stupid as to eat THAT!\n")
		return true
	}
	if G.ActVerb.Norm == "throw" || G.ActVerb.Norm == "attack" || G.ActVerb.Norm == "mung" {
		Queue("ICyclops", -1).Run = true
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
	}
	if G.ActVerb.Norm == "take" {
		Printf("The cyclops doesn't take kindly to being grabbed.\n")
		return true
	}
	if G.ActVerb.Norm == "tie" {
		Printf("You cannot tie the cyclops, though he is fit to be tied.\n")
		return true
	}
	if G.ActVerb.Norm == "listen" {
		Printf("You can hear his stomach rumbling.\n")
		return true
	}
	return false
}

// ================================================================
// THIEF / ROBBER
// ================================================================

func DumbContainerFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "open" || G.ActVerb.Norm == "close" || G.ActVerb.Norm == "look inside" {
		Printf("You can't do that.\n")
		return true
	}
	if G.ActVerb.Norm == "examine" {
		Printf("It looks pretty much like a %s.\n", G.DirObj.Desc)
		return true
	}
	return false
}

func ChaliceFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "take" {
		if G.DirObj.IsIn(&TreasureRoom) && Thief.IsIn(&TreasureRoom) && Thief.Has(FlgFight) && !Thief.Has(FlgInvis) && Thief.LongDesc != RobberUDesc {
			Printf("You'd be stabbed in the back first.\n")
			return true
		}
		return false
	}
	if G.ActVerb.Norm == "put" && G.IndirObj == &Chalice {
		Printf("You can't. It's not a very good chalice, is it?\n")
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
		Printf("The %s are safely inside; there's no need to do that.\n", str)
		return true
	}
	if G.ActVerb.Norm == "look inside" || G.ActVerb.Norm == "examine" {
		Printf("There are lots of %s in there.\n", str)
		return true
	}
	if G.ActVerb.Norm == "put" && G.IndirObj == obj {
		Printf("Don't be silly. It wouldn't be a %s anymore.\n", obj.Desc)
		return true
	}
	return false
}

func GarlicFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "eat" {
		RemoveCarefully(G.DirObj)
		Printf("What the heck! You won't make friends this way, but nobody around here is too friendly anyhow. Gulp!\n")
		return true
	}
	return false
}

func BatDescFcn(arg ActArg) bool {
	if Garlic.Location() == G.Winner || Garlic.IsIn(G.Here) {
		Printf("In the corner of the room on the ceiling is a large vampire bat who is obviously deranged and holding his nose.\n")
	} else {
		Printf("A large vampire bat, hanging from the ceiling, swoops down at you!\n")
	}
	return true
}

func TrophyCaseFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "take" && G.DirObj == &TrophyCase {
		Printf("The trophy case is securely fastened to the wall.\n")
		return true
	}
	return false
}

func BoardedWindowFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "open" {
		Printf("The windows are boarded and can't be opened.\n")
		return true
	}
	if G.ActVerb.Norm == "mung" {
		Printf("You can't break the windows open.\n")
		return true
	}
	return false
}

func NailsPseudo(arg ActArg) bool {
	if G.ActVerb.Norm == "take" {
		Printf("The nails, deeply imbedded in the door, cannot be removed.\n")
		return true
	}
	return false
}

func CliffObjectFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "leap" || (G.ActVerb.Norm == "put" && G.DirObj == &Me) {
		Printf("That would be very unwise. Perhaps even fatal.\n")
		return true
	}
	if G.IndirObj == &ClimbableCliff {
		if G.ActVerb.Norm == "put" || G.ActVerb.Norm == "throw off" {
			Printf("The %s tumbles into the river and is seen no more.\n", G.DirObj.Desc)
			RemoveCarefully(G.DirObj)
			return true
		}
	}
	return false
}

func WhiteCliffFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "climb up" || G.ActVerb.Norm == "climb down" || G.ActVerb.Norm == "climb" {
		Printf("The cliff is too steep for climbing.\n")
		return true
	}
	return false
}

func RainbowFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "cross" || G.ActVerb.Norm == "through" {
		if G.Here == &CanyonView {
			Printf("From here?!?\n")
			return true
		}
		if G.RainbowFlag {
			if G.Here == &AragainFalls {
				Goto(&EndOfRainbow, true)
			} else if G.Here == &EndOfRainbow {
				Goto(&AragainFalls, true)
			} else {
				Printf("You'll have to say which way...\n")
			}
		} else {
			Printf("Can you walk on water vapor?\n")
		}
		return true
	}
	if G.ActVerb.Norm == "look under" {
		Printf("The Frigid River flows under the rainbow.\n")
		return true
	}
	return false
}

func RopeFcn(arg ActArg) bool {
	if G.Here != &DomeRoom {
		G.DomeFlag = false
		if G.ActVerb.Norm == "tie" {
			Printf("You can't tie the rope to that.\n")
			return true
		}
		return false
	}
	if G.ActVerb.Norm == "tie" {
		if G.IndirObj == &Railing {
			if G.DomeFlag {
				Printf("The rope is already tied to it.\n")
			} else {
				Printf("The rope drops over the side and comes within ten feet of the floor.\n")
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
				Printf("Your attempt to tie up the %s awakens him.", G.DirObj.Desc)
				Awaken(G.DirObj)
			} else {
				Printf("The %s struggles and you cannot tie him up.\n", G.DirObj.Desc)
			}
		} else {
			Printf("Why would you tie up a %s?\n", G.DirObj.Desc)
		}
		return true
	}
	if G.ActVerb.Norm == "untie" {
		if G.DomeFlag {
			G.DomeFlag = false
			Rope.Take(FlgNoDesc)
			Printf("The rope is now untied.\n")
		} else {
			Printf("It is not tied to anything.\n")
		}
		return true
	}
	if G.ActVerb.Norm == "drop" && G.Here == &DomeRoom && !G.DomeFlag {
		Rope.MoveTo(&TorchRoom)
		Printf("The rope drops gently to the floor below.\n")
		return true
	}
	if G.ActVerb.Norm == "take" {
		if G.DomeFlag {
			Printf("The rope is tied to the railing.\n")
			return true
		}
	}
	return false
}

func EggObjectFcn(arg ActArg) bool {
	if (G.ActVerb.Norm == "open" || G.ActVerb.Norm == "mung") && G.DirObj == &Egg {
		if G.DirObj.Has(FlgOpen) {
			Printf("The egg is already open.\n")
			return true
		}
		if G.IndirObj == nil {
			Printf("You have neither the tools nor the expertise.\n")
			return true
		}
		if G.IndirObj == &Hands {
			Printf("I doubt you could do that without damaging it.\n")
			return true
		}
		if G.IndirObj.Has(FlgWeapon) || G.IndirObj.Has(FlgTool) || G.ActVerb.Norm == "mung" {
			Printf("The egg is now open, but the clumsiness of your attempt has seriously compromised its esthetic appeal.")
			BadEgg()
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
	if G.ActVerb.Norm == "climb on" || G.ActVerb.Norm == "hatch" {
		Printf("There is a noticeable crunch from beneath you, and inspection reveals that the egg is lying open, badly damaged.")
		BadEgg()
		Printf("\n")
		return true
	}
	if G.ActVerb.Norm == "open" || G.ActVerb.Norm == "mung" || G.ActVerb.Norm == "throw" {
		if G.ActVerb.Norm == "throw" {
			G.DirObj.MoveTo(G.Here)
		}
		Printf("Your rather indelicate handling of the egg has caused it some damage, although you have succeeded in opening it.")
		BadEgg()
		Printf("\n")
		return true
	}
	return false
}

func CanaryObjectFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "wind" {
		if G.DirObj == &Canary {
			if !G.SingSong && ForestRoomQ() {
				Printf("The canary chirps, slightly off-key, an aria from a forgotten opera. From out of the greenery flies a lovely songbird. It perches on a limb just over your head and opens its beak to sing. As it does so a beautiful brass bauble drops from its mouth, bounces off the top of your head, and lands glimmering in the grass. As the canary winds down, the songbird flies away.\n")
				G.SingSong = true
				dest := G.Here
				if G.Here == &UpATree {
					dest = &Path
				}
				Bauble.MoveTo(dest)
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

func RugFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "raise" {
		Printf("The rug is too heavy to lift")
		if G.RugMoved {
			Printf(".\n")
		} else {
			Printf(", but in trying to take it you have noticed an irregularity beneath it.\n")
		}
		return true
	}
	if G.ActVerb.Norm == "move" || G.ActVerb.Norm == "push" {
		if G.RugMoved {
			Printf("Having moved the carpet previously, you find it impossible to move it again.\n")
		} else {
			Printf("With a great effort, the rug is moved to one side of the room, revealing the dusty cover of a closed trap door.\n")
			TrapDoor.Take(FlgInvis)
			ThisIsIt(&TrapDoor)
			G.RugMoved = true
		}
		return true
	}
	if G.ActVerb.Norm == "take" {
		Printf("The rug is extremely heavy and cannot be carried.\n")
		return true
	}
	if G.ActVerb.Norm == "look under" && !G.RugMoved && !TrapDoor.Has(FlgOpen) {
		Printf("Underneath the rug is a closed trap door. As you drop the corner of the rug, the trap door is once again concealed from view.\n")
		return true
	}
	if G.ActVerb.Norm == "climb on" {
		if !G.RugMoved && !TrapDoor.Has(FlgOpen) {
			Printf("As you sit, you notice an irregularity underneath it. Rather than be uncomfortable, you stand up again.\n")
		} else {
			Printf("I suppose you think it's a magic carpet?\n")
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
				Printf("You can see a scarab here in the sand.\n")
				ThisIsIt(&Scarab)
				Scarab.Take(FlgInvis)
			}
		} else {
			Printf("%s\n", BDigs[G.BeachDig])
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
		Printf("You are in the kitchen of the white house. A table seems to have been used recently for the preparation of food. A passage leads to the west and a dark staircase can be seen leading upward. A dark chimney leads down and to the east is a small window which is ")
		if KitchenWindow.Has(FlgOpen) {
			Printf("open.\n")
		} else {
			Printf("slightly ajar.\n")
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
		Printf("You are in the living room. There is a doorway to the east")
		if G.MagicFlag {
			Printf(". To the west is a cyclops-shaped opening in an old wooden door, above which is some strange gothic lettering, ")
		} else {
			Printf(", a wooden door with strange gothic lettering to the west, which appears to be nailed shut, ")
		}
		Printf("a trophy case, ")
		if G.RugMoved && TrapDoor.Has(FlgOpen) {
			Printf("and a rug lying beside an open trap door.")
		} else if G.RugMoved {
			Printf("and a closed trap door at your feet.")
		} else if TrapDoor.Has(FlgOpen) {
			Printf("and an open trap door at your feet.")
		} else {
			Printf("and a large oriental rug in the center of the room.")
		}
		Printf("\n")
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
		Printf("You are in a dark and damp cellar with a narrow passageway leading north, and a crawlway to the south. On the west is the bottom of a steep metal ramp which is unclimbable.\n")
		return true
	}
	if arg == ActEnter {
		if TrapDoor.Has(FlgOpen) && !TrapDoor.Has(FlgTouch) {
			TrapDoor.Take(FlgOpen)
			TrapDoor.Give(FlgTouch)
			Printf("The trap door crashes shut, and you hear someone barring it.\n\n")
		}
		return false
	}
	return false
}

func StoneBarrowFcn(arg ActArg) bool {
	if arg == ActBegin {
		if G.ActVerb.Norm == "enter" || (G.ActVerb.Norm == "walk" && G.Params.HasWalkDir && (G.Params.WalkDir == West || G.Params.WalkDir == In)) || (G.ActVerb.Norm == "through" && G.DirObj == &Barrow) {
			Printf("Inside the Barrow\nAs you enter the barrow, the door closes inexorably behind you. Around you it is dark, but ahead is an enormous cavern, brightly lit. Through its center runs a wide stream. Spanning the stream is a small wooden footbridge, and beyond a path leads into a dark tunnel. Above the bridge, floating in the air, is a large sign. It reads:  All ye who stand before this bridge have completed a great and perilous adventure which has tested your wit and courage. You have mastered the first part of the ZORK trilogy. Those who pass over this bridge must be prepared to undertake an even greater adventure that will severely test your skill and bravery!\n\nThe ZORK trilogy continues with \"ZORK II: The Wizard of Frobozz\" and is completed in \"ZORK III: The Dungeon Master.\"\n")
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
		Printf("You are in a clearing, with a forest surrounding you on all sides. A path leads south.")
		if Grate.Has(FlgOpen) {
			Printf("\nThere is an open grating, descending into darkness.")
		} else if G.GrateRevealed {
			Printf("\nThere is a grating securely fastened into the ground.")
		}
		Printf("\n")
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
		Printf("You are in a small room near the maze. There are twisty passages in the immediate vicinity.\n")
		if Grate.Has(FlgOpen) {
			Printf("Above you is an open grating with sunlight pouring in.")
		} else if G.GrUnlock {
			Printf("Above you is a grating.")
		} else {
			Printf("Above you is a grating locked with a skull-and-crossbones lock.")
		}
		Printf("\n")
		return true
	}
	return false
}

func CyclopsRoomFcn(arg ActArg) bool {
	if arg == ActLook {
		Printf("This room has an exit on the northwest, and a staircase leading up.\n")
		if G.CyclopsFlag && !G.MagicFlag {
			Printf("The cyclops is sleeping blissfully at the foot of the stairs.\n")
		} else if G.MagicFlag {
			Printf("The east wall, previously solid, now has a cyclops-sized opening in it.\n")
		} else if G.CycloWrath == 0 {
			Printf("A cyclops, who looks prepared to eat horses (much less mere adventurers), blocks the staircase. From his state of health, and the bloodstains on the walls, you gather that he is not very friendly, though he likes people.\n")
		} else if G.CycloWrath > 0 {
			Printf("The cyclops is standing in the corner, eyeing you closely. I don't think he likes you very much. He looks extremely hungry, even for a cyclops.\n")
		} else {
			Printf("The cyclops, having eaten the hot peppers, appears to be gasping. His enflamed tongue protrudes from his man-sized mouth.\n")
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
			Printf("You hear a scream of anguish as you violate the robber's hideaway. Using passages unknown to you, he rushes to its defense.\n")
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
			Printf("You are in a long room, to the north of which was formerly a lake. However, with the water level lowered, there is merely a wide stream running through the center of the room.")
		} else if G.GatesOpen {
			Printf("You are in a long room. To the north is a large lake, too deep to cross. You notice, however, that the water level appears to be dropping at a rapid rate. Before long, it might be possible to cross to the other side from here.")
		} else if G.LowTide {
			Printf("You are in a long room, to the north of which is a wide area which was formerly a reservoir, but now is merely a stream. You notice, however, that the level of the stream is rising quickly and that before long it will be impossible to cross here.")
		} else {
			Printf("You are in a long room on the south shore of a large lake, far too deep and wide for crossing.")
		}
		Printf("\nThere is a path along the stream to the east or west, a steep pathway climbing southwest along the edge of a chasm, and a path leading into a canyon to the southeast.\n")
		return true
	}
	return false
}

func ReservoirFcn(arg ActArg) bool {
	if arg == ActEnd && !G.Winner.Location().Has(FlgVeh) && !G.GatesOpen && G.LowTide {
		Printf("You notice that the water level here is rising rapidly. The currents are also becoming stronger. Staying here seems quite perilous!\n")
		return true
	}
	if arg == ActLook {
		if G.LowTide {
			Printf("You are on what used to be a large lake, but which is now a large mud pile. There are \"shores\" to the north and south.")
		} else {
			Printf("You are on the lake. Beaches can be seen north and south. Upstream a small stream enters the lake through a narrow cleft in the rocks. The dam can be seen downstream.")
		}
		Printf("\n")
		return true
	}
	return false
}

func ReservoirNorthFcn(arg ActArg) bool {
	if arg == ActLook {
		if G.LowTide && G.GatesOpen {
			Printf("You are in a large cavernous room, the south of which was formerly a lake. However, with the water level lowered, there is merely a wide stream running through there.")
		} else if G.GatesOpen {
			Printf("You are in a large cavernous area. To the south is a wide lake, whose water level appears to be falling rapidly.")
		} else if G.LowTide {
			Printf("You are in a cavernous area, to the south of which is a very wide stream. The level of the stream is rising rapidly, and it appears that before long it will be impossible to cross to the other side.")
		} else {
			Printf("You are in a large cavernous room, north of a large lake.")
		}
		Printf("\nThere is a slimy stairway leaving the room to the north.\n")
		return true
	}
	return false
}

func MirrorRoomFcn(arg ActArg) bool {
	if arg == ActLook {
		Printf("You are in a large square room with tall ceilings. On the south wall is an enormous mirror which fills the entire wall. There are exits on the other three sides of the room.\n")
		if G.MirrorMung {
			Printf("Unfortunately, the mirror has been destroyed by your recklessness.\n")
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
			Printf("A gust of wind blows out your candles!\n")
			G.Lit = IsLit(G.Here, true)
			if !G.Lit {
				Printf("It is now completely dark.\n")
			}
			return true
		}
	}
	return false
}

func LLDRoomFcn(arg ActArg) bool {
	if arg == ActLook {
		Printf("You are outside a large gateway, on which is inscribed\n\n  Abandon every hope\nall ye who enter here!\n\nThe gate is open; through it you can see a desolation, with a pile of mangled bodies in one corner. Thousands of voices, lamenting some hideous fate, can be heard.\n")
		if !G.LLDFlag && !G.Dead {
			Printf("The way through the gate is barred by evil spirits, who jeer at your attempts to pass.\n")
		}
		return true
	}
	if arg == ActBegin {
		if G.ActVerb.Norm == "exorcise" && !G.LLDFlag {
			if Bell.IsIn(G.Winner) && Book.IsIn(G.Winner) && Candles.IsIn(G.Winner) {
				Printf("You must perform the ceremony.\n")
			} else {
				Printf("You aren't equipped for an exorcism.\n")
			}
			return true
		}
		if !G.LLDFlag && G.ActVerb.Norm == "ring" && G.DirObj == &Bell {
			G.XB = true
			RemoveCarefully(&Bell)
			ThisIsIt(&HotBell)
			HotBell.MoveTo(G.Here)
			Printf("The bell suddenly becomes red hot and falls to the ground. The wraiths, as if paralyzed, stop their jeering and slowly turn to face you. On their ashen faces, the expression of a long-forgotten terror takes shape.\n")
			if Candles.IsIn(G.Winner) {
				Printf("In your confusion, the candles drop to the ground (and they are out).\n")
				Candles.MoveTo(G.Here)
				Candles.Take(FlgOn)
				QueueInt("ICandles", false).Run = false
			}
			Queue("IXb", 6).Run = true
			Queue("IXbh", 20).Run = true
			return true
		}
		if G.XC && G.ActVerb.Norm == "read" && G.DirObj == &Book && !G.LLDFlag {
			Printf("Each word of the prayer reverberates through the hall in a deafening confusion. As the last word fades, a voice, loud and commanding, speaks: \"Begone, fiends!\" A heart-stopping scream fills the cavern, and the spirits, sensing a greater power, flee through the walls.\n")
			RemoveCarefully(&Ghosts)
			G.LLDFlag = true
			QueueInt("IXc", false).Run = false
			return true
		}
	}
	if arg == ActEnd {
		if G.XB && Candles.IsIn(G.Winner) && Candles.Has(FlgOn) && !G.XC {
			G.XC = true
			Printf("The flames flicker wildly and appear to dance. The earth beneath your feet trembles, and your legs nearly buckle beneath you. The spirits cower at your unearthly power.\n")
			QueueInt("IXb", false).Run = false
			Queue("IXc", 3).Run = true
		}
	}
	return false
}

func DomeRoomFcn(arg ActArg) bool {
	if arg == ActLook {
		Printf("You are at the periphery of a large dome, which forms the ceiling of another room below. Protecting you from a precipitous drop is a wooden railing which circles the dome.\n")
		if G.DomeFlag {
			Printf("Hanging down from the railing is a rope which ends about ten feet from the floor below.\n")
		}
		return true
	}
	if arg == ActEnter {
		if G.Dead {
			Printf("As you enter the dome you feel a strong pull as if from a wind drawing you over the railing and down.\n")
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
		Printf("This is a large room with a prominent doorway leading to a down staircase. Above you is a large dome. Up around the edge of the dome (20 feet up) is a wooden railing. In the center of the room sits a white marble pedestal.\n")
		if G.DomeFlag {
			Printf("A piece of rope descends from the railing above, ending some five feet above your head.\n")
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
		Printf("This is a large, cold room whose sole exit is to the north. In one corner there is a machine which is reminiscent of a clothes dryer. On its face is a switch which is labelled \"START\". The switch does not appear to be manipulable by any human hand (unless the fingers are about 1/16 by 1/4 inch). On the front of the machine is a large lid, which is ")
		if Machine.Has(FlgOpen) {
			Printf("open.\n")
		} else {
			Printf("closed.\n")
		}
		return true
	}
	return false
}

func LoudRoomFcn(arg ActArg) bool {
	if arg == ActLook {
		Printf("This is a large room with a ceiling which cannot be detected from the ground. There is a narrow passage from east to west and a stone stairway leading upward.")
		if G.LoudFlag || (!G.GatesOpen && G.LowTide) {
			Printf(" The room is eerie in its quietness.")
		} else {
			Printf(" The room is deafeningly loud with an undetermined rushing sound. The sound seems to reverberate from all of the walls, making it difficult even to think.")
		}
		Printf("\n")
		return true
	}
	if arg == ActEnd && G.GatesOpen && !G.LowTide {
		Printf("It is unbearably loud here, with an ear-splitting roar seeming to come from all around you. There is a pounding in your head which won't stop. With a tremendous effort, you scramble out of the room.\n\n")
		dest := LoudRuns[G.Rand.Intn(len(LoudRuns))]
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
			Printf("The rest of your commands have been lost in the noise.\n")
			G.Params.Continue = NumUndef
		}
		// In the original, this has a special read loop. We simplify.
		return false
	}
	if G.ActVerb.Norm == "echo" {
		if G.LoudFlag || (!G.GatesOpen && G.LowTide) {
			// Room is already quiet
			Printf("echo echo ...\n")
			return true
		}
		Printf("The acoustics of the room change subtly.\n")
		G.LoudFlag = true
		return true
	}
	return false
}

func DeepCanyonFcn(arg ActArg) bool {
	if arg == ActLook {
		Printf("You are on the south edge of a deep canyon. Passages lead off to the east, northwest and southwest. A stairway leads down.")
		if G.GatesOpen && !G.LowTide {
			Printf(" You can hear a loud roaring sound, like that of rushing water, from below.")
		} else if !G.GatesOpen && G.LowTide {
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
				Printf("How sad for an aspiring adventurer to light a %s in a room which reeks of gas. Fortunately, there is justice in the world.\n", G.DirObj.Desc)
			} else {
				Printf("Oh dear. It appears that the smell coming from this room was coal gas. I would have thought twice about carrying flaming objects in here.\n")
			}
			JigsUp("\n      ** BOOOOOOOOOOOM **", false)
			return true
		}
	}
	return false
}

func BatsRoomFcn(arg ActArg) bool {
	if arg == ActLook {
		Printf("You are in a small room which has doors only to the east and south.\n")
		return true
	}
	if arg == ActEnter && !G.Dead {
		if Garlic.Location() != G.Winner && !Garlic.IsIn(G.Here) {
			VLook(ActUnk)
			Printf("\n")
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
		Printf("You are about 10 feet above the ground nestled among some large branches. The nearest branch above you is above your reach.\n")
		if Path.HasChildren() && len(Path.Children) > 0 {
			Printf("On the ground below you can see:  ")
			PrintContents(&Path)
			Printf(".\n")
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
				Printf("The nest falls to the ground, and the egg spills out of it, seriously damaged.\n")
				RemoveCarefully(&Egg)
				BrokenEgg.MoveTo(&Path)
				return true
			}
			if G.DirObj == &Egg {
				Printf("The egg falls to the ground and springs open, seriously damaged.")
				Egg.MoveTo(&Path)
				BadEgg()
				Printf("\n")
				return true
			}
			if G.DirObj != G.Winner && G.DirObj != &Tree {
				G.DirObj.MoveTo(&Path)
				Printf("The %s falls to the ground.\n", G.DirObj.Desc)
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
			Printf("You cannot enter in your condition.\n")
			return true
		}
		return false
	}
	if G.ActVerb.Norm == "brief" || G.ActVerb.Norm == "verbose" || G.ActVerb.Norm == "super-brief" || G.ActVerb.Norm == "version" {
		return false
	}
	if G.ActVerb.Norm == "attack" || G.ActVerb.Norm == "mung" || G.ActVerb.Norm == "alarm" || G.ActVerb.Norm == "swing" {
		Printf("All such attacks are vain in your condition.\n")
		return true
	}
	if G.ActVerb.Norm == "open" || G.ActVerb.Norm == "close" || G.ActVerb.Norm == "eat" || G.ActVerb.Norm == "drink" ||
		G.ActVerb.Norm == "inflate" || G.ActVerb.Norm == "deflate" || G.ActVerb.Norm == "turn" || G.ActVerb.Norm == "burn" ||
		G.ActVerb.Norm == "tie" || G.ActVerb.Norm == "untie" || G.ActVerb.Norm == "rub" {
		Printf("Even such an action is beyond your capabilities.\n")
		return true
	}
	if G.ActVerb.Norm == "wait" {
		Printf("Might as well. You've got an eternity.\n")
		return true
	}
	if G.ActVerb.Norm == "lamp on" {
		Printf("You need no light to guide you.\n")
		return true
	}
	if G.ActVerb.Norm == "score" {
		Printf("You're dead! How can you think of your score?\n")
		return true
	}
	if G.ActVerb.Norm == "take" || G.ActVerb.Norm == "rub" {
		Printf("Your hand passes through its object.\n")
		return true
	}
	if G.ActVerb.Norm == "drop" || G.ActVerb.Norm == "throw" || G.ActVerb.Norm == "inventory" {
		Printf("You have no possessions.\n")
		return true
	}
	if G.ActVerb.Norm == "diagnose" {
		Printf("You are dead.\n")
		return true
	}
	if G.ActVerb.Norm == "look" {
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
			Printf("From the distance the sound of a lone trumpet is heard. The room becomes very bright and you feel disembodied. In a moment, the brightness fades and you find yourself rising as if from a long sleep, deep in the woods. In the distance you can faintly hear a songbird and the sounds of the forest.\n\n")
			Goto(&Forest1, true)
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
			Printf("You'd better have more light than from the %s.\n", obj.Desc)
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
					Printf("%s\n", msg)
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
	Printf("The match has gone out.\n")
	Match.Take(FlgFlame)
	Match.Take(FlgOn)
	G.Lit = IsLit(G.Here, true)
	return true
}

func IXb() bool {
	if !G.XC {
		if G.Here == &EnteranceToHades {
			Printf("The tension of this ceremony is broken, and the wraiths, amused but shaken at your clumsy attempt, resume their hideous jeering.\n")
		}
	}
	G.XB = false
	return true
}

func IXbh() bool {
	RemoveCarefully(&HotBell)
	Bell.MoveTo(&EnteranceToHades)
	if G.Here == &EnteranceToHades {
		Printf("The bell appears to have cooled down.\n")
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
			Printf("%s\n", Cyclomad[idx])
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
		Printf("You hear in the distance the chirping of a song bird.\n")
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
		Printf("The grating is closed!\n")
		ThisIsIt(&Grate)
		return nil
	}
	Printf("You can't go that way.\n")
	return nil
}

func TrapDoorExitFcn() *Object {
	if G.RugMoved {
		if TrapDoor.Has(FlgOpen) {
			return &Cellar
		}
		Printf("The trap door is closed.\n")
		ThisIsIt(&TrapDoor)
		return nil
	}
	Printf("You can't go that way.\n")
	return nil
}

func UpChimneyFcn() *Object {
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
	if count <= 2 && Lamp.IsIn(G.Winner) {
		if !TrapDoor.Has(FlgOpen) {
			TrapDoor.Take(FlgTouch)
		}
		return &Kitchen
	}
	Printf("You can't get up there with what you're carrying.\n")
	return nil
}

func MazeDiodesFcn() *Object {
	Printf("You won't be able to get back up to the tunnel you are going through when it gets to the next room.\n\n")
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
		Printf("You look before leaping, and realize that you would never survive.\n")
		return true
	}
	if G.ActVerb.Norm == "cross" {
		Printf("It's too far to jump, and there's no bridge.\n")
		return true
	}
	if (G.ActVerb.Norm == "put" || G.ActVerb.Norm == "throw off") && G.IndirObj == &PseudoObject {
		Printf("The %s drops out of sight into the chasm.\n", G.DirObj.Desc)
		RemoveCarefully(G.DirObj)
		return true
	}
	return false
}

func LakePseudo(arg ActArg) bool {
	if G.LowTide {
		Printf("There's not much lake left....\n")
		return true
	}
	if G.ActVerb.Norm == "cross" {
		Printf("It's too wide to cross.\n")
		return true
	}
	if G.ActVerb.Norm == "through" {
		Printf("You can't swim in this lake.\n")
		return true
	}
	return false
}

func StreamPseudo(arg ActArg) bool {
	if G.ActVerb.Norm == "swim" || G.ActVerb.Norm == "through" {
		Printf("You can't swim in the stream.\n")
		return true
	}
	if G.ActVerb.Norm == "cross" {
		Printf("The other side is a sheer rock cliff.\n")
		return true
	}
	return false
}

func DomePseudo(arg ActArg) bool {
	if G.ActVerb.Norm == "kiss" {
		Printf("No.\n")
		return true
	}
	return false
}

func GatePseudo(arg ActArg) bool {
	if G.ActVerb.Norm == "through" {
		DoWalk(In)
		return true
	}
	Printf("The gate is protected by an invisible force. It makes your teeth ache to touch it.\n")
	return true
}

func DoorPseudo(arg ActArg) bool {
	if G.ActVerb.Norm == "open" || G.ActVerb.Norm == "close" {
		Printf("The door won't budge.\n")
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
		Printf("Some paint chips away, revealing more paint.\n")
		return true
	}
	return false
}

func GasPseudo(arg ActArg) bool {
	if G.ActVerb.Norm == "breathe" {
		Printf("There is too much gas to blow away.\n")
		return true
	}
	if G.ActVerb.Norm == "smell" {
		Printf("It smells like coal gas in here.\n")
		return true
	}
	return false
}

func ChainPseudo(arg ActArg) bool {
	if G.ActVerb.Norm == "take" || G.ActVerb.Norm == "move" {
		Printf("The chain is secure.\n")
		return true
	}
	if G.ActVerb.Norm == "raise" || G.ActVerb.Norm == "lower" {
		Printf("Perhaps you should do that to the basket.\n")
		return true
	}
	if G.ActVerb.Norm == "examine" {
		Printf("The chain secures a basket within the shaft.\n")
		return true
	}
	return false
}

func BarrowDoorFcn2(arg ActArg) bool {
	return false
}
