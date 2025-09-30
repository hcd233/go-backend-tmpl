package storage

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hcd233/go-backend-tmpl/internal/config"
	"github.com/hcd233/go-backend-tmpl/internal/logger"
	"github.com/samber/lo"
	"github.com/tencentyun/cos-go-sdk-v5"
	"go.uber.org/zap"
)

var cosClient *cos.Client

func initCosClient() {
	endpoint := lo.Must1(url.Parse(fmt.Sprintf("https://%s-%s.cos.%s.myqcloud.com", config.CosBucketName, config.CosAppID, config.CosRegion)))

	uri := &cos.BaseURL{
		BucketURL: endpoint,
	}

	credential := &cos.Credential{
		SecretID:  config.CosSecretID,
		SecretKey: config.CosSecretKey,
	}
	cosClient = cos.NewClient(uri, &http.Client{
		Transport: &cos.CredentialTransport{
			Credential: credential,
		},
	})

	_, _ = lo.Must2(cosClient.Bucket.Get(context.Background(), &cos.BucketGetOptions{}))

	logger.Logger().Info("[Object Storage] Connected to COS", zap.String("endpoint", endpoint.String()))
}

// GetCosClient 获取 COS 客户端
func GetCosClient() *cos.Client {
	return cosClient
}
