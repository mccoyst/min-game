// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
#include "world.hpp"
#include "game.hpp"
#include "io.hpp"
#include "ui.hpp"
#include <limits>

using std::vector;

bool drawHeights = false;

static string readLine(std::istream&);

const Fixed World::TileW(16);
const Fixed World::TileH(16);
const Vec2 World::TileSz(TileW, TileH);

namespace{
	// terrain is an array of Terrain indexed by the
	// character representation of the Terrain.
	class TerrainType {
		std::map<char, Terrain> t;
		vector<unique_ptr<Img>> htImg;
	public:
		TerrainType();

		// operator[] returns the terrain with the given character
		// representation.
		Terrain operator [] (char c) const { return t.at(c); }

		// contains returns true iff a terrain type is defined
		// for the given char.
		bool contains(char c) const { return t.find(c) != t.end(); }

		// heightImg returns an image containing the text for
		// the given height value.
		Img &heightImg(int ht) { return *htImg[ht]; }
	};
}

class World::Impl{
public:
	TerrainType terrain;
	vector<Loc> locs;
	int width, height;

	// The world's dimensions in numer of tiles.
	Vec2 dim;

	// The world's size, dim*World::TileSz.
	Vec2 size;

	int x0, y0;

	Impl(std::istream &in);
};

World::Impl::Impl(std::istream &in){
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

TerrainType::TerrainType()
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

World::World(std::istream &in)
	: impl(new Impl(in)) {
}

World::~World(){
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
		view.SetTile(x.whole()+1, y.whole()+1,
			impl->terrain[l.terrain].tile, l.Shade());
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
		auto &txt = impl->terrain.heightImg(l.height);
		Vec2 pt = Vec2(x*TileW, y*TileH);
		pt.x += offs.x % TileW;
		pt.y += offs.y % TileH;
		ui.Draw(pt, txt);
	}
	}
}

World::Loc &World::At(unsigned int x, unsigned int y){
	return impl->locs.at(x*impl->dim.y.whole()+y);
}

const World::Loc &World::At(unsigned int x, unsigned int y) const{
	return const_cast<World*>(this)->At(x,y);
}

Terrain World::TerrainAt(Vec2 p) const{
	return impl->terrain[AtPoint(p).terrain];
}

Vec2 World::Size() const{
	return impl->size;
}

Vec2 World::Start() const{
	return Vec2(Fixed{impl->x0}*TileW, Fixed{impl->y0}*TileH);
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

const World::Loc& World::AtCoord(int x, int y) const{
	if (x < 0)
		x = impl->width - -x%impl->width;
	else
		x %= impl->width;
	if (y < 0)
		y = impl->height - -y%impl->height;
	else
		y %= impl->height;
	return At(x, y);
	
}

const World::Loc &World::AtPoint(Vec2 pt) const{
	int x = Floor(pt.x / World::TileW).whole();
	int y = Floor(pt.y / World::TileH).whole();
	return AtCoord(x, y);
}
