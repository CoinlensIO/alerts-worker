package main

import (
	"alerts-worker/internal/config"
	"alerts-worker/internal/constants"
	"alerts-worker/internal/event_handler"
	"alerts-worker/internal/service"
	"alerts-worker/pkg/metrics"
	"alerts-worker/pkg/worker"
	"context"
	"errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/samber/do"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	appBase := config.New(
		config.Init(),
		config.WithDependencyInjector(),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := appBase.Injector.HealthCheck()
	for name, err := range err {
		if err != nil {
			log.Panic().Err(err).Msgf("%s failed to initialize", name)
		}
	}

	redisClient := do.MustInvokeNamed[*redis.Client](appBase.Injector, "BinanceMarkPriceAlerts")
	logger := do.MustInvoke[*zerolog.Logger](appBase.Injector)

	// Configure retry behavior
	retryConfig := &event_handler.RetryConfig{
		MaxRetries:     10,
		InitialBackoff: 1 * time.Second,
		MaxBackoff:     3 * time.Minute,
		BackoffFactor:  1.5,
	}

	metricsServer := &http.Server{
		Addr:    ":2113",
		Handler: promhttp.Handler(),
	}
	go func() {
		if err := metricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Metrics %s error: %v", appBase.Config.ServiceName, err)
		}
	}()

	svc := do.MustInvoke[service.AlertService](appBase.Injector)
	workerMetrics := do.MustInvoke[*metrics.WorkerMetrics](appBase.Injector)

	handlerOpts := &event_handler.EventHandlerOptions{
		RetryConfig: retryConfig,
	}

	eventHandler := event_handler.NewEventHandler(svc, logger, handlerOpts)

	klinesSyncWorker := worker.NewWorker(redisClient, string(constants.BinanceMarkPriceAlertsQueue), eventHandler.HandleEvent, &worker.WorkerOptions{
		WorkerCount: 600,
	}, workerMetrics)

	if err := klinesSyncWorker.Start(ctx); err != nil {
		log.Fatal().Err(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Info().Msg("Received shutdown signal, gracefully stopping...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	eventHandler.Stop()
	klinesSyncWorker.Stop(10 * time.Second)

	if err := metricsServer.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("Error shutting down metrics server")
	}

	if err := redisClient.Close(); err != nil {
		log.Error().Err(err).Msg("Error closing Redis client")
	}

	cancel()
}
