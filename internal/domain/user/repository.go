// Package user 用户仓储接口
//
//	@update 2025-09-30 00:00:00
package user

import "context"

// Repository 用户仓储接口
//
//	@update 2025-09-30 00:00:00
type Repository interface {
	// Save 保存用户（新增或更新）
	Save(ctx context.Context, user *User) error

	// FindByID 根据ID查找用户
	FindByID(ctx context.Context, id uint) (*User, error)

	// FindByEmail 根据邮箱查找用户
	FindByEmail(ctx context.Context, email Email) (*User, error)

	// FindByName 根据用户名查找用户
	FindByName(ctx context.Context, name string) (*User, error)

	// FindByGithubBindID 根据Github绑定ID查找用户
	FindByGithubBindID(ctx context.Context, githubID string) (*User, error)

	// FindByGoogleBindID 根据Google绑定ID查找用户
	FindByGoogleBindID(ctx context.Context, googleID string) (*User, error)

	// Exists 判断用户是否存在
	Exists(ctx context.Context, id uint) (bool, error)

	// Delete 删除用户
	Delete(ctx context.Context, id uint) error
}

