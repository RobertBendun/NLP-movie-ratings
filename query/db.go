package main

import (
	"fmt"
	"os"
	"encoding/csv"
)

type kind uint

const (
	kindString kind = iota
	kindFloat
	kindNatural
)

type field struct {
	file   string
	name   string
	kind   kind
	column uint
}

func (f field) String() string {
	return fmt.Sprintf("<field %s/%s : %d>", f.file, f.name, f.kind)
}

func loadDatabase(file string) (fields []field) {
	src, err := os.Open(file)
	ensure(err == nil, fmt.Sprintf("Cannot read database file %s", file))
	defer src.Close()

	reader := csv.NewReader(src)
	reader.Comma = '\t'
	entries, err := reader.ReadAll()
	ensure(err == nil, fmt.Sprintf("Error while reading TSV database schema %s: %s", file, err))

	columns := make(map[string]uint)
	for _, entry := range entries[1:] {
		ensure(len(entry) == 3, "Invalid database format")

		fd := field{}
		fd.file = entry[0]
		fd.name = entry[1]
		fd.kind = 0 // TODO add kind checking

		if v, ok := columns[fd.file]; ok {
			columns[fd.file], fd.column = v+1, v+1
		} else {
			fd.column = 0
			columns[fd.file] = 0
		}

		fields = append(fields, fd)
	}
	return
}
