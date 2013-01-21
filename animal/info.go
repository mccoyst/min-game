// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package animal

import (
	"encoding/json"
	"os"

	"code.google.com/p/min-game/ai"
	"code.google.com/p/min-game/resrc"
	"code.google.com/p/min-game/sprite"
)

// Info contains the properties for an entire species of animal.
type Info struct {
	Name     string
	Sheet    sprite.Sheet
	Affinity map[string]float64
	BoidInfo ai.BoidInfo
}

var finder = resrc.NewPkgFinder()

func LoadInfo(s string) (Info, error) {
	var i Info

	f, err := os.Open(finder.Find(s + ".info"))
	if err != nil {
		return i, err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	err = dec.Decode(&i)
	return i, err
}
