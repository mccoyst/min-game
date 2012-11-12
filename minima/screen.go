// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"image/color"
	"time"

	"code.google.com/p/min-game/ui"
)

// A Drawer can draw things and change colors.
type Drawer interface {
	Draw(interface{}, ui.Point) (ui.Point, error)
	SetFont(name string, szPts float64) error
	SetColor(color.Color)
	TextSize(string) ui.Point
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
}

// A ScreenStack holds the stack of game screens.
type ScreenStack struct {
	stk       []Screen
	win       *ui.Ui
	nFrames   uint
	meanFrame float64 // milliseconds
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
			if _, ok := e.(ui.Quit); ok {
				return
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
