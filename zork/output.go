package zork

import "fmt"

type EndType int

const (
	Newline EndType = iota
	NoNewline
)

func NewLine() {
	fmt.Fprintln(G.GameOutput, "")
}

func PrintObject(obj *Object) {
	fmt.Fprintf(G.GameOutput, "%v", obj.Desc)
}

func PrintNumber(num int) {
	fmt.Fprintf(G.GameOutput, "%v", num)
}

func Print(msg string, end EndType) {
	if end == NoNewline {
		fmt.Fprintf(G.GameOutput, "%s", msg)
	} else {
		fmt.Fprintln(G.GameOutput, msg)
	}
}
