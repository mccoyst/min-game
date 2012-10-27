// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import "fmt"

const (
	FixedPoint = 4
	FixedMask = 1<<FixedPoint - 1
)

type Fxpt struct {
	int32
}

// F returns a whole fixed point value.
func F(whole int32) Fxpt {
	return Fxpt{whole << FixedPoint}
}

// Fp returns a fixed point value where whole
// represents the whole portion and frac
// represents the fractional portion in units
// of 1/FixedPoint.
func Fp(whole, frac int32) Fxpt {
	return Fxpt{whole<<FixedPoint | frac&FixedMask}
}

// Whole returns the whole component of the fixed point.
func (n Fxpt) Whole() int {
	return int(n.int32 >> FixedPoint)
}

// Strings a human-readable string representing a
// fixed point value.
func (n Fxpt) String() string {
	return fmt.Sprintf("%d.%02x", n.int32>>FixedPoint, n.int32&FixedMask)
}

// Add returns the sum of two fixed points.
func (n Fxpt) Add(m Fxpt) Fxpt {
	return Fxpt{n.int32 + m.int32}
}

// Sub returns the difference when the second
// fixed point is subtracted from the first.
func (n Fxpt) Sub(m Fxpt) Fxpt {
	return Fxpt{n.int32 - m.int32}
}

// Mul returns the product of two fixed points.
func (n Fxpt) Mul(m Fxpt) Fxpt {
	b := int64(n.int32) * int64(m.int32)
	return Fxpt{int32(b >> FixedPoint)}
}

// Div returns the quotient of two fixed points.
func (n Fxpt) Div(m Fxpt) Fxpt {
	b := int64(n.int32) << FixedPoint
	return Fxpt{int32(b / int64(m.int32))}
}

// Rem returns the remainder when dividing two fixed points.
func (n Fxpt) Rem(m Fxpt) Fxpt {
	return Fxpt{n.int32 % m.int32}
}

// Eq returns true if two fixed points are equal.
func (n Fxpt) Eq(m Fxpt) bool {
	return n.int32 == m.int32
}

// Lt returns true if the first value is less than the second.
func (n Fxpt) Lt(m Fxpt) bool {
	return n.int32 < m.int32
}

// Gt returns true if the first value is greater than the second.
func (n Fxpt) Gt(m Fxpt) bool {
	return n.int32 > m.int32
}
