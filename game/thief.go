package game

import . "github.com/ajdnik/gozork/engine"

func rob(what, where *Object, prob int) bool {
	robbed := false
	for _, x := range what.Children {
		if x.Has(FlgInvis) || x.Has(FlgSacred) {
			continue
		}
		if x.GetTValue() <= 0 {
			continue
		}
		if prob > 0 && !Prob(prob, false) {
			continue
		}
		x.MoveTo(where)
		x.Give(FlgTouch)
		if where == &thief {
			x.Give(FlgInvis)
		}
		robbed = true
	}
	return robbed
}

func recoverStiletto() {
	if stiletto.IsIn(thief.Location()) {
		stiletto.Give(FlgNoDesc)
		stiletto.MoveTo(&thief)
	}
}

func hackTreasures() {
	recoverStiletto()
	thief.Give(FlgInvis)
	for _, x := range treasureRoom.Children {
		x.Take(FlgInvis)
	}
}

func depositBooty(rm *Object) bool {
	flg := false
	var toMove []*Object
	for _, x := range thief.Children {
		if x == &stiletto || x == &largeBag {
			continue
		}
		if x.GetTValue() > 0 {
			toMove = append(toMove, x)
			flg = true
			if x == &egg {
				gD().EggSolve = true
				egg.Give(FlgOpen)
			}
		}
	}
	for _, x := range toMove {
		x.MoveTo(rm)
	}
	return flg
}

func robMaze(rm *Object) bool {
	for _, x := range rm.Children {
		if x.Has(FlgTake) && !x.Has(FlgInvis) && Prob(40, false) {
			Printf("You hear, off in the distance, someone saying \"My, I wonder what this fine %s is doing here.\"\n", x.Desc)
			if Prob(60, true) {
				x.MoveTo(&thief)
				x.Give(FlgTouch)
				x.Give(FlgInvis)
			}
			return true
		}
	}
	return false
}

func moveAll(from, to *Object) {
	toMove := append([]*Object{}, from.Children...)
	for _, x := range toMove {
		x.Take(FlgInvis)
		x.MoveTo(to)
	}
}

func thiefInTreasure() {
	if len(G.Here.Children) > 1 {
		Printf("The thief gestures mysteriously, and the treasures in the room suddenly vanish.\n\n")
	}
	for _, f := range G.Here.Children {
		if f != &chalice && f != &thief {
			f.Give(FlgInvis)
		}
	}
}

func infested(r *Object) bool {
	for _, f := range r.Children {
		if f.Has(FlgActor) && !f.Has(FlgInvis) {
			return true
		}
	}
	return false
}

func robberFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "tell":
		Printf("The thief is a strong, silent type.\n")
		G.Params.Continue = NumUndef
		return true
	}
	if arg == ActBusy {
		if stiletto.IsIn(&thief) {
			return false
		}
		if stiletto.IsIn(thief.Location()) {
			stiletto.MoveTo(&thief)
			stiletto.Give(FlgNoDesc)
			if thief.IsIn(G.Here) {
				Printf("The robber, somewhat surprised at this turn of events, nimbly retrieves his stiletto.\n")
			}
			return true
		}
		return false
	}
	if arg == ActDead {
		stiletto.MoveTo(G.Here)
		stiletto.Take(FlgNoDesc)
		x := depositBooty(G.Here)
		if G.Here == &treasureRoom {
			flg := false
			for _, obj := range G.Here.Children {
				if obj == &chalice || obj == &thief || obj == &adventurer {
					continue
				}
				obj.Take(FlgInvis)
				if !flg {
					flg = true
					Printf("As the thief dies, the power of his magic decreases, and his treasures reappear:\n")
				}
				Printf("  A %s", obj.Desc)
				if obj.HasChildren() && canSeeInside(obj) {
					Printf(", with ")
					printContents(obj)
				}
				Printf("\n")
			}
			if !flg {
				Printf("The chalice is now safe to take.\n")
			}
		} else if x {
			Printf("His booty remains.\n")
		}
		QueueInt("iThief", false).Run = false
		return true
	}
	if arg == ActFirst {
		if gD().ThiefHere && !thief.Has(FlgInvis) && Prob(20, false) {
			thief.Give(FlgFight)
			G.Params.Continue = NumUndef
			return true
		}
		return false
	}
	if arg == ActUnconscious {
		QueueInt("iThief", false).Run = false
		thief.Take(FlgFight)
		stiletto.MoveTo(G.Here)
		stiletto.Take(FlgNoDesc)
		thief.LongDesc = robberUDesc
		return true
	}
	if arg == ActConscious {
		if thief.Location() == G.Here {
			thief.Give(FlgFight)
			Printf("The robber revives, briefly feigning continued unconsciousness, and, when he sees his moment, scrambles away from you.\n")
		}
		Queue("iThief", -1).Run = true
		thief.LongDesc = robberCDesc
		recoverStiletto()
		return true
	}

	// Default (no special mode)
	switch G.ActVerb.Norm {
	case "hello":
		if thief.LongDesc == robberUDesc {
			Printf("The thief, being temporarily incapacitated, is unable to acknowledge your greeting with his usual graciousness.\n")
			return true
		}
	case "throw":
		if G.DirObj == &knife && !thief.Has(FlgFight) {
			G.DirObj.MoveTo(G.Here)
			if Prob(10, false) {
				Printf("You evidently frightened the robber, though you didn't hit him. He flees")
				largeBag.Remove()
				hasStiletto := false
				if stiletto.IsIn(&thief) {
					stiletto.Remove()
					hasStiletto = true
				}
				if thief.HasChildren() {
					moveAll(&thief, G.Here)
					Printf(", but the contents of his bag fall on the floor.")
				} else {
					Printf(".")
				}
				largeBag.MoveTo(&thief)
				if hasStiletto {
					stiletto.MoveTo(&thief)
				}
				Printf("\n")
				thief.Give(FlgInvis)
			} else {
				Printf("You missed. The thief makes no attempt to take the knife, though it would be a fine addition to the collection in his bag. He does seem angered by your attempt.\n")
				thief.Give(FlgFight)
			}
			return true
		}
		fallthrough
	case "give":
		if G.DirObj != nil && G.DirObj != &thief && G.IndirObj == &thief {
			if thief.GetStrength() < 0 {
				thief.SetStrength(-thief.GetStrength())
				Queue("iThief", -1).Run = true
				recoverStiletto()
				thief.LongDesc = robberCDesc
				Printf("Your proposed victim suddenly recovers consciousness.\n")
			}
			G.DirObj.MoveTo(&thief)
			if G.DirObj.GetTValue() > 0 {
				gD().ThiefEngrossed = true
				Printf("The thief is taken aback by your unexpected generosity, but accepts the %s and stops to admire its beauty.\n", G.DirObj.Desc)
			} else {
				Printf("The thief places the %s in his bag and thanks you politely.\n", G.DirObj.Desc)
			}
			return true
		}
	case "take":
		Printf("Once you got him, what would you do with him?\n")
		return true
	case "examine", "look inside":
		Printf("The thief is a slippery character with beady eyes that flit back and forth. He carries, along with an unmistakable arrogance, a large bag over his shoulder and a vicious stiletto, whose blade is aimed menacingly in your direction. I'd watch out if I were you.\n")
		return true
	case "listen":
		Printf("The thief says nothing, as you have not been formally introduced.\n")
		return true
	}
	return false
}

