package entity

type ShortUrlCreateRequest struct {
	Url    string `json:"url" binding:"required" example:"http://xx.xx/a/b"`
	Ttl    uint   `json:"ttl"`
	Domain string `json:"domain"  example:"http://xxx.xx"`
}

type ShortUrlCreateResponse struct {
	Url string
}
