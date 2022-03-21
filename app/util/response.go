package util

import (
	"github.com/gin-gonic/gin"
	"github.com/zhaoyang1214/ginco/app/entity"
	"net/http"
)

func JsonError(c *gin.Context, msg string, code ...int) {
	reCode := 500
	if len(code) > 0 {
		reCode = code[0]
	}
	c.JSON(http.StatusOK, entity.JSONResult{
		Code:    reCode,
		Message: msg,
	})
}

func JsonResponse(c *gin.Context, data interface{}, code int, msg ...string) {
	reMsg := "OK"
	if len(msg) > 0 {
		reMsg = msg[0]
	}
	c.JSON(http.StatusOK, entity.JSONResult{
		Code:    code,
		Message: reMsg,
		Data:    data,
	})
}
