package action

import "strings"

func (h *HeadAction) keyRun(key, value string) {
	if _, ok := h.called[key]; !ok {
		h.called[key] = Run(h.Action, value)
	}
}

func (h HeadAction) yield(key string) string {
	if v, ok := h.called[key]; ok {
		return v
	}
	return ""
}

func (j *JoinAction) keyRun(key, value string) {
	v, ok := j.joined[key]
	if !ok {
		v = &strings.Builder{}
		j.joined[key] = v
	}
	if v.Len() != 0 {
		v.WriteString(j.Delim)
	}
	v.WriteString(Run(j.Action, value))
}

func (j JoinAction) yield(key string) string {
	if v, ok := j.joined[key]; ok {
		return v.String()
	}
	return ""
}
