package engine

import "fmt"

// Printf writes formatted output to the game's output writer.
func Printf(format string, a ...interface{}) {
	fmt.Fprintf(G.GameOutput, format, a...)
}
