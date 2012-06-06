// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
#pragma once
#include "screen.hpp"
#include <stdexcept>

class Img;
class Font;
class World;

constexpr Vec2 ScreenDims{ Fixed{800}, Fixed{600} };

class Failure : public std::runtime_error{
public:
	Failure(const std::string &msg)
		: runtime_error(msg){
	}
};

class ExploreScreen : public Screen {
public:
	enum {
		// ScrollSpd is the amount to scroll per-frame
		// when an arrow key is held.
		ScrollSpd = 10,
	};

	ExploreScreen(World &);
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
	TileView view;
};

class Title : public Screen{
	std::unique_ptr<Font> menu;
	std::unique_ptr<Img> title, start, copyright;
	std::unique_ptr<World> world;
	bool loading;
public:
	Title();
	virtual void Update(ScreenStack&);
	virtual void Draw(Ui&);
	virtual void Handle(ScreenStack&, Event&);
};
