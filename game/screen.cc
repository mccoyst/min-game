#include "screen.hpp"
#include <cassert>
#include <cstdio>

Screen::~Screen() { }

ScreenStack::ScreenStack(std::shared_ptr<Ui> w, std::shared_ptr<Screen> screen0)
		: win(w), nFrames(0), meanFrame(0) {
	stk.push_back(std::shared_ptr<Screen>(screen0));
}

void ScreenStack::Run() {
	for ( ; ; ) {
		unsigned long t0 = win->Ticks();

		stk.back()->Draw(win);

		Event event;
		while (win->PollEvent(event)) {
			if (event.type == Event::Closed)
				goto out;
			stk.back()->Handle(*this, event);
		}

		stk.back()->Update(*this);

		unsigned long t1 = win->Ticks();
		if (t0 + FrameMsec > t1)
			win->Delay(t0 + FrameMsec - t1);
		nFrames++;
		meanFrame = meanFrame + (t1-t0 - meanFrame)/nFrames;
	}
out:
	return;

}

ScreenStack::~ScreenStack() {
	fprintf(stdout, "Mean Frame Time: %g msec\n", meanFrame);
}

void ScreenStack::Push(std::shared_ptr<Screen> s) {
	stk.push_back(s);
}

void ScreenStack::Pop() {
	assert (stk.size() > 1);
	stk.pop_back();
}