package zork

import (
	"bufio"
	"os"
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

// LexItm is an object that is returned after lexing
type LexItm struct {
	Norm string
	Orig string
	Type WordTyp
}

// Matches compares lex item to another lex item
func (e *LexItm) Matches(itm LexItm) bool {
	return e.Norm == itm.Norm && e.Type == itm.Type
}

// Set coppies one lex item into another
func (e *LexItm) Set(itm LexItm) {
	e.Norm = itm.Norm
	e.Orig = itm.Orig
	e.Type = itm.Type
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
	e.Type = WordUnk
}

// WordItm represents words stored in the global
// map called Vocabulary which is used when tagging
// parts of speach
type WordItm struct {
	Norm string
	Type WordTyp
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
				Norm: val.Norm,
				Orig: tok,
				Type: val.Type,
			})
		} else {
			itms = append(itms, LexItm{
				Norm: tok,
				Orig: tok,
				Type: WordUnk,
			})
		}
	}
	return itms
}
