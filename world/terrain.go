// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
package world

import (
	"image/color"
)

// TerrainType holds information on a given type of
// terrain.
type TerrainType struct {
	Char  uint8      // character representation
	Name string	// 
	Color color.RGBA // color at the highest altitude
}

// Terrain is an array with the canonical terrain
// representations, indexed by the terrain type's
// unique character.
var Terrain = []TerrainType{
	// Grass-land
	int('g'): {'g', "grass", color.RGBA{0, 255, 0, 255}},
	// Mountain
	int('m'): {'m', "mountain", color.RGBA{196, 196, 196, 255}},
	// Water
	int('w'): {'w', "water", color.RGBA{0, 0, 255, 255}},
	// Lava
	int('l'): {'l', "lava", color.RGBA{255, 0, 0, 255}},
	// Desert
	int('d'): {'d', "desert", color.RGBA{255, 255, 0, 255}},
	// Forrest
	int('f'): {'f', "forest", color.RGBA{0, 200, 128, 255}},
}
