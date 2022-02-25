package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func registerApi(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/info", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  "ok",
				"data": gin.H{
					"name": "ginco",
					"age":  0,
				},
			})
		})
	}

}
