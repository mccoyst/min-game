CXX := clang++
CXXFLAGS := -std=c++0x -stdlib=libc++ -I./include -fno-color-diagnostics

OS := $(shell uname | sed 's/.*MINGW.*/win/')

ifeq ($(OS),Darwin)
CXXFLAGS += -framework sfml-graphics
else
CXXFLAGS += -lsfml-graphics -lsfml-window -lsfml-system
endif

UIDEP := sfml

game/minima: game/main.cc ui/ui.a
	@echo prog $@
	@$(CXX) -o $@ $(CXXFLAGS) $^

ui/ui.a: ui/ui.o ui/impl_$(UIDEP).o
	@echo lib ui
	@ar rsc $@ $^

%.o: %.cc
	@$(CXX) -c -o $@ $(CXXFLAGS) $^
