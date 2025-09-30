// Package user 用户应用服务
//
//	@update 2025-09-30 00:00:00
package user

import (
	"context"
	"errors"

	"github.com/hcd233/go-backend-tmpl/internal/domain/user"
	"github.com/hcd233/go-backend-tmpl/internal/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	// ErrUserNotFound 用户不存在
	ErrUserNotFound = errors.New("user not found")
	// ErrInternalError 内部错误
	ErrInternalError = errors.New("internal error")
)

// Service 用户应用服务
//
//	@update 2025-09-30 00:00:00
type Service struct {
	userRepo user.Repository
}

// NewService 创建用户应用服务
func NewService(userRepo user.Repository) *Service {
	return &Service{
		userRepo: userRepo,
	}
}

// GetUserInfo 获取用户信息
func (s *Service) GetUserInfo(ctx context.Context, query *GetUserInfoQuery) (*UserInfoResponse, error) {
	logger := logger.WithCtx(ctx)

	// 从仓储获取用户
	u, err := s.userRepo.FindByID(ctx, query.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[UserAppService] user not found", zap.Uint("userID", query.UserID))
			return nil, ErrUserNotFound
		}
		logger.Error("[UserAppService] failed to get user", zap.Error(err))
		return nil, ErrInternalError
	}

	// 转换为DTO
	userDTO := &UserDTO{
		UserID:    u.ID(),
		Name:      u.Name(),
		Email:     u.Email().Value(),
		Avatar:    u.Avatar(),
		CreatedAt: FormatTime(u.CreatedAt()),
		LastLogin: FormatTime(u.LastLogin()),
	}

	logger.Info("[UserAppService] get user info success",
		zap.Uint("userID", u.ID()),
		zap.String("name", u.Name()))

	return &UserInfoResponse{User: userDTO}, nil
}

// GetCurUserInfo 获取当前用户信息（包含权限）
func (s *Service) GetCurUserInfo(ctx context.Context, query *GetCurUserInfoQuery) (*CurUserInfoResponse, error) {
	logger := logger.WithCtx(ctx)

	// 从仓储获取用户
	u, err := s.userRepo.FindByID(ctx, query.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[UserAppService] user not found", zap.Uint("userID", query.UserID))
			return nil, ErrUserNotFound
		}
		logger.Error("[UserAppService] failed to get user", zap.Error(err))
		return nil, ErrInternalError
	}

	// 转换为DTO
	curUserDTO := &CurUserDTO{
		UserID:     u.ID(),
		Name:       u.Name(),
		Email:      u.Email().Value(),
		Avatar:     u.Avatar(),
		CreatedAt:  FormatTime(u.CreatedAt()),
		LastLogin:  FormatTime(u.LastLogin()),
		Permission: string(u.Permission()),
	}

	logger.Info("[UserAppService] get current user info success",
		zap.Uint("userID", u.ID()),
		zap.String("permission", string(u.Permission())))

	return &CurUserInfoResponse{User: curUserDTO}, nil
}

// UpdateUserInfo 更新用户信息
func (s *Service) UpdateUserInfo(ctx context.Context, cmd *UpdateUserInfoCommand) (*UpdateUserInfoResponse, error) {
	logger := logger.WithCtx(ctx)

	// 从仓储获取用户
	u, err := s.userRepo.FindByID(ctx, cmd.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[UserAppService] user not found", zap.Uint("userID", cmd.UserID))
			return nil, ErrUserNotFound
		}
		logger.Error("[UserAppService] failed to get user", zap.Error(err))
		return nil, ErrInternalError
	}

	// 调用领域模型方法更新
	if err := u.UpdateName(cmd.UpdatedName); err != nil {
		logger.Error("[UserAppService] failed to update user name", zap.Error(err))
		return nil, err
	}

	// 保存到仓储
	if err := s.userRepo.Save(ctx, u); err != nil {
		logger.Error("[UserAppService] failed to save user", zap.Error(err))
		return nil, ErrInternalError
	}

	logger.Info("[UserAppService] update user info success",
		zap.Uint("userID", cmd.UserID),
		zap.String("newName", cmd.UpdatedName))

	return &UpdateUserInfoResponse{}, nil
}
