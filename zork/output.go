package zork

import (
	"fmt"
	"io"
	"os"
)

type EndType int

const (
	Newline EndType = iota
	NoNewline
)

// GameOutput is the writer all game output goes to. Defaults to os.Stdout.
// Tests can replace this to capture output.
var GameOutput io.Writer = os.Stdout

func NewLine() {
	fmt.Fprintln(GameOutput, "")
}

func PrintObject(obj *Object) {
	fmt.Fprintf(GameOutput, "%v", obj.Desc)
}

func PrintNumber(num int) {
	fmt.Fprintf(GameOutput, "%v", num)
}

func Print(msg string, end EndType) {
	if end == NoNewline {
		fmt.Fprintf(GameOutput, "%s", msg)
	} else {
		fmt.Fprintln(GameOutput, msg)
	}
}
