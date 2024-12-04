package repository

import (
	"database/sql"
	"errors"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	domain_repository "github.com/QBC8-Team1/magic-survey/domain/repository"
	"gorm.io/gorm"
)

type RbacRepo struct {
	db *gorm.DB
}

func NewRbacRepository(db *gorm.DB) domain_repository.IRbacRepository {
	return &RbacRepo{
		db: db,
	}
}

func (rr *RbacRepo) GetPermissions() []model.Permission {
	var permissions []model.Permission
	rr.db.Find(&permissions)
	return permissions
}

func (rr *RbacRepo) IsUserExist(userID uint) error {
	var user model.User
	err := rr.db.First(&user, "id = ?", userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.ErrorNotFoundUser
	}
	return err
}

func (rr *RbacRepo) GetUserWithRoles(userID uint) (model.User, error) {
	user := new(model.User)
	err := rr.db.Preload("Roles").First(user, "id = ?", userID).Error
	return *user, err
}

func (rr *RbacRepo) IsQuestionnaireExist(questionnaireID uint) error {
	var questionnaire model.Questionnaire
	err := rr.db.First(&questionnaire, "id = ?", questionnaireID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.ErrorNotFoundQuestionnaire
	}
	return err
}

func (rr *RbacRepo) IsOwnerOfQuestionnaire(userID uint, questionnaireID uint) error {
	var questionnaire model.Questionnaire
	err := rr.db.First(&questionnaire, "id = ? and owner_id = ?", questionnaireID, userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.ErrorIsNotOwnerOfQuestionnaire
	}
	return err
}

func (rr *RbacRepo) FindPermission(permissionName string) (uint, error) {
	var permission model.Permission
	err := rr.db.First(&permission, "name = ?", permissionName).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, model.ErrorNotFoundPermission
		}
		return 0, err
	}

	return permission.ID, nil
}

func (rr *RbacRepo) GetUserRoles(userID uint) ([]model.Role, error) {
	var user = new(model.User)
	err := rr.db.Preload("Roles").First(user, "id = ?", userID).Error
	return user.Roles, err
}

func (rr *RbacRepo) HasPermission(questionnaireID uint, permissionID uint) error {
	var rolePermission model.RolePermission
	err := rr.db.First(&rolePermission, "questionnaire_id = ? and permission_id = ?", questionnaireID, permissionID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.ErrorNotHavePermission
	}
	return err
}

func (rr *RbacRepo) MakeNewRole(tx *gorm.DB, name string) (uint, error) {
	role := model.Role{
		Name: name,
	}

	err := tx.Create(&role).Error
	if err != nil {
		return 0, err
	}

	return role.ID, nil
}

func (rr *RbacRepo) MakeRolePermission(tx *gorm.DB, roleID uint, questionnaireID uint, permissionID uint, expireDate sql.NullTime) (uint, error) {
	var rolePermission = model.RolePermission{
		RoleID:          roleID,
		QuestionnaireID: questionnaireID,
		PermissionID:    permissionID,
		ExpireAt:        expireDate,
	}

	result := tx.Create(&rolePermission)
	return rolePermission.ID, result.Error
}

func (rr *RbacRepo) InsertUsersWithVisibleAnswers(tx *gorm.DB, rolePermissionID uint, selectedUsersIDs []uint) error {
	for _, userID := range selectedUsersIDs {
		if err := tx.Create(&model.UsersWithVisibleAnswers{
			RolePermissionID: rolePermissionID,
			UserID:           userID,
		}).Error; err != nil {
			return err
		}
	}

	return nil
}

func (rr *RbacRepo) MakeRoleUser(tx *gorm.DB, roleID uint, userID uint) error {
	return tx.Create(&model.RoleUser{
		RoleID: roleID,
		UserID: userID,
	}).Error
}

func (rr *RbacRepo) Transaction(f func(tx *gorm.DB) error) error {
	return rr.db.Transaction(func(tx *gorm.DB) error {
		return f(tx)
	})
}

func (rr *RbacRepo) MakeUser(user model.User) (model.User, error) {
	err := rr.db.Where("id = ?", user.ID).FirstOrCreate(&user).Error
	return user, err
}

func (rr *RbacRepo) MakeQuestionnaire(questionnaire model.Questionnaire) (model.Questionnaire, error) {
	err := rr.db.Where("id = ?", questionnaire.ID).FirstOrCreate(&questionnaire).Error
	return questionnaire, err
}

