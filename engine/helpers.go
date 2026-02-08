package engine

// IsInGlobal checks if obj1 appears in obj2's Global list.
func IsInGlobal(obj1, obj2 *Object) bool {
	if obj2.Global == nil {
		return false
	}
	for _, o := range obj2.Global {
		if o == obj1 {
			return true
		}
	}
	return false
}

// IsHeld checks if an object is held (directly or indirectly) by the winner.
func IsHeld(obj *Object) bool {
	for {
		obj = obj.Location()
		if obj == nil {
			return false
		}
		if obj == G.Winner {
			return true
		}
	}
}
