package service

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
	"gorm.io/gorm"
)

var (
	ErrorSelectedUsersIdsFieldIsRequired = errors.New("field selected_users_ids is required")
	ErrorPermissionsFieldIsRequired      = errors.New("field permissions is required")
)

type IRbacService interface {
	GetAllPermissions() []string
	GivePermissions(giverUserID uint, receiverUserId uint, questionnaireID uint, permissions []PermissionType) error
	CanDo(userID uint, questionnaireID uint, permissionName string) (bool, error)
	RevokePermission(revokerUserID uint, targetUserID uint, questionnaireID uint, permissionName string) error
	HasPermission(userID uint, questionnaireID uint, permissionName string) (bool, error)
	MakeFakeUser() (model.User, error)
	MakeFakeQuestionnaire(userID uint) (model.Questionnaire, error)
	GetUser(userID uint) (model.User, error)
	GetUserRolesWithPermissions(userID uint) ([]domain_repository.RoleWithPermissions, error)
	CanDoAsSuperadmin(userID uint, permissionName string) (bool, error)
	GetUsersIDsWithVisibleAnswers(questionnaireID uint, userID uint) ([]uint, error)
}

type RbacService struct {
	repo domain_repository.IRbacRepository
}

func NewRbacService(repo domain_repository.IRbacRepository) *RbacService {
	return &RbacService{
		repo: repo,
	}
}

type PermissionType struct {
	ID               uint
	Name             string `json:"name"`
	ExpireDate       string `json:"expire_date"`
	SelectedUsersIds []uint `json:"selected_users_ids"`
}

func (o *RbacService) GetAllPermissions() []string {
	permissions := o.repo.GetPermissions()
	permissionNames := make([]string, 0, len(permissions))
	for _, permission := range permissions {
		permissionNames = append(permissionNames, permission.Name)
	}
	return permissionNames
}

func (o *RbacService) GivePermissions(giverUserID uint, receiverUserId uint, questionnaireID uint, permissions []PermissionType) error {
	if err := o.repo.IsUserExist(giverUserID); err != nil {
		return err
	}

	if err := o.repo.IsQuestionnaireExist(questionnaireID); err != nil {
		return err
	}

	can, err := o.CanDo(giverUserID, questionnaireID, model.PERMISSION_GIVE_OR_TAKE_PERMISSION)
	if err != nil {
		return err
	}

	if !can {
		return model.ErrorUserCanNotDoThis
	}

	if err := o.repo.IsUserExist(receiverUserId); err != nil {
		return err
	}

	roleName := ""
	for index, permission := range permissions {
		id, err := o.repo.FindPermission(permission.Name)
		if err != nil {
			return err
		}

		if permission.Name == model.PERMISSION_SEE_SELECTED_USERS_ANSWERS && len(permission.SelectedUsersIds) == 0 {
			return ErrorSelectedUsersIdsFieldIsRequired
		}

		roleName = makeAbbreviation(permission.Name) + "_"
		permissions[index].ID = id
	}

	roleName = fmt.Sprintf("%s%d_%d_%d_%d", roleName, questionnaireID, receiverUserId, giverUserID, MakeRandomNumber(10000000))

	err = o.repo.Transaction(func(tx *gorm.DB) error {
		roleId, err := o.repo.MakeNewRole(tx, roleName)
		if err != nil {
			return err // rollback will be triggered
		}

		for _, permission := range permissions {
			o.RevokePermission(giverUserID, receiverUserId, questionnaireID, permission.Name)

			rolePermissionID, err := o.repo.MakeRolePermission(tx, roleId, questionnaireID, permission.ID, ParseAsSqlNullTime(permission.ExpireDate))
			if err != nil {
				return err // rollback will be triggered
			}

			if permission.Name == model.PERMISSION_SEE_SELECTED_USERS_ANSWERS {
				err = o.repo.InsertUsersWithVisibleAnswers(tx, rolePermissionID, permission.SelectedUsersIds)
				if err != nil {
					return err
				}
			}
		}

		if err := o.repo.MakeRoleUser(tx, roleId, receiverUserId); err != nil {
			return err // rollback will be triggered
		}

		// Commit is automatically called if no error is returned
		return nil
	})

	return err
}

func ParseAsSqlNullTime(stringDate string) sql.NullTime {
	if stringDate != "" {

		location, err := time.LoadLocation("Asia/Tehran")
		if err != nil {
			return sql.NullTime{Valid: false}
		}

		t, err := time.ParseInLocation("2006-01-02 15:04:05", stringDate, location)
		if err != nil {
			return sql.NullTime{Valid: false}
		}

		return sql.NullTime{Time: t, Valid: true}
	} else {
		return sql.NullTime{Valid: false}
	}
}

func makeAbbreviation(name string) string {
	result := ""
	parts := strings.Split(name, "_")
	for _, part := range parts {
		result = fmt.Sprintf("%s%s", result, string(part[0]))
	}
	return result
}

func (o *RbacService) CanDo(userID uint, questionnaireID uint, permissionName string) (bool, error) {
	err := o.repo.IsOwnerOfQuestionnaire(userID, questionnaireID)
	if err == nil {
		return true, nil
	}

	return o.HasPermission(userID, questionnaireID, permissionName)
}

