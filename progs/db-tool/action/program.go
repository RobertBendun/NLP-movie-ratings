package action

type Program []Action

func (prog Program) Group() (int, *GroupAction) {
	for i, action := range prog {
		if v, ok := action.(*GroupAction); ok {
			return i, v
		}
	}
	return -1, nil
}
