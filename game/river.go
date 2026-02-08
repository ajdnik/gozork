package game

import . "github.com/ajdnik/gozork/engine"

func fixBoat() {
	Printf("Well done. The boat is repaired.\n")
	inflatableBoat.MoveTo(puncturedBoat.Location())
	removeCarefully(&puncturedBoat)
}

func fixMaintLeak() {
	gD().WaterLevel = -1
	QueueInt("iMaintRoom", false).Run = false
	Printf("By some miracle of Zorkian technology, you have managed to stop the leak in the dam.\n")
}

func waterFcn(arg ActionArg) bool {
	if G.ActVerb.Norm == "sgive" {
		return false
	}
	if G.ActVerb.Norm == "through" || G.ActVerb.Norm == "board" {
		Printf("%s\n", PickOne(swimYuks))
		return true
	}
	// Simplified water handling
	if G.ActVerb.Norm == "take" || G.ActVerb.Norm == "put" {
		w := G.DirObj
		if w == &globalWater {
			w = &water
		}
		if G.ActVerb.Norm == "take" {
			if w.IsIn(&bottle) && G.IndirObj == nil {
				Printf("it's in the bottle. Perhaps you should take that instead.\n")
				return true
			}
			if bottle.IsIn(G.Winner) {
				if !bottle.Has(FlgOpen) {
					Printf("The bottle is closed.\n")
					thisIsIt(&bottle)
					return true
				}
				if !bottle.HasChildren() {
					water.MoveTo(&bottle)
					Printf("The bottle is now full of water.\n")
					return true
				}
				Printf("The water slips through your fingers.\n")
				return true
			}
			Printf("The water slips through your fingers.\n")
			return true
		}
	}
	if G.ActVerb.Norm == "drop" || G.ActVerb.Norm == "give" {
		if G.ActVerb.Norm == "drop" && water.IsIn(&bottle) && !bottle.Has(FlgOpen) {
			Printf("The bottle is closed.\n")
			return true
		}
		removeCarefully(&water)
		av := G.Winner.Location()
		if av.Has(FlgVeh) {
			Printf("There is now a puddle in the bottom of the %s.\n", av.Desc)
			water.MoveTo(av)
		} else {
			Printf("The water spills to the floor and evaporates immediately.\n")
		}
		return true
	}
	if G.ActVerb.Norm == "throw" {
		Printf("The water splashes on the walls and evaporates immediately.\n")
		removeCarefully(&water)
		return true
	}
	return false
}

func boltFcn(arg ActionArg) bool {
	if G.ActVerb.Norm == "turn" {
		if G.IndirObj == &wrench {
			if gD().GateFlag {
				reservoirSouth.Take(FlgTouch)
				if gD().GatesOpen {
					gD().GatesOpen = false
					loudRoom.Take(FlgTouch)
					Printf("The sluice gates close and water starts to collect behind the dam.\n")
					Queue("iRfill", 8).Run = true
					QueueInt("iRempty", false).Run = false
				} else {
					gD().GatesOpen = true
					Printf("The sluice gates open and water pours through the dam.\n")
					Queue("iRempty", 8).Run = true
					QueueInt("iRfill", false).Run = false
				}
			} else {
				Printf("The bolt won't turn with your best effort.\n")
			}
		} else {
			Printf("The bolt won't turn using the %s.\n", G.IndirObj.Desc)
		}
		return true
	}
	if G.ActVerb.Norm == "take" {
		integralPart()
		return true
	}
	if G.ActVerb.Norm == "oil" {
		Printf("Hmm. it appears the tube contained glue, not oil. Turning the bolt won't get any easier....\n")
		return true
	}
	return false
}

func bubbleFcn(arg ActionArg) bool {
	if G.ActVerb.Norm == "take" {
		integralPart()
		return true
	}
	return false
}

func damFunction(arg ActionArg) bool {
	if G.ActVerb.Norm == "open" || G.ActVerb.Norm == "close" {
		Printf("Sounds reasonable, but this isn't how.\n")
		return true
	}
	if G.ActVerb.Norm == "plug" {
		if G.IndirObj == &hands {
			Printf("Are you the little Dutch boy, then? Sorry, this is a big dam.\n")
		} else {
			Printf("With a %s? Do you know how big this dam is? You could only stop a tiny leak with that.\n", G.IndirObj.Desc)
		}
		return true
	}
	return false
}

