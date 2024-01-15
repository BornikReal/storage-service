package internal

import (
	"context"
	"github.com/BornikReal/storage-service/internal/cron_jobs"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/BornikReal/storage-component/pkg/ss_storage/iterator"
	"github.com/BornikReal/storage-component/pkg/ss_storage/kv_file"
	"github.com/BornikReal/storage-component/pkg/ss_storage/ss_manager"
	"github.com/BornikReal/storage-component/pkg/storage"
	"github.com/BornikReal/storage-component/pkg/tree_with_clone"
	"github.com/BornikReal/storage-service/internal/config"
	"github.com/BornikReal/storage-service/internal/server"
	"github.com/BornikReal/storage-service/internal/storage_service"
	"github.com/BornikReal/storage-service/pkg/logger"
	"github.com/BornikReal/storage-service/pkg/storage-service/pb"
	"github.com/emirpasic/gods/trees/avltree"
	"github.com/go-co-op/gocron"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type serve func() error

type App struct {
	wg        sync.WaitGroup
	config    *config.Config
	startHttp serve
	startGrpc serve
}

func NewApp() *App {
	return &App{}
}

func (app *App) Init() error {
	logger.InitLogger()
	logger.Info("init service")
	defer logger.Info("init finished")

	ctx := context.Background()

	conf := config.New()

	if err := conf.LoadFromEnv(); err != nil {
		return err
	}

	var (
		mt storage_service.Storage

		ssManager      cron_jobs.SSManager
		storageWithWal *storage.MemTableWithWal
	)
	switch conf.GetStorageType() {
	case config.LSMStorage:
		storageWithWal, ssManager = initLSMStorage(conf)
		mt = storageWithWal
	case config.RedisStorage:
		mt = initRedis(conf)
	}

	if mt == nil {
		logger.Fatal("storage not init",
			zap.String("storage_type", conf.GetStorageType()))
	}

	storageService := storage_service.NewStorageService(mt, conf.GetPassword())
	impl := server.NewImplementation(storageService, conf.GetStorageType(), conf.GetReplicaType(), conf.GetPassword())

	app.config = conf

	initCronJobs(conf, ssManager, storageWithWal, storageService)
	app.initGrpc(impl)
	app.initHttp(ctx)

	return nil
}

func (app *App) Run() {
	logger.Info("service is starting")
	defer logger.Info("service shutdown")
	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		if err := app.startGrpc(); err != nil {
			logger.Fatal("starting grpc storage_service ended with error",
				zap.String("error", err.Error()), zap.String("port", app.config.GetGrpcAddress()))
		}
	}()

	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		if err := app.startHttp(); err != nil {
			logger.Fatal("starting http storage_service ended with error",
				zap.String("error", err.Error()), zap.String("port", app.config.GetHttpAddress()))
		}
	}()

	logger.Infof("Service successfully started. Addresses: HTTP - %s, GRPC - %s",
		app.config.GetHttpAddress(), app.config.GetGrpcAddress())
	app.wg.Wait()
}

func (app *App) initGrpc(service *server.Implementation) {
	logger.Info("init grpc storage_service")
	grpcServer := grpc.NewServer()
	pb.RegisterStorageServiceServer(grpcServer, service)
	reflection.Register(grpcServer)
	lsn, err := net.Listen("tcp", app.config.GetGrpcAddress())
	if err != nil {
		logger.Fatal("listening port ended with error",
			zap.String("error", err.Error()), zap.String("port", app.config.GetGrpcAddress()))
	}

	app.startGrpc = func() error {
		return grpcServer.Serve(lsn)
	}
}

func (app *App) initHttp(ctx context.Context) {
	logger.Info("init http storage_service")
	serveMux := runtime.NewServeMux()
	opt := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pb.RegisterStorageServiceHandlerFromEndpoint(ctx, serveMux, app.config.GetGrpcAddress(), opt)
	if err != nil {
		logger.Fatal("can't create http storage_service from grpc endpoint",
			zap.String("error", err.Error()))
	}

	app.startHttp = func() error {
		return http.ListenAndServe(app.config.GetHttpAddress(), serveMux)
	}
}

