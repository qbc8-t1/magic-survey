package model

import "errors"

// Permission represents the permissions table
type Permission struct {
	ID            uint           `gorm:"primaryKey"`
	Name          PermissionName `gorm:"unique;size:255"`
	Description   string         `gorm:"size:500"`
	ForSuperadmin bool           `gorm:"default:false"`
}

type PermissionName string

const (
	PERMISSION_QUESTIONNAIRE_VIEW   PermissionName = "questionnaire_view"
	PERMISSION_QUESTIONNAIRE_EDIT   PermissionName = "questionnaire_edit"
	PERMISSION_QUESTIONNAIRE_UPDATE PermissionName = "questionnaire_update"
	PERMISSION_QUESTIONNAIRE_DELETE PermissionName = "questionnaire_delete"
	PERMISSION_QUESTIONNAIRE_CANCEL PermissionName = "questionnaire_cancel"
	PERMISSION_QUESTIONNAIRE_CLOSE  PermissionName = "questionnaire_close"

	PERMISSION_REPORTS_VIEW PermissionName = "reports_view"

	PERMISSION_GIVE_OR_TAKE_PERMISSION PermissionName = "give_or_take_permission"

	PERMISSION_SEE_SELECTED_USERS_ANSWERS PermissionName = "see_selected_users_answers"

	PERMISSION_QUESTION_CREATE PermissionName = "question_create"
	PERMISSION_QUESTION_UPDATE PermissionName = "question_update"
	PERMISSION_QUESTION_VIEW   PermissionName = "question_view"
	PERMISSION_QUESTION_DELETE PermissionName = "question_delete"

	PERMISSION_OPTION_CREATE PermissionName = "option_create"
	PERMISSION_OPTION_UPDATE PermissionName = "option_update"
	PERMISSION_OPTION_VIEW   PermissionName = "option_view"
	PERMISSION_OPTION_DELETE PermissionName = "option_delete"

	PERMISSION_ANSWER_GIVE     PermissionName = "answer_give"
	PERMISSION_ANSWER_GIVEBACK PermissionName = "answer_giveback"
	PERMISSION_ANSWER_UPDATE   PermissionName = "answer_update"
	PERMISSION_ANSWER_CREATE   PermissionName = "answer_create"
	PERMISSION_ANSWER_DELETE   PermissionName = "answer_delete"

	PERMISSION_LIMIT_USER_QUESTIONNAIRES_COUNT PermissionName = "limit_user_questionnaires_count"
	PERMISSION_MAKE_NEW_SUPERADMIN             PermissionName = "make_new_superadmin"
)

var PermissionsForUser = []PermissionName{
	PERMISSION_QUESTIONNAIRE_VIEW,
	PERMISSION_QUESTIONNAIRE_EDIT,
	PERMISSION_QUESTIONNAIRE_UPDATE,
	PERMISSION_QUESTIONNAIRE_DELETE,
	PERMISSION_REPORTS_VIEW,
	PERMISSION_GIVE_OR_TAKE_PERMISSION,
	PERMISSION_SEE_SELECTED_USERS_ANSWERS,
	PERMISSION_QUESTION_CREATE,
	PERMISSION_QUESTION_UPDATE,
	PERMISSION_QUESTION_VIEW,
	PERMISSION_QUESTION_DELETE,
	PERMISSION_OPTION_CREATE,
	PERMISSION_OPTION_UPDATE,
	PERMISSION_OPTION_VIEW,
	PERMISSION_OPTION_DELETE,
	PERMISSION_ANSWER_GIVE,
	PERMISSION_ANSWER_GIVEBACK,
	PERMISSION_ANSWER_UPDATE,
	PERMISSION_ANSWER_CREATE,
	PERMISSION_ANSWER_DELETE,
}

var PermissionsForSuperadmin = []PermissionName{
	PERMISSION_LIMIT_USER_QUESTIONNAIRES_COUNT,
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
