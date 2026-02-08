package engine

import (
	"bufio"
	"sort"
	"strings"
)

// Vocabulary maps known words to their lexical data. Populated during initialization.
var Vocabulary = make(map[string]WordItem)

// WordTyp represents the part-of-speech type of a word.
type WordTyp int

const (
	// WordUnk is the zero value for an unrecognized word.
	WordUnk WordTyp = iota
	// WordDir indicates a compass direction (north, south, etc.).
	WordDir
	// WordVerb indicates a verb (take, open, etc.).
	WordVerb
	// WordPrep indicates a preposition (in, on, with, etc.).
	WordPrep
	// WordAdj indicates an adjective (brass, large, etc.).
	WordAdj
	// WordObj indicates a noun / object name (lamp, sword, etc.).
	WordObj
	// WordBuzz indicates a filler word that is ignored (the, a, etc.).
	WordBuzz
)

// WordTypes is a list of part-of-speech tags for a single word.
type WordTypes []WordTyp

// Has returns true if typ appears in the list.
func (wt WordTypes) Has(typ WordTyp) bool {
	for _, t := range wt {
		if t == typ {
			return true
		}
	}
	return false
}

// Len implements sort.Interface.
func (wt WordTypes) Len() int { return len(wt) }

// Less implements sort.Interface.
func (wt WordTypes) Less(i, j int) bool { return int(wt[i]) < int(wt[j]) }

// Swap implements sort.Interface.
func (wt WordTypes) Swap(i, j int) { wt[i], wt[j] = wt[j], wt[i] }

// Equals returns true if both type lists contain the same elements.
func (wt WordTypes) Equals(tt WordTypes) bool {
	if len(wt) != len(tt) {
		return false
	}
	sort.Sort(wt)
	sort.Sort(tt)
	for idx := range wt {
		if wt[idx] != tt[idx] {
			return false
		}
	}
	return true
}

// LexItem is an object that is returned after lexing
type LexItem struct {
	Norm  string
	Orig  string
	Types WordTypes
}

// Matches returns true if both items have the same norm and types.
func (e *LexItem) Matches(itm LexItem) bool {
	return e.Norm == itm.Norm && e.Types.Equals(itm.Types)
}

// Set copies another LexItem's data into this one.
func (e *LexItem) Set(itm LexItem) {
	e.Norm = itm.Norm
	e.Orig = itm.Orig
	e.Types = append(WordTypes{}, itm.Types...)
}

// IsSet returns true if the item has been populated with data.
func (e *LexItem) IsSet() bool {
	return len(e.Norm) != 0 && len(e.Orig) != 0
}

// Is returns true if the normalized form equals wrd.
func (e *LexItem) Is(wrd string) bool {
	return e.Norm == wrd
}

// IsAny returns true if the normalized form matches any of the given words.
func (e *LexItem) IsAny(wrds ...string) bool {
	for _, wrd := range wrds {
		if e.Norm == wrd {
			return true
		}
	}
	return false
}

// Clear resets the item to its zero state.
func (e *LexItem) Clear() {
	e.Norm = ""
	e.Orig = ""
	e.Types = nil
}

// WordItem is the vocabulary entry for a known word.
type WordItem struct {
	Norm  string
	Types WordTypes
}

// InitReader initializes the buffered reader from GameInput.
func InitReader() {
	G.Reader = bufio.NewReader(G.GameInput)
}

// Read reads input from the game input, tokenizes the input and tags parts-of-speech.
func Read() (string, []LexItem) {
	if G.Reader == nil {
		InitReader()
	}
	txt, err := G.Reader.ReadString('\n')
	if err != nil {
		G.InputExhausted = true
		if len(txt) == 0 {
			return "", nil
		}
	}
	txt = strings.ReplaceAll(txt, "\n", "")
	txt = strings.ReplaceAll(txt, "\r", "")
	txt = strings.ToLower(txt)
	toks := Tokenize(txt)
	itms := Lex(toks)
	return txt, itms
}

func isLetter(c rune) bool { return 'a' <= c && c <= 'z' }
func isNum(c rune) bool    { return '0' <= c && c <= '9' }

// Tokenize splits raw input into whitespace-delimited tokens,
// separating letters, digits, and punctuation.
func Tokenize(buf string) []string {
	toks := []string{}
	var cur string
	for _, c := range buf {
		if c == ' ' {
			toks = append(toks, cur)
			cur = ""
			continue
		}
		if len(cur) == 0 {
			cur += string(c)
			continue
		}
		p := rune(cur[len(cur)-1])
		if isLetter(c) && !isLetter(p) {
			toks = append(toks, cur)
			cur = ""
		} else if isNum(c) && !isNum(p) {
			toks = append(toks, cur)
			cur = ""
		} else if !isNum(c) && !isLetter(c) && (isLetter(p) || isNum(p)) {
			toks = append(toks, cur)
			cur = ""
		}
		cur += string(c)
	}
	if len(cur) != 0 {
		toks = append(toks, cur)
	}
	return toks
}

// Lex looks up each token in the Vocabulary and returns tagged LexItem entries.
func Lex(toks []string) []LexItem {
	itms := []LexItem{}
	for _, tok := range toks {
		if val, ok := Vocabulary[tok]; ok {
			itms = append(itms, LexItem{
				Norm:  val.Norm,
				Orig:  tok,
				Types: append(WordTypes{}, val.Types...),
			})
		} else {
			itms = append(itms, LexItem{
				Norm:  tok,
				Orig:  tok,
				Types: nil,
			})
		}
	}
	return itms
}
