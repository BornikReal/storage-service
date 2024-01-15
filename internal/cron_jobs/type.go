package cron_jobs

import "github.com/BornikReal/storage-service/internal/storage_service"

type StorageWithWal interface {
	GetWalElements(init bool) (map[string]string, error)
}

type StorageService interface {
	GetReplicaList() []storage_service.ReplicaInfo
}

type SSManager interface {
	CompressSS() error
}

type CronJobInfo struct {
	Cron                  string
	SupportedStorageTypes []string
	SupportedReplicaTypes []string
	Run                   func()
	JobName               string
}
