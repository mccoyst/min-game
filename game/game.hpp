// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
#pragma once
#include "fixed.hpp"

using std::string;
using std::unique_ptr;

constexpr Vec2 ScreenDims{ Fixed{800}, Fixed{600} };

class Failure : public std::runtime_error{
public:
	Failure(const string &msg)
		: runtime_error(msg){
	}
};

template<class T, class ...Args>
unique_ptr<T> make_unique(Args&&... args){
	return unique_ptr<T>(new T{std::forward<Args>(args)...});
}
