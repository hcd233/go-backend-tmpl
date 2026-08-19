// Package query 定义 Identity 域查询处理器。
package query

import (
	"context"

	identityport "github.com/hcd233/aris-api-tmpl/internal/application/identity/port"
	"github.com/hcd233/aris-api-tmpl/internal/common/ierr"
	"github.com/hcd233/aris-api-tmpl/internal/domain/identity"
	"github.com/hcd233/aris-api-tmpl/internal/logger"
	"go.uber.org/zap"
)

type getCurrentUserHandler struct {
	repo identity.UserRepository
}

// NewGetCurrentUserHandler 构造查询处理器。
func NewGetCurrentUserHandler(repo identity.UserRepository) identityport.GetCurrentUserHandler {
	return &getCurrentUserHandler{repo: repo}
}

// Handle 执行当前用户查询。
func (h *getCurrentUserHandler) Handle(ctx context.Context, q identityport.GetCurrentUserQuery) (*identityport.UserView, error) {
	log := logger.WithCtx(ctx)
	user, err := h.repo.FindByID(ctx, q.UserID)
	if err != nil {
		log.Error("[IdentityQuery] find user failed", zap.Error(err), zap.Uint("userID", q.UserID))
		return nil, err
	}
	if user == nil {
		return nil, ierr.New(ierr.ErrDataNotExists, "user not found")
	}
	return &identityport.UserView{
		ID:         user.AggregateID(),
		Name:       user.Name().String(),
		Email:      user.Email().String(),
		Avatar:     user.Avatar().String(),
		Permission: user.Permission(),
		CreatedAt:  user.CreatedAt(),
		LastLogin:  user.LastLogin(),
	}, nil
}
