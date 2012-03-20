package world

import (
	"image/color"
)

// TerrainType holds information on a given type of
// terrain.
type TerrainType struct {
	Char  uint8      // character representation
	Color color.RGBA // color at the highest altitude
}

// Terrain is an array with the canonical terrain
// representations, indexed by the terrain type's
// unique character.
var Terrain = []TerrainType{
	// Grass-land
	int('g'): {'g', color.RGBA{0, 255, 0, 255}},
	// Mountain
	int('m'): {'m', color.RGBA{196, 196, 196, 255}},
	// Water
	int('w'): {'w', color.RGBA{0, 0, 255, 255}},
	// Lava
	int('l'): {'l', color.RGBA{255, 0, 0, 255}},
	// Desert
	int('d'): {'d', color.RGBA{255, 255, 0, 255}},
}
