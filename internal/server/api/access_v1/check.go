package access

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"

	descAccess "github.com/igortoigildin/goph-keeper/pkg/access_v1"
	"github.com/igortoigildin/goph-keeper/pkg/logger"
)

func (i *Implementation) Check(ctx context.Context, req *descAccess.CheckRequest) (*emptypb.Empty, error) {
	err := i.accessService.Check(ctx, req.EndpointAddress)
	if err != nil {
		logger.Error("endpoint access error", zap.Error(err))
		return nil, fmt.Errorf("endpoint access error: %w", err)
	}

	return nil, nil
}
