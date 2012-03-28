#include "world.hpp"
#include "game.hpp"
#include "ui.hpp"
#include <limits>
#include <climits>
#include <cstdlib>

World::Terrain::Terrain(char c, const char *resrc) : ch(c) {
	char buf[PATH_MAX];
	char *path = realpath(resrc, buf);
	if (!path)
		throw Failure("Failed to get the realpath for %s", resrc);
	img = ui::LoadImg(path);
	if (!img)
		throw Failure("Failed to load %s", path);
}

const Fixed World::TileW(16);
const Fixed World::TileH(16);

World::TerrainType::TerrainType() {
	t.resize(255);
	t['w'] = Terrain('w', "resrc/Water.png");
	t['g'] = Terrain('g', "resrc/Grass.png");
	t['m'] = Terrain('m', "resrc/Mountain.png");
}

World::World(FILE *in) : xoff(0), yoff(0) {
	int n = fscanf(in, "%d %d\n", &width, &height);
	if (n != 2)
		throw Failure("Failed to read width and height");
	if (width <= 0 || height <= 0)
		throw Failure("%d by %d is an invalid world size", width, height);
	if (std::numeric_limits<int>::max() / width < height)
		throw Failure("%d by %d is too big", width, height);

	locs.resize(width*height);
	for (int i = 0; i < width*height; i++) {
		char c;
		int h, d;
		n = fscanf(in, " %c %d %d", &c, &h, &d);
		if (n != 3)
			throw Failure("Failed to read a location %u", i);
		if (h > MaxHeight || h < 0)
			throw Failure("Location %u has invalid height %d", i, h);
		if (d < 0 || d > h)
			throw Failure("Location %u of height %d has invalid depth %d", i, h, d);
		locs[i].height = h;
		locs[i].depth = d;
		locs[i].terrain = &terrain[c];
	}
}

void World::Draw(ui::Ui &ui) {
	Fixed w(ui.width / TileW);
	Fixed h(ui.height / TileH);
	Vec3 offs(xoff%TileW, yoff%TileW);

	for (Fixed x(-1); x <= w; x += Fixed(1)) {
	for (Fixed y(-1); y <= h; y += Fixed(1)) {
		int xcoord = (x - xoff/TileW).whole();
		int ycoord = (y - yoff/TileH).whole();
		const Loc &l = AtCoord(xcoord, ycoord);
		Vec3 v = Vec3(x*TileW, y*TileH, Fixed(0)) + offs;
		ui.Draw(v, l.terrain->img);

		float f = (l.height-l.depth+MaxHeight) / (2.0*MaxHeight);
		ui.Shade(v, Vec3(TileW, TileH), f);
	}
	}
}