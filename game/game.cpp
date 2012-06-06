// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
#include "game.hpp"
#include "world.hpp"
#include "ui.hpp"

ExploreScreen::ExploreScreen(World &w)
	: world(w), mul(1), x0(0), y0(0), drag(false),
	view((ScreenDims.x/World::TileW).whole() + 2,
		(ScreenDims.y/World::TileH).whole() + 3,
		World::TileW.whole(),
		World::TileH.whole(),
		std::unique_ptr<Img>(LoadImg("resrc/tiles.png"))){
}

ExploreScreen::~ExploreScreen() { }

void ExploreScreen::Update(ScreenStack&) {
}

void ExploreScreen::Draw(Ui &win) {
	win.MoveCam(mscroll + scroll*mul);
	mscroll = Vec2{};
	win.Clear();
	world.Draw(win, view);
	win.Flip();
}

void ExploreScreen::Handle(ScreenStack &stk, Event &e) {
	Fixed amt(0);
	switch (e.type) {
	case Event::MouseDown:
		scroll = Vec2{};
		drag = true;
		x0 = e.x;
		y0 = e.y;
		break;

	case Event::MouseUp:
		drag = false;
		break;

	case Event::MouseMoved:
		if (!drag)
			break;
		mscroll += Vec2(Fixed(e.x - x0), Fixed(y0 - e.y));
		x0 = e.x;
		y0 = e.y;
		break;

	case Event::KeyDown:
	case Event::KeyUp:
		if (e.type == Event::KeyDown) amt = Fixed(ScrollSpd);
		else amt = Fixed(0);
		switch (e.button) {
		case Event::DownArrow:
			scroll.y = amt;
			break;
		case Event::UpArrow:
			scroll.y = -amt;
			break;
		case Event::LeftArrow:
			scroll.x = amt;
			break;
		case Event::RightArrow:
			scroll.x = -amt;
			break;
		case Event::LShift:
		case Event::RShift:
			if (e.type == Event::KeyDown)
				mul = Fixed(5);
			else
				mul = Fixed(1);
		case Event::Action:
			if(e.type == Event::KeyDown) stk.Pop();
			break;
		case Event::None:
		default:
			// ignore
			break;
		}
		break;

	default:
		// ignore
		break;
	}

}
