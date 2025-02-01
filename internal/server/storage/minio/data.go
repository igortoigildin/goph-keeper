package minio

import (
	"context"
	"os"
	"path/filepath"

	fl "github.com/igortoigildin/goph-keeper/pkg/file"
	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

const (
	accessKey = "minioaccesskey"
	secretKey = "miniosecretkey"
	useSSL    = false // Whether to use SSL
	endpoint  = "localhost:9000"
)

type DataRepository struct{}

func NewRepository() *DataRepository {
	return &DataRepository{}
}

func (d *DataRepository) SaveData(ctx context.Context, file *fl.File, bucketName string) error {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		logger.Error("error while creating minio client: ", zap.Error(err))

		return err
	}

	// Define the file to upload and the destination bucket
	objectName := filepath.Base(file.FilePath) // The name for the object in MinIO

	// Ensure the bucket exists (or create it)
	err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		if exists, errBucketExists := client.BucketExists(context.Background(), bucketName); errBucketExists == nil && exists {
			logger.Info("Bucket already exists")
		} else {
			logger.Info("Failed to create bucket:", zap.Error(err))
		}
	}

	// Open the file to upload
	f, err := os.Open(file.FilePath)
	if err != nil {
		logger.Error("error while uploading file to minio: ", zap.Error(err))
		return err
	}
	defer f.Close()

	// Upload the file to MinIO
	_, err = client.PutObject(
		context.Background(),
		bucketName,
		objectName,
		f,
		-1, // -1 means the file size will be determined automatically
		minio.PutObjectOptions{ContentType: "application/octet-stream"},
	)
	if err != nil {
		logger.Error("error while uploading file to MinIO", zap.Error(err))
		return err
	}

	logger.Info("File uploaded successfully")

	return nil
}
