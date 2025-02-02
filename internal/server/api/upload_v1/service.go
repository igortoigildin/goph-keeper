package upload

import (
	"github.com/igortoigildin/goph-keeper/internal/server/service"
	desc "github.com/igortoigildin/goph-keeper/pkg/upload_v1"
)



type Implementation struct {
	desc.UnimplementedUploadV1Server
	uploadService service.UploadService
}

func NewImplementation(uploadService service.UploadService) *Implementation {
	return &Implementation{
		uploadService: uploadService,
	}
}
