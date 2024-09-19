package pagination

import (
	"github.com/gin-gonic/gin"
	"math"
	"strconv"
)

type HasPaginationRequest struct {
	PaginationRequestDto RequestDto
}

func (p *HasPaginationRequest) BindPagination(ctx *gin.Context) {
	page := ctx.Query("page")
	if page == "" {
		page = "1"
	}
	vpage, _ := strconv.Atoi(page)
	if vpage < 1 {
		vpage = 1
	}

	perPage := ctx.Query("per_page")
	if perPage == "" {
		perPage = "10"
	}
	vperpage, _ := strconv.Atoi(perPage)
	if vperpage < 1 {
		vperpage = 1
	} else if vperpage > 100 {
		// maximum per page is 100
		vperpage = 100
	}

	p.PaginationRequestDto.Page = int64(vpage)
	p.PaginationRequestDto.PerPage = int64(vperpage)
}

type RequestDto struct {
	Page    int64
	PerPage int64
}

func (p RequestDto) Offset() int64 {
	return (p.Page - 1) * p.PerPage
}

type Meta struct {
	Total       int64 `json:"total"`
	PerPage     int64 `json:"per_page"`
	From        int64 `json:"from"`
	To          int64 `json:"to"`
	CurrentPage int64 `json:"current_page"`
	LastPage    int64 `json:"last_page"`
}

func NewMeta(total int64, currentPage int64, perPage int64) Meta {
	var from int64
	var to int64
	lastPage := int64(math.Max(math.Ceil(float64(total)/float64(perPage)), 1.0))

	if currentPage > 0 && currentPage <= lastPage {
		f := (currentPage-1)*perPage + 1
		from = f

		t := f + perPage - 1
		if currentPage == lastPage {
			t = total
		}
		to = t
	}

	if total == 0 {
		from = 0
		to = 0
	}

	return Meta{
		CurrentPage: currentPage,
		PerPage:     perPage,
		Total:       total,
		From:        from,
		To:          to,
		LastPage:    lastPage,
	}
}
