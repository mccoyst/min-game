// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
#pragma once
#include "ui.hpp"
#include <SDL_opengl.h>
#include <vector>

class OpenGLUi {
public:
	OpenGLUi(Fixed, Fixed);

	void Clear();

	void DrawLine(const Vec2&, const Vec2&, const Color&);
	void FillRect(const Vec2&, const Vec2&, const Color&);
	void DrawRect(const Vec2&, const Vec2&, const Color&);
	void Draw(const Vec2&, Img*, float shade = 1);

	void InitTiles(int w, int h, int tw, int th, std::unique_ptr<Img>);

	void SetTile(int x, int y, int tile, float shade) {
		tiles.at(x * sheeth + y) = tile;
		shades.at(x * sheeth + y) = shade;
	}

	void DrawTiles(const Vec2&);

private:
	Img *tileImgs;
	int sheetw, sheeth, tilew, tileh;

	// Tile associated with each x,y.
	std::vector<int> tiles;

	// The shade associated with each x,y.
	std::vector<float> shades;
};

class OpenGLImg : public Img {
public:
	GLuint texid;
	virtual ~OpenGLImg();
};
