// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

#include "geom.hpp"

void test_isection(Testo &t){
	t.Assert(Isection{} == false, "default");
}

void test_bbox(Testo &t){
	Bbox b{ Vec2{}, Vec2{Fixed{5}, Fixed{5}} };
	Isection self = b.Isect(b);
	t.Assert(self, "self");
}
