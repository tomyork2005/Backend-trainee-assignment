package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env            string `yaml:"env" env-default:"local"`
	HttpConfig     `yaml:"http_server"`
	PostgresConfig `yaml:"postgres"`
}

type HttpConfig struct {
	Address      string        `yaml:"address" address-default:"localhost:8080"`
	WriteTimeout time.Duration `yaml:"write_timeout" write-timeout-default:"4s"`
	ReadTimeout  time.Duration `yaml:"read_timeout" read-timeout-default:"4s"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" idle-timeout-default:"60s"`
}

type PostgresConfig struct {
	Host     string `yaml:"host" host-default:"localhost"`
	Port     int    `yaml:"port" port-default:"5432"`
	User     string `yaml:"user" user-default:"postgres"`
	Password string `yaml:"password" password-default:"postgres"`
}

func MustLoadConfig() *Config {
	configEnv := os.Getenv("CONFIG_PATH")
	if configEnv == "" {
		log.Fatal("configs path can`t be empty")
	}

	if _, err := os.Stat(configEnv); os.IsNotExist(err) {
		log.Fatalf("configs file doesn't exist %s", err)
	}

	var config Config
	if err := cleanenv.ReadConfig(configEnv, &config); err != nil {
		log.Fatalf("fail with read configs %s", err)
	}

	return &config
}
