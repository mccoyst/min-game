#pragma once
#include "ui.hpp"
#include <SDL_opengl.h>
#include <vector>

struct OpenGLUi : public Ui {

	OpenGLUi(Fixed w, Fixed h) : Ui(w, h) { }

	// InitOpenGL initializes OpenGL for drawing stuff.
	void InitOpenGL();

	// Clear clears the screen.
	void Clear();

	// Draw draws the image to the back-buffer of the window.
	// This image will not appear until the Flip() method is called.
	// The shade argument is an alpha value between 0 (dark) and
	// 1 (light).
	virtual void Draw(const Vec2&, std::shared_ptr<Img> img, float shade = 1);

	// InitTiles initializes the tiles.
	virtual void InitTiles(int w, int h, int tw, int th, std::shared_ptr<Img> img) {
		tileImgs = img;
		sheetw = w;
		sheeth = h;
		tilew = tw;
		tileh = th;
		tiles.resize(w*h);
		shades.resize(w*h);
	}

	// SetTile sets the tile image number for the given tile.
	virtual void SetTile(int x, int y, int tile, float shade) {
		tiles.at(x * sheeth + y) = tile;
		shades.at(x * sheeth + y) = shade;
	}

	// DrawTiles draws the tiles at the given offset.
	virtual void DrawTiles(const Vec2&);

private:
	std::shared_ptr<Img> tileImgs;
	int sheetw, sheeth, tilew, tileh;

	// Tile associated with each x,y.
	std::vector<int> tiles;

	// The shade associated with each x,y.
	std::vector<float> shades;
};

struct OpenGLImg : public Img {
	GLuint texid;
	virtual ~OpenGLImg();
};
