// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

#include "game.hpp"
#include "world.hpp"
#include <istream>

extern std::unique_ptr<std::istream> Popen(const std::string&);
static void loadingText(Ui &, Font*);

Title::Title()
	: menu(LoadFont("resrc/prstartk.ttf", 12, White)){
	auto tfont (LoadFont("resrc/prstartk.ttf", 64, White));

	title.reset(tfont->Render("MINIMA"));
	start.reset(menu->Render("Press f to Start"));
	copyright.reset(menu->Render("Copyright © 2012 The Minima Authors"));
}

void Title::Update(ScreenStack &stk){
	if(loading){
		auto in (Popen("wgen"));
		world.reset(new World(*in.get()));
		stk.Push(std::shared_ptr<ExploreScreen>(new ExploreScreen(*world.get())));
		loading = false;
	}
}

void Title::Draw(Ui &ui){
	if(loading){
		loadingText(ui, menu.get());
		return;
	}

	ui.Clear();

	Vec2 tpos{ Fixed{64}, ui.height - Fixed{128} };
	ui.Draw(tpos, title.get());

	auto spos = tpos - Vec2{ Fixed{0}, Fixed{128 } };
	ui.Draw(spos, start.get());

	ui.Draw(Vec2{}, copyright.get());

	ui.Flip();
}

void Title::Handle(ScreenStack &stk, Event &e){
	//BUG? This is triggering right after the ExploreScreen pops itself.
	if(e.type == Event::KeyDown && e.button == Event::Action)
		loading = true;
}

static void loadingText(Ui &win, Font *font) {
	auto img = std::unique_ptr<Img>(font->Render("Generating World"));
	win.Clear();
	win.Draw(Vec2(Fixed(0), Fixed(0)), img.get());
	win.Flip();
}
