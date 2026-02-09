package game

import (
	"os"
	"strings"
	"testing"

	. "github.com/ajdnik/gozork/engine"
)

func TestSaveRestoreRoundTrip(t *testing.T) {
	savePath := "/tmp/gozork-save-roundtrip.sav"
	defer os.Remove(savePath)

	setupTestGame(t, savePath+"\n")

	// Change some state to ensure restore actually reverts it.
	G.Score = 42
	G.Moves = 7
	G.Here = &kitchen
	G.Winner.MoveTo(G.Here)
	lamp.MoveTo(G.Winner)

	savedHere := G.Here
	savedScore := G.Score
	savedMoves := G.Moves

	if err := doSave(); err != nil {
		t.Fatalf("doSave failed: %v", err)
	}

	// Mutate state again before restore.
	G.Score = 99
	G.Moves = 999
	G.Here = &cellar
	G.Winner.MoveTo(G.Here)
	removeCarefully(&lamp)

	// Restore from the saved file.
	G.GameInput = strings.NewReader(savePath + "\n")
	G.Reader = nil
	if err := doRestore(); err != nil {
		t.Fatalf("doRestore failed: %v", err)
	}

	if G.Here != savedHere {
		t.Fatalf("expected Here to be %s after restore, got %s", savedHere.Desc, G.Here.Desc)
	}
	if G.Score != savedScore || G.Moves != savedMoves {
		t.Fatalf("expected score/moves %d/%d after restore, got %d/%d", savedScore, savedMoves, G.Score, G.Moves)
	}
	if !lamp.IsIn(G.Winner) {
		t.Fatalf("expected lamp to be held after restore")
	}
}

func TestRestartResetsState(t *testing.T) {
	setupTestGame(t, "")

	// Mutate state.
	G.Score = 123
	G.Moves = 55
	gD().GrateRevealed = true
	G.Here = &cellar
	G.Winner.MoveTo(G.Here)

	if err := doRestart(); err != nil {
		t.Fatalf("doRestart failed: %v", err)
	}

	if G.Here != &westOfHouse {
		t.Fatalf("expected restart to place player at West of House, got %s", G.Here.Desc)
	}
	if G.Score != 0 || G.Moves != 0 {
		t.Fatalf("expected score/moves reset to 0, got %d/%d", G.Score, G.Moves)
	}
	if gD().GrateRevealed {
		t.Fatalf("expected grate state reset")
	}
}

