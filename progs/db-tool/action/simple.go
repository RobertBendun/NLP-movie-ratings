package action

import (
	"strings"
	"unicode"
)

func (IdAction) run(cell string) string {
	return cell
}

func (BagOfWordsAction) run(text string) string {
	bag := strings.Builder{}
	word := strings.Builder{}

	for _, char := range text {
		if char == '\'' {
			continue
		}

		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			word.WriteRune(unicode.ToLower(char))
			continue
		}

		if word.Len() > 0 {
			bag.WriteRune(' ')
			bag.WriteString(word.String())
			word = strings.Builder{}
		}
	}

	if word.Len() > 0 {
		bag.WriteRune(' ')
		bag.WriteString(word.String())
	}

	return bag.String()
}
