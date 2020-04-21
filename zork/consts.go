package zork

var (
	Turns         = 0
	Location      *Object
	Lit           = false
	Player        *Object
	Actor         *Object
	LastNoun      *Object
	LastNounPlace *Object
	DollarZero    = [24]byte{0, 0, 0, 0, 88, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 56, 52, 48, 55, 50, 54}
	Verbose       = true
	Superbrief    = false
	GrueRepellent = true
	WonFlag       = false
	Item4         *Object
	Rooms         = []*Object{WestOfHouse}
	Indents       = [6]string{
		"  ",
		"  ",
		"    ",
		"      ",
		"        ",
		"          ",
	}
	LastObLongdesc *Object
	Verb           = ""
	Noun           *Object

	Def3C = [10]int{1, 6, 6, 4, 4, 4, 4, 5, 5, 5}
	Def3B = [11]int{1, 1, 1, 6, 6, 4, 4, 4, 5, 5, 5}
	Def3A = [11]int{1, 1, 1, 1, 1, 6, 6, 4, 4, 5, 5}
	Def2B = [12]int{1, 1, 1, 6, 6, 4, 4, 4, 2, 3, 3, 3}
	Def2A = [10]int{1, 1, 1, 1, 1, 6, 6, 4, 4, 2}
	Def1  = [13]int{1, 1, 1, 1, 6, 6, 2, 2, 3, 3, 3, 3, 3}

	Def3Res = [5]int{Def3A[0], 0, Def3B[0], 0, Def3C[0]}
	Def2Res = [4]int{Def2A[0], Def2B[0], 0, 0}
	Def1Res = [3]int{Def1[0], 0, 0}
)

type Property string

const (
	Visited     Property = "visited"
	MazeRoom    Property = "maze_room"
	Vehicle     Property = "vehicle"
	DryLand     Property = "dry_land"
	Light       Property = "light"
	Sacred      Property = "sacred"
	Concealed   Property = "concealed"
	Scenery     Property = "scenery"
	Transparent Property = "transparent"
	Open        Property = "open"
	Supporter   Property = "supporter"
	Animate     Property = "animate"
	Clothing    Property = "clothing"
	Container   Property = "container"
	TryTakeBit  Property = "trytakebit"
)