func TestRiverHelpersAndActions(t *testing.T) {
	out := setupTestGame(t, "")

	// ---- fixBoat ----
	inflatableBoat.MoveTo(&damBase)
	puncturedBoat.MoveTo(&damBase)
	fixBoat()
	if puncturedBoat.In != nil {
		t.Fatalf("expected punctured boat removed after fix")
	}
	if !inflatableBoat.IsIn(&damBase) {
		t.Fatalf("expected inflatable boat moved to dam base after fix")
	}

	// ---- fixMaintLeak ----
	QueueInt("iMaintRoom", false).Run = true
	gD().WaterLevel = 10
	fixMaintLeak()
	if gD().WaterLevel != -1 {
		t.Fatalf("expected water level reset, got %d", gD().WaterLevel)
	}
	if QueueInt("iMaintRoom", false).Run {
		t.Fatalf("expected maintenance leak queue disabled")
	}

	// ---- waterFcn: take with bottle open ----
	bottle.MoveTo(G.Winner)
	bottle.Give(FlgOpen)
	G.ActVerb = ActionVerb{Norm: "take"}
	G.DirObj = &water
	G.IndirObj = nil
	water.MoveTo(G.Here)
	out.Reset()
	if !waterFcn(ActUnk) {
		t.Fatalf("expected water take handled")
	}
	if !water.IsIn(&bottle) {
		t.Fatalf("expected water moved into bottle")
	}
	if !strings.Contains(out.String(), "bottle is now full") {
		t.Fatalf("expected bottle fill message")
	}

	// ---- waterFcn: drop with bottle closed ----
	bottle.Take(FlgOpen)
	G.ActVerb = ActionVerb{Norm: "drop"}
	out.Reset()
	water.MoveTo(&bottle)
	if !waterFcn(ActUnk) {
		t.Fatalf("expected water drop handled")
	}
	if !strings.Contains(out.String(), "bottle is closed") {
		t.Fatalf("expected bottle closed message")
	}

	// ---- waterFcn: throw ----
	bottle.Give(FlgOpen)
	water.MoveTo(G.Winner)
	G.ActVerb = ActionVerb{Norm: "throw"}
	out.Reset()
	if !waterFcn(ActUnk) {
		t.Fatalf("expected water throw handled")
	}
	if water.In != nil {
		t.Fatalf("expected water removed after throw")
	}

	// ---- boltFcn ----
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "turn"}
	G.IndirObj = &screwdriver
	boltFcn(ActUnk)
	if !strings.Contains(out.String(), "bolt won't turn") {
		t.Fatalf("expected bolt tool rejection")
	}

	out.Reset()
	G.IndirObj = &wrench
	gD().GateFlag = false
	boltFcn(ActUnk)
	if !strings.Contains(out.String(), "won't turn") {
		t.Fatalf("expected bolt immovable message")
	}

	out.Reset()
	G.ActVerb = ActionVerb{Norm: "take"}
	boltFcn(ActUnk)
	if !strings.Contains(out.String(), "integral part") {
		t.Fatalf("expected integral part message")
	}

	out.Reset()
	G.ActVerb = ActionVerb{Norm: "oil"}
	boltFcn(ActUnk)
	if !strings.Contains(out.String(), "contained glue") {
		t.Fatalf("expected oil warning message")
	}

	// ---- damFunction ----
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "open"}
	if !damFunction(ActUnk) {
		t.Fatalf("expected dam open handled")
	}

	out.Reset()
	G.ActVerb = ActionVerb{Norm: "plug"}
	G.IndirObj = &hands
	if !damFunction(ActUnk) {
		t.Fatalf("expected dam plug handled")
	}

	// ---- puncturedBoatFcn ----
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "inflate"}
	puncturedBoat.MoveTo(G.Here)
	if !puncturedBoatFcn(ActUnk) {
		t.Fatalf("expected punctured boat inflate handled")
	}
	if !strings.Contains(out.String(), "No chance") {
		t.Fatalf("expected punctured boat inflate message")
	}

	out.Reset()
	G.ActVerb = ActionVerb{Norm: "plug"}
	G.IndirObj = &putty
	puncturedBoat.MoveTo(G.Here)
	if !puncturedBoatFcn(ActUnk) {
		t.Fatalf("expected punctured boat plug handled")
	}
	if puncturedBoat.In != nil {
		t.Fatalf("expected punctured boat repaired")
	}

	// ---- inflatableBoatFcn ----
	out.Reset()
	inflatableBoat.MoveTo(G.Here)
	G.ActVerb = ActionVerb{Norm: "inflate"}
	G.IndirObj = &pump
	if !inflatableBoatFcn(ActUnk) {
		t.Fatalf("expected inflatable boat inflate handled")
	}
	if !inflatedBoat.IsIn(G.Here) {
		t.Fatalf("expected inflated boat in room")
	}

	// ---- riverFcn ----
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "put"}
	G.DirObj = &advertisement
	G.IndirObj = &river
	advertisement.MoveTo(G.Winner)
	if !riverFcn(ActUnk) {
		t.Fatalf("expected river put handled")
	}
	if advertisement.In != nil {
		t.Fatalf("expected leaflet removed by river")
	}

	out.Reset()
	G.ActVerb = ActionVerb{Norm: "leap"}
	if !riverFcn(ActUnk) {
		t.Fatalf("expected river leap handled")
	}

	// ---- damRoomFcn ----
	out.Reset()
	if !damRoomFcn(ActLook) {
		t.Fatalf("expected dam room look handled")
	}

	// ---- bubbleFcn ----
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "take"}
	bubbleFcn(ActUnk)

	// ---- whiteCliffsFcn ----
	out.Reset()
	G.Here = &whiteCliffsNorth
	inflatedBoat.MoveTo(G.Winner)
	whiteCliffsFcn(ActEnd)

	// ---- fallsRoomFcn ----
	out.Reset()
	if !fallsRoomFcn(ActLook) {
		t.Fatalf("expected falls room look handled")
	}

	// ---- rivr4RoomFcn ----
	out.Reset()
	buoy.MoveTo(G.Winner)
	gD().BuoyFlag = true
	if rivr4RoomFcn(ActEnd) {
		// no-op
	}
	if gD().BuoyFlag {
		t.Fatalf("expected buoy flag cleared")
	}

	// ---- rBoatFcn deflate ----
	out.Reset()
	inflatedBoat.MoveTo(G.Here)
	G.ActVerb = ActionVerb{Norm: "deflate"}
	G.Winner.MoveTo(&westOfHouse)
	if !rBoatFcn(ActUnk) {
		t.Fatalf("expected boat deflate handled")
	}
	if !inflatableBoat.IsIn(G.Here) {
		t.Fatalf("expected inflatable boat after deflate")
	}

	// ---- iRfill + iMaintRoom ----
	out.Reset()
	G.Here = &reservoir
	gD().LowTide = true
	Queue("iRfill", 1).Run = true
	iRfill()

	out.Reset()
	G.Here = &maintenanceRoom
	gD().WaterLevel = 0
	iMaintRoom()
}

func TestActionCoverageBasics(t *testing.T) {
	out := setupTestGame(t, "")

	// leavesAppear
	G.ActVerb = ActionVerb{Norm: "move"}
	gD().GrateRevealed = false
	grate.Give(FlgInvis)
	leavesAppear()
	if !gD().GrateRevealed || grate.Has(FlgInvis) {
		t.Fatalf("expected leaves to reveal the grate")
	}

	// fweep + flyMe
	out.Reset()
	G.Rand = newSeededRNG(1)
	G.Here = &batRoom
	G.Winner.MoveTo(G.Here)
	fweep(2)
	if !strings.Contains(out.String(), "fweep") {
		t.Fatalf("expected fweep output")
	}
	out.Reset()
	flyMe()
	if G.Here == &batRoom {
		t.Fatalf("expected flyMe to move player")
	}

	// withTell
	out.Reset()
	withTell(&sword)
	if !strings.Contains(out.String(), "With a") {
		t.Fatalf("expected withTell output")
	}

	// badEgg
	egg.MoveTo(&livingRoom)
	canary.MoveTo(&egg)
	out.Reset()
	badEgg()
	if egg.In != nil {
		t.Fatalf("expected egg removed")
	}
	if brokenEgg.In == nil {
		t.Fatalf("expected broken egg placed")
	}

	// slider
	out.Reset()
	advertisement.MoveTo(G.Here)
	slider(&advertisement)
	if advertisement.In != &cellar {
		t.Fatalf("expected item to slide to cellar")
	}

	// randomizeObjects + killInterrupts
	lamp.MoveTo(G.Winner)
	coffin.MoveTo(G.Winner)
	randomizeObjects()
	killInterrupts()
	if match.Has(FlgOn) {
		t.Fatalf("expected match to be extinguished")
	}
}

