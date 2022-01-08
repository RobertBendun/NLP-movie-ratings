package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type vwExport struct {
	Database   string `name:"db" help:"Path to database file"`
	Query      string `name:"query" help:"File with SELECT statement"`
	OutputPath string `name:"out" help:"Output file path"`
}

func (params vwExport) Execute() {
	db, err := sql.Open("sqlite3", params.Database)
	ensure(err, "Opening SQLite database")

	queryBytes, err := ioutil.ReadFile(params.Query)
	ensure(err, "Opening Query file")

	rows, err := db.Query(string(queryBytes))
	ensure(err, "Executing provided query")
	defer rows.Close()

	types, err := rows.ColumnTypes()
	ensure(err, "Retriving column types")

	data := make([]string, len(types))
	interfaces := make([]interface{}, len(types))
	for i := range data {
		interfaces[i] = &data[i]
	}

	var out *os.File
	if params.OutputPath == "-" {
		out = os.Stdout
	} else {
		out, err = os.Create(params.OutputPath)
		ensure(err, "Creation of train file")
		defer out.Close()
	}

	rowsLoop: for rows.Next() {
		ensure(rows.Scan(interfaces...), "Scanning next row")

		line := strings.Builder{}
		for i, typ := range types {
			switch typ.Name() {
			case "str":
				words := strings.TrimSpace(bagOfWords(data[i]))
				if len(words) == 0 {
					continue rowsLoop
				}
				line.WriteString(words)
			case "nat":
			case "float":
				line.WriteString(data[i])
			default:
				fmt.Fprintf(os.Stderr, "Unrecognized column's type %s", typ.Name())
				os.Exit(1)
			}

			if i == 0 {
				line.WriteString(" | ")
			} else {
				line.WriteString(" ")
			}
		}

		fmt.Fprintln(out, line.String())
	}
}
