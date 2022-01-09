package action

type Program []Action

func (prog Program) Group() (int, Action) {
	for i, action := range prog {
		if _, ok := action.(GroupAction); ok {
			return i, action
		}
	}
	return -1, nil
}
