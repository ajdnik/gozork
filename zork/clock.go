package zork

type RtnFunc func() bool

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
	for i := len(QueueItms) - 1; i > QueueInts; i-- {
		// Compare two function pointers,
		// might not be the best way to do it but I wanted
		// to stay as true to the original code as possible.
		if &QueueItms[i].Rtn == &rtn {
			return &QueueItms[i]
		}
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
