package model

import "errors"

// Permission represents the permissions table
type Permission struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"unique;size:255"`
	Description string `gorm:"size:500"`
}

const (
	PERMISSION_QUESTIONNAIRE_VIEW                       = "questionnaire_view"
	PERMISSION_QUESTIONNAIRE_ANSWER                     = "questionnaire_answer"
	PERMISSION_QUESTIONNAIRE_EDIT                       = "questionnaire_edit"
	PERMISSION_QUESTIONNAIRE_REPORTS_VIEW               = "questionnaire_reporst_view"
	PERMISSION_QUESTIONNAIRE_GIVE_OR_TAKE_PERMISSION    = "questionnaire_give_or_take_permission"
	PERMISSION_QUESTIONNAIRE_SEE_SELECTED_USERS_ANSWERS = "questionnaire_see_selected_users_answers"
)

var Permissions = []string{
	PERMISSION_QUESTIONNAIRE_VIEW,
	PERMISSION_QUESTIONNAIRE_ANSWER,
	PERMISSION_QUESTIONNAIRE_EDIT,
	PERMISSION_QUESTIONNAIRE_REPORTS_VIEW,
	PERMISSION_QUESTIONNAIRE_GIVE_OR_TAKE_PERMISSION,
	PERMISSION_QUESTIONNAIRE_SEE_SELECTED_USERS_ANSWERS,
}

var (
	ErrorNotFoundUser               = errors.New("user not found")
	ErrorNotFoundQuestionnaire      = errors.New("questionnaire not found")
	ErrorIsNotOwnerOfQuestionnaire  = errors.New("is not owner of questionnaire")
	ErrorNotHavePermission          = errors.New("not have permission")
	ErrorNotFoundPermission         = errors.New("not found permission")
	ErrorCanNotGiveOrTakePermission = errors.New("user can not give or take permission")
)
