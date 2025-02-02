package minio

import (
	"bytes"
	"context"
	"encoding/json"
	"os"

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

func (d *DataRepository) SaveFile(ctx context.Context, file *fl.File, login string, id string) error {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		logger.Error("error while creating minio client: ", zap.Error(err))

		return err
	}

	// Define the file to upload and the destination bucket
	objectName := id // The name for the object in MinIO
	bucketName := login	// Bucket name in MinIO

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


func (d *DataRepository) SaveTextData(ctx context.Context, data any, login string, id string) error {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		logger.Error("error while creating minio client: ", zap.Error(err))

		return err
	}

	// Serialize the map to JSON
	serializedData, err := json.Marshal(data)
	if err != nil {
		logger.Error("error while serializing the map: ", zap.Error(err))

		return err
	}

	// Create a buffer from the serialized data
	buf := bytes.NewReader(serializedData)

	// Define the file to upload and the destination bucket
	objectName := id // The name for the object in MinIO
	bucketName := login	// Bucket name in MinIO

	// Ensure the bucket exists (or create it)
	err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		if exists, errBucketExists := client.BucketExists(context.Background(), bucketName); errBucketExists == nil && exists {
			logger.Info("Bucket already exists")
		} else {
			logger.Info("Failed to create bucket:", zap.Error(err))
		}
	}

	_, err = client.PutObject(ctx, bucketName, objectName, buf, int64(buf.Len()), minio.PutObjectOptions{ContentType: "application/json"})
	if err != nil {
		logger.Error("error while uploading object to minio: ", zap.Error(err))

		return err
	}

	logger.Info("string data uploaded to Minio successfully")

	return nil
}