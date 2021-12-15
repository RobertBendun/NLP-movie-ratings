package main

import (
	"log"
	"strings"
	"fmt"
)

type context struct {
	fields []field
}

func (ctx context) resolve(path, name string) (resolved field) {
	if len(ctx.fields) == 0 {
		log.Fatalln("Context uninitialized")
	}

	foundOne := false

	for _, fd := range ctx.fields {
		if fd.file == path && fd.name == name {
			return fd
		}

		if (len(path) == 0 || strings.Index(fd.file, path) >= 0) && name == fd.name {
			ensure(!foundOne, fmt.Sprintf("Symbol %s/%s is ambigious", path, name))
			resolved, foundOne = fd, true
			continue
		}
	}

	ensure(foundOne, fmt.Sprintf("Symbol %s/%s is not defined", path, name))
	return
}
