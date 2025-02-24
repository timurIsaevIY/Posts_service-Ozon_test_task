package config

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"time"
)

type Config struct {
	ConfigPath  string `env:"CONFIG_PATH" env-default:"config/config.yaml"`
	Database    Database
	HttpServer  HttpServer `yaml:"HttpServer"`
	StorageType string     `env:"STORAGE_TYPE" env-default:"postgres"`
}

type Database struct {
	DbHost string `env:"DB_HOST" env-required:"true"`
	DbPort int    `env:"DB_PORT" env-required:"true"`
	DbUser string `env:"DB_USER" env-required:"true"`
	DbPass string `env:"DB_PASS" env-required:"true"`
	DbName string `env:"DB_NAME" env-required:"true"`
}

type HttpServer struct {
	Address      string        `yaml:"Address" yaml-default:"8080"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" yaml-default:"60s"`
	ReadTimeout  time.Duration `yaml:"read_timeout" yaml-default:"10s"`
	WriteTimeout time.Duration `yaml:"write_timeout" yaml-default:"10s"`
}

func Load() (config Config, err error) {

	storageType := flag.String("storage_type", "", "Storage type (e.g., 'postgres', 'in-memory')")
	flag.Parse()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("failed to read config file: %v", err)
	}

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("failed to load .env file: %v", err)
	}

	// Bind environment variables to viper
	viper.AutomaticEnv()
	viper.BindEnv("Database.DbHost", "DB_HOST")
	viper.BindEnv("Database.DbPort", "DB_PORT")
	viper.BindEnv("Database.DbUser", "DB_USER")
	viper.BindEnv("Database.DbPass", "DB_PASS")
	viper.BindEnv("Database.DbName", "DB_NAME")
	viper.BindEnv("StorageType", "STORAGE_TYPE")

	// Unmarshal config into struct
	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config: %v", err)
	}
	if *storageType != "" {
		config.StorageType = *storageType
	}
	log.Printf("STORAGE_TYPE from env: %s", config.StorageType)
	return config, nil
}
