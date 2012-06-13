// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

#include "geom.hpp"
#include "ui.hpp"

Isection Bbox::Isect(const Bbox &o) const {
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

Isection Bbox::IsectWorld(Vec2 size, Bbox &o) {
	if (!Wraps(size) && !o.Wraps(size))
		return Isect(o);

	if (o.Wraps(size)) {
		o.WrapMin(size);
		Isection isect = Isect(o);
		if (isect)
			return isect;

		o.WrapMax(size);
		return Isect(o);
	}

	WrapMin(size);
	Isection isect = Isect(o);
	if (isect)
		return isect;

	WrapMax(size);
	return Isect(o);
}

void Bbox::WrapMin(const Vec2 &sz) {
	if (min >= Vec2{} && min < Vec2{} + sz)
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

void Bbox::WrapMax(const Vec2 &sz) {
	Vec2 max = min + sz;

	if (max >= Vec2{} && max < Vec2{} + sz)
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

bool Bbox::Wraps(const Vec2 &sz) const {
	if (min < Vec2{} || min >= Vec2{} + sz)
		return true;
	Vec2 max = min + sz;
	return max < Vec2{} && max >= Vec2{} + sz;
}

void Bbox::Draw(Ui &win, const Color& c) const {
	win.DrawRect(min, sz, c);
}

void Bbox::Fill(Ui &win, const Color& c) const {
	win.FillRect(min, sz, c);
}