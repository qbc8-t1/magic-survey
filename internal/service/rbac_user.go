package service

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	repository "github.com/QBC8-Team1/magic-survey/persistance"
	"gorm.io/gorm"
)

func (o *RbacService) MakeFakeUser() (model.User, error) {
	index := MakeRandomNumber(400)
	return o.repo.MakeUser(model.User{
		FirstName:    "ufn" + strconv.Itoa(index),
		LastName:     "u2ln" + strconv.Itoa(index),
		Email:        "u2@email" + strconv.Itoa(index),
		NationalCode: strconv.Itoa(MakeRandomNumber(9999999999)),
		Password:     "111111111" + strconv.Itoa(index),
	})
}

func (o *RbacService) MakeFakeQuestionnaire(userID uint) (model.Questionnaire, error) {
	return o.repo.MakeQuestionnaire(model.Questionnaire{
		OwnerID:                    uint(userID),
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

func (o *RbacService) GetUserRolesWithPermissions(userID uint) ([]repository.RoleWithPermissions, error) {
	return o.repo.GetUserRolesWithPermissions(userID)
}

func (o *RbacService) MakeSuperadmin(giverUserID uint, userID uint, permissionNames []string) error {
	if err := o.repo.IsUserExist(giverUserID); err != nil {
		return err
	}

	superadmin, err := o.repo.IsSuperadmin(giverUserID)
	if err != nil {
		return err
	}
	if superadmin.ID == 0 {
		return model.ErrorUserCanNotDoThis
	}

	permissionID, err := o.repo.FindPermission(model.PERMISSION_MAKE_NEW_SUPERADMIN)
	if err != nil {
		return err
	}

	hasSuperadminPermission, err := o.repo.FindSuperadminPermission(superadmin.ID, permissionID)
	if err != nil {
		return err
	}
	if !hasSuperadminPermission {
		return model.ErrorUserCanNotDoThis
	}

	if err := o.repo.IsUserExist(userID); err != nil {
		return err
	}

	if len(permissionNames) == 0 {
		return ErrorPermissionNamesFieldIsRequired
	}

	permissionIDs := make([]uint, 0, len(permissionNames))
	for _, permissionName := range permissionNames {
		id, err := o.repo.FindPermission(permissionName)
		if err != nil {
			return err
		}

		permissionIDs = append(permissionIDs, id)
	}

	oldSuperadmin, err := o.repo.IsSuperadmin(userID)
	if err != nil {
		return err
	}

	err = o.repo.Transaction(func(tx *gorm.DB) error {
		userSuperadminID := oldSuperadmin.ID
		if userSuperadminID == 0 {
			userSuperadminID, err = o.repo.GiveSuperadminRole(tx, userID, giverUserID)
			if err != nil {
				return err
			}
		}

		for _, permissionID := range permissionIDs {
			err := o.repo.MakeSuperadminPermission(tx, userSuperadminID, permissionID)

			if err != nil {
				return err
			}
		}

		return nil
	})

	return nil
}
