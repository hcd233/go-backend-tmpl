// Package command 定义 Identity 域命令处理器。
package command

import (
	"context"

	identityport "github.com/hcd233/aris-api-tmpl/internal/application/identity/port"
	"github.com/hcd233/aris-api-tmpl/internal/common/ierr"
	"github.com/hcd233/aris-api-tmpl/internal/domain/identity"
	"github.com/hcd233/aris-api-tmpl/internal/domain/identity/vo"
	"github.com/hcd233/aris-api-tmpl/internal/logger"
	"go.uber.org/zap"
)

type updateProfileHandler struct {
	repo identity.UserRepository
}

// NewUpdateProfileHandler 构造更新资料处理器。
func NewUpdateProfileHandler(repo identity.UserRepository) identityport.UpdateProfileHandler {
	return &updateProfileHandler{repo: repo}
}

// Handle 执行资料更新。
func (h *updateProfileHandler) Handle(ctx context.Context, cmd identityport.UpdateProfileCommand) error {
	log := logger.WithCtx(ctx)
	user, err := h.repo.FindByID(ctx, cmd.UserID)
	if err != nil {
		log.Error("[IdentityCommand] find user failed", zap.Error(err), zap.Uint("userID", cmd.UserID))
		return err
	}
	if user == nil {
		log.Warn("[IdentityCommand] user not found for profile update", zap.Uint("userID", cmd.UserID))
		return ierr.New(ierr.ErrDataNotExists, "user not found")
	}
	if err := user.UpdateProfile(vo.UserName(cmd.Name), vo.Email(cmd.Email), vo.Avatar(cmd.Avatar)); err != nil {
		log.Warn("[IdentityCommand] update profile validation failed", zap.Error(err), zap.Uint("userID", cmd.UserID))
		return err
	}
	if err := h.repo.Save(ctx, user); err != nil {
		log.Error("[IdentityCommand] save user failed", zap.Error(err), zap.Uint("userID", cmd.UserID))
		return err
	}
	return nil
}
