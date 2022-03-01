package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/zhaoyang1214/ginco/app/model"
	"gorm.io/gorm"
)

func init() {
	migrations = append(migrations, &gormigrate.Migration{
		ID: "20220301160000",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.Domain{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(model.Domain{}.TableName())
		},
	})
}
