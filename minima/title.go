// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/mccoyst/min-game/geom"
	"github.com/mccoyst/min-game/ui"
	"mccoy.space/g/pipeline"
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
		d.SetFont(DialogFont, 8)
		genSz := d.TextSize(t.genTxt)
		d.Draw(t.genTxt, geom.Pt(0, ScreenDims.Y-genSz.Y))
		return
	}

	d.SetColor(White)
	d.SetFont(TitleFont, 96)
	titleTxt := "MINIMA"
	titleSz := d.TextSize(titleTxt)
	titlePos := geom.Pt(ScreenDims.X/2-titleSz.X/2,
		ScreenDims.Y/2-titleSz.Y)
	wh := d.Draw(titleTxt, titlePos)

	d.SetFont(DialogFont, 16)
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
			t.wgenErr <- "Reading the world"
			close(t.wgenErr)
			in := bufio.NewReader(os.Stdin)
			g, err := ReadGame(in)
			if err != nil {
				panic(err)
			}
			t.gameChan <- g
			return
		}
		cmds := []*exec.Cmd{
			gen("wgen"),
			gen("herbgen 25 Gull 10 Guppy 10 Guppy 10 Guppy 10 Guppy 10 Guppy 10 Guppy 10 Guppy 10 Guppy 10 Guppy 10 Guppy 25 Cow 25 Cow 25 Cow 25 Cow 10 Chicken 10 Chicken 10 Chicken 10 Chicken 10 Chicken"),
			gen("itemnear -num 2 -name Uranium"),
			gen("itemnear -num 2 -name Scrap"),
			gen("itemnear -num 1 -name Flippers"),
		}

		stderrin, stderrout, err := os.Pipe()
		if err != nil {
			panic(err)
		}
		defer stderrin.Close()
		defer stderrout.Close()
		for _, c := range cmds {
			c.Stderr = stderrout
		}
		go readErr(stderrin, t.wgenErr)

		if *debug {
			cmds = append(cmds, exec.Command("wimg", "-e", "-o", "cur.png"))
			cmds = append(cmds, exec.Command("tee", "cur.world"))
		}

		p, err := pipeline.New(cmds...)
		if err != nil {
			panic(err)
		}
		if *debug {
			os.Stderr.Write([]byte("executing: " + p.String() + "\n"))
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
			s := ""
			for _, e := range errs {
				s += e.Error() + "\n"
			}
			panic(s)
		}
		t.gameChan <- g
	}()
}

// Gen returns a command for an XXXgen program.  The command
// is given by a string, followed by "-seed" and the random seed.
func gen(s string) *exec.Cmd {
	fs := strings.Fields(s)
	cmd := fs[0]
	fs = append([]string{"-seed", strconv.FormatInt(*seed, 10)}, fs[1:]...)
	*seed++
	return exec.Command(cmd, fs...)
}

// ReadErr reads wgen's standard error, picks out
// what it is currently doing and sends the strings
// to the channel.
func readErr(r io.Reader, strs chan<- string) {
	in := bufio.NewReader(io.TeeReader(r, os.Stderr))
	for {
		s, err := readRunes(in, "…\n")
		if err != nil {
			break
		}
		r, sz := utf8.DecodeLastRuneInString(s)
		if r == '\n' {
			continue
		}
		strs <- s[:len(s)-sz]
	}
	strs <- "Reading the world"
	close(strs)
}

// readRunes returns all runes until a delimiter
// is read or an error occurs.  If a delimiter is
// read then it is included in the returned string.
func readRunes(in *bufio.Reader, delims string) (string, error) {
	var err error
	var runes []rune
	for {
		var r rune
		r, _, err = in.ReadRune()
		if err != nil {
			break
		}
		runes = append(runes, r)
		if strings.ContainsRune(delims, r) {
			break
		}
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