func TestActionsMoreCoverage(t *testing.T) {
	out := setupTestGame(t, "")

	// boardFcn
	G.ActVerb = ActionVerb{Norm: "take"}
	boardFcn(ActUnk)

	// teethFcn (no tool)
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "brush"}
	G.DirObj = &teeth
	G.IndirObj = nil
	teethFcn(ActUnk)
	if !strings.Contains(out.String(), "Dental hygiene") {
		t.Fatalf("expected teeth hygiene message")
	}

	// graniteWallFcn across locations
	out.Reset()
	G.Here = &northTemple
	G.ActVerb = ActionVerb{Norm: "find"}
	graniteWallFcn(ActUnk)
	G.ActVerb = ActionVerb{Norm: "take"}
	graniteWallFcn(ActUnk)
	G.Here = &slideRoom
	G.ActVerb = ActionVerb{Norm: "read"}
	graniteWallFcn(ActUnk)
	G.Here = &westOfHouse
	G.ActVerb = ActionVerb{Norm: "find"}
	graniteWallFcn(ActUnk)

	// songbirdFcn
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "listen"}
	songbirdFcn(ActUnk)

	// mountainRangeFcn
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "climb"}
	mountainRangeFcn(ActUnk)

	// forestFcn
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "walk around"}
	G.Here = &westOfHouse
	forestFcn(ActUnk)
	G.Here = &forest1
	forestFcn(ActUnk)
	G.ActVerb = ActionVerb{Norm: "listen"}
	forestFcn(ActUnk)

	// chimneyFcn
	out.Reset()
	G.Here = &kitchen
	G.ActVerb = ActionVerb{Norm: "examine"}
	chimneyFcn(ActUnk)

	// ghostsFcn
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "tell"}
	ghostsFcn(ActUnk)
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "attack"}
	G.DirObj = &ghosts
	ghostsFcn(ActUnk)

	// basketFcn raise/lower/take
	out.Reset()
	G.Here = &shaftRoom
	gD().CageTop = false
	G.ActVerb = ActionVerb{Norm: "raise"}
	basketFcn(ActUnk)
	G.ActVerb = ActionVerb{Norm: "lower"}
	basketFcn(ActUnk)
	G.ActVerb = ActionVerb{Norm: "take"}
	G.DirObj = &raisedBasket
	basketFcn(ActUnk)

	// batFcn (no fly if garlic present)
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "take"}
	garlic.MoveTo(G.Winner)
	batFcn(ActUnk)

	// bellFcn + hotBellFcn
	out.Reset()
	G.Here = &northTemple
	G.ActVerb = ActionVerb{Norm: "ring"}
	bellFcn(ActUnk)
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "take"}
	hotBellFcn(ActUnk)
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "rub"}
	G.IndirObj = &hands
	hotBellFcn(ActUnk)
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "ring"}
	G.IndirObj = nil
	hotBellFcn(ActUnk)
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "pour on"}
	G.DirObj = &water
	G.IndirObj = &hotBell
	water.MoveTo(G.Here)
	hotBellFcn(ActUnk)

	// axeFcn (troll already dead)
	gD().TrollFlag = true
	axeFcn(ActUnk)

	// trapDoorFcn
	out.Reset()
	G.Here = &livingRoom
	G.DirObj = &trapDoor
	G.ActVerb = ActionVerb{Norm: "open"}
	trapDoorFcn(ActUnk)
	G.ActVerb = ActionVerb{Norm: "look under"}
	trapDoorFcn(ActUnk)
	G.Here = &cellar
	G.ActVerb = ActionVerb{Norm: "open"}
	trapDoorFcn(ActUnk)

	// frontDoorFcn + barrowDoorFcn + barrowFcn
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "open"}
	frontDoorFcn(ActUnk)
	out.Reset()
	barrowDoorFcn(ActUnk)
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "through"}
	G.Here = &stoneBarrow
	barrowFcn(ActUnk)

	// bottleFcn
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "throw"}
	G.DirObj = &bottle
	bottle.MoveTo(G.Here)
	water.MoveTo(&bottle)
	bottle.Give(FlgOpen)
	bottleFcn(ActUnk)

	// crackFcn
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "through"}
	crackFcn(ActUnk)

	// grateFcn unlock/open
	out.Reset()
	G.Here = &clearing
	G.DirObj = &grate
	G.IndirObj = &keys
	G.ActVerb = ActionVerb{Norm: "unlock"}
	grateFcn(ActUnk)
	G.ActVerb = ActionVerb{Norm: "open"}
	grateFcn(ActUnk)

	// knifeFcn
	G.ActVerb = ActionVerb{Norm: "take"}
	knifeFcn(ActUnk)

	// skeletonFcn
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "take"}
	skeletonFcn(ActUnk)

	// torchFcn
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "examine"}
	torchFcn(ActUnk)
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "lamp off"}
	G.DirObj = &torch
	torch.Give(FlgOn)
	torchFcn(ActUnk)

	// rustyKnifeFcn
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "take"}
	sword.MoveTo(G.Winner)
	rustyKnifeFcn(ActUnk)

	// leafPileFcn
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "count"}
	leafPileFcn(ActUnk)
	out.Reset()
	gD().GrateRevealed = false
	G.ActVerb = ActionVerb{Norm: "move"}
	leafPileFcn(ActUnk)

	// matchFcn
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "lamp on"}
	G.DirObj = &match
	gD().MatchCount = 2
	matchFcn(ActUnk)
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "lamp off"}
	matchFcn(ActUnk)
}

