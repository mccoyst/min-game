OBJS:=\
	game/main.o\
	game/game.o\
	game/world.o\
	game/ui.o\
	game/opengl.o\
	game/impl_sdl.o\

CXXFLAGS:=-std=c++0x -Wall -O3
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

all: wgen/wgen wimg/wimg game/minima

fetch:
	go get -v $(shell go list ./...)

game/minima: $(OBJS)
	@echo $@
	@$(CXX) -o $@ $^ $(LDFLAGS)

wgen/wgen: wgen/*.go world/*.go
	go build -o wgen/wgen ./wgen

wimg/wimg: wimg/*.go world/*.go
	go build -o wimg/wimg ./wimg

include $(OBJS:.o=.d)

%.d: %.cc
	@echo $@
	@./dep.sh $(CXX) $(shell dirname $<) $(HEADERFLAGS) $< > $@

%.d: %.m
	@echo $@
	@./dep.sh $(OBJCC) $(shell dirname $<) $(HEADERFLAGS) $< > $@

%.o: %.cc
	@echo $@
	@$(CXX) -c -o $@ $(CXXFLAGS) $<

%.o: %.m
	@echo $@
	@$(OBJCC) -c -o $@ $(OBJCFLAGS) $<

clean:
	rm -f $(OBJS) game/minima wgen/wgen wimg/wimg

nuke: clean
	rm -f $(shell find . -not -iwholename \*.hg\* -name \*.d)
