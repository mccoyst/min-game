// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
#pragma once

// Fixed is a fixed-point numeric type, scaled by Fixed::Scale;
class Fixed{
	int n;
public:
	enum{ Scale = 16 };

	Fixed();
	explicit Fixed(int n);
	Fixed(int n, int frac);
	Fixed(const Fixed&);

	// Between returns true if the 3rd argument is
	// between the first two.
	static bool Between(Fixed, Fixed, Fixed);

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

// Vec is a tuple of Fixeds, representing a vector in 2D coords.
class Vec2{
public:
	Fixed x, y;

	Vec2();
	Vec2(Fixed x, Fixed y);
	Vec2(const Vec2 &);

	Vec2 &operator = (const Vec2&);
	Vec2 &operator += (const Vec2&);
	Vec2 &operator -= (const Vec2&);
	Vec2 &operator *= (Fixed);
	Vec2 &operator /= (Fixed);

	static const Vec2 Zero;
};

Vec2 operator + (Vec2, const Vec2&);
Vec2 operator - (Vec2, const Vec2&);
Vec2 operator * (Vec2, Fixed);
Vec2 operator / (Vec2, Fixed);
bool operator == (const Vec2&, const Vec2&);
bool operator != (const Vec2&, const Vec2&);
bool operator < (const Vec2&, const Vec2&);
bool operator <= (const Vec2&, const Vec2&);
bool operator > (const Vec2&, const Vec2&);
bool operator >= (const Vec2&, const Vec2&);

// Template and other inline definitions reside below.

inline Fixed::Fixed(){
}

inline Fixed::Fixed(int n)
	: n(n*Scale){
}

inline Fixed::Fixed(int n, int frac)
	: n(n*Scale + frac){
}

inline Fixed::Fixed(const Fixed &b)
	: n(b.n){
}

static inline bool Between(Fixed min, Fixed max, Fixed x) {
	return x >= min && x <= max;
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

inline Vec2::Vec2(){
}

inline Vec2::Vec2(Fixed x, Fixed y)
	: x(x), y(y){
}

inline Vec2::Vec2(const Vec2 &v)
	: x(v.x), y(v.y){
}

inline Vec2 &Vec2::operator = (const Vec2 &v){
	this->x = v.x;
	this->y = v.y;
	return *this;
}

inline Vec2 &Vec2::operator += (const Vec2 &v){
	this->x += v.x;
	this->y += v.y;
	return *this;
}

inline Vec2 &Vec2::operator -= (const Vec2 &v){
	this->x -= v.x;
	this->y -= v.y;
	return *this;
}

inline Vec2 &Vec2::operator *= (Fixed n){
	this->x *= n;
	this->y *= n;
	return *this;
}

inline Vec2 &Vec2::operator /= (Fixed n){
	this->x /= n;
	this->y /= n;
	return *this;
}

inline Vec2 operator + (Vec2 a, const Vec2 &b){
	return a += b;
}

inline Vec2 operator - (Vec2 a, const Vec2 &b){
	return a -= b;
}

inline Vec2 operator * (Vec2 v, Fixed n){
	return v *= n;
}

inline Vec2 operator / (Vec2 v, Fixed n){
	return v /= n;
}

inline bool operator == (const Vec2 &a, const Vec2 &b){
	return a.x == b.x && a.y == b.y;
}

inline bool operator != (const Vec2 &a, const Vec2 &b){
	return !(a == b);
}

inline bool operator < (const Vec2 &a, const Vec2& b) {
	return a.x < b.x && a.y < b.y;
}

inline bool operator <= (const Vec2 &a, const Vec2& b) {
	return a.x <= b.x && a.y <= b.y;
}

inline bool operator > (const Vec2 &a, const Vec2& b) {
	return a.x > b.x && a.y > b.y;
}

inline bool operator >= (const Vec2 &a, const Vec2& b) {
	return a.x >= b.x && a.y >= b.y;
}

