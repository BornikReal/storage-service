package cron_jobs

import (
	"context"

	"github.com/BornikReal/storage-service/pkg/logger"
	desc "github.com/BornikReal/storage-service/pkg/storage-service/pb"
	"go.uber.org/zap"
)

type SyncWithReplicaJob struct {
	storageWithWal StorageWithWal
	storageService StorageService

	password string
}

func NewSyncWithReplicaJob(storageWithWal StorageWithWal, storageService StorageService, password string) *SyncWithReplicaJob {
	return &SyncWithReplicaJob{
		storageWithWal: storageWithWal,
		storageService: storageService,
		password:       password,
	}
}

func (j *SyncWithReplicaJob) Name() string {
	return "sync with replica"
}

func (j *SyncWithReplicaJob) Run(ctx context.Context) {
	defer logger.Info("job finished",
		zap.String(logger.JobNameField, j.Name()),
	)
	logger.Info("job started",
		zap.String(logger.JobNameField, j.Name()),
	)

	res, err := j.storageWithWal.GetWalElements(true)
	if err != nil {
		logger.Error("Getting wal elements finished with error",
			zap.String(logger.ErrorField, err.Error()),
			zap.String(logger.JobNameField, j.Name()),
		)
	}

	replicas := j.storageService.GetReplicaList()
	for _, repl := range replicas {
		if !repl.IsAsync {
			continue
		}
		_, err = repl.Client.SendData(ctx, &desc.SendDataRequest{
			Data:     res,
			Password: j.password,
		})
		if err != nil {
			logger.Error("Setting wal elements finished with error",
				zap.String(logger.ErrorField, err.Error()),
				zap.String(logger.JobNameField, j.Name()),
			)
			continue
		}
	}
}
