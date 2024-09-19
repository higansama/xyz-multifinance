package request

import (
	"regexp"
)

type HasFullTextSearchDto struct {
	Query string `form:"q" mod:"trim"`
}

type HasSortDto struct {
	Sorts []string `form:"sorts" mod:"dive,trim"`
}

// GetSorts @TODO: useless code
func (h *HasSortDto) GetSorts() []string {
	rs := make([]string, 0)

	for _, v := range h.Sorts {
		field := regexp.MustCompile(`[^a-zA-Z0-9_\-]+`).
			ReplaceAllString(v, "")
		rs = append(rs, field)
	}

	return rs
}

type HasDateRangeDto struct {
	FromDate  string `form:"from_date" mod:"trim"`
	UntilDate string `form:"until_date" mod:"trim"`
}
