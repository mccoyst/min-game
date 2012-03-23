#include <cstdio>
#include <vector>
#include <memory>
#include <ui.hpp>

struct World {
	
	// A Terrain represents a type of terrain in the world.
	struct Terrain {
		Terrain() : ch(0), resrc(0), img(0) { }
		Terrain(char c, const char *r) : ch(c), resrc(r), img(0) { }

		// Img returns the image for this terrain.
		//
		// The images are loaded lazily.
		std::shared_ptr<ui::Img> Img(ui::Ui&);
	
		char ch;
		const char *resrc;
	private:
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
		int height, depth;
		Terrain *terrain;
	};

	// World constructs a new world by reading it from
	// the given file stream.
	World(FILE*);

	// Draw draws the world to the given window.
	void Draw(ui::Ui&);

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

private:

	std::vector<Loc> locs;
	unsigned int width, height;
};