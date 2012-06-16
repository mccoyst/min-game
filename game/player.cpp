// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

#include "entities.hpp"
#include "ui.hpp"
#include "world.hpp"
#include <cassert>

static Fixed meanHeight(Bbox, World&);

Body::Body(Bbox b) : box(b) { }

void Body::MoveTo(Vec2 pos){
	box.min = pos;
}

void Body::Move(World &w) {
	if (vel == Vec2{})
		return;
	
	Bbox box1(box);
	box1.Move(vel);
	const Fixed MaxDh (0, 6);
	Fixed dh = Min(Abs(meanHeight(box1, w) - meanHeight(box, w)), MaxDh);
	Fixed slope = (Fixed{0, 1} - Fixed{1}) / MaxDh;
	Fixed scale = dh*slope + Fixed{1};
	Vec2 v = vel * scale * w.terrain[w.AtPoint(box.Center()).terrain].velScale;
	if (vel.x < Fixed{} && v.x == Fixed{})
		v.x = -Fixed{0, 1};
	if (vel.x > Fixed{} && v.x == Fixed{})
		v.x = Fixed{0, 1};
	if (vel.y < Fixed{} && v.y == Fixed{})
		v.y = -Fixed{0, 1};
	if (vel.y > Fixed{} && v.y == Fixed{})
		v.y = Fixed{0, 1};
	box.Move(v);
}

Fixed meanHeight(Bbox box, World &w) {
	const World::Loc &l = w.AtPoint(box.Center());
	const World::Loc *locs[9];
	int i = 0;
	for (int dx = -1; dx <= 1; dx++) {
	for (int dy = -1; dy <= 1; dy++) {
		locs[i++] = &w.AtCoord(l.x + dx, l.y + dy);
	}
	}

	Fixed sum (0), area (0);
	for (int i = 0; i < 9; i++) {
		Bbox lbox = locs[i]->Box();
		Isection is = Bbox(box).IsectWorld(w.Size(), lbox);
		if (!is)
			continue;
		sum += Fixed{locs[i]->height}*is.Area();
		area += is.Area();
	}
	return sum / area;
}

Fixed Astro::Speed{2};

Astro::Astro(Img *i)
	: Body(Bbox(Vec2{Fixed{}, Fixed{}}, Vec2{Fixed{16},Fixed{16}})),
	sprite(i){
}

void Astro::Draw(Ui &ui) const{
	ui.DrawCam(Box().min, *sprite);
}
