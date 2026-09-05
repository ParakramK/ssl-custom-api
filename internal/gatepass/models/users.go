package models

// User maps the users table.
type User struct {
	Id       int64   `gorm:"column:id;primaryKey"`
	Name     *string `gorm:"column:name"`
	Username *string `gorm:"column:username"`
	Password *string `gorm:"column:password"`
	RoleId   *int    `gorm:"column:role_id"`
	KataNo   *int    `gorm:"column:kata_no"`
	CompId   *int    `gorm:"column:comp_id"`
}

func (User) TableName() string { return "users" }
