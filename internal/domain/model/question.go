package model

import "gorm.io/gorm"

type Questions struct {
	gorm.Model
	Title           string         `gorm:"type:varchar;not null"`
	Type            string         `gorm:"type:enum('multichoice','descriptive');not null"`
	ShowIf          map[string]int `gorm:"type:json"`
	QuestionnaireID uint           `gorm:"column:questionnaire_id;not null"`
}
