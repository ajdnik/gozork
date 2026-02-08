package game

import (
	"errors"
	. "github.com/ajdnik/gozork/engine"
	"strings"
	"testing"
)

// ---- verbs_meta_test.go ----

func TestVerbosityModes(t *testing.T) {
	out := setupTestGame(t, "")

	G.Verbose = false
	G.SuperBrief = true
	vVerbose(ActUnk)
	if !G.Verbose || G.SuperBrief {
		t.Fatalf("expected verbose on and superbrief off")
	}
	if !strings.Contains(out.String(), "Maximum verbosity") {
		t.Fatalf("expected verbose output, got %q", out.String())
	}

	out.Reset()
	vBrief(ActUnk)
	if G.Verbose || G.SuperBrief {
		t.Fatalf("expected brief to disable verbose/superbrief")
	}
	if !strings.Contains(out.String(), "Brief descriptions") {
		t.Fatalf("expected brief output, got %q", out.String())
	}

	out.Reset()
	vSuperBrief(ActUnk)
	if !G.SuperBrief {
		t.Fatalf("expected superbrief on")
	}
	if !strings.Contains(out.String(), "Superbrief descriptions") {
		t.Fatalf("expected superbrief output, got %q", out.String())
	}
}

func TestScriptToggles(t *testing.T) {
	out := setupTestGame(t, "")

	G.Script = false
	vScript(ActUnk)
	if !G.Script {
		t.Fatalf("expected script to be enabled")
	}
	if !strings.Contains(out.String(), "transcript") {
		t.Fatalf("expected transcript start message, got %q", out.String())
	}

	out.Reset()
	vUnscript(ActUnk)
	if G.Script {
		t.Fatalf("expected script to be disabled")
	}
	if !strings.Contains(out.String(), "transcript") {
		t.Fatalf("expected transcript end message, got %q", out.String())
	}
}

func TestVerifyAndWait(t *testing.T) {
	out := setupTestGame(t, "")

	vVerify(ActUnk)
	if !strings.Contains(out.String(), "Verifying disk") || !strings.Contains(out.String(), "disk is correct") {
		t.Fatalf("expected verify output, got %q", out.String())
	}

	G.ClockWait = false
	vWait(ActUnk)
	if !G.ClockWait {
		t.Fatalf("expected wait to set ClockWait")
	}
}

func TestSaveRestore(t *testing.T) {
	out := setupTestGame(t, "")

	saveCalled := false
	restoreCalled := false
	G.Save = func() error {
		saveCalled = true
		return nil
	}
	G.Restore = func() error {
		restoreCalled = true
		return nil
	}

	out.Reset()
	vSave(ActUnk)
	if !saveCalled || !strings.Contains(out.String(), "Ok.") {
		t.Fatalf("expected save to succeed, got %q", out.String())
	}

	out.Reset()
	vRestore(ActUnk)
	if !restoreCalled || !strings.Contains(out.String(), "Ok.") {
		t.Fatalf("expected restore to succeed, got %q", out.String())
	}
}

func TestSaveFailure(t *testing.T) {
	out := setupTestGame(t, "")

	G.Save = func() error {
		return errors.New("nope")
	}
	vSave(ActUnk)
	if !strings.Contains(out.String(), "Failed:") {
		t.Fatalf("expected save failure message, got %q", out.String())
	}
}

func TestQuitYesNo(t *testing.T) {
	out := setupTestGame(t, "n\ny\n")

	G.QuitRequested = false
	vQuit(ActUnk)
	if G.QuitRequested {
		t.Fatalf("expected quit to be rejected")
	}
	if !strings.Contains(out.String(), "Ok.") {
		t.Fatalf("expected ok response, got %q", out.String())
	}

	out.Reset()
	vQuit(ActUnk)
	if !G.QuitRequested {
		t.Fatalf("expected quit to be accepted")
	}
}

func TestRestartNo(t *testing.T) {
	out := setupTestGame(t, "n\n")

	restartCalled := false
	G.Restart = func() error {
		restartCalled = true
		return nil
	}

	if vRestart(ActUnk) {
		t.Fatalf("expected restart to return false when declined")
	}
	if restartCalled {
		t.Fatalf("did not expect restart to be called")
	}
	if !strings.Contains(out.String(), "Do you wish to restart") {
		t.Fatalf("expected restart prompt, got %q", out.String())
	}
}

func TestIsYesAndFinishQuit(t *testing.T) {
	out := setupTestGame(t, "y\nquit\n")

	if !isYes() {
		t.Fatalf("expected isYes to accept y")
	}
	if out.Len() == 0 {
		t.Fatalf("expected prompt output")
	}

	out.Reset()
	G.QuitRequested = false
	if !finish() {
		t.Fatalf("expected finish to return true")
	}
	if !G.QuitRequested {
		t.Fatalf("expected finish to request quit")
	}
}

