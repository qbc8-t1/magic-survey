package service

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	repository "github.com/QBC8-Team1/magic-survey/persistance"
	"gorm.io/gorm"
)

var ErrorSelectedUsersSliceIsEmpty = errors.New("selected users slice is empty")

type RbacService struct {
	repo *repository.RbacRepo
}

func NewRbacService(repo *repository.RbacRepo) *RbacService {
	return &RbacService{
		repo: repo,
	}
}

type PermissionType struct {
	ID               uint
	Name             string       `json:"name"`
	ExpireDate       sql.NullTime `json:"expire_date"`
	SelectedUsersIds []uint       `json:"selected_users_ids"`
}

func (o *RbacService) GetAllPermissions() []string {
	permissions := o.repo.GetAllPermissions()
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

	if err := o.CanGiveOrTakePermission(giverUserID, questionnaireID); err != nil {
		return err
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

		if permission.Name == model.PERMISSION_QUESTIONNAIRE_SEE_SELECTED_USERS_ANSWERS && len(permission.SelectedUsersIds) == 0 {
			return ErrorSelectedUsersSliceIsEmpty
		}

		roleName = makeAbbreviation(permission.Name) + "_"
		permissions[index].ID = id
	}

	roleName = fmt.Sprintf("%s%d_%d_%d", roleName, questionnaireID, receiverUserId, giverUserID)

	err := o.repo.Transaction(func(tx *gorm.DB) error {
		roleId, err := o.repo.MakeNewRole(tx, roleName)
		if err != nil {
			return err // rollback will be triggered
		}

		for _, permission := range permissions {
			o.TakePermission(receiverUserId, questionnaireID, permission.Name)
			rolePermissionID, err := o.repo.MakeRolePermission(tx, roleId, questionnaireID, permission.ID, permission.ExpireDate)
			if err != nil {
				return err // rollback will be triggered
			}

			if permission.Name == model.PERMISSION_QUESTIONNAIRE_SEE_SELECTED_USERS_ANSWERS {
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

func makeAbbreviation(name string) string {
	result := ""
	parts := strings.Split(name, "_")
	for _, part := range parts {
		result = fmt.Sprintf("%s%s", result, string(part[0]))
	}
	return result
}

func (o *RbacService) CanGiveOrTakePermission(userID uint, questionnaireID uint) error {
	err := o.repo.IsOwnerOfQuestionnaire(userID, questionnaireID)
	if err == nil {
		return nil
	}

	return o.HasPermission(userID, questionnaireID, model.PERMISSION_QUESTIONNAIRE_GIVE_OR_TAKE_PERMISSION)
}

func (o *RbacService) TakePermission(userID uint, questionnaireID uint, permissionName string) error {
	// check if user is exist
	// check if questionnaire is exist
	// check if all permission is exist
	// loop on all roles of user
	// loop on all permissions of roles
	// if permission was for questionnaireID
	// delete permission
	return nil
}

func (o *RbacService) HasPermission(userID uint, questionnaireID uint, permissionName string) error {
	// 	// check if user is exist
	// 	// check if questionnaire is exist
	// 	// check if all permission is exist
	// 	// loop on all roles of user
	// 	// loop on all permissions of roles
	// 	// if permission was for questionnaireID and expiretime wasn't reached
	// 	// return true

	// return false
	return model.ErrorNotHavePermission
}

func (o *RbacService) GetUsersWithVisibleAnswers(questionnaireID uint, userID uint) []uint {
	return nil
}
