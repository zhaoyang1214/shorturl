package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/zhaoyang1214/ginco/framework/foundation/app"
	"net/http"
)

func Index(c *gin.Context) {
	c.String(http.StatusOK, "Hello Ginco v"+app.Get().Version()+"\n")
}
