package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	imSet := flag.NewFlagSet("import", flag.ExitOnError)
	database := imSet.String("db", "", "Path to database file")
	table := imSet.String("table", "", "Table where data will be inserted")
	tsv := imSet.String("tsv", "", "Path to TSV file that gonna be imported")

	_ = flag.NewFlagSet("export", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Expected import|export subcommand")
		return
	}

	switch os.Args[1] {
	case "import":
		imSet.Parse(os.Args[2:])
		nonEmpty(database, "Missing --db database option")
		nonEmpty(table, "Missing --table table option")
		nonEmpty(tsv, "Missing --tsv TSV file option")
		importDatabase(*database, *table, *tsv)
	default:
		fmt.Fprintln(os.Stderr, "Expected import|export subcommand")
		return
	}
}