// ---- verbs_movement_test.go ----

func TestPreBoardNotInRoom(t *testing.T) {
	out := setupTestGame(t, "")

	boat := &Object{Desc: "boat", Flags: FlgVeh}
	boat.In = &northOfHouse
	G.DirObj = boat
	G.PerformFatal = false

	preBoard(ActUnk)
	if !G.PerformFatal {
		t.Fatalf("expected preBoard to signal fatal when vehicle not in room")
	}
	if !strings.Contains(out.String(), "must be on the ground") {
		t.Fatalf("expected not-in-room message, got %q", out.String())
	}
}

func TestPreBoardAlreadyInVehicle(t *testing.T) {
	out := setupTestGame(t, "")

	boat := &Object{Desc: "boat", Flags: FlgVeh}
	boat.In = G.Here
	G.DirObj = boat
	G.Winner.MoveTo(boat)
	G.PerformFatal = false

	preBoard(ActUnk)
	if !G.PerformFatal {
		t.Fatalf("expected preBoard to signal fatal when already in vehicle")
	}
	if !strings.Contains(out.String(), "already in the boat") {
		t.Fatalf("expected already-in message, got %q", out.String())
	}
}

func TestVBoardMovesWinner(t *testing.T) {
	out := setupTestGame(t, "")

	boat := &Object{Desc: "boat", Flags: FlgVeh, In: G.Here}
	G.DirObj = boat
	called := false
	boat.Action = func(ActionArg) bool {
		called = true
		return true
	}

	vBoard(ActUnk)
	if !G.Winner.IsIn(boat) {
		t.Fatalf("expected winner to move into boat")
	}
	if !called {
		t.Fatalf("expected vehicle enter action to be called")
	}
	if !strings.Contains(out.String(), "now in the boat") {
		t.Fatalf("expected board message, got %q", out.String())
	}
}

func TestVClimbOnNonVehicle(t *testing.T) {
	out := setupTestGame(t, "")

	rock := &Object{Desc: "rock", In: G.Here}
	G.DirObj = rock

	vClimbOn(ActUnk)
	if !strings.Contains(out.String(), "can't climb onto") {
		t.Fatalf("expected climb-on rejection, got %q", out.String())
	}
}

func TestVClimbFcnNoExit(t *testing.T) {
	out := setupTestGame(t, "")

	G.DirObj = nil
	vClimbFcn(Up, nil)
	if !strings.Contains(out.String(), "can't go that way") {
		t.Fatalf("expected no-exit message, got %q", out.String())
	}
}

func TestVDisembarkNotInThat(t *testing.T) {
	out := setupTestGame(t, "")

	boat := &Object{Desc: "boat", Flags: FlgVeh}
	G.DirObj = boat
	G.PerformFatal = false

	vDisembark(ActUnk)
	if !G.PerformFatal {
		t.Fatalf("expected disembark to be fatal when not in vehicle")
	}
	if !strings.Contains(out.String(), "You're not in that") {
		t.Fatalf("expected not-in message, got %q", out.String())
	}
}

func TestVSwimNoWater(t *testing.T) {
	out := setupTestGame(t, "")

	G.DirObj = nil
	vSwim(ActUnk)
	if !strings.Contains(out.String(), "Go jump in a lake") {
		t.Fatalf("expected no-water swim message, got %q", out.String())
	}
}

func TestThroughHeadBump(t *testing.T) {
	out := setupTestGame(t, "")

	box := &Object{Desc: "box", In: G.Here}
	G.DirObj = box
	vThrough(ActUnk)
	if !strings.Contains(out.String(), "hit your head") {
		t.Fatalf("expected head-bump message, got %q", out.String())
	}
}

func TestOtherSideFindsDoorExit(t *testing.T) {
	setupTestGame(t, "")

	door := &Object{Desc: "door", Flags: FlgDoor}
	room := &Object{Desc: "room", Flags: FlgLand}
	G.Here = room
	room.SetExit(North, ExitProps{DExit: door, RExit: room})

	dir, ok := otherSide(door)
	if !ok || dir != North {
		t.Fatalf("expected otherSide to find north door")
	}
}

func TestVWalkAroundAndWalkTo(t *testing.T) {
	out := setupTestGame(t, "")

	vWalkAround(ActUnk)
	if !strings.Contains(out.String(), "Use compass directions") {
		t.Fatalf("expected walk-around message, got %q", out.String())
	}

	out.Reset()
	G.DirObj = &mailbox
	vWalkTo(ActUnk)
	if !strings.Contains(out.String(), "it's here") {
		t.Fatalf("expected walk-to here message, got %q", out.String())
	}

	out.Reset()
	G.DirObj = nil
	vWalkTo(ActUnk)
	if !strings.Contains(out.String(), "supply a direction") {
		t.Fatalf("expected walk-to direction message, got %q", out.String())
	}
}

