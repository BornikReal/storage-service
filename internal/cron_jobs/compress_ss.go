package cron_jobs

import (
	"context"

	"github.com/BornikReal/storage-service/pkg/logger"
	"go.uber.org/zap"
)

type CompressSSJob struct {
	SsManager SSManager
}

func NewCompressSSJob(ssManager SSManager) *CompressSSJob {
	return &CompressSSJob{
		SsManager: ssManager,
	}
}

func (j *CompressSSJob) Name() string {
	return "compress ss"
}

func (j *CompressSSJob) Run(_ context.Context) {
	defer logger.Info("job finished",
		zap.String(logger.JobNameField, j.Name()),
	)
	logger.Info("job started",
		zap.String(logger.JobNameField, j.Name()),
	)

	if err := j.SsManager.CompressSS(); err != nil {
		logger.Error("ssManager.CompressSS error",
			zap.String(logger.ErrorField, err.Error()),
			zap.String(logger.JobNameField, j.Name()),
		)
	}
}
