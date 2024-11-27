package model

type Questions struct {
	BaseModel           `gorm:"embedded"`
	Title               string `gorm:"type:varchar;not null"`
	Type                string `gorm:"type:enum('multichoice','descriptive');not null"`
	FilePath            string `gorm:"type:varchar;column:file_path"`
	Order               int    `gorm:"type:int"`
	QuestionnaireID     uint   `gorm:"column:questionnaire_id;not null"`
	DependsOnQuestionId uint   `gorm:"column:depends_on_question_id"`
	DependsOnOptionId   uint   `gorm:"column:depends_on_option_id"`
}
