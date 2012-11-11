// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"bufio"
	"io"
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
		if err := d.SetFont("prstartk", 12); err != nil {
			return err
		}
		genTxt := "Generating World"
		genSz := d.TextSize(genTxt)
		_, err := d.Draw(genTxt, ui.Pt(0, ScreenDims.Y-genSz.Y))
		return err
	}

	d.SetColor(255, 255, 255, 255)
	if err := d.SetFont("prstartk", 64); err != nil {
		return err
	}
	titleTxt := "MINIMA"
	titleSz := d.TextSize(titleTxt)
	titlePos := ui.Pt(ScreenDims.X/2-titleSz.X/2,
		ScreenDims.Y/2-titleSz.Y)
	wh, err := d.Draw(titleTxt, titlePos)
	if err != nil {
		return err
	}

	if err := d.SetFont("prstartk", 12); err != nil {
		return err
	}
	startTxt := "Press " + actionKey() + " to Start"
	startSz := d.TextSize(startTxt)
	startPos := ui.Pt(titlePos.X, titlePos.Y+wh.Y+startSz.Y)
	if wh, err = d.Draw(startTxt, startPos); err != nil {
		return err
	}

	crTxt := "© 2012 The Minima Authors"
	crSz := d.TextSize(crTxt)
	crPos := ui.Pt(ScreenDims.X/2-crSz.X/2, ScreenDims.Y-crSz.Y)
	_, err = d.Draw(crTxt, crPos)
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
			if stderr, err := cmd.StderrPipe(); err != nil {
				return err
			} else {
				go io.Copy(os.Stderr, stderr)
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
