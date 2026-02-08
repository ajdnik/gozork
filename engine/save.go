package engine

// ObjIndex maps each *Object in G.AllObjects to its index. Built lazily.
var ObjIndex map[*Object]int

func BuildObjIndex() {
	if ObjIndex != nil {
		return
	}
	ObjIndex = make(map[*Object]int, len(G.AllObjects)+1)
	for i, o := range G.AllObjects {
		ObjIndex[o] = i
	}
}

func ObjToIdx(o *Object) int {
	if o == nil {
		return -1
	}
	if idx, ok := ObjIndex[o]; ok {
		return idx
	}
	return -1
}

func IdxToObj(idx int) *Object {
	if idx < 0 || idx >= len(G.AllObjects) {
		return nil
	}
	return G.AllObjects[idx]
}