func puncturedBoatFcn(arg ActionArg) bool {
	if (G.ActVerb.Norm == "put" || G.ActVerb.Norm == "put on") && G.DirObj == &putty {
		fixBoat()
		return true
	}
	if G.ActVerb.Norm == "inflate" || G.ActVerb.Norm == "fill" {
		Printf("No chance. Some moron punctured it.\n")
		return true
	}
	if G.ActVerb.Norm == "plug" {
		if G.IndirObj == &putty {
			fixBoat()
			return true
		}
		withTell(G.IndirObj)
		return true
	}
	return false
}

func inflatableBoatFcn(arg ActionArg) bool {
	if G.ActVerb.Norm == "inflate" || G.ActVerb.Norm == "fill" {
		if !inflatableBoat.IsIn(G.Here) {
			Printf("The boat must be on the ground to be inflated.\n")
			return true
		}
		if G.IndirObj == &pump {
			Printf("The boat inflates and appears seaworthy.\n")
			if !boatLabel.Has(FlgTouch) {
				Printf("A tan label is lying inside the boat.\n")
			}
			gD().Deflate = false
			removeCarefully(&inflatableBoat)
			inflatedBoat.MoveTo(G.Here)
			thisIsIt(&inflatedBoat)
			return true
		}
		if G.IndirObj == &lungs {
			Printf("You don't have enough lung power to inflate it.\n")
			return true
		}
		Printf("With a %s? Surely you jest!\n", G.IndirObj.Desc)
		return true
	}
	return false
}

func riverFcn(arg ActionArg) bool {
	if G.ActVerb.Norm == "put" && G.IndirObj == &river {
		if G.DirObj == &me {
			jigsUp("You splash around for a while, fighting the current, then you drown.", false)
			return true
		}
		if G.DirObj == &inflatedBoat {
			Printf("You should get in the boat then launch it.\n")
			return true
		}
		if G.DirObj.Has(FlgBurn) {
			removeCarefully(G.DirObj)
			Printf("The %s floats for a moment, then sinks.\n", G.DirObj.Desc)
			return true
		}
		removeCarefully(G.DirObj)
		Printf("The %s splashes into the water and is gone forever.\n", G.DirObj.Desc)
		return true
	}
	if G.ActVerb.Norm == "leap" || G.ActVerb.Norm == "through" {
		Printf("A look before leaping reveals that the river is wide and dangerous, with swift currents and large, half-hidden rocks. You decide to forgo your swim.\n")
		return true
	}
	return false
}

func damRoomFcn(arg ActionArg) bool {
	if arg == ActLook {
		Printf("You are standing on the top of the Flood Control dam #3, which was quite a tourist attraction in times far distant. There are paths to the north, south, and west, and a scramble down.\n")
		if gD().LowTide && gD().GatesOpen {
			Printf("The water level behind the dam is low: The sluice gates have been opened. water rushes through the dam and downstream.\n")
		} else if gD().GatesOpen {
			Printf("The sluice gates are open, and water rushes through the dam. The water level behind the dam is still high.\n")
		} else if gD().LowTide {
			Printf("The sluice gates are closed. The water level in the reservoir is quite low, but the level is rising quickly.\n")
		} else {
			Printf("The sluice gates on the dam are closed. Behind the dam, there can be seen a wide reservoir. water is pouring over the top of the now abandoned dam.\n")
		}
		Printf("There is a control panel here, on which a large metal bolt is mounted. Directly above the bolt is a small green plastic bubble")
		if gD().GateFlag {
			Printf(" which is glowing serenely")
		}
		Printf(".\n")
		return true
	}
	return false
}

func whiteCliffsFcn(arg ActionArg) bool {
	if arg == ActEnd {
		if inflatedBoat.IsIn(G.Winner) {
			gD().Deflate = false
		} else {
			gD().Deflate = true
		}
	}
	return false
}

func fallsRoomFcn(arg ActionArg) bool {
	if arg == ActLook {
		Printf("You are at the top of Aragain Falls, an enormous waterfall with a drop of about 450 feet. The only path here is on the north end.\n")
		if gD().RainbowFlag {
			Printf("A solid rainbow spans the falls.\n")
		} else {
			Printf("A beautiful rainbow can be seen over the falls and to the west.\n")
		}
		return true
	}
	return false
}

