// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
#pragma once

#include "fixed.hpp"

using std::string;
using std::unique_ptr;

class Ui;
class Img;
struct Color;
struct Event;
class TileView;
class Camera;

// Ui is the interface to a user interface window, including graphics,
// device input, and sound.
class Ui{
public:
	// width and height are the dimensions of the window.
	Fixed width, height;

	virtual ~Ui();

	// DrawLine draws the given line on the screen.
	//
	// This is probably slow so just use it for debugging stuff.
	virtual void DrawLine(const Vec2&, const Vec2&, const Color&) = 0;

	// FillRect fills the given rectangle (specified by lower left
	// vertex, and width/height) with some color.
	//
	// This is probably slow so just use it for debugging stuff.
	virtual void FillRect(const Vec2&, const Vec2&, const Color&) = 0;

	// DrawRect outlines the given rectangle (specified by lower left
	// vertex, and width/height) with some color.
	//
	// This is probably slow so just use it for debugging stuff.
	virtual void DrawRect(const Vec2&, const Vec2&, const Color&) = 0;

	// Draw draws the image to the back-buffer of the window.
	// This image will not appear until the Flip() method is called.
	// The shade argument is an alpha value between 0 (dark) and
	// 1 (light).
	virtual void Draw(const Vec2&, Img&, float shade = 1) = 0;

	// Draw draws the tiles at the given offset.
	virtual void Draw(const Vec2&, const TileView&) = 0;

	// Flip swaps the back buffer with the screen buffer, effectively
	// displaying everything that has been drawn to the Ui.
	virtual void Flip() = 0;

	// Clear draws black over the entire screen.
	virtual void Clear() = 0;

	// Delay waits for the specified number of milli-seconds
	// before returning.
	virtual void Delay(unsigned long msec) = 0;

	// Ticks returns the number of milliseconds since the Ui
	// was created.
	virtual unsigned long Ticks() = 0;

	// PollEvent polls for events, returning true if the
	// event was filled in and false if there were no
	// events.
	virtual bool PollEvent(Event&) = 0;
};

// NewUi constructs a real user interface that consists
// of a window with the given width and height.
unique_ptr<Ui> NewUi(Fixed w, Fixed h, const string &title);

class Camera{
	Vec2 cam;
public:
	Camera() = default;

	// MoveCam adds v to the camera's current position.
	void Move(Vec2 v);

	// CenterCam moves the camera so that v is in the center.
	void Center(Vec2 v);

	// CamPos returns the camera's current position.
	Vec2 Pos() const;

	// DrawCam draws the image from the camera's point of view.
	void Draw(Vec2, Ui&, Img&, float shade = 1);
};

class TileView{
public:
	virtual ~TileView();

	// SetTile sets the tile image and shade for the given tile.
	virtual void SetTile(int x, int y, int tile, float shade) = 0;
};

unique_ptr<TileView> NewTileView(int w, int h, int tw, int th, unique_ptr<Img> &&img);

struct Color {
	constexpr Color(unsigned char r, unsigned char g, unsigned char b,
		unsigned char a = 255) : r(r), g(g), b(b), a(a) { }
	unsigned char r, g, b, a;
};

constexpr Color White{ 255, 255, 255 };
constexpr Color Gray{ 128, 128, 128 };
constexpr Color Black{ 0, 0, 0 };
constexpr Color Red{ 255, 0, 0 };
constexpr Color Green{ 0, 255, 0 };
constexpr Color Blue{ 0, 0, 255 };

// Img is the interface to a 2D image.
class Img{
public:
	virtual ~Img() = 0;
	virtual Vec2 Size() const = 0;
};

// LoadImg returns an image pointer that has been
// loaded from the given file path.  This pointer can
// be used to draw to the window.
// 
// This is a "low-level" function. Use FindImg instead.
unique_ptr<Img> LoadImg(const string &);

// FindImg returns an image pointer that has been
// loaded from the application's installed resources.
unique_ptr<Img> FindImg(const string &);

// Font describes a text font, color, etc.
class Font{
public:
	virtual ~Font() = 0;

	// Render renders the given text to an image using this font.
	virtual unique_ptr<Img> Render(const string&) = 0;
};

// LoadFont loads a font from a file with the given size and color.
// 
// This is a "low-level" function. Use FindFont instead.
unique_ptr<Font> LoadFont(const string&, int, Color);

// FindFont returns a Font pointer that has been
// loaded from the application's installed resources.
unique_ptr<Font> FindFont(const string &, int, Color);

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
		Action,

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
