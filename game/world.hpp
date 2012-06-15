// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
#pragma once

#include "fixed.hpp"
#include <iosfwd>
#include <map>
#include <vector>
#include <memory>

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
	
	// terrain is an array of Terrain indexed by the
	// character representation of the Terrain.
	class TerrainType {
		std::map<char, int> t;
		vector<unique_ptr<Img>> htImg;
	public:
		TerrainType();

		// operator[] returns the terrain with the given character
		// representation.
		int operator [] (char c) const { return t.at(c); }

		// contains returns true iff a terrain type is defined for the given char.
		bool contains(char c) const { return t.find(c) != t.end(); }

		// heightImg returns an image containing the text for
		// the given height value.
		Img &heightImg(int ht) { return *htImg[ht]; }
	} terrain;

	// A Loc represents a single cell of the world.
	struct Loc {
		unsigned char height, depth;
		char terrain;

		float Shade() const;
	};

	// World constructs a new world by reading it from
	// the given file stream.
	World(std::istream&);

	// Draw draws the world to the given window.
	void Draw(Ui &, TileView &);

	// at returns the location at the given x,y in the grid.
	//
	// This routine doesn't wrap around at the limits of
	// the world.
	Loc &At(unsigned int x, unsigned int y) {
		return locs.at(x*size.y.whole()+y);
	}

	const Loc &At(unsigned x, unsigned y) const{
		return locs.at(x*size.y.whole() + y);
	}

	// atcoord returns the location at the given world
	// coordinate taking into account wrapping around
	// the ends.
	Loc &AtCoord(int x, int y) {
		if (x < 0)
			x = width - -x%width;
		else
			x %= width;
		if (y < 0)
			y = height - -y%height;
		else
			y %= height;
		return At(x, y);
	}

	// The indices of the start tile.
	int x0, y0;

	// The world's dimensions.
	Vec2 size;

private:

	vector<Loc> locs;
	int width, height;
};
