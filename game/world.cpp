// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
#include "world.hpp"
#include "game.hpp"
#include "io.hpp"
#include "ui.hpp"
#include <string>
#include <istream>
#include <limits>

bool drawHeights = false;

static string readLine(std::istream&);

const Fixed World::TileW(16);
const Fixed World::TileH(16);
const Vec2 World::TileSz(TileW, TileH);

World::TerrainType::TerrainType()
	: t({
		{ 'g', Terrain(0, Fixed(1)) },
		{ 'w', Terrain(1, Fixed(0, 1)) },
		{ 'm', Terrain(2, Fixed(0, 5)) },
		{ 'f', Terrain(3, Fixed(0, 10)) },
		{ 'd', Terrain(4, Fixed(0, 5)) },
		{ 'i', Terrain(5, Fixed(0, 5)) },
	}){
	auto f = LoadFont("resrc/retganon.ttf", 12, Gray);
	htImg.resize(World::MaxHeight+1);
	for (int i = 0; i <= World::MaxHeight; i++)
		htImg[i] = unique_ptr<Img>(f->Render(std::to_string(i)));
}

Bbox World::Loc::Box() const {
	Vec2 min(Fixed{x*World::TileW.whole()}, Fixed{y*World::TileH.whole()});
	return Bbox(min, World::TileSz);
}

World::World(std::istream &in) : dim(Fixed(0), Fixed(0)) {
	int n = sscanf(readLine(in).c_str(), "%d %d\n", &width, &height);
	if (n != 2)
		throw Failure("Failed to read width and height");
	if (width <= 0 || height <= 0)
		throw Failure(sprintf("%v by %v is an invalid world size", width, height));
	if (std::numeric_limits<int>::max() / width < height)
		throw Failure(sprintf("%v by %v is too big", width, height));

	dim = Vec2(Fixed(width), Fixed(height));
	size = dim;
	size.x *= TileW;
	size.y *= TileH;

	locs.resize(width*height);
	for (int i = 0; i < width*height; i++) {
		char c;
		int h, d;
		auto line = readLine(in);
		n = sscanf(line.c_str(), " %c %d %d", &c, &h, &d);
		if (n != 3)
			throw Failure(sprintf("Failed to read a location %v,line [%v]", i, line));
		if (h > MaxHeight || h < 0)
			throw Failure(sprintf("Location %v has invalid height %v", i, h));
		if (d < 0 || d > h)
			throw Failure(sprintf("Location %v of height %v has invalid depth %v", i, h, d));
		if (!terrain.contains(c))
			throw Failure(sprintf("Unknown terrain type %v", c));
		locs[i].height = h;
		locs[i].depth = d;
		locs[i].terrain = c;
		locs[i].x = i/height;
		locs[i].y = i%height;
	}

	n = sscanf(readLine(in).c_str(), " %d %d", &x0, &y0);
	if (n != 2)
		throw Failure("Failed to read the start location");
}

void World::Draw(Ui &ui, TileView &view) {
	Fixed w(ui.width / TileW);
	Fixed h(ui.height / TileH);
	Vec2 offs = ui.CamPos();

	for (Fixed y(-1); y <= h + Fixed(1); ++y) {
	for (Fixed x(-1); x <= w; ++x) {
		int xcoord = (x - Trunc(offs.x/TileW)).whole();
		int ycoord = (y - Trunc(offs.y/TileH)).whole();
		const Loc &l = AtCoord(xcoord, ycoord);
		view.SetTile(x.whole()+1, y.whole()+1, terrain[l.terrain].tile, l.Shade());
	}
	}
	offs.x %= TileW;
	offs.y %= TileH;
	ui.Draw(offs - Vec2(TileW, TileH), view);

	if (!drawHeights)
		return;

	offs = ui.CamPos();
	for (Fixed y(-1); y <= h + Fixed(1); ++y) {
	for (Fixed x(-1); x <= w; ++x) {
		int xcoord = (x - Trunc(offs.x/TileW)).whole();
		int ycoord = (y - Trunc(offs.y/TileH)).whole();
		const Loc &l = AtCoord(xcoord, ycoord);
		auto &txt = terrain.heightImg(l.height);
		Vec2 pt = Vec2(x*TileW, y*TileH);
		pt.x += offs.x % TileW;
		pt.y += offs.y % TileH;
		ui.Draw(pt, txt);
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

static string readLine(std::istream &in) {
	string line;
	for (;;) {
		if (!getline(in, line)){
			if(in.eof())
				throw Failure("Unexpected EOF reading the world");

			throw Failure("Error reading the world");
		}
		if (line.empty() || line[0] == '#' || line[0] == '\r') {
			continue;
		}
		return line;
	}
}

World::Loc& World::AtCoord(int x, int y) {
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

World::Loc &World::AtPoint(Vec2 pt) {
	int x = Floor(pt.x / World::TileW).whole();
	int y = Floor(pt.y / World::TileH).whole();
	return AtCoord(x, y);
}
