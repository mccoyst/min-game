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
	TileView view;
	unique_ptr<Img> astroimg;
	Astro astro;
};

shared_ptr<Screen> NewExploreScreen(World &w){
	return std::make_shared<ExploreScreen>(&w);
}

ExploreScreen::ExploreScreen(World *w)
	: world(w),
	view((ScreenDims.x/World::TileW).whole() + 2,
		(ScreenDims.y/World::TileH).whole() + 3,
		World::TileW.whole(),
		World::TileH.whole(),
		LoadImg("resrc/tiles.png")),
	astroimg(LoadImg("resrc/Astronaut.png")),
	astro(astroimg.get()){
}

ExploreScreen::~ExploreScreen() { }

void ExploreScreen::Update(ScreenStack&) {
	astro.Move(*world);
}

void ExploreScreen::Draw(Ui &win) {
	win.CenterCam(astro.Box().min);
	win.Clear();
	world->Draw(win, view);
	astro.Draw(win);
	win.Flip();
}

void ExploreScreen::Handle(ScreenStack &stk, Event &e) {
	if(e.type != Event::KeyDown && e.type != Event::KeyUp)
		return;

	int move = e.type == Event::KeyDown;

	switch (e.button) {
	case Event::DownArrow:
		astro.AccelY(-move);
		break;
	case Event::UpArrow:
		astro.AccelY(move);
		break;
	case Event::LeftArrow:
		astro.AccelX(-move);
		break;
	case Event::RightArrow:
		astro.AccelX(move);
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
