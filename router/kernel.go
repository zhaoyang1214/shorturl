package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zhaoyang1214/ginco/app/http/controller"
	"github.com/zhaoyang1214/ginco/framework/contract"
)

func Register(app contract.Application) {
	routerServer, err := app.Get("router")
	if err != nil {
		panic(err)
	}
	r := routerServer.(*gin.Engine)

	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/", controller.Index)
	r.GET("/:hash", controller.UrlRedirect(app))

	registerApi(app)
}
