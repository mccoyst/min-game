// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
#pragma once
#include "screen.hpp"
#include "ui.hpp"

class World;
class Astro;

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
	TileView view;
	std::unique_ptr<Astro> astro;
	std::unique_ptr<Img> astroimg;
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
