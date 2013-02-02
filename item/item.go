// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package item

import (
	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/world"
	"fmt"
)

const (
	ETele    = "E-Tele"
	Uranium  = "Uranium"
	Flippers = "Flippers"
	Scrap    = "Scrap"
)

// InitialUses returns the initial value for the Uses field
// of items, indexed by their name.
var initialUses = map[string]int{
	ETele: 3,
}

// Descs maps item names to functions that return
// an item description.
var descs = map[string]func(*Item) string{
	ETele: withUses(`The Emergency Teleporter reacts to
		critical condition by sending you back to home base.`),

	Uranium: func(*Item) string {
		return `Uranium is of great interest because of its application to
		nuclear power and nuclear weapons. Uranium contamination is
		an emotive environmental problem. It is not particularly rare and
		is more common than beryllium or tungsten for instance.
		[br] [br] http://www.webelements.com/uranium`
	},

	Flippers: func(*Item) string {
		return "A flat rubber attachment worn on the foot for underwater swimming." // via Googling for "define flippers."  Is this copyrighted or something?
	},
}

// WithUses returns a function that prints the description
// followed by the remaining number of uses.
func withUses(desc string) func(*Item) string {
	return func(i *Item) string {
		return fmt.Sprintf("%s  This %s currently has %d uses remaining.",
			desc, i.Name, i.Uses)
	}
}

// Bonus maps each item to it's terrain bonus.  A terrain bonus is
// the scale to set for the given terrain type.  If an item is not in this
// map then it doesn't give a terrain bonus.
var Bonus = map[string]struct {
	Terrain string
	Scale   float64
}{
	Flippers: {"w", 1.0},
}

// An Item is something that the player can collect and possibly use.
type Item struct {
	Name string
	Uses int
}

// New returns a new item of the given name.
func New(name string) *Item {
	return &Item{name, initialUses[name]}
}

// Desc returns the item's description.
func (it *Item) Desc() string {
	if d, ok := descs[it.Name]; ok {
		return d(it)
	}
	return "<No Description for " + it.Name + ">"
}

// A Treasure is an Item located somewhere in the world.
type Treasure struct {
	Item *Item
	Box  geom.Rectangle
}

// NewTreasure returns a new treasure.
func NewTreasure(x, y float64, it *Item) *Treasure {
	return &Treasure{
		Item: it,
		Box:  geom.Rect(x, y, x+world.TileSize.X, y+world.TileSize.Y),
	}
}
