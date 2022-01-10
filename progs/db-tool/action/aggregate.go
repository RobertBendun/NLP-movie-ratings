package action

import (
	"fmt"
	"strings"
)

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
	if result := Run(j.Action, value); len(result) > 0 {
		if v.Len() != 0 {
			v.WriteString(j.Delim)
		}
		v.WriteString(result)
	}
}

func (j *LimitedJoinAction) keyRun(key, value string) {
	if j.counts[key] > j.limit {
		return
	}
	v, ok := j.joined[key]
	if !ok {
		v = &strings.Builder{}
		j.joined[key] = v
	}

	if result := Run(j.Action, value); len(result) > 0 {
		if v.Len() != 0 {
			v.WriteString(j.Delim)
		}
		v.WriteString(result)
		j.counts[key]++
	}
}

func (j *LimitedCountedJoinAction) keyRun(key, value string) {
	if j.counts[key] > j.limit {
		return
	}

	v, ok := j.joined[key]
	if !ok {
		v = &strings.Builder{}
		j.joined[key] = v
	}
	if result := Run(j.Action, value); len(result) > 0 {
		if v.Len() != 0 {
			v.WriteString(j.Delim)
		}
		v.WriteString(fmt.Sprintf("%s-%d", result, j.counts[key]))
		j.counts[key]++
	}
}

func (j JoinAction) yield(key string) string {
	if v, ok := j.joined[key]; ok {
		return v.String()
	}
	return ""
}
