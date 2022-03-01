package model

type User struct {
	ID    uint
	Name  string `gorm:"size:255;not null"`
	Email string `gorm:"uniqueIndex:uk_email;not null"`
}
