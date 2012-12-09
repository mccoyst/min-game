// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package uitil

import (
	"strings"

	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/ui"
)

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

		if wp.X+wsz.X > bounds.Max.X {
			wp.Y += wsz.Y * 1.5
			if wp.Y > bounds.Max.Y {
				return
			}
			wp.X = left
			spword = word
		}

		wsz = d.Draw(spword, wp)
		wp.X += wsz.X
	}
}
