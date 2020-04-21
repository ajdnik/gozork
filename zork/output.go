package zork

import "fmt"

func NewLine() {
	fmt.Printf("\n")
}

func Print(msg string) {
	fmt.Printf(msg)
}

func PrintNumber(n uint32) {
	fmt.Printf("%v", n)
}

func PrintChar(n byte) {
	fmt.Printf("%v", string(n))
}

func PrintObj(o *Object) {
	fmt.Printf(o.Name)
}
