// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
#pragma once

#include "fixed.hpp"
#include "geom.hpp"

using std::unique_ptr;

class Img;
class TileView;
class Ui;

class Terrain {
public:
	Terrain(int i, Fixed s) : tile(i), velScale(s) { }

	// tile is the index number in the tile
	// sheet of the image for this terrain.
	int tile;

	// velScale is the amount by which to scale
	// a body's base velocity on this terrain.
	Fixed velScale;
};

class World {
public:
	// TileW and TileH are the size of the world cells
	// in pixels.
	static const Fixed TileW, TileH;

	// A Vec2 with the width and height of a tile.
	static const Vec2 TileSz;

	enum {
		// MaxHeight is the maximum value for the height
		// of any location.  Heights are 0..MaxHeight.
		MaxHeight = 19,
	};

	// A Loc represents a single cell of the world.
	struct Loc {
		int x, y;
		unsigned char height, depth;
		char terrain;

		Bbox Box() const;
		float Shade() const;
	};

	// World constructs a new world by reading it from
	// the given file stream.
	World(std::istream&);
	~World();

	// Draw draws the world to the given window.
	void Draw(Ui &, TileView &);

	// At returns the location at the given x,y in the grid.
	//
	// This routine doesn't wrap around at the limits of
	// the world.
	Loc &At(unsigned int x, unsigned int y);

	const Loc &At(unsigned x, unsigned y) const;

	// AtCoord returns the location at the given world
	// coordinate taking into account wrapping around
	// the ends.
	const Loc &AtCoord(int x, int y) const;

	// AtPoint returns the location at the given point.
	const Loc &AtPoint(Vec2) const;

	Terrain TerrainAt(Vec2) const;

	// Size returns the size of the world, suitable for
	// using with the geom intersection functions.
	Vec2 Size() const;

	// Start returns the start location for the player.
	Vec2 Start() const;

private:
	class Impl;
	unique_ptr<Impl> impl;
};
