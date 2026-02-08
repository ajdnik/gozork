package zork

import "math/rand"

// ================================================================
// COMBAT SYSTEM TYPES AND TABLES
// ================================================================

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

// Combat mode constants
const (
	FBusy       = 1
	FDead       = 2
	FUnconscious = 3
	FConscious  = 4
	FFirst      = 5
)

// Combat strength constants
const (
	StrengthMax = 7
	StrengthMin = 2
	CureWait    = 30
)

// Villain table field indices
const (
	VVillain  = 0
	VBest     = 1
	VBestAdv  = 2
	VProb     = 3
	VMsgs     = 4
)

// Melee message marker constants
const (
	FWep = 0 // means print weapon name
	FDef = 1 // means print defender name
)

// MeleePart is a fragment of a melee message (either a string or FWep/FDef marker)
type MeleePart struct {
	Text   string
	Marker int // -1 = normal text, FWep = weapon, FDef = defender
}

func mp(s string) MeleePart { return MeleePart{Text: s, Marker: -1} }
func mw() MeleePart         { return MeleePart{Marker: FWep} }
func md() MeleePart         { return MeleePart{Marker: FDef} }

// A MeleeMsg is a sequence of parts that combine to form one message
type MeleeMsg []MeleePart

// MeleeSet is a set of alternative messages for one outcome type
type MeleeSet []MeleeMsg

// MeleeTable is indexed by BlowRes-1 (Missed=0, Unconscious=1, Killed=2, LightWound=3, HeavyWound=4, Stagger=5, LoseWeapon=6, Hesitate=7, SittingDuck=8)
type MeleeTable [9]MeleeSet

// VillainEntry holds a villain's combat parameters
type VillainEntry struct {
	Villain *Object
	Best    *Object // best weapon against this villain
	BestAdv int     // advantage conferred by best weapon
	Prob    int     // probability of waking if unconscious
	Msgs    *MeleeTable
}

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

	// Hero melee messages (indexed by outcome: missed, unconscious, killed, light-wound, heavy-wound, stagger, lose-weapon)
	HeroMelee = MeleeTable{
		// BlowMissed
		{
			{mp("Your "), mw(), mp(" misses the "), md(), mp(" by an inch.")},
			{mp("A good slash, but it misses the "), md(), mp(" by a mile.")},
			{mp("You charge, but the "), md(), mp(" jumps nimbly aside.")},
			{mp("Clang! Crash! The "), md(), mp(" parries.")},
			{mp("A quick stroke, but the "), md(), mp(" is on guard.")},
			{mp("A good stroke, but it's too slow; the "), md(), mp(" dodges.")},
		},
		// BlowUncon
		{
			{mp("Your "), mw(), mp(" crashes down, knocking the "), md(), mp(" into dreamland.")},
			{mp("The "), md(), mp(" is battered into unconsciousness.")},
			{mp("A furious exchange, and the "), md(), mp(" is knocked out!")},
			{mp("The haft of your "), mw(), mp(" knocks out the "), md(), mp(".")},
			{mp("The "), md(), mp(" is knocked out!")},
		},
		// BlowKill
		{
			{mp("It's curtains for the "), md(), mp(" as your "), mw(), mp(" removes his head.")},
			{mp("The fatal blow strikes the "), md(), mp(" square in the heart: He dies.")},
			{mp("The "), md(), mp(" takes a fatal blow and slumps to the floor dead.")},
		},
		// BlowLightWnd
		{
			{mp("The "), md(), mp(" is struck on the arm; blood begins to trickle down.")},
			{mp("Your "), mw(), mp(" pinks the "), md(), mp(" on the wrist, but it's not serious.")},
			{mp("Your stroke lands, but it was only the flat of the blade.")},
			{mp("The blow lands, making a shallow gash in the "), md(), mp("'s arm!")},
		},
		// BlowHeavyWnd
		{
			{mp("The "), md(), mp(" receives a deep gash in his side.")},
			{mp("A savage blow on the thigh! The "), md(), mp(" is stunned but can still fight!")},
			{mp("Slash! Your blow lands! That one hit an artery, it could be serious!")},
			{mp("Slash! Your stroke connects! This could be serious!")},
		},
		// BlowStag
		{
			{mp("The "), md(), mp(" is staggered, and drops to his knees.")},
			{mp("The "), md(), mp(" is momentarily disoriented and can't fight back.")},
			{mp("The force of your blow knocks the "), md(), mp(" back, stunned.")},
			{mp("The "), md(), mp(" is confused and can't fight back.")},
			{mp("The quickness of your thrust knocks the "), md(), mp(" back, stunned.")},
		},
		// BlowLoseWpn
		{
			{mp("The "), md(), mp("'s weapon is knocked to the floor, leaving him unarmed.")},
			{mp("The "), md(), mp(" is disarmed by a subtle feint past his guard.")},
		},
		// BlowHesitate (not used for hero)
		{},
		// BlowSitDuck (not used for hero)
		{},
	}

	// Cyclops melee messages
	CyclopsMelee = MeleeTable{
		// BlowMissed
		{
			{mp("The Cyclops misses, but the backwash almost knocks you over.")},
			{mp("The Cyclops rushes you, but runs into the wall.")},
		},
		// BlowUncon
		{
			{mp("The Cyclops sends you crashing to the floor, unconscious.")},
		},
		// BlowKill
		{
			{mp("The Cyclops breaks your neck with a massive smash.")},
		},
		// BlowLightWnd
		{
			{mp("A quick punch, but it was only a glancing blow.")},
			{mp("A glancing blow from the Cyclops' fist.")},
		},
		// BlowHeavyWnd
		{
			{mp("The monster smashes his huge fist into your chest, breaking several ribs.")},
			{mp("The Cyclops almost knocks the wind out of you with a quick punch.")},
		},
		// BlowStag
		{
			{mp("The Cyclops lands a punch that knocks the wind out of you.")},
			{mp("Heedless of your weapons, the Cyclops tosses you against the rock wall of the room.")},
		},
		// BlowLoseWpn
		{
			{mp("The Cyclops grabs your "), mw(), mp(", tastes it, and throws it to the ground in disgust.")},
			{mp("The monster grabs you on the wrist, squeezes, and you drop your "), mw(), mp(" in pain.")},
		},
		// BlowHesitate
		{
			{mp("The Cyclops seems unable to decide whether to broil or stew his dinner.")},
		},
		// BlowSitDuck
		{
			{mp("The Cyclops, no sportsman, dispatches his unconscious victim.")},
		},
	}

	// Troll melee messages
	TrollMelee = MeleeTable{
		// BlowMissed
		{
			{mp("The troll swings his axe, but it misses.")},
			{mp("The troll's axe barely misses your ear.")},
			{mp("The axe sweeps past as you jump aside.")},
			{mp("The axe crashes against the rock, throwing sparks!")},
		},
		// BlowUncon
		{
			{mp("The flat of the troll's axe hits you delicately on the head, knocking you out.")},
		},
		// BlowKill
		{
			{mp("The troll neatly removes your head.")},
			{mp("The troll's axe stroke cleaves you from the nave to the chops.")},
			{mp("The troll's axe removes your head.")},
		},
		// BlowLightWnd
		{
			{mp("The axe gets you right in the side. Ouch!")},
			{mp("The flat of the troll's axe skins across your forearm.")},
			{mp("The troll's swing almost knocks you over as you barely parry in time.")},
			{mp("The troll swings his axe, and it nicks your arm as you dodge.")},
		},
		// BlowHeavyWnd
		{
			{mp("The troll charges, and his axe slashes you on your "), mw(), mp(" arm.")},
			{mp("An axe stroke makes a deep wound in your leg.")},
			{mp("The troll's axe swings down, gashing your shoulder.")},
		},
		// BlowStag
		{
			{mp("The troll hits you with a glancing blow, and you are momentarily stunned.")},
			{mp("The troll swings; the blade turns on your armor but crashes broadside into your head.")},
			{mp("You stagger back under a hail of axe strokes.")},
			{mp("The troll's mighty blow drops you to your knees.")},
		},
		// BlowLoseWpn
		{
			{mp("The axe hits your "), mw(), mp(" and knocks it spinning.")},
			{mp("The troll swings, you parry, but the force of his blow knocks your "), mw(), mp(" away.")},
			{mp("The axe knocks your "), mw(), mp(" out of your hand. It falls to the floor.")},
		},
		// BlowHesitate
		{
			{mp("The troll hesitates, fingering his axe.")},
			{mp("The troll scratches his head ruminatively:  Might you be magically protected, he wonders?")},
		},
		// BlowSitDuck
		{
			{mp("Conquering his fears, the troll puts you to death.")},
		},
	}

	// Thief melee messages
	ThiefMelee = MeleeTable{
		// BlowMissed
		{
			{mp("The thief stabs nonchalantly with his stiletto and misses.")},
			{mp("You dodge as the thief comes in low.")},
			{mp("You parry a lightning thrust, and the thief salutes you with a grim nod.")},
			{mp("The thief tries to sneak past your guard, but you twist away.")},
		},
		// BlowUncon
		{
			{mp("Shifting in the midst of a thrust, the thief knocks you unconscious with the haft of his stiletto.")},
			{mp("The thief knocks you out.")},
		},
		// BlowKill
		{
			{mp("Finishing you off, the thief inserts his blade into your heart.")},
			{mp("The thief comes in from the side, feints, and inserts the blade into your ribs.")},
			{mp("The thief bows formally, raises his stiletto, and with a wry grin, ends the battle and your life.")},
		},
		// BlowLightWnd
		{
			{mp("A quick thrust pinks your left arm, and blood starts to trickle down.")},
			{mp("The thief draws blood, raking his stiletto across your arm.")},
			{mp("The stiletto flashes faster than you can follow, and blood wells from your leg.")},
			{mp("The thief slowly approaches, strikes like a snake, and leaves you wounded.")},
		},
		// BlowHeavyWnd
		{
			{mp("The thief strikes like a snake! The resulting wound is serious.")},
			{mp("The thief stabs a deep cut in your upper arm.")},
			{mp("The stiletto touches your forehead, and the blood obscures your vision.")},
			{mp("The thief strikes at your wrist, and suddenly your grip is slippery with blood.")},
		},
		// BlowStag
		{
			{mp("The butt of his stiletto cracks you on the skull, and you stagger back.")},
			{mp("The thief rams the haft of his blade into your stomach, leaving you out of breath.")},
			{mp("The thief attacks, and you fall back desperately.")},
		},
		// BlowLoseWpn
		{
			{mp("A long, theatrical slash. You catch it on your "), mw(), mp(", but the thief twists his knife, and the "), mw(), mp(" goes flying.")},
			{mp("The thief neatly flips your "), mw(), mp(" out of your hands, and it drops to the floor.")},
			{mp("You parry a low thrust, and your "), mw(), mp(" slips out of your hand.")},
		},
		// BlowHesitate
		{
			{mp("The thief, a man of superior breeding, pauses for a moment to consider the propriety of finishing you off.")},
			{mp("The thief amuses himself by searching your pockets.")},
			{mp("The thief entertains himself by rifling your pack.")},
		},
		// BlowSitDuck
		{
			{mp("The thief, forgetting his essentially genteel upbringing, cuts your throat.")},
			{mp("The thief, a pragmatist, dispatches you as a threat to his livelihood.")},
		},
	}

	// Villain table - initialized in FinalizeGameObjects to avoid init cycles
	Villains []*VillainEntry
)

