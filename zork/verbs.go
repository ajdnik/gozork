package zork

import (
	"encoding/binary"
	"os"
)

type RndSelect struct {
	Unselected []string
	Selected   []string
}

var (
	Moves      int
	Score      int
	BaseScore  int
	Lit        bool
	SuperBrief bool
	Verbose    bool
	Dead       bool
	// IsSprayed is a flag indicating if the player is wearing Grue repellent
	IsSprayed bool
	Version   = [24]byte{0, 0, 0, 0, 119, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 56, 56, 48, 52, 50, 57}
	WonGame   bool
	// DescObj stores the last object which was described
	DescObj *Object
	Indents = [6]string{
		"",
		"  ",
		"    ",
		"      ",
		"        ",
		"          ",
	}
	Yuks = RndSelect{
		Unselected: []string{
			"A valiant attempt.",
			"You can't be serious.",
			"An interesting idea...",
			"What a concept!",
		},
	}
	Lucky        = true
	FumbleNumber = 7
	FumbleProb   = 8
)

func VVersion(arg ActArg) bool {
	Print("ZORK I: The Great Underground Empire", Newline)
	Print("Infocom interactive fiction - a fantasy story", Newline)
	Print("Copyright (c) 1981, 1982, 1983, 1984, 1985, 1986", NoNewline)
	Print(" Infocom, Inc. All rights reserved.", Newline)
	Print("ZORK is a registered trademark of Infocom, Inc.", Newline)
	Print("Release ", NoNewline)
	num := binary.LittleEndian.Uint16(Version[4:6])
	PrintNumber(int(num & 2047))
	Print(" / Serial number ", NoNewline)
	for offset := 18; offset <= 23; offset++ {
		Print(string(Version[offset]), NoNewline)
	}
	NewLine()
	return true
}

func VLook(arg ActArg) bool {
	if DescribeRoom(true) {
		return DescribeObjects(true)
	}
	return false
}

func VLookBehind(arg ActArg) bool {
	Print("There is nothing behind the ", NoNewline)
	PrintObject(DirObj)
	Print(".", Newline)
	return true
}

func VLookInside(arg ActArg) bool {
	if DirObj.Has(FlgDoor) {
		Print("The ", NoNewline)
		PrintObject(DirObj)
		if DirObj.Has(FlgOpen) {
			Print(" is open, but I can't tell what's beyond it.", NoNewline)
		} else {
			Print(" is closed.", NoNewline)
		}
		NewLine()
		return true
	}
	if DirObj.Has(FlgCont) {
		if DirObj.Has(FlgPerson) {
			Print("There is nothing special to be seen.", Newline)
			return true
		}
		if !CanSeeInside(DirObj) {
			Print("The ", NoNewline)
			PrintObject(DirObj)
			Print(" is closed.", Newline)
			return true
		}
		if DirObj.HasChildren() && PrintCont(DirObj, false, 0) {
			return true
		}
		Print("The ", NoNewline)
		PrintObject(DirObj)
		Print(" is empty.", Newline)
		return true
	}
	Print("You can't look inside a ", NoNewline)
	PrintObject(DirObj)
	Print(".", Newline)
	return true
}

func VLookOn(arg ActArg) bool {
	if DirObj.Has(FlgSurf) {
		Perform("look inside", DirObj, nil)
		return true
	}
	Print("Look on a ", NoNewline)
	PrintObject(DirObj)
	Print("???", Newline)
	return true
}

func VLookUnder(arg ActArg) bool {
	Print("There is nothing but dust there.", Newline)
	return true
}

func VExamine(arg ActArg) bool {
	if len(DirObj.Text) > 0 {
		Print(DirObj.Text, Newline)
		return true
	}
	if DirObj.Has(FlgCont) || DirObj.Has(FlgDoor) {
		return VLookInside(ActUnk)
	}
	Print("There's nothing special about the ", NoNewline)
	PrintObject(DirObj)
	Print(".", Newline)
	return true
}

func VWalk(arg ActArg) bool {
	return false
}

func VQuit(arg ActArg) bool {
	VScore(arg)
	Print("Do you wish to leave the game? (Y is affirmative): ", NoNewline)
	if IsYes() {
		os.Exit(0)
	} else {
		Print("Ok.", Newline)
	}
	return true
}

