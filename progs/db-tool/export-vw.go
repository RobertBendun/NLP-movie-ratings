package main

import (
	"database/sql"
	"db-tool/action"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type vwExport struct {
	Database   string `name:"db" help:"Path to database file"`
	Query  string `name:"query" help:"File with SELECT statement"`
	OutputPath string `name:"out" help:"Output file path"`
}

func (params vwExport) Execute() {
	db, err := sql.Open("sqlite3", params.Database)
	ensure(err, "Opening SQLite database")
	defer db.Close()

	main := readFile(params.Query, "Opening main query file")

	var out *os.File
	if params.OutputPath == "-" {
		out = os.Stdout
	} else {
		out, err = os.Create(params.OutputPath)
		ensure(err, "Creation of out file")
		defer out.Close()
	}

	rows, err := db.Query(main)
	ensure(err, "Executing main query")
	defer rows.Close()

	types, err := rows.ColumnTypes()
	ensure(err, "Retriving column types")

	actions := make(action.Program, len(types))

	for i, typ := range types {
		actions[i], err = action.From(typ.Name())
		ensure(err, "Creating action")
	}

	data, itfs := newStringInterfaceArray(len(types))

	start := time.Now()
	count := 0

	groupIndex, group := actions.Group()

	if groupIndex < 0 {
		nextScan:
		for rows.Next() {
			count++
			ensure(rows.Scan(itfs...), "Row scanning")

			line := strings.Builder{}
			for i, act := range actions {
				result := strings.TrimSpace(action.Run(act, data[i]))
				if len(result) == 0 {
					continue nextScan
				}
				line.WriteString(result)
				if i == 0 {
					line.WriteString(" | ")
				} else {
					line.WriteString(" ")
				}
			}
			final := line.String()
			if len(final) > 0 {
				fmt.Fprintln(out, final)
			}
		}
	} else {
		for rows.Next() {
			count++
			ensure(rows.Scan(itfs...), "Group row scanning")
			key := data[groupIndex]
			group.Keys[key] = struct{}{}
			for i, act := range actions {
				if i == groupIndex {
					continue
				}
				action.KeyRun(act, key, data[i])
			}
		}

		for key := range group.Keys {
			line := strings.Builder{}
			for i, act := range actions {
				if i == groupIndex {
					continue
				}
				result := strings.TrimSpace(action.Yield(act, key))
				if len(result) > 0 {
					line.WriteString(result)
					if i == 0 {
						line.WriteString(" | ")
					} else {
						line.WriteString(" ")
					}
				}
			}

			fmt.Fprintln(out, line.String())
		}
	}

	fmt.Printf("Processed %d rows in %v\n", count, time.Since(start))
}