func (o *RbacService) RevokePermission(revokerUserID uint, targetUserID uint, questionnaireID uint, permissionName string) error {
	if err := o.repo.IsUserExist(revokerUserID); err != nil {
		return err
	}

	if err := o.repo.IsQuestionnaireExist(questionnaireID); err != nil {
		return err
	}

	can, err := o.CanDo(revokerUserID, questionnaireID, model.PERMISSION_GIVE_OR_TAKE_PERMISSION)
	if err != nil {
		return err
	}

	if !can {
		return model.ErrorUserCanNotDoThis
	}

	targetUser, err := o.repo.GetUserWithRoles(targetUserID)
	if err != nil {
		return err
	}

	permissionID, err := o.repo.FindPermission(permissionName)
	if err != nil {
		return err
	}

	for _, role := range targetUser.Roles {
		err = o.repo.DeleteRolePermissions(role.ID, questionnaireID, permissionID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *RbacService) HasPermission(userID uint, questionnaireID uint, permissionName string) (bool, error) {
	user, err := o.repo.GetUserWithRoles(userID)
	if err != nil {
		return false, err
	}

	if err := o.repo.IsQuestionnaireExist(questionnaireID); err != nil {
		return false, err
	}

	permissionID, err := o.repo.FindPermission(permissionName)
	if err != nil {
		return false, err
	}

	superadmin, err := o.repo.GetSuperadmin(userID)
	if err != nil {
		return false, err
	}
	if superadmin.ID != 0 {
		hasSuperadminPermission, err := o.repo.FindSuperadminPermission(superadmin.ID, permissionID)
		if err != nil {
			return false, err
		}
		if hasSuperadminPermission {
			return true, nil
		}
	}

	foundRolePermission := false
	for _, role := range user.Roles {
		found, err := o.repo.HasRolePermission(role.ID, questionnaireID, permissionID)
		if err != nil {
			return false, err
		}
		if found {
			foundRolePermission = true
			break
		}
	}

	if foundRolePermission {
		return true, nil
	}

	return false, nil
}

func (o *RbacService) MakeFakeUser() (model.User, error) {
	index := MakeRandomNumber(400)
	return o.repo.MakeUser(model.User{
		FirstName:    "ufn" + strconv.Itoa(index),
		LastName:     "u2ln" + strconv.Itoa(index),
		Email:        "u2@email" + strconv.Itoa(index),
		NationalCode: makeRandomNationalCode(),
		Password:     "111111111" + strconv.Itoa(index),
	})
}

func makeRandomNationalCode() string {
	min := int64(1000000000)
	max := int64(9999999999)
	randomNumber := min + rand.Int63n(max-min+1)
	return fmt.Sprintf("%d", randomNumber)
}

func (o *RbacService) MakeFakeQuestionnaire(userID uint) (model.Questionnaire, error) {
	return o.repo.MakeQuestionnaire(model.Questionnaire{
		OwnerID:                    model.UserID(userID),
		Status:                     model.QuestionnaireStatusOpen,
		MaxMinutesToResponse:       MakeRandomNumber(100),
		MaxMinutesToChangeAnswer:   MakeRandomNumber(100),
		MaxMinutesToGivebackAnswer: MakeRandomNumber(100),
		RandomOrSequential:         model.QuestionnaireTypeSequential,
		CanBackToPreviousQuestion:  false,
		Title:                      "test q" + strconv.Itoa(MakeRandomNumber(200000000)),
		MaxAllowedSubmissionsCount: 0,
		AnswersVisibleFor:          model.QuestionnaireVisibilityEverybody,
		CreatedAt:                  time.Now(),
	})
}

func MakeRandomNumber(max int) int {
	return rand.Intn(max)
}

func (o *RbacService) GetUser(userID uint) (model.User, error) {
	return o.repo.GetUserWithQuestionnaires(userID)
}

func (o *RbacService) GetUserRolesWithPermissions(userID uint) ([]domain_repository.RoleWithPermissions, error) {
	return o.repo.GetUserRolesWithPermissions(userID)
}

func (o *RbacService) CanDoAsSuperadmin(userID uint, permissionName string) (bool, error) {
	if err := o.repo.IsUserExist(userID); err != nil {
		return false, err
	}

	superadmin, err := o.repo.GetSuperadmin(userID)
	if err != nil {
		return false, err
	}
	if superadmin.ID == 0 {
		return false, model.ErrorUserCanNotDoThis
	}

	permissionID, err := o.repo.FindPermission(permissionName)
	if err != nil {
		return false, err
	}

	hasSuperadminPermission, err := o.repo.FindSuperadminPermission(superadmin.ID, permissionID)
	if err != nil {
		return false, err
	}

	if !hasSuperadminPermission {
		return false, model.ErrorUserCanNotDoThis
	}

	return true, nil
}

func (o *RbacService) GetUsersIDsWithVisibleAnswers(questionnaireID uint, userID uint) ([]uint, error) {
	roles, err := o.repo.GetUserRoles(userID)
	if err != nil {
		return nil, err
	}

	permissionID, err := o.repo.FindPermission(model.PERMISSION_SEE_SELECTED_USERS_ANSWERS)
	if err != nil {
		return nil, err
	}

	rolePermission := new(model.RolePermission)
	for _, role := range roles {
		rolePermission, err = o.repo.FindRolePermission(role.ID, questionnaireID, permissionID)
		if err != nil {
			return nil, err
		}

		if rolePermission.ID != 0 {
			break
		}
	}

	if rolePermission.ID == 0 {
		return nil, model.ErrorNotHavePermission
	}

	usersIDs, err := o.repo.FindUsersWithVisibleAnswers(rolePermission.ID)
	if err != nil {
		return nil, err
	}

	return usersIDs, nil
}
