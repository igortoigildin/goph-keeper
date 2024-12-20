package access

import (
	"context"

	desc "github.com/igortoigildin/goph-keeper/pkg/access_v1"
)

type AccessService interface {
	Check(ctx context.Context, endpoint string) error
}

type Implementation struct {
	desc.UnimplementedAccessV1Server
	accessService AccessService
}

func NewImplementation(accessService AccessService) *Implementation {
	return &Implementation{
		accessService: accessService,
	}
}
