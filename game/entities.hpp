// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
#pragma once

#include "geom.hpp"

class Img;
class Ui;
class World;

class Body {
	Bbox box;

public:
	Body(Bbox box);
	void MoveTo(Vec2);
	void Move(World&);
	Bbox Box() const { return box; }

	Vec2 vel;
};

class Astro : public Body {
	Img *sprite;
public:
	static Fixed Speed;

	Astro(Img*);
	Astro(const Astro&) = default;
	void Draw(Ui&) const;
};
