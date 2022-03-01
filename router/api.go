package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/zhaoyang1214/ginco/app/http/controller"
	"github.com/zhaoyang1214/ginco/app/http/middleware"
	"github.com/zhaoyang1214/ginco/docs"
	"github.com/zhaoyang1214/ginco/framework/contract"
	"net/http"
)

func registerApi(app contract.Application) {
	r := app.GetI("router").(*gin.Engine)
	authMiddleware := middleware.JWT(app)
	docs.SwaggerInfo_swagger.BasePath = "/api"
	api := r.Group("/api")
	api.Use(authMiddleware.MiddlewareFunc())
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
		api.GET("/helloworld", controller.Helloworld)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}
