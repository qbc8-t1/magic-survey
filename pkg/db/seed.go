package db

import (
	"strings"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	"gorm.io/gorm"
)

func seed(db *gorm.DB) error {
	seedPermissions(db)
	return nil
}

func seedPermissions(db *gorm.DB) error {
	for _, permissionName := range model.Permissions {
		err := db.Create(&model.Permission{
			Name:        permissionName,
			Description: strings.Replace(permissionName, "_", " ", -1),
		}).Error
		if err != nil {
			return err
		}
	}
	return nil
}
