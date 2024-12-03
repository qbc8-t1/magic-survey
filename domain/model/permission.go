package model

import "errors"

// Permission represents the permissions table
type Permission struct {
	ID            uint   `gorm:"primaryKey"`
	Name          string `gorm:"unique;size:255"`
	Description   string `gorm:"size:500"`
	ForSuperadmin bool   `gorm:"default:false"`
}

const (
	PERMISSION_QUESTIONNAIRE_VIEW                       = "questionnaire_view"
	PERMISSION_QUESTIONNAIRE_ANSWER                     = "questionnaire_answer"
	PERMISSION_QUESTIONNAIRE_EDIT                       = "questionnaire_edit"
	PERMISSION_QUESTIONNAIRE_REPORTS_VIEW               = "questionnaire_reports_view"
	PERMISSION_QUESTIONNAIRE_GIVE_OR_TAKE_PERMISSION    = "questionnaire_give_or_take_permission"
	PERMISSION_QUESTIONNAIRE_SEE_SELECTED_USERS_ANSWERS = "questionnaire_see_selected_users_answers"
	PERMISSION_LIMIT_USERS_QUESTIONNAIRES_COUNT         = "limit_users_questionnaire_count"
	PERMISSION_MAKE_NEW_SUPERADMIN                      = "make_new_superadmin"
)

var PermissionsForUser = []string{
	PERMISSION_QUESTIONNAIRE_VIEW,
	PERMISSION_QUESTIONNAIRE_ANSWER,
	PERMISSION_QUESTIONNAIRE_EDIT,
	PERMISSION_QUESTIONNAIRE_REPORTS_VIEW,
	PERMISSION_QUESTIONNAIRE_GIVE_OR_TAKE_PERMISSION,
	PERMISSION_QUESTIONNAIRE_SEE_SELECTED_USERS_ANSWERS,
}

var PermissionsForSuperadmin = []string{
	PERMISSION_LIMIT_USERS_QUESTIONNAIRES_COUNT,
	PERMISSION_MAKE_NEW_SUPERADMIN,
}

var (
	ErrorNotFoundUser               = errors.New("user not found")
	ErrorNotFoundQuestionnaire      = errors.New("questionnaire not found")
	ErrorIsNotOwnerOfQuestionnaire  = errors.New("is not owner of questionnaire")
	ErrorNotHavePermission          = errors.New("user doesn't have permission")
	ErrorNotFoundPermission         = errors.New("permission not found")
	ErrorUserCanNotDoThis           = errors.New("user can not do this")
	ErrorCanNotGiveOrTakePermission = errors.New("user can not give or take permission")
)
