#ifndef _UI_HPP_
#define _UI_HPP_

#include <memory>

namespace ui{

class Screen;
class ScreenStk;
class Ui;
class Control;
class Img;
struct Event;

// Len is a fixed-point numeric type, scaled by Len::Scale;
class Len{
	int n;
public:
	enum{ Scale = 16 };

	explicit Len(int n);
	Len(int n, int frac);
	Len(const Len&);

	Len &operator = (const Len&);
	Len &operator += (Len);
	Len &operator -= (Len);
	Len &operator *= (Len);
	Len &operator /= (Len);
	Len &operator %= (Len);

	int value() const;
	int whole() const;
};

Len operator + (Len, Len);
Len operator - (Len, Len);
Len operator * (Len, Len);
Len operator / (Len, Len);
Len operator % (Len, Len);
bool operator == (Len, Len);
bool operator != (Len, Len);
bool operator < (Len, Len);
bool operator <= (Len, Len);
bool operator > (Len, Len);
bool operator >= (Len, Len);

// Vec is a trio of Lens, representing a vector in 2D coords.
struct Vec3{
	Len x, y, z;

	Vec3(Len x, Len y, Len z);
	Vec3(const Vec3 &);

	Vec3 &operator = (const Vec3&);
	Vec3 &operator += (Vec3);
	Vec3 &operator -= (Vec3);
	Vec3 &operator *= (Len);
};

Vec3 operator + (Vec3, Vec3);
Vec3 operator - (Vec3, Vec3);
Vec3 operator * (Vec3, Len);
bool operator == (Vec3, Vec3);
bool operator != (Vec3, Vec3);

// Ui is the interface to a user interface window, including graphics,
// device input, and sound.
class Ui{
public:
	// width and height are the dimensions of the window.
	Len width, height;

	// Ui constructs a new user interface that consists
	// of a window with the given width and height.
	Ui(const Len &w, const Len &h) : width(w), height(h) { }

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
std::unique_ptr<Ui> OpenWindow(Len w, Len h, const char *title);

// Img is the interface to a 2D image.
class Img{
public:
	virtual ~Img() = 0;
};

// Template and other inline definitions reside below.

inline Len::Len(int n)
	: n(n*Scale){
}

inline Len::Len(int n, int frac)
	: n(n*Scale + frac){
}

inline Len::Len(const Len &b)
	: n(b.n){
}

inline Len &Len::operator = (const Len &b){
	this->n = b.n;
	return *this;
}

inline Len &Len::operator += (Len n){
	this->n += n.n;
	return *this;
}

inline Len &Len::operator -= (Len n){
	this->n -= n.n;
	return *this;
}

inline Len &Len::operator *= (Len n){
	auto m = static_cast<long long>(this->n) * n.n;
	m /= Scale;
	this->n *= m;
	return *this;
}

inline Len &Len::operator /= (Len n){
	auto m = static_cast<long long>(this->n) * Scale;
	m /= n.n;
	this->n = m;
	return *this;
}

inline Len &Len::operator %= (Len n){
	this->n %= n.n;
	this->n *= Scale;
	return *this;
}

inline int Len::value() const{
	return this->n;
}

inline int Len::whole() const{
	return this->n / Scale;
}

inline Len operator + (Len a, Len b){
	return a += b;
}

inline Len operator - (Len a, Len b){
	return a -= b;
}

inline Len operator * (Len a, Len b){
	return a *= b;
}

inline Len operator / (Len a, Len b){
	return a /= b;
}

inline Len operator % (Len a, Len b){
	return a %= b;
}

inline bool operator == (Len a, Len b){
	return a.value() == b.value();
}

inline bool operator != (Len a, Len b){
	return a.value() != b.value();
}

inline bool operator < (Len a, Len b){
	return a.value() < b.value();
}

inline bool operator <= (Len a, Len b){
	return a.value() <= b.value();
}

inline bool operator > (Len a, Len b){
	return a.value() > b.value();
}

inline bool operator >= (Len a, Len b){
	return a.value() >= b.value();
}

inline Vec3::Vec3(Len x, Len y, Len z)
	: x(x), y(y), z(z){
}

inline Vec3::Vec3(const Vec3 &v)
	: x(v.x), y(v.y), z(v.z){
}

inline Vec3 &Vec3::operator = (const Vec3 &v){
	this->x = v.x;
	this->y = v.y;
	this->z = v.z;
	return *this;
}

inline Vec3 &Vec3::operator += (Vec3 v){
	this->x += v.x;
	this->y += v.y;
	this->z += v.z;
	return *this;
}

inline Vec3 &Vec3::operator -= (Vec3 v){
	this->x -= v.x;
	this->y -= v.y;
	this->z -= v.z;
	return *this;
}

inline Vec3 &Vec3::operator *= (Len n){
	this->x *= n;
	this->y *= n;
	this->z *= n;
	return *this;
}

inline Vec3 operator + (Vec3 a, Vec3 b){
	return a += b;
}

inline Vec3 operator - (Vec3 a, Vec3 b){
	return a -= b;
}

inline Vec3 operator * (Vec3 v, Len n){
	return v *= n;
}

inline bool operator == (Vec3 a, Vec3 b){
	return a.x == b.x && a.y == b.y && a.z == b.z;
}

inline bool operator != (Vec3 a, Vec3 b){
	return !(a == b);
}

}

#endif	// _UI_HPP_