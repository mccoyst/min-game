# Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

OS := $(shell uname | sed 's/.*MINGW.*/win/')

TARGS :=\
	wgen/wgen\
	wimg/wimg\
	game/minima\
	mksheet/mksheet\
	resrc/tiles.png\

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
	popen.o\
	player.o\
	geom.o\
	resrc.o\

TILES :=\
	Grass.png\
	Water.png\
	Mountain.png\
	Tree.png\
	Desert.png\
	Glacier.png\

CXXFLAGS :=\
	-std=c++0x\
	-Wall\
	-O3\

LDFLAGS :=\

ifeq ($(OS),Darwin)

CXX := clang++
OBJCC := clang

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

CXXFLAGS +=\
	-fno-color-diagnostics\
	-stdlib=libc++\
	$(HEADERFLAGS)\

OBJCFLAGS :=\
	-Wall\
	-O3\
	-Wno-objc-protocol-method-implementation\
	$(HEADERFLAGS)

OBJS += SDLMain.o

else

CXX := g++

LDFLAGS +=\
	-lSDL\
	-lSDL_image\
	-lSDL_ttf\
	-lGLU\
	-lGL\

HEADERFLAGS :=\
	-I/usr/include/SDL\

CXXFLAGS +=\
	-Werror\
	-DGL_GLEXT_PROTOTYPES\
	$(HEADERFLAGS)\

endif



all: test

fetch:
	go get -v -u github.com/mccoyst/runt
	go get -v -u $(shell go list ./...)

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

.PHONY: test clean nuke

clean:
	rm -f _work/*.d
	rm -f _work/*.o

nuke: clean
	rm -f $(TARGS)

test: $(TARGS)
	@runt "-cxx=$(CXX)"\
		"-cxxflags=$(CXXFLAGS)"\
		"-ldflags=$(LDFLAGS)"\
		-testdir=game\
		$(filter-out _work/main.o _work/SDLMain.o,$(OBJS:%=_work/%))
