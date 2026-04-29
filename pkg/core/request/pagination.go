package request

type Pagination struct {
	Page     int `json:"page" form:"page" binding:"omitempty,min=1"`
	PageSize int `json:"page_size" form:"page_size" binding:"omitempty,min=1"`
}

func (p *Pagination) GetPage() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) GetPageSize() int {
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	return p.PageSize
}
