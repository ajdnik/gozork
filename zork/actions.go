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
	return false
}

func WestHouseFcn(arg ActArg) bool {
	return false
}

func WhiteHouseFcn(arg ActArg) bool {
	return false
}

func ForestFcn(arg ActArg) bool {
	return false
}

func BoardFcn(arg ActArg) bool {
	return false
}
