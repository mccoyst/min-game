// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

#include "geom.hpp"

void test_isection(Testo &t){
	t.Assert(Isection{} == false, "default");
}

void test_bbox(Testo &t){
	Bbox b{ Vec2{}, Vec2{Fixed{5}, Fixed{5}} };
	Isection self = b.Isect(b);
	t.Assert(self, "self");
	t.Assert(self.overlap.x == b.sz.x, "self.x");
	t.Assert(self.overlap.y == b.sz.y, "self.y");

	Bbox c{ Vec2{Fixed{1}, Fixed{1}}, Vec2{Fixed{5}, Fixed{5}} };
	Isection bc = b.Isect(c);
	Isection cb = c.Isect(b);
	t.Assert(bc, "bc");
	t.Assert(cb, "cb");
	t.Assert(bc == cb, "bc == cb");
	t.Assert(bc.overlap.x == b.sz.x - Fixed{1}, "bc.x == b.x-1");
	t.Assert(bc.overlap.y == b.sz.y - Fixed{1}, "bc.y == b.y-1");

	Vec2 world{ Fixed{10}, Fixed{10} };
	Isection wself = b.IsectWorld(world, b);
	t.Assert(wself, "wrap self");
	t.Assert(wself.overlap.x == b.sz.x, "wrap self.x");
	t.Assert(wself.overlap.y == b.sz.y, "wrap self.y");

	Isection wbc = b.IsectWorld(world, c);
	Isection wcb = c.IsectWorld(world, b);
	t.Assert(wbc, "wrap bc");
	t.Assert(wcb, "wrap cb");
	t.Assert(wbc == wcb, "wrap bc == cb");
	t.Assert(wbc.overlap.x == b.sz.x - Fixed{1}, "wrap bc.x == b.x-1");
	t.Assert(wbc.overlap.y == b.sz.y - Fixed{1}, "wrap bc.y == b.y-1");

	Bbox d{ Vec2{Fixed{-1}, Fixed{-1}}, Vec2{Fixed{5}, Fixed{5}} };
	Isection bd = b.Isect(d);
	Isection db = d.Isect(b);
	Isection cd = c.Isect(d);
	Isection dc = d.Isect(c);
	t.Assert(bd, "bd");
	t.Assert(db, "db");
	t.Assert(bd == db, "bd == db");
	t.Assert(!(bd != db), "bd doesn't not equal db");
	t.Assert(bd.overlap.x == b.sz.x - Fixed{1}, "bd.x == b.x-1");
	t.Assert(bd.overlap.y == b.sz.y - Fixed{1}, "bd.y == b.y-1");
	t.Assert(cd, "cd");
	t.Assert(dc, "dc");
	t.Assert(cd == dc, "cd == dc");
	t.Assert(cd.overlap.x == c.sz.x - Fixed{2}, "cd.x == c.x-2");
	t.Assert(cd.overlap.y == c.sz.y - Fixed{2}, "cd.y == c.y-2");

	Isection wbd = b.IsectWorld(world, d);
	Isection wdb = d.IsectWorld(world, b);
	Isection wcd = c.IsectWorld(world, d);
	Isection wdc = d.IsectWorld(world, c);
	t.Assert(wbd, "bd");
	t.Assert(wdb, "db");
	t.Assert(wbd == wdb, "bd == db");
	t.Assert(wbd.overlap.x == b.sz.x - Fixed{1}, "wrap bd.x == b.x-1");
	t.Assert(wbd.overlap.y == b.sz.y - Fixed{1}, "wrap bd.y == b.y-1");
	t.Assert(wcd, "cd");
	t.Assert(wdc, "dc");
	t.Assert(wcd == wdc, "wrap cd == dc");
	t.Assert(wcd.overlap.x == c.sz.x - Fixed{2}, "wrap cd.x == c.x-2");
	t.Assert(wcd.overlap.y == c.sz.y - Fixed{2}, "wrap cd.y == c.y-2");
}
