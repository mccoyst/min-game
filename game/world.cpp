// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
#include "world.hpp"
#include "fixed.hpp"
#include "game.hpp"
#include "ui.hpp"
#include <limits>

static char *readLine(FILE*);

const Fixed World::TileW(16);
const Fixed World::TileH(16);
const Vec2 World::TileSz(TileW, TileH);

World::TerrainType::TerrainType() {
	t.resize(255);
	t['g'] = Terrain('g', 0);
	t['w'] = Terrain('w', 1);
	t['m'] = Terrain('m', 2);
	t['f'] = Terrain('f', 3);
	t['d'] = Terrain('d', 4);
	t['i'] = Terrain('i', 5);

	auto f = std::unique_ptr<Font>(LoadFont("resrc/retganon.ttf", 12, 128, 128, 128));
	htImg.resize(World::MaxHeight+1);
	for (int i = 0; i <= World::MaxHeight; i++)
		htImg[i] = std::unique_ptr<Img>(f->Render("%d", i));
}

World::World(FILE *in) : size(Fixed(0), Fixed(0)) {
	int n = sscanf(readLine(in), "%d %d\n", &width, &height);
	if (n != 2)
		throw Failure("Failed to read width and height");
	if (width <= 0 || height <= 0)
		throw Failure("%d by %d is an invalid world size", width, height);
	if (std::numeric_limits<int>::max() / width < height)
		throw Failure("%d by %d is too big", width, height);

	size = Vec2(Fixed(width), Fixed(height));

	locs.resize(width*height);
	for (int i = 0; i < width*height; i++) {
		char c;
		int h, d;
		char *line = readLine(in);
		n = sscanf(line, " %c %d %d", &c, &h, &d);
		if (n != 3)
			throw Failure("Failed to read a location %u,line [%s]", i, line);
		if (h > MaxHeight || h < 0)
			throw Failure("Location %u has invalid height %d", i, h);
		if (d < 0 || d > h)
			throw Failure("Location %u of height %d has invalid depth %d", i, h, d);
		if (!terrain[c].ch)
			throw Failure("Unknown terrain type %c", c);
		locs[i].height = h;
		locs[i].depth = d;
		locs[i].terrain = c;
	}

	n = sscanf(readLine(in), " %d %d", &x0, &y0);
	if (n != 2)
		throw Failure("Failed to read the start location");
}

void World::Draw(Ui &ui) {
	extern bool drawHeights;

	Fixed w(ui.width / TileW);
	Fixed h(ui.height / TileH);
	Vec2 offs = ui.CamPos();

	for (Fixed x(-1); x <= w; ++x) {
	for (Fixed y(-1); y <= h + Fixed(1); ++y) {
		int xcoord = (x - offs.x/TileW).whole();
		int ycoord = (y - offs.y/TileH).whole();
		const Loc &l = AtCoord(xcoord, ycoord);
		ui.SetTile(x.whole()+1, y.whole()+1, terrain[l.terrain].tile, l.Shade());
	}
	}
	offs.x %= TileW;
	offs.y %= TileH;
	ui.DrawTiles(offs - Vec2(TileW, TileH));

	if (!drawHeights)
		return;

	offs = ui.CamPos();
	for (Fixed x(-1); x <= w; ++x) {
	for (Fixed y(-1); y <= h + Fixed(1); ++y) {
		int xcoord = (x - offs.x/TileW).whole();
		int ycoord = (y - offs.y/TileH).whole();
		const Loc &l = AtCoord(xcoord, ycoord);
		Vec2 v = Vec2(x*TileW, y*TileH) + offs;
		auto txt = terrain.heightImg(l.height);
		ui.DrawCam(v, txt);
	}
	}
}

// shade returns the shade value for the given location
// which is based on its height and depth.
//
// The shade value is computed by linear interpolation
// between 0=minSh and MaxHeight=1.
float World::Loc::Shade() const{
	// minSh is the minimum shade value.
	static const float minSh = 0.25;
	static const float slope = (1 - minSh) / World::MaxHeight;
	return slope*(this->height - this->depth) + minSh;
}

static char *readLine(FILE *in) {
	static char line[4096];
	for (;;) {
		if (fgets(line, sizeof(line), in) == NULL) {
			if (ferror(in)) {
				throw Failure("Error reading the world");
			}
			throw Failure("Unexpected EOF reading the world");
		}
		if (line[0] == '#' || line[0] == '\n' || line[0] == '\r') {
			continue;
		}
		return line;
	}
}
