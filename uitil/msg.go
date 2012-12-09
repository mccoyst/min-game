// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package uitil

import (
	"strings"

	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/ui"
)

func WordWrap(d ui.Drawer, text string, bounds geom.Rectangle) {
	words := strings.Fields(text)

	left := bounds.Min.X
	wp := bounds.Min
	for _, word := range words {
		word += " "
		wsz := d.TextSize(word)
		if wp.X+wsz.X > bounds.Dx() {
			wp.Y += wsz.Y * 1.5
			wp.X = left
		}

		d.Draw(word, wp)
		wp.X += wsz.X
	}
}
