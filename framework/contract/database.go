package contract

import "gorm.io/gorm"

type Database interface {
	Connection(names ...string) *gorm.DB
	Resolve(name string) *gorm.DB
}