func TestActionsSweep(t *testing.T) {
	out := setupTestGame(t, "")

	// gunkFcn, bodyFcn, blackBookFcn
	out.Reset()
	gunk.MoveTo(G.Here)
	G.ActVerb = ActionVerb{Norm: "take"}
	gunkFcn(ActUnk)
	G.ActVerb = ActionVerb{Norm: "take"}
	bodyFcn(ActUnk)
	G.ActVerb = ActionVerb{Norm: "open"}
	blackBookFcn(ActUnk)

	// sceptreFcn (rainbow solidify)
	out.Reset()
	G.Here = &aragainFalls
	G.ActVerb = ActionVerb{Norm: "wave"}
	gD().RainbowFlag = false
	sceptreFcn(ActUnk)

	// slideFcn
	out.Reset()
	G.Here = &cellar
	G.ActVerb = ActionVerb{Norm: "through"}
	slideFcn(ActUnk)

	// sandwichBagFcn
	out.Reset()
	lunch.MoveTo(&sandwichBag)
	G.ActVerb = ActionVerb{Norm: "smell"}
	sandwichBagFcn(ActUnk)

	// toolChestFcn
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "examine"}
	toolChestFcn(ActUnk)
	G.ActVerb = ActionVerb{Norm: "take"}
	toolChestFcn(ActUnk)

	// buttonFcn (blue + yellow)
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "read"}
	buttonFcn(ActUnk)
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "push"}
	G.DirObj = &blueButton
	gD().WaterLevel = 0
	buttonFcn(ActUnk)
	out.Reset()
	G.DirObj = &yellowButton
	buttonFcn(ActUnk)

	// leakFcn + putty/tube
	out.Reset()
	gD().WaterLevel = 1
	G.ActVerb = ActionVerb{Norm: "put"}
	G.DirObj = &putty
	G.IndirObj = &leak
	leakFcn(ActUnk)
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "oil"}
	G.IndirObj = &putty
	puttyFcn(ActUnk)
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "squeeze"}
	tube.Give(FlgOpen)
	putty.MoveTo(&tube)
	G.DirObj = &tube
	tubeFcn(ActUnk)

	// machine + switch
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "open"}
	machine.Give(FlgOpen)
	machine.Take(FlgOpen)
	coal.MoveTo(&machine)
	machineFcn(ActUnk)
	machine.Take(FlgOpen)
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "turn"}
	G.IndirObj = &screwdriver
	machineSwitchFcn(ActUnk)
	if !diamond.IsIn(&machine) && !gunk.IsIn(&machine) {
		t.Fatalf("expected machine output to produce diamond or gunk")
	}

	// swordFcn examine
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "examine"}
	sword.SetTValue(1)
	G.DirObj = &sword
	swordFcn(ActUnk)

	// boardedWindowFcn, nailsPseudo, cliffObjectFcn, whiteCliffFcn
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "open"}
	boardedWindowFcn(ActUnk)
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "take"}
	nailsPseudo(ActUnk)
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "put"}
	G.DirObj = &advertisement
	G.IndirObj = &climbableCliff
	advertisement.MoveTo(G.Winner)
	cliffObjectFcn(ActUnk)
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "climb"}
	whiteCliffFcn(ActUnk)

	// rainbowFcn
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "look under"}
	rainbowFcn(ActUnk)
	out.Reset()
	G.Here = &aragainFalls
	gD().RainbowFlag = true
	G.ActVerb = ActionVerb{Norm: "cross"}
	rainbowFcn(ActUnk)

	// ropeFcn
	out.Reset()
	G.Here = &domeRoom
	G.ActVerb = ActionVerb{Norm: "tie"}
	G.IndirObj = &railing
	ropeFcn(ActUnk)
	G.ActVerb = ActionVerb{Norm: "untie"}
	ropeFcn(ActUnk)

	// eggObjectFcn
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "open"}
	G.DirObj = &egg
	G.IndirObj = &hands
	eggObjectFcn(ActUnk)

	// canaryObjectFcn
	out.Reset()
	G.Here = &forest1
	G.ActVerb = ActionVerb{Norm: "wind"}
	gD().SingSong = false
	G.DirObj = &canary
	canaryObjectFcn(ActUnk)

	// rugFcn
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "move"}
	gD().RugMoved = false
	rugFcn(ActUnk)
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "look under"}
	rugFcn(ActUnk)

	// sandFunction
	out.Reset()
	G.Here = &sandyCave
	G.ActVerb = ActionVerb{Norm: "dig"}
	G.IndirObj = &shovel
	gD().BeachDig = 0
	sandFunction(ActUnk)
	sandFunction(ActUnk)

	// batDescFcn
	out.Reset()
	garlic.MoveTo(G.Winner)
	batDescFcn(ActUnk)
	out.Reset()
	removeCarefully(&garlic)
	batDescFcn(ActUnk)

	// trophyCaseFcn, trunk/bagOfCoins, garlicFcn
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "take"}
	trophyCaseFcn(ActUnk)
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "open"}
	trunkFcn(ActUnk)
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "open"}
	bagOfCoinsFcn(ActUnk)
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "eat"}
	G.DirObj = &garlic
	garlic.MoveTo(G.Winner)
	garlicFcn(ActUnk)
}