func rivr4RoomFcn(arg ActionArg) bool {
	if arg == ActEnd {
		if buoy.IsIn(G.Winner) && gD().BuoyFlag {
			Printf("You notice something funny about the feel of the buoy.\n")
			gD().BuoyFlag = false
		}
	}
	return false
}

func rBoatFcn(arg ActionArg) bool {
	if arg == ActEnter || arg == ActEnd || arg == ActLook {
		return false
	}
	if arg == ActBegin {
		if G.ActVerb.Norm == "walk" && G.Params.HasWalkDir {
			if G.Params.WalkDir == Land || G.Params.WalkDir == East || G.Params.WalkDir == West {
				return false
			}
			if G.Here == &reservoir && (G.Params.WalkDir == North || G.Params.WalkDir == South) {
				return false
			}
			if G.Here == &inStream && G.Params.WalkDir == South {
				return false
			}
			Printf("Read the label for the boat's instructions.\n")
			return true
		}
		if G.ActVerb.Norm == "launch" {
			if G.Here == &river1 || G.Here == &river2 || G.Here == &river3 || G.Here == &river4 || G.Here == &reservoir || G.Here == &inStream {
				Printf("You are on the ")
				if G.Here == &reservoir {
					Printf("reservoir")
				} else if G.Here == &inStream {
					Printf("stream")
				} else {
					Printf("river")
				}
				Printf(", or have you forgotten?\n")
				return true
			}
			tmp := goNext(riverLaunch)
			if tmp == 1 {
				// ZIL: <ENABLE <QUEUE I-RIVER <LKP ,HERE ,RIVER-SPEEDS>>>
				// After goNext, HERE is the destination room. Use its speed.
				if spd, ok := riverSpeedMap[G.Here]; ok {
					Queue("iRiver", spd).Run = true
				}
				return true
			}
			if tmp != 2 {
				Printf("You can't launch it here.\n")
				return true
			}
			return true
		}
		if (G.ActVerb.Norm == "drop" && G.DirObj.Has(FlgWeapon)) ||
			(G.ActVerb.Norm == "put" && G.DirObj.Has(FlgWeapon) && G.IndirObj == &inflatedBoat) ||
			((G.ActVerb.Norm == "attack" || G.ActVerb.Norm == "mung") && G.IndirObj != nil && G.IndirObj.Has(FlgWeapon)) {
			removeCarefully(&inflatedBoat)
			puncturedBoat.MoveTo(G.Here)
			rob(&inflatedBoat, G.Here, 0)
			G.Winner.MoveTo(G.Here)
			Printf("it seems that the ")
			if G.ActVerb.Norm == "drop" || G.ActVerb.Norm == "put" {
				Printf("%s", G.DirObj.Desc)
			} else {
				Printf("%s", G.IndirObj.Desc)
			}
			Printf(" didn't agree with the boat, as evidenced by the loud hissing noise issuing therefrom. With a pathetic sputter, the boat deflates, leaving you without.\n")
			if G.Here.Has(FlgNonLand) {
				Printf("\n")
				if G.Here == &reservoir || G.Here == &inStream {
					jigsUp("Another pathetic sputter, this time from you, heralds your drowning.", false)
				} else {
					jigsUp("In other words, fighting the fierce currents of the Frigid river. You manage to hold your own for a bit, but then you are carried over a waterfall and into some nasty rocks. Ouch!", false)
				}
			}
			return true
		}
		return false
	}
	if G.ActVerb.Norm == "board" {
		if sceptre.IsIn(G.Winner) || knife.IsIn(G.Winner) || sword.IsIn(G.Winner) || rustyKnife.IsIn(G.Winner) || axe.IsIn(G.Winner) || stiletto.IsIn(G.Winner) {
			Printf("Oops! Something sharp seems to have slipped and punctured the boat. The boat deflates to the sounds of hissing, sputtering, and cursing.\n")
			removeCarefully(&inflatedBoat)
			puncturedBoat.MoveTo(G.Here)
			thisIsIt(&puncturedBoat)
			return true
		}
		return false
	}
	if G.ActVerb.Norm == "inflate" || G.ActVerb.Norm == "fill" {
		Printf("Inflating it further would probably burst it.\n")
		return true
	}
	if G.ActVerb.Norm == "deflate" {
		if G.Winner.Location() == &inflatedBoat {
			Printf("You can't deflate the boat while you're in it.\n")
			return true
		}
		if !inflatedBoat.IsIn(G.Here) {
			Printf("The boat must be on the ground to be deflated.\n")
			return true
		}
		Printf("The boat deflates.\n")
		gD().Deflate = true
		removeCarefully(&inflatedBoat)
		inflatableBoat.MoveTo(G.Here)
		thisIsIt(&inflatableBoat)
		return true
	}
	return false
}

