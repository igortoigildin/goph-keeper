package download

import (
	"github.com/igortoigildin/goph-keeper/internal/server/service"
	desc "github.com/igortoigildin/goph-keeper/pkg/download_v1"
)


const (
	filePath = "./"
)
type Implementation struct {
	desc.UnimplementedDownloadV1Server
	downloadService service.DownloadService
}


func NewImplementation(downloadService service.DownloadService) *Implementation {
	return &Implementation{
		downloadService: downloadService,
	}
}
