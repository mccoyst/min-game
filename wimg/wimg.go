// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
package main

import (
	"bufio"
	"flag"
	"minima/world"
	"os"
)

var (
	outFile = flag.String("o", "world.png", "The output file")
	echo    = flag.Bool("e", false, "Echo the world to standard output")
	depth = flag.Bool("d", true, "Draw water depth")
)

func main() {
	flag.Parse()

	w, err := world.Read(bufio.NewReader(os.Stdin))
	if err != nil {
		panic(err)
	}

	if err := w.SavePng(*outFile, *depth); err != nil {
		panic(err)
	}

	if *echo {
		out := bufio.NewWriter(os.Stdout)
		defer out.Flush()
		if err := w.Write(out); err != nil {
			panic(err)
		}
	}
}

