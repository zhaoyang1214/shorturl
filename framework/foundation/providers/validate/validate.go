package validate

import (
	"github.com/go-playground/validator/v10"
	"github.com/zhaoyang1214/ginco/framework/contract"
)

type Validate struct {
}

var _ contract.Provider = (*Validate)(nil)

func (v *Validate) Build(container contract.Container, params ...interface{}) (interface{}, error) {
	return validator.New(), nil
}