// ================================================================
// HELPER FUNCTIONS
// ================================================================

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

func Rob(what, where *Object, prob int) bool {
	robbed := false
	for _, x := range what.Children {
		if x.Has(FlgInvis) || x.Has(FlgSacred) {
			continue
		}
		if x.TValue <= 0 {
			continue
		}
		if prob > 0 && !Prob(prob, false) {
			continue
		}
		x.MoveTo(where)
		x.Give(FlgTouch)
		if where == &Thief {
			x.Give(FlgInvis)
		}
		robbed = true
	}
	return robbed
}

func WeaponFunction(w, v *Object) bool {
	if !v.IsIn(Here) {
		return false
	}
	if ActVerb.Norm == "take" {
		if w.IsIn(v) {
			Print("The ", NoNewline)
			PrintObject(v)
			Print(" swings it out of your reach.", Newline)
		} else {
			Print("The ", NoNewline)
			PrintObject(w)
			Print(" seems white-hot. You can't hold on to it.", Newline)
		}
		return true
	}
	return false
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

func FixBoat() {
	Print("Well done. The boat is repaired.", Newline)
	InflatableBoat.MoveTo(PuncturedBoat.Location())
	RemoveCarefully(&PuncturedBoat)
}


func FixMaintLeak() {
	WaterLevel = -1
	QueueInt(IMaintRoom, false).Run = false
	Print("By some miracle of Zorkian technology, you have managed to stop the leak in the dam.", Newline)
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

func StolenLight() bool {
	oLit := Lit
	Lit = IsLit(Here, true)
	if !Lit && oLit {
		Print("The thief seems to have left you in the dark.", Newline)
	}
	return true
}

func RecoverStiletto() {
	if Stiletto.IsIn(Thief.Location()) {
		Stiletto.Give(FlgNoDesc)
		Stiletto.MoveTo(&Thief)
	}
}

func HackTreasures() {
	RecoverStiletto()
	Thief.Give(FlgInvis)
	for _, x := range TreasureRoom.Children {
		x.Take(FlgInvis)
	}
}

func DepositBooty(rm *Object) bool {
	flg := false
	var toMove []*Object
	for _, x := range Thief.Children {
		if x == &Stiletto || x == &LargeBag {
			continue
		}
		if x.TValue > 0 {
			toMove = append(toMove, x)
			flg = true
			if x == &Egg {
				EggSolve = true
				Egg.Give(FlgOpen)
			}
		}
	}
	for _, x := range toMove {
		x.MoveTo(rm)
	}
	return flg
}

func RobMaze(rm *Object) bool {
	for _, x := range rm.Children {
		if x.Has(FlgTake) && !x.Has(FlgInvis) && Prob(40, false) {
			Print("You hear, off in the distance, someone saying \"My, I wonder what this fine ", NoNewline)
			PrintObject(x)
			Print(" is doing here.\"", Newline)
			if Prob(60, true) {
				x.MoveTo(&Thief)
				x.Give(FlgTouch)
				x.Give(FlgInvis)
			}
			return true
		}
	}
	return false
}

// StealJunk and DropJunk are defined in the IThief section below

func MoveAll(from, to *Object) {
	var toMove []*Object
	for _, x := range from.Children {
		toMove = append(toMove, x)
	}
	for _, x := range toMove {
		x.Take(FlgInvis)
		x.MoveTo(to)
	}
}

func ThiefInTreasure() {
	if len(Here.Children) > 1 {
		Print("The thief gestures mysteriously, and the treasures in the room suddenly vanish.", Newline)
		NewLine()
	}
	for _, f := range Here.Children {
		if f != &Chalice && f != &Thief {
			f.Give(FlgInvis)
		}
	}
}

func Infested(r *Object) bool {
	for _, f := range r.Children {
		if f.Has(FlgActor) && !f.Has(FlgInvis) {
			return true
		}
	}
	return false
}

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

func FightStrength(adjust bool) int {
	if ScoreMax == 0 {
		return StrengthMin
	}
	s := StrengthMin + Score/(ScoreMax/(StrengthMax-StrengthMin))
	if adjust {
		s += Winner.Strength
	}
	return s
}

func FindWeapon(o *Object) *Object {
	for _, w := range o.Children {
		if w == &Stiletto || w == &Axe || w == &Sword || w == &Knife || w == &RustyKnife {
			return w
		}
	}
	return nil
}

func Winning(v *Object) bool {
	vs := v.Strength
	ps := vs - FightStrength(true)
	if ps > 3 {
		return Prob(90, false)
	}
	if ps > 0 {
		return Prob(75, false)
	}
	if ps == 0 {
		return Prob(50, false)
	}
	if vs > 1 {
		return Prob(25, false)
	}
	return Prob(10, false)
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

func WaterFcn(arg ActArg) bool {
	if ActVerb.Norm == "sgive" {
		return false
	}
	if ActVerb.Norm == "through" || ActVerb.Norm == "board" {
		Print(PickOne(SwimYuks), Newline)
		return true
	}
	// Simplified water handling
	if ActVerb.Norm == "take" || ActVerb.Norm == "put" {
		w := DirObj
		if w == &GlobalWater {
			w = &Water
		}
		if ActVerb.Norm == "take" {
			if w.IsIn(&Bottle) && IndirObj == nil {
				Print("It's in the bottle. Perhaps you should take that instead.", Newline)
				return true
			}
			if Bottle.IsIn(Winner) {
				if !Bottle.Has(FlgOpen) {
					Print("The bottle is closed.", Newline)
					ThisIsIt(&Bottle)
					return true
				}
				if !Bottle.HasChildren() {
					Water.MoveTo(&Bottle)
					Print("The bottle is now full of water.", Newline)
					return true
				}
				Print("The water slips through your fingers.", Newline)
				return true
			}
			Print("The water slips through your fingers.", Newline)
			return true
		}
	}
	if ActVerb.Norm == "drop" || ActVerb.Norm == "give" {
		if ActVerb.Norm == "drop" && Water.IsIn(&Bottle) && !Bottle.Has(FlgOpen) {
			Print("The bottle is closed.", Newline)
			return true
		}
		RemoveCarefully(&Water)
		av := Winner.Location()
		if av.Has(FlgVeh) {
			Print("There is now a puddle in the bottom of the ", NoNewline)
			PrintObject(av)
			Print(".", Newline)
			Water.MoveTo(av)
		} else {
			Print("The water spills to the floor and evaporates immediately.", Newline)
		}
		return true
	}
	if ActVerb.Norm == "throw" {
		Print("The water splashes on the walls and evaporates immediately.", Newline)
		RemoveCarefully(&Water)
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

func BoltFcn(arg ActArg) bool {
	if ActVerb.Norm == "turn" {
		if IndirObj == &Wrench {
			if GateFlag {
				ReservoirSouth.Take(FlgTouch)
				if GatesOpen {
					GatesOpen = false
					LoudRoom.Take(FlgTouch)
					Print("The sluice gates close and water starts to collect behind the dam.", Newline)
					Queue(IRfill, 8).Run = true
					QueueInt(IRempty, false).Run = false
				} else {
					GatesOpen = true
					Print("The sluice gates open and water pours through the dam.", Newline)
					Queue(IRempty, 8).Run = true
					QueueInt(IRfill, false).Run = false
				}
			} else {
				Print("The bolt won't turn with your best effort.", Newline)
			}
		} else {
			Print("The bolt won't turn using the ", NoNewline)
			PrintObject(IndirObj)
			Print(".", Newline)
		}
		return true
	}
	if ActVerb.Norm == "take" {
		IntegralPart()
		return true
	}
	if ActVerb.Norm == "oil" {
		Print("Hmm. It appears the tube contained glue, not oil. Turning the bolt won't get any easier....", Newline)
		return true
	}
	return false
}

func BubbleFcn(arg ActArg) bool {
	if ActVerb.Norm == "take" {
		IntegralPart()
		return true
	}
	return false
}

func DamFunction(arg ActArg) bool {
	if ActVerb.Norm == "open" || ActVerb.Norm == "close" {
		Print("Sounds reasonable, but this isn't how.", Newline)
		return true
	}
	if ActVerb.Norm == "plug" {
		if IndirObj == &Hands {
			Print("Are you the little Dutch boy, then? Sorry, this is a big dam.", Newline)
		} else {
			Print("With a ", NoNewline)
			PrintObject(IndirObj)
			Print("? Do you know how big this dam is? You could only stop a tiny leak with that.", Newline)
		}
		return true
	}
	return false
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

func PuncturedBoatFcn(arg ActArg) bool {
	if (ActVerb.Norm == "put" || ActVerb.Norm == "put on") && DirObj == &Putty {
		FixBoat()
		return true
	}
	if ActVerb.Norm == "inflate" || ActVerb.Norm == "fill" {
		Print("No chance. Some moron punctured it.", Newline)
		return true
	}
	if ActVerb.Norm == "plug" {
		if IndirObj == &Putty {
			FixBoat()
			return true
		}
		WithTell(IndirObj)
		return true
	}
	return false
}

func InflatableBoatFcn(arg ActArg) bool {
	if ActVerb.Norm == "inflate" || ActVerb.Norm == "fill" {
		if !InflatableBoat.IsIn(Here) {
			Print("The boat must be on the ground to be inflated.", Newline)
			return true
		}
		if IndirObj == &Pump {
			Print("The boat inflates and appears seaworthy.", Newline)
			if !BoatLabel.Has(FlgTouch) {
				Print("A tan label is lying inside the boat.", Newline)
			}
			Deflate = false
			RemoveCarefully(&InflatableBoat)
			InflatedBoat.MoveTo(Here)
			ThisIsIt(&InflatedBoat)
			return true
		}
		if IndirObj == &Lungs {
			Print("You don't have enough lung power to inflate it.", Newline)
			return true
		}
		Print("With a ", NoNewline)
		PrintObject(IndirObj)
		Print("? Surely you jest!", Newline)
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

func RobberFcn(arg ActArg) bool {
	if ActVerb.Norm == "tell" {
		Print("The thief is a strong, silent type.", Newline)
		Params.Continue = NumUndef
		return true
	}
	if arg == ActArg(FBusy) {
		if Stiletto.IsIn(&Thief) {
			return false
		}
		if Stiletto.IsIn(Thief.Location()) {
			Stiletto.MoveTo(&Thief)
			Stiletto.Give(FlgNoDesc)
			if Thief.IsIn(Here) {
				Print("The robber, somewhat surprised at this turn of events, nimbly retrieves his stiletto.", Newline)
			}
			return true
		}
		return false
	}
	if arg == ActArg(FDead) {
		Stiletto.MoveTo(Here)
		Stiletto.Take(FlgNoDesc)
		x := DepositBooty(Here)
		if Here == &TreasureRoom {
			flg := false
			for _, obj := range Here.Children {
				if obj == &Chalice || obj == &Thief || obj == &Adventurer {
					continue
				}
				obj.Take(FlgInvis)
				if !flg {
					flg = true
					Print("As the thief dies, the power of his magic decreases, and his treasures reappear:", Newline)
				}
				Print("  A ", NoNewline)
				PrintObject(obj)
				if obj.HasChildren() && CanSeeInside(obj) {
					Print(", with ", NoNewline)
					PrintContents(obj)
				}
				NewLine()
			}
			if !flg {
				Print("The chalice is now safe to take.", Newline)
			}
		} else if x {
			Print("His booty remains.", Newline)
		}
		QueueInt(IThief, false).Run = false
		return true
	}
	if arg == ActArg(FFirst) {
		if ThiefHere && !Thief.Has(FlgInvis) && Prob(20, false) {
			Thief.Give(FlgFight)
			Params.Continue = NumUndef
			return true
		}
		return false
	}
	if arg == ActArg(FUnconscious) {
		QueueInt(IThief, false).Run = false
		Thief.Take(FlgFight)
		Stiletto.MoveTo(Here)
		Stiletto.Take(FlgNoDesc)
		Thief.LongDesc = RobberUDesc
		return true
	}
	if arg == ActArg(FConscious) {
		if Thief.Location() == Here {
			Thief.Give(FlgFight)
			Print("The robber revives, briefly feigning continued unconsciousness, and, when he sees his moment, scrambles away from you.", Newline)
		}
		Queue(IThief, -1).Run = true
		Thief.LongDesc = RobberCDesc
		RecoverStiletto()
		return true
	}

	// Default (no special mode)
	if ActVerb.Norm == "hello" && Thief.LongDesc == RobberUDesc {
		Print("The thief, being temporarily incapacitated, is unable to acknowledge your greeting with his usual graciousness.", Newline)
		return true
	}
	if DirObj == &Knife && ActVerb.Norm == "throw" && !Thief.Has(FlgFight) {
		DirObj.MoveTo(Here)
		if Prob(10, false) {
			Print("You evidently frightened the robber, though you didn't hit him. He flees", NoNewline)
			LargeBag.Remove()
			hasStiletto := false
			if Stiletto.IsIn(&Thief) {
				Stiletto.Remove()
				hasStiletto = true
			}
			if Thief.HasChildren() {
				MoveAll(&Thief, Here)
				Print(", but the contents of his bag fall on the floor.", NoNewline)
			} else {
				Print(".", NoNewline)
			}
			LargeBag.MoveTo(&Thief)
			if hasStiletto {
				Stiletto.MoveTo(&Thief)
			}
			NewLine()
			Thief.Give(FlgInvis)
		} else {
			Print("You missed. The thief makes no attempt to take the knife, though it would be a fine addition to the collection in his bag. He does seem angered by your attempt.", Newline)
			Thief.Give(FlgFight)
		}
		return true
	}
	if (ActVerb.Norm == "throw" || ActVerb.Norm == "give") && DirObj != nil && DirObj != &Thief && IndirObj == &Thief {
		if Thief.Strength < 0 {
			Thief.Strength = -Thief.Strength
			Queue(IThief, -1).Run = true
			RecoverStiletto()
			Thief.LongDesc = RobberCDesc
			Print("Your proposed victim suddenly recovers consciousness.", Newline)
		}
		DirObj.MoveTo(&Thief)
		if DirObj.TValue > 0 {
			ThiefEngrossed = true
			Print("The thief is taken aback by your unexpected generosity, but accepts the ", NoNewline)
			PrintObject(DirObj)
			Print(" and stops to admire its beauty.", Newline)
		} else {
			Print("The thief places the ", NoNewline)
			PrintObject(DirObj)
			Print(" in his bag and thanks you politely.", Newline)
		}
		return true
	}
	if ActVerb.Norm == "take" {
		Print("Once you got him, what would you do with him?", Newline)
		return true
	}
	if ActVerb.Norm == "examine" || ActVerb.Norm == "look inside" {
		Print("The thief is a slippery character with beady eyes that flit back and forth. He carries, along with an unmistakable arrogance, a large bag over his shoulder and a vicious stiletto, whose blade is aimed menacingly in your direction. I'd watch out if I were you.", Newline)
		return true
	}
	if ActVerb.Norm == "listen" {
		Print("The thief says nothing, as you have not been formally introduced.", Newline)
		return true
	}
	return false
}

func LargeBagFcn(arg ActArg) bool {
	if ActVerb.Norm == "take" {
		if Thief.LongDesc == RobberUDesc {
			Print("Sadly for you, the robber collapsed on top of the bag. Trying to take it would wake him.", Newline)
		} else {
			Print("The bag will be taken over his dead body.", Newline)
		}
		return true
	}
	if ActVerb.Norm == "put" && IndirObj == &LargeBag {
		Print("It would be a good trick.", Newline)
		return true
	}
	if ActVerb.Norm == "open" || ActVerb.Norm == "close" {
		Print("Getting close enough would be a good trick.", Newline)
		return true
	}
	if ActVerb.Norm == "examine" || ActVerb.Norm == "look inside" {
		Print("The bag is underneath the thief, so one can't say what, if anything, is inside.", Newline)
		return true
	}
	return false
}

func StiletteFcn(arg ActArg) bool {
	return WeaponFunction(&Stiletto, &Thief)
}

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

func RiverFcn(arg ActArg) bool {
	if ActVerb.Norm == "put" && IndirObj == &River {
		if DirObj == &Me {
			JigsUp("You splash around for a while, fighting the current, then you drown.", false)
			return true
		}
		if DirObj == &InflatedBoat {
			Print("You should get in the boat then launch it.", Newline)
			return true
		}
		if DirObj.Has(FlgBurn) {
			RemoveCarefully(DirObj)
			Print("The ", NoNewline)
			PrintObject(DirObj)
			Print(" floats for a moment, then sinks.", Newline)
			return true
		}
		RemoveCarefully(DirObj)
		Print("The ", NoNewline)
		PrintObject(DirObj)
		Print(" splashes into the water and is gone forever.", Newline)
		return true
	}
	if ActVerb.Norm == "leap" || ActVerb.Norm == "through" {
		Print("A look before leaping reveals that the river is wide and dangerous, with swift currents and large, half-hidden rocks. You decide to forgo your swim.", Newline)
		return true
	}
	return false
}

func TreasureInsideFcn(arg ActArg) bool {
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

func DamRoomFcn(arg ActArg) bool {
	if arg == ActLook {
		Print("You are standing on the top of the Flood Control Dam #3, which was quite a tourist attraction in times far distant. There are paths to the north, south, and west, and a scramble down.", Newline)
		if LowTide && GatesOpen {
			Print("The water level behind the dam is low: The sluice gates have been opened. Water rushes through the dam and downstream.", Newline)
		} else if GatesOpen {
			Print("The sluice gates are open, and water rushes through the dam. The water level behind the dam is still high.", Newline)
		} else if LowTide {
			Print("The sluice gates are closed. The water level in the reservoir is quite low, but the level is rising quickly.", Newline)
		} else {
			Print("The sluice gates on the dam are closed. Behind the dam, there can be seen a wide reservoir. Water is pouring over the top of the now abandoned dam.", Newline)
		}
		Print("There is a control panel here, on which a large metal bolt is mounted. Directly above the bolt is a small green plastic bubble", NoNewline)
		if GateFlag {
			Print(" which is glowing serenely", NoNewline)
		}
		Print(".", Newline)
		return true
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

func WhiteCliffsFcn(arg ActArg) bool {
	if arg == ActEnd {
		if InflatedBoat.IsIn(Winner) {
			Deflate = false
		} else {
			Deflate = true
		}
	}
	return false
}

func FallsRoomFcn(arg ActArg) bool {
	if arg == ActLook {
		Print("You are at the top of Aragain Falls, an enormous waterfall with a drop of about 450 feet. The only path here is on the north end.", Newline)
		if RainbowFlag {
			Print("A solid rainbow spans the falls.", Newline)
		} else {
			Print("A beautiful rainbow can be seen over the falls and to the west.", Newline)
		}
		return true
	}
	return false
}

func Rivr4RoomFcn(arg ActArg) bool {
	if arg == ActEnd {
		if Buoy.IsIn(Winner) && BuoyFlag {
			Print("You notice something funny about the feel of the buoy.", Newline)
			BuoyFlag = false
		}
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

func RBoatFcn(arg ActArg) bool {
	if arg == ActEnter || arg == ActEnd || arg == ActLook {
		return false
	}
	if arg == ActBegin {
		if ActVerb.Norm == "walk" {
			if DirObj == ToDirObj("land") || DirObj == ToDirObj("east") || DirObj == ToDirObj("west") {
				return false
			}
			if Here == &Reservoir && (DirObj == ToDirObj("north") || DirObj == ToDirObj("south")) {
				return false
			}
			if Here == &InStream && DirObj == ToDirObj("south") {
				return false
			}
			Print("Read the label for the boat's instructions.", Newline)
			return true
		}
		if ActVerb.Norm == "launch" {
			if Here == &River1 || Here == &River2 || Here == &River3 || Here == &River4 || Here == &Reservoir || Here == &InStream {
				Print("You are on the ", NoNewline)
				if Here == &Reservoir {
					Print("reservoir", NoNewline)
				} else if Here == &InStream {
					Print("stream", NoNewline)
				} else {
					Print("river", NoNewline)
				}
				Print(", or have you forgotten?", Newline)
				return true
			}
			tmp := GoNext(RiverLaunch)
			if tmp == 1 {
				// ZIL: <ENABLE <QUEUE I-RIVER <LKP ,HERE ,RIVER-SPEEDS>>>
				// After GoNext, HERE is the destination room. Use its speed.
				if spd, ok := RiverSpeedMap[Here]; ok {
					Queue(IRiver, spd).Run = true
				}
				return true
			}
			if tmp != 2 {
				Print("You can't launch it here.", Newline)
				return true
			}
			return true
		}
		if (ActVerb.Norm == "drop" && DirObj.Has(FlgWeapon)) ||
			(ActVerb.Norm == "put" && DirObj.Has(FlgWeapon) && IndirObj == &InflatedBoat) ||
			((ActVerb.Norm == "attack" || ActVerb.Norm == "mung") && IndirObj != nil && IndirObj.Has(FlgWeapon)) {
			RemoveCarefully(&InflatedBoat)
			PuncturedBoat.MoveTo(Here)
			Rob(&InflatedBoat, Here, 0)
			Winner.MoveTo(Here)
			Print("It seems that the ", NoNewline)
			if ActVerb.Norm == "drop" || ActVerb.Norm == "put" {
				PrintObject(DirObj)
			} else {
				PrintObject(IndirObj)
			}
			Print(" didn't agree with the boat, as evidenced by the loud hissing noise issuing therefrom. With a pathetic sputter, the boat deflates, leaving you without.", Newline)
			if Here.Has(FlgNonLand) {
				NewLine()
				if Here == &Reservoir || Here == &InStream {
					JigsUp("Another pathetic sputter, this time from you, heralds your drowning.", false)
				} else {
					JigsUp("In other words, fighting the fierce currents of the Frigid River. You manage to hold your own for a bit, but then you are carried over a waterfall and into some nasty rocks. Ouch!", false)
				}
			}
			return true
		}
		return false
	}
	if ActVerb.Norm == "board" {
		if Sceptre.IsIn(Winner) || Knife.IsIn(Winner) || Sword.IsIn(Winner) || RustyKnife.IsIn(Winner) || Axe.IsIn(Winner) || Stiletto.IsIn(Winner) {
			Print("Oops! Something sharp seems to have slipped and punctured the boat. The boat deflates to the sounds of hissing, sputtering, and cursing.", Newline)
			RemoveCarefully(&InflatedBoat)
			PuncturedBoat.MoveTo(Here)
			ThisIsIt(&PuncturedBoat)
			return true
		}
		return false
	}
	if ActVerb.Norm == "inflate" || ActVerb.Norm == "fill" {
		Print("Inflating it further would probably burst it.", Newline)
		return true
	}
	if ActVerb.Norm == "deflate" {
		if Winner.Location() == &InflatedBoat {
			Print("You can't deflate the boat while you're in it.", Newline)
			return true
		}
		if !InflatedBoat.IsIn(Here) {
			Print("The boat must be on the ground to be deflated.", Newline)
			return true
		}
		Print("The boat deflates.", Newline)
		Deflate = true
		RemoveCarefully(&InflatedBoat)
		InflatableBoat.MoveTo(Here)
		ThisIsIt(&InflatableBoat)
		return true
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

func Awaken(o *Object) bool {
	s := o.Strength
	if s < 0 {
		o.Strength = -s
		if o.Action != nil {
			o.Action(ActArg(FConscious))
		}
	}
	return true
}

// ================================================================
// INTERRUPT ROUTINES
// ================================================================

func IFight() bool {
	if Dead {
		return false
	}
	fightQ := false
	numVillains := len(Villains)
	for cnt := 0; cnt < numVillains; cnt++ {
		oo := Villains[cnt]
		o := oo.Villain
		if o.IsIn(Here) && !o.Has(FlgInvis) {
			if o == &Thief && ThiefEngrossed {
				ThiefEngrossed = false
			} else if o.Strength < 0 {
				p := oo.Prob
				if p != 0 && Prob(p, false) {
					oo.Prob = 0
					Awaken(o)
				} else {
					oo.Prob = p + 25
				}
			} else if o.Has(FlgFight) || (o.Action != nil && o.Action(ActArg(FFirst))) {
				fightQ = true
			}
		} else {
			if o.Has(FlgFight) {
				if o.Action != nil {
					o.Action(ActArg(FBusy))
				}
			}
			if o == &Thief {
				ThiefEngrossed = false
			}
			Winner.Take(FlgStagg)
			o.Take(FlgStagg)
			o.Take(FlgFight)
			Awaken(o)
		}
	}
	if !fightQ {
		return false
	}
	return DoFight(numVillains)
}

func ISword() bool {
	if Sword.IsIn(&Adventurer) {
		ng := 0
		if Infested(Here) {
			ng = 2
		} else {
			// Check adjacent rooms for monsters
			dirs := []*DirProps{
				&Here.North, &Here.South, &Here.East, &Here.West,
				&Here.NorthEast, &Here.NorthWest, &Here.SouthEast, &Here.SouthWest,
				&Here.Up, &Here.Down, &Here.Into, &Here.Out,
			}
			for _, dp := range dirs {
				if dp.IsSet() && dp.RExit != nil {
					if Infested(dp.RExit) {
						ng = 1
						break
					}
				}
			}
		}
		g := Sword.TValue
		if ng == g {
			return false
		}
		if ng == 2 {
			Print("Your sword has begun to glow very brightly.", Newline)
		} else if ng == 1 {
			Print("Your sword is glowing with a faint blue glow.", Newline)
		} else if ng == 0 {
			Print("Your sword is no longer glowing.", Newline)
		}
		Sword.TValue = ng
		return true
	}
	// Sword not held - disable the interrupt
	QueueInt(ISword, false).Run = false
	return false
}

func IThief() bool {
	rm := Thief.Location()
	hereQ := !Thief.Has(FlgInvis)
	if hereQ {
		rm = Thief.Location()
	}
	flg := false
	once := false
	for {
		if rm == &TreasureRoom && rm != Here {
			if hereQ {
				HackTreasures()
				hereQ = false
			}
			DepositBooty(&TreasureRoom)
		} else if rm == Here && !Here.Has(FlgOn) && !Troll.IsIn(Here) {
			if ThiefVsAdventurer(hereQ) {
				return true
			}
			if Thief.Has(FlgInvis) {
				hereQ = false
			}
		} else {
			if Thief.IsIn(rm) && !Thief.Has(FlgInvis) {
				// Leave if victim left
				Thief.Give(FlgInvis)
				hereQ = false
			}
			if rm != nil && rm.Has(FlgTouch) {
				Rob(rm, &Thief, 75)
				if rm.Has(FlgMaze) && Here.Has(FlgMaze) {
					flg = RobMaze(rm)
				} else {
					flg = StealJunk(rm)
				}
			}
		}
		if !once && !hereQ {
			once = true
			// Move to next room
			RecoverStiletto()
			found := false
			for _, r := range Rooms.Children {
				if !r.Has(FlgSacred) && r.Has(FlgRLand) {
					Thief.MoveTo(r)
					Thief.Take(FlgFight)
					Thief.Give(FlgInvis)
					ThiefHere = false
					rm = r
					found = true
					break
				}
			}
			if !found {
				break
			}
			continue
		}
		break
	}
	if rm != &TreasureRoom {
		DropJunk(rm)
	}
	return flg
}

func ThiefVsAdventurer(hereQ bool) bool {
	if !Dead && Here == &TreasureRoom {
		return false
	}
	if !ThiefHere {
		if !Dead && !hereQ && Prob(30, false) {
			if Stiletto.IsIn(&Thief) {
				Thief.Take(FlgInvis)
				Print("Someone carrying a large bag is casually leaning against one of the walls here. He does not speak, but it is clear from his aspect that the bag will be taken only over his dead body.", Newline)
				ThiefHere = true
				return true
			}
		}
		if hereQ && Thief.Has(FlgFight) && !Winning(&Thief) {
			Print("Your opponent, determining discretion to be the better part of valor, decides to terminate this little contretemps. With a rueful nod of his head, he steps backward into the gloom and disappears.", Newline)
			Thief.Give(FlgInvis)
			Thief.Take(FlgFight)
			RecoverStiletto()
			return true
		}
	}
	return false
}

// DropJunk - thief drops valueless items from his bag
func DropJunk(rm *Object) bool {
	if rm == nil {
		return false
	}
	flg := false
	for _, x := range Thief.Children {
		if x == &Stiletto || x == &LargeBag {
			continue
		}
		if x.TValue == 0 && Prob(30, true) {
			x.Take(FlgInvis)
			x.MoveTo(rm)
			if !flg && rm == Here {
				Print("The robber, rummaging through his bag, dropped a few items he found valueless.", Newline)
				flg = true
			}
		}
	}
	return flg
}

// StealJunk - thief steals worthless items from a room
func StealJunk(rm *Object) bool {
	if rm == nil {
		return false
	}
	for _, x := range rm.Children {
		if x.TValue == 0 && x.Has(FlgTake) && !x.Has(FlgSacred) && !x.Has(FlgInvis) {
			if x == &Stiletto || Prob(10, true) {
				x.MoveTo(&Thief)
				x.Give(FlgTouch)
				x.Give(FlgInvis)
				if x == &Rope {
					DomeFlag = false
				}
				if rm == Here {
					Print("You suddenly notice that the ", NoNewline)
					PrintObject(x)
					Print(" vanished.", Newline)
					return true
				}
				return false
			}
		}
	}
	return false
}

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

func IRfill() bool {
	Reservoir.Give(FlgNonLand)
	Reservoir.Take(FlgRLand)
	DeepCanyon.Take(FlgTouch)
	LoudRoom.Take(FlgTouch)
	if Trunk.IsIn(&Reservoir) {
		Trunk.Give(FlgInvis)
	}
	LowTide = false
	if Here == &Reservoir {
		if Winner.Location().Has(FlgVeh) {
			Print("The boat lifts gently out of the mud and is now floating on the reservoir.", Newline)
		} else {
			JigsUp("You are lifted up by the rising river! You try to swim, but the currents are too strong. You come closer, closer to the awesome structure of Flood Control Dam #3. The dam beckons to you. The roar of the water nearly deafens you, but you remain conscious as you tumble over the dam toward your certain doom among the rocks at its base.", false)
		}
	} else if Here == &DeepCanyon {
		Print("A sound, like that of flowing water, starts to come from below.", Newline)
	} else if Here == &LoudRoom {
		Print("All of a sudden, an alarmingly loud roaring sound fills the room. Filled with fear, you scramble away.", Newline)
		dest := LoudRuns[rand.Intn(len(LoudRuns))]
		Goto(dest, true)
	} else if Here == &ReservoirNorth || Here == &ReservoirSouth {
		Print("You notice that the water level has risen to the point that it is impossible to cross.", Newline)
	}
	return true
}

func IRempty() bool {
	Reservoir.Give(FlgRLand)
	Reservoir.Take(FlgNonLand)
	DeepCanyon.Take(FlgTouch)
	LoudRoom.Take(FlgTouch)
	Trunk.Take(FlgInvis)
	LowTide = true
	if Here == &Reservoir && Winner.Location().Has(FlgVeh) {
		Print("The water level has dropped to the point at which the boat can no longer stay afloat. It sinks into the mud.", Newline)
	} else if Here == &DeepCanyon {
		Print("The roar of rushing water is quieter now.", Newline)
	} else if Here == &ReservoirNorth || Here == &ReservoirSouth {
		Print("The water level is now quite low here and you could easily cross over to the other side.", Newline)
	}
	return true
}

func IMaintRoom() bool {
	hereQ := Here == &MaintenanceRoom
	if hereQ {
		Print("The water level here is now ", NoNewline)
		idx := WaterLevel / 2
		if idx >= 0 && idx < len(Drownings) {
			Print(Drownings[idx], NoNewline)
		}
		NewLine()
	}
	WaterLevel++
	if WaterLevel >= 14 {
		MungRoom(&MaintenanceRoom, "The room is full of water and cannot be entered.")
		QueueInt(IMaintRoom, false).Run = false
		if hereQ {
			JigsUp("I'm afraid you have done drowned yourself.", false)
		}
	} else if InflatedBoat.IsIn(Winner) && (Here == &MaintenanceRoom || Here == &DamRoom || Here == &DamLobby) {
		JigsUp("The rising water carries the boat over the dam, down the river, and over the falls. Tsk, tsk.", false)
	}
	return true
}

func IRiver() bool {
	if Here != &River1 && Here != &River2 && Here != &River3 && Here != &River4 && Here != &River5 {
		QueueInt(IRiver, false).Run = false
		return false
	}
	rm := Lkp(Here, RiverNext)
	if rm != nil {
		Print("The flow of the river carries you downstream.", Newline)
		NewLine()
		Goto(rm, true)
		// ZIL: <ENABLE <QUEUE I-RIVER <LKP ,HERE ,RIVER-SPEEDS>>>
		// After Goto, HERE is the new room. Use its speed for the next leg.
		if spd, ok := RiverSpeedMap[Here]; ok {
			Queue(IRiver, spd).Run = true
		}
		return true
	}
	JigsUp("Unfortunately, the magic boat doesn't provide protection from the rocks and boulders one meets at the bottom of waterfalls. Including this one.", false)
	return true
}

// Remark prints a combat message, substituting weapon/defender names for markers
func Remark(msg MeleeMsg, defender, weapon *Object) {
	for _, part := range msg {
		switch part.Marker {
		case FWep:
			if weapon != nil {
				PrintObject(weapon)
			}
		case FDef:
			if defender != nil {
				PrintObject(defender)
			}
		default:
			Print(part.Text, NoNewline)
		}
	}
	NewLine()
}

// RandomMeleeMsg picks a random message from a MeleeSet
func RandomMeleeMsg(set MeleeSet) MeleeMsg {
	if len(set) == 0 {
		return nil
	}
	return set[rand.Intn(len(set))]
}

// VillainStrength calculates a villain's effective combat strength
func VillainStrength(oo *VillainEntry) int {
	od := oo.Villain.Strength
	if od >= 0 {
		if oo.Villain == &Thief && ThiefEngrossed {
			if od > 2 {
				od = 2
			}
			ThiefEngrossed = false
		}
		if IndirObj != nil && IndirObj.Has(FlgWeapon) && oo.Best == IndirObj {
			tmp := od - oo.BestAdv
			if tmp < 1 {
				tmp = 1
			}
			od = tmp
		}
	}
	return od
}

// VillainResult applies the combat result to a villain
func VillainResult(villain *Object, def int, res BlowRes) BlowRes {
	villain.Strength = def
	if def == 0 {
		villain.Take(FlgFight)
		Print("Almost as soon as the ", NoNewline)
		PrintObject(villain)
		Print(" breathes his last breath, a cloud of sinister black fog envelops him, and when the fog lifts, the carcass has disappeared.", Newline)
		RemoveCarefully(villain)
		if villain.Action != nil {
			villain.Action(ActArg(FDead))
		}
		return res
	}
	if res == BlowUncon {
		if villain.Action != nil {
			villain.Action(ActArg(FUnconscious))
		}
		return res
	}
	return res
}

// WinnerResult applies the combat result to the player
func WinnerResult(def int, res BlowRes, od int) {
	if def == 0 {
		Winner.Strength = -10000
	} else {
		Winner.Strength = def - od
	}
	if def-od < 0 {
		Queue(ICure, CureWait).Run = true
	}
	if FightStrength(true) <= 0 {
		Winner.Strength = 1 - FightStrength(false)
		JigsUp("It appears that that last blow was too much for you. I'm afraid you are dead.", false)
	}
}

// VillainBlow processes a villain's attack on the player
func VillainBlow(oo *VillainEntry, out bool) BlowRes {
	villain := oo.Villain
	remarks := oo.Msgs

	Winner.Take(FlgStagg)
	if villain.Has(FlgStagg) {
		Print("The ", NoNewline)
		PrintObject(villain)
		Print(" slowly regains his feet.", Newline)
		villain.Take(FlgStagg)
		return BlowMissed
	}
	att := VillainStrength(oo)
	oa := att
	def := FightStrength(true)
	if def <= 0 {
		return BlowMissed
	}
	od := FightStrength(false)
	dweapon := FindWeapon(Winner)
	var res BlowRes
	if def < 0 {
		res = BlowKill
	} else {
		if def == 1 {
			if att > 2 {
				att = 3
			}
			idx := att - 1
			if idx < 0 {
				idx = 0
			}
			if idx >= len(Def1Res) {
				idx = len(Def1Res) - 1
			}
			tbl := Def1Res[idx]
			_ = tbl
			// Lookup from the actual Def arrays
			switch {
			case att <= 1:
				res = Def1[rand.Intn(len(Def1))]
			case att == 2:
				if int(Def1Res[1]) < len(Def1) {
					start := int(Def1Res[1])
					_ = start
				}
				res = Def1[rand.Intn(len(Def1))]
			default:
				res = Def1[rand.Intn(len(Def1))]
			}
		} else if def == 2 {
			if att > 3 {
				att = 4
			}
			switch {
			case att <= 1:
				res = Def2A[rand.Intn(len(Def2A))]
			default:
				res = Def2B[rand.Intn(len(Def2B))]
			}
		} else {
			att = att - def
			if att < -1 {
				att = -2
			}
			if att > 1 {
				att = 2
			}
			idx := att + 2
			if idx < 0 {
				idx = 0
			}
			if idx > 4 {
				idx = 4
			}
			switch {
			case idx <= 1:
				res = Def3A[rand.Intn(len(Def3A))]
			case idx == 2:
				res = Def3B[rand.Intn(len(Def3B))]
			default:
				res = Def3C[rand.Intn(len(Def3C))]
			}
		}
		if out {
			if res == BlowStag {
				res = BlowHesitate
			} else {
				res = BlowSitDuck
			}
		}
		if res == BlowStag && dweapon != nil && Prob(25, true) {
			res = BlowLoseWpn
		}
		msg := RandomMeleeMsg((*remarks)[res-1])
		if msg != nil {
			Remark(msg, Winner, dweapon)
		}
	}
	_ = oa
	// Apply results
	switch {
	case res == BlowMissed || res == BlowHesitate:
		// Nothing
	case res == BlowUncon:
		// Nothing extra
	case res == BlowKill || res == BlowSitDuck:
		def = 0
	case res == BlowLightWnd:
		def--
		if def < 0 {
			def = 0
		}
		if LoadAllowed > 50 {
			LoadAllowed -= 10
		}
	case res == BlowHeavyWnd:
		def -= 2
		if def < 0 {
			def = 0
		}
		if LoadAllowed > 50 {
			LoadAllowed -= 20
		}
	case res == BlowStag:
		Winner.Give(FlgStagg)
	default:
		// BlowLoseWpn
		if dweapon != nil {
			dweapon.MoveTo(Here)
			nweapon := FindWeapon(Winner)
			if nweapon != nil {
				Print("Fortunately, you still have a ", NoNewline)
				PrintObject(nweapon)
				Print(".", Newline)
			}
		}
	}
	WinnerResult(def, res, od)
	return res
}

// DoFight runs the villain combat loop
func DoFight(numVillains int) bool {
	for {
		cnt := 0
		var res BlowRes
		out := false
		for cnt < numVillains {
			oo := Villains[cnt]
			cnt++
			o := oo.Villain
			if !o.Has(FlgFight) {
				continue
			}
			if o.Action != nil && o.Action(ActArg(FBusy)) {
				continue
			}
			res = VillainBlow(oo, out)
			if res == BlowUnk {
				return false
			}
			if res == BlowUncon {
				out = true
			}
		}
		if res != BlowUnk {
			if !out {
				return true
			}
			// unconscious rounds
			if out {
				return true
			}
		}
		return true
	}
}

func HeroBlow() bool {
	if DirObj == nil {
		return false
	}
	// Find the villain entry for the target
	var oo *VillainEntry
	for _, v := range Villains {
		if v.Villain == DirObj {
			oo = v
			break
		}
	}
	DirObj.Give(FlgFight)
	if Winner.Has(FlgStagg) {
		Print("You are still recovering from that last blow, so your attack is ineffective.", Newline)
		Winner.Take(FlgStagg)
		return true
	}
	att := FightStrength(true)
	if att < 1 {
		att = 1
	}
	oa := att
	_ = oa
	villain := DirObj
	if oo == nil {
		return false
	}
	od := VillainStrength(oo)
	def := od
	if def == 0 {
		if DirObj == Winner {
			JigsUp("Well, you really did it that time. Is suicide painless?", false)
			return true
		}
		Print("Attacking the ", NoNewline)
		PrintObject(villain)
		Print(" is pointless.", Newline)
		return true
	}
	dweapon := FindWeapon(villain)
	var res BlowRes
	if dweapon == nil || def < 0 {
		Print("The ", NoNewline)
		if def < 0 {
			Print("unconscious", NoNewline)
		} else {
			Print("unarmed", NoNewline)
		}
		Print(" ", NoNewline)
		PrintObject(villain)
		Print(" cannot defend himself: He dies.", Newline)
		res = BlowKill
	} else {
		if def == 1 {
			if att > 2 {
				att = 3
			}
			res = Def1[rand.Intn(len(Def1))]
		} else if def == 2 {
			if att > 3 {
				att = 4
			}
			res = Def2B[rand.Intn(len(Def2B))]
		} else {
			att = att - def
			if att < -1 {
				att = -2
			}
			if att > 1 {
				att = 2
			}
			idx := att + 2
			if idx < 0 {
				idx = 0
			}
			switch {
			case idx <= 1:
				res = Def3A[rand.Intn(len(Def3A))]
			case idx == 2:
				res = Def3B[rand.Intn(len(Def3B))]
			default:
				res = Def3C[rand.Intn(len(Def3C))]
			}
		}
		if res == BlowStag && dweapon != nil && Prob(25, false) {
			res = BlowLoseWpn
		}
		msg := RandomMeleeMsg(HeroMelee[res-1])
		if msg != nil {
			Remark(msg, DirObj, IndirObj)
		}
	}
	// Apply hero blow results
	switch {
	case res == BlowMissed || res == BlowHesitate:
		// Nothing
	case res == BlowUncon:
		def = -def
	case res == BlowKill || res == BlowSitDuck:
		def = 0
	case res == BlowLightWnd:
		def--
		if def < 0 {
			def = 0
		}
	case res == BlowHeavyWnd:
		def -= 2
		if def < 0 {
			def = 0
		}
	case res == BlowStag:
		DirObj.Give(FlgStagg)
	default:
		// BlowLoseWpn
		if dweapon != nil {
			dweapon.Take(FlgNoDesc)
			dweapon.Give(FlgWeapon)
			dweapon.MoveTo(Here)
			ThisIsIt(dweapon)
		}
	}
	VillainResult(DirObj, def, res)
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
