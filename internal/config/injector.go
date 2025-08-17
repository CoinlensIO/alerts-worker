package config

import (
	"alerts-worker/internal/constants"
	"alerts-worker/internal/repository"
	"alerts-worker/internal/service"
	"alerts-worker/pkg/metrics"
	"alerts-worker/pkg/queue"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"time"

	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/samber/do"
)

func NewInjector(cfg *Config) *do.Injector {
	injector := do.New()

	do.Provide(
		injector, func(i *do.Injector) (*zerolog.Logger, error) {
			logLevel, err := zerolog.ParseLevel(cfg.LogLevel)
			if err != nil {
				return nil, err
			}

			logger := zerolog.New(os.Stdout).
				Level(logLevel).
				With().
				Timestamp().
				Logger()

			return &logger, nil
		},
	)

	//Databases
	do.Provide(injector, func(i *do.Injector) (*pgxpool.Pool, error) {
		dsn := fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.PostgresqlUsername, cfg.PostgresqlPassword, cfg.PostgresqlHost, cfg.PostgresqlPort, cfg.PostgresqlDatabaseName,
		)

		dbpool, err := pgxpool.New(context.Background(), dsn)
		if err != nil {
			log.Fatalf("Unable to connect to database: %v\n", err)
		}

		return dbpool, nil
	})

	do.Provide(injector, func(i *do.Injector) (*gorm.DB, error) {
		pool := do.MustInvoke[*pgxpool.Pool](i)

		sqlDB := stdlib.OpenDBFromPool(pool)

		gormDB, err := gorm.Open(postgres.New(postgres.Config{
			Conn: sqlDB,
		}), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("unable to connect to database with GORM: %w", err)
		}

		return gormDB, nil
	})

	do.Provide(injector, func(i *do.Injector) (*repository.Repository, error) {
		db := do.MustInvoke[*gorm.DB](i)
		return repository.NewRepository(db), nil
	})

	do.ProvideNamed(injector, "BinanceMarkPriceAlerts", func(i *do.Injector) (*redis.Client, error) {
		client := redis.NewClient(&redis.Options{
			Addr:         cfg.BinanceMarkPricesRedis,
			DialTimeout:  60 * time.Second,
			ReadTimeout:  60 * time.Second,
			WriteTimeout: 60 * time.Second,
			PoolSize:     100000,
			MaxRetries:   5,
		})

		return client, nil
	})

	do.Provide(injector, func(i *do.Injector) (*metrics.WorkerMetrics, error) {
		workerMetrics := metrics.InitWorkerMetrics()

		return workerMetrics, nil
	})

	do.ProvideNamed(injector, string(constants.BinanceMarkPriceAlertsQueue), func(i *do.Injector) (*queue.Queue, error) {
		redisClient := do.MustInvokeNamed[*redis.Client](i, "BinanceMarkPriceAlerts")
		batchQueue := queue.NewQueue(redisClient, string(constants.BinanceMarkPriceAlertsQueue))

		return batchQueue, nil
	})

	do.Provide(injector, func(i *do.Injector) (service.AlertService, error) {
		userRepo := do.MustInvoke[*repository.Repository](i)

		return service.New(userRepo), nil
	})

	return injector
}
