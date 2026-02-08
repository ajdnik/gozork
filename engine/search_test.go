package engine

import (
	"testing"
)

func TestSearchListFindTopAndBottom(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	G.Search.Syn.Norm = "gem"
	G.Search.Syn.Orig = "gem"
	G.Search.Syn.Types = WordTypes{WordObj}

	root := &Object{Desc: "root"}
	container := &Object{Desc: "box", Flags: FlgCont | FlgOpen | FlgSearch, In: root}
	itemTop := &Object{Desc: "gem", Synonyms: []string{"gem"}, In: root}
	itemNested := &Object{Desc: "gem", Synonyms: []string{"gem"}, In: container}

	root.AddChild(container)
	root.AddChild(itemTop)
	container.AddChild(itemNested)

	foundTop := SearchList(root, FindTop)
	if len(foundTop) != 2 {
		t.Fatalf("expected FindTop to return direct child gem and container contents")
	}
	if (foundTop[0] != itemTop && foundTop[1] != itemTop) || (foundTop[0] != itemNested && foundTop[1] != itemNested) {
		t.Fatalf("expected FindTop to include both gems")
	}

	foundAll := SearchList(root, FindAll)
	if len(foundAll) != 2 {
		t.Fatalf("expected FindAll to return both gems, got %d", len(foundAll))
	}

	foundBottom := SearchList(root, FindBottom)
	if len(foundBottom) != 1 || foundBottom[0] != itemNested {
		t.Fatalf("expected FindBottom to return nested gem")
	}
}

func TestDoSLWithLocFlags(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	G.Search.Syn.Norm = "coin"
	G.Search.Syn.Orig = "coin"
	G.Search.Syn.Types = WordTypes{WordObj}

	root := &Object{Desc: "root"}
	table := &Object{Desc: "table", Flags: FlgSurf | FlgOpen, In: root}
	coin := &Object{Desc: "coin", Synonyms: []string{"coin"}, In: table}
	root.AddChild(table)
	table.AddChild(coin)

	G.Search.LocFlags = LocSet(LocOnGrnd, LocInRoom)
	if res := DoSL(root, LocOnGrnd, LocInRoom); len(res) != 1 || res[0] != coin {
		t.Fatalf("expected DoSL to find coin with both flags")
	}

	G.Search.LocFlags = LocSet(LocOnGrnd)
	if res := DoSL(root, LocOnGrnd, LocInRoom); len(res) != 1 || res[0] != coin {
		t.Fatalf("expected DoSL to find coin with LocOnGrnd")
	}

	G.Search.LocFlags = LocSet(LocInRoom)
	if res := DoSL(root, LocOnGrnd, LocInRoom); len(res) != 1 || res[0] != coin {
		t.Fatalf("expected DoSL to find coin with LocInRoom")
	}

	G.Search.LocFlags = 0
	if res := DoSL(root, LocOnGrnd, LocInRoom); len(res) != 0 {
		t.Fatalf("expected DoSL to return empty when no flags set")
	}
}

// ---- search_flags_test.go ----

func TestSearchListRespectsOpenAndTransparent(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	G.Search.Syn = LexItem{Norm: "key", Orig: "key", Types: WordTypes{WordObj}}

	root := &Object{Desc: "root"}
	closed := &Object{Desc: "box", Flags: FlgCont, In: root}
	open := &Object{Desc: "chest", Flags: FlgCont | FlgOpen, In: root}
	transparent := &Object{Desc: "case", Flags: FlgCont | FlgTrans, In: root}
	keyClosed := &Object{Desc: "key", Synonyms: []string{"key"}, In: closed}
	keyOpen := &Object{Desc: "key", Synonyms: []string{"key"}, In: open}
	keyTrans := &Object{Desc: "key", Synonyms: []string{"key"}, In: transparent}

	root.AddChild(closed)
	root.AddChild(open)
	root.AddChild(transparent)
	closed.AddChild(keyClosed)
	open.AddChild(keyOpen)
	transparent.AddChild(keyTrans)

	found := SearchList(root, FindAll)
	if len(found) != 2 {
		t.Fatalf("expected SearchList to skip closed container, got %d", len(found))
	}
	if found[0] != keyOpen && found[1] != keyOpen {
		t.Fatalf("expected open container key to be found")
	}
	if found[0] != keyTrans && found[1] != keyTrans {
		t.Fatalf("expected transparent container key to be found")
	}
}

// ---- search_surface_test.go ----

func TestSearchListSurfaceSearchFlag(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	G.Search.Syn = LexItem{Norm: "coin", Orig: "coin", Types: WordTypes{WordObj}}

	root := &Object{Desc: "root"}
	table := &Object{Desc: "table", Flags: FlgSurf | FlgOpen, In: root}
	box := &Object{Desc: "box", Flags: FlgSearch | FlgOpen, In: root}
	coinOnTable := &Object{Desc: "coin", Synonyms: []string{"coin"}, In: table}
	coinInBox := &Object{Desc: "coin", Synonyms: []string{"coin"}, In: box}

	root.AddChild(table)
	root.AddChild(box)
	table.AddChild(coinOnTable)
	box.AddChild(coinInBox)

	found := SearchList(root, FindAll)
	if len(found) != 2 {
		t.Fatalf("expected FindTop to return coins in surf/search containers, got %d", len(found))
	}
	if (found[0] != coinOnTable && found[1] != coinOnTable) || (found[0] != coinInBox && found[1] != coinInBox) {
		t.Fatalf("expected to find both coins")
	}
}

func TestSearchListSearchFlagClosedContainer(t *testing.T) {
	oldG := G
	G = NewGameState()
	t.Cleanup(func() { G = oldG })

	G.Search.Syn = LexItem{Norm: "coin", Orig: "coin", Types: WordTypes{WordObj}}

	root := &Object{Desc: "root"}
	box := &Object{Desc: "box", Flags: FlgSearch, In: root}
	coin := &Object{Desc: "coin", Synonyms: []string{"coin"}, In: box}
	root.AddChild(box)
	box.AddChild(coin)

	found := SearchList(root, FindAll)
	if len(found) != 0 {
		t.Fatalf("expected closed search container to be skipped, got %d", len(found))
	}
}