func TestActionsInterruptsAndPseudo(t *testing.T) {
	setupTestGame(t, "quit\n")

	// stoneBarrowFcn (finish reads quit)
	G.ActVerb = ActionVerb{Norm: "enter"}
	stoneBarrowFcn(ActBegin)

	// deadFunction branches
	G.ActVerb = ActionVerb{Norm: "take"}
	deadFunction(ActUnk)
	G.ActVerb = ActionVerb{Norm: "wait"}
	deadFunction(ActUnk)
	G.ActVerb = ActionVerb{Norm: "look"}
	deadFunction(ActUnk)
	G.ActVerb = ActionVerb{Norm: "pray"}
	G.Here = &westOfHouse
	deadFunction(ActUnk)

	// iXb/iXc/iCyclops
	gD().XC = false
	G.Here = &enteranceToHades
	iXb()
	iXc()
	gD().CyclopsFlag = true
	iCyclops()

	// gratingExitFcn
	gD().GrateRevealed = false
	gratingExitFcn()
	gD().GrateRevealed = true
	grate.Take(FlgOpen)
	gratingExitFcn()
	grate.Give(FlgOpen)
	gratingExitFcn()

	// pseudo functions
	G.ActVerb = ActionVerb{Norm: "leap"}
	chasmPseudo(ActUnk)
	G.ActVerb = ActionVerb{Norm: "cross"}
	lakePseudo(ActUnk)
	G.ActVerb = ActionVerb{Norm: "swim"}
	streamPseudo(ActUnk)
	G.ActVerb = ActionVerb{Norm: "kiss"}
	domePseudo(ActUnk)
	G.ActVerb = ActionVerb{Norm: "through"}
	gatePseudo(ActUnk)
	G.ActVerb = ActionVerb{Norm: "open"}
	doorPseudo(ActUnk)
	G.ActVerb = ActionVerb{Norm: "mung"}
	paintPseudo(ActUnk)
	G.ActVerb = ActionVerb{Norm: "smell"}
	gasPseudo(ActUnk)
	G.ActVerb = ActionVerb{Norm: "examine"}
	chainPseudo(ActUnk)

	// trapDoorExitFcn + upChimneyFcn
	gD().RugMoved = false
	trapDoorExitFcn()
	gD().RugMoved = true
	trapDoor.Take(FlgOpen)
	trapDoorExitFcn()

	G.Winner.Children = nil
	upChimneyFcn()
	lamp.MoveTo(G.Winner)
	upChimneyFcn()

	// mazeDiodesFcn
	G.Here = &maze2
	mazeDiodesFcn()

	// iForestRandom (ensure no panic)
	G.Here = &forest1
	G.Rand = newSeededRNG(1)
	iForestRandom()

	// iMatch
	match.Give(FlgOn)
	match.Give(FlgFlame)
	iMatch()
}

func TestActionsMoreBranches(t *testing.T) {
	setupTestGame(t, "")

	// paintingFcn
	G.ActVerb = ActionVerb{Norm: "mung"}
	G.DirObj = &painting
	paintingFcn(ActUnk)

	// candlesFcn (light with match, then lamp off)
	match.Give(FlgOn)
	G.ActVerb = ActionVerb{Norm: "lamp on"}
	G.IndirObj = &match
	candlesFcn(ActUnk)
	G.ActVerb = ActionVerb{Norm: "lamp off"}
	candles.Give(FlgOn)
	candlesFcn(ActUnk)

	// eggObjectFcn (open with tool)
	G.ActVerb = ActionVerb{Norm: "open"}
	G.DirObj = &egg
	G.IndirObj = &knife
	eggObjectFcn(ActUnk)

	// trollFcn branches
	G.Here = &trollRoom
	troll.MoveTo(G.Here)
	G.ActVerb = ActionVerb{Norm: "tell"}
	trollFcn(ActUnk)
	trollFcn(ActBusy)
	trollFcn(ActUnconscious)
	trollFcn(ActConscious)
	G.ActVerb = ActionVerb{Norm: "examine"}
	trollFcn(ActUnk)

	// cyclopsFcn branches
	G.Here = &cyclopsRoom
	cyclops.MoveTo(G.Here)
	gD().CyclopsFlag = false
	G.Winner = &cyclops
	G.ActVerb = ActionVerb{Norm: "odysseus"}
	cyclopsFcn(ActUnk)
	G.Winner = &adventurer
	gD().CyclopsFlag = true
	G.ActVerb = ActionVerb{Norm: "examine"}
	cyclopsFcn(ActUnk)
}

func TestCombatHelpersCoverage(t *testing.T) {
	out := setupTestGame(t, "")
	_ = out

	// weaponFunction
	G.ActVerb = ActionVerb{Norm: "take"}
	G.Here = &trollRoom
	troll.MoveTo(G.Here)
	axe.MoveTo(&troll)
	if !weaponFunction(&axe, &troll) {
		t.Fatalf("expected weaponFunction to handle take")
	}

	// fightStrength
	G.Score = 10
	G.Winner.SetStrength(1)
	if fightStrength(false) <= 0 {
		t.Fatalf("expected fightStrength positive")
	}
	_ = fightStrength(true)

	// findWeapon
	sword.MoveTo(G.Winner)
	if findWeapon(G.Winner) == nil {
		t.Fatalf("expected to find weapon")
	}

	// winning (deterministic)
	G.Rand = newSeededRNG(1)
	troll.SetStrength(2)
	_ = winning(&troll)

	// awaken
	awakened := false
	thief.Action = func(ActionArg) bool {
		awakened = true
		return true
	}
	thief.SetStrength(-2)
	awaken(&thief)
	if !awakened {
		t.Fatalf("expected awaken to call action")
	}

	// randomMeleeMsg + remark
	msg := randomMeleeMsg(MeleeSet{
		{mp("hit "), md(), mp(" with "), mw()},
	})
	if msg == nil {
		t.Fatalf("expected melee msg")
	}
	remark(msg, &troll, &sword)

	// villainStrength
	entry := &VillainEntry{Villain: &thief, Best: &sword, BestAdv: 1, Prob: 10, Msgs: &trollMelee}
	thief.SetStrength(3)
	G.IndirObj = &sword
	_ = villainStrength(entry)

	// villainResult (non-kill)
	thief.SetStrength(2)
	villainResult(&thief, 2, blowLightWnd)

	// winnerResult (no death)
	G.Winner.SetStrength(2)
	winnerResult(2, blowLightWnd, 0)

	// iSword (held and not held)
	sword.MoveTo(G.Winner)
	G.Here = &trollRoom
	iSword()
	removeCarefully(&sword)
	iSword()
}

