package main

import "log"

func ensure(b bool, m string) {
	if !b {
		log.Fatalln(m)
	}
}
