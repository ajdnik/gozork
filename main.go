package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ajdnik/gozork/game"
)

var version = "dev"

func main() {
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Println("gozork version " + version)
		os.Exit(0)
	}

	game.Run()
}
