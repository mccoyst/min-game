// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
#pragma once

#include "world.hpp"
#include "ui.hpp"
#include <vector>
#include <cassert>

// An Isection holds information about an intersection.
class Isection {
public:
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
class Bbox {
public:
	Vec2 min, sz;

	// Bbox constructs a new bounding box with the given
	// lower left corner and size.
	Bbox(const Vec2 &p, const Vec2 &s) : min(p), sz(s) { }

	// Isect returns the intersection between the two bounding boxes.
	Isection Isect(const Bbox &o) const {
		Isection isect;
		Vec2 max = min + sz, omax = o.min + o.sz;
		if (max.x >= o.min.x && max.x <= omax.x)
			isect.overlap.x = max.x - o.min.x;
		else if (omax.x >= min.x && omax.x <= max.x)
			isect.overlap.x = omax.x - min.x;

		if (max.y >= o.min.y && max.y <= omax.y)
			isect.overlap.y = max.y - o.min.y;
		else if (omax.y >= min.y && omax.y <= max.y)
			isect.overlap.y = omax.y - min.y;

		return isect;
	}

	// IsectWorld returns the intersection between two bounding
	// boxes that may be wrapping around the edge of the world.
	Isection IsectWorld(const World &w, Bbox &o) {
		if (!Wraps(w.size) && !o.Wraps(w.size))
			return Isect(o);

		if (o.Wraps(w.size)) {
			o.WrapMin(w.size);
			Isection isect = Isect(o);
			if (isect)
				return isect;
	
			o.WrapMax(w.size);
			return Isect(o);
		}

		WrapMin(w.size);
		Isection isect = Isect(o);
		if (isect)
			return isect;

		WrapMax(w.size);
		return Isect(o);
	}

	// Move moves the bounding box.
	void Move(const Vec2 &d) {
		min += d;
	}

	// WrapMin wraps the bounding box's point so that it's
	// minimum point is within the rectangle defined by the
	// point 0,0 and the given size.
	void WrapMin(const Vec2 &sz) {
		if (min >= Vec2::Zero && min < Vec2::Zero + sz)
			return;

		if (min.x < Fixed(0))
			min.x = sz.x - (-min.x % sz.x);
		else if (min.x >= sz.x)
			min.x %= sz.x;

		if (min.y < Fixed(0))
			min.y = sz.y - (-min.y % sz.y);
		else if (min.y >= sz.y)
			min.y %= sz.y;
	}

	// WrapMax wraps the bounding box's point so that it's
	// maximum point is within the rectangle defined by the
	// point 0,0 and the given size.
	void WrapMax(const Vec2 &sz) {
		Vec2 max = min + sz;

		if (max >= Vec2::Zero && max < Vec2::Zero + sz)
			return;

		if (max.x < Fixed(0))
			max.x = sz.x - (-max.x % sz.x);
		else if (max.x >= sz.x)
			max.x %= sz.x;

		if (max.y < Fixed(0))
			max.y = sz.y - (-max.y % sz.y);
		else if (max.y >= sz.y)
			max.y %= sz.y;

		min = max - sz;
	}

	// Wraps returns true if the given bounding box would wrap
	// around the rectangle defined by 0,0 and the given size.
	bool Wraps(const Vec2 &sz) const {
		if (min < Vec2::Zero || min >= Vec2::Zero + sz)
			return true;
		Vec2 max = min + sz;
		return max < Vec2::Zero && max >= Vec2::Zero + sz;
	}

	// Draw draws the bounding box outline.
	void Draw(std::shared_ptr<Ui> win, const Color& c) const {
		win->DrawRect(min, sz, c);
	}

	// Fill draws the bounding box filled.
	void Fill(std::shared_ptr<Ui> win, const Color& c) const {
		win->FillRect(min, sz, c);
	}
};
