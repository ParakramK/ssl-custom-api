package models

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username string    `gorm:"type:varchar(100);unique;not null"`
	Email    string    `gorm:"type:varchar(100);unique;not null"`
	Password string    `gorm:"type:varchar(255);not null"`
}
