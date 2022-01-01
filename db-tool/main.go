package main

import "db-tool/subcmd"

func main() {
	subcmd.Run(
		subcmd.New(&tsvImport{}, "import"),
		subcmd.New(&vwExport{}, "vw"))
}
