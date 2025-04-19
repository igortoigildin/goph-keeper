package minio

import (
	"bytes"
	"context"
	"encoding/json"
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
)

type DataRepository struct{}

func NewRepository() *DataRepository {
	return &DataRepository{}
}

func (d *DataRepository) SaveFile(ctx context.Context, file *fl.File, login string, id string, meta string) error {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		logger.Error("error while creating minio client: ", zap.Error(err))

		return fmt.Errorf("error instantiating Minio client with options: %w", err)
	}

	// Define the file to upload and the destination bucket
	objectName := id    // The name for the object in MinIO
	bucketName := login // Bucket name in MinIO

	// Ensure the bucket exists (or create it)
	err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		if exists, errBucketExists := client.BucketExists(context.Background(), bucketName); errBucketExists == nil && exists {
			logger.Info("Bucket already exists")
		} else {
			logger.Info("Failed to create bucket:", zap.Error(err))

			return fmt.Errorf("Minio error: %w", err)
		}
	}

	meatadata := map[string]string{
		"meta": meta,
	}

	// Open the file to upload
	f, err := os.Open(file.FilePath)
	if err != nil {
		logger.Error("error opening targeted file: ", zap.Error(err))

		return fmt.Errorf("error opening targeted file: %w", err)
	}
	defer f.Close()

	// Upload the file to MinIO
	_, err = client.PutObject(
		context.Background(),
		bucketName,
		objectName,
		f,
		-1,
		minio.PutObjectOptions{ContentType: "application/octet-stream", UserMetadata: meatadata},
	)
	if err != nil {
		logger.Error("error while uploading file to MinIO", zap.Error(err))

		return fmt.Errorf("error while uploading file to MinIO: %w", err)
	}

	logger.Info("File uploaded to Minio successfully", zap.String("id:", id))

	return nil
}

func (d *DataRepository) DownloadFile(ctx context.Context, bucketName, objectName string) (*bytes.Buffer, string, error) {
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

func (d *DataRepository) SaveTextData(ctx context.Context, data any, login string, id string, info string) error {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		logger.Error("error while creating minio client: ", zap.Error(err))

		return fmt.Errorf("error instantiating Minio client with options: %w", err)
	}

	// Serialize the map to JSON
	serializedData, err := json.Marshal(data)
	if err != nil {
		logger.Error("error while serializing the map: ", zap.Error(err))

		return fmt.Errorf("serialization error: %w", err)
	}

	// Create a buffer from the serialized data
	buf := bytes.NewReader(serializedData)

	// Define the file to upload and the destination bucket
	objectName := id    // The name for the object in MinIO
	bucketName := login // Bucket name in MinIO

	// Ensure the bucket exists (or create it)
	err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		if exists, errBucketExists := client.BucketExists(context.Background(), bucketName); errBucketExists == nil && exists {
			logger.Info("Bucket already exists")
		} else {
			logger.Error("Failed to create bucket:", zap.Error(err))

			return fmt.Errorf("Minio error: %w", err)
		}
	}

	// Save additional info about data to be saved
	meatadata := map[string]string{
		"info": info,
	}

	_, err = client.PutObject(ctx, bucketName, objectName, buf,
		int64(buf.Len()),
		minio.PutObjectOptions{ContentType: "application/json", UserMetadata: meatadata})
	if err != nil {
		logger.Error("error while uploading object to minio: ", zap.Error(err))

		return fmt.Errorf("Minio error: %w", err)
	}

	logger.Info("String data uploaded to Minio successfully:", zap.String("id:", id))

	return nil
}

func (d *DataRepository) DownloadTextData(ctx context.Context, bucketName, objectName string) ([]byte, string, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		logger.Error("error while creating minio client: ", zap.Error(err))

		return nil, "", fmt.Errorf("error instantiating Minio client with options: %w", err)
	}

	obj, err := client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		logger.Error("error opening targeted file: ", zap.Error(err))

		return nil, "", fmt.Errorf("error opening targeted file: %w", err)
	}
	defer obj.Close()

	info, err := client.StatObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		logger.Error("error getting object metadata: ", zap.Error(err))

		return nil, "", fmt.Errorf("error getting object metadata: %w", err)
	}

	metadata := info.UserMetadata["info"]

	// Read the object data into a byte buffer
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, obj)
	if err != nil {
		logger.Error("error copying targeted file: ", zap.Error(err))

		return nil, "", fmt.Errorf("error copying targeted file: %w", err)
	}

	res := buf.Bytes()

	logger.Info("Object downloaded successfully: ", zap.String("id:", objectName))

	return res, metadata, nil
}
