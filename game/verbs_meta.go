package game

import (
	"encoding/binary"
	. "github.com/ajdnik/gozork/engine"
)

func vVerbose(arg ActionArg) bool {
	G.Verbose = true
	G.SuperBrief = false
	Printf("Maximum verbosity.\n")
	return true
}

func vBrief(arg ActionArg) bool {
	G.Verbose = false
	G.SuperBrief = false
	Printf("Brief descriptions.\n")
	return true
}

func vSuperBrief(arg ActionArg) bool {
	G.SuperBrief = true
	Printf("Superbrief descriptions.\n")
	return true
}

func vInventory(arg ActionArg) bool {
	if G.Winner.HasChildren() {
		return printCont(G.Winner, false, 0)
	}
	Printf("You are empty-handed.\n")
	return true
}

func vRestart(arg ActionArg) bool {
	vScore(ActUnk)
	Printf("Do you wish to restart? (Y is affirmative): ")
	if isYes() {
		Printf("Restarting.\n")
		if err := G.Restart(); err != nil {
			Printf("Failed: %s\n", err)
			return true
		}
		vVersion(ActUnk)
		Printf("\n")
		vFirstLook(ActUnk)
		return true
	}
	return false
}

func vRestore(arg ActionArg) bool {
	if err := G.Restore(); err != nil {
		Printf("Failed: %s\n", err)
		return true
	}
	Printf("Ok.\n")
	return vFirstLook(ActUnk)
}

func vSave(arg ActionArg) bool {
	if err := G.Save(); err != nil {
		Printf("Failed: %s\n", err)
		return true
	}
	Printf("Ok.\n")
	return true
}

func vScript(arg ActionArg) bool {
	// This code turns on the first bit in the 8th word from the beginning
	// Put(0, 8, Bor(Get(0, 8), 1))
	G.Script = true
	Printf("Here begins a transcript of interaction with\n")
	vVersion(ActUnk)
	return true
}

func vUnscript(arg ActionArg) bool {
	Printf("Here ends a transcript of interaction with\n")
	vVersion(ActUnk)
	// This code turns off the first bit in the 8th word from the beginning
	// Put(0, 8, Band(Get(0, 8), -2))
	G.Script = false
	return true
}

func vVerify(arg ActionArg) bool {
	Printf("Verifying disk...\n")
	if Verify() {
		Printf("The disk is correct.\n")
		return true
	}
	Printf("\n** Disk Failure **\n")
	return true
}

func vVersion(arg ActionArg) bool {
	Printf("ZORK I: The Great Underground Empire\nInfocom interactive fiction - a fantasy story\nCopyright (c) 1981, 1982, 1983, 1984, 1985, 1986 Infocom, Inc. All rights reserved.\nZORK is a registered trademark of Infocom, Inc.\nRelease ")
	num := binary.LittleEndian.Uint16(version[4:6])
	Printf("%d / Serial number ", int(num&2047))
	for offset := 18; offset <= 23; offset++ {
		Printf("%s", string(version[offset]))
	}
	Printf("\n")
	return true
}

func vQuit(arg ActionArg) bool {
	vScore(arg)
	Printf("Do you wish to leave the game? (Y is affirmative): ")
	if isYes() {
		Quit()
	} else {
		Printf("Ok.\n")
	}
	return true
}

func vBug(arg ActionArg) bool {
	Printf("Bug? Not in a flawless program like this! (Cough, cough).\n")
	return true
}

func vWait(arg ActionArg) bool {
	return wait(3)
}

func wait(num int) bool {
	Printf("Time passes...\n")
	for i := num; i > 0; i-- {
		Clocker()
		if G.QuitRequested {
			return true
		}
	}
	G.ClockWait = true
	return true
}

func vWin(arg ActionArg) bool {
	Printf("Naturally!\n")
	return true
}

func vWish(arg ActionArg) bool {
	Printf("With luck, your wish will come true.\n")
	return true
}

func vZork(arg ActionArg) bool {
	Printf("At your service!\n")
	return true
}

func isYes() bool {
	Printf(">")
	_, lex := Read()
	if len(lex) > 0 && lex[0].IsAny("yes", "y") {
		return true
	}
	return false
}

func finish() bool {
	vScore(ActUnk)
	for {
		Printf("\nWould you like to restart the game from the beginning, restore a saved game position, or end this session of the game?\n(Type RESTART, RESTORE, or QUIT):\n>")
		_, lex := Read()
		if len(lex) == 0 {
			if G.InputExhausted {
				Quit()
				return true
			}
			continue
		}
		wrd := lex[0]
		if wrd.Norm == "restart" {
			if err := G.Restart(); err != nil {
				Printf("Failed: %s\n", err)
			}
			return true
		}
		if wrd.Norm == "restore" {
			if err := G.Restore(); err != nil {
				Printf("Failed: %s\n", err)
				return true
			}
			Printf("Ok.\n")
			return true
		}
		if wrd.Norm == "quit" {
			Quit()
			return true
		}
	}
}

func vAdvent(arg ActionArg) bool {
	Printf("A hollow voice says \"Fool.\"\n")
	return true
}
