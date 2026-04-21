package engine

import "testing"

func FuzzTokenize(f *testing.F) {
	f.Add("go north")
	f.Add("take the brass lamp")
	f.Add("go2west, now!")
	f.Add("a1b..c")
	f.Add("")
	f.Add("   ")
	f.Add("12:30")
	f.Add("open the mailbox and read the leaflet")
	f.Add("put sword in case")

	f.Fuzz(func(t *testing.T, input string) {
		// Verify Tokenize does not panic on arbitrary input.
		_ = Tokenize(input)
	})
}
