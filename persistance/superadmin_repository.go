package repository

import (
	"errors"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
	"gorm.io/gorm"
)

type SuperadminRepository struct {
	db *gorm.DB
}

func NewSuperadminRepository(db *gorm.DB) domain_repository.ISuperadminRepository {
	return &SuperadminRepository{db: db}
}

func (s *SuperadminRepository) MakeSuperadminRole(tx *gorm.DB, userID uint, giverUserID uint) (uint, error) {
	superadmin := model.Superadmin{
		UserID:    userID,
		GrantedBy: &giverUserID,
	}
	err := tx.Create(&superadmin).Error
	if err != nil {
		return 0, err
	}
	return superadmin.ID, nil
}

func (s *SuperadminRepository) MakeSuperadminPermission(tx *gorm.DB, superadminID uint, permissionID uint) error {
	var exists bool
	err := tx.Model(&model.SuperadminPermission{}).
		Select("count(*) > 0").
		Where("superadmin_id = ? AND permission_id = ?", superadminID, permissionID).
		Find(&exists).Error

	if err != nil {
		return err
	}

	if !exists {
		return tx.Create(&model.SuperadminPermission{
			SuperadminID: superadminID,
			PermissionID: permissionID,
		}).Error
	}

	return nil
}

func (s *SuperadminRepository) GetSuperadmin(userID uint) (model.Superadmin, error) {
	superadmin := new(model.Superadmin)
	err := s.db.First(superadmin, "user_id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Superadmin{}, nil
		}

		return model.Superadmin{}, err
	}

	return *superadmin, nil
}

func (s *SuperadminRepository) FindSuperadminPermission(superadminID uint, permissionID uint) (bool, error) {
	err := s.db.First(&model.SuperadminPermission{}, "superadmin_id = ? and permission_id = ?", superadminID, permissionID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (s *SuperadminRepository) FindPermission(permissionName string) (uint, error) {
	var permission model.Permission
	err := s.db.First(&permission, "name = ?", permissionName).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, model.ErrorNotFoundPermission
		}
		return 0, err
	}

	return permission.ID, nil
}

func (s *SuperadminRepository) Transaction(f func(tx *gorm.DB) error) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		return f(tx)
	})
}

func (s *SuperadminRepository) IsUserExist(userID uint) error {
	var user model.User
	err := s.db.First(&user, "id = ?", userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.ErrorNotFoundUser
	}
	return err
}
