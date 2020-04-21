package zork

func SetLastObject(o *Object) *Object {
	LastNoun = o
	LastNounPlace = Location
	return LastNounPlace
}

func IsRoom(o *Object) bool {
	for _, r := range Rooms {
		if o.Name == r.Name {
			return true
		}
	}
	return false
}

func DescribeRoom(look bool) bool {
	if Location == nil {
		return false
	}
	v := Verbose
	if look {
		v = look
	}
	if !Lit {
		Print("It is pitch black.")
		if !GrueRepellent {
			Print(" You are likely to be eaten by a grue.")
		}
		NewLine()
		return false
	}
	if !Location.Has(Visited) {
		Location.Give(Visited)
		v = true
	}
	if Location.Has(MazeRoom) {
		Location.Take(Visited)
	}
	environment := Player.Parent
	if IsRoom(Location) {
		PrintObj(Location)
		if environment != nil && environment.Has(Vehicle) {
			Print(", in the ")
			PrintObj(environment)
		}
		NewLine()
	}
	if !look && Superbrief {
		return true
	}
	if v {
		if Location.Action != nil && Location.Action(MLook) {
			return true
		}
		if Location.Description != NoDescription {
			Print(Location.Description)
			NewLine()
		}
	} else if Location.Action != nil {
		Location.Action(MWake)
	}
	if environment != nil {
		if environment.Name == Location.Name {
			return true
		}
		if !environment.Has(Vehicle) {
			return true
		}
	}
	if Location.Action != nil {
		Location.Action(MLook)
	}
	return true
}

func DescribeObjects(vrb bool) bool {
	if Location == nil {
		return false
	}
	if !Lit {
		Print("Only bats can see in the dark. And you're not one.")
		return true
	}
	if Location.Contains == nil || len(Location.Contains) == 0 {
		return false
	}
	v := Verbose
	if vrb {
		v = vrb
	}
	return PrintCont(Location, v, -1)
}

func CanSeeContents(obj *Object) bool {
	if obj.Has(Concealed) {
		return false
	}
	if obj.Has(Transparent) || obj.Has(Open) {
		return true
	}
	return false
}

func Firster(obj *Object, level int) bool {
	if Player == nil {
		return false
	}
	if obj == nil {
		return false
	}
	if obj.Name == TrophyCase.Name {
		Print("Your collection of treasures consists of:")
		return true
	}
	if obj.Name == Player.Name {
		Print("You are carrying:")
		return true
	}
	if IsRoom(obj) {
		return false
	}
	if level > 0 {
		Print(Indents[level])
	}
	if obj.Has(Supporter) {
		Print("Sitting on the ")
		PrintObj(obj)
		Print(" is: ")
		return true
	}
	if obj.Has(Animate) {
		Print("The ")
		PrintObj(obj)
		Print(" is holding: ")
		return true
	}
	Print("The ")
	PrintObj(obj)
	Print(" contains:")
	return true
}

func DescribeOb(obj *Object, v bool, level int) bool {
	LastObLongdesc = obj
	if level == 0 {
		if obj.Initial2 != nil && obj.Initial2(MFight) {
			return true
		}
		if !obj.Has(Visited) && obj.Initial != NoInitial {
			Print(obj.Initial)
		} else if obj.Description != NoDescription {
			Print(obj.Description)
		} else {
			Print("There is a ")
			PrintObj(obj)
			Print(" here")
			if obj.Has(Light) {
				Print(" (providing light)")
			}
			Print(".")
		}
	} else {
		Print(Indents[level])
		Print("A ")
		PrintObj(obj)
		if obj.Has(Light) {
			Print(" (providing light)")
		} else if obj.Has(Clothing) {
			Print(" (being worn)")
		}
	}
	if level == 0 {
		av := Player.Parent
		if av != nil && av.Has(Vehicle) {
			Print(" (outside the ")
			PrintObj(av)
			Print(")")
		}
	}
	NewLine()
	if !CanSeeContents(obj) {
		return false
	}
	if obj.Contains == nil || len(obj.Contains) == 0 {
		return false
	}
	return PrintCont(obj, v, level)
}

func PrintCont(obj *Object, v bool, level int) bool {
	if Player == nil {
		return false
	}
	if obj == nil {
		return false
	}
	if obj.Contains == nil || len(obj.Contains) == 0 {
		return true
	}
	environment := Player.Parent
	if environment != nil && !environment.Has(Vehicle) {
		environment = nil
	}
	pv := false
	av := true
	first := true
	inv := false
	if Player.Name != obj.Name && obj.Parent != nil && obj.Parent.Name != Player.Name {
		for _, itm := range obj.Contains {
			if itm.Name == environment.Name {
				pv = true
				continue
			}
			if itm.Name == Player.Name {
				continue
			}
			if itm.Has(Concealed) || itm.Has(Visited) {
				continue
			}
			if itm.Initial == NoInitial {
				continue
			}
			if !itm.Has(Scenery) {
				Print(itm.Initial)
				NewLine()
				av = false
			}
			if !CanSeeContents(itm) || (itm.Parent != nil && itm.Parent.Initial2 != nil) || itm.Contains == nil || len(itm.Contains) == 0 || !PrintCont(itm, v, 0) {
				continue
			}
			first = false
		}
	} else {
		inv = true
	}
	for _, itm := range obj.Contains {
		if environment != nil && itm.Name == environment.Name {
			continue
		}
		if Item4 != nil && itm.Name == Item4.Name {
			continue
		}
		if itm.Has(Concealed) {
			continue
		}
		if !inv && !itm.Has(Visited) && itm.Initial != NoInitial {
			continue
		}
		if itm.Has(Scenery) {
			if itm.Contains != nil && len(itm.Contains) > 0 && CanSeeContents(itm) {
				level++
				PrintCont(itm, v, level)
			}
			continue
		}
		if first {
			if Firster(obj, level) && level < 0 {
				level = 0
			}
			level++
			first = false
		}
		if level < 0 {
			level = 0
		}
		DescribeOb(itm, v, level)
	}
	if pv && environment != nil && environment.Contains != nil && len(environment.Contains) > 0 {
		level++
		PrintCont(environment, v, level)
	}
	if !first {
		return true
	}
	if !av {
		return true
	}
	return false
}

func CommandLoop() {

}

func StartGame() {
	for {
		Queue("FightDaemon", FightDaemon, -1).Enabled = true
		Queue("SwordDaemon", SwordDaemon, -1)
		Queue("ThiefDaemon", ThiefDaemon, -1).Enabled = true
		Queue("CandleDaemon", CandleDaemon, 40)
		Queue("LampDaemon", LampDaemon, 200)
		MagicBoat.VType = 8
		Def1Res[1] = Def1[2]
		Def1Res[2] = Def1[4]
		Def2Res[2] = Def2B[2]
		Def2Res[3] = Def2B[4]
		Def3Res[1] = Def3A[2]
		Def3Res[3] = Def3B[2]
		Location = WestOfHouse
		SetLastObject(SmallMailbox)
		if !Location.Has(Visited) {
			VersionSub()
			NewLine()
		}
		Lit = true
		Player = Cretin
		Actor = Player
		Player.MoveTo(Location)
		LookSub()
		CommandLoop()
		break
	}
}
