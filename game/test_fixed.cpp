#include "fixed.hpp"

void test_trivial(Testo &t){
	Fixed zero{};
	Fixed one{1};

	t.Assert(zero == zero, "zero != zero");
	t.Assert(zero != one, "zero == one");
	t.Assert(zero + zero == zero, "zero + zero != zero");
	t.Assert(one + one == Fixed{2}, "1 + 1 != 2");
	t.Assert(zero - one == Fixed{-1}, "0 - 1 != -1");
	t.Assert(one*one == one, "1*1 != 1");
}

void test_fractions(Testo &t){
	Fixed one{1};

	t.Assert(one/Fixed{1} == one, "1/1 != 1");
	t.Assert(one/Fixed{2} == Fixed{0,8}, "1/2 != ½");
	t.Assert(Fixed{}/Fixed{1} == Fixed{}, "0/1 != 0");
	t.Assert(one/Fixed{16} == Fixed{0,1}, "1/16");
	t.Assert(one/Fixed{17} == Fixed{}, "1/17");

	t.Assert(Fixed{2}/Fixed{16} == Fixed{0,2}, "2/16");
	t.Assert(Fixed{16}/Fixed{16} == one, "16/16");
	t.Assert(Fixed{17}/Fixed{16} == Fixed{1,1}, "17/16");
}

void test_multiply(Testo &t){
	t.Assert(Fixed{2}*Fixed{2} == Fixed{4}, "2*2");
	t.Assert(Fixed{2}*Fixed{3} == Fixed{6}, "2*3");
	t.Assert(Fixed{3}*Fixed{3} == Fixed{9}, "3*3");

	t.Assert(Fixed{3}*(Fixed{2}/Fixed{16}) == Fixed{0,6}, "3*(2/16)");
}
