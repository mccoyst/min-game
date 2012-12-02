// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package animal

import (
	"code.google.com/p/min-game/phys"
	"code.google.com/p/min-game/sprite"
	"code.google.com/p/min-game/ui"
	"code.google.com/p/min-game/world"
)

// Herbivore is a struct representing basic herding, non-agressive, herbivorous animals.
type Herbivore struct {
	Body phys.Body
	Anim sprite.Anim
	info *Info
}

func (h *Herbivore) Move(w *world.World) {
	h.Anim.Move(&h.info.Sheet, h.Body.Vel)
	h.Body.Move(w, h.info.Affinity)
}

func (h *Herbivore) Draw(d ui.Drawer, cam ui.Camera) error {
	_, err := cam.Draw(d, ui.Sprite{
		Name:   h.info.Sheet.Name,
		Bounds: h.info.Sheet.Frame(h.Anim.Face, h.Anim.Frame),
		Shade:  1.0,
	}, h.Body.Box.Min)
	return err
}