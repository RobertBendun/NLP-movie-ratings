package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	var source string
	var database string

	flag.StringVar(&source, "c", "", "Uses query string procided in command line rather then script")
	flag.StringVar(&database, "d", "", "Provides path to the database description file in TSV format")
	// printSchema := flag.Bool("schema", false, "Print loaded fields from database schema and quit")
	// printQuery := flag.Bool("query", false, "Print loaded query from script file or -c argument and quit")

	flag.Parse()

	if len(source) == 0 && len(flag.Arg(0)) == 0 {
		log.Fatalln("Neither script file path nor -c flag was specified")
	}

	if len(database) == 0 {
		log.Fatalln("Database file was not specified")
	}

	schema := loadDatabase(database)
	fmt.Println(schema)

	ctx := context{fields: schema}

	query, _ := read(source)
	fmt.Println(query)

	result := eval(ctx, query)
	fmt.Println(result)

	readers := make(map[string]*csv.Reader)

	ensure(result.kind == listKind, "Only list are supported for TSV output now!")
	for _, entry := range result.list {
		file := entry.fd.file
		if _, ok := readers[file]; !ok {
			f, err := os.Open("../imdb-datasets/" + file + ".tsv")
			ensure(err == nil, fmt.Sprintf("Cannot open `%s`", file))
			readers[file] = csv.NewReader(f)
			readers[file].Comma = '\t'
			readers[file].Read() // skip header
		}
	}

	for {
		cache := make(map[string][]string)
		tsv := strings.Builder{}

		for _, entry := range result.list {
			if tsv.Len() != 0 {
				tsv.WriteByte('\t')
			}

			loaded, ok := cache[entry.fd.file]
			if !ok {
				var err error
				cache[entry.fd.file], err = readers[entry.fd.file].Read()
				if err == io.EOF {
					return
				}
				loaded = cache[entry.fd.file]
			}
			tsv.WriteString(loaded[entry.fd.column])
		}
		fmt.Println(tsv.String())
	}
}
