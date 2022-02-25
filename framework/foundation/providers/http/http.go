package http

import (
	"github.com/gin-gonic/gin"
	"github.com/zhaoyang1214/ginco/framework/contract"
)

type Http struct {
}

var _ contract.Provider = (*Http)(nil)

func (h *Http) Build(container contract.Container, params ...interface{}) (interface{}, error) {
	configServer, err := container.Get("config")
	if err != nil {
		return nil, err
	}
	config := configServer.(contract.Config)
	mode := config.GetString("http.mode")

	gin.SetMode(mode)
	r := gin.New()
	return r, nil
}
