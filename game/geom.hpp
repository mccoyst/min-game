#pragma once

#include "world.hpp"
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
	Vec2 min, max;

	// Bbox constructs a new bounding box with the given
	// minimum and maximum points.
	Bbox(const World &w, const Vec2 &mn, const Vec2 &mx) : min(mn), max(mx) {
		if (min.x > max.x) {
			Fixed t = max.x;
			max.x = min.x;
			min.x = t;
		}
		if (min.y > max.y) {
			Fixed t = max.y;
			max.y = min.y;
			min.y = t;
		}
		normalize(w);
	}

	// Isect returns the intersection between the two bounding boxes.
	Isection Isect(const World &w, const Bbox &o) const {
		assert (min < w.size);
		assert (max >= Vec2::Zero);
		assert (o.min < w.size);
		assert (o.max >= Vec2::Zero);

		Vec2 overlap(Fixed(-1), Fixed(-1));
		if (Fixed::Between(o.min.x, o.max.x, max.x))
			overlap.x = max.x - o.min.x;
		else if (Fixed::Between(min.x, max.x, o.max.x))
			overlap.x = o.max.x - min.x;
		else if (max.x >= w.size.x && o.max.x < w.size.x &&
			Fixed::Between(o.min.x, o.max.x, max.x % w.size.x))
			overlap.x = o.max.x - max.x % w.size.x;
		else if (o.max.x >= w.size.x && max.x < w.size.x &&
			Fixed::Between(min.x, max.x, o.max.x % w.size.x))
			overlap.x = o.max.x - min.x % w.size.x;

		if (Fixed::Between(o.min.y, o.max.y, max.y))
			overlap.y = max.y - o.min.y;
		else if (Fixed::Between(min.y, max.y, o.max.y))
			overlap.y = o.max.y - min.y;
		else if (max.y >= w.size.y && o.max.y < w.size.y &&
			Fixed::Between(o.min.y, o.max.y, max.y % w.size.y))
			overlap.y = max.y - o.min.y % w.size.y;
		else if (o.max.y >= w.size.y && max.y < w.size.y &&
			Fixed::Between(min.y, max.y, o.max.y % w.size.y))
			overlap.y = o.max.y - max.y % w.size.y;

		return Isection(overlap);
	}

	// Move moves the bounding box, wrapping it so that the minimum
	// point is always within the world's (0,0)--(width-1,height-1).
	void Move(const World &w, const Vec2 &d) {
		min += d;
		max.x += d.x;
		normalize(w);
	}

private:

	// normalize ensures that the minimum coordinate is
	// within the (0,0)--w.size while the maximum coordinate
	// is allowed to pass off of the end of the world.
	void normalize(const World &w) {
		if (min.x >= w.size.x) {
			Fixed width = max.x - min.x;
			min.x %= w.size.x;
			max.x = min.x + width;
		} else if (min.x < Fixed(0)) {
			Fixed width = max.x - min.x;
			min.x = w.size.x + (min.x % w.size.x);
			max.x = min.x + width;
		}

		if (min.y >= w.size.y) {
			Fixed height = max.y - min.y;
			min.y %= w.size.y;
			max.y = min.y + height;
		} else if (min.y < Fixed(0)) {
			Fixed height = max.y - min.y;
			min.y = w.size.y + (min.y % w.size.y);
			max.y = min.y + height;
		}
		assert (min < w.size);
	}
};
