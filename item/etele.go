// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package item

import (
	"fmt"
)

// Etele is the emergency self-teleporter.
type Etele struct {
	Uses int
}

func (e *Etele) Name() string {
	return "E-Tele"
}

func (e *Etele) Desc() string {
	return fmt.Sprintf("The Emergency Teleporter reacts to critical condition by sending "+
		"you back to home base. This E-Tele currently has %d uses remaining.",
		e.Uses)
}

// Use return true iff the item could be used.
func (e *Etele) Use() bool {
	if e.Uses > 0 {
		e.Uses--
		return true
	}
	return false
}
