CXXFLAGS :=-std=c++0x -I./include
LDFLAGS :=

OS := $(shell uname | sed 's/.*MINGW.*/win/')

ifeq ($(OS),Darwin)
CXX := clang++ -fno-color-diagnostics -stdlib=libc++
CXXFLAGS += \
	-framework sfml-graphics \
	-framework sfml-window \
	-framework sfml-system \

else
CXX := g++
LDFLAGS += -lsfml-graphics -lsfml-window -lsfml-system
endif

UIDEP := sfml

game/minima:\
	game/main.o\
	game/game.o\
	game/world.o\
	ui/ui.a
	@echo prog $@
	@$(CXX) -o $@ $(CXXFLAGS) $^ $(LDFLAGS)

game/world.o: game/world.hpp
game/game.o: game/game.hpp

ui/ui.a: ui/ui.o ui/impl_$(UIDEP).o
	@echo lib ui
	@ar rsc $@ $^

%.o: %.cc
	$(CXX) -c -o $@ $(CXXFLAGS) $*.cc

clean:
	rm -f $(shell find . -name \*.o)
	rm -f $(shell find . -name \*.a)
	rm -f game/minima