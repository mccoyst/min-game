// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"

	"code.google.com/p/min-game/ui"
)

// Command flags
var (
	drawHeights  = flag.Bool("heights", false, "draw tile height values — SLOW")
	worldOnStdin = flag.Bool("stdin", false, "read the world from stdin")
	profile      = flag.Bool("profile", false, "enable CPU profiling to ./prof.txt")
	dvorak       = flag.Bool("dvorak", false, "use a Dvorak key map")
	locInfo      = flag.Bool("locinfo", false, "display current location info")
)

var ScreenDims = ui.Pt(640, 480)

func init() {
	runtime.LockOSThread()
}

func main() {
	flag.Parse()

	if *profile {
		p, err := os.Create("./prof.txt")
		if err != nil {
			os.Stderr.WriteString("oops: " + err.Error() + "\n")
			os.Exit(1)
		}
		pprof.StartCPUProfile(p)
		defer pprof.StopCPUProfile()
	}

	if *dvorak {
		ui.DefaultKeymap = ui.DvorakKeymap
	}

	u, err := ui.New("minima", int(ScreenDims.X), int(ScreenDims.Y))
	if err != nil {
		os.Stderr.WriteString("oops: " + err.Error() + "\n")
		os.Exit(1)
	}
	defer u.Close()

	stk := NewScreenStack(u, NewTitleScreen())
	stk.Run()
	fmt.Printf("mean frame time: %4.1fms\n", stk.meanFrame)
}
