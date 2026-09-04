package models

// Agent represents a sales agent or representative.
type Agent struct {
	ID        int64    `gorm:"column:id;primaryKey;autoIncrement"`
	Agent     *string  `gorm:"column:agent"`
	BonusRate *float64 `gorm:"column:bonus_rate;default:0"`
	Status    *string  `gorm:"column:status"`
	CreatedAt *string  `gorm:"column:created_at"`
	UpdatedAt *string  `gorm:"column:updated_at"`
	IsActive  *bool    `gorm:"column:is_active;default:true"`
	Remarks   *string  `gorm:"column:remarks;default:-"`
}

func (Agent) TableName() string {
	return "agents"
}
