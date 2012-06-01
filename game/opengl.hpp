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
	void Draw(const Vec2&, const TileView&);
};

class OpenGLImg : public Img {
public:
	GLuint texid;
	virtual ~OpenGLImg();
};
