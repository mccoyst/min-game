// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

#include "entities.hpp"
#include "ui.hpp"
#include "world.hpp"

Body::Body(Bbox b) : box(b) { }

void Body::MoveTo(Vec2 pos){
	box.min = pos;
}

void Body::Move(World &w) {
	const World::Loc &l0 = w.AtPoint(box.Center());
	box.Move(vel * w.terrain[l0.terrain].velScale);
}

Fixed Astro::Speed{2};

Astro::Astro(Img *i)
	: Body({ {}, {World::TileW,World::TileH}}),
	sprite(i){
}

void Astro::Draw(Ui &ui) const{
	ui.DrawCam(Box().min, *sprite);
}
