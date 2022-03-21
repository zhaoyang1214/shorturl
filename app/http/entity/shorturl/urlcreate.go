package shorturl

type UrlCreateRequest struct {
	Url       string    `json:"url" binding:"required" example:"http://xx.xx/a/b"`
	Ttl       uint      `json:"ttl"`
	Domain    string    `json:"domain"  example:"http://xxx.xx"`
}

type UrlCreateResponse struct {
	Url string
}