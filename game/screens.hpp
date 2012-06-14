// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
#pragma once
#include "ui.hpp"
#include "fixed.hpp"
#include <memory>
#include <vector>

using std::shared_ptr;
using std::unique_ptr;
using std::vector;

struct Event;
class Screen;
class Ui;
class World;
class Astro;

class ScreenStack {
public:
	// Creates a new screen stack with the given initial screen.
	ScreenStack(Ui &, const shared_ptr<Screen>&);

	~ScreenStack();

	// Run runs the main loop of the program, calling the
	// Draw(), Handle(), and Update() methods on the top
	// screen on the stack.
	void Run();

	// Push pushes a new screen onto the top of the stack.
	void Push(const shared_ptr<Screen>&);

	// Pop pops the current screen off of the top of the stack.
	void Pop();

	enum {
		// FrameMsec is the minimum frame time in msec.
		FrameMsec = 20,
	};

private:
	vector<shared_ptr<Screen>> stk;
	Ui &win;
	unsigned long nFrames;
	double meanFrame;
};

class Screen {
public:
	virtual ~Screen();

	// Draw draws the screen.  The Screen is responsible for
	// calling Clear() and Flip().
	virtual void Draw(Ui &) = 0;

	// Handle is called for each event coming from the
	// Ui, with the exception of the Close event which is
	// intercepted by the ScreenStack to exit the program.
	virtual void Handle(ScreenStack&, Event&) = 0;

	// Update is called after all of the events are handled and before
	// the next frame is drawn in order to allow the screen to update
	// its state based on the events.
	virtual void Update(ScreenStack&) = 0;
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
	TileView view;
	unique_ptr<Astro> astro;
	unique_ptr<Img> astroimg;
};

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
