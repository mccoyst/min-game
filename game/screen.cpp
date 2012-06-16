// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
#include "screens.hpp"
#include "io.hpp"
#include "ui.hpp"
#include <cassert>

Screen::~Screen() { }

ScreenStack::ScreenStack(Ui &w, const shared_ptr<Screen> &screen0)
	: win(w), nFrames(0), meanFrame(0) {
	Push(screen0);
}

void ScreenStack::Run() {
	for ( ; ; ) {
		unsigned long t0 = win.Ticks();

		Event event;
		while (win.PollEvent(event)) {
			if (event.type == Event::Closed)
				return;
			stk.top()->Handle(*this, event);
			if(stk.empty())
				return;
		}

		stk.top()->Draw(win);

		stk.top()->Update(*this);
		if(stk.empty())
			return;

		unsigned long t1 = win.Ticks();
		if (t0 + FrameMsec > t1)
			win.Delay(t0 + FrameMsec - t1);
		nFrames++;
		meanFrame = meanFrame + (t1-t0 - meanFrame)/nFrames;
	}
}

ScreenStack::~ScreenStack() {
	printf(cout(), "Mean Frame Time: %v msec\n", meanFrame);
}

void ScreenStack::Push(const shared_ptr<Screen> &s) {
	stk.push(s);
}

void ScreenStack::Pop() {
	assert(!stk.empty());
	stk.pop();
}
