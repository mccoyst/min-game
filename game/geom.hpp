// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
#pragma once

#include "fixed.hpp"

struct Color;
class Ui;

// An Isection holds information about an intersection.
class Isection {
public:
	Vec2 overlap;

	// This constructor creates an intersection that represents
	// no intersection.
	Isection() { }

	// This constructor creates an intersection with the given
	// amount of overlap.
	Isection(const Vec2 &o) : overlap(o) { }

	// Converting an intersection to a bool tests if there even
	// was an intersection.
	operator bool() {
		return overlap.x != Fixed(0) && overlap.y != Fixed(0);
	}

	// Area is the area of overlap for this intersection.
	Fixed Area() const {
		return overlap.x * overlap.y;
	}
};

inline bool operator == (const Isection &a, const Isection &b){
	return a.overlap == b.overlap;
}

inline bool operator != (const Fixed &a, const Fixed &b){
	return !(a == b);
}

// A Bbox is a bounding box that can be used for collision
// detection, among other things.
class Bbox {
public:
	Vec2 min, sz;

	// Bbox constructs a new bounding box with the given
	// lower left corner and size.
	Bbox(const Vec2 &p, const Vec2 &s) : min(p), sz(s) { }

	// Isect returns the intersection between the two bounding boxes.
	Isection Isect(const Bbox &o) const;

	// IsectWorld returns the intersection between two bounding
	// boxes that may be wrapping around the edge of the world.
	Isection IsectWorld(Vec2 size, Bbox &o);

	// Center returns the center point of the box.
	Vec2 Center() const {
		const Fixed two(2);
		return Vec2((two*min.x + sz.x) / two, (two*min.y + sz.y) / two);
	}

	// Move moves the bounding box.
	void Move(const Vec2 &d) {
		min += d;
	}

	// WrapMin wraps the bounding box's point so that it's
	// minimum point is within the rectangle defined by the
	// point 0,0 and the given size.
	void WrapMin(const Vec2 &sz);

	// WrapMax wraps the bounding box's point so that it's
	// maximum point is within the rectangle defined by the
	// point 0,0 and the given size.
	void WrapMax(const Vec2 &sz);

	// Wraps returns true if the given bounding box would wrap
	// around the rectangle defined by 0,0 and the given size.
	bool Wraps(const Vec2 &sz) const;

	// Draw draws the bounding box outline.
	void Draw(Ui &win, const Color& c) const;

	// Fill draws the bounding box filled.
	void Fill(Ui &win, const Color& c) const;
};

inline bool operator == (const Bbox &a, const Bbox &b){
	return a.min == b.min && a.sz == b.sz;
}

inline bool operator != (const Bbox &a, const Bbox &b){
	return !(a == b);
}
