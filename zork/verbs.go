package zork

import (
	"encoding/binary"
	"os"
)

var (
	Moves      = 0
	Lit        = false
	SuperBrief = false
	Version    = [24]byte{0, 0, 0, 0, 88, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 56, 52, 48, 55, 50, 54}
)

func VVersion(arg ActArg) bool {
	Print("ZORK I: The Great Underground Empire", Newline)
	Print("Infocom interactive fiction - a fantasy story", Newline)
	Print("Copyright (c) 1981, 1982, 1983, 1984, 1985, 1986", Newline)
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
	return false
}

func VWalk(arg ActArg) bool {
	return false
}

func VQuit(arg ActArg) bool {
	VScore(arg)
	Print("Do you wish to restart? (Y is affirmative): ", NoNewline)
	if IsYes() {
		os.Exit(0)
	} else {
		Print("Ok.", Newline)
	}
	return true
}

func VScore(arg ActArg) bool {
	return false
}

func IsYes() bool {
	Print(">", NoNewline)
	_, lex := Read()
	if len(lex) > 0 && lex[0].IsAny("yes", "y") {
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

func ITake(vb bool) bool {
	// TODO: gverbs.zil:1900
	return false
}

func ToDirObj(dir string) *Object {
	return nil
}
