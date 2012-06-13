// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
#pragma once

#include "geom.hpp"

class Img;
class Ui;

class Astro{
	Bbox box;
	Img *sprite;
public:
	Astro(Img*);
	Astro(const Astro&) = default;

	Bbox Box() const;
	void Move(); //TODO: take the local landscape
	void MoveTo(Vec2);
	void Draw(Ui&) const;

	Vec2 vel;
};
