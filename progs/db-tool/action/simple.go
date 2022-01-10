package action

import (
	"strings"
	"unicode"
	"fmt"
)

func (IdAction) run(cell string) string {
	return cell
}

func (p PrefixAction) run(cell string) string {
	return fmt.Sprintf("%s%s", p.prefix, cell)
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
