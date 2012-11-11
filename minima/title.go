// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"bufio"
	"fmt"
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

	// strChan receives strings from wgen's stderr.
	strChan <-chan string

	// wChan receieves the world from the reader.
	wChan <-chan interface{}

	frame int
}

func NewTitleScreen() *TitleScreen {
	t := &TitleScreen{}
	if *worldOnStdin {
		t.startLoading()
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
			t.startLoading()
		}
	}
	return nil
}

func (t *TitleScreen) Update(stk *ScreenStack) error {
	if t.loading {
		select {
		case s, ok := <-t.strChan:
			if !ok {
				t.strChan = nil
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

func (t *TitleScreen) startLoading() {
	t.wChan, t.strChan = readWorld()
	t.loading = true
}

func readWorld() (<-chan interface{}, <-chan string) {
	ch := make(chan interface{})
	strChan := make(chan string, 1)
	go func() {
		var r io.Reader = os.Stdin
		var err error
		var cmd *exec.Cmd

		if !*worldOnStdin {
			cmd = exec.Command("wgen")
			if r, err = cmd.StdoutPipe(); err != nil {
				ch <- err
				return
			}
			stderr, err := cmd.StderrPipe()
			if err != nil {
				ch <- err
				return
			}
			go readErr(stderr, strChan)
			cmd.Start()
		} else {
			fmt.Println("Sending Reading World")
			strChan <- "Reading World"
			close(strChan)
		}
		w, err := world.Read(bufio.NewReader(r))
		if err != nil {
			ch <- err
			return
		}

		if !*worldOnStdin {
			if err = cmd.Wait(); err != nil {
				ch <- err
				return
			}
		}
		ch <- w
	}()
	return ch, strChan
}

func readErr(r io.Reader, strs chan<- string) {
	var err error
	var runes []rune
	in := bufio.NewReader(r)
	for {
		var r rune
		if r, _, err = in.ReadRune(); err != nil {
			break
		}
		os.Stderr.WriteString(string([]rune{r}))
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

	// Junk the rest
	for err == nil {
		var b byte
		b, err = in.ReadByte()
		os.Stderr.Write([]byte{b})
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
