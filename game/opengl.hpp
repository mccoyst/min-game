#pragma once
#include "ui.hpp"
#include <SDL_opengl.h>

struct OpenGLUi : public Ui {

	OpenGLUi(Fixed w, Fixed h) : Ui(w, h) { }

	// InitOpenGL initializes OpenGL for drawing stuff.
	void InitOpenGL();

	// Draw draws the image to the back-buffer of the window.
	// This image will not appear until the Flip() method is called.
	// The shade argument is an alpha value between 0 (dark) and
	// 1 (light).
	virtual void Draw(const Vec2&, std::shared_ptr<Img> img, float shade = 1);
};

struct OpenGLImg : public Img {
	GLuint texid;
	virtual ~OpenGLImg();
};
