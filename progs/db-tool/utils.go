package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
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

func nonEmpty(val string, errMessage string) {
	if len(val) == 0 {
		log.Fatalln(errMessage)
	}
}

func unimplemented() {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		fmt.Fprintf(os.Stderr, "%s:%d: error: unimplemented\n", file, line)
	} else {
		fmt.Fprintln(os.Stderr, "unimplemented at unknown location")
	}
	os.Exit(1)
}

func readFile(file string, message string) string {
	bytes, err := ioutil.ReadFile(file)
	ensure(err, message)
	return string(bytes)
}

func newStringInterfaceArray(count int) (data []string, interfaces []interface{}) {
	data = make([]string, count)
	interfaces = make([]interface{}, count)
	for i := range data {
		interfaces[i] = &data[i]
	}
	return
}
