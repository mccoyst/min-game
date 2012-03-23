#include "game.hpp"
#include <cstdarg>
#include <cstdio>

Failure::Failure(const char *fmt, ...) {
	va_list ap;
	va_start(ap, fmt);
	vsnprintf(msg, sizeof(msg), fmt, ap);
	va_end(ap);
}