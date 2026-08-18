// Package port 定义 Identity 域的命令/查询端口（依赖倒置接口）。
// handler 层只依赖本包，不依赖 command/query 的具体实现。
package port

import (
	"context"
	"time"

	"github.com/hcd233/aris-api-tmpl/internal/common/enum"
	"github.com/hcd233/aris-api-tmpl/internal/domain/identity/vo"
)

// RefreshTokensCommand 刷新 token 对命令。
type RefreshTokensCommand struct {
	RefreshToken string
}

// RefreshTokensHandler 刷新命令处理器。
type RefreshTokensHandler interface {
	Handle(ctx context.Context, cmd RefreshTokensCommand) (*vo.TokenPair, error)
}

// UpdateProfileCommand 更新用户资料命令。
type UpdateProfileCommand struct {
	UserID uint
	Name   string
	Email  string
	Avatar string
}

// UpdateProfileHandler 更新资料命令处理器。
type UpdateProfileHandler interface {
	Handle(ctx context.Context, cmd UpdateProfileCommand) error
}

// UserView 用户详情只读投影。
type UserView struct {
	ID         uint
	Name       string
	Email      string
	Avatar     string
	Permission enum.Permission
	CreatedAt  time.Time
	LastLogin  time.Time
}

// GetCurrentUserQuery 查询当前用户命令。
type GetCurrentUserQuery struct {
	UserID uint
}

// GetCurrentUserHandler 查询处理器。
type GetCurrentUserHandler interface {
	Handle(ctx context.Context, q GetCurrentUserQuery) (*UserView, error)
}
