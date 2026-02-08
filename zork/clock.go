package zork

import "reflect"

type RtnFunc func() bool

// funcAddr returns the function pointer address for comparison.
func funcAddr(f RtnFunc) uintptr {
	return reflect.ValueOf(f).Pointer()
}

type QueueItm struct {
	Run  bool
	Tick int
	Rtn  RtnFunc
}


func Queue(rtn RtnFunc, tick int) *QueueItm {
	itm := QueueInt(rtn, false)
	itm.Tick = tick
	return itm
}

func QueueInt(rtn RtnFunc, dmn bool) *QueueItm {
	// Search existing queue entries for a matching function.
	// We use reflect to compare function pointers since Go doesn't
	// allow direct comparison of function values.
	rtnPtr := funcAddr(rtn)
	for i := len(G.QueueItms) - 1; i >= G.QueueInts; i-- {
		if G.QueueItms[i].Rtn != nil && funcAddr(G.QueueItms[i].Rtn) == rtnPtr {
			return &G.QueueItms[i]
		}
	}
	if G.QueueInts <= 0 {
		// Queue is full, reuse the last slot
		return &G.QueueItms[0]
	}
	G.QueueInts--
	if dmn {
		G.QueueDmns--
	}
	G.QueueItms[G.QueueInts].Rtn = rtn
	G.QueueItms[G.QueueInts].Run = false
	G.QueueItms[G.QueueInts].Tick = 0
	return &G.QueueItms[G.QueueInts]
}

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
		if !G.QueueItms[i].Run {
			continue
		}
		if G.QueueItms[i].Tick == 0 {
			continue
		}
		G.QueueItms[i].Tick--
		if G.QueueItms[i].Tick <= 0 && G.QueueItms[i].Rtn() {
			flg = true
		}
	}
	G.Moves++
	return flg
}
