# Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
OBJS:=\
	game/main.o\
	game/game.o\
	game/world.o\
	game/ui.o\
	game/screen.o\
	game/opengl.o\
	game/ui_sdl.o\
	game/io.o\
	game/title.o\

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

OBJS += game/SDLMain.o

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

game/minima: $(OBJS)
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

include $(OBJS:.o=.d)

%.d: %.cpp
	@echo $@
	@./dep.sh $(CXX) $(shell dirname $<) $(CXXFLAGS) $< > $@

%.d: %.m
	@echo $@
	@./dep.sh $(OBJCC) $(shell dirname $<) $(OBJCFLAGS) $< > $@

%.o: %.cpp
	@echo $@
	@$(CXX) -c -o $@ $(CXXFLAGS) $<

%.o: %.m
	@echo $@
	@$(OBJCC) -c -o $@ $(OBJCFLAGS) $<

clean:
	rm -f $(OBJS) $(TARGS)

nuke: clean
	rm -f $(shell find . -not -iwholename \*.hg\* -name \*.d)
