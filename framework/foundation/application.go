package foundation

import (
	"github.com/zhaoyang1214/ginco/framework/container"
	"github.com/zhaoyang1214/ginco/framework/contract"
	"os"
	"path"
)

type Application struct {
	contract.Container
	version  string
	basePath string
}

var _ contract.Application = (*Application)(nil)

func NewApplication(basePath string) *Application {
	var c = container.NewContainer()
	a := &Application{
		c,
		"1.0.6",
		"",
	}

	if basePath == "" {
		a.basePath, _ = os.Getwd()
	} else {
		a.basePath = basePath
	}
	a.Set("container", a)
	a.Alias("container", "app")
	return a
}

func (a *Application) Version() string {
	return a.version
}

func (a *Application) BasePath(joinPath string) string {
	return path.Join(a.basePath, joinPath)
}

func (a *Application) RuntimePath() string {
	return a.BasePath("runtime")
}

func (a *Application) GetI(name string) interface{} {
	s, _ := a.Get(name)
	return s
}
