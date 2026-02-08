package zork

import "encoding/binary"


func VVerbose(arg ActArg) bool {
	G.Verbose = true
	G.SuperBrief = false
	Print("Maximum verbosity.", Newline)
	return true
}

func VBrief(arg ActArg) bool {
	G.Verbose = false
	G.SuperBrief = false
	Print("Brief descriptions.", Newline)
	return true
}

func VSuperBrief(arg ActArg) bool {
	G.SuperBrief = true
	Print("Superbrief descriptions.", Newline)
	return true
}

func VInventory(arg ActArg) bool {
	if G.Winner.HasChildren() {
		return PrintCont(G.Winner, false, 0)
	}
	Print("You are empty-handed.", Newline)
	return true
}

func VRestart(arg ActArg) bool {
	VScore(ActUnk)
	Print("Do you wish to restart? (Y is affirmative): ", NoNewline)
	if IsYes() {
		Print("Restarting.", Newline)
		if G.Restart() {
			VVersion(ActUnk)
			NewLine()
			VFirstLook(ActUnk)
			return true
		}
		Print("Failed.", Newline)
		return true
	}
	return false
}

func VRestore(arg ActArg) bool {
	if G.Restore() {
		Print("Ok.", Newline)
		return VFirstLook(ActUnk)
	}
	Print("Failed.", Newline)
	return true
}

func VSave(arg ActArg) bool {
	if G.Save() {
		Print("Ok.", Newline)
		return true
	}
	Print("Failed.", Newline)
	return true
}

func VScript(arg ActArg) bool {
	// This code turns on the first bit in the 8th word from the beginning
	// Put(0, 8, Bor(Get(0, 8), 1))
	G.Script = true
	Print("Here begins a transcript of interaction with", Newline)
	VVersion(ActUnk)
	return true
}

func VUnscript(arg ActArg) bool {
	Print("Here ends a transcript of interaction with", Newline)
	VVersion(ActUnk)
	// This code turns off the first bit in the 8th word from the beginning
	// Put(0, 8, Band(Get(0, 8), -2))
	G.Script = false
	return true
}

func VVerify(arg ActArg) bool {
	Print("Verifying disk...", Newline)
	if Verify() {
		Print("The disk is correct.", Newline)
		return true
	}
	NewLine()
	Print("** Disk Failure **", Newline)
	return true
}

// VCommandFile switches input to come from a file. In the original Z-machine
// this uses the DIRIN opcode. Not applicable to this Go port.
func VCommandFile(arg ActArg) bool {
	return true
}

// VRandom reseeds the random number generator. In the original Z-machine
// this uses the RANDOM opcode with a negative argument to set the seed.
func VRandom(arg ActArg) bool {
	Print("Illegal call to #RND.", Newline)
	return true
}

// VRecord starts recording input to a file. In the original Z-machine
// this uses the DIROUT opcode. Not applicable to this Go port.
func VRecord(arg ActArg) bool {
	return true
}

// VUnrecord stops recording input. In the original Z-machine this uses
// the DIROUT opcode. Not applicable to this Go port.
func VUnrecord(arg ActArg) bool {
	return true
}

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

func VQuit(arg ActArg) bool {
	VScore(arg)
	Print("Do you wish to leave the game? (Y is affirmative): ", NoNewline)
	if IsYes() {
		Quit()
	} else {
		Print("Ok.", Newline)
	}
	return true
}

func VBug(arg ActArg) bool {
	Print("Bug? Not in a flawless program like this! (Cough, cough).", Newline)
	return true
}

func VWait(arg ActArg) bool {
	return Wait(3)
}

func Wait(num int) bool {
	Print("Time passes...", Newline)
	for i := num; i > 0; i-- {
		Clocker()
	}
	G.ClockWait = true
	return true
}

func VWin(arg ActArg) bool {
	Print("Naturally!", Newline)
	return true
}

func VWish(arg ActArg) bool {
	Print("With luck, your wish will come true.", Newline)
	return true
}

func VZork(arg ActArg) bool {
	Print("At your service!", Newline)
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

func Finish() bool {
	VScore(ActUnk)
	for {
		NewLine()
		Print("Would you like to restart the game from the beginning, restore a saved game position, or end this session of the game?", Newline)
		Print("(Type RESTART, RESTORE, or QUIT):", Newline)
		Print(">", NoNewline)
		_, lex := Read()
		if len(lex) == 0 {
			continue
		}
		wrd := lex[0]
		if wrd.Norm == "restart" {
			if !G.Restart() {
				Print("Failed.", Newline)
			}
			return true
		}
		if wrd.Norm == "restore" {
			if G.Restore() {
				Print("Ok.", Newline)
				return true
			}
			Print("Failed.", Newline)
			return true
		}
		if wrd.Norm == "quit" {
			Quit()
		}
	}
}

func VAdvent(arg ActArg) bool {
	Print("A hollow voice says \"Fool.\"", Newline)
	return true
}
