// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"image/color"
	"time"

	"code.google.com/p/min-game/ui"
	"code.google.com/p/min-game/geom"
)

// A Drawer can draw things and change colors.
type Drawer interface {
	Draw(interface{}, geom.Point) (geom.Point, error)
	SetFont(name string, szPts float64) error
	SetColor(color.Color)
	TextSize(string) geom.Point
}

// A Screen represents some game screen. E.g. the title, the main gameplay, etc.
type Screen interface {
	// Draw should send draw commands via the given Writer.
	Draw(Drawer) error

	// Handle is called for each event coming from the
	// Ui, with the exception of the Close event which is
	// intercepted by the ScreenStack to exit the program.
	Handle(*ScreenStack, ui.Event) error

	// Update is called after all of the events are handled and after
	// the next frame is drawn in order to allow the screen to update
	// its state based on the events.
	Update(*ScreenStack) error

	// Transparent returns true iff the screen doesn't fill the window.
	Transparent() bool
}

// A ScreenStack holds the stack of game screens.
type ScreenStack struct {
	stk       []Screen
	win       *ui.Ui
	nFrames   uint
	meanFrame float64 // milliseconds

	// Buttons is a bit set of the currently pressed buttons.
	buttons ui.Button
}

// NewScreenStack returns a new screen stack with the given initial screen.
func NewScreenStack(win *ui.Ui, first Screen) *ScreenStack {
	return &ScreenStack{
		stk:       []Screen{first},
		win:       win,
		nFrames:   0,
		meanFrame: 0.0,
	}
}

const FrameMsec = 16 * time.Millisecond

// Run runs the main loop of the program, calling the
// Draw(), Handle(), then Update() methods on the top
// screen on the stack.
func (s *ScreenStack) Run() {
	for {
		frameStart := time.Now()
		for {
			e := s.win.PollEvent()
			if e == nil {
				break
			}

			switch k := e.(type) {
			case ui.Quit:
				return

			case ui.Key:
				if k.Button == ui.Unknown {
					break
				} else if k.Down {
					s.buttons |= k.Button
				} else {
					s.buttons &^= k.Button
				}
			}
			if err := s.top().Handle(s, e); err != nil {
				panic(err)
			}
			if len(s.stk) == 0 {
				return
			}
		}

		s.win.SetColor(color.Black)
		s.win.Clear()
		if s.top().Transparent() && len(s.stk) > 1 {
			if err := s.stk[len(s.stk)-2].Draw(s.win); err != nil {
				panic(err)
			}
		}

		if err := s.top().Draw(s.win); err != nil {
			panic(err)
		}
		s.win.Sync()

		if err := s.top().Update(s); err != nil {
			panic(err)
		}
		if len(s.stk) == 0 {
			return
		}

		frameLen := time.Now().Sub(frameStart)
		if frameLen < FrameMsec {
			time.Sleep(FrameMsec - frameLen)
		}
		s.nFrames++
		ms := frameLen.Seconds() * 1000
		s.meanFrame += (ms - s.meanFrame) / float64(s.nFrames)
	}
}

// Push pushes a new screen onto the top of the stack.
func (s *ScreenStack) Push(screen Screen) {
	s.stk = append(s.stk, screen)
}

// Pop pops the current screen off of the top of the stack.
func (s *ScreenStack) Pop() {
	last := len(s.stk) - 1
	s.stk[last] = nil
	s.stk = s.stk[:last]
}

// top returns the top screen.
func (s *ScreenStack) top() Screen {
	return s.stk[len(s.stk)-1]
}
