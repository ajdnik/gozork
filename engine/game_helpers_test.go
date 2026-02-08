package engine

import (
	"testing"
)

func TestCallHandlerNotHandled(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	res, handled := callHandler(func(ActionArg) bool { return false }, ActBegin)
	if handled {
		t.Fatalf("expected handler to report not handled")
	}
	if res != PerfNotHndld {
		t.Fatalf("expected PerfNotHndld, got %v", res)
	}
}

func TestCallHandlerHandled(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	res, handled := callHandler(func(ActionArg) bool { return true }, ActBegin)
	if !handled {
		t.Fatalf("expected handler to be handled")
	}
	if res != PerfHndld {
		t.Fatalf("expected PerfHndld, got %v", res)
	}
}

func TestCallHandlerQuit(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	res, handled := callHandler(func(ActionArg) bool {
		Quit()
		return false
	}, ActBegin)
	if !handled {
		t.Fatalf("expected quit to be handled")
	}
	if res != PerfQuit {
		t.Fatalf("expected PerfQuit, got %v", res)
	}
}

func TestCallHandlerFatal(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	res, handled := callHandler(func(ActionArg) bool { return RFatal() }, ActBegin)
	if !handled {
		t.Fatalf("expected fatal to be handled")
	}
	if res != PerfFatal {
		t.Fatalf("expected PerfFatal, got %v", res)
	}
}

func TestVerify(t *testing.T) {
	if !Verify() {
		t.Fatalf("expected Verify to return true")
	}
}
