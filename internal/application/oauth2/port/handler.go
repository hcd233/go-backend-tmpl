// Package port 定义 OAuth2 域的命令端口（依赖倒置接口）。
// handler 层只依赖本包，不依赖 command 的具体实现。
package port

import (
	"context"

	identityvo "github.com/hcd233/aris-api-tmpl/internal/domain/identity/vo"
)

// InitiateLoginCommand 发起 OAuth 登录命令。
type InitiateLoginCommand struct {
	Platform string
}

// InitiateLoginResult 登录发起结果。
type InitiateLoginResult struct {
	RedirectURL string
}

// InitiateLoginHandler 登录发起命令处理器。
type InitiateLoginHandler interface {
	Handle(ctx context.Context, cmd InitiateLoginCommand) (*InitiateLoginResult, error)
}

// HandleCallbackCommand OAuth2 回调处理命令。
type HandleCallbackCommand struct {
	Platform string
	Code     string
	State    string
}

// HandleCallbackResult 回调处理结果。
type HandleCallbackResult struct {
	TokenPair *identityvo.TokenPair
	UserID    uint
	IsNewUser bool
}

// ObjectStorageDirCreator 对象存储目录创建器。
type ObjectStorageDirCreator interface {
	CreateDir(ctx context.Context, userID uint) error
}

// HandleCallbackHandler 回调命令处理器。
type HandleCallbackHandler interface {
	Handle(ctx context.Context, cmd HandleCallbackCommand) (*HandleCallbackResult, error)
}
