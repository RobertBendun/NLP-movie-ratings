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

func (cj LimitedJoinAction) String() string {
	return fmt.Sprintf("limited join\"%s\": %s", cj.Delim, cj.Action.String())
}

func (cj LimitedCountedJoinAction) String() string {
	return fmt.Sprintf("limited counted join\"%s\": %s", cj.Delim, cj.Action.String())
}
