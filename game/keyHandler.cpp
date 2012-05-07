// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
#include "keyHandler.hpp"
#include <cassert>
#include <SDL.h>

KeyHandler::KeyHandler(){
	for(int i = 0; i < Event::NumKeys; i++)
		keyState[i] = false;
}


int KeyHandler::KeysDown(){
	return pressedOrder.size();
}


bool KeyHandler::IsPressed(int i){
	if(i >= 0 && i < Event::NumKeys)
		return keyState[i];
	else return false;
}


int KeyHandler::ActiveKey(){
	if(pressedOrder.empty())
		return Event::None;
	else
		return pressedOrder.top();
}

void KeyHandler::PrintKey(int k){
	switch (k){
	case Event::UpArrow:
		fprintf(stderr, "UP\n");
		break;
	case Event::DownArrow:
		fprintf(stderr, "DOWN\n");
		break;
	case Event::LeftArrow:
		fprintf(stderr, "LEFT\n");
		break;
	case Event::RightArrow:
		fprintf(stderr, "RIGHT\n");
		break;
	case Event::LShift:
	case Event::RShift:
		fprintf(stderr, "SHIFT\n");
		break;
	case Event::None:
		fprintf(stderr, "No Key!\n");
		break;
	default:
		fprintf(stderr, "Invalid Key!\n");
		break;
	}
}

int KeyHandler::HandleStroke(SDL_Event &sdle, bool keydown){
	int key = Event::None;

	switch(sdle.key.keysym.sym){
        case SDLK_UP:
		key = Event::UpArrow;
		break;
	case SDLK_DOWN:
		key = Event::DownArrow;
		break;
	case SDLK_LEFT:
		key = Event::LeftArrow;
		break;
	case SDLK_RIGHT:
		key = Event::RightArrow;
		break;
	case SDLK_RSHIFT:
		key = Event::RShift;
		break;
	case SDLK_LSHIFT:
		key = Event::LShift;
		break;
	default:
		return Event::None;
	}

	if (key < Event::NumKeys) keyState[key] = keydown;
	if(keydown && IsStackable(key)) pressedOrder.push(key);
	else FixStack();
	return key;
}


void KeyHandler::PollKeyboard(){
	Uint8 *keystate = SDL_GetKeyState(NULL);

	keyState[Event::LShift] = keystate[SDLK_LSHIFT];
	keyState[Event::RShift] = keystate[SDLK_RSHIFT];
	keyState[Event::RightArrow] = keystate[SDLK_RIGHT];
	keyState[Event::LeftArrow] = keystate[SDLK_LEFT];
	keyState[Event::UpArrow] = keystate[SDLK_UP];
	keyState[Event::DownArrow] = keystate[SDLK_DOWN];
}

void KeyHandler::FixStack(){
	//assumes that the keyState array is correct
	while ((not pressedOrder.empty()) &&
	       (not keyState[pressedOrder.top()]))
	       pressedOrder.pop();
}


bool KeyHandler::IsStackable(int k){
	if (k == Event::LShift || k == Event::RShift)
		return false;
	else return true;
}