func initLSMStorage(conf *config.Config) (*storage.MemTableWithWal, cron_jobs.SSManager) {
	ssManager := ss_manager.NewSSManager(conf.GetSSDirectory(), conf.GetBlockSize(), conf.GetBatch())
	if err := ssManager.Init(); err != nil {
		panic(err)
	}
	tree := avltree.NewWithStringComparator()
	wal := kv_file.NewKVFile(conf.GetWalPath(), conf.GetWalName())
	if err := wal.Init(); err != nil {
		panic(err)
	}

	dumper := make(chan iterator.Iterator, conf.SSChanSize())
	mt := storage.NewMemTableWithWal(
		storage.NewMemTableWithSS(
			storage.NewMemTable(
				tree_with_clone.NewTreeWithClone(
					tree,
				),
				dumper,
				conf.GetMaxTreeSize(),
			),
			ssManager,
		),
		wal,
	)

	errorCh := make(chan error, 1)
	ssProcessor := storage.NewSSProcessor(ssManager, errorCh)
	go ssProcessor.Start(dumper)
	go func() {
		for err := range errorCh {
			logger.Error("SS processor encounters with error while saving tree", zap.String("error", err.Error()))
		}
	}()

	//s := gocron.NewScheduler(time.UTC)
	//_, err := s.Cron(conf.GetCompressCronJob()).Do(ssManager.CompressSS)
	//if err != nil {
	//	panic(err)
	//}
	//s.StartAsync()

	return mt, ssManager
}

func initRedis(conf *config.Config) storage_service.Storage {
	rdbMaster := redis.NewClient(&redis.Options{
		Addr:     conf.GetRedisHost(),
		Password: conf.GetRedisPassword(),
	})

	return storage.NewRedisStorage(rdbMaster, nil)
}

func initCronJobs(conf *config.Config,
	ssManager cron_jobs.SSManager,
	storageWithWal cron_jobs.StorageWithWal,
	storageService cron_jobs.StorageService,
) {
	syncWithReplicaJob := cron_jobs.NewSyncWithReplicaJob(storageWithWal, storageService, conf.GetPassword())
	compressSSJob := cron_jobs.NewCompressSSJob(ssManager)
	//compressSSJobName := "compress ss"

	jobsList := []cron_jobs.CronJobInfo{
		{
			Cron:                  conf.GetCompressCronJob(),
			SupportedStorageTypes: []string{config.LSMStorage},
			Run: func() {
				compressSSJob.Run(context.Background())
			},
			JobName: compressSSJob.Name(),
		},
		//conf.GetCompressCronJob(): {
		//	SupportedStorageTypes: []string{config.LSMStorage},
		//	Run: func() {
		//		defer logger.Info("job finished",
		//			zap.String(logger.JobNameField, compressSSJobName),
		//		)
		//		logger.Info("job started",
		//			zap.String(logger.JobNameField, compressSSJobName),
		//		)
		//
		//		if err := ssManager.CompressSS(); err != nil {
		//			logger.Error("ssManager.CompressSS error",
		//				zap.String(logger.ErrorField, err.Error()),
		//				zap.String(logger.JobNameField, compressSSJobName),
		//			)
		//		}
		//	},
		//	JobName: compressSSJobName,
		//},
		{
			Cron:                  conf.GetSyncWithReplicaCronJob(),
			SupportedStorageTypes: []string{config.LSMStorage},
			SupportedReplicaTypes: []string{config.Master},
			Run: func() {
				syncWithReplicaJob.Run(context.Background())
			},
			JobName: syncWithReplicaJob.Name(),
		},
	}

	s := gocron.NewScheduler(time.UTC)
	for _, jobInfo := range jobsList {
		ok := false
		for _, st := range jobInfo.SupportedStorageTypes {
			if st == conf.GetStorageType() {
				ok = true
			}
		}

		if len(jobInfo.SupportedStorageTypes) == 0 {
			ok = true
		}

		if !ok {
			continue
		}

		ok = false
		for _, st := range jobInfo.SupportedReplicaTypes {
			if st == conf.GetReplicaType() {
				ok = true
			}
		}

		if len(jobInfo.SupportedReplicaTypes) == 0 {
			ok = true
		}

		if !ok {
			continue
		}

		_, err := s.Cron(jobInfo.Cron).Do(jobInfo.Run)
		if err != nil {
			panic(err)
		}
		logger.Info("registered job",
			zap.String("job_name", jobInfo.JobName),
			zap.String("cron", jobInfo.Cron),
		)
	}

	s.StartAsync()
}
