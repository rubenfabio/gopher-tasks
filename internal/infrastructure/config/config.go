package config

import (
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/rubenfabio/gopher-tasks/internal/infrastructure/logger"
	"github.com/spf13/viper"
)

// Config agrupa todas as configurações da aplicação.
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Auth     AuthConfig
    Log      LogConfig
}

type ServerConfig struct {
    Port         int           `mapstructure:"port"`
    ReadTimeout  time.Duration `mapstructure:"readtimeout"`
    WriteTimeout time.Duration `mapstructure:"writetimeout"`
}

type DatabaseConfig struct {
    Driver string `mapstructure:"driver"`
    DSN    string `mapstructure:"dsn"`
}

type AuthConfig struct {
    JWTSecret          string `mapstructure:"jwtsecret"`
    TokenExpiryMinutes int    `mapstructure:"tokenexpiryminutes"`
}

type LogConfig struct {
    Level  string `mapstructure:"level"`
    Format string `mapstructure:"format"`
}

// Load carrega .env.local, config YAML e ENVs via Viper e registra logs.
func Load(path string) (*Config, error) {
    // logger temporário
    log := logger.NewDefault()

    // 1) Tenta carregar .env.local (variáveis expostas ao os.Getenv)
    if err := godotenv.Load(".env.local"); err != nil {
        log.WithField("file", ".env.local").Warn("No .env.local file found, skipping")
    } else {
        log.Info("Loaded .env.local")
    }

    // 2) Config file (YAML) e ENV vars
    log.Infof("Loading configuration from %s", path)
    v := viper.New()
    v.SetConfigFile(path)
    v.AutomaticEnv()                                      // lê variáveis de ambiente
    v.SetEnvPrefix("APP")                                 // APP_SERVER_PORT, APP_DATABASE_DSN…
    v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))    // server.port → APP_SERVER_PORT
    v.AllowEmptyEnv(true)

    if err := v.ReadInConfig(); err != nil {
        log.WithField("error", err).Errorf("Failed to read config file: %s", path)
        return nil, err
    }

    // 3) Defaults
    v.SetDefault("server.readtimeout", 5*time.Second)
    v.SetDefault("server.writetimeout", 10*time.Second)

    // 4) Unmarshal em struct
    var cfg Config
    if err := v.Unmarshal(&cfg); err != nil {
        log.WithField("error", err).Error("Failed to parse configuration")
        return nil, err
    }

    log.WithField("config", cfg).Info("Configuration loaded successfully")
    return &cfg, nil
}