func TestGlobalsCoverage(t *testing.T) {
	out := setupTestGame(t, "")

	// sailorFcn
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "tell"}
	if !sailorFcn(ActUnk) {
		t.Fatalf("expected sailor tell handled")
	}
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "hello"}
	sailorFcn(ActUnk)
	if !strings.Contains(out.String(), "Nothing happens") && !strings.Contains(out.String(), "repeating") && !strings.Contains(out.String(), "worn out") {
		t.Fatalf("expected sailor hello output")
	}

	// groundFunction (dig + put on ground)
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "dig"}
	if !groundFunction(ActUnk) {
		t.Fatalf("expected ground dig handled")
	}
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "put"}
	G.DirObj = &advertisement
	G.IndirObj = &ground
	groundFunction(ActUnk)

	// grueFunction
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "examine"}
	if !grueFunction(ActUnk) {
		t.Fatalf("expected grue examine handled")
	}

	// cretinFcn
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "eat"}
	if !cretinFcn(ActUnk) {
		t.Fatalf("expected cretin eat handled")
	}
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "attack"}
	G.IndirObj = &sword
	sword.Give(FlgWeapon)
	cretinFcn(ActUnk)

	// pathObject
	out.Reset()
	G.ActVerb = ActionVerb{Norm: "take"}
	if !pathObject(ActUnk) {
		t.Fatalf("expected path take handled")
	}
}

func TestLowCoverageTargets(t *testing.T) {
	setupTestGame(t, "")

	// verbs_meta
	vWin(ActUnk)
	vAdvent(ActUnk)

	// verbs_movement climb down/up
	G.DirObj = &tree
	vClimbDown(ActUnk)
	vClimbUp(ActUnk)

	// verbs: put on, ring, rub
	G.DirObj = &advertisement
	G.IndirObj = &ground
	vPutOn(ActUnk)
	vRing(ActUnk)
	vRub(ActUnk)

	// findIn + vSay
	G.Here = &trollRoom
	troll.MoveTo(G.Here)
	troll.Give(FlgPerson)
	_ = findIn(G.Here, FlgPerson)
	G.Params.Continue = NumUndef
	vSay(ActUnk)
	G.Params.Continue = 0
	G.LexRes = []LexItem{{Norm: "hello", Orig: "hello"}}
	vSay(ActUnk)

	// vLookInside
	G.DirObj = &kitchenWindow
	vLookInside(ActUnk)

	// iCyclops (safe branch)
	gD().CyclopsFlag = false
	gD().Dead = false
	G.Here = &cyclopsRoom
	gD().CycloWrath = 0
	iCyclops()

	// troll/cyclops functions
	G.ActVerb = ActionVerb{Norm: "attack"}
	G.DirObj = &troll
	trollFcn(ActUnk)
	G.DirObj = &cyclops
	cyclopsFcn(ActUnk)

	// whiteHouseFcn
	G.Here = &kitchen
	G.ActVerb = ActionVerb{Norm: "find"}
	whiteHouseFcn(ActUnk)
	G.Here = &clearing
	G.ActVerb = ActionVerb{Norm: "find"}
	whiteHouseFcn(ActUnk)
}

func TestObjectActionVerbSweep(t *testing.T) {
	setupTestGame(t, "")

	verbs := []string{"take", "open", "close", "move", "look under", "rub", "ring", "count", "smell"}

	for _, obj := range G.AllObjects {
		if obj == nil || obj.Action == nil {
			continue
		}
		if obj == &bones {
			continue
		}
		if obj.IsIn(&rooms) || obj == &rooms {
			continue
		}
		if obj.Location() != nil {
			G.Here = obj.Location()
		} else {
			G.Here = &westOfHouse
		}
		G.DirObj = obj
		G.IndirObj = nil
		for _, v := range verbs {
			if obj.Location() != nil {
				G.Here = obj.Location()
			} else {
				G.Here = &westOfHouse
			}
			G.ActVerb = ActionVerb{Norm: v}
			obj.Action(ActUnk)
		}
	}
}

func TestRoomAndObjectActionSweep(t *testing.T) {
	setupTestGame(t, "")

	// Sweep room actions with ActLook/ActEnter to cover room descriptions.
	for _, rm := range rooms.Children {
		if rm.Action == nil {
			continue
		}
		G.Here = rm
		rm.Action(ActLook)
		rm.Action(ActEnter)
	}

	// Sweep non-room object actions with "examine" to cover simple branches.
	G.ActVerb = ActionVerb{Norm: "examine"}
	for _, obj := range G.AllObjects {
		if obj == nil || obj.Action == nil {
			continue
		}
		if obj.IsIn(&rooms) || obj == &rooms {
			continue
		}
		if obj.Location() != nil {
			G.Here = obj.Location()
		} else {
			G.Here = &westOfHouse
		}
		G.DirObj = obj
		obj.Action(ActUnk)
	}
}

