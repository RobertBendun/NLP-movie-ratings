package main

import (
	"log"
	"database/sql"
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

func nonEmpty(val *string, errMessage string) {
	if len(*val) == 0 {
		log.Fatalln(errMessage)
	}
}
