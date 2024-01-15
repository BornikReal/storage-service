package server

import (
	"context"

	desc "github.com/BornikReal/storage-service/pkg/storage-service/pb"
)

type StorageService interface {
	Get(key string) (string, error)
	Set(ctx context.Context, key string, value string) error
	Subscribe(ip string, isAsync bool) error
}

type Implementation struct {
	desc.UnsafeStorageServiceServer

	storageType string
	replicaType string
	password    string
	kvService   StorageService
}

func NewImplementation(kvService StorageService, storageType string, replicaType string, password string) *Implementation {
	return &Implementation{
		storageType: storageType,
		replicaType: replicaType,
		password:    password,
		kvService:   kvService,
	}
}
