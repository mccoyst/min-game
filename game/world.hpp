#include <cstdio>
#pragma once

#include <vector>
#include <memory>
#include "ui.hpp"

struct World {

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
	
	// A Terrain represents a type of terrain in the world.
	struct Terrain {
		Terrain() : ch(0), img(0) { }

		Terrain(char, const char*);
	
		char ch;
		std::shared_ptr<Img> img;
	};
	
	// terrain is an array of Terrain indexed by the
	// character representation of the Terrain.
	struct TerrainType {
		std::vector<Terrain> t;
		std::vector< std::shared_ptr<Img> > htImg;
	public:
		TerrainType();

		// operator[] returns the terrain with the given character
		// representation.
		Terrain &operator[](int i) { return t[i]; }

		// heightImg returns an image containing the text for
		// the given height value.
		std::shared_ptr<Img> heightImg(int ht) { return htImg[ht]; }
	} terrain;

	// A Loc represents a single cell of the world.
	struct Loc {
		unsigned char height, depth;
		char terrain;

		float Shade() const;
	};

	// World constructs a new world by reading it from
	// the given file stream.
	World(FILE*);

	// Draw draws the world to the given window.
	void Draw(std::shared_ptr<Ui>);

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
		x %= width;
		if (x < 0)
			x = width + x;
		y %= height;
		if (y < 0)
			y = height + y;
		return At(x, y);
	}

	// Offset returns the current world offset.
	Vec2 Offset() const {
		return Vec2(xoff, yoff);
	}

	// Scroll scrolls the world by the given delta;
	void Scroll(Fixed dx, Fixed dy) {
		xoff = (xoff + dx) % (Fixed(width) * TileW);
		yoff = (yoff + dy) % (Fixed(height) * TileH);
	}

	// Center changes the world's scroll offset so that the location
	// at the given x,y coordinate is centered.
	void Center(std::shared_ptr<Ui> win, int x, int y) {
		x %= width;
		if (x < 0)
			x = width + x;
		y %= height;
		if (y < 0)
			y = height + y;
		xoff = win->width/Fixed(2) - (Fixed(x) * TileW);
		yoff = win->height/Fixed(2) - (Fixed(y) * TileH);
	}

	// The indices of the start tile.
	int x0, y0;

	// The world's dimensions.
	Vec2 size;

private:

	std::vector<Loc> locs;

	int width, height;

	// x and y offset of the viewport.
	Fixed xoff, yoff;
};