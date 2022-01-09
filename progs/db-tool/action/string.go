package action

import "fmt"

func (IdAction) String() string {
	return "id"
}

func (GroupAction) String() string {
	return "group"
}

func (BagOfWordsAction) String() string {
	return "bag"
}

func (head HeadAction) String() string {
	return "head: " + head.Action.String()
}

func (join JoinAction) String() string {
	return fmt.Sprintf("join\"%s\": %s", join.Delim, join.Action.String())
}
