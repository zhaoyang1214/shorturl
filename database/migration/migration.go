package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/zhaoyang1214/ginco/framework/contract"
	"github.com/zhaoyang1214/ginco/framework/database"
)

var migrations []*gormigrate.Migration

func Init(a contract.Application) map[string]*gormigrate.Gormigrate {
	m := make(map[string]*gormigrate.Gormigrate)
	// 默认数据库
	db := a.GetI("db").(*database.Database).DB

	m["default"] = gormigrate.New(db, gormigrate.DefaultOptions, migrations)

	//m["other"] = gormigrate.New(otherDb, gormigrate.DefaultOptions, OtherMigrations)

	return m
}
