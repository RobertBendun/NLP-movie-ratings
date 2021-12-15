package main

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type valueKind uint
const (
	nilKind valueKind = iota
	symbolKind
	listKind
	fieldKind
)

type value struct {
	kind valueKind
	sval string
	list []value
	fd field
}

func (v value) String() string {
	switch v.kind {
	case symbolKind:
		return v.sval
	case listKind:
		return fmt.Sprint(v.list)
	case fieldKind:
		return v.fd.String()
	}

	return "nil"
}

const validSymbolChars = "+-*/_-!@#$%^&{}|:;,."

func read(src string) (value, string) {
	src = strings.TrimSpace(src)

	if src == "" {
		return value{kind: nilKind}, ""
	}

	c, cl := utf8.DecodeRune([]byte(src))

	if unicode.IsLetter(c) || unicode.IsNumber(c) || strings.ContainsRune(validSymbolChars, c) {
		var v value
		v.kind = symbolKind
		length := 0
		for _, r := range src[cl:] {
			if unicode.IsLetter(r) || unicode.IsNumber(r) || strings.ContainsRune(validSymbolChars, r) {
				length++
			} else {
				break
			}
		}
		v.sval = src[0:length+1]
		return v, src[length:]
	}

	if c == '(' {
		list := value{}
		list.kind = listKind
		for {
			elem, after := read(src[cl:])
			if elem.kind == nilKind {
				break
			}
			list.list = append(list.list, elem)
			src = after
		}

		return list, src
	}

	if c == ')' {
		return value{kind: nilKind}, src[cl:]
	}

	return value{kind:nilKind}, ""
}

func eval(ctx context, v value) value {
	switch v.kind {
	case symbolKind:
		resolved := value{kind:fieldKind}
		if pos := strings.IndexRune(v.sval, '/'); pos >= 0 {
			resolved.fd = ctx.resolve(v.sval[:pos], v.sval[pos+1:])
		} else {
			resolved.fd = ctx.resolve("", v.sval)
		}
		return resolved

	case listKind:
		if len(v.list) == 0 {
			return value{kind:nilKind}
		}

		f := v.list[0]
		if f.kind != symbolKind {
			return f
		}

		v.list = v.list[1:]
		for i, arg := range v.list {
			v.list[i] = eval(ctx, arg)
		}
		return v
	}

	return v
}
