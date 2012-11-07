// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package ui

import (
	"bytes"
	"encoding/csv"
	"strconv"

	"github.com/mccoyst/errorlist"
)

func (ui *Ui) Write(b []byte) (int, error) {
	r := csv.NewReader(bytes.NewReader(b))
	r.Comma = ' '
	r.FieldsPerRecord = -1
	records, err := r.ReadAll()
	if err != nil {
		return 0, err
	}

	n := 0
	var errs []error
	for _, record := range records {
		e := dispatch[record[0]](ui, record[1:])
		if e != nil {
			errs = append(errs, e)
		}
		n += reclen(record)
	}

	return n, errorlist.New(errs...)
}

func reclen(r []string) int {
	n := 0
	if len(r) == 1 {
		n += len(r[0])
	}
	for _, f := range r[1:] {
		n += len(f) + 1
	}
	return n
}

var dispatch = map[string]func(*Ui, []string)error{
	"rectfill": doRectfill,
	"img": doImg,
	"color": doColor,
}

func doRectfill(ui *Ui, args []string) error {
	var x, y, w, h int
	err := parseInts(args, &x, &y, &w, &h)
	if err != nil {
		return err
	}
	fillRect(ui, x, y, w, h)
	return nil
}

func doImg(ui *Ui, args []string) error {
	name := args[0]
	var x, y, subx, suby, w, h int
	err := parseInts(args[1:5], &x, &y, &w, &h)
	if err != nil {
		return err
	}
	shade, err := strconv.ParseFloat(args[5], 32)
	if err != nil {
		return err
	}

	if len(args) >= 7 {
		subx, err = strconv.Atoi(args[6])
		if err != nil {
			return err
		}
	}
	if len(args) == 8 {
		suby, err = strconv.Atoi(args[6])
		if err != nil {
			return err
		}
	}

	drawImg(ui, name, x, y, subx, suby, w, h, float32(shade))
	return nil	
}

func doColor(ui *Ui, args []string) error {
	var r, g, b int
	err := parseInts(args[:3], &r, &g, &b) // meow
	if err != nil {
		return err
	}

	a := 255
	if len(args) == 4 {
		a, err = strconv.Atoi(args[3])
		if err != nil {
			return err
		}
	}

	setColor(ui, r, g, b, a)
	return nil
}

func parseInts(args []string, n ...*int) error {
	var err error
	for i, p := range n {
		*p, err = strconv.Atoi(args[i])
		if err != nil {
			return err
		}
	}
	return nil
}
