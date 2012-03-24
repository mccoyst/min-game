#include "world.hpp"
#include "game.hpp"
#include <ui.hpp>
#include <limits>

std::shared_ptr<ui::Img> World::Terrain::Img(ui::Ui &ui) {
	if (!img.get())
		img = ui.LoadImg(resrc);
	return img;
}

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
		locs[i].height = h;
		locs[i].depth = d;
		locs[i].terrain = &terrain[c];
	}
}

void World::Draw(ui::Ui &ui) {
	int w = ui.width.whole() / TileW;
	int h = ui.height.whole() / TileH;

	for (int x = -1; x <= w; x++) {
	for (int y = -1; y <= h; y++) {
		const Loc &l = AtCoord(
			x - xoff/TileW,
			y - yoff/TileH
		);

		ui::Vec3 v(
			ui::Len(x*TileW + xoff%TileW),
			ui::Len(y*TileH + yoff%TileH),
			ui::Len(0)
		);

		ui.Draw(v, l.terrain->Img(ui));
	}
	}
}