func IsYes() bool {
	Print(">", NoNewline)
	_, lex := Read()
	if len(lex) > 0 && lex[0].IsAny("yes", "y") {
		return true
	}
	return false
}

func DescribeRoom(isLook bool) bool {
	v := Verbose
	if isLook {
		v = isLook
	}
	if !Lit {
		Print("It is pitch black.", NoNewline)
		if !IsSprayed {
			Print(" You are likely to be eaten by a grue.", NoNewline)
		}
		NewLine()
		return false
	}
	if !Here.Has(FlgTouch) {
		Here.Give(FlgTouch)
		v = true
	}
	if Here.Has(FlgMaze) {
		Here.Take(FlgTouch)
	}
	if Here.IsIn(&Rooms) {
		PrintObject(Here)
		if av := Winner.Location(); av.Has(FlgVeh) {
			Print(", in the ", NoNewline)
			PrintObject(av)
		}
		NewLine()
	}
	if !isLook && SuperBrief {
		return true
	}
	av := Winner.Location()
	if v && Here.Action != nil && Here.Action(ActLook) {
		return true
	}
	if v && len(Here.LongDesc) != 0 {
		Print(Here.LongDesc, Newline)
	} else if Here.Action != nil {
		Here.Action(ActFlash)
	}
	if Here != av && av.Has(FlgVeh) && av.Action != nil {
		av.Action(ActLook)
	}
	return true
}

func DescribeObjects(v bool) bool {
	if !Lit {
		Print("Only bats can see in the dark. And you're not one.", Newline)
		return true
	}
	if !Here.HasChildren() {
		return false
	}
	if !v {
		v = Verbose
	}
	return PrintCont(Here, v, -1)
}

func PrintCont(obj *Object, v bool, lvl int) bool {
	if !obj.HasChildren() {
		return true
	}
	var av *Object
	if Winner.Location().Has(FlgVeh) {
		av = Winner.Location()
	}
	isInv := false
	isPv := false
	shit := true
	isFirst := true
	if Winner == obj || Winner == obj.Location() {
		isInv = true
	} else {
		for _, child := range obj.Children {
			if child == av {
				isPv = true
			} else if Winner == child {
				continue
			} else if !child.Has(FlgInvis) && !child.Has(FlgTouch) && len(child.FirstDesc) > 0 {
				if !child.Has(FlgNoDesc) {
					Print(child.FirstDesc, Newline)
					shit = false
				}
				if CanSeeInside(child) && child.Location().DescFcn == nil && child.HasChildren() {
					if PrintCont(child, v, 0) {
						isFirst = false
					}
				}
			}
		}
	}
	for _, child := range obj.Children {
		if child == av || child == &Adventurer {
			continue
		}
		if child.Has(FlgInvis) || (!isInv && !child.Has(FlgTouch) && len(child.FirstDesc) > 0) {
			continue
		}
		if !child.Has(FlgNoDesc) {
			if isFirst {
				if Firster(obj, lvl) && lvl < 0 {
					lvl = 0
				}
				lvl++
				isFirst = false
			}
			if lvl < 0 {
				lvl = 0
			}
			DescribeObject(child, v, lvl)
		} else if child.HasChildren() && CanSeeInside(child) {
			lvl++
			PrintCont(child, v, lvl)
			lvl--
		}
	}
	if isPv && av != nil && av.HasChildren() {
		lvl++
		PrintCont(av, v, lvl)
	}
	if isFirst && shit {
		return false
	}
	return true
}

func Firster(obj *Object, lvl int) bool {
	if obj == &TrophyCase {
		Print("Your collection of treasures consists of:", Newline)
		return true
	}
	if obj == Winner {
		Print("You are carrying:", Newline)
		return true
	}
	if obj.IsIn(&Rooms) {
		return false
	}
	if lvl > 0 {
		Print(Indents[lvl], NoNewline)
	}
	if obj.Has(FlgSurf) {
		Print("Sitting on the ", NoNewline)
		PrintObject(obj)
		Print(" is: ", Newline)
		return true
	}
	if obj.Has(FlgPerson) {
		Print("The ", NoNewline)
		PrintObject(obj)
		Print(" is holding: ", Newline)
		return true
	}
	Print("The ", NoNewline)
	PrintObject(obj)
	Print(" contains:", Newline)
	return true
}

