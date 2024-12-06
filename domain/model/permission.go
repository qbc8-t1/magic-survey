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
	PERMISSION_QUESTIONNAIRE_VIEW = "questionnaire_view"
	PERMISSION_QUESTIONNAIRE_EDIT = "questionnaire_edit"

	PERMISSION_REPORTS_VIEW = "reports_view"

	PERMISSION_GIVE_OR_TAKE_PERMISSION = "give_or_take_permission"

	PERMISSION_SEE_SELECTED_USERS_ANSWERS = "see_selected_users_answers"

	PERMISSION_QUESTION_CREATE = "question_create"
	PERMISSION_QUESTION_UPDATE = "question_update"
	PERMISSION_QUESTION_VIEW   = "question_view"
	PERMISSION_QUESTION_DELETE = "question_delete"

	PERMISSION_QUESTION_OPTION_CREATE = "question_option_create"
	PERMISSION_QUESTION_OPTION_UPDATE = "question_option_update"
	PERMISSION_QUESTION_OPTION_VIEW   = "question_option_view"
	PERMISSION_QUESTION_OPTION_DELETE = "question_option_delete"

	PERMISSION_ANSWER_GIVE     = "answer_give"
	PERMISSION_ANSWER_GIVEBACK = "answer_giveback"
	PERMISSION_ANSWER_UPDATE   = "answer_update"

	PERMISSION_LIMIT_USERS_QUESTIONNAIRES_COUNT = "limit_users_questionnaire_count"
	PERMISSION_MAKE_NEW_SUPERADMIN              = "make_new_superadmin"
)

var PermissionsForUser = []string{
	PERMISSION_QUESTIONNAIRE_VIEW,
	PERMISSION_QUESTIONNAIRE_EDIT,
	PERMISSION_REPORTS_VIEW,
	PERMISSION_GIVE_OR_TAKE_PERMISSION,
	PERMISSION_SEE_SELECTED_USERS_ANSWERS,
	PERMISSION_QUESTION_CREATE,
	PERMISSION_QUESTION_UPDATE,
	PERMISSION_QUESTION_VIEW,
	PERMISSION_QUESTION_DELETE,
	PERMISSION_QUESTION_OPTION_CREATE,
	PERMISSION_QUESTION_OPTION_UPDATE,
	PERMISSION_QUESTION_OPTION_VIEW,
	PERMISSION_QUESTION_OPTION_DELETE,
	PERMISSION_ANSWER_GIVE,
	PERMISSION_ANSWER_GIVEBACK,
	PERMISSION_ANSWER_UPDATE,
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
