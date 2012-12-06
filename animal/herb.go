// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package animal

import (
	"code.google.com/p/min-game/phys"
	"code.google.com/p/min-game/sprite"
	"code.google.com/p/min-game/ui"
	"code.google.com/p/min-game/ai"
	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/world"
)

// Herbivore is a struct representing basic herding, non-agressive, herbivorous animals.
type Herbivore struct {
	Body phys.Body
	Anim sprite.Anim
}

type Herbivores struct {
	info *Info
	herbs []*Herbivore
}

func MakeHerbivores(name string) (Herbivores, error) {
	i, err := LoadInfo(name)
	if err != nil {
		return Herbivores{}, err
	}
	return Herbivores{ &i, nil }, err
}

func (hs Herbivores) Move(w *world.World) {
	for _, h := range hs.herbs {
		h.Anim.Move(&hs.info.Sheet, h.Body.Vel)
		h.Body.Move(w, hs.info.Affinity)
	}
}

func (hs Herbivores) Draw(d ui.Drawer, cam ui.Camera) error {
	for _, h := range hs.herbs {	
		_, err := cam.Draw(d, ui.Sprite{
			Name:   hs.info.Sheet.Name,
			Bounds: hs.info.Sheet.Frame(h.Anim.Face, h.Anim.Frame),
			Shade:  1.0,
		}, h.Body.Box.Min)
		if err != nil {
			return err
		}
	}
	return nil
}

// Spawn spawns a new Herbivore for this Herbivores collection.
func (hs *Herbivores) Spawn(p, v geom.Point) {
	sz := float64(hs.info.Sheet.FrameSize)
	hs.herbs = append(hs.herbs, &Herbivore{
		Body: phys.Body{
			Box: geom.Rect(p.X, p.Y, p.X+sz, p.Y+sz),
			Vel: v,
		},
	})
}

func (hs Herbivores) Len() int {
	return len(hs.herbs)
}

func (hs Herbivores) Boid(n int) ai.Boid {
	return ai.Boid{&hs.herbs[n].Body}
}

func (hs Herbivores) BoidInfo() ai.BoidInfo {
	return hs.info.BoidInfo
}
