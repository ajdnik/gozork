package zork

func PickOne(data RndSelect) string {
	if len(data.Unselected) == 0 {
		data.Unselected = data.Selected
		data.Selected = []string{}
	}
	if data.Selected == nil {
		data.Selected = []string{}
	}
	rnd := G.Rand.Intn(len(data.Unselected))
	msg := data.Unselected[rnd]
	data.Selected = append(data.Selected, msg)
	data.Unselected[rnd] = data.Unselected[len(data.Unselected)-1]
	data.Unselected[len(data.Unselected)-1] = ""
	data.Unselected = data.Unselected[:len(data.Unselected)-1]
	return msg
}

func Random(tbl []*Object) *Object {
	return tbl[G.Rand.Intn(len(tbl))]
}

// Prob returns true with the given probability (0-100).
// ZIL's PROB macro: <G? .BASE? <RANDOM 100>> where RANDOM returns 1..100.
// We use G.Rand.Intn(100)+1 to match the 1..100 range.
func Prob(base int, isLooser bool) bool {
	if isLooser {
		return Zprob(base)
	}
	return base > G.Rand.Intn(100)+1
}

// Zprob is the "loser" version of Prob that accounts for the Lucky flag.
// ZIL: <G? .BASE <RANDOM 100>> when lucky, <G? .BASE <RANDOM 300>> otherwise.
func Zprob(base int) bool {
	if G.Lucky {
		return base > G.Rand.Intn(100)+1
	}
	return base > G.Rand.Intn(300)+1
}

func IsFlaming(obj *Object) bool {
	return obj.Has(FlgFlame) && obj.Has(FlgOn)
}

func IsOpenable(obj *Object) bool {
	return obj.Has(FlgDoor) || obj.Has(FlgCont)
}
