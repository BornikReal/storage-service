package config

import (
	"os"
	"strconv"

	"github.com/BornikReal/storage-service/pkg/logger"
)

const (
	LSMStorage   = "lsm"
	RedisStorage = "redis"
)

const (
	Master       = "master"
	SyncReplica  = "sync_replica"
	AsyncReplica = "async_replica"
)

type Config struct {
	httpAddress string
	grpcAddress string

	ssDirectory string
	walPath     string
	walName     string

	maxTreeSize int
	blockSize   int64
	batch       int64
	ssChanSize  int64

	compressCronJob        string
	syncWithReplicaCronJob string

	redisHost     string
	redisPassword string

	storageType string
	replicaType string

	password string
}

func New() *Config {
	return &Config{}
}

func logUseDefault(varName string, stdValue interface{}) {
	logger.Infof("%s not found, using standard value: %v", varName, stdValue)
}

func (c *Config) LoadFromEnv() error {
	c.httpAddress = os.Getenv("SERVICE_HTTP_ADDRESS")
	if c.httpAddress == "" {
		c.httpAddress = "localhost:7001"
		logUseDefault("SERVICE_HTTP_ADDRESS", c.httpAddress)
	}

	c.grpcAddress = os.Getenv("SERVICE_GRPC_ADDRESS")
	if c.grpcAddress == "" {
		c.grpcAddress = "localhost:7002"
		logUseDefault("SERVICE_GRPC_ADDRESS", c.grpcAddress)
	}

	c.ssDirectory = os.Getenv("DB_DIR")
	if c.ssDirectory == "" {
		c.ssDirectory = "db"
		logUseDefault("DB_DIR", c.ssDirectory)
	}

	c.compressCronJob = os.Getenv("COMPRESS_CRON_JOB")
	if c.compressCronJob == "" {
		c.compressCronJob = "0 */1 * * *"
		logUseDefault("COMPRESS_CRON_JOB", c.compressCronJob)
	}

	c.syncWithReplicaCronJob = os.Getenv("SYNC_WITH_REPLICA_CRON_JOB")
	if c.syncWithReplicaCronJob == "" {
		c.syncWithReplicaCronJob = "* * * * *"
		logUseDefault("SYNC_WITH_REPLICA_CRON_JOB", c.syncWithReplicaCronJob)
	}

	value := os.Getenv("MAX_TREE_SIZE")
	if value == "" {
		c.maxTreeSize = 5
		logUseDefault("MAX_TREE_SIZE", c.maxTreeSize)
	} else {
		valueInt, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			c.maxTreeSize = 5
			logUseDefault("MAX_TREE_SIZE", c.maxTreeSize)
		} else {
			c.maxTreeSize = int(valueInt)
		}
	}

	value = os.Getenv("BLOCK_SIZE")
	if value == "" {
		c.blockSize = 5
		logUseDefault("BLOCK_SIZE", c.blockSize)
	} else {
		valueInt, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			c.blockSize = 5
			logUseDefault("BLOCK_SIZE", c.blockSize)
		} else {
			c.blockSize = valueInt
		}
	}

	value = os.Getenv("BATCH")
	if value == "" {
		c.batch = 1
		logUseDefault("BATCH", c.batch)
	} else {
		valueInt, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			c.batch = 1
			logUseDefault("BATCH", c.batch)
		} else {
			c.batch = valueInt
		}
	}

	c.walPath = os.Getenv("WAL_PATH")
	if c.walPath == "" {
		c.walPath = ""
		logUseDefault("WAL_PATH", c.walPath)
	}

	c.walName = os.Getenv("WAL_NAME")
	if c.walName == "" {
		c.walName = "wal"
		logUseDefault("WAL_NAME", c.walName)
	}

	value = os.Getenv("SS_CHAN_SIZE")
	if value == "" {
		c.ssChanSize = 5
		logUseDefault("SS_CHAN_SIZE", c.ssChanSize)
	} else {
		valueInt, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			c.ssChanSize = 5
			logUseDefault("SS_CHAN_SIZE", c.ssChanSize)
		} else {
			c.ssChanSize = valueInt
		}
	}

	c.redisHost = os.Getenv("REDIS_HOST")
	if c.redisHost == "" {
		c.redisHost = "172.28.1.4:6380"
		logUseDefault("REDIS_HOST", c.redisHost)
	}

	c.redisPassword = os.Getenv("REDIS_PASSWORD")
	if c.redisPassword == "" {
		c.redisPassword = "1234"
		logUseDefault("REDIS_PASSWORD", c.redisPassword)
	}

	c.storageType = os.Getenv("STORAGE_TYPE")
	if c.storageType == "" {
		c.storageType = LSMStorage
		logUseDefault("STORAGE_TYPE", c.storageType)
	}

	c.replicaType = os.Getenv("REPLICA_TYPE")
	if c.replicaType == "" {
		c.replicaType = Master
		logUseDefault("STORAGE_TYPE", c.replicaType)
	}

	c.password = os.Getenv("PASSWORD")
	if c.password == "" {
		c.password = "1234"
		logUseDefault("PASSWORD", c.password)
	}

	return nil
}

func (c *Config) GetHttpAddress() string {
	return c.httpAddress
}

func (c *Config) GetGrpcAddress() string {
	return c.grpcAddress
}

func (c *Config) GetSSDirectory() string {
	return c.ssDirectory
}

func (c *Config) GetWalPath() string {
	return c.walPath
}

func (c *Config) GetWalName() string {
	return c.walName
}

func (c *Config) GetCompressCronJob() string {
	return c.compressCronJob
}

func (c *Config) GetSyncWithReplicaCronJob() string {
	return c.syncWithReplicaCronJob
}

func (c *Config) GetMaxTreeSize() int {
	return c.maxTreeSize
}

func (c *Config) GetBlockSize() int64 {
	return c.blockSize
}

func (c *Config) GetBatch() int64 {
	return c.batch
}

func (c *Config) SSChanSize() int64 {
	return c.ssChanSize
}

func (c *Config) GetRedisHost() string {
	return c.redisHost
}

func (c *Config) GetRedisPassword() string {
	return c.redisPassword
}

func (c *Config) GetStorageType() string {
	return c.storageType
}

func (c *Config) GetReplicaType() string {
	return c.replicaType
}

func (c *Config) GetPassword() string {
	return c.password
}
