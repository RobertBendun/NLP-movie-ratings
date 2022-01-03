package main

import (
	"strings"
	"unicode"
)

func bagOfWords(text string) string {
	bag := strings.Builder{}

	word := strings.Builder{}

	for _, char := range text {
		if char == '\'' {
			continue
		}

		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			_, err := word.WriteRune(char)
			ensure(err, "word.WriteRune failed")
			continue
		}

		if word.Len() > 0 {
			if bag.Len() > 0 {
				_, err := bag.WriteRune(' ')
				ensure(err, "bag.WriteRune(' ') failed")
			}
			bag.WriteString(word.String())
			word = strings.Builder{}
		}
	}

	if word.Len() > 0 {
		bag.WriteString(word.String())
	}

	return bag.String()
}
