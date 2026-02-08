package engine

// PickOne selects a random string from data without repeating until all have been used.
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

// Random returns a uniformly random element from the given object slice.
func Random(tbl []*Object) *Object {
	return tbl[G.Rand.Intn(len(tbl))]
}

// Prob returns true with the given probability (0-100).
func Prob(base int, isLooser bool) bool {
	if isLooser {
		return Zprob(base)
	}
	return base > G.Rand.Intn(100)+1
}

// Zprob is the "loser" version of Prob that accounts for the Lucky flag.
func Zprob(base int) bool {
	if G.Lucky {
		return base > G.Rand.Intn(100)+1
	}
	return base > G.Rand.Intn(300)+1
}

// IsFlaming returns true if the object is an active fire source.
func IsFlaming(obj *Object) bool {
	return obj.Has(FlgFlame) && obj.Has(FlgOn)
}

// IsOpenable returns true if the object is a door or container.
func IsOpenable(obj *Object) bool {
	return obj.Has(FlgDoor) || obj.Has(FlgCont)
}
