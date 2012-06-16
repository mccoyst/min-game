// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

#include "game.hpp"
#include "ui.hpp"
#include <array>
#include <functional>

namespace{
	std::array<const string, 2> roots = {{
		"resrc",
		"../Resources",
	}};

	template<class T>
	unique_ptr<T> find(const string&, std::function<unique_ptr<T>(string)>);
}

unique_ptr<Img> FindImg(const string &name){
	return find<Img>(name, LoadImg);
}

unique_ptr<Font> FindFont(const string &name, int size, Color c){
	return find<Font>(name, [size,c](const string &path){
		return LoadFont(path, size, c);
	});
}

namespace{

template<class T>
unique_ptr<T> 
find(const string &name, std::function<unique_ptr<T>(string)> load){
	//TODO: cache these
	for(auto root : roots)
		try{
			return load(root + "/" + name);
		}catch(const Failure &f){
			// Okay, try the next root
		}
	throw Failure("Failed to find resource \"" + name + "\"");
}

}
