package domain_repository

import (
	"database/sql"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	"gorm.io/gorm"
)

type RoleWithPermissions struct {
	ID          uint                    `json:"id"`
	Name        string                  `json:"name"`
	Permissions []PermissionWithDetails `json:"permissions"`
}

type PermissionWithDetails struct {
	PermissionID    uint         `json:"permission_id"`
	Name            string       `json:"name"`
	QuestionnaireID uint         `json:"questionnaire_id"`
	ExpireAt        sql.NullTime `json:"expire_at"`
}

type IRbacRepository interface {
	GetPermissions() []model.Permission
	GetUserWithRoles(userID uint) (model.User, error)
	IsUserExist(userID uint) error
	IsQuestionnaireExist(questionnaireID model.QuestionnaireID) error
	IsOwnerOfQuestionnaire(userID uint, questionnaireID model.QuestionnaireID) error
	FindPermission(permissionName model.PermissionName) (uint, error)
	FindPermissionForRegularUsers(permissionName model.PermissionName) (uint, error)
	GetUserRoles(userID uint) ([]model.Role, error)
	HasPermission(questionnaireID model.QuestionnaireID, permissionID uint) error
	MakeNewRole(tx *gorm.DB, name string) (uint, error)
	MakeRolePermission(tx *gorm.DB, roleID uint, questionnaireID model.QuestionnaireID, permissionID uint, expireDate sql.NullTime) (uint, error)
	InsertUsersWithVisibleAnswers(tx *gorm.DB, rolePermissionID uint, selectedUsersIDs []uint) error
	MakeRoleUser(tx *gorm.DB, roleID uint, userID uint) error
	Transaction(f func(tx *gorm.DB) error) error
	MakeUser(user model.User) (model.User, error)
	MakeQuestionnaire(questionnaire model.Questionnaire) (model.Questionnaire, error)
	GetUserWithQuestionnaires(userID uint) (model.User, error)
	GetUserRolesWithPermissions(userID uint) ([]RoleWithPermissions, error)
	DeleteRolePermissions(roleID uint, questionnaireID model.QuestionnaireID, permissionID uint) error
	HasRolePermission(roleID uint, questionnaireID model.QuestionnaireID, permissionID uint) (bool, error)
	GetSuperadmin(userID uint) (model.Superadmin, error)
	FindSuperadminPermission(superadminID uint, permissionID uint) (bool, error)
	FindRolePermission(roleID uint, questionnaireID model.QuestionnaireID, permissionID uint) (*model.RolePermission, error)
	FindUsersWithVisibleAnswers(rolePermissionID uint) ([]uint, error)
	GetQuestionByAnswerID(answerID model.AnswerID) (model.Question, error)
	GetQuestionByOptionID(optionID model.OptionID) (model.Question, error)
	GetQuestionByID(questionID model.QuestionID) (model.Question, error)
	GetAnswersForQuestionnaire(questionnaireID model.QuestionnaireID) []AnswersResult
}

type AnswersResult struct {
	QuestionID   uint
	QuestionText string
	SubmissionID uint
	UserID       uint
	AnswerText   string
	OptionID     uint
}
