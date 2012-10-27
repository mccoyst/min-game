// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import "fmt"

const FixedPoint = 4

type Fxpt struct {
	int32
}

func F(whole int32) Fxpt {
	return Fxpt{whole << FixedPoint}
}

func Fp(whole, frac int32) Fxpt {
	return Fxpt{whole<<FixedPoint | frac&0x0f}
}

func (n Fxpt) Whole() int {
	return int(n.int32 >> FixedPoint)
}

func (n Fxpt) String() string {
	return fmt.Sprintf("%d.%02x", n.int32>>FixedPoint, n.int32&0x0f)
}

func (n Fxpt) Add(m Fxpt) Fxpt {
	return Fxpt{n.int32 + m.int32}
}

func (n Fxpt) Sub(m Fxpt) Fxpt {
	return Fxpt{n.int32 - m.int32}
}

func (n Fxpt) Mul(m Fxpt) Fxpt {
	b := int64(n.int32) * int64(m.int32)
	return Fxpt{int32(b >> FixedPoint)}
}

func (n Fxpt) Div(m Fxpt) Fxpt {
	b := int64(n.int32) << FixedPoint
	return Fxpt{int32(b / int64(m.int32))}
}

func (n Fxpt) Eq(m Fxpt) bool {
	return n.int32 == m.int32
}

func (n Fxpt) Lt(m Fxpt) bool {
	return n.int32 < m.int32
}

func (n Fxpt) Gt(m Fxpt) bool {
	return n.int32 > m.int32
}
