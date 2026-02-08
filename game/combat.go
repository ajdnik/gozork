package game

import . "github.com/ajdnik/gozork/engine"

// ================================================================
// COMBAT SYSTEM TYPES AND TABLES
// ================================================================

type BlowRes int

const (
	blowUnk BlowRes = iota
	blowMissed
	blowUncon
	blowKill
	blowLightWnd
	blowHeavyWnd
	blowStag
	blowLoseWpn
	blowHesitate
	blowSitDuck
)

// Combat mode constants are defined in engine as ActionArg values:
// ActBusy, ActDead, ActUnconscious, ActConscious, ActFirst.

// Combat strength constants
const (
	strengthMax = 7
	strengthMin = 2
	cureWait    = 30
)

// Melee message marker constants
const (
	fWep = 0 // means print weapon name
	fDef = 1 // means print defender name
)

// MeleePart is a fragment of a melee message (either a string or fWep/fDef marker)
type MeleePart struct {
	Text   string
	Marker int // -1 = normal text, fWep = weapon, fDef = defender
}

func mp(s string) MeleePart { return MeleePart{Text: s, Marker: -1} }
func mw() MeleePart         { return MeleePart{Marker: fWep} }
func md() MeleePart         { return MeleePart{Marker: fDef} }

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
	def1    = [13]BlowRes{blowMissed, blowMissed, blowMissed, blowMissed, blowStag, blowStag, blowUncon, blowUncon, blowKill, blowKill, blowKill, blowKill, blowKill}
	def2A   = [10]BlowRes{blowMissed, blowMissed, blowMissed, blowMissed, blowMissed, blowStag, blowStag, blowLightWnd, blowLightWnd, blowUncon}
	def2B   = [12]BlowRes{blowMissed, blowMissed, blowMissed, blowStag, blowStag, blowLightWnd, blowLightWnd, blowLightWnd, blowUncon, blowKill, blowKill, blowKill}
	def3A   = [11]BlowRes{blowMissed, blowMissed, blowMissed, blowMissed, blowMissed, blowStag, blowStag, blowLightWnd, blowLightWnd, blowHeavyWnd, blowHeavyWnd}
	def3B   = [11]BlowRes{blowMissed, blowMissed, blowMissed, blowStag, blowStag, blowLightWnd, blowLightWnd, blowLightWnd, blowHeavyWnd, blowHeavyWnd, blowHeavyWnd}
	def3C   = [10]BlowRes{blowMissed, blowStag, blowStag, blowLightWnd, blowLightWnd, blowLightWnd, blowLightWnd, blowHeavyWnd, blowHeavyWnd, blowHeavyWnd}
	def1Res = [4]BlowRes{def1[0], blowUnk, blowUnk}
	def2Res = [4]BlowRes{def2A[0], def2B[0], blowUnk, blowUnk}
	def3Res = [5]BlowRes{def3A[0], blowUnk, def3B[0], blowUnk, def3C[0]}

	// Hero melee messages (indexed by outcome: missed, unconscious, killed, light-wound, heavy-wound, stagger, lose-weapon)
	heroMelee = MeleeTable{
		// blowMissed
		{
			{mp("Your "), mw(), mp(" misses the "), md(), mp(" by an inch.")},
			{mp("A good slash, but it misses the "), md(), mp(" by a mile.")},
			{mp("You charge, but the "), md(), mp(" jumps nimbly aside.")},
			{mp("Clang! Crash! The "), md(), mp(" parries.")},
			{mp("A quick stroke, but the "), md(), mp(" is on guard.")},
			{mp("A good stroke, but it's too slow; the "), md(), mp(" dodges.")},
		},
		// blowUncon
		{
			{mp("Your "), mw(), mp(" crashes down, knocking the "), md(), mp(" into dreamland.")},
			{mp("The "), md(), mp(" is battered into unconsciousness.")},
			{mp("A furious exchange, and the "), md(), mp(" is knocked out!")},
			{mp("The haft of your "), mw(), mp(" knocks out the "), md(), mp(".")},
			{mp("The "), md(), mp(" is knocked out!")},
		},
		// blowKill
		{
			{mp("it's curtains for the "), md(), mp(" as your "), mw(), mp(" removes his head.")},
			{mp("The fatal blow strikes the "), md(), mp(" square in the heart: He dies.")},
			{mp("The "), md(), mp(" takes a fatal blow and slumps to the floor dead.")},
		},
		// blowLightWnd
		{
			{mp("The "), md(), mp(" is struck on the arm; blood begins to trickle down.")},
			{mp("Your "), mw(), mp(" pinks the "), md(), mp(" on the wrist, but it's not serious.")},
			{mp("Your stroke lands, but it was only the flat of the blade.")},
			{mp("The blow lands, making a shallow gash in the "), md(), mp("'s arm!")},
		},
		// blowHeavyWnd
		{
			{mp("The "), md(), mp(" receives a deep gash in his side.")},
			{mp("A savage blow on the thigh! The "), md(), mp(" is stunned but can still fight!")},
			{mp("Slash! Your blow lands! That one hit an artery, it could be serious!")},
			{mp("Slash! Your stroke connects! This could be serious!")},
		},
		// blowStag
		{
			{mp("The "), md(), mp(" is staggered, and drops to his knees.")},
			{mp("The "), md(), mp(" is momentarily disoriented and can't fight back.")},
			{mp("The force of your blow knocks the "), md(), mp(" back, stunned.")},
			{mp("The "), md(), mp(" is confused and can't fight back.")},
			{mp("The quickness of your thrust knocks the "), md(), mp(" back, stunned.")},
		},
		// blowLoseWpn
		{
			{mp("The "), md(), mp("'s weapon is knocked to the floor, leaving him unarmed.")},
			{mp("The "), md(), mp(" is disarmed by a subtle feint past his guard.")},
		},
		// blowHesitate (not used for hero)
		{},
		// blowSitDuck (not used for hero)
		{},
	}

	// cyclops melee messages
	cyclopsMelee = MeleeTable{
		// blowMissed
		{
			{mp("The cyclops misses, but the backwash almost knocks you over.")},
			{mp("The cyclops rushes you, but runs into the wall.")},
		},
		// blowUncon
		{
			{mp("The cyclops sends you crashing to the floor, unconscious.")},
		},
		// blowKill
		{
			{mp("The cyclops breaks your neck with a massive smash.")},
		},
		// blowLightWnd
		{
			{mp("A quick punch, but it was only a glancing blow.")},
			{mp("A glancing blow from the cyclops' fist.")},
		},
		// blowHeavyWnd
		{
			{mp("The monster smashes his huge fist into your chest, breaking several ribs.")},
			{mp("The cyclops almost knocks the wind out of you with a quick punch.")},
		},
		// blowStag
		{
			{mp("The cyclops lands a punch that knocks the wind out of you.")},
			{mp("Heedless of your weapons, the cyclops tosses you against the rock wall of the room.")},
		},
		// blowLoseWpn
		{
			{mp("The cyclops grabs your "), mw(), mp(", tastes it, and throws it to the ground in disgust.")},
			{mp("The monster grabs you on the wrist, squeezes, and you drop your "), mw(), mp(" in pain.")},
		},
		// blowHesitate
		{
			{mp("The cyclops seems unable to decide whether to broil or stew his dinner.")},
		},
		// blowSitDuck
		{
			{mp("The cyclops, no sportsman, dispatches his unconscious victim.")},
		},
	}

	// troll melee messages
	trollMelee = MeleeTable{
		// blowMissed
		{
			{mp("The troll swings his axe, but it misses.")},
			{mp("The troll's axe barely misses your ear.")},
			{mp("The axe sweeps past as you jump aside.")},
			{mp("The axe crashes against the rock, throwing sparks!")},
		},
		// blowUncon
		{
			{mp("The flat of the troll's axe hits you delicately on the head, knocking you out.")},
		},
		// blowKill
		{
			{mp("The troll neatly removes your head.")},
			{mp("The troll's axe stroke cleaves you from the nave to the chops.")},
			{mp("The troll's axe removes your head.")},
		},
		// blowLightWnd
		{
			{mp("The axe gets you right in the side. Ouch!")},
			{mp("The flat of the troll's axe skins across your forearm.")},
			{mp("The troll's swing almost knocks you over as you barely parry in time.")},
			{mp("The troll swings his axe, and it nicks your arm as you dodge.")},
		},
		// blowHeavyWnd
		{
			{mp("The troll charges, and his axe slashes you on your "), mw(), mp(" arm.")},
			{mp("An axe stroke makes a deep wound in your leg.")},
			{mp("The troll's axe swings down, gashing your shoulder.")},
		},
		// blowStag
		{
			{mp("The troll hits you with a glancing blow, and you are momentarily stunned.")},
			{mp("The troll swings; the blade turns on your armor but crashes broadside into your head.")},
			{mp("You stagger back under a hail of axe strokes.")},
			{mp("The troll's mighty blow drops you to your knees.")},
		},
		// blowLoseWpn
		{
			{mp("The axe hits your "), mw(), mp(" and knocks it spinning.")},
			{mp("The troll swings, you parry, but the force of his blow knocks your "), mw(), mp(" away.")},
			{mp("The axe knocks your "), mw(), mp(" out of your hand. it falls to the floor.")},
		},
		// blowHesitate
		{
			{mp("The troll hesitates, fingering his axe.")},
			{mp("The troll scratches his head ruminatively:  Might you be magically protected, he wonders?")},
		},
		// blowSitDuck
		{
			{mp("Conquering his fears, the troll puts you to death.")},
		},
	}

	// thief melee messages
	thiefMelee = MeleeTable{
		// blowMissed
		{
			{mp("The thief stabs nonchalantly with his stiletto and misses.")},
			{mp("You dodge as the thief comes in low.")},
			{mp("You parry a lightning thrust, and the thief salutes you with a grim nod.")},
			{mp("The thief tries to sneak past your guard, but you twist away.")},
		},
		// blowUncon
		{
			{mp("Shifting in the midst of a thrust, the thief knocks you unconscious with the haft of his stiletto.")},
			{mp("The thief knocks you out.")},
		},
		// blowKill
		{
			{mp("Finishing you off, the thief inserts his blade into your heart.")},
			{mp("The thief comes in from the side, feints, and inserts the blade into your ribs.")},
			{mp("The thief bows formally, raises his stiletto, and with a wry grin, ends the battle and your life.")},
		},
		// blowLightWnd
		{
			{mp("A quick thrust pinks your left arm, and blood starts to trickle down.")},
			{mp("The thief draws blood, raking his stiletto across your arm.")},
			{mp("The stiletto flashes faster than you can follow, and blood wells from your leg.")},
			{mp("The thief slowly approaches, strikes like a snake, and leaves you wounded.")},
		},
		// blowHeavyWnd
		{
			{mp("The thief strikes like a snake! The resulting wound is serious.")},
			{mp("The thief stabs a deep cut in your upper arm.")},
			{mp("The stiletto touches your forehead, and the blood obscures your vision.")},
			{mp("The thief strikes at your wrist, and suddenly your grip is slippery with blood.")},
		},
		// blowStag
		{
			{mp("The butt of his stiletto cracks you on the skull, and you stagger back.")},
			{mp("The thief rams the haft of his blade into your stomach, leaving you out of breath.")},
			{mp("The thief attacks, and you fall back desperately.")},
		},
		// blowLoseWpn
		{
			{mp("A long, theatrical slash. You catch it on your "), mw(), mp(", but the thief twists his knife, and the "), mw(), mp(" goes flying.")},
			{mp("The thief neatly flips your "), mw(), mp(" out of your hands, and it drops to the floor.")},
			{mp("You parry a low thrust, and your "), mw(), mp(" slips out of your hand.")},
		},
		// blowHesitate
		{
			{mp("The thief, a man of superior breeding, pauses for a moment to consider the propriety of finishing you off.")},
			{mp("The thief amuses himself by searching your pockets.")},
			{mp("The thief entertains himself by rifling your pack.")},
		},
		// blowSitDuck
		{
			{mp("The thief, forgetting his essentially genteel upbringing, cuts your throat.")},
			{mp("The thief, a pragmatist, dispatches you as a threat to his livelihood.")},
		},
	}
)

