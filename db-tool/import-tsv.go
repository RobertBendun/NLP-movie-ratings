package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type tsvImport struct {
	database, table, tsv string
}

func (params tsvImport) execute() {
	tsvData, err := os.Open(params.tsv)
	ensure(err, "Opening file")

	r := csv.NewReader(tsvData)
	r.Comma = '\t'
	r.LazyQuotes = true
	r.FieldsPerRecord = -1
	header, err := r.Read()
	ensure(err, "Reading header")

	db, err := sql.Open("sqlite3", params.database)
	ensure(err, "Opening SQLite database")
	defer db.Close()

	args := strings.Repeat("?,", len(header))
	insertStatementText := fmt.Sprintf(
		"insert into %s values (%s)",
		params.table, args[:len(args)-1])

	tx, err := db.Begin()
	ensure(err, "Start of transaction")
	stmt, err := tx.Prepare(insertStatementText)
	ensure(err, "Prepering SQL statement")
	defer stmt.Close()

	count := 0
	for {
		row, err := r.Read()
		if err != nil && err == io.EOF {
			break
		}
		ensure(err, "Reading row from TSV")
		if len(row) != len(header) {
			continue
		}
		_, err = stmt.Exec(str2sql(row)...)
		ensure(err, "Executing SQL statement")
		count++
	}

	err = tx.Commit()
	ensure(err, "Commit")
	fmt.Printf("Inserted %d rows\n", count)
}
