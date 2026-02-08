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

var (
	QueueItms [30]QueueItm
	QueueInts = 30
	QueueDmns = 30
	ClockWait = false
)

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
	for i := len(QueueItms) - 1; i >= QueueInts; i-- {
		if QueueItms[i].Rtn != nil && funcAddr(QueueItms[i].Rtn) == rtnPtr {
			return &QueueItms[i]
		}
	}
	if QueueInts <= 0 {
		// Queue is full, reuse the last slot
		return &QueueItms[0]
	}
	QueueInts--
	if dmn {
		QueueDmns--
	}
	QueueItms[QueueInts].Rtn = rtn
	QueueItms[QueueInts].Run = false
	QueueItms[QueueInts].Tick = 0
	return &QueueItms[QueueInts]
}

func Clocker() bool {
	if ClockWait {
		ClockWait = false
		return false
	}
	end := QueueDmns
	if ParserOk {
		end = QueueInts
	}
	flg := false
	for i := len(QueueItms) - 1; i > end; i-- {
		if !QueueItms[i].Run {
			continue
		}
		if QueueItms[i].Tick == 0 {
			continue
		}
		QueueItms[i].Tick--
		if QueueItms[i].Tick <= 1 && QueueItms[i].Rtn() {
			flg = true
		}
	}
	Moves++
	return flg
}
