package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Storage Storage `yaml:"storage" env-required:"true"`
	Kafka   Kafka   `yaml:"kafka" env-required:"true"`
}

type Storage struct {
	Host         string `yaml:"host" env:"POSTGRES_HOST" env-required:"true"`
	Port         string `yaml:"port" env:"POSTGRES_PORT" env-required:"true"`
	DatabaseName string `yaml:"database_name" env:"POSTGRES_DATABASE_NAME" env-required:"true"`
	User         string `yaml:"user" env:"POSTGRES_USER" env-required:"true"`
	Password     string `yaml:"password" env:"POSTGRES_PASSWORD" env-required:"true"`
}

type Kafka struct {
	BrokerAddress string `yaml:"broker_address" env:"KAFKA_BROKER_ADDRESS" env-required:"true"`
}

func NewConfig(filePath string) (*Config, error) {
	config := &Config{}

	if err := cleanenv.ReadConfig(filePath, config); err != nil {
		return nil, err
	}

	return config, nil
}

func (cfg *Config) GetPostgresConnectionString() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Storage.User,
		cfg.Storage.Password,
		cfg.Storage.Host,
		cfg.Storage.Port,
		cfg.Storage.DatabaseName,
	)
}
