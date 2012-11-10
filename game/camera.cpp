// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
#include "ui.hpp"
#include "game.hpp"

void Camera::Move(Vec2 v){
	cam += v;
}

void Camera::Center(Vec2 v){
	cam.x = ScreenDims.x/Fixed{2} - v.x;
	cam.y = ScreenDims.y/Fixed{2} - v.y;
}

Vec2 Camera::Pos() const{
	return cam;
}

void Camera::Draw(Vec2 p, Ui &ui, Img &i, float shade){
	ui.Draw(p + cam, i, shade);
}
