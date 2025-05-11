package config

import (
	"time"

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

// Load carrega config de YAML/ENV via Viper e registra logs do processo.
func Load(path string) (*Config, error) {
    // Logger temporário para os logs de configuração
    log := logger.NewDefault()
    log.Infof("Loading configuration from %s", path)

    v := viper.New()
    v.SetConfigFile(path)
    v.AutomaticEnv()           // lê variáveis de ambiente
    v.SetEnvPrefix("APP")      // prefixo APP_SERVER_PORT, etc.
    v.AllowEmptyEnv(true)

    // Leitura do arquivo de configuração
    if err := v.ReadInConfig(); err != nil {
        log.WithField("error", err).Errorf("Failed to read config file: %s", path)
        return nil, err
    }

    // Valores padrão
    v.SetDefault("server.readtimeout", 5*time.Second)
    v.SetDefault("server.writetimeout", 10*time.Second)

    // Unmarshal para struct
    var cfg Config
    if err := v.Unmarshal(&cfg); err != nil {
        log.WithField("error", err).Error("Failed to parse configuration")
        return nil, err
    }

    log.WithField("config", cfg).Info("Configuration loaded successfully")
    return &cfg, nil
}
