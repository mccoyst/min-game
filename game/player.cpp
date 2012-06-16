// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

#include "entities.hpp"
#include "ui.hpp"
#include "world.hpp"
#include <cmath>

static World::Loc &loc(Bbox, World&);

Body::Body(Bbox b) : box(b) { }

void Body::MoveTo(Vec2 pos){
	box.min = pos;
}

void Body::Move(World &w) {
	const World::Loc &l0 = loc(box, w);
	World::Terrain t0 = w.terrain[l0.terrain];
	box.Move(vel*t0.velScale);
}

World::Loc &loc(Bbox box, World &w) {
	Vec2 c = box.Center();
	int x = floor((double) c.x.whole() / World::TileW.whole());
	int y = floor((double) c.y.whole() / World::TileH.whole());
	World::Loc &l = w.AtCoord(x, y);
	return l;
}

Fixed Astro::Speed{2};

Astro::Astro(Img *i)
	: Body(Bbox(Vec2{Fixed{}, Fixed{}}, Vec2{Fixed{16},Fixed{16}})),
	sprite(i){
}

void Astro::Draw(Ui &ui) const{
	ui.DrawCam(Box().min, *sprite);
}
