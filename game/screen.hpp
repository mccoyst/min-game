#pragma once
#include "ui.hpp"
#include <vector>

struct Screen;

struct ScreenStack {
	ScreenStack(std::shared_ptr<Ui>, std::shared_ptr<Screen>);
	~ScreenStack();
	void Run();
	void Push(std::shared_ptr<Screen>);
	void Pop();

	enum {
		// FrameMsec is the minimum frame time in msec.
		FrameMsec = 20,
	};

private:
	std::vector< std::shared_ptr<Screen> > stk;
	std::shared_ptr<Ui> win;
	unsigned long nFrames;
	double meanFrame;
};

struct Screen {
	virtual ~Screen();

	virtual void Update(ScreenStack&) = 0;

	// Draw draws the screen.
	virtual void Draw(std::shared_ptr<Ui>) = 0;

	virtual void Handle(ScreenStack&, Event&) = 0;
}; 
