package model

import (
	"time"
)

type ShortUrl struct {
	ID        uint      `gorm:"primarykey"`
	Hash      uint32    `gorm:"not null;comment:短链MurmurHash;uniqueIndex:uk_hash"`
	Url       string    `gorm:"size:2000;not null;comment:url"`
	Ttl       uint      `gorm:"not null;comment:有效期,单位s,0不限制"`
	Domain    string    `gorm:"size:50;not null;comment:域名"`
	CreatedAt time.Time `gorm:"type:DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP;index:idx_created_at"`
	UpdatedAt time.Time `gorm:"type:DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

func (ShortUrl) TableName() string {
	return "short_url"
}
