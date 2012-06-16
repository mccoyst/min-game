// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
#pragma once

// Fixed is a fixed-point numeric type, scaled by Fixed::Scale;
class Fixed{
	int n;
public:
	enum{ Scale = 16 };

	constexpr Fixed();
	constexpr explicit Fixed(int n);
	constexpr Fixed(int n, int frac);

	Fixed &operator = (Fixed);
	Fixed &operator += (Fixed);
	Fixed &operator -= (Fixed);
	Fixed &operator *= (Fixed);
	Fixed &operator /= (Fixed);
	Fixed &operator %= (Fixed);
	Fixed &operator ++ ();
	Fixed &operator -- ();

	constexpr int value() const;
	constexpr int whole() const;
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

Fixed Trunc(Fixed);
Fixed Round(Fixed);
Fixed Floor(Fixed);

// Vec is a tuple of Fixeds, representing a vector in 2D coords.
class Vec2{
public:
	Fixed x, y;

	constexpr Vec2();
	constexpr Vec2(Fixed x, Fixed y);

	Vec2 &operator = (const Vec2&);
	Vec2 &operator += (const Vec2&);
	Vec2 &operator -= (const Vec2&);
	Vec2 &operator *= (Fixed);
	Vec2 &operator /= (Fixed);
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

constexpr Fixed::Fixed()
	: n(0){
}

constexpr Fixed::Fixed(int n)
	: n(n*Scale){
}

constexpr Fixed::Fixed(int n, int frac)
	: n(n*Scale + frac){
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
	this->n = m / Scale;
	return *this;
}

inline Fixed &Fixed::operator /= (Fixed n){
	long long m = this->n;
	m *= Scale;
	this->n = m / n.n;
	return *this;
}

inline Fixed &Fixed::operator %= (Fixed n){
	this->n %= n.n;
	return *this;
}

inline Fixed &Fixed::operator ++ (){
	return *this += Fixed(1);
}

inline Fixed &Fixed::operator -- (){
	return *this -= Fixed(1);
}

constexpr int Fixed::value() const{
	return this->n;
}

constexpr int Fixed::whole() const{
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

inline Fixed Trunc(Fixed f){
	return Fixed{f.whole()};
}

inline Fixed Round(Fixed f){
	return Trunc(f + Fixed{0,Fixed::Scale/2});
}

inline Fixed Floor(Fixed f){
	if (f >= Fixed{0})
		return Trunc(f);
	if (Fixed{f.whole()} == f)
		return f;
	return Trunc(f) - Fixed{1};
}

constexpr Vec2::Vec2()
	: x(0), y(0){
}

constexpr Vec2::Vec2(Fixed x, Fixed y)
	: x(x), y(y){
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
