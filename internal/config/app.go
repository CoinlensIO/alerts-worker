package config

import (
	"alerts-worker/pkg/trace_metrics"

	"github.com/rs/zerolog/log"
	"github.com/samber/do"
)

type AppBase struct {
	Config       *Config
	Injector     *do.Injector
	TraceMetrics *trace_metrics.TraceMetrics
}

func New(options ...func(*AppBase)) *AppBase {
	appBase := &AppBase{}

	for _, o := range options {
		o(appBase)
	}

	return appBase
}

func Init() func(*AppBase) {
	return func(appBase *AppBase) {
		cfg, err := LoadConfig()
		if err != nil {
			panic(err)
		}

		appBase.Config = cfg
	}
}

func WithDependencyInjector() func(*AppBase) {
	return func(appBase *AppBase) {
		appBase.Injector = NewInjector(appBase.Config)
	}
}

func (appBase *AppBase) Shutdown() {
	err := appBase.Injector.Shutdown()
	if err != nil {
		log.Panic().Err(err).Msg("injector's shutdown failed")
	}
}
