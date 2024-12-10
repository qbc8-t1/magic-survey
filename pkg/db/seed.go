package db

import (
	"strings"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	"github.com/QBC8-Team1/magic-survey/pkg/utils"
	"gorm.io/gorm"
)

func seed(db *gorm.DB) error {
	if err := seedPermissions(db, model.PermissionsForUser, false); err != nil {
		return err
	}

	if err := seedPermissions(db, model.PermissionsForSuperadmin, true); err != nil {
		return err
	}

	if err := makeSuperadminUser(db); err != nil {
		return err
	}
	return nil
}

func seedPermissions(db *gorm.DB, permissions []model.PermissionName, forSuperadmin bool) error {
	for _, permissionName := range permissions {
		err := db.Create(&model.Permission{
			Name:          permissionName,
			Description:   strings.Replace(string(permissionName), "_", " ", -1),
			ForSuperadmin: forSuperadmin,
		}).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func makeSuperadminUser(db *gorm.DB) error {
	// make a user
	pass, err := utils.HashPassword("password")
	if err != nil {
		return err
	}

	user := model.User{
		FirstName: "Super",
		LastName:  "Admin",
		Email:     "super@admin.com",
		Password:  pass,
		IsActive:  true,
	}

	err = db.Create(&user).Error
	if err != nil {
		return err
	}

	// add it to superadmin table
	superadmin := model.Superadmin{
		UserID: user.ID,
	}
	err = db.Create(&superadmin).Error
	if err != nil {
		return err
	}

	// make superadmin_permission records for that user with all permissions
	permissions := []model.Permission{}
	err = db.Find(&permissions).Error
	if err != nil {
		return err
	}

	for _, permission := range permissions {
		err := db.Create(&model.SuperadminPermission{
			SuperadminID: superadmin.ID,
			PermissionID: permission.ID,
		}).Error

		if err != nil {
			return err
		}
	}
	return nil
}