func DescribeObject(obj *Object, v bool, lvl int) bool {
	DescObj = obj
	if lvl == 0 && obj.DescFcn != nil && obj.DescFcn(ActObjDesc) {
		return true
	}
	if lvl == 0 && ((!obj.Has(FlgTouch) && len(obj.FirstDesc) > 0) || len(obj.LongDesc) > 0) {
		if !obj.Has(FlgTouch) && len(obj.FirstDesc) > 0 {
			Print(obj.FirstDesc, NoNewline)
		} else {
			Print(obj.LongDesc, NoNewline)
		}
	} else if lvl == 0 {
		Print("There is a ", NoNewline)
		PrintObject(obj)
		Print(" here", NoNewline)
		if obj.Has(FlgOn) {
			Print(" (providing light)", NoNewline)
		}
		Print(".", NoNewline)
	} else {
		Print(Indents[lvl], NoNewline)
		Print("A ", NoNewline)
		PrintObject(obj)
		if obj.Has(FlgOn) {
			Print(" (providing light)", NoNewline)
		} else if obj.Has(FlgWear) && obj.IsIn(Winner) {
			Print(" (being worn)", NoNewline)
		}
	}
	if av := Winner.Location(); lvl == 0 && av != nil && av.Has(FlgVeh) {
		Print(" (outside the ", NoNewline)
		PrintObject(av)
		Print(")", NoNewline)
	}
	NewLine()
	if CanSeeInside(obj) && obj.HasChildren() {
		return PrintCont(obj, v, lvl)
	}
	return true
}

func CanSeeInside(obj *Object) bool {
	if obj.Has(FlgInvis) {
		return false
	}
	if obj.Has(FlgTrans) || obj.Has(FlgOpen) {
		return true
	}
	return false
}

func ThisIsIt(obj *Object) {
	Params.ItObj = obj
}

func IsInGlobal(obj1, obj2 *Object) bool {
	if obj2.Global == nil {
		return false
	}
	for _, o := range obj2.Global {
		if o == obj1 {
			return true
		}
	}
	return false
}

func IsHeld(obj *Object) bool {
	for {
		obj := obj.Location()
		if obj == nil {
			return false
		}
		if obj == Winner {
			return true
		}
	}
}

func ScoreUpd(num int) bool {
	BaseScore += num
	Score += num
	if Score == 350 && !WonGame {
		WonGame = true
		Map.Take(FlgInvis)
		WestOfHouse.Take(FlgTouch)
		Print("An almost inaudible voice whispers in your ear, \"Look to your treasures for the final secret.\"", Newline)
	}
	return true
}

func ScoreObj(obj *Object) {
	if obj.Value <= 0 {
		return
	}
	ScoreUpd(obj.Value)
	obj.Value = 0
}

func CCount(obj *Object) int {
	cnt := 0
	for _, child := range obj.Children {
		if !child.Has(FlgWear) {
			cnt++
		}
	}
	return cnt
}

func Weight(obj *Object) int {
	wt := 0
	for _, child := range obj.Children {
		if obj == Player && child.Has(FlgWear) {
			wt++
		} else {
			wt += Weight(child)
		}
	}
	return wt + obj.Size
}

func ITake(vb bool) bool {
	if Dead {
		if vb {
			Print("Your hand passes through its object.", Newline)
		}
		return false
	}
	if !DirObj.Has(FlgTake) {
		if vb {
			Print(PickOne(Yuks), Newline)
		}
		return false
	}
	if DirObj.Has(FlgCont) && !DirObj.Has(FlgOpen) {
		return false
	}
	if !DirObj.Location().IsIn(Winner) && Weight(DirObj)+Weight(Winner) > LoadAllowed {
		if vb {
			Print("Your load is too heavy", NoNewline)
			if LoadAllowed < LoadMax {
				Print(", especially in light of your condition.", NoNewline)
			} else {
				Print(".", NoNewline)
			}
			NewLine()
		}
		// TODO: rfatal, not rfalse
		return false
	}
	cnt := CCount(Winner)
	if ActVerb == "tell" && cnt > FumbleNumber && Prob(cnt*FumbleProb, false) {
		Print("You're holding too many things already!", Newline)
		return false
	}
	DirObj.MoveTo(Winner)
	DirObj.Take(FlgNoDesc)
	DirObj.Give(FlgTouch)
	ScoreObj(DirObj)
	return true
}

func ToDirObj(dir string) *Object {
	return nil
}
