package subcmd

import (
	"flag"
	"fmt"
	"os"
	"path"
	"reflect"
	"strings"
)

type SubCommand interface {
	Execute()
}

type subCmd struct {
	sc          SubCommand
	fields      []reflect.StructField
	names, help []string
	flagSet     *flag.FlagSet
	subcmd      string
}

func New(sc SubCommand, name string) (info subCmd) {
	info.sc = sc

	t := reflect.TypeOf(sc).Elem()
	fieldNames := []string{"name", "help"}

fieldsLoop:
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		var fields [2]string
		for i, fieldName := range fieldNames {
			f, ok := f.Tag.Lookup(fieldName)
			if !ok {
				continue fieldsLoop
			}
			fields[i] = f
		}
		info.fields = append(info.fields, f)
		info.help = append(info.help, fields[1])
		info.names = append(info.names, fields[0])
	}

	info.flagSet = createFlags(info, name)
	return
}

func Run(subcmds ...subCmd) {
	options := listOptions(subcmds)

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "[ERROR] Expected <%s> subcommand\n", options)
		usage(options, subcmds...)
	}

	for _, subcmd := range subcmds {
		if os.Args[1] == subcmd.flagSet.Name() {
			subcmd.flagSet.Parse(os.Args[2:])
			v := reflect.ValueOf(subcmd.sc).Elem()

			for i, field := range subcmd.fields {
				if len(v.FieldByName(field.Name).String()) == 0 {
					fmt.Fprintf(os.Stderr, "[ERROR] Missing argument --%s <value>\n", subcmd.names[i])
					os.Exit(1)
				}
			}

			subcmd.sc.Execute()
			return
		}
	}

	usage(options, subcmds...)
}

func Usage(subcmds ...subCmd) {
	usage(listOptions(subcmds), subcmds...)
}

func createFlags(info subCmd, name string) (flags *flag.FlagSet) {
	flags = flag.NewFlagSet(name, flag.ExitOnError)

	v := reflect.ValueOf(info.sc).Elem()

	for i, field := range info.fields {
		flags.StringVar(
			v.FieldByName(field.Name).Addr().Interface().(*string),
			info.names[i],
			"",
			info.help[i])
	}
	return
}

func listOptions(infos []subCmd) string {
	options := make([]string, 0, len(infos))
	for _, subcmd := range infos {
		options = append(options, subcmd.flagSet.Name())
	}
	return strings.Join(options, "|")
}

func usage(options string, subcmds ...subCmd) {
	fmt.Fprintf(os.Stderr, "usage: %s <%s> [OPTIONS]\n", path.Base(os.Args[0]), options)

	for _, subcmd := range subcmds {
		fmt.Fprintln(os.Stderr)
		subcmd.flagSet.Usage()
	}

	os.Exit(1)
}
