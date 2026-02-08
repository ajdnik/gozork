package game

import . "github.com/ajdnik/gozork/engine"



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

func StolenLight() bool {
	oLit := G.Lit
	G.Lit = IsLit(G.Here, true)
	if !G.Lit && oLit {
		Printf("The thief seems to have left you in the dark.\n")
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
				GD().EggSolve = true
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
			Printf("You hear, off in the distance, someone saying \"My, I wonder what this fine %s is doing here.\"\n", x.Desc)
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
	if len(G.Here.Children) > 1 {
		Printf("The thief gestures mysteriously, and the treasures in the room suddenly vanish.\n\n")
	}
	for _, f := range G.Here.Children {
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

func RobberFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "tell" {
		Printf("The thief is a strong, silent type.\n")
		G.Params.Continue = NumUndef
		return true
	}
	if arg == ActArg(FBusy) {
		if Stiletto.IsIn(&Thief) {
			return false
		}
		if Stiletto.IsIn(Thief.Location()) {
			Stiletto.MoveTo(&Thief)
			Stiletto.Give(FlgNoDesc)
			if Thief.IsIn(G.Here) {
				Printf("The robber, somewhat surprised at this turn of events, nimbly retrieves his stiletto.\n")
			}
			return true
		}
		return false
	}
	if arg == ActArg(FDead) {
		Stiletto.MoveTo(G.Here)
		Stiletto.Take(FlgNoDesc)
		x := DepositBooty(G.Here)
		if G.Here == &TreasureRoom {
			flg := false
			for _, obj := range G.Here.Children {
				if obj == &Chalice || obj == &Thief || obj == &Adventurer {
					continue
				}
				obj.Take(FlgInvis)
				if !flg {
					flg = true
					Printf("As the thief dies, the power of his magic decreases, and his treasures reappear:\n")
				}
				Printf("  A %s", obj.Desc)
				if obj.HasChildren() && CanSeeInside(obj) {
					Printf(", with ")
					PrintContents(obj)
				}
				Printf("\n")
			}
			if !flg {
				Printf("The chalice is now safe to take.\n")
			}
		} else if x {
			Printf("His booty remains.\n")
		}
		QueueInt("IThief", false).Run = false
		return true
	}
	if arg == ActArg(FFirst) {
		if GD().ThiefHere && !Thief.Has(FlgInvis) && Prob(20, false) {
			Thief.Give(FlgFight)
			G.Params.Continue = NumUndef
			return true
		}
		return false
	}
	if arg == ActArg(FUnconscious) {
		QueueInt("IThief", false).Run = false
		Thief.Take(FlgFight)
		Stiletto.MoveTo(G.Here)
		Stiletto.Take(FlgNoDesc)
		Thief.LongDesc = RobberUDesc
		return true
	}
	if arg == ActArg(FConscious) {
		if Thief.Location() == G.Here {
			Thief.Give(FlgFight)
			Printf("The robber revives, briefly feigning continued unconsciousness, and, when he sees his moment, scrambles away from you.\n")
		}
		Queue("IThief", -1).Run = true
		Thief.LongDesc = RobberCDesc
		RecoverStiletto()
		return true
	}

	// Default (no special mode)
	if G.ActVerb.Norm == "hello" && Thief.LongDesc == RobberUDesc {
		Printf("The thief, being temporarily incapacitated, is unable to acknowledge your greeting with his usual graciousness.\n")
		return true
	}
	if G.DirObj == &Knife && G.ActVerb.Norm == "throw" && !Thief.Has(FlgFight) {
		G.DirObj.MoveTo(G.Here)
		if Prob(10, false) {
			Printf("You evidently frightened the robber, though you didn't hit him. He flees")
			LargeBag.Remove()
			hasStiletto := false
			if Stiletto.IsIn(&Thief) {
				Stiletto.Remove()
				hasStiletto = true
			}
			if Thief.HasChildren() {
				MoveAll(&Thief, G.Here)
				Printf(", but the contents of his bag fall on the floor.")
			} else {
				Printf(".")
			}
			LargeBag.MoveTo(&Thief)
			if hasStiletto {
				Stiletto.MoveTo(&Thief)
			}
			Printf("\n")
			Thief.Give(FlgInvis)
		} else {
			Printf("You missed. The thief makes no attempt to take the knife, though it would be a fine addition to the collection in his bag. He does seem angered by your attempt.\n")
			Thief.Give(FlgFight)
		}
		return true
	}
	if (G.ActVerb.Norm == "throw" || G.ActVerb.Norm == "give") && G.DirObj != nil && G.DirObj != &Thief && G.IndirObj == &Thief {
		if Thief.Strength < 0 {
			Thief.Strength = -Thief.Strength
			Queue("IThief", -1).Run = true
			RecoverStiletto()
			Thief.LongDesc = RobberCDesc
			Printf("Your proposed victim suddenly recovers consciousness.\n")
		}
		G.DirObj.MoveTo(&Thief)
		if G.DirObj.TValue > 0 {
			GD().ThiefEngrossed = true
			Printf("The thief is taken aback by your unexpected generosity, but accepts the %s and stops to admire its beauty.\n", G.DirObj.Desc)
		} else {
			Printf("The thief places the %s in his bag and thanks you politely.\n", G.DirObj.Desc)
		}
		return true
	}
	if G.ActVerb.Norm == "take" {
		Printf("Once you got him, what would you do with him?\n")
		return true
	}
	if G.ActVerb.Norm == "examine" || G.ActVerb.Norm == "look inside" {
		Printf("The thief is a slippery character with beady eyes that flit back and forth. He carries, along with an unmistakable arrogance, a large bag over his shoulder and a vicious stiletto, whose blade is aimed menacingly in your direction. I'd watch out if I were you.\n")
		return true
	}
	if G.ActVerb.Norm == "listen" {
		Printf("The thief says nothing, as you have not been formally introduced.\n")
		return true
	}
	return false
}

func LargeBagFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "take" {
		if Thief.LongDesc == RobberUDesc {
			Printf("Sadly for you, the robber collapsed on top of the bag. Trying to take it would wake him.\n")
		} else {
			Printf("The bag will be taken over his dead body.\n")
		}
		return true
	}
	if G.ActVerb.Norm == "put" && G.IndirObj == &LargeBag {
		Printf("It would be a good trick.\n")
		return true
	}
	if G.ActVerb.Norm == "open" || G.ActVerb.Norm == "close" {
		Printf("Getting close enough would be a good trick.\n")
		return true
	}
	if G.ActVerb.Norm == "examine" || G.ActVerb.Norm == "look inside" {
		Printf("The bag is underneath the thief, so one can't say what, if anything, is inside.\n")
		return true
	}
	return false
}

func StiletteFcn(arg ActArg) bool {
	return WeaponFunction(&Stiletto, &Thief)
}

func TreasureInsideFcn(arg ActArg) bool {
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
		if rm == &TreasureRoom && rm != G.Here {
			if hereQ {
				HackTreasures()
				hereQ = false
			}
			DepositBooty(&TreasureRoom)
		} else if rm == G.Here && !G.Here.Has(FlgOn) && !Troll.IsIn(G.Here) {
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
				if rm.Has(FlgMaze) && G.Here.Has(FlgMaze) {
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
					GD().ThiefHere = false
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
	if !GD().Dead && G.Here == &TreasureRoom {
		return false
	}
	if !GD().ThiefHere {
		if !GD().Dead && !hereQ && Prob(30, false) {
			if Stiletto.IsIn(&Thief) {
				Thief.Take(FlgInvis)
				Printf("Someone carrying a large bag is casually leaning against one of the walls here. He does not speak, but it is clear from his aspect that the bag will be taken only over his dead body.\n")
				GD().ThiefHere = true
				return true
			}
		}
		if hereQ && Thief.Has(FlgFight) && !Winning(&Thief) {
			Printf("Your opponent, determining discretion to be the better part of valor, decides to terminate this little contretemps. With a rueful nod of his head, he steps backward into the gloom and disappears.\n")
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
			if !flg && rm == G.Here {
				Printf("The robber, rummaging through his bag, dropped a few items he found valueless.\n")
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
					GD().DomeFlag = false
				}
				if rm == G.Here {
					Printf("You suddenly notice that the %s vanished.\n", x.Desc)
					return true
				}
				return false
			}
		}
	}
	return false
}
