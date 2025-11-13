package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type dbConfig struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-required:"true"`
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	DbName   string `yaml:"db_name" env-required:"true"`
}

type apiConfig struct {
	Host string `yaml:"host" env-required:"true"`
	Port int    `yaml:"port" env-required:"true"`
}

type Config struct {
	db       dbConfig  `yaml:"db" env-required:"true"`
	api      apiConfig `yaml:"api" env-required:"true"`
	LogLevel string    `env:"LOGLEVEL"`
}

func Load(cfgPath string) (*Config, error) {
	var cfg Config
	err := cleanenv.ReadConfig(cfgPath, &cfg)
	if err != nil {
		return nil, err
	}

	logLevel := os.Getenv("LOGLEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	cfg.LogLevel = logLevel
	return &cfg, nil
}

func (cfg *Config) ConnString() string {
	return fmt.Sprintf(`user=%s password=%s host=%s port=%d dbname=%s sslmode=disable`,
		cfg.db.User,
		cfg.db.Password,
		cfg.db.Host,
		cfg.db.Port,
		cfg.db.DbName)
}

func (cfg *Config) ListenAddr() string {
	return fmt.Sprintf("%s:%d",
		cfg.api.Host,
		cfg.api.Port)
}
