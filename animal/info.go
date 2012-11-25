// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package animal

import (
	"encoding/json"
	"os"

	"code.google.com/p/min-game/ai"
	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/phys"
	"code.google.com/p/min-game/sprite"
)

// Info contains the properties for an entire species of animal.
type Info struct {
	Name     string
	Sheet    sprite.Sheet
	Affinity map[string]float64
	BoidInfo ai.BoidInfo
}

func LoadInfo(s string) (Info, error) {
	var i Info

	f, err := os.Open("resrc/" + s + ".info")
	if err != nil {
		return i, err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	err = dec.Decode(&i)
	return i, err
}

func (i *Info) SpawnHerbivore(p, v geom.Point) Herbivore {
	sz := float64(i.Sheet.FrameSize)
	return Herbivore{
		Body: phys.Body{
			Box: geom.Rect(p.X, p.Y, p.X+sz, p.Y+sz),
			Vel: v,
		},
		info: i,
	}
}
