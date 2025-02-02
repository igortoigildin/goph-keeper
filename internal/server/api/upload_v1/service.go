package upload

import (
	"context"

	desc "github.com/igortoigildin/goph-keeper/pkg/upload_v1"
)

type Saver interface {
	SaveFile(stream desc.UploadV1_UploadFileServer) error
	SaveBankData(ctx context.Context, data map[string]string) error
	SaveText(ctx context.Context, text string) error
	SaveLoginPassword(ctx context.Context, data map[string]string) error
}

type Implementation struct {
	desc.UnimplementedUploadV1Server
	uploadService Saver
}

func NewImplementation(uploadService Saver) *Implementation {
	return &Implementation{
		uploadService: uploadService,
	}
}
