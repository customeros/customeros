package utils

import (
	"math"
)

type PaginationRequestBody struct {
	Limit int
	Page  int
}

type Pagination struct {
	Limit      int
	Page       int
	TotalRows  int64
	TotalPages int
	Rows       interface{}
}

func (p *Pagination) GetSkip() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.Limit <= 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Pagination) GetPage() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) SetTotalRows(totalRows int64) {
	p.TotalRows = totalRows
	p.TotalPages = int(math.Ceil(float64(totalRows) / float64(p.GetLimit())))
}

func (p *Pagination) SetRows(rows interface{}) {
	p.Rows = rows
}
