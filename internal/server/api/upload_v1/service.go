package upload

import (
	desc "github.com/igortoigildin/goph-keeper/pkg/upload_v1"
)

type UploadService interface {
	Upload(stream desc.UploadV1_UploadFileServer) error
}

type Implementation struct {
	desc.UnimplementedUploadV1Server
	uploadService UploadService
}

func NewImplementation(uploadService UploadService) *Implementation {
	return &Implementation{
		uploadService: uploadService,
	}
}
