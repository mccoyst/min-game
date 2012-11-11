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

	// GenTxt is the last string from the level generator.
	genTxt string

	// wgenErr receives strings from wgen's stderr.
	wgenErr chan string

	// wChan receieves the world from the reader.
	wChan chan interface{}

	frame int
}

func NewTitleScreen() *TitleScreen {
	t := &TitleScreen{}
	if *worldOnStdin {
		t.loadWorld()
	}
	return t
}

func (t *TitleScreen) Draw(d Drawer) error {
	if t.loading {
		d.SetColor(Lime)
		if err := d.SetFont("prstartk", 12); err != nil {
			return err
		}
		if t.genTxt == "" {
			t.genTxt = "Generating World"
		}
		genSz := d.TextSize(t.genTxt)
		_, err := d.Draw(t.genTxt, ui.Pt(0, ScreenDims.Y-genSz.Y))
		return err
	}

	d.SetColor(White)
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
			t.loadWorld()
		}
	}
	return nil
}

func (t *TitleScreen) Update(stk *ScreenStack) error {
	if t.loading {
		select {
		case s, ok := <-t.wgenErr:
			if !ok {
				t.wgenErr = nil
				break
			}
			t.genTxt = s
		case i := <-t.wChan:
			switch w := i.(type) {
			case error:
				return w
			case *world.World:
				for _ = range t.genTxt { // junk it
				}
				*worldOnStdin = false
				stk.Push(NewExploreScreen(w))
				t.loading = false
			}
		default:
		}
	}
	return nil
}

func (t *TitleScreen) loadWorld() {
	t.wChan = make(chan interface{})
	t.wgenErr = make(chan string, 1)
	t.loading = true

	go func() {
		if *worldOnStdin {
			t.wgenErr <- "Reading World"
			close(t.wgenErr)
			if w, err := world.Read(bufio.NewReader(os.Stdin)); err != nil {
				t.wChan <- err
			} else {
				t.wChan <- w
			}
			return
		}

		cmd := exec.Command("wgen")
		stdout, stderr, err := pipes(cmd)
		if err != nil {
			t.wChan <- err
			return
		}
		go readErr(stderr, t.wgenErr)

		cmd.Start()
		w, err := world.Read(bufio.NewReader(stdout))
		if err != nil {
			t.wChan <- err
			return
		}
		if err = cmd.Wait(); err != nil {
			t.wChan <- err
			return
		}
		t.wChan <- w
	}()
}

// Pipes returns the standard output and error
// pipes for a command.
func pipes(cmd *exec.Cmd) (stdout, stderr io.Reader, err error) {
	if stdout, err = cmd.StdoutPipe(); err != nil {
		return
	}
	stderr, err = cmd.StderrPipe()
	return
}

// ReadErr reads wgen's standard error, picks out
// what it is currently doing and sends the strings
// to the channel.
// BUG(eaburns): readErr is pretty ugly.
func readErr(r io.Reader, strs chan<- string) {
	var err error
	var runes []rune
	in := bufio.NewReader(io.TeeReader(r, os.Stderr))
	for {
		var r rune
		if r, _, err = in.ReadRune(); err != nil {
			break
		}
		if len(runes) > 0 && runes[0] == '…' {
			if r == '\n' {
				runes = runes[:0]
			}
			continue
		}
		if r == '…' || r == '\n' {
			s := string(runes)
			if s == "Writing the world" {
				strs <- "Reading the world"
				break
			}
			strs <- s
			runes = runes[:0]
		}
		if r != '\n' {
			runes = append(runes, r)
		}
	}

	close(strs)
	for err != nil {
		_, err = in.ReadByte()
	}
}

func actionKey() string {
	for k, b := range ui.DefaultKeymap {
		if b == ui.Action {
			return k.String()
		}
	}
	panic("Somebody broke the keymap")
}
