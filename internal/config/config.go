package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Env                    string `env:"ENV" env-default:"development"`
	BinanceMarkPricesRedis string `env:"BINANCE_MARK_PRICES_REDIS_URL"`
	PostgresqlUsername     string `env:"POSTGRESQL_USERNAME" env-default:"postgres"`
	PostgresqlPassword     string `env:"POSTGRESQL_PASSWORD" env-default:"postgres"`
	PostgresqlDatabaseName string `env:"POSTGRESQL_DATABASE_NAME" env-default:"postgres"`
	PostgresqlHost         string `env:"POSTGRESQL_HOST" env-default:"localhost"`
	PostgresqlPort         string `env:"POSTGRESQL_PORT" env-default:"5432"`
	LogLevel               string `env:"LOG_LEVEL"`
	ServerHost             string `env:"SERVER_HOST" env-default:"0.0.0.0"`
	ServerPort             string `env:"SERVER_PORT" env-default:"3000"`
	ServiceName            string `env:"SERVICE_NAME"`
	HTTPTimeout            int32  `env:"HTTP_TIMEOUT" env-default:"175"`
}

func (c *Config) HTTPTimeoutDuration() time.Duration {
	return time.Duration(c.HTTPTimeout) * time.Second
}

type goEnv struct {
	GoMod string `json:"GOMOD"`
}

func LoadConfig() (*Config, error) {
	c := new(Config)

	err := loadConfig(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func loadConfig(cfg interface{}) error {
	envFilePath, fileExists := initLocal()

	var readCfgErr error

	if fileExists {
		readCfgErr = cleanenv.ReadConfig(envFilePath, cfg)
	} else {
		readCfgErr = cleanenv.ReadEnv(cfg)
	}

	if readCfgErr != nil {
		return readCfgErr
	}

	updateEnvErr := cleanenv.UpdateEnv(cfg)
	if updateEnvErr != nil {
		return updateEnvErr
	}

	return nil
}

func initLocal() (string, bool) {
	if !goExists() {
		return "", false
	}

	modRoot := getModuleRoot()
	localEnvFilePath := fmt.Sprintf("%s/.local.env", modRoot)
	envFilePath := fmt.Sprintf("%s/.env", modRoot)

	if envFileExists(localEnvFilePath) {
		return localEnvFilePath, true
	} else if envFileExists(envFilePath) {
		return envFilePath, true
	}

	return "", false
}

func envFileExists(envFilePath string) bool {
	_, err := os.Stat(envFilePath)

	return err == nil
}

func goExists() bool {
	_, err := exec.LookPath("go")

	return err == nil
}

func getModuleRoot() string {
	goEnvRaw, err := exec.Command("go", "env", "-json").Output()
	if err != nil {
		log.Fatal().Err(err).Msg("go env command failed")

		return ""
	}

	env := new(goEnv)

	err = json.Unmarshal(goEnvRaw, env)
	if err != nil {
		log.Fatal().Err(err).Msg("go mod unmarshalling failed")

		return ""
	}

	return strings.TrimSuffix(env.GoMod, "/go.mod")
}
