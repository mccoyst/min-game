// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
#pragma once
#include "fixed.hpp"
#include <stdexcept>

constexpr Vec2 ScreenDims{ Fixed{800}, Fixed{600} };

class Failure : public std::runtime_error{
public:
	Failure(const std::string &msg)
		: runtime_error(msg){
	}
};
