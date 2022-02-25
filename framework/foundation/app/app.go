package app

import (
	"github.com/zhaoyang1214/ginco/framework/contract"
)

var app contract.Application

func Set(a contract.Application) {
	app = a
}

func Get() contract.Application {
	return app
}
