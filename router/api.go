package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/zhaoyang1214/ginco/app/http/controller"
	"github.com/zhaoyang1214/ginco/docs"
	"github.com/zhaoyang1214/ginco/framework/contract"
)

func registerApi(app contract.Application, s *controller.ShortUrl) {
	r := app.GetI("router").(*gin.Engine)
	docs.SwaggerInfo_swagger.BasePath = "/api"
	api := r.Group("/api")
	//authMiddleware := middleware.JWT(app)
	//api.Use(authMiddleware.MiddlewareFunc())
	{
		api.POST("/url", s.Create)
		api.GET("/url", s.List)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}
