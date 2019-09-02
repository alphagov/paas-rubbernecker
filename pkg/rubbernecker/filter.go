package rubbernecker

import (
	"strings"
)

type Filter interface {
	QueryText() string
	DisplayText() string
	IsApplied(queries []string) bool
}

func DefaultFilterSet() []Filter {
	return []Filter{
		&BlockedFilter{},
		&ScheduledFilter{},
		&CommentsToResolveFilter{},
		&SmallTaskFilter{},
		&DocumentationFilter{},
		&PairingFilter{},
		&NonTechFilter{},
		&TechFilter{},
	}
}

type BlockedFilter struct{}

func (*BlockedFilter) QueryText() string {
	return "filter=sticker:blocked"
}

func (*BlockedFilter) DisplayText() string {
	return "Blocked"
}

func (*BlockedFilter) IsApplied(queries []string) bool {
	return isApplied(queries, "sticker:blocked")
}

type ScheduledFilter struct {}

func (*ScheduledFilter) QueryText() string {
	return "filter=sticker:scheduled"
}

func (*ScheduledFilter) DisplayText() string {
	return "Scheduled"
}

func (*ScheduledFilter) IsApplied(queries []string) bool {
	return isApplied(queries, "sticker:scheduled")
}

type CommentsToResolveFilter struct {}

func (*CommentsToResolveFilter) QueryText() string {
	return "filter=sticker:comments-to-resolve"
}

func (*CommentsToResolveFilter) DisplayText() string {
	return "Comments to resolve"
}

func (*CommentsToResolveFilter) IsApplied(queries []string) bool {
	return isApplied(queries, "sticker:comments-to-resolve")
}

type SmallTaskFilter struct {}

func (SmallTaskFilter) QueryText() string {
	return "filter=sticker:'small' task"
}

func (SmallTaskFilter) DisplayText() string {
	return "Small tasks"
}

func (SmallTaskFilter) IsApplied(queries []string) bool {
	return isApplied(queries, "sticker:'small' task")
}

type DocumentationFilter struct{}

func (*DocumentationFilter) QueryText() string {
	return "filter=sticker:documentation"
}

func (*DocumentationFilter) DisplayText() string {
	return "Documentation"
}

func (*DocumentationFilter) IsApplied(queries []string) bool {
	return isApplied(queries, "sticker:documentation")
}

type PairingFilter struct{}

func (*PairingFilter) QueryText() string {
	return "filter=sticker:pairing"
}

func (*PairingFilter) DisplayText() string {
	return "Pairing"
}

func (*PairingFilter) IsApplied(queries []string) bool {
	return isApplied(queries, "sticker:pairing")
}

type NonTechFilter struct {}

func (*NonTechFilter) QueryText() string {
	return "filter=sticker:non-tech"
}

func (*NonTechFilter) DisplayText() string {
	return "Non-tech"
}

func (*NonTechFilter) IsApplied(queries []string) bool {
	return isApplied(queries, "sticker:non-tech")
}

type TechFilter struct {}

func (*TechFilter) QueryText() string {
	return "filter=not-sticker:non-tech"
}

func (*TechFilter) DisplayText() string {
	return "Tech"
}

func (*TechFilter) IsApplied(queries []string) bool {
	return isApplied(queries, "not-sticker:non-tech")
}

func isApplied(queries []string, term string) bool {
	term = strings.ToLower(term)
	for _, q := range queries {
		if strings.ToLower(q) == term {
			return true
		}
	}

	return false
}

