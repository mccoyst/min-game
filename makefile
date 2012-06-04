# Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
OBJS:=\
	main.o\
	game.o\
	world.o\
	ui.o\
	screen.o\
	opengl.o\
	ui_sdl.o\
	io.o\
	title.o\

TILES :=\
	Grass.png\
	Water.png\
	Mountain.png\
	Tree.png\
	Desert.png\
	Glacier.png\

CXXFLAGS:=-std=c++11 -Wall -O3
OBJCFLAGS :=
LDFLAGS:=

OS := $(shell uname | sed 's/.*MINGW.*/win/')

ifeq ($(OS),Darwin)
CXX:=clang++

HEADERFLAGS+=\
	-I/Library/Frameworks/SDL.framework/Headers\
	-I/Library/Frameworks/SDL_image.framework/Headers\
	-I/Library/Frameworks/SDL_ttf.framework/Headers\

LDFLAGS +=\
	-stdlib=libc++\
	-framework SDL\
	-framework SDL_image\
	-framework SDL_ttf\
	-framework OpenGL\
	-framework Foundation\
	-framework Cocoa\

CXXFLAGS += -fno-color-diagnostics -stdlib=libc++ $(HEADERFLAGS)

OBJCFLAGS := $(HEADERFLAGS)

OBJCC := clang

OBJS += SDLMain.o

else

CXX:=g++
LDFLAGS+=-lSDL -lSDL_image -lSDL_ttf -lGLU -lGL
HEADERFLAGS := -I/usr/include/SDL
CXXFLAGS+=-Werror -DGL_GLEXT_PROTOTYPES $(HEADERFLAGS)

endif

TARGS :=\
	wgen/wgen\
	wimg/wimg\
	game/minima\
	mksheet/mksheet\
	resrc/tiles.png\

all: $(TARGS)

fetch:
	go get -v $(shell go list ./...)

game/minima: $(OBJS:%=_work/%)
	@echo $@
	@$(CXX) -o $@ $^ $(LDFLAGS)

wgen/wgen: wgen/*.go world/*.go
	go install ./wgen && touch $@

wimg/wimg: wimg/*.go world/*.go
	go install ./wimg && touch $@

mksheet/mksheet: mksheet/*.go
	go install ./mksheet && touch $@

resrc/tiles.png: $(TILES:%=resrc/%)
	mksheet $^ $@

include $(OBJS:%.o=_work/%.d)

_work/%.d: game/%.cpp
	@echo $@
	@./dep.sh $(CXX) $(shell dirname $<) $(CXXFLAGS) $< > $@

_work/%.d: game/%.m
	@echo $@
	@./dep.sh $(OBJCC) $(shell dirname $<) $(OBJCFLAGS) $< > $@

_work/%.o: game/%.cpp
	@echo $@
	@$(CXX) -c -o $@ $(CXXFLAGS) $<

_work/%.o: game/%.m
	@echo $@
	@$(OBJCC) -c -o $@ $(OBJCFLAGS) $<

clean:
	rm -f _work/*.d
	rm -f _work/*.o

nuke: clean
	rm -f $(TARGS)
