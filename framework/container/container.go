package container

import (
	"errors"
	"github.com/zhaoyang1214/ginco/framework/contract"
)

type Container struct {
	instances map[string]interface{}
	providers map[string]contract.Provider
	aliases   map[string]string
}

var _ contract.Container = (*Container)(nil)

func NewContainer() *Container {
	return &Container{
		instances: make(map[string]interface{}),
		providers: make(map[string]contract.Provider),
		aliases:   make(map[string]string),
	}
}

func (c *Container) Set(name string, entry interface{}) {
	if _, ok := c.aliases[name]; ok {
		delete(c.aliases, name)
	}
	c.instances[name] = entry
}

func (c *Container) Get(name string) (interface{}, error) {
	name = c.GetNameByAlias(name)
	if entry, ok := c.instances[name]; ok {
		return entry, nil
	}
	entry, err := c.Make(name)
	if err != nil {
		return nil, err
	}
	c.instances[name] = entry
	return entry, nil
}

func (c *Container) Make(name string, params ...interface{}) (interface{}, error) {
	provider, ok := c.providers[name]
	if !ok {
		return nil, errors.New("No provider for " + name)
	}
	return provider.Build(c, params...)
}

func (c *Container) Has(name string) bool {
	if c.HasInstance(name) {
		return true
	}
	if c.HasProvider(name) {
		return true
	}
	if c.HasAlias(name) {
		return true
	}
	return false
}

func (c *Container) Bind(name string, provider contract.Provider) {
	if _, ok := c.instances[name]; ok {
		delete(c.instances, name)
	}
	if _, ok := c.aliases[name]; ok {
		delete(c.aliases, name)
	}

	c.providers[name] = provider
}

func (c *Container) GetNameByAlias(alias string) string {
	name, ok := c.aliases[alias]
	if !ok {
		return alias
	}
	return c.GetNameByAlias(name)
}

func (c *Container) Alias(name, alias string) {
	if c.GetNameByAlias(name) == c.GetNameByAlias(alias) {
		panic(errors.New(name + " is aliased to itself."))
	}
	c.aliases[alias] = name
}

func (c *Container) HasInstance(name string) bool {
	if _, ok := c.instances[name]; ok {
		return true
	}
	return false
}

func (c *Container) HasAlias(alias string) bool {
	if _, ok := c.aliases[alias]; ok {
		return true
	}
	return false
}

func (c *Container) HasProvider(name string) bool {
	if _, ok := c.providers[name]; ok {
		return true
	}
	return false
}

func (c *Container) ForgetInstances() {
	c.instances = make(map[string]interface{})
}

func (c *Container) ForgetInstance(name string) {
	delete(c.instances, name)
}
