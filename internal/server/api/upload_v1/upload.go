package upload

import (
	desc "github.com/igortoigildin/goph-keeper/pkg/upload_v1"
)

func (i *Implementation) UploadFile(stream desc.UploadV1_UploadFileServer) error {
	err := i.uploadService.Upload(stream)
	if err != nil {
		return err
	}

	return nil
}
