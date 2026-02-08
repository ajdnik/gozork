package engine

import (
	"bufio"
	"sort"
	"strings"
)

var Vocabulary = make(map[string]WordItm)

// WordTyp defined the part-of-speech type
type WordTyp int

const (
	WordUnk WordTyp = iota
	WordDir
	WordVerb
	WordPrep
	WordAdj
	WordObj
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

func (wt WordTypes) Len() int           { return len(wt) }
func (wt WordTypes) Less(i, j int) bool { return int(wt[i]) < int(wt[j]) }
func (wt WordTypes) Swap(i, j int)      { wt[i], wt[j] = wt[j], wt[i] }

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

func (e *LexItm) Matches(itm LexItm) bool {
	return e.Norm == itm.Norm && e.Types.Equals(itm.Types)
}

func (e *LexItm) Set(itm LexItm) {
	e.Norm = itm.Norm
	e.Orig = itm.Orig
	e.Types = append(WordTypes{}, itm.Types...)
}

func (e *LexItm) IsSet() bool {
	return len(e.Norm) != 0 && len(e.Orig) != 0
}

func (e *LexItm) Is(wrd string) bool {
	return e.Norm == wrd
}

func (e *LexItm) IsAny(wrds ...string) bool {
	for _, wrd := range wrds {
		if e.Norm == wrd {
			return true
		}
	}
	return false
}

func (e *LexItm) Clear() {
	e.Norm = ""
	e.Orig = ""
	e.Types = nil
}

type WordItm struct {
	Norm  string
	Types WordTypes
}

// InitReader initializes the buffered reader from GameInput.
func InitReader() {
	G.Reader = bufio.NewReader(G.GameInput)
}

// Read reads input from the game input, tokenizes the input and tags parts-of-speech.
func Read() (string, []LexItm) {
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
	txt = strings.Replace(txt, "\n", "", -1)
	txt = strings.Replace(txt, "\r", "", -1)
	txt = strings.ToLower(txt)
	toks := Tokenize(txt)
	itms := Lex(toks)
	return txt, itms
}

func isLetter(c rune) bool { return 'a' <= c && c <= 'z' }
func isNum(c rune) bool    { return '0' <= c && c <= '9' }

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
