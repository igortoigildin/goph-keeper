package upload

import (
	"context"

	desc "github.com/igortoigildin/goph-keeper/pkg/upload_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) UploadFile(stream desc.UploadV1_UploadFileServer) error {
	err := i.uploadService.SaveFile(stream)
	if err != nil {
		return err
	}

	return nil
}

func (i *Implementation) UploadBankData(ctx context.Context, req *desc.UploadBankDataRequest) (*emptypb.Empty, error) {
	err := i.uploadService.SaveBankData(ctx, req.GetData())
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (i *Implementation) UploadPassword(ctx context.Context, req *desc.UploadPasswordRequest) (*emptypb.Empty, error) {
	err := i.uploadService.SaveLoginPassword(ctx, req.GetData())

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (i *Implementation) UploadText(ctx context.Context, req *desc.UploadTextRequest) (*emptypb.Empty, error) {
	err := i.uploadService.SaveText(ctx, req.GetText())
	if err != nil {
		return nil, err
	}

	return nil, nil
}