func largeBagFcn(arg ActionArg) bool {
	switch G.ActVerb.Norm {
	case "take":
		if thief.LongDesc == robberUDesc {
			Printf("Sadly for you, the robber collapsed on top of the bag. Trying to take it would wake him.\n")
		} else {
			Printf("The bag will be taken over his dead body.\n")
		}
		return true
	case "put":
		if G.IndirObj == &largeBag {
			Printf("it would be a good trick.\n")
			return true
		}
	case "open", "close":
		Printf("Getting close enough would be a good trick.\n")
		return true
	case "examine", "look inside":
		Printf("The bag is underneath the thief, so one can't say what, if anything, is inside.\n")
		return true
	}
	return false
}

func stiletteFcn(arg ActionArg) bool {
	return weaponFunction(&stiletto, &thief)
}

func treasureInsideFcn(arg ActionArg) bool {
	return false
}

func iThief() bool {
	rm := thief.Location()
	hereQ := !thief.Has(FlgInvis)
	if hereQ {
		rm = thief.Location()
	}
	flg := false
	once := false
	for {
		if rm == &treasureRoom && rm != G.Here {
			if hereQ {
				hackTreasures()
				hereQ = false
			}
			depositBooty(&treasureRoom)
		} else if rm == G.Here && !G.Here.Has(FlgOn) && !troll.IsIn(G.Here) {
			if thiefVsAdventurer(hereQ) {
				return true
			}
			if thief.Has(FlgInvis) {
				hereQ = false
			}
		} else {
			if thief.IsIn(rm) && !thief.Has(FlgInvis) {
				// Leave if victim left
				thief.Give(FlgInvis)
				hereQ = false
			}
			if rm != nil && rm.Has(FlgTouch) {
				rob(rm, &thief, 75)
				if rm.Has(FlgMaze) && G.Here.Has(FlgMaze) {
					flg = robMaze(rm)
				} else {
					flg = stealJunk(rm)
				}
			}
		}
		if !once && !hereQ {
			once = true
			// Move to next room
			recoverStiletto()
			found := false
			for _, r := range rooms.Children {
				if !r.Has(FlgSacred) && r.Has(FlgRLand) {
					thief.MoveTo(r)
					thief.Take(FlgFight)
					thief.Give(FlgInvis)
					gD().ThiefHere = false
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
	if rm != &treasureRoom {
		dropJunk(rm)
	}
	return flg
}

func thiefVsAdventurer(hereQ bool) bool {
	if !gD().Dead && G.Here == &treasureRoom {
		return false
	}
	if !gD().ThiefHere {
		if !gD().Dead && !hereQ && Prob(30, false) {
			if stiletto.IsIn(&thief) {
				thief.Take(FlgInvis)
				Printf("Someone carrying a large bag is casually leaning against one of the walls here. He does not speak, but it is clear from his aspect that the bag will be taken only over his dead body.\n")
				gD().ThiefHere = true
				return true
			}
		}
		if hereQ && thief.Has(FlgFight) && !winning(&thief) {
			Printf("Your opponent, determining discretion to be the better part of valor, decides to terminate this little contretemps. With a rueful nod of his head, he steps backward into the gloom and disappears.\n")
			thief.Give(FlgInvis)
			thief.Take(FlgFight)
			recoverStiletto()
			return true
		}
	}
	return false
}

// dropJunk - thief drops valueless items from his bag
func dropJunk(rm *Object) bool {
	if rm == nil {
		return false
	}
	flg := false
	for _, x := range thief.Children {
		if x == &stiletto || x == &largeBag {
			continue
		}
		if x.GetTValue() == 0 && Prob(30, true) {
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

// stealJunk - thief steals worthless items from a room
func stealJunk(rm *Object) bool {
	if rm == nil {
		return false
	}
	for _, x := range rm.Children {
		if x.GetTValue() == 0 && x.Has(FlgTake) && !x.Has(FlgSacred) && !x.Has(FlgInvis) {
			if x == &stiletto || Prob(10, true) {
				x.MoveTo(&thief)
				x.Give(FlgTouch)
				x.Give(FlgInvis)
				if x == &rope {
					gD().DomeFlag = false
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
