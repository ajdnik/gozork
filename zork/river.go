package zork

import "math/rand"


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
