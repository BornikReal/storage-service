package server

import (
	"context"
	"github.com/BornikReal/storage-service/internal/config"
	"github.com/BornikReal/storage-service/pkg/logger"
	"go.uber.org/zap"

	desc "github.com/BornikReal/storage-service/pkg/storage-service/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) SendData(ctx context.Context, req *desc.SendDataRequest) (*emptypb.Empty, error) {
	if i.replicaType == config.Master {
		return nil, status.Error(codes.Aborted, "master doesn't support SendData operation")
	}
	if i.password != req.Password {
		return nil, status.Error(codes.PermissionDenied, "incorrect password")
	}

	for k, v := range req.Data {
		err := i.kvService.Set(ctx, k, v)
		if err != nil {
			logger.Error("Set: error",
				zap.String("error", err.Error()),
				zap.String("key", k), zap.String("value", v))
			return nil, status.Errorf(codes.Internal, "Set: %v", err)
		}
	}

	return &emptypb.Empty{}, nil
}
