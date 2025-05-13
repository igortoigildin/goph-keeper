package minio

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

func (d *DataRepository) SaveTextData(ctx context.Context, data any, login string, id string, info string, datatype string) (string, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		logger.Error("error while creating minio client: ", zap.Error(err))

		return "", fmt.Errorf("error instantiating Minio client with options: %w", err)
	}

	// Serialize the map to JSON
	serializedData, err := json.Marshal(data)
	if err != nil {
		logger.Error("error while serializing the map: ", zap.Error(err))

		return "", fmt.Errorf("serialization error: %w", err)
	}

	// Create a buffer from the serialized data
	buf := bytes.NewReader(serializedData)

	// Define the file to upload and the destination bucket
	objectName := datatype + "_" + id // The name for the object in MinIO
	bucketName := login               // Bucket name in MinIO

	// Ensure the bucket exists (or create it)
	err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		if exists, errBucketExists := client.BucketExists(ctx, bucketName); errBucketExists == nil && exists {
			logger.Info("Bucket already exists")
		} else {
			logger.Error("Failed to create bucket:", zap.Error(err))

			return "", fmt.Errorf("Minio error: %w", err)
		}
	}

	// Save additional info about data to be saved
	metadata := map[string]string{
		"info":     info,
		"datatype": datatype,
	}

	objInfo, err := client.PutObject(ctx, bucketName, objectName, buf,
		int64(buf.Len()),
		minio.PutObjectOptions{ContentType: "application/json", UserMetadata: metadata})

	if err != nil {
		logger.Error("error while uploading object to minio: ", zap.Error(err))

		return "", fmt.Errorf("Minio error: %w", err)
	}

	logger.Info("String data uploaded to Minio successfully:", zap.String("id:", id))

	return objInfo.ETag, nil
}

func (d *DataRepository) DownloadTextData(ctx context.Context, bucketName, objectName, dataType string) ([]byte, string, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		logger.Error("error while creating minio client: ", zap.Error(err))
		return nil, "", fmt.Errorf("error instantiating Minio client with options: %w", err)
	}

	objectName = dataType + "_" + objectName

	fmt.Println("OBJECTNAME", objectName)

	// Get the object
	obj, err := client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		logger.Error("error opening targeted file: ", zap.Error(err))
		return nil, "", fmt.Errorf("error opening targeted file: %w", err)
	}
	defer obj.Close()

	// Read the object data into a byte buffer
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, obj)
	if err != nil {
		logger.Error("error copying targeted file: ", zap.Error(err))
		return nil, "", fmt.Errorf("error copying targeted file: %w", err)
	}

	res := buf.Bytes()

	return res, "", nil
}
