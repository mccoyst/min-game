// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
#pragma once

#include "fixed.hpp"
#include "geom.hpp"
#include <iosfwd>
#include <map>
#include <vector>
#include <memory>
#include <cstdio>

using std::unique_ptr;
using std::vector;

class Img;
class TileView;
class Ui;

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

	class Terrain {
	public:
		Terrain(int i, Fixed s) : tile(i), velScale(s) { }

		// tile is the index number in the tile
		// sheet of the image for this terrain.
		int tile;

		// velScale is the amount by which to scale
		// down a body's base velocity on this
		// terrain.
		Fixed velScale;
	};
	
	// terrain is an array of Terrain indexed by the
	// character representation of the Terrain.
	class TerrainType {
		std::map<char, Terrain> t;
		vector<unique_ptr<Img>> htImg;
	public:
		TerrainType();

		// operator[] returns the terrain with the given character
		// representation.
		Terrain operator [] (char c) const { return t.at(c); }

		// contains returns true iff a terrain type is defined
		// for the given char.
		bool contains(char c) const { return t.find(c) != t.end(); }

		// heightImg returns an image containing the text for
		// the given height value.
		Img &heightImg(int ht) { return *htImg[ht]; }
	} terrain;

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

	// Draw draws the world to the given window.
	void Draw(Ui &, TileView &);

	// At returns the location at the given x,y in the grid.
	//
	// This routine doesn't wrap around at the limits of
	// the world.
	Loc &At(unsigned int x, unsigned int y) {
		return locs.at(x*dim.y.whole()+y);
	}

	const Loc &At(unsigned x, unsigned y) const{
		return locs.at(x*dim.y.whole() + y);
	}

	// AtCoord returns the location at the given world
	// coordinate taking into account wrapping around
	// the ends.
	Loc &AtCoord(int x, int y);

	// AtPoint returns the location at the given point.
	Loc &AtPoint(Vec2);

	// Size returns the size of the world, suitable for
	// using with the geom intersection functions.
	const Vec2 &Size() const {
		return size;
	};

	// Start returns the start location for the player.
	Vec2 Start() const {
		return Vec2(Fixed{x0}*TileW, Fixed{y0}*TileH);
	}

private:

	vector<Loc> locs;
	int width, height;

	// The world's dimensions in numer of tiles.
	Vec2 dim;

	// The world's size, dim*World::TileSz.
	Vec2 size;

	int x0, y0;
};
