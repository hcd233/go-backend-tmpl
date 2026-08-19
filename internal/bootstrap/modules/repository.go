package modules

import (
	oauth2port "github.com/hcd233/aris-api-tmpl/internal/application/oauth2/port"
	"github.com/hcd233/aris-api-tmpl/internal/config"
	"github.com/hcd233/aris-api-tmpl/internal/domain/identity"
	"github.com/hcd233/aris-api-tmpl/internal/infrastructure/repository"
	"github.com/hcd233/aris-api-tmpl/internal/infrastructure/storage"
	"go.uber.org/fx"
)

// RepositoryModule 仓储模块：领域仓储接口的实现。
var RepositoryModule = fx.Module("repository",
	fx.Provide(
		NewUserRepository,
		NewAudioDirCreator,
	),
)

// NewUserRepository 构造用户仓储。
func NewUserRepository() identity.UserRepository {
	return repository.NewUserRepository()
}

// NewAudioDirCreator 构造对象存储目录创建器；未配置对象存储时返回 nil。
func NewAudioDirCreator() oauth2port.ObjectStorageDirCreator {
	if config.CosAppID == "" && config.MinioEndpoint == "" {
		return nil
	}
	storage.InitObjectStorage()
	return repository.NewAudioDirCreator()
}
