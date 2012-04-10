#include "screen.hpp"
#include <cassert>
#include <cstdio>

Screen::~Screen() { }

ScreenStack::ScreenStack(std::shared_ptr<Ui> w, std::shared_ptr<Screen> screen0)
		: win(w), drawFps(false), nFrames(0), meanFrame(0) {
	stk.push_back(std::shared_ptr<Screen>(screen0));
}

void ScreenStack::Run() {
	auto fpsFont = LoadFont("resrc/prstartk.ttf", 12, 255, 255, 255);
	std::shared_ptr<Img> fps(0);
	unsigned long lastFpsTime = 0;
	unsigned long lastFpsFrames = 0;

	for ( ; ; ) {
		unsigned long t0 = win->Ticks();

		win->Clear();

		stk.back()->Draw(win);

		if (drawFps && lastFpsTime + FpsTime <= t0) {
			unsigned long rate = (nFrames - lastFpsFrames)/(FpsTime/1000.0);
			fps = fpsFont->Render("%lu fps", rate);
			lastFpsTime = t0;
			lastFpsFrames = nFrames;
		}
		if (drawFps && fps)
			win->Draw(Vec2::Zero, fps);

		win->Flip();

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

void ScreenStack::SetDrawFps(bool b) {
	drawFps = b;
}

void ScreenStack::Push(std::shared_ptr<Screen> s) {
	stk.push_back(s);
}

void ScreenStack::Pop() {
	assert (stk.size() > 1);
	stk.pop_back();
}