// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

#include "entities.hpp"
#include "ui.hpp"

Astro::Astro(Img *i)
	: box(Vec2{Fixed{}, Fixed{}}, Vec2{Fixed{16},Fixed{16}}),
	sprite(i){
}

Bbox Astro::Box() const{
	return box;
}

void Astro::Move(){
	box.Move(vel);
}

void Astro::MoveTo(Vec2 pos){
	box.min = pos;
}

void Astro::Draw(Ui &ui) const{
	ui.DrawCam(box.min, sprite);
}
