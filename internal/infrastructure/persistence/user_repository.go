// Package persistence 持久化层实现
//
//	@update 2025-09-30 00:00:00
package persistence

import (
	"context"

	"github.com/hcd233/go-backend-tmpl/internal/domain/user"
	"github.com/hcd233/go-backend-tmpl/internal/resource/database"
	"github.com/hcd233/go-backend-tmpl/internal/resource/database/dao"
	"github.com/hcd233/go-backend-tmpl/internal/resource/database/model"
)

// UserRepository 用户仓储实现
type UserRepository struct {
	userDAO *dao.UserDAO
}

// NewUserRepository 创建用户仓储
func NewUserRepository() user.Repository {
	return &UserRepository{
		userDAO: dao.GetUserDAO(),
	}
}

// Save 保存用户
func (r *UserRepository) Save(ctx context.Context, u *user.User) error {
	db := database.GetDBInstance(ctx)

	// 转换为持久化模型
	po := r.toPO(u)

	// 如果ID为0，创建新用户；否则更新
	if u.ID() == 0 {
		if err := r.userDAO.Create(db, po); err != nil {
			return err
		}
		// 回填ID
		u.SetID(po.ID)
		return nil
	}

	// 更新用户
	updateData := map[string]interface{}{
		"name":           po.Name,
		"email":          po.Email,
		"avatar":         po.Avatar,
		"permission":     po.Permission,
		"last_login":     po.LastLogin,
		"github_bind_id": po.GithubBindID,
		"qq_bind_id":     po.QQBindID,
		"google_bind_id": po.GoogleBindID,
	}
	return r.userDAO.Update(db, &model.User{BaseModel: model.BaseModel{ID: u.ID()}}, updateData)
}

// FindByID 根据ID查找用户
func (r *UserRepository) FindByID(ctx context.Context, id uint) (*user.User, error) {
	db := database.GetDBInstance(ctx)

	po, err := r.userDAO.GetByID(db, id, []string{"*"}, []string{})
	if err != nil {
		return nil, err
	}

	return r.toDO(po)
}

// FindByEmail 根据邮箱查找用户
func (r *UserRepository) FindByEmail(ctx context.Context, email user.Email) (*user.User, error) {
	db := database.GetDBInstance(ctx)

	po, err := r.userDAO.GetByEmail(db, email.Value(), []string{"*"}, []string{})
	if err != nil {
		return nil, err
	}

	return r.toDO(po)
}

// FindByName 根据用户名查找用户
func (r *UserRepository) FindByName(ctx context.Context, name string) (*user.User, error) {
	db := database.GetDBInstance(ctx)

	po, err := r.userDAO.GetByName(db, name, []string{"*"}, []string{})
	if err != nil {
		return nil, err
	}

	return r.toDO(po)
}

// FindByGithubBindID 根据Github绑定ID查找用户
func (r *UserRepository) FindByGithubBindID(ctx context.Context, githubID string) (*user.User, error) {
	db := database.GetDBInstance(ctx)

	var po *model.User
	err := db.Where("github_bind_id = ?", githubID).First(&po).Error
	if err != nil {
		return nil, err
	}

	return r.toDO(po)
}

// FindByGoogleBindID 根据Google绑定ID查找用户
func (r *UserRepository) FindByGoogleBindID(ctx context.Context, googleID string) (*user.User, error) {
	db := database.GetDBInstance(ctx)

	var po *model.User
	err := db.Where("google_bind_id = ?", googleID).First(&po).Error
	if err != nil {
		return nil, err
	}

	return r.toDO(po)
}

// Exists 判断用户是否存在
func (r *UserRepository) Exists(ctx context.Context, id uint) (bool, error) {
	db := database.GetDBInstance(ctx)

	var count int64
	err := db.Model(&model.User{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Delete 删除用户
func (r *UserRepository) Delete(ctx context.Context, id uint) error {
	db := database.GetDBInstance(ctx)
	return r.userDAO.Delete(db, &model.User{BaseModel: model.BaseModel{ID: id}})
}

// toPO 将领域模型转换为持久化模型
func (r *UserRepository) toPO(u *user.User) *model.User {
	return &model.User{
		BaseModel: model.BaseModel{
			ID:        u.ID(),
			CreatedAt: u.CreatedAt(),
			UpdatedAt: u.UpdatedAt(),
		},
		Name:         u.Name(),
		Email:        u.Email().Value(),
		Avatar:       u.Avatar(),
		Permission:   model.Permission(u.Permission()),
		LastLogin:    u.LastLogin(),
		GithubBindID: u.GithubBindID(),
		QQBindID:     u.QQBindID(),
		GoogleBindID: u.GoogleBindID(),
	}
}

// toDO 将持久化模型转换为领域模型
func (r *UserRepository) toDO(po *model.User) (*user.User, error) {
	email, err := user.NewEmail(po.Email)
	if err != nil {
		return nil, err
	}

	return user.ReconstructUser(
		po.ID,
		po.Name,
		*email,
		po.Avatar,
		user.Permission(po.Permission),
		po.LastLogin,
		po.CreatedAt,
		po.UpdatedAt,
		po.GithubBindID,
		po.QQBindID,
		po.GoogleBindID,
	), nil
}

var _ user.Repository = (*UserRepository)(nil) // 确保实现了接口
