package server

import (
	"context"

	"github.com/BornikReal/storage-service/internal/config"
	"github.com/BornikReal/storage-service/pkg/logger"
	desc "github.com/BornikReal/storage-service/pkg/storage-service/pb"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) Subscribe(_ context.Context, req *desc.SubscribeRequest) (*emptypb.Empty, error) {
	if i.replicaType != config.Master {
		return nil, status.Error(codes.InvalidArgument, "Can't subscribe to replica")
	}
	if i.storageType != config.LSMStorage {
		return &emptypb.Empty{}, nil
	}
	if i.password != req.Password {
		return nil, status.Error(codes.PermissionDenied, "incorrect password")
	}

	err := i.kvService.Subscribe(req.Ip, req.IsAsync)
	if err != nil {
		logger.Error("Subscribe: error", zap.String("error", err.Error()), zap.String("ip", req.Ip), zap.Bool("is_async", req.IsAsync))
		return nil, status.Errorf(codes.Internal, "Subscribe: %v", err)
	}
	return &emptypb.Empty{}, nil
}
