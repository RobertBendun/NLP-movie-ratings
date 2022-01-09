package action

import (
	"fmt"
	"strings"
)

type (
	Action interface {
		fmt.Stringer
	}

	SimpleAction interface {
		Action
		run(cell string) string
	}

	AggregateAction interface {
		Action
		keyRun(key, value string)
		yield(key string) string
	}

	// Pass value without any transformation
	IdAction struct{}

	// Transforms string using bag of words method
	BagOfWordsAction struct{}

	// When agregating yields first given value as result
	HeadAction struct {
		Action Action
		called map[string]string
	}

	// Group
	GroupAction struct {
		Keys map[string]struct{}
	}

	// When agregating yields concatenation of given values with Delim delimiter
	JoinAction struct {
		Action Action
		Delim  string
		joined map[string]*strings.Builder
	}
)

func Run(action Action, cell string) string {
	if v, ok := action.(SimpleAction); ok {
		return strings.TrimSpace(v.run(cell))
	}
	panic(fmt.Sprint("Group action called in non-group context: ", action))
}

func KeyRun(action Action, key, value string) {
	if v, ok := action.(AggregateAction); ok {
		v.keyRun(key, value)
		return
	}
	panic(fmt.Sprint("Non-group action called in group context: ", action))
}

func Yield(action Action, key string) string {
	if v, ok := action.(AggregateAction); ok {
		return v.yield(key)
	}
	panic("Non-group action called in group context")
}
