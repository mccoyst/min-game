// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"bufio"
	"github.com/mccoyst/pipeline"
	"io"
	"os"
	"os/exec"

	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/ui"
)

type TitleScreen struct {
	loading bool

	// GenTxt is the last string from the level generator.
	genTxt string

	// wgenErr receives strings from wgen's stderr.
	wgenErr chan string

	// gameChan receieves the *Game from the reader.
	gameChan chan *Game

	frame int
}

func NewTitleScreen() *TitleScreen {
	t := &TitleScreen{}
	if *worldOnStdin {
		t.loadWorld()
	}
	return t
}

func (t *TitleScreen) Transparent() bool {
	return false
}

func (t *TitleScreen) Draw(d ui.Drawer) {
	if t.loading {
		if t.genTxt == "" {
			t.genTxt = "Reticulating Splines"
		}
		d.SetColor(Lime)
		d.SetFont("prstartk", 8)
		genSz := d.TextSize(t.genTxt)
		d.Draw(t.genTxt, geom.Pt(0, ScreenDims.Y-genSz.Y))
		return
	}

	d.SetColor(White)
	d.SetFont("bit_outline", 96)
	titleTxt := "MINIMA"
	titleSz := d.TextSize(titleTxt)
	titlePos := geom.Pt(ScreenDims.X/2-titleSz.X/2,
		ScreenDims.Y/2-titleSz.Y)
	wh := d.Draw(titleTxt, titlePos)

	d.SetFont("prstartk", 16)
	startTxt := "Press " + actionKey() + " to Start"
	startSz := d.TextSize(startTxt)
	startPos := geom.Pt(ScreenDims.X/2-startSz.X/2, titlePos.Y+wh.Y+startSz.Y)
	wh = d.Draw(startTxt, startPos)

	crTxt := "© 2012 The Minima Authors"
	crSz := d.TextSize(crTxt)
	crPos := geom.Pt(ScreenDims.X/2-crSz.X/2, ScreenDims.Y-crSz.Y)
	d.Draw(crTxt, crPos)
}

func (t *TitleScreen) Handle(stk *ui.ScreenStack, e ui.Event) error {
	switch k := e.(type) {
	case ui.Key:
		if k.Down && k.Button == ui.Action {
			t.loadWorld()
		}
	}
	return nil
}

func (t *TitleScreen) Update(stk *ui.ScreenStack) error {
	if !t.loading {
		return nil
	}
	select {
	case s, ok := <-t.wgenErr:
		if !ok {
			t.wgenErr = nil
			break
		}
		t.genTxt = s
	case g := <-t.gameChan:
		for _ = range t.genTxt {
		} // junk it
		t.loading = false
		stk.Push(g)
	default:
	}
	return nil
}

func (t *TitleScreen) loadWorld() {
	t.gameChan = make(chan *Game)
	t.wgenErr = make(chan string, 1)
	t.loading = true

	go func() {
		if *worldOnStdin {
			*worldOnStdin = false
			t.wgenErr <- "Reading World"
			close(t.wgenErr)
			in := bufio.NewReader(os.Stdin)
			g, err := ReadGame(in)
			if err != nil {
				panic(err)
			}
			t.gameChan <- g
			return
		}

		wgen := exec.Command("wgen")
		stderr, err := wgen.StderrPipe()
		if err != nil {
			panic(err)
		}
		go readErr(stderr, t.wgenErr)

		p, err := pipeline.New(
			wgen,
			exec.Command("herbnear", "-num", "25", "-name", "Cow"),
			exec.Command("herbnear", "-num", "25", "-name", "Gull"),
			exec.Command("tee", "cur"),
		)
		if err != nil {
			panic(err)
		}

		stdout, err := p.Last().StdoutPipe()
		if err != nil {
			panic(err)
		}
		if err := p.Start(); err != nil {
			panic(err)
		}
		in := bufio.NewReader(stdout)
		g, err := ReadGame(in)
		if err != nil {
			panic(err)
		}
		if errs := p.Wait(); len(errs) > 0 {
			panic(errs)
		}
		t.gameChan <- g
	}()
}

// ReadErr reads wgen's standard error, picks out
// what it is currently doing and sends the strings
// to the channel.
func readErr(r io.Reader, strs chan<- string) {
	in := bufio.NewReader(io.TeeReader(r, os.Stderr))
	for {
		_, err := readRunes(in, '\n')
		if err != nil {
			break
		}
		line, err := readRunes(in, '…')
		if err != nil {
			break
		}
		if line == "Writing the world" {
			line = "Reading the world"
		}
		strs <- line
	}
	close(strs)
}

// readRunes returns all runes until the delimiter
// is read or an error occurs.  The delimiter is not
// included in the returned slice.
func readRunes(in *bufio.Reader, delim rune) (string, error) {
	var err error
	var runes []rune
	for {
		var r rune
		r, _, err = in.ReadRune()
		if err != nil || r == delim {
			break
		}
		runes = append(runes, r)
	}
	return string(runes), err
}

func actionKey() string {
	for k, b := range ui.CurrentKeymap {
		if b == ui.Action {
			return k.String()
		}
	}
	panic("Somebody broke the keymap")
}