func TestMoveToRoomVehicleChecks(t *testing.T) {
	out := setupTestGame(t, "")

	start := &Object{Desc: "start", Flags: FlgLand}
	target := &Object{Desc: "lake"}
	start.In = &rooms
	target.In = &rooms

	G.Here = start
	G.Winner.MoveTo(start)

	out.Reset()
	if moveToRoom(target, false) {
		t.Fatalf("expected moveToRoom to fail without vehicle")
	}
	if !strings.Contains(out.String(), "without a vehicle") {
		t.Fatalf("expected no-vehicle message, got %q", out.String())
	}

	out.Reset()
	boat := &Object{Desc: "boat", Flags: FlgVeh}
	boat.SetVehType(FlgWater)
	boat.MoveTo(start)
	G.Winner.MoveTo(boat)

	if moveToRoom(target, false) {
		t.Fatalf("expected moveToRoom to fail for wrong vehicle type")
	}
	if !strings.Contains(out.String(), "in a boat") {
		t.Fatalf("expected wrong-vehicle message, got %q", out.String())
	}
}

func TestVWalkUExitAndNExit(t *testing.T) {
	out := setupTestGame(t, "")

	roomA := &Object{Desc: "roomA", Flags: FlgLand}
	roomB := &Object{Desc: "roomB", Flags: FlgLand, LongDesc: "B room."}
	roomA.In = &rooms
	roomB.In = &rooms
	roomA.SetExit(North, ExitProps{UExit: true, RExit: roomB})
	roomA.SetExit(South, ExitProps{NExit: "Nope."})

	G.Here = roomA
	G.Winner.MoveTo(roomA)
	G.Params.HasWalkDir = true
	G.Params.WalkDir = North
	G.Lit = true

	if !vWalk(ActUnk) {
		t.Fatalf("expected vWalk to succeed for UExit")
	}
	if G.Here != roomB {
		t.Fatalf("expected to move to roomB")
	}

	out.Reset()
	G.Here = roomA
	G.Winner.MoveTo(roomA)
	G.Params.WalkDir = South
	G.PerformFatal = false
	vWalk(ActUnk)
	if !strings.Contains(out.String(), "Nope.") {
		t.Fatalf("expected NExit message, got %q", out.String())
	}
	if !G.PerformFatal {
		t.Fatalf("expected NExit to be fatal")
	}
}

func TestVWalkNoExit(t *testing.T) {
	out := setupTestGame(t, "")

	room := &Object{Desc: "room", Flags: FlgLand}
	room.In = &rooms
	G.Here = room
	G.Winner.MoveTo(room)
	G.Lit = true

	G.Params.HasWalkDir = true
	G.Params.WalkDir = East
	G.PerformFatal = false

	vWalk(ActUnk)
	if !G.PerformFatal {
		t.Fatalf("expected no-exit walk to be fatal")
	}
	if !strings.Contains(out.String(), "can't go that way") {
		t.Fatalf("expected no-exit message, got %q", out.String())
	}
}

func TestVWalkFExitAndCExitAndDExit(t *testing.T) {
	out := setupTestGame(t, "")

	roomA := &Object{Desc: "roomA", Flags: FlgLand}
	roomB := &Object{Desc: "roomB", Flags: FlgLand}
	roomC := &Object{Desc: "roomC", Flags: FlgLand}
	roomA.In = &rooms
	roomB.In = &rooms
	roomC.In = &rooms

	roomA.SetExit(North, ExitProps{FExit: func() *Object { return roomB }})
	roomA.SetExit(East, ExitProps{CExit: func() bool { return false }, CExitStr: "Blocked.", RExit: roomB})
	door := &Object{Desc: "door"}
	roomA.SetExit(West, ExitProps{DExit: door, RExit: roomC})

	G.Here = roomA
	G.Winner.MoveTo(roomA)
	G.Lit = true

	G.Params.HasWalkDir = true
	G.Params.WalkDir = North
	if !vWalk(ActUnk) || G.Here != roomB {
		t.Fatalf("expected functional exit to move to roomB")
	}

	out.Reset()
	G.Here = roomA
	G.Winner.MoveTo(roomA)
	G.Params.WalkDir = East
	G.PerformFatal = false
	vWalk(ActUnk)
	if !strings.Contains(out.String(), "Blocked.") {
		t.Fatalf("expected conditional exit message, got %q", out.String())
	}
	if !G.PerformFatal {
		t.Fatalf("expected conditional exit to be fatal")
	}

	out.Reset()
	G.Here = roomA
	G.Winner.MoveTo(roomA)
	G.Params.WalkDir = West
	G.PerformFatal = false
	vWalk(ActUnk)
	if !strings.Contains(out.String(), "closed") {
		t.Fatalf("expected closed door message, got %q", out.String())
	}
	if !G.PerformFatal {
		t.Fatalf("expected closed door to be fatal")
	}
}
