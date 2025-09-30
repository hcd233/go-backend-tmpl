// Package user 用户领域服务
//
//	@update 2025-09-30 00:00:00
package user

import (
	"context"
	"errors"
	"fmt"
	"regexp"
)

var (
	// ErrUserNotFound 用户不存在
	ErrUserNotFound = errors.New("user not found")
	// ErrUserAlreadyExists 用户已存在
	ErrUserAlreadyExists = errors.New("user already exists")
	// ErrInvalidUserName 无效的用户名
	ErrInvalidUserName = errors.New("invalid user name")
)

// Service 用户领域服务
//
//	@update 2025-09-30 00:00:00
type Service struct {
	repo Repository
}

// NewService 创建用户领域服务
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// ValidateUserName 验证用户名
func ValidateUserName(name string) error {
	if len(name) < 3 || len(name) > 20 {
		return fmt.Errorf("%w: length must be between 3 and 20", ErrInvalidUserName)
	}

	// 用户名只能包含字母、数字、下划线和连字符
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, name)
	if !matched {
		return fmt.Errorf("%w: can only contain letters, numbers, underscores and hyphens", ErrInvalidUserName)
	}

	return nil
}

// CreateUser 创建用户
func (s *Service) CreateUser(ctx context.Context, name string, emailStr string, avatar string, permission Permission) (*User, error) {
	// 验证用户名
	if err := ValidateUserName(name); err != nil {
		return nil, err
	}

	// 创建邮箱值对象
	email, err := NewEmail(emailStr)
	if err != nil {
		return nil, err
	}

	// 检查邮箱是否已存在
	existingUser, _ := s.repo.FindByEmail(ctx, *email)
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// 创建用户实体
	user, err := NewUser(name, *email, avatar, permission)
	if err != nil {
		return nil, err
	}

	// 保存用户
	if err := s.repo.Save(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetOrCreateUser 获取或创建用户（用于OAuth2登录）
func (s *Service) GetOrCreateUser(ctx context.Context, emailStr string, name string, avatar string) (*User, error) {
	// 创建邮箱值对象
	email, err := NewEmail(emailStr)
	if err != nil {
		return nil, err
	}

	// 尝试查找现有用户
	user, err := s.repo.FindByEmail(ctx, *email)
	if err == nil && user != nil {
		return user, nil
	}

	// 用户不存在，创建新用户
	return s.CreateUser(ctx, name, emailStr, avatar, PermissionReader)
}

