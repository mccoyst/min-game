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
CXX:=clang++ -fno-color-diagnostics -stdlib=libc++

CXXFLAGS+=\
	-I/Library/Frameworks/SDL.framework/Headers\
	-I/Library/Frameworks/SDL_image.framework/Headers\

LDFLAGS +=\
	-framework SDL\
	-framework SDL_image\
	-framework OpenGL\
	-framework Foundation\
	-framework Cocoa\

OBJCFLAGS := $(LDFLAGS)

OBJCC := clang

OBJS+=game/impl_sdl.o game/SDLMain.o

else

CXX:=g++
LDFLAGS+=-lSDL -lSDL_image -lGLU -lGL
OBJS+=game/impl_sdl.o

CXXFLAGS+=-Wall -Werror -std=c++0x\
	-I/usr/local/include/SDL\

endif

all: wgen/wgen wimg/wimg game/minima

fetch:
	go get -v $(shell go list ./...)

game/minima: $(OBJS)
	@echo $@
	$(CXX) -o $@ $(CXXFLAGS) $^ $(LDFLAGS)

wgen/wgen: wgen/*.go
	go build -o wgen/wgen ./wgen

wimg/wimg: wimg/*.go
	go build -o wimg/wimg ./wimg

include $(OBJS:.o=.d)

%.d: %.cc
	@echo $@
	./dep.sh g++ $(shell dirname $*) $(CXXFLAGS) $*.cc > $@

%.d: %.c
	@echo $@
	./dep.sh gcc $(shell dirname $*) $(CFLAGS) $*.c > $@

%.d: %.m
	@echo $@
	./dep.sh gcc $(shell dirname $*) $(CFLAGS) $*.m > $@

%.o: %.cc
	@echo $@
	$(CXX) -c -o $@ $(CXXFLAGS) $*.cc

%.o: %.m
	@echo $@
	$(OBJCC) -c -o $@ $(OBJCFLAGS) $*.m

clean:
	rm -f $(OBJS) game/minima wgen/wgen wimg/wimg

nuke: clean
	rm -f $(shell find . -not -iwholename \*.hg\* -name \*.d)
