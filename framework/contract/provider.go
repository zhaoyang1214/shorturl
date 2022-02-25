package contract

type Provider interface {
	Build(Container, ...interface{}) (interface{}, error)
}
