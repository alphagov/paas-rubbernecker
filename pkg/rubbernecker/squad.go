package rubbernecker

import (
	"fmt"
	"strings"
)

type Squad struct {
	QueryText   string
	DisplayText string
}

func (s *Squad) IsApplied(queries []string) bool {
	fmt.Println("queries", queries)
	fmt.Println("s.QueryText", strings.Replace(s.QueryText, "filter=", "", -1))
	return isApplied(queries, strings.Replace(s.QueryText, "filter=", "", -1))
}
