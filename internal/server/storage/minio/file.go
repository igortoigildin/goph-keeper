package minio

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
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
	useSSL    = false
	endpoint  = "localhost:9000"
	binData   = "bin_data"
)

type DataRepository struct{}

func NewRepository() *DataRepository {
	return &DataRepository{}
}

func (d *DataRepository) SaveFile(ctx context.Context, file *fl.File, login string, id string, meta string) (string, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		logger.Error("error while creating minio client: ", zap.Error(err))

		return "", fmt.Errorf("error instantiating Minio client with options: %w", err)
	}

	// Define the file to upload and the destination bucket
	objectName := binData + "_" + id // The name for the object in MinIO
	bucketName := login              // Bucket name in MinIO

	// Ensure the bucket exists (or create it)
	err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		if exists, errBucketExists := client.BucketExists(ctx, bucketName); errBucketExists == nil && exists {
			logger.Info("Bucket already exists")
		} else {
			logger.Info("Failed to create bucket:", zap.Error(err))

			return "", fmt.Errorf("Minio error: %w", err)
		}
	}

	meatadata := map[string]string{
		"meta":     meta,
		"dataType": binData,
	}

	// Open the file to upload
	f, err := os.Open(file.FilePath)
	if err != nil {
		logger.Error("error opening targeted file: ", zap.Error(err))

		return "", fmt.Errorf("error opening targeted file: %w", err)
	}
	defer f.Close()

	// Upload the file to MinIO
	objectInfo, err := client.PutObject(
		context.Background(),
		bucketName,
		objectName,
		f,
		-1,
		minio.PutObjectOptions{ContentType: "application/octet-stream", UserMetadata: meatadata},
	)
	if err != nil {
		logger.Error("error while uploading file to MinIO", zap.Error(err))

		return "", fmt.Errorf("error while uploading file to MinIO: %w", err)
	}

	logger.Info("File uploaded to Minio successfully", zap.String("id:", id))

	return objectInfo.ETag, nil
}

func (d *DataRepository) DownloadFile(ctx context.Context, bucketName, objectName string) (*bytes.Buffer, string, error) {
	objectName = binData + "_" + objectName // The name for the object in MinIO

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		logger.Error("error creating minio client: ", zap.Error(err))

		return nil, "", errors.New("error instantiating Minio client with options")
	}

	obj, err := client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, "", fmt.Errorf("error downloading object from Minio: %w", err)
	}
	defer obj.Close()

	// Read the object data into a byte buffer
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, obj)
	if err != nil {
		logger.Error("copy file error: ", zap.Error(err))

		return nil, "", fmt.Errorf("copy file error: %w", err)
	}

	info, err := client.StatObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		logger.Error("error getting object metadata: ", zap.Error(err))

		return nil, "", fmt.Errorf("error getting object metadata: %w", err)
	}

	metadata := info.UserMetadata["info"]

	logger.Info("Object downloaded successfully:", zap.String("id:", objectName))

	return buf, metadata, nil
}
