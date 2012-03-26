OBJS:=\
	game/main.o\
	game/game.o\
	game/world.o\
	game/ui.o\

CXXFLAGS:=-std=c++0x
LDFLAGS:=

OS := $(shell uname | sed 's/.*MINGW.*/win/')

ifeq ($(OS),Darwin)
CXX:=clang++ -fno-color-diagnostics -stdlib=libc++
CXXFLAGS+=\
	-framework sfml-graphics\
	-framework sfml-window\
	-framework sfml-system \

OBJS+=game/impl_sfml.o

else

CXX:=g++
#LDFLAGS+=-lsfml-graphics -lsfml-window -lsfml-system
#OBJS+=game/impl_sfml.o

LDFLAGS+=-lSDL -lSDL_image -lGLU -lGL
OBJS+=game/impl_sdl.o

CXXFLAGS+=-Wall -Werror -std=c++0x

endif

all: wgen/wgen wimg/wimg game/minima

fetch:
	go get -v $(shell go list ./...)

game/minima: $(OBJS)
	@echo $@
	@$(CXX) -o $@ $(CXXFLAGS) $^ $(LDFLAGS)

wgen/wgen: wgen/*.go
	go build -o wgen/wgen ./wgen

wimg/wimg: wimg/*.go
	go build -o wimg/wimg ./wimg

include $(OBJS:.o=.d)

%.d: %.cc
	@echo $@
	@./dep.sh g++ $(shell dirname $*) $(CXXFLAGS) $*.cc > $@

%.d: %.c
	@echo $@
	@./dep.sh gcc $(shell dirname $*) $(CFLAGS) $*.c > $@

%.o: %.cc
	@echo $@
	@$(CXX) -c -o $@ $(CXXFLAGS) $*.cc

clean:
	rm -f $(OBJS) game/minima wgen/wgen wimg/wimg

nuke: clean
	rm -f $(shell find . -not -iwholename \*.hg\* -name \*.d)