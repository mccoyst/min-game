// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
#pragma once
#include "screen.hpp"
#include <iosfwd>
class World;

class Failure {
public:
	Failure(const char *, ...);
	char msg[128];
};

std::ostream &operator << (std::ostream&, const Failure&);

class ExploreScreen : public Screen {
public:
	enum {
		// ScrollSpd is the amount to scroll per-frame
		// when an arrow key is held.
		ScrollSpd = 10,
	};

	ExploreScreen(Ui &, World &);

	virtual ~ExploreScreen();

	virtual void Update(ScreenStack&);

	virtual void Draw(Ui &);

	virtual void Handle(ScreenStack&, Event&);

private:
	World &world;
	Vec2 mscroll, scroll;	// mouse and keyboard scroll.
	Fixed mul;
	int x0, y0;
	bool drag;
};
