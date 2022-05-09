package entity

type ShortUrlCreateRequest struct {
	Url    string `json:"url" binding:"required" example:"http://xx.xx/a/b"`
	Ttl    uint   `json:"ttl"`
	Domain string `json:"domain"  example:"http://xxx.xx"`
}

type ShortUrlCreateResponse struct {
	Url string
}

type ShortUrlListRequest struct {
	Pagination
}

type ShortUrlListResponseWithList struct {
	Hash      string `json:"hash"`
	Url       string `json:"url"`
	Ttl       uint   `json:"ttl"`
	Domain    string `json:"domain"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ShortUrlListResponse struct {
	Total int64                          `json:"total"`
	List  []ShortUrlListResponseWithList `json:"list"`
}
