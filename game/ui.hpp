// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
#pragma once

#include "fixed.hpp"
#include <memory>

class Ui;
class Img;
struct Color;
struct Event;

// Ui is the interface to a user interface window, including graphics,
// device input, and sound.
class Ui{
public:
	struct Impl;
	std::unique_ptr<Impl> impl;

	// width and height are the dimensions of the window.
	Fixed width, height;

	// Ui constructs a new user interface that consists
	// of a window with the given width and height.
	Ui(Fixed w, Fixed h, const char *title);
	~Ui();

	// DrawLine draws the given line on the screen.
	//
	// This is probably slow so just use it for debugging stuff.
	void DrawLine(const Vec2&, const Vec2&, const Color&);

	// FillRect fills the given rectangle (specified by lower left
	// vertex, and width/height) with some color.
	//
	// This is probably slow so just use it for debugging stuff.
	void FillRect(const Vec2&, const Vec2&, const Color&);

	// DrawRect outlines the given rectangle (specified by lower left
	// vertex, and width/height) with some color.
	//
	// This is probably slow so just use it for debugging stuff.
	void DrawRect(const Vec2&, const Vec2&, const Color&);

	// Draw draws the image to the back-buffer of the window.
	// This image will not appear until the Flip() method is called.
	// The shade argument is an alpha value between 0 (dark) and
	// 1 (light).
	void Draw(const Vec2&, Img*, float shade = 1);

	// InitTiles initializes the tiles sheet.
	void InitTiles(int w, int h, int tw, int th, std::unique_ptr<Img>&&);

	// SetTile sets the tile image and shade for the given tile.
	void SetTile(int x, int y, int tile, float shade);

	// DrawTiles draws the tiles at the given offset.
	void DrawTiles(const Vec2&);

	// MoveCam adds v to the camera's current position.
	void MoveCam(Vec2 v);

	// CamPos returns the camera's current position.
	Vec2 CamPos() const;

	// DrawCam draws the image from the camera's point of view.
	void DrawCam(Vec2, Img*, float shade = 1);

	// DrawTilesCam draws the tiles from the camera's point of view.
	void DrawTilesCam(Vec2);

	// Flip swaps the back buffer with the screen buffer, effectively
	// displaying everything that has been drawn to the Ui.
	void Flip();

	// Clear draws black over the entire screen.
	void Clear();

	// Delay waits for the specified number of milli-seconds
	// before returning.
	void Delay(unsigned long msec);

	// Ticks returns the number of milliseconds since the Ui
	// was created.
	unsigned long Ticks();

	// PollEvent polls for events, returning true if the
	// event was filled in and false if there were no
	// events.
	bool PollEvent(Event&);
};

struct Color {
	Color(unsigned char r, unsigned char g, unsigned char b,
		unsigned char a = 255) : r(r), g(g), b(b), a(a) { }
	unsigned char r, g, b, a;
};

// Img is the interface to a 2D image.
class Img{
public:
	virtual ~Img() = 0;
	virtual Vec2 Size() const = 0;
};

// LoadImg returns an image pointer that has been
// loaded from the given file path.  This pointer can
// be used to draw to the window.
Img *LoadImg(const char*);

// Font describes a text font, color, etc.
struct Font{
	virtual ~Font() = 0;

	// Render renders the given text to an image using this font.
	virtual Img *Render(const char*, ...) = 0;
};

// LoadFont loads a font from a file with the given size and color.
Font *LoadFont(const char*, int, char, char, char);

// An Event is a user input event handed back from the
// Ui's PollEvent method.
//
// Check the type field to see what event type has been
// returned and then based on the type you may use the
// information stored in some of the other fields.
struct Event {

	// The names of mouse buttons.
	enum {
		MouseLeft,
		MouseRight,
		MouseCenter,
	};

	enum {
		None,
		UpArrow,
		DownArrow,
		LeftArrow,
		RightArrow,
		LShift,
		RShift,

		NumKeys,
	};

	enum Type {
		// A Closed event indicates that the window
		// has been closed.  This event has no other
		// information.
		Closed,

		// A MouseMoved event indicates that the
		// mouse has moved.  The x and y fields
		// indicate the new mouse position within
		// the window.
		MouseMoved,

		// A MouseDown event indicates that a mouse
		// button has been pressed. The button field
		// contains the name of the mouse button and
		// the x and y fields have the location of the press.
		MouseDown,

		// A MouseUp event is the same as a MouseDown
		// event except that it indicates that the button
		// has been released.
		MouseUp,

		// A KeyDown event indicates that a key has been
		// pressed.  The button field indicates name of the
		// key.  Letter keys are associated with their character
		// representation and other keys have specific symbols
		// that are defined in an enumeration.
		KeyDown,

		// fakes a key down because a key has been held.
		SimulatedKeyDown,

		// A KeyUp event is the same as a KeyDown event
		// except that it indicates that a key has been released.
		KeyUp,


	};

	Type type;
	int x, y;
	int button;
};
