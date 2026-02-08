package engine

// ClockEvent represents a timed or daemon interrupt in the game's event queue.
type ClockEvent struct {
	Key  string      // unique identifier for this event
	Run  bool        // whether the event is active
	Tick int         // countdown; fires when it reaches 0
	Fn   func() bool // the function to call when the event fires
}

// Queue finds (or creates) a clock event by key and sets its tick value.
func Queue(key string, tick int) *ClockEvent {
	ev := QueueInt(key, false)
	ev.Tick = tick
	return ev
}

// QueueInt finds an existing clock event by key, or allocates a new slot.
// If dmn is true, the event is a daemon (always ticked, even on bad parses).
func QueueInt(key string, dmn bool) *ClockEvent {
	for i := len(G.QueueItms) - 1; i >= G.QueueInts; i-- {
		if G.QueueItms[i].Key == key {
			return &G.QueueItms[i]
		}
	}
	if G.QueueInts <= 0 {
		return &G.QueueItms[0]
	}
	G.QueueInts--
	if dmn {
		G.QueueDmns--
	}
	G.QueueItms[G.QueueInts] = ClockEvent{Key: key, Fn: G.ClockFuncs[key]}
	return &G.QueueItms[G.QueueInts]
}

// Clocker runs all active clock events. Called once per turn.
func Clocker() bool {
	if G.ClockWait {
		G.ClockWait = false
		return false
	}
	end := G.QueueDmns
	if G.ParserOk {
		end = G.QueueInts
	}
	flg := false
	for i := len(G.QueueItms) - 1; i >= end; i-- {
		if G.QuitRequested {
			return flg
		}
		if !G.QueueItms[i].Run {
			continue
		}
		if G.QueueItms[i].Tick == 0 {
			continue
		}
		G.QueueItms[i].Tick--
		if G.QueueItms[i].Tick <= 0 && G.QueueItms[i].Fn() {
			flg = true
		}
	}
	G.Moves++
	return flg
}
