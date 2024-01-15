package storage_service

import (
	"context"
	"fmt"
	desc "github.com/BornikReal/storage-service/pkg/storage-service/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Storage interface {
	Get(key string) (string, error)
	Set(key string, value string) error
}

type StorageService struct {
	replicas []ReplicaInfo

	storage  Storage
	password string
}

func NewStorageService(storage Storage, password string) *StorageService {
	return &StorageService{
		storage:  storage,
		password: password,
	}
}

func (r *StorageService) Get(key string) (string, error) {
	value, err := r.storage.Get(key)
	if err != nil {
		return "", fmt.Errorf("resposiotry.Get: %w", err)
	}
	return value, nil
}

func (r *StorageService) Set(ctx context.Context, key string, value string) error {
	err := r.storage.Set(key, value)
	if err != nil {
		return fmt.Errorf("resposiotry.Set: %w", err)
	}

	for _, repl := range r.replicas {
		if repl.IsAsync {
			continue
		}
		_, err = repl.Client.SendData(ctx, &desc.SendDataRequest{
			Data:     map[string]string{key: value},
			Password: r.password,
		})
		if err != nil {
			return fmt.Errorf("repl.Client.Set: %w", err)
		}
	}
	return nil
}

func (r *StorageService) Subscribe(ip string, isAsync bool) error {
	conn, err := grpc.Dial(
		ip,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("grpc.Dial: %w", err)
	}

	r.replicas = append(r.replicas, ReplicaInfo{
		Client:  desc.NewStorageServiceClient(conn),
		IP:      ip,
		IsAsync: isAsync,
	})
	return nil
}

func (r *StorageService) GetReplicaList() []ReplicaInfo {
	return r.replicas
}
