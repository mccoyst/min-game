// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
/**
 * Handle input from the keyboard
 *
 * Take keydown events, store keys pressed in a stack, with the last key
 * pressed active (ie hold left, tap right, then left is still active)
 * Could be used for key combinations as well, maybe also combos by passing
 * a window over the top of the stack.
 */

#pragma once

#include "keyMap.hpp"
#include <stack>

union SDL_Event;

class KeyHandler {
public:
	KeyHandler();

	//returns the number of pressed keys
	int KeysDown();

	//is the given key pressed
	bool IsPressed(int i);

	//returns the active key
	int ActiveKey();

	//handles a single Key Stroke
	int HandleStroke(SDL_Event &sdle, bool keydown);

	//prints english thing for key
	void PrintKey(int k);

private:
	const static int MAX_PRESS = 3;
	bool keyState[keymap::NumKeys];
	std::stack<int> pressedOrder;

	/* in the event that more than n-keys are depressed, we need to start
	   polling keyboard state because modern keyboards suck. */
	void PollKeyboard();

	// fixes the activeKey by assuring top of the stack is still pressed
	void FixStack();

	//does key k go to the stack?
	bool IsStackable(int k);
};



