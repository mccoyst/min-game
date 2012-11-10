// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"bufio"
	"os"
	"os/exec"

	"code.google.com/p/min-game/ui"
	"code.google.com/p/min-game/world"
)

type TitleScreen struct {
	loading bool
	frame   int
}

func NewTitleScreen() *TitleScreen {
	return &TitleScreen{}
}

func (t *TitleScreen) Draw(d Drawer) error {
	if t.loading {
		d.SetColor(109, 170, 44, 255)
		_, err := d.Draw(ui.Text{
			Font:   "prstartk",
			Pts:    12,
			string: "Generating World",
		}, ui.Pt(0, ScreenDims.Y-12))
		return err
	}

	// BUG(mccoyst): need to expose Font.Width() in Drawer
	titlePos := ui.Pt(ScreenDims.X/2-float64(64*len("MINIMA"))/2, ScreenDims.Y/2-64)

	d.SetColor(255, 255, 255, 255)

	wh, err := d.Draw(ui.Text{
		Font:   "prstartk",
		Pts:    64,
		string: "MINIMA",
	}, titlePos)

	startPos := ui.Pt(titlePos.X, titlePos.Y+wh.Y+12)
	wh, err = d.Draw(ui.Text{
		Font:   "prstartk",
		Pts:    12,
		string: "Press " + actionKey() + " to Start",
	}, startPos)

	cr := "© 2012 The Minima Authors"
	crPos := ui.Pt(ScreenDims.X/2-float64(12*len(cr))/2, ScreenDims.Y-12)
	_, err = d.Draw(ui.Text{
		Font:   "prstartk",
		Pts:    12,
		string: cr,
	}, crPos)

	return err
}

func (t *TitleScreen) Handle(stk *ScreenStack, e ui.Event) error {
	switch k := e.(type) {
	case ui.Key:
		if k.Down && ui.DefaultKeymap[k.Code] == ui.Action {
			t.loading = true
		}
	}
	return nil
}

func (t *TitleScreen) Update(stk *ScreenStack) error {
	if t.loading {
		var w world.World
		var err error
		if *worldOnStdin {
			w, err = world.Read(bufio.NewReader(os.Stdin))
			if err != nil {
				return err
			}
		} else {
			cmd := exec.Command("wgen")
			stdout, err := cmd.StdoutPipe()
			if err != nil {
				return err
			}
			if err = cmd.Start(); err != nil {
				return err
			}

			w, err = world.Read(bufio.NewReader(stdout))
			if err != nil {
				return err
			}

			if err = cmd.Wait(); err != nil {
				return err
			}
		}

		*worldOnStdin = false
		stk.Push(NewExploreScreen(w))
		t.loading = false
	}
	return nil
}

func actionKey() string {
	for k, b := range ui.DefaultKeymap {
		if b == ui.Action {
			return k.String()
		}
	}
	panic("Somebody broke the keymap")
}
