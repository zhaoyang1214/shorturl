package model

import (
	"gorm.io/gorm"
	"time"
)

type Domain struct {
	ID        uint           `gorm:"primarykey"`
	Name      string         `gorm:"size:255;not null;comment:名称"`
	Addr      string         `gorm:"size:100;not null;comment:地址;uniqueIndex:uk_addr_deleted_at"`
	Remark    string         `gorm:"size:255;not null;default:'';comment:备注"`
	CreatedAt time.Time      `gorm:"type:DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"type:DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"type:DATETIME NULL;uniqueIndex:uk_addr_deleted_at"`
}

func (Domain) TableName() string {
	return "domain"
}
