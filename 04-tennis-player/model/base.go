package model

// Page represents a Page of entities
type Page struct {
	Items      interface{} `json:"items"`
	Page       int         `json:"page"`
	PageSize   int         `json:"pageSize"`
	TotalPages int         `json:"totalPages"`
	TotalCount int         `json:"totalCount"`
}

// CalculateTotalPages calculates the total number of pages based on item count and page size
func (p *Page) CalculateTotalPages() {
	pages := p.TotalCount / p.PageSize
	remainder := p.TotalCount - (p.PageSize * pages)
	if remainder > 0 {
		pages++
	}
	p.TotalPages = pages
}