func (rr *RbacRepo) GetUserWithQuestionnaires(userID uint) (model.User, error) {
	var user model.User
	err := rr.db.Preload("Questionnaires").First(&user, "id = ?", userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, model.ErrorNotFoundUser
	}
	return user, err
}

func (rr *RbacRepo) GetUserRolesWithPermissions(userID uint) ([]domain_repository.RoleWithPermissions, error) {
	var rolesWithPermissions []domain_repository.RoleWithPermissions

	// Query user roles with permissions and additional details
	rows, err := rr.db.Table("roles").
		Select(`
            roles.id AS role_id,
            roles.name AS role_name,
            permissions.id AS permission_id,
            permissions.name AS permission_name,
            role_permissions.questionnaire_id,
            role_permissions.expire_at
        `).
		Joins("JOIN role_users ON role_users.role_id = roles.id").
		Joins("JOIN role_permissions ON role_permissions.role_id = roles.id").
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
		Where("role_users.user_id = ?", userID).
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roleMap := make(map[uint]*domain_repository.RoleWithPermissions)

	for rows.Next() {
		var roleID uint
		var permission domain_repository.PermissionWithDetails
		var roleName string

		err := rows.Scan(&roleID, &roleName, &permission.PermissionID, &permission.Name, &permission.QuestionnaireID, &permission.ExpireAt)
		if err != nil {
			return nil, err
		}

		// Check if the role already exists in the map
		if _, exists := roleMap[roleID]; !exists {
			roleMap[roleID] = &domain_repository.RoleWithPermissions{
				ID:          roleID,
				Name:        roleName,
				Permissions: []domain_repository.PermissionWithDetails{},
			}
		}

		// Append the permission details to the role
		roleMap[roleID].Permissions = append(roleMap[roleID].Permissions, permission)
	}

	// Convert map to slice
	for _, role := range roleMap {
		rolesWithPermissions = append(rolesWithPermissions, *role)
	}

	return rolesWithPermissions, nil
}

func (rr *RbacRepo) DeleteRolePermissions(roleID uint, questionnaireID uint, permissionID uint) error {
	return rr.db.Delete(&model.RolePermission{}, "role_id = ? and questionnaire_id = ? and permission_id = ?", roleID, questionnaireID, permissionID).Error
}

func (rr *RbacRepo) HasRolePermission(roleID uint, questionnaireID uint, permissionID uint) (bool, error) {
	rolePermission := new(model.RolePermission)
	err := rr.db.First(rolePermission, "role_id = ? and questionnaire_id = ? and permission_id = ? and (expire_at IS NULL OR expire_at > NOW())", roleID, questionnaireID, permissionID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (rr *RbacRepo) IsSuperadmin(userID uint) (model.Superadmin, error) {
	superadmin := new(model.Superadmin)
	err := rr.db.First(superadmin, "user_id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Superadmin{}, nil
		}

		return model.Superadmin{}, err
	}

	return *superadmin, nil
}

func (rr *RbacRepo) FindSuperadminPermission(superadminID uint, permissionID uint) (bool, error) {
	err := rr.db.First(&model.SuperadminPermission{}, "superadmin_id = ? and permission_id = ?", superadminID, permissionID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (rr *RbacRepo) GiveSuperadminRole(tx *gorm.DB, userID uint, giverUserID uint) (uint, error) {
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

func (rr *RbacRepo) MakeSuperadminPermission(tx *gorm.DB, superadminID uint, permissionID uint) error {
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

func (rr *RbacRepo) FindRolePermission(roleID, questionnaireID, permissionID uint) (*model.RolePermission, error) {
	rolePermission := new(model.RolePermission)
	err := rr.db.First(&rolePermission, "role_id = ? and questionnaire_id = ? and permission_id = ? and (expire_at IS NULL OR expire_at > NOW())", roleID, questionnaireID, permissionID).Error
	return rolePermission, err
}

func (rr *RbacRepo) FindUsersWithVisibleAnswers(rolePermissionID uint) ([]uint, error) {
	var userIDs []uint
	err := rr.db.Model(&model.UsersWithVisibleAnswers{}).Where("role_permission_id = ?", rolePermissionID).Pluck("user_id", &userIDs).Error
	return userIDs, err
}