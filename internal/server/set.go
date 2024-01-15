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

func (i *Implementation) Set(ctx context.Context, req *desc.SetRequest) (*emptypb.Empty, error) {
	if i.replicaType != config.Master {
		return nil, status.Error(codes.Aborted, "replica doesn't support Set operation")
	}

	err := i.kvService.Set(ctx, req.Key, req.Value)
	if err != nil {
		logger.Error("Set: error",
			zap.String("error", err.Error()),
			zap.String("key", req.Key), zap.String("value", req.Value))
		return nil, status.Errorf(codes.Internal, "Set: %v", err)
	}
	return &emptypb.Empty{}, nil
}
