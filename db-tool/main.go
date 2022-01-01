package main

import (
	"flag"
	"fmt"
	"os"
)


func main() {
	var imData tsvImport
	im := flag.NewFlagSet("import", flag.ExitOnError)
	im.StringVar(&imData.database, "db", "", "Path to database file")
	im.StringVar(&imData.table, "table", "", "Table where data will be inserted")
	im.StringVar(&imData.tsv, "tsv", "", "Path to TSV file that gonna be imported")

	var vwData vwExport
	vw := flag.NewFlagSet("vw", flag.ExitOnError)
	vw.StringVar(&vwData.database, "db", "", "Path to database file")
	vw.StringVar(&vwData.outputPath, "out", "", "Output file path")
	vw.StringVar(&vwData.query, "query", "", "File with SELECT statement")

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Expected import|export subcommand")
		return
	}

	switch os.Args[1] {
	case "import":
		im.Parse(os.Args[2:])
		nonEmpty(imData.database, "Missing --db database option")
		nonEmpty(imData.table, "Missing --table table option")
		nonEmpty(imData.tsv, "Missing --tsv TSV file option")
		imData.execute()
	case "vw":
		vw.Parse(os.Args[2:])
		nonEmpty(vwData.database, "Missing --db database option")
		nonEmpty(vwData.query, "Missing --query query_file option")
		nonEmpty(vwData.outputPath, "Missing --out output_file option")
		vwData.execute()
	default:
		fmt.Fprintln(os.Stderr, "Expected import|export subcommand")
		return
	}
}
