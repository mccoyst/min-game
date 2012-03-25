OBJS:=\
	game/main.o\
	game/game.o\
	game/world.o\
	game/ui.o\
	game/impl_sfml.o\

CXXFLAGS:=-std=c++0x
LDFLAGS:=

OS := $(shell uname | sed 's/.*MINGW.*/win/')

ifeq ($(OS),Darwin)
CXX:=clang++ -fno-color-diagnostics -stdlib=libc++
CXXFLAGS+=\
	-framework sfml-graphics\
	-framework sfml-windo\
	-framework sfml-system \

else
CXX:=g++
LDFLAGS+=-lsfml-graphics -lsfml-window -lsfml-system
CXXFLAGS+=-Wall -Werror -std=c++0x
endif

all: game/minima

game/minima: $(OBJS)
	@echo $@
	@$(CXX) -o $@ $(CXXFLAGS) $^ $(LDFLAGS)

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
	rm -f $(OBJS) game/minima

nuke: clean
	rm -f $(shell find . -not -iwholename \*.hg\* -name \*.d)