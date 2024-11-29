package model

// Option represents the options table
type Option struct {
	ID         uint `gorm:"primaryKey"`
	QuestionID uint
	Order      int
	Caption    string `gorm:"size:255"`
	IsCorrect  *bool
	Question   Question `gorm:"foreignKey:QuestionID"`
}
