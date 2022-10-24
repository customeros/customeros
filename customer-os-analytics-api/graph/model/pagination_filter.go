package model

func (p *PaginationFilter) GetPage() int {
	if p == nil {
		return 1
	}
	return p.Page
}

func (p *PaginationFilter) GetLimit() int {
	if p == nil {
		return 10
	}
	return p.Limit
}
