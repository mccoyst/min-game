#include <cstdio>
#include <vector>
#include <memory>
#include "ui.hpp"

struct World {

	// TileW and TileH are the size of the world cells
	// in pixels.
	static const Fixed TileW, TileH;

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
		std::shared_ptr<ui::Img> img;
	};
	
	// terrain is an array of Terrain indexed by the
	// character representation of the Terrain.
	struct TerrainType {
		std::vector<Terrain> t;
	public:
		TerrainType();
		Terrain &operator[](int i) { return t[i]; }
	} terrain;

	// A Loc represents a single cell of the world.
	struct Loc {
		unsigned char height, depth;
		char terrain;
	};

	// World constructs a new world by reading it from
	// the given file stream.
	World(FILE*);

	// Draw draws the world to the given window.
	void Draw(std::shared_ptr<ui::Ui>);

	// at returns the location at the given x,y in the grid.
	//
	// This routine doesn't wrap around at the limits of
	// the world.
	Loc &At(unsigned int x, unsigned int y) {
		return locs.at(x*height+y);
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
	std::pair<Fixed,Fixed> Offset() const {
		return std::pair<Fixed,Fixed>(xoff, yoff);
	}

	// Scroll scrolls the world by the given delta;
	void Scroll(Fixed dx, Fixed dy) {
		xoff = (xoff + dx) % (Fixed(width) * TileW);
		yoff = (yoff + dy) % (Fixed(height) * TileH);
	}

private:

	std::vector<Loc> locs;

	int width, height;

	// x and y offset of the viewport.
	Fixed xoff, yoff;
};