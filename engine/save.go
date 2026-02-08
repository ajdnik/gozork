package engine

// ObjIndex maps each *Object in G.AllObjects to its index. Built lazily.
var ObjIndex map[*Object]int

// BuildObjIndex populates the ObjIndex map from G.AllObjects. It is
// idempotent and only builds the index once.
func BuildObjIndex() {
	if ObjIndex != nil {
		return
	}
	ObjIndex = make(map[*Object]int, len(G.AllObjects)+1)
	for i, o := range G.AllObjects {
		ObjIndex[o] = i
	}
}

// ObjToIdx returns the index of an object in G.AllObjects, or -1 if nil/unknown.
func ObjToIdx(o *Object) int {
	if o == nil {
		return -1
	}
	if idx, ok := ObjIndex[o]; ok {
		return idx
	}
	return -1
}

// IdxToObj returns the object at the given index in G.AllObjects, or nil if out of range.
func IdxToObj(idx int) *Object {
	if idx < 0 || idx >= len(G.AllObjects) {
		return nil
	}
	return G.AllObjects[idx]
}
