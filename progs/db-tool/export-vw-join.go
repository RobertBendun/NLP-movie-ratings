package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type vwJoinExport struct {
	Database   string `name:"db" help:"Path to database file"`
	MainQuery  string `name:"main" help:"File with SELECT statement"`
	EachQuery  string `name:"each" help:"File with SELECT statement"`
	OutputPath string `name:"out" help:"Output file path"`
}

const connectionsLimit = 10_000

func (params vwJoinExport) Execute() {
	db, err := sql.Open("sqlite3", params.Database)
	ensure(err, "Opening SQLite database")
	defer db.Close()

	main := readFile(params.MainQuery, "Opening main query file")
	each := readFile(params.EachQuery, "Opening each query file")

	rows, err := db.Query(main)
	ensure(err, "Executing main query")
	defer rows.Close()

	eachStmt, err := db.Prepare(each)
	ensure(err, "Preparing each query")

	types, err := rows.ColumnTypes()
	ensure(err, "Retriving column types")

	data, interfaces := newStringInterfaceArray(len(types))

	for rows.Next() {
		ensure(rows.Scan(interfaces...), "main: Scanning next row")

		subQueryIndex := data[0]
		// line := queryResultToVWString(types[1:], data[1:])

		// TODO determine size of new array

		_, subIfs := newStringInterfaceArray(1)
		subRows, err := eachStmt.Query(sql.Named("ID", subQueryIndex))
		ensure(err, "Executing each query")

		count := 0
		for subRows.Next() {
			ensure(subRows.Scan(subIfs...), "each: Scanning next row")
			count += 1
		}
		fmt.Println(count)
	}

}
