package contract

type Application interface {
	Container
	Version() string
	BasePath(path string) string
	RuntimePath() string
	GetI(name string) interface{}
}
