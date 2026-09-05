package models

// DocumentSeries tracks running numbers for document numbering series
// (prefix + current number). Migrated from the SQLAlchemy model of the
// same name (table: document_series).
type DocumentSeries struct {
	ID            int64   `gorm:"column:id;primaryKey;autoIncrement"`
	Prefix        string  `gorm:"column:prefix;unique;notNull"`
	CurrentNumber int64   `gorm:"column:current_number;notNull"`
	Description   *string `gorm:"column:description;default:-"`
	KataNo        *int    `gorm:"column:kata_no"`
}

func (DocumentSeries) TableName() string {
	return "document_series"
}
