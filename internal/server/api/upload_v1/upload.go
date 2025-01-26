package upload

import (
	"context"

	desc "github.com/igortoigildin/goph-keeper/pkg/upload_v1"
)

func (i *Implementation) UploadFile(ctx context.Context, stream desc.UploadV1_UploadFileServer) error {
	err := i.uploadService.Upload(ctx, stream)
	if err != nil {
		return err
	}

	return nil
}
