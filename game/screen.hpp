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
	void SetDrawFps(bool);
	
	
	enum {
		// FrameMsec is the minimum frame time in msec.
		FrameMsec = 20,
	
		// FpsTime is the time in msec between redrawing the
		// frame rate.
		FpsTime = 500,
	};

private:
	std::vector< std::shared_ptr<Screen> > stk;
	std::shared_ptr<Ui> win;
	bool drawFps;
	unsigned long nFrames;
	double meanFrame;
};

struct Screen {
	virtual ~Screen();

	virtual void Update(ScreenStack&) = 0;

	// Draw draws the screen.  No need to call Clear() or Flip()
	// on the Ui, the screen stack does this for you.
	virtual void Draw(std::shared_ptr<Ui>) = 0;

	virtual void Handle(ScreenStack&, Event&) = 0;
}; 
