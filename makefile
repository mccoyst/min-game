OBJS:=\
	game/main.o\
	game/game.o\
	game/world.o\
	game/ui.o\

CXXFLAGS:=-std=c++0x
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

OBJS+=game/impl_sdl.o game/SDLMain.o

else

CXX:=g++
LDFLAGS+=-lSDL -lSDL_image -lSDL_ttf -lGLU -lGL
OBJS+=game/impl_sdl.o

CXXFLAGS+=-Wall -Werror -std=c++0x\
	-I/usr/include/SDL\

endif

all: wgen/wgen wimg/wimg game/minima

fetch:
	go get -v $(shell go list ./...)

game/minima: $(OBJS)
	@echo $@
	@$(CXX) -o $@ $^ $(LDFLAGS)

wgen/wgen: wgen/*.go
	go build -o wgen/wgen ./wgen

wimg/wimg: wimg/*.go
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
