// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"os"
	"runtime"
	"time"
)

var ScreenDims = Pt(F(640), F(480))

func init() {
	runtime.LockOSThread()
}

func main() {
	ui, err := NewUi("minima", ScreenDims.X.Whole(), ScreenDims.Y.Whole())
	if err != nil {
		os.Stderr.WriteString("oops: " + err.Error() + "\n")
		os.Exit(1)
	}
	defer ui.Quit()

	for running := true; running; {
		frameStart := time.Now()
		for {
			e := ui.PollEvent()
			if e == nil {
				break
			}
			switch event := e.(type) {
			case Quit:
				running = false
				break
			case Key:
				if event.Repeat {
					continue
				}
			}
		}

		ui.SetColor(0, 0, 0, 255)
		ui.Clear()
		ui.Show()

		frameLen := time.Now().Sub(frameStart)
		if frameLen < 16*time.Millisecond {
			time.Sleep(16*time.Millisecond - frameLen)
		}
	}
}
