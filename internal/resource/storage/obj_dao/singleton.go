package objdao

import (
	"github.com/hcd233/go-backend-tmpl/internal/config"
	"github.com/hcd233/go-backend-tmpl/internal/resource/storage"
)

var (
	// ImageObjDAOSingleton 图片对象DAO单例
	//	update 2025-01-05 22:45:54
	ImageObjDAOSingleton ObjDAO

	// ThumbnailObjDAOSingleton 缩略图对象DAO单例
	//	update 2025-01-05 22:45:54
	ThumbnailObjDAOSingleton ObjDAO
)

func init() {
	ImageObjDAOSingleton = createObjectStorageDAO(ObjectTypeImage)
	ThumbnailObjDAOSingleton = createObjectStorageDAO(ObjectTypeThumbnail)
}

// createObjectStorageDAO 创建对象存储DAO
func createObjectStorageDAO(objectType ObjectType) ObjDAO {
	switch storage.GetProvider() {
	case storage.ProviderMinio:
		return &MinioObjDAO{
			ObjectType: objectType,
			BucketName: config.MinioBucketName,
			client:     storage.GetMinioStorage(),
		}
	case storage.ProviderCOS:
		return &CosObjDAO{
			ObjectType: objectType,
			BucketName: config.CosBucketName,
			client:     storage.GetCosClient(),
		}
	default:
		panic("unsupported storage type")
	}
}

// GetImageObjDAO 获取图片对象DAO单例
//
//	return ObjDAO
//	author centonhuang
//	update 2024-10-18 01:10:28
func GetImageObjDAO() ObjDAO {
	return ImageObjDAOSingleton
}

// GetThumbnailObjDAO 获取缩略图对象DAO单例
//
//	return ObjDAO
//	author centonhuang
//	update 2024-10-18 01:09:59
func GetThumbnailObjDAO() ObjDAO {
	return ThumbnailObjDAOSingleton
}
