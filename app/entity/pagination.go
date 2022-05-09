package entity

type Pagination struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
}

func (p *Pagination) GetPage() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) GetSize() int {
	if p.Size <= 0 {
		p.Size = 10
	}
	return p.Size
}
