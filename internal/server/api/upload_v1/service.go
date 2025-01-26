package upload

import (
	"context"

	desc "github.com/igortoigildin/goph-keeper/pkg/upload_v1"
)

type UploadService interface {
	Upload(ctx context.Context, stream desc.UploadV1_UploadFileServer) error
}

type Implementation struct {
	desc.UnimplementedUploadV1Server
	uploadService UploadService
	ctx context.Context
}

func NewImplementation(ctx context.Context, uploadService UploadService) *Implementation {
	return &Implementation{
		ctx: ctx,
		uploadService: uploadService,
	}
}
