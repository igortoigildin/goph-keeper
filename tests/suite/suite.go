package suite

import (
	"context"
	"os"
	"testing"

	config "github.com/igortoigildin/goph-keeper/internal/server/config"
	auth "github.com/igortoigildin/goph-keeper/pkg/auth_v1"
	download "github.com/igortoigildin/goph-keeper/pkg/download_v1"
	upload "github.com/igortoigildin/goph-keeper/pkg/upload_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	T *testing.T
	Cfg *config.Config
	AuthClient	auth.AuthV1Client
	UploadClient upload.UploadV1Client
	DownloadClient download.DownloadV1Client
}

const (
	grpcHost = "localhost"
)

// New creates new test suite.
func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadPath(configPath()) 

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.DialContext(context.Background(), 
		cfg.GRPC.Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal("grpc server connection failed:", err)
	}

	return ctx, &Suite{
		T: 		t,
		Cfg: 	cfg,
		AuthClient: auth.NewAuthV1Client(cc),
		UploadClient: upload.NewUploadV1Client(cc),
		DownloadClient: download.NewDownloadV1Client(cc),
	}
}

func configPath() string {
	const key = "CONFIG_PATH"

	if v := os.Getenv(key); v != "" {
		return v
	}

	return "../config/local_tests.yaml"
}

