package model

func (p *Pagination) GetPage() int {
	if p == nil {
		return 1
	}
	return p.Page
}

func (p *Pagination) GetLimit() int {
	if p == nil {
		return 10
	}
	return p.Limit
}
