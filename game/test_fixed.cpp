#include "fixed.hpp"

void test_trivial(Testo &t){
	Fixed zero{};
	Fixed one{1};

	t.Assert(zero == zero, "zero != zero");
	t.Assert(zero != one, "zero == one");
	t.Assert(zero + zero == zero, "zero + zero != zero");
}
