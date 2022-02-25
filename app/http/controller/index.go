package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/zhaoyang1214/ginco/app/model"
	"github.com/zhaoyang1214/ginco/framework/contract"
	"github.com/zhaoyang1214/ginco/framework/foundation/app"
	"net/http"
)

func Index(c *gin.Context) {
	c.String(http.StatusOK, "Hello Ginco v"+app.Get().Version()+"\n")
}

func Name(app contract.Application) func(*gin.Context) {
	return func(c *gin.Context) {
		if userVal, ok := c.Get("user"); ok {
			user := userVal.(*model.User)
			c.String(http.StatusOK, "Hello userid %d username %s ", user.ID, user.Name)
			return
		}
		c.String(http.StatusOK, "My name is "+app.GetI("config").(contract.Config).GetString("app.name")+"\n")
	}
}

// Hello World
// @Summary test
// @Schemes
// @Description test
// @Tags
// @Accept json
// @Produce json
// @Success 200 {object} entity.JSONResult{code=int,message=string,data=string} "helloworld"
// @Router /helloworld [get]
func Helloworld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "OK",
		"data":    "Hello World",
	})
}
