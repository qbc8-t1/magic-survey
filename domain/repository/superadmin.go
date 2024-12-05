package domain_repository

import (
	"github.com/QBC8-Team1/magic-survey/domain/model"
	"gorm.io/gorm"
)

type ISuperadminRepository interface {
	MakeSuperadminRole(tx *gorm.DB, userID uint, giverUserID uint) (uint, error)
	MakeSuperadminPermission(tx *gorm.DB, superadminID uint, permissionID uint) error
	GetSuperadmin(userID uint) (model.Superadmin, error)
	FindSuperadminPermission(superadminID uint, permissionID uint) (bool, error)
	FindPermission(permissionName string) (uint, error)
	Transaction(f func(tx *gorm.DB) error) error
	IsUserExist(userID uint) error
}
