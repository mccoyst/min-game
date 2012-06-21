// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
#include "screens.hpp"
#include "io.hpp"
#include "ui.hpp"
#include <cassert>
#include <stack>

class ScreenStack::Impl{
public:
	std::stack<shared_ptr<Screen>> stk;
	Ui &win;
	unsigned long nFrames;
	double meanFrame;
	Impl(Ui &);
};

ScreenStack::Impl::Impl(Ui &w)
	: win(w), nFrames(0), meanFrame(0){
}

Screen::~Screen(){
}

ScreenStack::ScreenStack(Ui &w, const shared_ptr<Screen> &screen0)
	: impl(new Impl(w)){
	Push(screen0);
}

void ScreenStack::Run() {
	for ( ; ; ) {
		unsigned long t0 = impl->win.Ticks();

		Event event;
		while (impl->win.PollEvent(event)) {
			if (event.type == Event::Closed)
				return;
			impl->stk.top()->Handle(*this, event);
			if(impl->stk.empty())
				return;
		}

		impl->stk.top()->Draw(impl->win);

		impl->stk.top()->Update(*this);
		if(impl->stk.empty())
			return;

		unsigned long t1 = impl->win.Ticks();
		if (t0 + FrameMsec > t1)
			impl->win.Delay(t0 + FrameMsec - t1);
		impl->nFrames++;
		impl->meanFrame += (t1-t0 - impl->meanFrame)/impl->nFrames;
	}
}

ScreenStack::~ScreenStack() {
	printf(cout(), "Mean Frame Time: %v msec\n", impl->meanFrame);
}

void ScreenStack::Push(const shared_ptr<Screen> &s) {
	impl->stk.push(s);
}

void ScreenStack::Pop() {
	assert(!impl->stk.empty());
	impl->stk.pop();
}
