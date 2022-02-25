package contract

type Container interface {
	Set(name string, entry interface{})
	Get(name string) (interface{}, error)
	Make(name string, params ...interface{}) (interface{}, error)
	Has(name string) bool
	Bind(name string, provider Provider)
	GetNameByAlias(alias string) string
	Alias(name, alias string)
	HasInstance(name string) bool
	HasAlias(name string) bool
	HasProvider(name string) bool
}
