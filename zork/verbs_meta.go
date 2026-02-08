package zork

import "encoding/binary"


func VVerbose(arg ActArg) bool {
	G.Verbose = true
	G.SuperBrief = false
	Printf("Maximum verbosity.\n")
	return true
}

func VBrief(arg ActArg) bool {
	G.Verbose = false
	G.SuperBrief = false
	Printf("Brief descriptions.\n")
	return true
}

func VSuperBrief(arg ActArg) bool {
	G.SuperBrief = true
	Printf("Superbrief descriptions.\n")
	return true
}

func VInventory(arg ActArg) bool {
	if G.Winner.HasChildren() {
		return PrintCont(G.Winner, false, 0)
	}
	Printf("You are empty-handed.\n")
	return true
}

func VRestart(arg ActArg) bool {
	VScore(ActUnk)
	Printf("Do you wish to restart? (Y is affirmative): ")
	if IsYes() {
		Printf("Restarting.\n")
		if G.Restart() {
			VVersion(ActUnk)
			Printf("\n")
			VFirstLook(ActUnk)
			return true
		}
		Printf("Failed.\n")
		return true
	}
	return false
}

func VRestore(arg ActArg) bool {
	if G.Restore() {
		Printf("Ok.\n")
		return VFirstLook(ActUnk)
	}
	Printf("Failed.\n")
	return true
}

func VSave(arg ActArg) bool {
	if G.Save() {
		Printf("Ok.\n")
		return true
	}
	Printf("Failed.\n")
	return true
}

func VScript(arg ActArg) bool {
	// This code turns on the first bit in the 8th word from the beginning
	// Put(0, 8, Bor(Get(0, 8), 1))
	G.Script = true
	Printf("Here begins a transcript of interaction with\n")
	VVersion(ActUnk)
	return true
}

func VUnscript(arg ActArg) bool {
	Printf("Here ends a transcript of interaction with\n")
	VVersion(ActUnk)
	// This code turns off the first bit in the 8th word from the beginning
	// Put(0, 8, Band(Get(0, 8), -2))
	G.Script = false
	return true
}

func VVerify(arg ActArg) bool {
	Printf("Verifying disk...\n")
	if Verify() {
		Printf("The disk is correct.\n")
		return true
	}
	Printf("\n** Disk Failure **\n")
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
	Printf("Illegal call to #RND.\n")
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
	Printf("ZORK I: The Great Underground Empire\nInfocom interactive fiction - a fantasy story\nCopyright (c) 1981, 1982, 1983, 1984, 1985, 1986 Infocom, Inc. All rights reserved.\nZORK is a registered trademark of Infocom, Inc.\nRelease ")
	num := binary.LittleEndian.Uint16(Version[4:6])
	Printf("%d / Serial number ", int(num & 2047))
	for offset := 18; offset <= 23; offset++ {
		Printf("%s", string(Version[offset]))
	}
	Printf("\n")
	return true
}

func VQuit(arg ActArg) bool {
	VScore(arg)
	Printf("Do you wish to leave the game? (Y is affirmative): ")
	if IsYes() {
		Quit()
	} else {
		Printf("Ok.\n")
	}
	return true
}

func VBug(arg ActArg) bool {
	Printf("Bug? Not in a flawless program like this! (Cough, cough).\n")
	return true
}

func VWait(arg ActArg) bool {
	return Wait(3)
}

func Wait(num int) bool {
	Printf("Time passes...\n")
	for i := num; i > 0; i-- {
		Clocker()
	}
	G.ClockWait = true
	return true
}

func VWin(arg ActArg) bool {
	Printf("Naturally!\n")
	return true
}

func VWish(arg ActArg) bool {
	Printf("With luck, your wish will come true.\n")
	return true
}

func VZork(arg ActArg) bool {
	Printf("At your service!\n")
	return true
}

func IsYes() bool {
	Printf(">")
	_, lex := Read()
	if len(lex) > 0 && lex[0].IsAny("yes", "y") {
		return true
	}
	return false
}

func Finish() bool {
	VScore(ActUnk)
	for {
		Printf("\nWould you like to restart the game from the beginning, restore a saved game position, or end this session of the game?\n(Type RESTART, RESTORE, or QUIT):\n>")
		_, lex := Read()
		if len(lex) == 0 {
			continue
		}
		wrd := lex[0]
		if wrd.Norm == "restart" {
			if !G.Restart() {
				Printf("Failed.\n")
			}
			return true
		}
		if wrd.Norm == "restore" {
			if G.Restore() {
				Printf("Ok.\n")
				return true
			}
			Printf("Failed.\n")
			return true
		}
		if wrd.Norm == "quit" {
			Quit()
		}
	}
}

func VAdvent(arg ActArg) bool {
	Printf("A hollow voice says \"Fool.\"\n")
	return true
}