func TestThiefHelpersCoverage(t *testing.T) {
	setupTestGame(t, "")

	// rob
	trunk.MoveTo(G.Winner)
	trunk.SetTValue(10)
	trunk.Take(FlgInvis)
	if !rob(G.Winner, &thief, 0) {
		t.Fatalf("expected rob to move treasure")
	}

	// recoverStiletto
	stiletto.MoveTo(thief.Location())
	recoverStiletto()
	if !stiletto.IsIn(&thief) {
		t.Fatalf("expected stiletto moved to thief")
	}

	// hackTreasures
	treasureRoom.Give(FlgOn)
	chalice.MoveTo(&treasureRoom)
	hackTreasures()

	// depositBooty
	egg.MoveTo(&thief)
	egg.SetTValue(5)
	if !depositBooty(&treasureRoom) {
		t.Fatalf("expected depositBooty to move treasure")
	}

	// robMaze (non-deterministic; just ensure no panic)
	maze5.Give(FlgMaze)
	bagOfCoins.MoveTo(&maze5)
	robMaze(&maze5)

	// moveAll
	advertisement.MoveTo(&maze5)
	moveAll(&maze5, &treasureRoom)

	// thiefInTreasure
	G.Here = &treasureRoom
	treasureRoom.Children = append(treasureRoom.Children, &adventurer)
	thiefInTreasure()

	// infested
	troll.MoveTo(&treasureRoom)
	troll.Give(FlgActor)
	if !infested(&treasureRoom) {
		t.Fatalf("expected room infested")
	}

	// largeBagFcn
	G.ActVerb = ActionVerb{Norm: "take"}
	largeBagFcn(ActUnk)

	// stiletteFcn
	G.ActVerb = ActionVerb{Norm: "take"}
	stiletteFcn(ActUnk)

	// treasureInsideFcn
	treasureInsideFcn(ActUnk)

	// stealJunk (use stiletto to force)
	stiletto.MoveTo(&maze5)
	maze5.Give(FlgTouch)
	G.Here = &maze5
	stealJunk(&maze5)
}

func TestVerbCoverageBasics(t *testing.T) {
	out := setupTestGame(t, "")

	// vAlarm
	out.Reset()
	G.DirObj = &advertisement
	vAlarm(ActUnk)

	// vAnswer
	out.Reset()
	vAnswer(ActUnk)
	if G.Params.InQuotes {
		t.Fatalf("expected answer to clear quotes")
	}

	// vAttack (non-person, bare hands, not holding weapon, wrong weapon)
	out.Reset()
	G.DirObj = &advertisement
	G.IndirObj = &hands
	vAttack(ActUnk)
	out.Reset()
	G.DirObj = &troll
	G.DirObj.Give(FlgPerson)
	G.IndirObj = &hands
	vAttack(ActUnk)
	out.Reset()
	G.IndirObj = &sword
	sword.Take(FlgWeapon)
	vAttack(ActUnk)
	out.Reset()
	sword.Give(FlgWeapon)
	vAttack(ActUnk)

	// vBack, vBlast, vBrush
	out.Reset()
	vBack(ActUnk)
	vBlast(ActUnk)
	vBrush(ActUnk)

	// preBurn + vBurn
	out.Reset()
	G.IndirObj = nil
	preBurn(ActUnk)
	out.Reset()
	G.IndirObj = &match
	match.Give(FlgOn)
	preBurn(ActUnk)
	out.Reset()
	G.DirObj = &advertisement
	vBurn(ActUnk)

	// vChomp, vClose (non-closable)
	out.Reset()
	vChomp(ActUnk)
	G.DirObj = &advertisement
	vClose(ActUnk)

	// vCommand
	out.Reset()
	G.DirObj = &advertisement
	vCommand(ActUnk)

	// vCount
	out.Reset()
	vCount(ActUnk)

	// vDrink (no water)
	out.Reset()
	G.DirObj = &water
	vDrink(ActUnk)

	// vEat (edible vs not)
	out.Reset()
	G.DirObj = &sword
	vEat(ActUnk)
	out.Reset()
	lunch.MoveTo(G.Winner)
	G.DirObj = &lunch
	vEat(ActUnk)
	if !strings.Contains(out.String(), "hit the spot") {
		t.Fatalf("expected eat output")
	}
}

func TestVerbsSweep(t *testing.T) {
	setupTestGame(t, "")

	G.Here = &livingRoom
	advertisement.MoveTo(G.Here)
	G.DirObj = &advertisement
	G.IndirObj = &hands

	vCurses(ActUnk)
	vCut(ActUnk)
	vDeflate(ActUnk)
	vDig(ActUnk)
	vDisenchant(ActUnk)
	vDrinkFrom(ActUnk)
	vEcho(ActUnk)
	vEnchant(ActUnk)
	vExcorcise(ActUnk)
	vFind(ActUnk)
	vFrobozz(ActUnk)
	vHatch(ActUnk)
	vHello(ActUnk)
	vIncant(ActUnk)
	vInflate(ActUnk)
	vKick(ActUnk)
	vKiss(ActUnk)
	vKnock(ActUnk)
	vLeanOn(ActUnk)
	vListen(ActUnk)
	vMelt(ActUnk)
	vMumble(ActUnk)
	vMung(ActUnk)
	vOil(ActUnk)
	vPick(ActUnk)
	vPlay(ActUnk)
	vPlug(ActUnk)
	vPushTo(ActUnk)
	vPutBehind(ActUnk)
	vPutUnder(ActUnk)
	vRaise(ActUnk)
	vRape(ActUnk)
	vReadPage(ActUnk)
	vRepent(ActUnk)
	vReply(ActUnk)
	vSay(ActUnk)
	vSearch(ActUnk)
	vSend(ActUnk)
	vSGive(ActUnk)
	vSkip(ActUnk)
	vSmell(ActUnk)
	vSpin(ActUnk)
	vSpray(ActUnk)
	vSqueeze(ActUnk)
	vStab(ActUnk)
	vStrike(ActUnk)
	vSwing(ActUnk)
	vThrowOff(ActUnk)
	vTreasure(ActUnk)
	vWear(ActUnk)
	vWind(ActUnk)
	vYell(ActUnk)

	// vOpen + vClose on bottle
	bottle.MoveTo(G.Winner)
	G.DirObj = &bottle
	vOpen(ActUnk)
	vClose(ActUnk)

	// vPut into bottle
	bottle.Give(FlgOpen)
	G.DirObj = &advertisement
	G.IndirObj = &bottle
	advertisement.MoveTo(G.Winner)
	vPut(ActUnk)

	// vTake
	advertisement.MoveTo(G.Here)
	G.DirObj = &advertisement
	vTake(ActUnk)

	// vThrow at person
	G.DirObj = &advertisement
	G.IndirObj = &troll
	troll.Give(FlgPerson)
	vThrow(ActUnk)

	// vGive
	G.DirObj = &troll
	G.IndirObj = &adventurer
	vGive(ActUnk)

	// vOverboard
	inflatedBoat.MoveTo(G.Here)
	G.Winner.MoveTo(&inflatedBoat)
	G.IndirObj = &teeth
	vOverboard(ActUnk)
	G.Winner.MoveTo(G.Here)

	// vPourOn water on torch
	water.MoveTo(G.Here)
	G.DirObj = &water
	G.IndirObj = &torch
	torch.Give(FlgOn)
	vPourOn(ActUnk)

	// vPray (both branches)
	G.Here = &westOfHouse
	vPray(ActUnk)
	G.Here = &southTemple
	vPray(ActUnk)

	// vPump
	pump.MoveTo(G.Winner)
	G.DirObj = &inflatableBoat
	G.IndirObj = &pump
	vPump(ActUnk)

	// vPush (generic)
	G.DirObj = &advertisement
	vPush(ActUnk)

	// vLock/vUnlock with keys
	G.Here = &clearing
	G.DirObj = &grate
	G.IndirObj = &keys
	vUnlock(ActUnk)
	vLock(ActUnk)

	// vLook/vExamine/vRead
	G.Here = &livingRoom
	G.DirObj = &advertisement
	vLook(ActUnk)
	vExamine(ActUnk)
	vRead(ActUnk)

	// vOdysseus in cyclops room
	G.Here = &cyclopsRoom
	cyclops.MoveTo(G.Here)
	gD().CyclopsFlag = false
	vOdysseus(ActUnk)

	// vTie/vTieUp/vUntie/vWave
	G.DirObj = &rope
	G.IndirObj = &railing
	vTie(ActUnk)
	G.IndirObj = &rope
	vTieUp(ActUnk)
	vUntie(ActUnk)
	vWave(ActUnk)
}

