// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
#include "game.hpp"
#include "world.hpp"
#include "ui.hpp"
#include "entities.hpp"
#include "screens.hpp"

class ExploreScreen : public Screen {
public:
	enum {
		// ScrollSpd is the amount to scroll per-frame
		// when an arrow key is held.
		ScrollSpd = 10,
	};

	ExploreScreen(World *);
	virtual ~ExploreScreen();
	virtual void Update(ScreenStack&);
	virtual void Draw(Ui &);
	virtual void Handle(ScreenStack&, Event&);

private:
	World *world;
	unique_ptr<TileView> view;
	unique_ptr<Img> astroimg;
	Astro astro;
};

shared_ptr<Screen> NewExploreScreen(World &w){
	return std::make_shared<ExploreScreen>(&w);
}

ExploreScreen::ExploreScreen(World *w)
	: world(w),
	view{NewTileView((ScreenDims.x/World::TileW).whole() + 2,
		(ScreenDims.y/World::TileH).whole() + 3,
		World::TileW.whole(),
		World::TileH.whole(),
		FindImg("tiles.png"))},
	astroimg(FindImg("Astronaut.png")),
	astro(astroimg.get()){
	astro.MoveTo(world->Start());
}

ExploreScreen::~ExploreScreen() { }

void ExploreScreen::Update(ScreenStack&) {
	astro.Move(*world);
}

void ExploreScreen::Draw(Ui &win) {
	Camera c;
	c.Center(astro.Box().min);
	win.Clear();
	world->Draw(c, win, *view);
	astro.Draw(c, win);
	win.Flip();
}

void ExploreScreen::Handle(ScreenStack &stk, Event &e) {
	if(e.type != Event::KeyDown && e.type != Event::KeyUp)
		return;

	Fixed mv = e.type == Event::KeyDown ? Astro::Speed : Fixed{};

	switch (e.button) {
	case Event::DownArrow:
		astro.vel.y = -mv;
		break;
	case Event::UpArrow:
		astro.vel.y = mv;
		break;
	case Event::LeftArrow:
		astro.vel.x = -mv;
		break;
	case Event::RightArrow:
		astro.vel.x = mv;
		break;
	case Event::Action:
		if(e.type == Event::KeyDown) stk.Pop();
		break;
	case Event::None:
	default:
		// ignore
		break;
	}

}
