#ifndef _LEN_HPP_
#define _LEN_HPP_

// Fixed is a fixed-point numeric type, scaled by Fixed::Scale;
class Fixed{
	int n;
public:
	enum{ Scale = 16 };

	explicit Fixed(int n);
	Fixed(int n, int frac);
	Fixed(const Fixed&);

	Fixed &operator = (Fixed);
	Fixed &operator += (Fixed);
	Fixed &operator -= (Fixed);
	Fixed &operator *= (Fixed);
	Fixed &operator /= (Fixed);
	Fixed &operator %= (Fixed);

	int value() const;
	int whole() const;
};

Fixed operator - (Fixed);
Fixed operator + (Fixed, Fixed);
Fixed operator - (Fixed, Fixed);
Fixed operator * (Fixed, Fixed);
Fixed operator / (Fixed, Fixed);
Fixed operator % (Fixed, Fixed);
bool operator == (Fixed, Fixed);
bool operator != (Fixed, Fixed);
bool operator < (Fixed, Fixed);
bool operator <= (Fixed, Fixed);
bool operator > (Fixed, Fixed);
bool operator >= (Fixed, Fixed);

// Vec is a trio of Fixeds, representing a vector in 2D coords.
struct Vec3{
	Fixed x, y, z;

	Vec3(Fixed x, Fixed y, Fixed z);
	Vec3(const Vec3 &);

	Vec3 &operator = (const Vec3&);
	Vec3 &operator += (const Vec3&);
	Vec3 &operator -= (const Vec3&);
	Vec3 &operator *= (Fixed);
	Vec3 &operator /= (Fixed);
};

Vec3 operator + (Vec3, const Vec3&);
Vec3 operator - (Vec3, const Vec3&);
Vec3 operator * (Vec3, Fixed);
Vec3 operator / (Vec3, Fixed);
bool operator == (const Vec3&, const Vec3&);
bool operator != (const Vec3&, const Vec3&);

// Template and other inline definitions reside below.

inline Fixed::Fixed(int n)
	: n(n*Scale){
}

inline Fixed::Fixed(int n, int frac)
	: n(n*Scale + frac){
}

inline Fixed::Fixed(const Fixed &b)
	: n(b.n){
}

inline Fixed &Fixed::operator = (Fixed b){
	this->n = b.n;
	return *this;
}

inline Fixed &Fixed::operator += (Fixed n){
	this->n += n.n;
	return *this;
}

inline Fixed &Fixed::operator -= (Fixed n){
	this->n -= n.n;
	return *this;
}

inline Fixed &Fixed::operator *= (Fixed n){
	long long m = this->n;
	m *= n.n;
	m /= Scale;
	this->n = m;
	return *this;
}

inline Fixed &Fixed::operator /= (Fixed n){
	long long m = this->n;
	m /= n.n;
	this->n = m * Scale;
	return *this;
}

inline Fixed &Fixed::operator %= (Fixed n){
	this->n %= n.n;
	return *this;
}

inline int Fixed::value() const{
	return this->n;
}

inline int Fixed::whole() const{
	return this->n / Scale;
}

inline Fixed operator - (Fixed a){
	return Fixed(0) - a;
}

inline Fixed operator + (Fixed a, Fixed b){
	return a += b;
}

inline Fixed operator - (Fixed a, Fixed b){
	return a -= b;
}

inline Fixed operator * (Fixed a, Fixed b){
	return a *= b;
}

inline Fixed operator / (Fixed a, Fixed b){
	return a /= b;
}

inline Fixed operator % (Fixed a, Fixed b){
	return a %= b;
}

inline bool operator == (Fixed a, Fixed b){
	return a.value() == b.value();
}

inline bool operator != (Fixed a, Fixed b){
	return a.value() != b.value();
}

inline bool operator < (Fixed a, Fixed b){
	return a.value() < b.value();
}

inline bool operator <= (Fixed a, Fixed b){
	return a.value() <= b.value();
}

inline bool operator > (Fixed a, Fixed b){
	return a.value() > b.value();
}

inline bool operator >= (Fixed a, Fixed b){
	return a.value() >= b.value();
}

inline Vec3::Vec3(Fixed x, Fixed y, Fixed z = Fixed(0))
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

inline Vec3 &Vec3::operator += (const Vec3 &v){
	this->x += v.x;
	this->y += v.y;
	this->z += v.z;
	return *this;
}

inline Vec3 &Vec3::operator -= (const Vec3 &v){
	this->x -= v.x;
	this->y -= v.y;
	this->z -= v.z;
	return *this;
}

inline Vec3 &Vec3::operator *= (Fixed n){
	this->x *= n;
	this->y *= n;
	this->z *= n;
	return *this;
}

inline Vec3 &Vec3::operator /= (Fixed n){
	this->x /= n;
	this->y /= n;
	this->z /= n;
	return *this;
}

inline Vec3 operator + (Vec3 a, const Vec3 &b){
	return a += b;
}

inline Vec3 operator - (Vec3 a, const Vec3 &b){
	return a -= b;
}

inline Vec3 operator * (Vec3 v, Fixed n){
	return v *= n;
}

inline Vec3 operator / (Vec3 v, Fixed n){
	return v /= n;
}

inline bool operator == (const Vec3 &a, const Vec3 &b){
	return a.x == b.x && a.y == b.y && a.z == b.z;
}

inline bool operator != (const Vec3 &a, const Vec3 &b){
	return !(a == b);
}

#endif	// _LEN_HPP_