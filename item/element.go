// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package item

type Element struct {
	string
}

func (e *Element) Name() string {
	return e.string
}

func (e *Element) Desc() string {
	return elemDesc[e.string]
}

var elemDesc = map[string]string{
	"Uranium": `Uranium is of great interest because of its application to nuclear power
		and nuclear weapons. Uranium contamination is an emotive environmental
		problem. It is not particularly rare and is more common than beryllium
		or tungsten for instance.`, // http://www.webelements.com/uranium/
}
