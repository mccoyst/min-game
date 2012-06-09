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
}
