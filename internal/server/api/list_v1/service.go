package list

import (
	"github.com/igortoigildin/goph-keeper/internal/server/service"
	desc "github.com/igortoigildin/goph-keeper/pkg/sync_v1"
)

type Implementation struct {
	desc.UnimplementedSyncV1Server
	listService service.ListService
}

func NewImplementation(listService service.ListService) *Implementation {
	return &Implementation{
		listService: listService,
	}
}