func iRfill() bool {
	reservoir.Give(FlgNonLand)
	reservoir.Take(FlgRLand)
	deepCanyon.Take(FlgTouch)
	loudRoom.Take(FlgTouch)
	if trunk.IsIn(&reservoir) {
		trunk.Give(FlgInvis)
	}
	gD().LowTide = false
	if G.Here == &reservoir {
		if G.Winner.Location().Has(FlgVeh) {
			Printf("The boat lifts gently out of the mud and is now floating on the reservoir.\n")
		} else {
			jigsUp("You are lifted up by the rising river! You try to swim, but the currents are too strong. You come closer, closer to the awesome structure of Flood Control dam #3. The dam beckons to you. The roar of the water nearly deafens you, but you remain conscious as you tumble over the dam toward your certain doom among the rocks at its base.", false)
		}
	} else if G.Here == &deepCanyon {
		Printf("A sound, like that of flowing water, starts to come from below.\n")
	} else if G.Here == &loudRoom {
		Printf("All of a sudden, an alarmingly loud roaring sound fills the room. Filled with fear, you scramble away.\n")
		dest := loudRuns[G.Rand.Intn(len(loudRuns))]
		gotoRoom(dest, true)
	} else if G.Here == &reservoirNorth || G.Here == &reservoirSouth {
		Printf("You notice that the water level has risen to the point that it is impossible to cross.\n")
	}
	return true
}

func iRempty() bool {
	reservoir.Give(FlgRLand)
	reservoir.Take(FlgNonLand)
	deepCanyon.Take(FlgTouch)
	loudRoom.Take(FlgTouch)
	trunk.Take(FlgInvis)
	gD().LowTide = true
	if G.Here == &reservoir && G.Winner.Location().Has(FlgVeh) {
		Printf("The water level has dropped to the point at which the boat can no longer stay afloat. it sinks into the mud.\n")
	} else if G.Here == &deepCanyon {
		Printf("The roar of rushing water is quieter now.\n")
	} else if G.Here == &reservoirNorth || G.Here == &reservoirSouth {
		Printf("The water level is now quite low here and you could easily cross over to the other side.\n")
	}
	return true
}

func iMaintRoom() bool {
	hereQ := G.Here == &maintenanceRoom
	if hereQ {
		Printf("The water level here is now ")
		idx := gD().WaterLevel / 2
		if idx >= 0 && idx < len(drownings) {
			Printf("%s", drownings[idx])
		}
		Printf("\n")
	}
	gD().WaterLevel++
	if gD().WaterLevel >= 14 {
		mungRoom(&maintenanceRoom, "The room is full of water and cannot be entered.")
		QueueInt("iMaintRoom", false).Run = false
		if hereQ {
			jigsUp("I'm afraid you have done drowned yourself.", false)
		}
	} else if inflatedBoat.IsIn(G.Winner) && (G.Here == &maintenanceRoom || G.Here == &damRoom || G.Here == &damLobby) {
		jigsUp("The rising water carries the boat over the dam, down the river, and over the falls. Tsk, tsk.", false)
	}
	return true
}

func iRiver() bool {
	if G.Here != &river1 && G.Here != &river2 && G.Here != &river3 && G.Here != &river4 && G.Here != &river5 {
		QueueInt("iRiver", false).Run = false
		return false
	}
	rm, ok := riverNext[G.Here]
	if ok {
		Printf("The flow of the river carries you downstream.\n\n")
		gotoRoom(rm, true)
		// ZIL: <ENABLE <QUEUE I-RIVER <LKP ,HERE ,RIVER-SPEEDS>>>
		// After gotoRoom, HERE is the new room. Use its speed for the next leg.
		if spd, ok := riverSpeedMap[G.Here]; ok {
			Queue("iRiver", spd).Run = true
		}
		return true
	}
	jigsUp("Unfortunately, the magic boat doesn't provide protection from the rocks and boulders one meets at the bottom of waterfalls. Including this one.", false)
	return true
}
