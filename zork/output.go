package zork

import "fmt"

type EndType int

const (
	Newline EndType = iota
	NoNewline
)

func NewLine() {
	fmt.Println("")
}

func PrintObject(obj *Object) {
	fmt.Printf("%v", obj.Desc)
}

func PrintNumber(num int) {
	fmt.Printf("%v", num)
}

func Print(msg string, end EndType) {
	if end == NoNewline {
		fmt.Printf(msg)
	} else {
		fmt.Println(msg)
	}
}
