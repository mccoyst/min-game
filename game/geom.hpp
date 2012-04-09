#pragma once

#include "world.hpp"
#include <vector>
#include <cassert>

// An Isection holds information about an intersection.
struct Isection {
	Vec2 overlap;

	// This constructor creates an intersection that represents
	// no intersection.
	Isection() : overlap(Fixed(-1), Fixed(-1)) { }

	// This constructor creates an intersection with the given
	// amount of overlap.
	Isection(const Vec2 &o) : overlap(o) { }

	// Converting an intersection to a bool tests if there even
	// was an intersection.
	operator bool() {
		return overlap.x < Fixed(0) || overlap.y < Fixed(0);
	}

	// Area is the area of overlap for this intersection.
	Fixed Area() const {
		return overlap.x * overlap.y;
	}
};

// A Bbox is a bounding box that can be used for collision
// detection, among other things.
struct Bbox {
	Vec2 pt, sz;

	// Bbox constructs a new bounding box with the given
	// lower left corner and size.
	Bbox(const Vec2 &p, const Vec2 &s) : pt(p), sz(s) { }

	// Isect returns the intersection between the two bounding boxes.
	Isection Isect(const Bbox &o) const {
		return Isection();
	}

	// Move moves the bounding box, wrapping it so that the minimum
	// point is always within the world's (0,0)--(width-1,height-1).
	void Move(const Vec2 &d) {
		pt += d;
	}
};
