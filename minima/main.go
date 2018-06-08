// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/mccoyst/min-game/geom"
	"github.com/mccoyst/min-game/resrc"
	"github.com/mccoyst/min-game/ui"
)

// Command flags
var (
	drawHeights  = flag.Bool("heights", false, "draw tile height values — SLOW")
	worldOnStdin = flag.Bool("stdin", false, "read the world from stdin")
	profile      = flag.Bool("profile", false, "enable CPU profiling to ./prof.txt")
	dvorak       = flag.Bool("dvorak", false, "use a Dvorak key map")
	debug        = flag.Bool("debug", false, "turn on debug printing")
	vsyncoff     = flag.Bool("vsyncoff", false, "turn off vsyncing")
	seed         = flag.Int64("seed", time.Now().UnixNano(), "The random seed")
)

var ScreenDims = geom.Pt(640, 480)

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
		ui.CurrentKeymap = ui.DvorakKeymap
	}

	u, err := ui.New("minima", int(ScreenDims.X), int(ScreenDims.Y), resrc.NewPkgFinder(), !*vsyncoff)
	if err != nil {
		os.Stderr.WriteString("oops: " + err.Error() + "\n")
		os.Exit(1)
	}
	defer u.Close()

	stk := ui.NewScreenStack(u, NewTitleScreen())
	stk.Run()
	fmt.Printf("mean frame time: %4.1fms\n", stk.MeanFrame)
}
