package main

import (
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func ensure(err error, why string) {
	if err != nil {
		log.Fatalf("%s: %s\n", why, err.Error())
	}
}

func str2sql(arr []string) (res []interface{}) {
	for _, s := range arr {
		if s == "\\N" {
			res = append(res, sql.NullString{})
		} else {
			res = append(res, s)
		}
	}
	return
}

func main() {
	database := flag.String("db", "", "Path to database file")
	table := flag.String("table", "", "Table where data will be inserted")
	tsv := flag.String("tsv", "", "Path to TSV file that gonna be imported")
	flag.Parse()

	if len(*database) == 0 {
		log.Fatalln("Missing --db database option")
	}

	if len(*tsv) == 0 {
		log.Fatalln("Missing --tsv TSV file option")
	}

	if len(*table) == 0 {
		log.Fatalln("Missing --table table option")
	}

	tsvData, err := os.Open(*tsv)
	ensure(err, "Opening file")

	r := csv.NewReader(tsvData)
	r.Comma = '\t'
	r.LazyQuotes = true
	r.FieldsPerRecord = -1
	header, err := r.Read()
	ensure(err, "Reading header")

	db, err := sql.Open("sqlite3", *database)
	ensure(err, "Opening SQLite database")
	defer db.Close()

	args := strings.Repeat("?,", len(header))
	insertStatementText := fmt.Sprintf("insert into %s values (%s)", *table, args[:len(args)-1])

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
