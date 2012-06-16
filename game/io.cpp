// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

#include "io.hpp"
#include <iostream>

std::ostream &cout(){
	return std::cout;
}

std::ostream &cerr(){
	return std::cerr;
}

std::istream &cin(){
	return std::cin;
}

int printf(std::ostream &out, const string &fmt){
	for(auto c = fmt.begin(), end = fmt.end(); c != end; c++){
		if(*c != '%'){
			out << *c;
			continue;
		}

		c++;
		if(c == end)
			out << "[invalid format string \"" << fmt << "\"]";
		else if(*c == '%')
			out << '%';
		else if(*c == 'v')
			out << "[no more arguments for format \"" << *c << "\"]";
		else
			out << "[invalid format specifier \"" << *c << "\"]";
	}
	return 0;
}
