#ifndef _UI_HPP_
#define _UI_HPP_

#include <memory>
#include "fixed.hpp"

namespace ui{

class Ui;
class Img;
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

	// LoadImg returns an image pointer that has been
	// loaded from the given file path.  This pointer can
	// be used to draw to the window.
	virtual std::shared_ptr<Img> LoadImg(const char *path) = 0;

	// Draw draws the image to the back-buffer of the window.
	// This image will not appear until the Flip() method is caled.
	virtual void Draw(const Vec3&, std::shared_ptr<Img> img) = 0;

	// Flip swaps the back buffer with the screen buffer, effectively
	// displaying everything that has been drawn to the Ui.
	virtual void Flip() = 0;

	// Clear draws black over the entire screen.
	virtual void Clear() = 0;

	// Delay waits for the specified number of seconds before
	// returning.
	virtual void Delay(float sec) = 0;

	// PollEvent polls for events, returning true if the
	// event was filled in and false if there were no
	// events.
	virtual bool PollEvent(Event&) = 0;
};

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

	// The names of non-letter keys.  Just fill this
	// in as needed.
	enum {
		KeyReturn,
		KeySpace,
		KeyArrowUp,
		KeyArrowDown,
		KeyArrowLeft,
		KeyArrowRight,
		KeyEscape,
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
/*,
		// A KeyDown event indicates that a key has been
		// pressed.  The button field indicates name of the
		// key.  Letter keys are associated with their character
		// representation and other keys have specific symbols
		// that are defined in an enumeration.
		KeyDown,

		// A KeyUp event is the same as a KeyDown event
		// except that it indicates that a key has been released.
		KeyUp,
*/
	};

	Type type;
	int x, y;
	int button;
};

// OpenWindow returns a new Ui object.
std::unique_ptr<Ui> OpenWindow(Fixed w, Fixed h, const char *title);

// Img is the interface to a 2D image.
class Img{
public:
	virtual ~Img() = 0;
};

}	// namespace uikx

#endif	// _UI_HPP_