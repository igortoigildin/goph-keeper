package service

import (
	desc "github.com/igortoigildin/goph-keeper/pkg/upload_v1"
)

type UploadService interface {
	Upload(stream desc.FileService_UploadServer) error
}
