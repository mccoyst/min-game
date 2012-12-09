// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package uitil

import (
	"image/color"
	"strings"

	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/ui"
)

// WordWrap draws text to the screen, bounded by the given rectangle. The special token "[br]" inserts a line break.
// Any text that does not fit beyond bounds.Max.Y is truncated.
func WordWrap(d ui.Drawer, text string, bounds geom.Rectangle) {
	words := strings.Fields(text)
	if len(words) == 0 {
		return
	}

	left := bounds.Min.X
	wp := bounds.Min

	wsz := d.Draw(words[0], wp)
	wp.X += wsz.X

	for _, word := range words[1:] {
		spword := " " + word
		wsz = d.TextSize(spword)

		if wp.X+wsz.X > bounds.Max.X || word == "[br]" {
			wp.Y += wsz.Y * 1.5
			if wp.Y > bounds.Max.Y {
				return
			}
			wp.X = left
			if word == "[br]" {
				continue
			}
		}

		if wp.X == left {
			spword = word
		}

		wsz = d.Draw(spword, wp)
		wp.X += wsz.X
	}
}

type MessageBox struct {
	Text    string //TODO(mccoyst): Implement paging
	Font    string
	Fontsz  float64
	Fg, Bg  color.Color
	Box     geom.Rectangle
	Pad     float64
	closing bool
}

func (mb *MessageBox) Transparent() bool {
	return true
}

func (mb *MessageBox) Draw(d ui.Drawer) {
	border := mb.Box.Pad(mb.Pad)
	d.SetColor(mb.Fg)
	d.Draw(border, geom.Pt(0, 0))

	d.SetFont(mb.Font, mb.Fontsz)
	d.SetColor(mb.Bg)
	d.Draw(mb.Box, geom.Pt(0, 0))

	d.SetColor(mb.Fg)
	WordWrap(d, mb.Text, mb.Box.Rpad(mb.Pad))
}

func (mb *MessageBox) Handle(stk *ui.ScreenStack, e ui.Event) error {
	if mb.closing {
		return nil
	}

	key, ok := e.(ui.Key)
	if !ok || !key.Down {
		return nil
	}

	mb.closing = true
	return nil
}

func (mb *MessageBox) Update(stk *ui.ScreenStack) error {
	if mb.closing {
		stk.Pop()
		return nil
	}

	return nil
}
