#include "world.hpp"
#include "game.hpp"
#include <limits>

World::TerrainType::TerrainType() {
	t.resize(255);
	t['w'] = Terrain('w');
	t['g'] = Terrain('g');
	t['m'] = Terrain('m');
}

World::World(FILE *in) {
	int n = fscanf(in, "%u %u\n", &width, &height);
	if (n != 2)
		throw Failure("Failed to read width and height");
	if (std::numeric_limits<unsigned int>::max() / width < height)
		throw Failure("%u by %u is too big", width, height);

	locs.resize(width*height);
	for (unsigned int i = 0; i < width*height; i++) {
		char c;
		int h, d;
		n = fscanf(in, " %c %d %d", &c, &h, &d);
		if (n != 3)
			throw Failure("Failed to read a location %u", i);
		locs[i].height = h;
		locs[i].depth = d;
		locs[i].terrain = &terrain[c];
	}
}
