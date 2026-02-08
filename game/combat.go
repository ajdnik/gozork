package game

import . "github.com/ajdnik/gozork/engine"

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

)

// ================================================================
// HELPER FUNCTIONS
// ================================================================



func WeaponFunction(w, v *Object) bool {
	if !v.IsIn(G.Here) {
		return false
	}
	if G.ActVerb.Norm == "take" {
		if w.IsIn(v) {
			Printf("The %s swings it out of your reach.\n", v.Desc)
		} else {
			Printf("The %s seems white-hot. You can't hold on to it.\n", w.Desc)
		}
		return true
	}
	return false
}

func FightStrength(adjust bool) int {
	if ScoreMax == 0 {
		return StrengthMin
	}
	s := StrengthMin + G.Score/(ScoreMax/(StrengthMax-StrengthMin))
	if adjust {
		s += G.Winner.GetStrength()
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
	vs := v.GetStrength()
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

func Awaken(o *Object) bool {
	s := o.GetStrength()
	if s < 0 {
		o.SetStrength(-s)
		if o.Action != nil {
			o.Action(ActArg(FConscious))
		}
	}
	return true
}

func IFight() bool {
	if GD().Dead {
		return false
	}
	fightQ := false
	numVillains := len(GD().Villains)
	for cnt := 0; cnt < numVillains; cnt++ {
		oo := GD().Villains[cnt]
		o := oo.Villain
		if o.IsIn(G.Here) && !o.Has(FlgInvis) {
			if o == &Thief && GD().ThiefEngrossed {
				GD().ThiefEngrossed = false
			} else if o.GetStrength() < 0 {
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
				GD().ThiefEngrossed = false
			}
			G.Winner.Take(FlgStagg)
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
		if Infested(G.Here) {
			ng = 2
		} else {
			// Check adjacent rooms for monsters
			for _, d := range AllDirections {
				dp := G.Here.GetExit(d)
				if dp != nil && dp.RExit != nil {
					if Infested(dp.RExit) {
						ng = 1
						break
					}
				}
			}
		}
		g := Sword.GetTValue()
		if ng == g {
			return false
		}
		if ng == 2 {
			Printf("Your sword has begun to glow very brightly.\n")
		} else if ng == 1 {
			Printf("Your sword is glowing with a faint blue glow.\n")
		} else if ng == 0 {
			Printf("Your sword is no longer glowing.\n")
		}
		Sword.SetTValue(ng)
		return true
	}
	// Sword not held - disable the interrupt
	QueueInt("ISword", false).Run = false
	return false
}

// Remark prints a combat message, substituting weapon/defender names for markers
func Remark(msg MeleeMsg, defender, weapon *Object) {
	for _, part := range msg {
		switch part.Marker {
		case FWep:
			if weapon != nil {
				Printf("%s", weapon.Desc)
			}
		case FDef:
			if defender != nil {
				Printf("%s", defender.Desc)
			}
		default:
			Printf("%s", part.Text)
		}
	}
	Printf("\n")
}

// RandomMeleeMsg picks a random message from a MeleeSet
func RandomMeleeMsg(set MeleeSet) MeleeMsg {
	if len(set) == 0 {
		return nil
	}
	return set[G.Rand.Intn(len(set))]
}

// VillainStrength calculates a villain's effective combat strength
func VillainStrength(oo *VillainEntry) int {
	od := oo.Villain.GetStrength()
	if od >= 0 {
		if oo.Villain == &Thief && GD().ThiefEngrossed {
			if od > 2 {
				od = 2
			}
			GD().ThiefEngrossed = false
		}
		if G.IndirObj != nil && G.IndirObj.Has(FlgWeapon) && oo.Best == G.IndirObj {
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
	villain.SetStrength(def)
	if def == 0 {
		villain.Take(FlgFight)
		Printf("Almost as soon as the %s breathes his last breath, a cloud of sinister black fog envelops him, and when the fog lifts, the carcass has disappeared.\n", villain.Desc)
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
		G.Winner.SetStrength(-10000)
	} else {
		G.Winner.SetStrength(def - od)
	}
	if def-od < 0 {
		Queue("ICure", CureWait).Run = true
	}
	if FightStrength(true) <= 0 {
		G.Winner.SetStrength(1 - FightStrength(false))
		JigsUp("It appears that that last blow was too much for you. I'm afraid you are dead.", false)
	}
}

// VillainBlow processes a villain's attack on the player
func VillainBlow(oo *VillainEntry, out bool) BlowRes {
	villain := oo.Villain
	remarks := oo.Msgs

	G.Winner.Take(FlgStagg)
	if villain.Has(FlgStagg) {
		Printf("The %s slowly regains his feet.\n", villain.Desc)
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
	dweapon := FindWeapon(G.Winner)
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
				res = Def1[G.Rand.Intn(len(Def1))]
			case att == 2:
				if int(Def1Res[1]) < len(Def1) {
					start := int(Def1Res[1])
					_ = start
				}
				res = Def1[G.Rand.Intn(len(Def1))]
			default:
				res = Def1[G.Rand.Intn(len(Def1))]
			}
		} else if def == 2 {
			if att > 3 {
				att = 4
			}
			switch {
			case att <= 1:
				res = Def2A[G.Rand.Intn(len(Def2A))]
			default:
				res = Def2B[G.Rand.Intn(len(Def2B))]
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
				res = Def3A[G.Rand.Intn(len(Def3A))]
			case idx == 2:
				res = Def3B[G.Rand.Intn(len(Def3B))]
			default:
				res = Def3C[G.Rand.Intn(len(Def3C))]
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
			Remark(msg, G.Winner, dweapon)
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
		if GD().LoadAllowed > 50 {
			GD().LoadAllowed -= 10
		}
	case res == BlowHeavyWnd:
		def -= 2
		if def < 0 {
			def = 0
		}
		if GD().LoadAllowed > 50 {
			GD().LoadAllowed -= 20
		}
	case res == BlowStag:
		G.Winner.Give(FlgStagg)
	default:
		// BlowLoseWpn
		if dweapon != nil {
			dweapon.MoveTo(G.Here)
			nweapon := FindWeapon(G.Winner)
			if nweapon != nil {
				Printf("Fortunately, you still have a %s.\n", nweapon.Desc)
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
			oo := GD().Villains[cnt]
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
	if G.DirObj == nil {
		return false
	}
	// Find the villain entry for the target
	var oo *VillainEntry
	for _, v := range GD().Villains {
		if v.Villain == G.DirObj {
			oo = v
			break
		}
	}
	G.DirObj.Give(FlgFight)
	if G.Winner.Has(FlgStagg) {
		Printf("You are still recovering from that last blow, so your attack is ineffective.\n")
		G.Winner.Take(FlgStagg)
		return true
	}
	att := FightStrength(true)
	if att < 1 {
		att = 1
	}
	oa := att
	_ = oa
	villain := G.DirObj
	if oo == nil {
		return false
	}
	od := VillainStrength(oo)
	def := od
	if def == 0 {
		if G.DirObj == G.Winner {
			JigsUp("Well, you really did it that time. Is suicide painless?", false)
			return true
		}
		Printf("Attacking the %s is pointless.\n", villain.Desc)
		return true
	}
	dweapon := FindWeapon(villain)
	var res BlowRes
	if dweapon == nil || def < 0 {
		Printf("The ")
		if def < 0 {
			Printf("unconscious")
		} else {
			Printf("unarmed")
		}
		Printf(" %s cannot defend himself: He dies.\n", villain.Desc)
		res = BlowKill
	} else {
		if def == 1 {
			if att > 2 {
				att = 3
			}
			res = Def1[G.Rand.Intn(len(Def1))]
		} else if def == 2 {
			if att > 3 {
				att = 4
			}
			res = Def2B[G.Rand.Intn(len(Def2B))]
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
				res = Def3A[G.Rand.Intn(len(Def3A))]
			case idx == 2:
				res = Def3B[G.Rand.Intn(len(Def3B))]
			default:
				res = Def3C[G.Rand.Intn(len(Def3C))]
			}
		}
		if res == BlowStag && dweapon != nil && Prob(25, false) {
			res = BlowLoseWpn
		}
		msg := RandomMeleeMsg(HeroMelee[res-1])
		if msg != nil {
			Remark(msg, G.DirObj, G.IndirObj)
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
		G.DirObj.Give(FlgStagg)
	default:
		// BlowLoseWpn
		if dweapon != nil {
			dweapon.Take(FlgNoDesc)
			dweapon.Give(FlgWeapon)
			dweapon.MoveTo(G.Here)
			ThisIsIt(dweapon)
		}
	}
	VillainResult(G.DirObj, def, res)
	return true
}
