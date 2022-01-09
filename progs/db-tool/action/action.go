package action

import (
	"fmt"
	"strings"
	"unicode"
)

type (
	Action interface {
		fmt.Stringer
		run(cell string) string
	}

	// Pass value without any transformation
	IdAction struct {}

	// Group
	GroupAction struct{}

	// Transforms string using bag of words method
	BagOfWordsAction struct {}

	// When agregating yields first given value as result
	HeadAction struct {
		Action Action
	}

	// When agregating yields concatenation of given values with Delim delimiter
	JoinAction struct {
		Action Action
		Delim string
	}
)

func Run(action Action, cell string) string {
	return strings.TrimSpace(action.run(cell))
}

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

func (HeadAction) run(string) string {
	panic("head.run unimplemented")
}

func (GroupAction) run(string) string {
	panic("group.run unimplemented")
}

func (JoinAction) run(string) string {
	panic("join.run unimplemented")
}
