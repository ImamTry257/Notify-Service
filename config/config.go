package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	App         AppConfig
	MySQL       MySQLConfig
	NATS        NATSConfig
	GRPC        GRPCConfig
	SMTP        SMTPConfig
}

type AppConfig struct {
	Name string
	Env  string
}

type SMTPConfig struct {
	Host        string
	Port        int
	Username    string
	Password    string
	SenderName  string
	SenderEmail string
}

type MySQLConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	DBName   string
}

type NATSConfig struct {
	URL        string
	StreamName string
	Subject    string
}

type GRPCConfig struct {
	Port int
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found, using environment variables")
		} else {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// Default values if not set
	setDefaults(&config)

	return &config, nil
}

func setDefaults(cfg *Config) {
	if cfg.GRPC.Port == 0 {
		cfg.GRPC.Port = 50051
	}
	if cfg.NATS.StreamName == "" {
		cfg.NATS.StreamName = "notify"
	}
	if cfg.NATS.Subject == "" {
		cfg.NATS.Subject = "notify.email"
	}
}