// ================================================================
// HELPER FUNCTIONS
// ================================================================

func weaponFunction(w, v *Object) bool {
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

func fightStrength(adjust bool) int {
	if ScoreMax == 0 {
		return strengthMin
	}
	s := strengthMin + G.Score/(ScoreMax/(strengthMax-strengthMin))
	if adjust {
		s += G.Winner.GetStrength()
	}
	return s
}

func findWeapon(o *Object) *Object {
	for _, w := range o.Children {
		if w == &stiletto || w == &axe || w == &sword || w == &knife || w == &rustyKnife {
			return w
		}
	}
	return nil
}

func winning(v *Object) bool {
	vs := v.GetStrength()
	ps := vs - fightStrength(true)
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

func awaken(o *Object) bool {
	s := o.GetStrength()
	if s < 0 {
		o.SetStrength(-s)
		if o.Action != nil {
			o.Action(ActConscious)
		}
	}
	return true
}

func iFight() bool {
	if gD().Dead {
		return false
	}
	fightQ := false
	numVillains := len(gD().Villains)
	for cnt := 0; cnt < numVillains; cnt++ {
		oo := gD().Villains[cnt]
		o := oo.Villain
		if o.IsIn(G.Here) && !o.Has(FlgInvis) {
			if o == &thief && gD().ThiefEngrossed {
				gD().ThiefEngrossed = false
			} else if o.GetStrength() < 0 {
				p := oo.Prob
				if p != 0 && Prob(p, false) {
					oo.Prob = 0
					awaken(o)
				} else {
					oo.Prob = p + 25
				}
			} else if o.Has(FlgFight) || (o.Action != nil && o.Action(ActFirst)) {
				fightQ = true
			}
		} else {
			if o.Has(FlgFight) {
				if o.Action != nil {
					o.Action(ActBusy)
				}
			}
			if o == &thief {
				gD().ThiefEngrossed = false
			}
			G.Winner.Take(FlgStaggered)
			o.Take(FlgStaggered)
			o.Take(FlgFight)
			awaken(o)
		}
	}
	if !fightQ {
		return false
	}
	return doFight(numVillains)
}

func iSword() bool {
	if sword.IsIn(&adventurer) {
		ng := 0
		if infested(G.Here) {
			ng = 2
		} else {
			// Check adjacent rooms for monsters
			for _, d := range AllDirections {
				dp, ok := G.Here.GetExit(d)
				if ok && dp.RExit != nil {
					if infested(dp.RExit) {
						ng = 1
						break
					}
				}
			}
		}
		g := sword.GetTValue()
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
		sword.SetTValue(ng)
		return true
	}
	// sword not held - disable the interrupt
	QueueInt("iSword", false).Run = false
	return false
}

// remark prints a combat message, substituting weapon/defender names for markers
func remark(msg MeleeMsg, defender, weapon *Object) {
	for _, part := range msg {
		switch part.Marker {
		case fWep:
			if weapon != nil {
				Printf("%s", weapon.Desc)
			}
		case fDef:
			if defender != nil {
				Printf("%s", defender.Desc)
			}
		default:
			Printf("%s", part.Text)
		}
	}
	Printf("\n")
}

// randomMeleeMsg picks a random message from a MeleeSet
func randomMeleeMsg(set MeleeSet) MeleeMsg {
	if len(set) == 0 {
		return nil
	}
	return set[G.Rand.Intn(len(set))]
}

// villainStrength calculates a villain's effective combat strength
func villainStrength(oo *VillainEntry) int {
	od := oo.Villain.GetStrength()
	if od >= 0 {
		if oo.Villain == &thief && gD().ThiefEngrossed {
			if od > 2 {
				od = 2
			}
			gD().ThiefEngrossed = false
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

// villainResult applies the combat result to a villain
func villainResult(villain *Object, def int, res BlowRes) BlowRes {
	villain.SetStrength(def)
	if def == 0 {
		villain.Take(FlgFight)
		Printf("Almost as soon as the %s breathes his last breath, a cloud of sinister black fog envelops him, and when the fog lifts, the carcass has disappeared.\n", villain.Desc)
		removeCarefully(villain)
		if villain.Action != nil {
			villain.Action(ActDead)
		}
		return res
	}
	if res == blowUncon {
		if villain.Action != nil {
			villain.Action(ActUnconscious)
		}
		return res
	}
	return res
}

// winnerResult applies the combat result to the player
func winnerResult(def int, res BlowRes, od int) {
	if def == 0 {
		G.Winner.SetStrength(-10000)
	} else {
		G.Winner.SetStrength(def - od)
	}
	if def-od < 0 {
		Queue("iCure", cureWait).Run = true
	}
	if fightStrength(true) <= 0 {
		G.Winner.SetStrength(1 - fightStrength(false))
		jigsUp("it appears that that last blow was too much for you. I'm afraid you are dead.", false)
	}
}

// villainBlow processes a villain's attack on the player
func villainBlow(oo *VillainEntry, out bool) BlowRes {
	villain := oo.Villain
	remarks := oo.Msgs

	G.Winner.Take(FlgStaggered)
	if villain.Has(FlgStaggered) {
		Printf("The %s slowly regains his feet.\n", villain.Desc)
		villain.Take(FlgStaggered)
		return blowMissed
	}
	att := villainStrength(oo)
	oa := att
	def := fightStrength(true)
	if def <= 0 {
		return blowMissed
	}
	od := fightStrength(false)
	dweapon := findWeapon(G.Winner)
	var res BlowRes
	if def < 0 {
		res = blowKill
	} else {
		if def == 1 {
			if att > 2 {
				att = 3
			}
			idx := att - 1
			if idx < 0 {
				idx = 0
			}
			if idx >= len(def1Res) {
				idx = len(def1Res) - 1
			}
			tbl := def1Res[idx]
			_ = tbl
			// Lookup from the actual Def arrays
			switch {
			case att <= 1:
				res = def1[G.Rand.Intn(len(def1))]
			case att == 2:
				if int(def1Res[1]) < len(def1) {
					start := int(def1Res[1])
					_ = start
				}
				res = def1[G.Rand.Intn(len(def1))]
			default:
				res = def1[G.Rand.Intn(len(def1))]
			}
		} else if def == 2 {
			if att > 3 {
				att = 4
			}
			switch {
			case att <= 1:
				res = def2A[G.Rand.Intn(len(def2A))]
			default:
				res = def2B[G.Rand.Intn(len(def2B))]
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
				res = def3A[G.Rand.Intn(len(def3A))]
			case idx == 2:
				res = def3B[G.Rand.Intn(len(def3B))]
			default:
				res = def3C[G.Rand.Intn(len(def3C))]
			}
		}
		if out {
			if res == blowStag {
				res = blowHesitate
			} else {
				res = blowSitDuck
			}
		}
		if res == blowStag && dweapon != nil && Prob(25, true) {
			res = blowLoseWpn
		}
		msg := randomMeleeMsg((*remarks)[res-1])
		if msg != nil {
			remark(msg, G.Winner, dweapon)
		}
	}
	_ = oa
	// Apply results
	switch {
	case res == blowMissed || res == blowHesitate:
		// Nothing
	case res == blowUncon:
		// Nothing extra
	case res == blowKill || res == blowSitDuck:
		def = 0
	case res == blowLightWnd:
		def--
		if def < 0 {
			def = 0
		}
		if gD().LoadAllowed > 50 {
			gD().LoadAllowed -= 10
		}
	case res == blowHeavyWnd:
		def -= 2
		if def < 0 {
			def = 0
		}
		if gD().LoadAllowed > 50 {
			gD().LoadAllowed -= 20
		}
	case res == blowStag:
		G.Winner.Give(FlgStaggered)
	default:
		// blowLoseWpn
		if dweapon != nil {
			dweapon.MoveTo(G.Here)
			nweapon := findWeapon(G.Winner)
			if nweapon != nil {
				Printf("Fortunately, you still have a %s.\n", nweapon.Desc)
			}
		}
	}
	winnerResult(def, res, od)
	return res
}

// doFight runs the villain combat loop
func doFight(numVillains int) bool {
	for {
		cnt := 0
		var res BlowRes
		out := false
		for cnt < numVillains {
			oo := gD().Villains[cnt]
			cnt++
			o := oo.Villain
			if !o.Has(FlgFight) {
				continue
			}
			if o.Action != nil && o.Action(ActBusy) {
				continue
			}
			res = villainBlow(oo, out)
			if res == blowUnk {
				return false
			}
			if res == blowUncon {
				out = true
			}
		}
		if res != blowUnk {
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

func heroBlow() bool {
	if G.DirObj == nil {
		return false
	}
	// Find the villain entry for the target
	var oo *VillainEntry
	for _, v := range gD().Villains {
		if v.Villain == G.DirObj {
			oo = v
			break
		}
	}
	G.DirObj.Give(FlgFight)
	if G.Winner.Has(FlgStaggered) {
		Printf("You are still recovering from that last blow, so your attack is ineffective.\n")
		G.Winner.Take(FlgStaggered)
		return true
	}
	att := fightStrength(true)
	if att < 1 {
		att = 1
	}
	oa := att
	_ = oa
	villain := G.DirObj
	if oo == nil {
		return false
	}
	od := villainStrength(oo)
	def := od
	if def == 0 {
		if G.DirObj == G.Winner {
			jigsUp("Well, you really did it that time. Is suicide painless?", false)
			return true
		}
		Printf("Attacking the %s is pointless.\n", villain.Desc)
		return true
	}
	dweapon := findWeapon(villain)
	var res BlowRes
	if dweapon == nil || def < 0 {
		Printf("The ")
		if def < 0 {
			Printf("unconscious")
		} else {
			Printf("unarmed")
		}
		Printf(" %s cannot defend himself: He dies.\n", villain.Desc)
		res = blowKill
	} else {
		if def == 1 {
			res = def1[G.Rand.Intn(len(def1))]
		} else if def == 2 {
			res = def2B[G.Rand.Intn(len(def2B))]
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
				res = def3A[G.Rand.Intn(len(def3A))]
			case idx == 2:
				res = def3B[G.Rand.Intn(len(def3B))]
			default:
				res = def3C[G.Rand.Intn(len(def3C))]
			}
		}
		if res == blowStag && dweapon != nil && Prob(25, false) {
			res = blowLoseWpn
		}
		msg := randomMeleeMsg(heroMelee[res-1])
		if msg != nil {
			remark(msg, G.DirObj, G.IndirObj)
		}
	}
	// Apply hero blow results
	switch {
	case res == blowMissed || res == blowHesitate:
		// Nothing
	case res == blowUncon:
		def = -def
	case res == blowKill || res == blowSitDuck:
		def = 0
	case res == blowLightWnd:
		def--
		if def < 0 {
			def = 0
		}
	case res == blowHeavyWnd:
		def -= 2
		if def < 0 {
			def = 0
		}
	case res == blowStag:
		G.DirObj.Give(FlgStaggered)
	default:
		// blowLoseWpn
		if dweapon != nil {
			dweapon.Take(FlgNoDesc)
			dweapon.Give(FlgWeapon)
			dweapon.MoveTo(G.Here)
			thisIsIt(dweapon)
		}
	}
	villainResult(G.DirObj, def, res)
	return true
}
