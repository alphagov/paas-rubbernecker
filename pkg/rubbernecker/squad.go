package rubbernecker

import (
	"strings"
)

type Squad struct {
	QueryText   string
	DisplayText string
}

func (s *Squad) IsApplied(queries []string) bool {
	return isApplied(queries, strings.Replace(s.QueryText, "filter=", "", -1))
}
