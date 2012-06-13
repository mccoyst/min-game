// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
#pragma once

#include <ostream>
#include <sstream>
#include <string>

using std::string;

std::ostream &cout();
std::ostream &cerr();

int printf(std::ostream &out, const string &fmt);

template<class Arg, class... Args>
int printf(std::ostream &out, const string &fmt, const Arg &arg, const Args&... args){
	auto c = fmt.begin();
	auto end = fmt.end();
	while(c != end && *c != '%'){
		out << *c;
		c++;
	}

	c++;
	if(c == end)
		out << "[invalid format string \"" << fmt << "\"]";
	else if(*c == '%')
		out << '%';
	else if(*c == 'v')
		out << arg;
	else
		out << "[invalid format specifier \"" << *c << "\"]";
	c++;
	return 1 + printf(out, fmt.substr(c-fmt.begin()), args...);
}

template<class... Args>
string sprintf(const string &fmt, const Args&... args){
	std::ostringstream out;
	printf(out, fmt, args...);
	return out.str();
}
