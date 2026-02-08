package zork



func FixBoat() {
	Printf("Well done. The boat is repaired.\n")
	InflatableBoat.MoveTo(PuncturedBoat.Location())
	RemoveCarefully(&PuncturedBoat)
}

func FixMaintLeak() {
	G.WaterLevel = -1
	QueueInt("IMaintRoom", false).Run = false
	Printf("By some miracle of Zorkian technology, you have managed to stop the leak in the dam.\n")
}

func WaterFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "sgive" {
		return false
	}
	if G.ActVerb.Norm == "through" || G.ActVerb.Norm == "board" {
		Printf("%s\n", PickOne(SwimYuks))
		return true
	}
	// Simplified water handling
	if G.ActVerb.Norm == "take" || G.ActVerb.Norm == "put" {
		w := G.DirObj
		if w == &GlobalWater {
			w = &Water
		}
		if G.ActVerb.Norm == "take" {
			if w.IsIn(&Bottle) && G.IndirObj == nil {
				Printf("It's in the bottle. Perhaps you should take that instead.\n")
				return true
			}
			if Bottle.IsIn(G.Winner) {
				if !Bottle.Has(FlgOpen) {
					Printf("The bottle is closed.\n")
					ThisIsIt(&Bottle)
					return true
				}
				if !Bottle.HasChildren() {
					Water.MoveTo(&Bottle)
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
		if G.ActVerb.Norm == "drop" && Water.IsIn(&Bottle) && !Bottle.Has(FlgOpen) {
			Printf("The bottle is closed.\n")
			return true
		}
		RemoveCarefully(&Water)
		av := G.Winner.Location()
		if av.Has(FlgVeh) {
			Printf("There is now a puddle in the bottom of the %s.\n", av.Desc)
			Water.MoveTo(av)
		} else {
			Printf("The water spills to the floor and evaporates immediately.\n")
		}
		return true
	}
	if G.ActVerb.Norm == "throw" {
		Printf("The water splashes on the walls and evaporates immediately.\n")
		RemoveCarefully(&Water)
		return true
	}
	return false
}

func BoltFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "turn" {
		if G.IndirObj == &Wrench {
			if G.GateFlag {
				ReservoirSouth.Take(FlgTouch)
				if G.GatesOpen {
					G.GatesOpen = false
					LoudRoom.Take(FlgTouch)
					Printf("The sluice gates close and water starts to collect behind the dam.\n")
					Queue("IRfill", 8).Run = true
					QueueInt("IRempty", false).Run = false
				} else {
					G.GatesOpen = true
					Printf("The sluice gates open and water pours through the dam.\n")
					Queue("IRempty", 8).Run = true
					QueueInt("IRfill", false).Run = false
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
		IntegralPart()
		return true
	}
	if G.ActVerb.Norm == "oil" {
		Printf("Hmm. It appears the tube contained glue, not oil. Turning the bolt won't get any easier....\n")
		return true
	}
	return false
}

func BubbleFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "take" {
		IntegralPart()
		return true
	}
	return false
}

func DamFunction(arg ActArg) bool {
	if G.ActVerb.Norm == "open" || G.ActVerb.Norm == "close" {
		Printf("Sounds reasonable, but this isn't how.\n")
		return true
	}
	if G.ActVerb.Norm == "plug" {
		if G.IndirObj == &Hands {
			Printf("Are you the little Dutch boy, then? Sorry, this is a big dam.\n")
		} else {
			Printf("With a %s? Do you know how big this dam is? You could only stop a tiny leak with that.\n", G.IndirObj.Desc)
		}
		return true
	}
	return false
}

func PuncturedBoatFcn(arg ActArg) bool {
	if (G.ActVerb.Norm == "put" || G.ActVerb.Norm == "put on") && G.DirObj == &Putty {
		FixBoat()
		return true
	}
	if G.ActVerb.Norm == "inflate" || G.ActVerb.Norm == "fill" {
		Printf("No chance. Some moron punctured it.\n")
		return true
	}
	if G.ActVerb.Norm == "plug" {
		if G.IndirObj == &Putty {
			FixBoat()
			return true
		}
		WithTell(G.IndirObj)
		return true
	}
	return false
}

func InflatableBoatFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "inflate" || G.ActVerb.Norm == "fill" {
		if !InflatableBoat.IsIn(G.Here) {
			Printf("The boat must be on the ground to be inflated.\n")
			return true
		}
		if G.IndirObj == &Pump {
			Printf("The boat inflates and appears seaworthy.\n")
			if !BoatLabel.Has(FlgTouch) {
				Printf("A tan label is lying inside the boat.\n")
			}
			G.Deflate = false
			RemoveCarefully(&InflatableBoat)
			InflatedBoat.MoveTo(G.Here)
			ThisIsIt(&InflatedBoat)
			return true
		}
		if G.IndirObj == &Lungs {
			Printf("You don't have enough lung power to inflate it.\n")
			return true
		}
		Printf("With a %s? Surely you jest!\n", G.IndirObj.Desc)
		return true
	}
	return false
}

func RiverFcn(arg ActArg) bool {
	if G.ActVerb.Norm == "put" && G.IndirObj == &River {
		if G.DirObj == &Me {
			JigsUp("You splash around for a while, fighting the current, then you drown.", false)
			return true
		}
		if G.DirObj == &InflatedBoat {
			Printf("You should get in the boat then launch it.\n")
			return true
		}
		if G.DirObj.Has(FlgBurn) {
			RemoveCarefully(G.DirObj)
			Printf("The %s floats for a moment, then sinks.\n", G.DirObj.Desc)
			return true
		}
		RemoveCarefully(G.DirObj)
		Printf("The %s splashes into the water and is gone forever.\n", G.DirObj.Desc)
		return true
	}
	if G.ActVerb.Norm == "leap" || G.ActVerb.Norm == "through" {
		Printf("A look before leaping reveals that the river is wide and dangerous, with swift currents and large, half-hidden rocks. You decide to forgo your swim.\n")
		return true
	}
	return false
}

func DamRoomFcn(arg ActArg) bool {
	if arg == ActLook {
		Printf("You are standing on the top of the Flood Control Dam #3, which was quite a tourist attraction in times far distant. There are paths to the north, south, and west, and a scramble down.\n")
		if G.LowTide && G.GatesOpen {
			Printf("The water level behind the dam is low: The sluice gates have been opened. Water rushes through the dam and downstream.\n")
		} else if G.GatesOpen {
			Printf("The sluice gates are open, and water rushes through the dam. The water level behind the dam is still high.\n")
		} else if G.LowTide {
			Printf("The sluice gates are closed. The water level in the reservoir is quite low, but the level is rising quickly.\n")
		} else {
			Printf("The sluice gates on the dam are closed. Behind the dam, there can be seen a wide reservoir. Water is pouring over the top of the now abandoned dam.\n")
		}
		Printf("There is a control panel here, on which a large metal bolt is mounted. Directly above the bolt is a small green plastic bubble")
		if G.GateFlag {
			Printf(" which is glowing serenely")
		}
		Printf(".\n")
		return true
	}
	return false
}

func WhiteCliffsFcn(arg ActArg) bool {
	if arg == ActEnd {
		if InflatedBoat.IsIn(G.Winner) {
			G.Deflate = false
		} else {
			G.Deflate = true
		}
	}
	return false
}

func FallsRoomFcn(arg ActArg) bool {
	if arg == ActLook {
		Printf("You are at the top of Aragain Falls, an enormous waterfall with a drop of about 450 feet. The only path here is on the north end.\n")
		if G.RainbowFlag {
			Printf("A solid rainbow spans the falls.\n")
		} else {
			Printf("A beautiful rainbow can be seen over the falls and to the west.\n")
		}
		return true
	}
	return false
}

func Rivr4RoomFcn(arg ActArg) bool {
	if arg == ActEnd {
		if Buoy.IsIn(G.Winner) && G.BuoyFlag {
			Printf("You notice something funny about the feel of the buoy.\n")
			G.BuoyFlag = false
		}
	}
	return false
}

func RBoatFcn(arg ActArg) bool {
	if arg == ActEnter || arg == ActEnd || arg == ActLook {
		return false
	}
	if arg == ActBegin {
		if G.ActVerb.Norm == "walk" && G.Params.HasWalkDir {
			if G.Params.WalkDir == Land || G.Params.WalkDir == East || G.Params.WalkDir == West {
				return false
			}
			if G.Here == &Reservoir && (G.Params.WalkDir == North || G.Params.WalkDir == South) {
				return false
			}
			if G.Here == &InStream && G.Params.WalkDir == South {
				return false
			}
			Printf("Read the label for the boat's instructions.\n")
			return true
		}
		if G.ActVerb.Norm == "launch" {
			if G.Here == &River1 || G.Here == &River2 || G.Here == &River3 || G.Here == &River4 || G.Here == &Reservoir || G.Here == &InStream {
				Printf("You are on the ")
				if G.Here == &Reservoir {
					Printf("reservoir")
				} else if G.Here == &InStream {
					Printf("stream")
				} else {
					Printf("river")
				}
				Printf(", or have you forgotten?\n")
				return true
			}
			tmp := GoNext(RiverLaunch)
			if tmp == 1 {
				// ZIL: <ENABLE <QUEUE I-RIVER <LKP ,HERE ,RIVER-SPEEDS>>>
				// After GoNext, HERE is the destination room. Use its speed.
				if spd, ok := RiverSpeedMap[G.Here]; ok {
					Queue("IRiver", spd).Run = true
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
			(G.ActVerb.Norm == "put" && G.DirObj.Has(FlgWeapon) && G.IndirObj == &InflatedBoat) ||
			((G.ActVerb.Norm == "attack" || G.ActVerb.Norm == "mung") && G.IndirObj != nil && G.IndirObj.Has(FlgWeapon)) {
			RemoveCarefully(&InflatedBoat)
			PuncturedBoat.MoveTo(G.Here)
			Rob(&InflatedBoat, G.Here, 0)
			G.Winner.MoveTo(G.Here)
			Printf("It seems that the ")
			if G.ActVerb.Norm == "drop" || G.ActVerb.Norm == "put" {
				Printf("%s", G.DirObj.Desc)
			} else {
				Printf("%s", G.IndirObj.Desc)
			}
			Printf(" didn't agree with the boat, as evidenced by the loud hissing noise issuing therefrom. With a pathetic sputter, the boat deflates, leaving you without.\n")
			if G.Here.Has(FlgNonLand) {
				Printf("\n")
				if G.Here == &Reservoir || G.Here == &InStream {
					JigsUp("Another pathetic sputter, this time from you, heralds your drowning.", false)
				} else {
					JigsUp("In other words, fighting the fierce currents of the Frigid River. You manage to hold your own for a bit, but then you are carried over a waterfall and into some nasty rocks. Ouch!", false)
				}
			}
			return true
		}
		return false
	}
	if G.ActVerb.Norm == "board" {
		if Sceptre.IsIn(G.Winner) || Knife.IsIn(G.Winner) || Sword.IsIn(G.Winner) || RustyKnife.IsIn(G.Winner) || Axe.IsIn(G.Winner) || Stiletto.IsIn(G.Winner) {
			Printf("Oops! Something sharp seems to have slipped and punctured the boat. The boat deflates to the sounds of hissing, sputtering, and cursing.\n")
			RemoveCarefully(&InflatedBoat)
			PuncturedBoat.MoveTo(G.Here)
			ThisIsIt(&PuncturedBoat)
			return true
		}
		return false
	}
	if G.ActVerb.Norm == "inflate" || G.ActVerb.Norm == "fill" {
		Printf("Inflating it further would probably burst it.\n")
		return true
	}
	if G.ActVerb.Norm == "deflate" {
		if G.Winner.Location() == &InflatedBoat {
			Printf("You can't deflate the boat while you're in it.\n")
			return true
		}
		if !InflatedBoat.IsIn(G.Here) {
			Printf("The boat must be on the ground to be deflated.\n")
			return true
		}
		Printf("The boat deflates.\n")
		G.Deflate = true
		RemoveCarefully(&InflatedBoat)
		InflatableBoat.MoveTo(G.Here)
		ThisIsIt(&InflatableBoat)
		return true
	}
	return false
}

func IRfill() bool {
	Reservoir.Give(FlgNonLand)
	Reservoir.Take(FlgRLand)
	DeepCanyon.Take(FlgTouch)
	LoudRoom.Take(FlgTouch)
	if Trunk.IsIn(&Reservoir) {
		Trunk.Give(FlgInvis)
	}
	G.LowTide = false
	if G.Here == &Reservoir {
		if G.Winner.Location().Has(FlgVeh) {
			Printf("The boat lifts gently out of the mud and is now floating on the reservoir.\n")
		} else {
			JigsUp("You are lifted up by the rising river! You try to swim, but the currents are too strong. You come closer, closer to the awesome structure of Flood Control Dam #3. The dam beckons to you. The roar of the water nearly deafens you, but you remain conscious as you tumble over the dam toward your certain doom among the rocks at its base.", false)
		}
	} else if G.Here == &DeepCanyon {
		Printf("A sound, like that of flowing water, starts to come from below.\n")
	} else if G.Here == &LoudRoom {
		Printf("All of a sudden, an alarmingly loud roaring sound fills the room. Filled with fear, you scramble away.\n")
		dest := LoudRuns[G.Rand.Intn(len(LoudRuns))]
		Goto(dest, true)
	} else if G.Here == &ReservoirNorth || G.Here == &ReservoirSouth {
		Printf("You notice that the water level has risen to the point that it is impossible to cross.\n")
	}
	return true
}

func IRempty() bool {
	Reservoir.Give(FlgRLand)
	Reservoir.Take(FlgNonLand)
	DeepCanyon.Take(FlgTouch)
	LoudRoom.Take(FlgTouch)
	Trunk.Take(FlgInvis)
	G.LowTide = true
	if G.Here == &Reservoir && G.Winner.Location().Has(FlgVeh) {
		Printf("The water level has dropped to the point at which the boat can no longer stay afloat. It sinks into the mud.\n")
	} else if G.Here == &DeepCanyon {
		Printf("The roar of rushing water is quieter now.\n")
	} else if G.Here == &ReservoirNorth || G.Here == &ReservoirSouth {
		Printf("The water level is now quite low here and you could easily cross over to the other side.\n")
	}
	return true
}

func IMaintRoom() bool {
	hereQ := G.Here == &MaintenanceRoom
	if hereQ {
		Printf("The water level here is now ")
		idx := G.WaterLevel / 2
		if idx >= 0 && idx < len(Drownings) {
			Printf("%s", Drownings[idx])
		}
		Printf("\n")
	}
	G.WaterLevel++
	if G.WaterLevel >= 14 {
		MungRoom(&MaintenanceRoom, "The room is full of water and cannot be entered.")
		QueueInt("IMaintRoom", false).Run = false
		if hereQ {
			JigsUp("I'm afraid you have done drowned yourself.", false)
		}
	} else if InflatedBoat.IsIn(G.Winner) && (G.Here == &MaintenanceRoom || G.Here == &DamRoom || G.Here == &DamLobby) {
		JigsUp("The rising water carries the boat over the dam, down the river, and over the falls. Tsk, tsk.", false)
	}
	return true
}

func IRiver() bool {
	if G.Here != &River1 && G.Here != &River2 && G.Here != &River3 && G.Here != &River4 && G.Here != &River5 {
		QueueInt("IRiver", false).Run = false
		return false
	}
	rm, ok := RiverNext[G.Here]
	if ok {
		Printf("The flow of the river carries you downstream.\n\n")
		Goto(rm, true)
		// ZIL: <ENABLE <QUEUE I-RIVER <LKP ,HERE ,RIVER-SPEEDS>>>
		// After Goto, HERE is the new room. Use its speed for the next leg.
		if spd, ok := RiverSpeedMap[G.Here]; ok {
			Queue("IRiver", spd).Run = true
		}
		return true
	}
	JigsUp("Unfortunately, the magic boat doesn't provide protection from the rocks and boulders one meets at the bottom of waterfalls. Including this one.", false)
	return true
}
