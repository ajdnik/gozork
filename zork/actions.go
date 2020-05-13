package zork

type BlowRes int

const (
	BlowUnk BlowRes = iota
	BlowMissed
	BlowUncon
	BlowKill
	BlowLightWnd
	BlowHeavyWnd
	BlowStag
	BlowLoseWpn
	BlowHesitate
	BlowSitDuck
)

// Tables of melee results

var (
	Def1    = [13]BlowRes{BlowMissed, BlowMissed, BlowMissed, BlowMissed, BlowStag, BlowStag, BlowUncon, BlowUncon, BlowKill, BlowKill, BlowKill, BlowKill, BlowKill}
	Def2A   = [10]BlowRes{BlowMissed, BlowMissed, BlowMissed, BlowMissed, BlowMissed, BlowStag, BlowStag, BlowLightWnd, BlowLightWnd, BlowUncon}
	Def2B   = [12]BlowRes{BlowMissed, BlowMissed, BlowMissed, BlowStag, BlowStag, BlowLightWnd, BlowLightWnd, BlowLightWnd, BlowUncon, BlowKill, BlowKill, BlowKill}
	Def3A   = [11]BlowRes{BlowMissed, BlowMissed, BlowMissed, BlowMissed, BlowMissed, BlowStag, BlowStag, BlowLightWnd, BlowLightWnd, BlowHeavyWnd, BlowHeavyWnd}
	Def3B   = [11]BlowRes{BlowMissed, BlowMissed, BlowMissed, BlowStag, BlowStag, BlowLightWnd, BlowLightWnd, BlowLightWnd, BlowHeavyWnd, BlowHeavyWnd, BlowHeavyWnd}
	Def3C   = [10]BlowRes{BlowMissed, BlowStag, BlowStag, BlowLightWnd, BlowLightWnd, BlowLightWnd, BlowLightWnd, BlowHeavyWnd, BlowHeavyWnd, BlowHeavyWnd}
	Def1Res = [4]BlowRes{Def1[0], BlowUnk, BlowUnk}
	Def2Res = [4]BlowRes{Def2A[0], Def2B[0], BlowUnk, BlowUnk}
	Def3Res = [5]BlowRes{Def3A[0], BlowUnk, Def3B[0], BlowUnk, Def3C[0]}
)

func IFight() bool {
	return false
}

func ISword() bool {
	return false
}

func IThief() bool {
	return false
}

func ICandles() bool {
	return false
}

func ILantern() bool {
	return false
}

func IXb() bool {
	return false
}

func IXc() bool {
	return false
}

func ICyclops() bool {
	return false
}

func IForestRandom() bool {
	return false
}

func IMatch() bool {
	return false
}

func RBoatFcn(arg ActArg) bool {
	return false
}

func MailboxFcn(arg ActArg) bool {
	if ActVerb.Norm == "take" && DirObj == &Mailbox {
		Print("It is securely anchored.", Newline)
		return true
	}
	return false
}

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

func ForestFcn(arg ActArg) bool {
	return false
}

func BoardFcn(arg ActArg) bool {
	if ActVerb.Norm == "take" || ActVerb.Norm == "examine" {
		Print("The boards are securely fastened.", Newline)
		return true
	}
	return false
}

func TrophyCaseFcn(arg ActArg) bool {
	return false
}

func LivingRoomFcn(arg ActArg) bool {
	return false
}

func BoardedWindowFcn(arg ActArg) bool {
	return false
}

func NailsPseudo(arg ActArg) bool {
	return false
}

func SandFunction(arg ActArg) bool {
	return false
}

func Awaken(o *Object) bool {
	return false
}

func HeroBlow() bool {
	return false
}

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

func DeadFunction(arg ActArg) bool {
	return false
}

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
	}
	TrapDoor.Take(FlgTouch)
	Params.Continue = NumUndef
	RandomizeObjects()
	KillInterrupts()
	// TODO: return fatal
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
	for _, child := range Winner.Children {
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
