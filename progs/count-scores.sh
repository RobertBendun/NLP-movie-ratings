#!/usr/bin/sh

[ "$1" ] || {
	echo "usage: $(basename "$0") <filename>"
	exit 1
}

cut -d' ' -f1 "$1" | sort -n | uniq -c
