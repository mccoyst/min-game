// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package animal

import (
	"code.google.com/p/min-game/ai"
	"code.google.com/p/min-game/phys"
	"code.google.com/p/min-game/ui"
	"code.google.com/p/min-game/world"
	"encoding/json"
	"io"
	"io/ioutil"
)

// Animals is a struct gathering all of the animals so
// that they can be easily marshaled and demarshaled.
type Animals struct {
	Gulls Herbivores
	Cows  Herbivores
}

// WriteTo writes the animals to the given writer.
func (a Animals) WriteTo(w io.Writer) (int64, error) {
	b, err := json.Marshal(a)
	if err != nil {
		return 0, err
	}
	n, err := w.Write(b)
	return int64(n), err
}

// Read returns Animals read from the given reader.
func Read(r io.Reader) (Animals, error) {
	var anims Animals
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return anims, err
	}
	err = json.Unmarshal(b, &anims)
	return anims, err
}

func (a Animals) Draw(d ui.Drawer, cam ui.Camera) error {
	if err := a.Gulls.Draw(d, cam); err != nil {
		return err
	}
	if err := a.Cows.Draw(d, cam); err != nil {
		return err
	}
	return nil
}

func (a Animals) Move(astro *phys.Body, w *world.World) {
	ai.UpdateBoids(a.Gulls, astro, w)
	a.Gulls.Move(w)
	ai.UpdateBoids(a.Cows, astro, w)
	a.Cows.Move(w)
}
