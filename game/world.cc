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

void World::Draw(ui::Ui &ui) {
	for (unsigned int x = 0; x < ui.width.whole(); x++) {
	for (unsigned int y = 0; y < ui.height.whole(); y++) {
		const Loc &l = At(x, y);
		ui.Draw(ui::Vec3(ui::Len(x), ui::Len(y), ui::Len(0)),
			l.terrain->Img(ui));
	}
	}
}