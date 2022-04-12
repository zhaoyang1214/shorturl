package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/zhaoyang1214/ginco/framework/contract"
	"net/http"
)

type Index struct {
	*Controller
}

func NewIndex(app contract.Application) *Index {
	return &Index{
		&Controller{
			app: app,
		},
	}
}

func (i Index) Index(c *gin.Context) {
	c.String(http.StatusOK, "Hello Ginco v"+i.app.Version()+"\n")
}
