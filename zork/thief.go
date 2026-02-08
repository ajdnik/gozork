package zork



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
