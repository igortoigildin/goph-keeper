package download

import (
	"context"

	desc "github.com/igortoigildin/goph-keeper/pkg/download_v1"
	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) DownloadBankData(ctx context.Context, req *desc.DownloadBankDataRequest) (*desc.DownloadBankDataResponse, error) {
	res, metadata, err := i.downloadService.DownloadBankData(ctx, req.GetUuid())
	if err != nil {
		logger.Error("error downloading card details:", zap.Error(err))

		return nil, status.Error(codes.Unknown, "failed to download bank data")
	}

	return &desc.DownloadBankDataResponse{
		Data: res,
		Metadata: metadata,
	}, nil
}

func (i *Implementation) DownloadPassword(ctx context.Context, req *desc.DownloadPasswordRequest) (*desc.DownloadPasswordResponse, error) {
	res, metadata, err := i.downloadService.DownloadLoginPassword(ctx, req.GetUuid())
	if err != nil {
		logger.Error("error downloading pass details:", zap.Error(err))

		return nil, status.Error(codes.Unknown, "failed to download credentials")
	}

	return &desc.DownloadPasswordResponse{
		Data: res,
		Metadata: metadata,
	}, nil
}

func (i *Implementation) DownloadText(ctx context.Context, req *desc.DownloadTextRequest) (*desc.DownloadTextResponse, error) {
	res, metadata, err := i.downloadService.DownloadText(ctx, req.GetUuid())
	if err != nil {
		logger.Error("error downloading text:", zap.Error(err))

		return nil, status.Error(codes.Unknown, "failed to download text")
	}

	return &desc.DownloadTextResponse{
		Text: res,
		Metadata: metadata,
	}, nil
}

func (i *Implementation) DownloadFile(req *desc.DownloadFileRequest, stream grpc.ServerStreamingServer[desc.DownloadFileResponse]) error {
	res, metadata, err := i.downloadService.DownloadFile(stream.Context(), req.GetUuid())
	if err != nil {
		logger.Error("error downloading file:", zap.Error(err))

		return status.Error(codes.Unknown, "failed to download bin file")
	}

	return stream.Send(&desc.DownloadFileResponse{Uuid: req.GetUuid(), Chunk: res, Metadata: metadata})
}
