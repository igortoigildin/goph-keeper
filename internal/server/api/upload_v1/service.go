package upload

import (
	desc "github.com/igortoigildin/goph-keeper/pkg/upload_v1"
)

type UploadService interface {
	Upload(stream desc.FileService_UploadServer) error
}

type Implementation struct {
	desc.UnimplementedFileServiceServer
	uploadService UploadService
}

func NewImplementation(uploadService UploadService) *Implementation {
	return &Implementation{
		uploadService: uploadService,
	}
}
