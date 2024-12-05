package service

import (
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
	"gorm.io/gorm"
)

type SuperadminService struct {
	repo domain_repository.ISuperadminRepository
}

func NewSuperadminService(repo domain_repository.ISuperadminRepository) *SuperadminService {
	return &SuperadminService{
		repo: repo,
	}
}

func (o *SuperadminService) MakeSuperadmin(giverUserID uint, userID uint, permissionsNames []string) error {
	if err := o.repo.IsUserExist(userID); err != nil {
		return err
	}

	if len(permissionsNames) == 0 {
		return ErrorPermissionsFieldIsRequired
	}

	permissionIDs := make([]uint, 0, len(permissionsNames))
	for _, permissionName := range permissionsNames {
		id, err := o.repo.FindPermission(permissionName)
		if err != nil {
			return err
		}

		permissionIDs = append(permissionIDs, id)
	}

	oldSuperadmin, err := o.repo.GetSuperadmin(userID)
	if err != nil {
		return err
	}

	err = o.repo.Transaction(func(tx *gorm.DB) error {
		// add permissions to exist superadmin record (if exist)
		userSuperadminID := oldSuperadmin.ID
		if userSuperadminID == 0 {
			userSuperadminID, err = o.repo.MakeSuperadminRole(tx, userID, giverUserID)
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
