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

func RBoatFcn(arg ActArg) bool {
	return false
}

func MailboxFcn(arg ActArg) bool {
	if ActVerb == "take" && DirObj == &Mailbox {
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
	return false
}

func ForestFcn(arg ActArg) bool {
	return false
}

func BoardFcn(arg ActArg) bool {
	if ActVerb == "take" || ActVerb == "examine" {
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

func NailsPseudo(arg ActArg) bool {
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
