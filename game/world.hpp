#include <cstdio>
#include <vector>

struct World {
	
	// A Terrain represents a type of terrain in the world.
	struct Terrain {
		Terrain() : ch(0) { }
		Terrain(char c) : ch(c) { }
	
		char ch;
	};
	
	// terrain is an array of Terrain indexed by the
	// character representation of the Terrain.
	const struct TerrainType {
		std::vector<Terrain> t;
	public:
		TerrainType();
		const Terrain &operator[](int i) const { return t[i]; }
	} terrain;

	// A Loc represents a single cell of the world.
	struct Loc {
		int height, depth;
		const Terrain *terrain;
	};

	// World constructs a new world by reading it from
	// the given file stream.
	World(FILE*);

	// at returns the location at the given x,y in the grid.
	//
	// This routine doesn't wrap around at the limits of
	// the world.
	Loc &At(unsigned int x, unsigned int y) {
		return locs[x*height+y];
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
		return locs[x*height+y];
	}

private:

	std::vector<Loc> locs;
	unsigned int width, height;
};