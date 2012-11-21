// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

type DumbFinder struct{
}

func (d DumbFinder) Find(s string) string {
	return "./resrc/" + s
}
