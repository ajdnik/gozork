package zork

import (
	"bufio"
	"os"
	"sort"
	"strings"
)

// WordTyp defined the part-of-speech type
type WordTyp int

const (
	// WordUnk is used when part-of-speach is unknown
	WordUnk WordTyp = iota
	// WordDir represents a direction
	WordDir
	// WordVerb represents a verb
	WordVerb
	// WordPrep represents a preposition
	WordPrep
	// WordAdj represents an adjective
	WordAdj
	// WordObj represents an object (game object)
	WordObj
	// WordBuzz represents a buzz word (a filler)
	WordBuzz
)

type WordTypes []WordTyp

func (wt WordTypes) Has(typ WordTyp) bool {
	for _, t := range wt {
		if t == typ {
			return true
		}
	}
	return false
}

func (wt WordTypes) Len() int {
	return len(wt)
}

func (wt WordTypes) Less(i, j int) bool {
	return int(wt[i]) < int(wt[j])
}

func (wt WordTypes) Swap(i, j int) {
	wt[i], wt[j] = wt[j], wt[i]
}

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

// LexItm is an object that is returned after lexing
type LexItm struct {
	Norm  string
	Orig  string
	Types WordTypes
}

// Matches compares lex item to another lex item
func (e *LexItm) Matches(itm LexItm) bool {
	return e.Norm == itm.Norm && e.Types.Equals(itm.Types)
}

// Set coppies one lex item into another
func (e *LexItm) Set(itm LexItm) {
	e.Norm = itm.Norm
	e.Orig = itm.Orig
	e.Types = append(WordTypes{}, itm.Types...)
}

// IsSet checks if lex item is defined
func (e *LexItm) IsSet() bool {
	return len(e.Norm) != 0 && len(e.Orig) != 0
}

// Is checks if lex item matches the word provided.
func (e *LexItm) Is(wrd string) bool {
	return e.Norm == wrd
}

// IsAny checks if the lex item matches any of the
// words provided as arguments.
func (e *LexItm) IsAny(wrds ...string) bool {
	for _, wrd := range wrds {
		if e.Norm == wrd {
			return true
		}
	}
	return false
}

// Clear resets the lex item
func (e *LexItm) Clear() {
	e.Norm = ""
	e.Orig = ""
	e.Types = nil
}

// WordItm represents words stored in the global
// map called Vocabulary which is used when tagging
// parts of speach
type WordItm struct {
	Norm  string
	Types WordTypes
}

var (
	Reader     = bufio.NewReader(os.Stdin)
	Vocabulary = make(map[string]WordItm)
)

// Read function reads input from stdin,
// tokenizes the input and taggs parts-of-speach
func Read() (string, []LexItm) {
	txt, _ := Reader.ReadString('\n')
	txt = strings.Replace(txt, "\n", "", -1)
	txt = strings.ToLower(txt)
	toks := Tokenize(txt)
	itms := Lex(toks)
	return txt, itms
}

func isLetter(c rune) bool {
	return 'a' <= c && c <= 'z'
}

func isNum(c rune) bool {
	return '0' <= c && c <= '9'
}

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

func Lex(toks []string) []LexItm {
	itms := []LexItm{}
	for _, tok := range toks {
		if val, ok := Vocabulary[tok]; ok {
			itms = append(itms, LexItm{
				Norm:  val.Norm,
				Orig:  tok,
				Types: append(WordTypes{}, val.Types...),
			})
		} else {
			itms = append(itms, LexItm{
				Norm:  tok,
				Orig:  tok,
				Types: nil,
			})
		}
	}
	return itms
}
