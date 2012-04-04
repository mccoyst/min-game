#pragma once

#include "world.hpp"

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
	Bbox(const Vec2 &mn, const Vec2 &mx) : min(mn), max(mx) {
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
	}

	// Isect returns the intersection between the two bounding boxes.
	Isection Isect(const Bbox &o) const {
		Vec2 overlap(Fixed(-1), Fixed(-1));
		if (Fixed::Between(o.min.x, o.max.x, max.x))
			overlap.x = max.x - o.min.x;
		else if (Fixed::Between(min.x, max.x, o.max.x))
			overlap.x = o.max.x - min.x;

		if (Fixed::Between(o.min.y, o.max.y, max.y))
			overlap.y = max.y - o.min.y;
		else if (Fixed::Between(min.y, max.y, o.max.y))
			overlap.y = o.max.y - min.y;

		return Isection(overlap);
	}

	// Move moves the bounding box, wrapping it so that the minimum
	// point is always within the world's (0,0)--(width-1,height-1).
	void Move(const World &w, const Vec2 &d) {
		min += d;
		if (min.x > w.size.x) {
			Fixed width = max.x - min.x;
			min.x %= w.size.x;
			max.x = min.x + width;
		} else if (min.x < Fixed(0)) {
			Fixed width = max.x - min.x;
			min.x = w.size.x + (min.x % w.size.x);
			max.x = min.x + width;
		} else {
			max.x += d.x;
		}

		if (min.y > w.size.y) {
			Fixed height = max.y - min.y;
			min.y %= w.size.y;
			max.y = min.y + height;
		} else if (min.y < Fixed(0)) {
			Fixed height = max.y - min.y;
			min.y = w.size.y + (min.y % w.size.y);
			max.y = min.y + height;
		} else {
			max.y += d.y;
		}
	}
};
