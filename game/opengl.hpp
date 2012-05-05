#pragma once
#include "ui.hpp"
#include <SDL_opengl.h>
#include <vector>

class OpenGLUi : public Ui {
public:
	OpenGLUi(Fixed w, Fixed h) : Ui(w, h) { }

	// InitOpenGL initializes OpenGL for drawing stuff.
	void InitOpenGL();

	void Clear();

	virtual void DrawLine(const Vec2&, const Vec2&, const Color&);
	virtual void FillRect(const Vec2&, const Vec2&, const Color&);
	virtual void DrawRect(const Vec2&, const Vec2&, const Color&);
	virtual void Draw(const Vec2&, std::shared_ptr<Img> img, float shade = 1);

	virtual void InitTiles(int w, int h, int tw, int th, std::shared_ptr<Img>);

	virtual void SetTile(int x, int y, int tile, float shade) {
		tiles.at(x * sheeth + y) = tile;
		shades.at(x * sheeth + y) = shade;
	}

	virtual void DrawTiles(const Vec2&);

private:
	std::shared_ptr<Img> tileImgs;
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
