package list

import (
	"context"
	"time"

	desc "github.com/igortoigildin/goph-keeper/pkg/sync_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) GetObjectList(ctx context.Context, req *desc.SyncRequest) (*desc.SyncResponse, error) {
	objects, err := i.listService.List(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to list objects")
	}

	objs := make([]*desc.ObjectInfo, len(objects))
	for i, obj := range objects {
		objs[i] = &desc.ObjectInfo{
			Key:          obj.Key,
			Size:         obj.Size,
			LastModified: obj.LastModified.Format(time.RFC3339),
			Etag:         obj.ETag,
			Datatype:     obj.Datatype,
		}
	}

	return &desc.SyncResponse{Objects: objs}, nil
}
