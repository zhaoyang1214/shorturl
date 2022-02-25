package bootstrap

import (
	"github.com/zhaoyang1214/ginco/app/console"
	"github.com/zhaoyang1214/ginco/app/providers"
	"github.com/zhaoyang1214/ginco/framework/contract"
	"github.com/zhaoyang1214/ginco/framework/foundation"
	"os"
)

func InitApp() contract.Application {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var app contract.Application = foundation.NewApplication(path)

	registerBaseProviders(app)
	providers.Register(app)
	registerCoreAliases(app)
	registerConfigAliases(app)
	registerBaseCommands(app)
	console.Register(app)
	return app
}
