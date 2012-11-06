// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"os"
	"runtime"
	"time"

	"code.google.com/p/min-game/ui"
)

var ScreenDims = ui.Pt(640, 480)

func init() {
	runtime.LockOSThread()
}

func main() {
	u, err := ui.NewUi("minima", int(ScreenDims.X), int(ScreenDims.Y))
	if err != nil {
		os.Stderr.WriteString("oops: " + err.Error() + "\n")
		os.Exit(1)
	}
	defer u.Quit()

	for running := true; running; {
		frameStart := time.Now()
		for {
			e := u.PollEvent()
			if e == nil {
				break
			}
			switch event := e.(type) {
			case ui.Quit:
				running = false
				break
			case ui.Key:
				if event.Repeat {
					continue
				}
			}
		}

		u.SetColor(0, 0, 0, 255)
		u.Clear()
		u.Show()

		frameLen := time.Now().Sub(frameStart)
		if frameLen < 16*time.Millisecond {
			time.Sleep(16*time.Millisecond - frameLen)
		}
	}
}
