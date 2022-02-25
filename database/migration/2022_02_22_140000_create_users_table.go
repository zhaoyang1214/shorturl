package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/zhaoyang1214/ginco/app/model"
	"gorm.io/gorm"
)

func init() {
	migrations = append(migrations, &gormigrate.Migration{
		ID: "20220222140000",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.User{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("users")
		},
	})
}
