#include <memory>

namespace ui{

class Screen;
class ScreenStk;
class Ui;
class Control;
class Img;

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
	virtual ~Ui();

	virtual void Flip() = 0;
	virtual void Clear() = 0;

	virtual std::shared_ptr<Img> LoadImg(const char *path) = 0;
	virtual void Draw(std::shared_ptr<Img> img) = 0;
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
	this->n *= n.n;
	return *this;
}

inline Len &Len::operator /= (Len n){
	this->n /= n.n;
	return *this;
}

inline Len &Len::operator %= (Len n){
	this->n %= n.n;
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
	return a += b;
}

inline Len operator * (Len a, Len b){
	return a += b;
}

inline Len operator / (Len a, Len b){
	return a += b;
}

inline Len operator % (Len a, Len b){
	return a += b;
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
