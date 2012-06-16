// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

#include "entities.hpp"
#include "ui.hpp"
#include "world.hpp"

static World::Loc &loc(Bbox, World&);

Body::Body(Bbox b) : box(b) { }

void Body::MoveTo(Vec2 pos){
	box.min = pos;
}

void Body::Move(World &w) {
	const World::Loc &l0 = loc(box, w);
	World::Terrain t0 = w.terrain[l0.terrain];
	box.Move(vel*Speed() * t0.velScale);
}

static World::Loc &loc(Bbox box, World &w) {
	Vec2 c = box.Center();
	int x = Floor(c.x / World::TileW).whole();
	int y = Floor(c.y / World::TileH).whole();
	World::Loc &l = w.AtCoord(x, y);
	return l;
}

void Body::AccelX(int sign) {
	if (sign == 0)
		vel.x = Fixed{};
	else if (sign < 0)
		vel.x = Fixed{-1};
	else
		vel.x = Fixed{1};
}

void Body::AccelY(int sign) {
	if (sign == 0)
		vel.y = Fixed{};
	else if (sign < 0)
		vel.y = Fixed{-1};
	else
		vel.y = Fixed{1};
}

Astro::Astro(Img *i)
	: Body({ {}, {World::TileW,World::TileH}}),
	sprite(i){
}

Fixed Astro::Speed() const{
	return Fixed{2};
}

void Astro::Draw(Ui &ui) const{
	ui.DrawCam(Box().min, *sprite);
}
