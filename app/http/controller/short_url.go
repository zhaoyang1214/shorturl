package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/zhaoyang1214/ginco/app/entity"
	"github.com/zhaoyang1214/ginco/app/service"
	"github.com/zhaoyang1214/ginco/framework/contract"
	"net/http"
)

type ShortUrl struct {
	*Controller
	service *service.ShortUrl
}

func NewShortUrl(app contract.Application) *ShortUrl {
	return &ShortUrl{
		Controller: &Controller{
			app: app,
		},
		service: service.NewShortUrl(app),
	}
}

// Url Create
// @Summary Url Create
// @Accept json
// @Produce json
// @Param root body entity.ShortUrlCreateRequest true "UrlCreate info"
// @Success 200 {object} entity.ResultJSON{code=int,message=string,data=entity.ShortUrlCreateResponse}
// @BasePath /api
// @Router /url [post]
func (s ShortUrl) Create(c *gin.Context) {
	var ucr entity.ShortUrlCreateRequest
	if err := c.ShouldBindJSON(&ucr); err != nil {
		s.responseJsonError(c, "Bind and validate params error, "+err.Error(), http.StatusBadRequest)
		return
	}

	code, data, err := s.service.Create(ucr)
	s.responseJson(c, code, err, data)
}

func (s ShortUrl) Redirect(c *gin.Context) {
	hash := c.Param("hash")
	code, url, err := s.service.GetUrlByHash(hash)
	if err != nil {
		c.Status(code)
		return
	}
	c.Redirect(code, url)
	return
}

func (s ShortUrl) List(c *gin.Context) {

}
