package main

type vwExport struct {
	Database string `name:"db" help:"Path to database file"`
	Query string `name:"query" help:"File with SELECT statement"`
	OutputPath string `name:"out" help:"Output file path"`
}

func (params vwExport) Execute() {
	unimplemented()
}
