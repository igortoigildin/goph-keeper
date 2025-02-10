package download

import (
	"context"
	"fmt"

	desc "github.com/igortoigildin/goph-keeper/pkg/download_v1"
	"github.com/igortoigildin/goph-keeper/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)


func (i *Implementation) DownloadBankData(ctx context.Context, req *desc.DownloadBankDataRequest) (*desc.DownloadBankDataResponse, error) {
	res, err := i.downloadService.DownloadBankData(ctx, req.GetUuid())
	if err != nil {
		logger.Error("error downloading card details:", zap.Error(err))
		return nil, err
	}

	return &desc.DownloadBankDataResponse{
		Data: res,
	}, nil
}

func (i *Implementation) DownloadPassword(ctx context.Context, req *desc.DownloadPasswordRequest) (*desc.DownloadPasswordResponse, error) {
	fmt.Println(req.GetUuid())

	res, err := i.downloadService.DownloadLoginPassword(ctx, req.GetUuid())
	if err != nil {
		logger.Error("error downloading pass details:", zap.Error(err))
		return nil, err
	}

	return &desc.DownloadPasswordResponse{
		Data: res,
	}, nil
}

func (i *Implementation) DownloadText(ctx context.Context, req *desc.DownloadTextRequest) (*desc.DownloadTextResponse, error) {
	res, err := i.downloadService.DownloadText(ctx, req.GetUuid())
	if err != nil {
		logger.Error("error downloading text:", zap.Error(err))
		return nil, err
	}

	return &desc.DownloadTextResponse{
		Text: res,
	}, nil
}

func (i *Implementation) DownloadFile(req *desc.DownloadFileRequest, stream grpc.ServerStreamingServer[desc.DownloadFileResponse]) error {
	res, err := i.downloadService.DownloadFile(stream.Context(), req.GetUuid())
	if err != nil {
		logger.Error("error downloading file:", zap.Error(err))
		return err
	}


	return stream.Send(&desc.DownloadFileResponse{Uuid:  req.GetUuid(), Chunk: res})
}