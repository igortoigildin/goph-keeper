package upload

import (
	"context"

	desc "github.com/igortoigildin/goph-keeper/pkg/upload_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) UploadFile(stream desc.UploadV1_UploadFileServer) (error) {
	err := i.uploadService.SaveFile(stream)
	if err != nil {
		return status.Error(codes.Unknown, "failed to upload file")
	}

	return nil
}

func (i *Implementation) UploadBankData(
	ctx context.Context,
	req *desc.UploadBankDataRequest,
) (*desc.UploadBankDataResponse, error) {
	etag, err := i.uploadService.SaveBankData(ctx, req.GetData(), req.Metadata)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to upload bank data")
	}

	return &desc.UploadBankDataResponse{Etag: etag}, nil
}

func (i *Implementation) UploadPassword(
	ctx context.Context,
	req *desc.UploadPasswordRequest,
) (*desc.UploadPasswordResponse, error) {
	etag, err := i.uploadService.SaveLoginPassword(ctx, req.GetData(), req.Metadata)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to upload credentials")
	}

	return &desc.UploadPasswordResponse{Etag: etag}, nil
}

func (i *Implementation) UploadText(
	ctx context.Context,
	req *desc.UploadTextRequest,
) (*desc.UploadTextResponse, error) {
	etag, err := i.uploadService.SaveText(ctx, req.GetText(), req.GetMetadata())
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to upload text")
	}

	return &desc.UploadTextResponse{Etag: etag}, nil
}
