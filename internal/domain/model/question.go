package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	ID        uuid.UUID `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime `gorm:"index"`
}

type Questions struct {
	BaseModel       `gorm:"embedded"`
	Title           string         `gorm:"type:varchar;not null"`
	Type            string         `gorm:"type:enum('multichoice','descriptive');not null"`
	ShowIf          map[string]int `gorm:"type:json"`
	QuestionnaireID uint           `gorm:"column:questionnaire_id;not null"`
}
