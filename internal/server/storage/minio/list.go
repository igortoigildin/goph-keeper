package minio

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

// Struct for storing object information
type ObjectInfo struct {
    Key          string    `json:"key"`
    Size         int64     `json:"size"`
    LastModified time.Time `json:"last_modified"`
    ETag         string    `json:"etag"`
}


func (d *DataRepository) ListObjects(ctx context.Context, bucketName string) ([]ObjectInfo, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		logger.Error("error while creating minio client: ", zap.Error(err))

		return nil, fmt.Errorf("error instantiating Minio client with options: %w", err)
	}

	allObjects := []ObjectInfo{}

	objectCh := client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
        Recursive: true,
    })

	for object := range objectCh {
        if object.Err != nil {
			logger.Error("error while listing objects: ", zap.Error(object.Err))
            continue
        }

        info := ObjectInfo{
            Key:          object.Key,
            Size:         object.Size,
            LastModified: object.LastModified,
            ETag:         object.ETag,
        }

        allObjects = append(allObjects, info)
    }

    file, err := os.Create("minio_objects.json")
    if err != nil {
		logger.Error("error while creating file: ", zap.Error(err))
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    if err := encoder.Encode(allObjects); err != nil {
		logger.Error("error while encoding JSON: ", zap.Error(err))
    }

	logger.Info("JSON saved to minio_objects.json")

	return allObjects, nil
}
