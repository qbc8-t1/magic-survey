package service

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	repository "github.com/QBC8-Team1/magic-survey/persistance"
	"gorm.io/gorm"
)

var ErrorSelectedUsersIdsFieldIsRequired = errors.New("selected users ids field is required")

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
	Name             string `json:"name"`
	ExpireDate       string `json:"expire_date"`
	SelectedUsersIds []uint `json:"selected_users_ids"`
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

	can, err := o.CanDo(giverUserID, questionnaireID, model.PERMISSION_QUESTIONNAIRE_GIVE_OR_TAKE_PERMISSION)
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

		if permission.Name == model.PERMISSION_QUESTIONNAIRE_SEE_SELECTED_USERS_ANSWERS && len(permission.SelectedUsersIds) == 0 {
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
	// Todo - check if user has super admin role

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

	can, err := o.CanDo(revokerUserID, questionnaireID, model.PERMISSION_QUESTIONNAIRE_GIVE_OR_TAKE_PERMISSION)
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

	foundRolePermission := false
	for _, role := range user.Roles {
		found, err := o.repo.FindRolePermission(role.ID, questionnaireID, permissionID)
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

func (o *RbacService) GetUsersWithVisibleAnswers(questionnaireID uint, userID uint) []uint {
	return nil
}