func TestVerbsMoreCoverage(t *testing.T) {
	setupTestGame(t, "")

	// vBreathe
	G.DirObj = &lungs
	vBreathe(ActUnk)

	// preFill + vFill + hitSpot
	G.Here = &westOfHouse
	G.DirObj = &bottle
	G.IndirObj = nil
	preFill(ActUnk)
	vFill(ActUnk)
	water.MoveTo(G.Winner)
	G.DirObj = &water
	hitSpot()

	// vShake + shakeLoop
	box := &Object{Desc: "box", Flags: FlgTake | FlgCont | FlgOpen, Item: &ItemData{Capacity: 10}}
	box.MoveTo(G.Here)
	advertisement.MoveTo(box)
	G.DirObj = box
	vShake(ActUnk)

	// preSGive
	G.DirObj = &troll
	G.IndirObj = &advertisement
	preSGive(ActUnk)

	// vLaunch (no boat)
	G.DirObj = &inflatedBoat
	vLaunch(ActUnk)

	// vLookBehind/On/Under
	G.DirObj = &advertisement
	vLookBehind(ActUnk)
	vLookOn(ActUnk)
	vLookUnder(ActUnk)

	// vMake, vMove
	vMake(ActUnk)
	G.DirObj = &advertisement
	vMove(ActUnk)

	// preMung + mungRoom
	G.DirObj = &advertisement
	G.IndirObj = nil
	preMung(ActUnk)
	mungRoom(&westOfHouse, "No entry.")

	// vTell
	G.DirObj = &advertisement
	vTell(ActUnk)

	// vTurn
	G.DirObj = &bolt
	G.IndirObj = &wrench
	vTurn(ActUnk)
}

func TestVerbBranches(t *testing.T) {
	setupTestGame(t, "")

	// vOpen / vClose on container
	mailbox.MoveTo(G.Here)
	G.DirObj = &mailbox
	vOpen(ActUnk)
	vClose(ActUnk)

	// vOpen on door
	G.DirObj = &kitchenWindow
	kitchenWindow.Give(FlgDoor)
	vOpen(ActUnk)

	// vPut into open bottle
	bottle.MoveTo(G.Winner)
	bottle.Give(FlgOpen)
	advertisement.MoveTo(G.Winner)
	G.DirObj = &advertisement
	G.IndirObj = &bottle
	vPut(ActUnk)

	// vDrop held item
	G.DirObj = &advertisement
	vDrop(ActUnk)

	// vTake from room
	advertisement.MoveTo(G.Here)
	G.DirObj = &advertisement
	vTake(ActUnk)

	// vLookUnder on room item
	G.DirObj = &rug
	vLookUnder(ActUnk)

	// vTreasure toggle
	G.Here = &northTemple
	vTreasure(ActUnk)
	G.Here = &treasureRoom
	vTreasure(ActUnk)
}

func TestVerbsMovementCoverage(t *testing.T) {
	setupTestGame(t, "")

	// vEnter/vExit
	G.Here = &westOfHouse
	G.DirObj = &rooms
	vEnter(ActUnk)
	vExit(ActUnk)

	// vFollow
	vFollow(ActUnk)

	// vLeap (no down exit)
	G.DirObj = nil
	vLeap(ActUnk)

	// vLeave
	vLeave(ActUnk)

	// vStand/vStay
	vStand(ActUnk)
	vStay(ActUnk)

	// vSwim
	G.Here = &westOfHouse
	G.DirObj = &water
	vSwim(ActUnk)

	// vClimbOn non-vehicle
	G.DirObj = &advertisement
	vClimbOn(ActUnk)

	// vClimbFcn with walls/tree fallbacks
	G.Here = &westOfHouse
	G.DirObj = &tree
	vClimbFcn(Up, &tree)
	G.DirObj = &graniteWall
	vClimbFcn(Up, &graniteWall)
}
