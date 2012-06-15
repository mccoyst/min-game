// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
#pragma once

#include "geom.hpp"

class Img;
class Ui;
class World;

class Body {
	Bbox box;
	Vec2 vel;

public:
	Body(Bbox box);
	void MoveTo(Vec2);
	void Move(World&);
	Bbox Box() const { return box; }

	// AccelX sets the direction which the Body is
	// moving along the X axis.  Only the sign is
	// used, the speed is computed based on
	// Speed() and the terrain.
	void AccelX(int sign);

	// AccelY is like AccelX, but in the Y direction.
	void AccelY(int sign);

	// Speed returns the base velocity for this body.
	// The base velocity is used along with terrain
	// to determine the total velocity of a given
	// movement.
	virtual Fixed Speed() const = 0;
};

class Astro : public Body {
	Img *sprite;
public:
	Astro(Img*);
	Astro(const Astro&) = default;

	virtual Fixed Speed() const { return Fixed(2); }

	void Draw(Ui&) const;
};
