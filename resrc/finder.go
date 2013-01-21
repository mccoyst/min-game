// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package resrc

import (
	"go/build"
	"path/filepath"
)

type PkgFinder struct {
	path string
}

func NewPkgFinder() PkgFinder {
	pkg, err := build.Import("code.google.com/p/min-game/resrc", "", build.FindOnly)
	if err != nil {
		panic(err)
	}

	return PkgFinder{
		path: pkg.Dir,
	}
}

func (p PkgFinder) Find(s string) string {
	return filepath.Join(p.path, s)
}
