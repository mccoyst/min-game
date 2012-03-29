#pragma once

struct Failure {
	Failure(const char *, ...);
	char msg[128];
};