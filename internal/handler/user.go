package handler

import (
	"context"
	"time"

	apiutil "github.com/hcd233/aris-api-tmpl/internal/api/util"
	identityport "github.com/hcd233/aris-api-tmpl/internal/application/identity/port"
	"github.com/hcd233/aris-api-tmpl/internal/common/constant"
	"github.com/hcd233/aris-api-tmpl/internal/common/ierr"
	"github.com/hcd233/aris-api-tmpl/internal/dto"
	"github.com/hcd233/aris-api-tmpl/internal/logger"
	"github.com/hcd233/aris-api-tmpl/internal/util"
	"go.uber.org/zap"
)

// UserHandler 用户处理器。
type UserHandler interface {
	HandleGetCurUser(ctx context.Context, req *dto.EmptyReq) (*dto.HTTPResponse[*dto.GetCurUserRsp], error)
	HandleUpdateUser(ctx context.Context, req *dto.UpdateUserReq) (*dto.HTTPResponse[*dto.EmptyRsp], error)
}

// UserDependencies UserHandler 依赖项。
type UserDependencies struct {
	GetCurrentUser identityport.GetCurrentUserHandler
	UpdateProfile  identityport.UpdateProfileHandler
}

type userHandler struct {
	getCurrentUser identityport.GetCurrentUserHandler
	updateProfile  identityport.UpdateProfileHandler
}

// NewUserHandler 创建用户处理器。
func NewUserHandler(deps UserDependencies) UserHandler {
	return &userHandler{getCurrentUser: deps.GetCurrentUser, updateProfile: deps.UpdateProfile}
}

// HandleGetCurUser 获取当前用户信息。
func (h *userHandler) HandleGetCurUser(ctx context.Context, _ *dto.EmptyReq) (*dto.HTTPResponse[*dto.GetCurUserRsp], error) {
	userID := util.CtxValueUint(ctx, constant.CtxKeyUserID)
	if userID == 0 {
		return nil, apiutil.NewHumaBizErrorFromModel(ctx, ierr.ErrUnauthorized.BizError())
	}
	view, err := h.getCurrentUser.Handle(ctx, identityport.GetCurrentUserQuery{UserID: userID})
	if err != nil {
		logger.WithCtx(ctx).Error("[UserHandler] get current user failed", zap.Error(err))
		return nil, apiutil.NewHumaBizError(ctx, err, ierr.ErrInternal.BizError())
	}
	rsp := &dto.GetCurUserRsp{
		User: &dto.DetailedUser{
			ID:         view.ID,
			CreatedAt:  view.CreatedAt.Format(time.DateTime),
			LastLogin:  view.LastLogin.Format(time.DateTime),
			Permission: string(view.Permission),
			User: dto.User{
				Name:   view.Name,
				Email:  view.Email,
				Avatar: view.Avatar,
			},
		},
	}
	return util.WrapHTTPResponse(rsp, nil)
}

// HandleUpdateUser 更新当前用户资料。
func (h *userHandler) HandleUpdateUser(ctx context.Context, req *dto.UpdateUserReq) (*dto.HTTPResponse[*dto.EmptyRsp], error) {
	userID := util.CtxValueUint(ctx, constant.CtxKeyUserID)
	if userID == 0 || req == nil || req.Body == nil || req.Body.User == nil {
		return nil, apiutil.NewHumaBizErrorFromModel(ctx, ierr.ErrBadRequest.BizError())
	}
	if err := h.updateProfile.Handle(ctx, identityport.UpdateProfileCommand{
		UserID: userID,
		Name:   req.Body.User.Name,
		Email:  req.Body.User.Email,
		Avatar: req.Body.User.Avatar,
	}); err != nil {
		logger.WithCtx(ctx).Error("[UserHandler] update user failed", zap.Error(err))
		return nil, apiutil.NewHumaBizError(ctx, err, ierr.ErrInternal.BizError())
	}
	return util.WrapHTTPResponse(&dto.EmptyRsp{}, nil)
}
