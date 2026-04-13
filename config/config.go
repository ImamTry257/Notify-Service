package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	App   AppConfig
	MySQL MySQLConfig
	NATS  NATSConfig
	GRPC  GRPCConfig
	SMTP  SMTPConfig
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

	// Check APP_ENV for conditional binding
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "staging" || appEnv == "production" {
		// MySQL bindings
		viper.BindEnv("mysql.host", "MYSQL_HOST")
		viper.BindEnv("mysql.user", "MYSQL_USER")
		viper.BindEnv("mysql.password", "MYSQL_PASSWORD")
		viper.BindEnv("mysql.port", "MYSQL_PORT")
		viper.BindEnv("mysql.dbname", "MYSQL_DBNAME")

		// SMTP bindings
		viper.BindEnv("smtp.host", "SMTP_HOST")
		viper.BindEnv("smtp.port", "SMTP_PORT")
		viper.BindEnv("smtp.username", "SMTP_USERNAME")
		viper.BindEnv("smtp.password", "SMTP_PASSWORD")
		viper.BindEnv("smtp.sendername", "SMTP_SENDER_NAME")
		viper.BindEnv("smtp.senderemail", "SMTP_SENDER_EMAIL")

		// NATS bindings
		viper.BindEnv("nats.url", "NATS_URL")
		viper.BindEnv("nats.streamname", "NATS_STREAM_NAME")
		viper.BindEnv("nats.subject", "NATS_SUBJECT")
	}

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
		cfg.GRPC.Port = 50056
	}
	if cfg.NATS.StreamName == "" {
		cfg.NATS.StreamName = "notify"
	}
	if cfg.NATS.Subject == "" {
		cfg.NATS.Subject = "notify.email"
	}
}
