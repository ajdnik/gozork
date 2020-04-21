package zork

type Daemon func() bool

type QueueEntry struct {
	Enabled bool
	Ticks   int
	Id      string
	Action  Daemon
}

var (
	CTable    [30]QueueEntry
	CInts     = 30
	CInts2    = 30
	ClockWait = false
)

func Queue(id string, rtn Daemon, ticks int) *QueueEntry {
	cint := QueueInterrupt(id, rtn, false)
	cint.Ticks = ticks
	return cint
}

func QueueInterrupt(id string, rtn Daemon, int2flag bool) *QueueEntry {
	for i := CInts; i < len(CTable); i++ {
		if CTable[i].Id == id {
			return &CTable[i]
		}
	}
	CInts--
	if int2flag {
		CInts2--
	}
	CTable[CInts].Enabled = false
	CTable[CInts].Ticks = 0
	CTable[CInts].Id = id
	CTable[CInts].Action = rtn
	return &CTable[CInts]
}

func Clocker() bool {
	if ClockWait {
		ClockWait = false
		return false
	}
	flag := false
	i := CInts2
	if PWon {
		i = CInts
	}
	for ; i < len(CTable); i++ {
		if !CTable[i].Enabled {
			continue
		}
		if CTable[i].Ticks == 0 {
			continue
		}
		CTable[i].Ticks--
		if CTable[i].Ticks > 1 || !CTable[i].Action() {
			continue
		}
		flag = true
	}
	Turns++
	if Turns > 999 {
		Turns = 0
	}
	return flag
}
