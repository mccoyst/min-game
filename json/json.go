package json

import (
	"encoding/json"
	"io"
)

// Encode writes the JSON encoding of name, followed
// by the JSON encoding of v to a writer.  The encoding
// uses a pretty, human-readable indented format.
func Encode(out io.Writer, name string, v interface{}) error {
	b, err := json.Marshal(name)
	if err != nil {
		return err
	}
	b = append(b, '\n')
	if _, err = out.Write(b); err != nil {
		return err
	}

	b, err = json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	_, err = out.Write(b)
	return err
}
