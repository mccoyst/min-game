// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

#include "game.hpp"
#include "world.hpp"
#include "screens.hpp"
#include "io.hpp"
#include <istream>

extern unique_ptr<std::istream> Popen(const string&);
static void loadingText(Ui &, Font&);

bool worldOnStdin;

class Title : public Screen{
	unique_ptr<Font> menu;
	unique_ptr<Img> title, start, copyright;
	unique_ptr<World> world;
	bool loading;
public:
	Title();
	virtual void Update(ScreenStack&);
	virtual void Draw(Ui&);
	virtual void Handle(ScreenStack&, Event&);
};

shared_ptr<Screen> NewTitleScreen(){
	return std::make_shared<Title>();
}

Title::Title()
	: menu(FindFont("prstartk.ttf", 12, White)){
	auto tfont (FindFont("prstartk.ttf", 64, White));

	title = move(tfont->Render("MINIMA"));
	start = move(menu->Render("Press f to Start"));
	copyright = move(menu->Render("Copyright © 2012 The Minima Authors"));
}

void Title::Update(ScreenStack &stk){
	if(loading){
		World *w;
		if (worldOnStdin)
			w = new World(cin());
		else
			w = new World(*Popen("wgen"));
		worldOnStdin = false;
		world.reset(w);
		stk.Push(NewExploreScreen(*world));
		loading = false;
	}
}

void Title::Draw(Ui &ui){
	if(loading){
		loadingText(ui, *menu);
		return;
	}

	ui.Clear();

	Vec2 tpos{ Fixed{64}, ui.height - Fixed{128} };
	ui.Draw(tpos, *title);

	auto spos = tpos - Vec2{ Fixed{0}, Fixed{128 } };
	ui.Draw(spos, *start);

	ui.Draw(Vec2{}, *copyright);

	ui.Flip();
}

void Title::Handle(ScreenStack &stk, Event &e){
	//BUG? This is triggering right after the ExploreScreen pops itself.
	if(e.type == Event::KeyDown && e.button == Event::Action)
		loading = true;
}

static void loadingText(Ui &win, Font &font) {
	auto img = font.Render("Generating World");
	win.Clear();
	win.Draw(Vec2(Fixed(0), Fixed(0)), *img);
	win.Flip();
}
