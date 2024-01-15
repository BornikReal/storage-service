package storage_service

import desc "github.com/BornikReal/storage-service/pkg/storage-service/pb"

type ReplicaInfo struct {
	Client  desc.StorageServiceClient
	IP      string
	IsAsync bool
}
