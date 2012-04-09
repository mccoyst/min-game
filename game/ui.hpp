#pragma once

#include <memory>
#include "fixed.hpp"

class Ui;
class Img;
struct Color;
struct Event;

// Ui is the interface to a user interface window, including graphics,
// device input, and sound.
class Ui{
public:
	// width and height are the dimensions of the window.
	Fixed width, height;

	// Ui constructs a new user interface that consists
	// of a window with the given width and height.
	Ui(const Fixed &w, const Fixed &h) : width(w), height(h) { }

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
	virtual void Draw(const Vec2&, std::shared_ptr<Img> img, float shade = 1) = 0;

	// InitTiles initializes the tiles sheet.
	virtual void InitTiles(int, int, int, int, std::shared_ptr<Img>) = 0;

	// SetTile sets the tile image and shade for the given tile.
	virtual void SetTile(int, int, int, float) = 0;

	// DrawTiles draws the tiles at the given offset.
	virtual void DrawTiles(const Vec2&) = 0;

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

// OpenWindow returns a new Ui object.
std::shared_ptr<Ui> OpenWindow(Fixed w, Fixed h, const char *title);

struct Color {
	Color(unsigned char _r, unsigned char _g, unsigned char _b,
		unsigned char _a = 255) : r(_r), g(_g), b(_b), a(_a) { }
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
std::shared_ptr<Img> LoadImg(const char*);

// Font describes a text font, color, etc.
struct Font{
	virtual ~Font() = 0;

	// Render renders the given text to an image using this font.
	virtual std::shared_ptr<Img> Render(const char*, ...) = 0;
};

// LoadFont loads a font from a file with the given size and color.
std::shared_ptr<Font> LoadFont(const char*, int, char, char, char);

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

	// The names of non-letter keys.  Just fill this in as needed.
	enum {
		KeyUpArrow,
		KeyDownArrow,
		KeyLeftArrow,
		KeyRightArrow,
		KeyLShift,
		KeyRShift,
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

		// A KeyUp event is the same as a KeyDown event
		// except that it indicates that a key has been released.
		KeyUp,
	};

	Type type;
	int x, y;
	int button;
};
