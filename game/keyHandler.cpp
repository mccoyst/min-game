#include "keyHandler.hpp"
#include <cassert>
#include <SDL.h>
using namespace std;

KeyHandler::KeyHandler(){
	for(int i = 0; i < keymap::NumKeys; i++)
		keyState[i] = false;
}


int KeyHandler::KeysDown(){
	return pressedOrder.size();
}


bool KeyHandler::IsPressed(int i){
	if(i >= 0 && i < keymap::NumKeys)
		return keyState[i];
	else return false;
}


int KeyHandler::ActiveKey(){
	if(pressedOrder.empty())
		return keymap::None;
	else
		return pressedOrder.top();
}

void KeyHandler::PrintKey(int k){
	using namespace keymap;
	switch (k){
	case UpArrow:
		fprintf(stderr, "UP\n");
		break;
	case DownArrow:
		fprintf(stderr, "DOWN\n");
		break;
	case LeftArrow:
		fprintf(stderr, "LEFT\n");
		break;
	case RightArrow:
		fprintf(stderr, "RIGHT\n");
		break;
	case LShift:
	case RShift:
		fprintf(stderr, "SHIFT\n");
		break;
	case Invalid:
		fprintf(stderr, "Invalid Key!\n");
		break;
	case None:
		fprintf(stderr, "No Key!\n");
		break;
	}
}

int KeyHandler::HandleStroke(SDL_Event &sdle, bool keydown){
	using namespace keymap;
	int key = Invalid;

	switch(sdle.key.keysym.sym){
        case SDLK_UP:
		key = UpArrow;
		break;
	case SDLK_DOWN:
		key = DownArrow;
		break;
	case SDLK_LEFT:
		key = LeftArrow;
		break;
	case SDLK_RIGHT:
		key = RightArrow;
		break;
	case SDLK_RSHIFT:
		key = RShift;
		break;
	case SDLK_LSHIFT:
		key = LShift;
		break;
	default:
		return Invalid;
	}

	if (key < NumKeys) keyState[key] = keydown;
	if(keydown && IsStackable(key)) pressedOrder.push(key);
	else FixStack();
	return key;
}


void KeyHandler::PollKeyboard(){
	using namespace keymap;
	Uint8 *keystate = SDL_GetKeyState(NULL);

	keyState[LShift] = keystate[SDLK_LSHIFT];
	keyState[RShift] = keystate[SDLK_RSHIFT];
	keyState[RightArrow] = keystate[SDLK_RIGHT];
	keyState[LeftArrow] = keystate[SDLK_LEFT];
	keyState[UpArrow] = keystate[SDLK_UP];
	keyState[DownArrow] = keystate[SDLK_DOWN];
}

void KeyHandler::FixStack(){
	//assumes that the keyState array is correct
	while ((not pressedOrder.empty()) &&
	       (not keyState[pressedOrder.top()]))
	       pressedOrder.pop();
}


bool KeyHandler::IsStackable(int k){
	if (k == keymap::LShift || k == keymap::RShift)
		return false;
	else return true;
}
