package config

import (
	"fmt"
	"os"
)

type AppConfig struct {
	Database   DatabaseConfig
	Keycloak   KeycloakConfig
	Kafka      KafkaConfig
	Jaeger     JaegerConfig
	Server     ServerConfig
	Telegram   TelegramConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func (d *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		d.User, d.Password, d.Host, d.Port, d.DBName)
}

type KeycloakConfig struct {
	Host         string
	Port         string
	Client       string
	Realm        string
	ClientSecret string
	Admin        string
	AdminPassword string
	MasterRealm  string
}

func (k *KeycloakConfig) BaseURL() string {
	return fmt.Sprintf("http://%s:%s", k.Host, k.Port)
}

type KafkaConfig struct {
	Broker string
	Port   string
	Topic  string
}

func (k *KafkaConfig) BrokerAddress() string {
	return fmt.Sprintf("%s:%s", k.Broker, k.Port)
}

type JaegerConfig struct {
	AgentPort string
	SendPort  string
	Host      string
}

type ServerConfig struct {
	Host string
	Port string
}

func (s *ServerConfig) Address() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}

type TelegramConfig struct {
	APIKey string
}

func LoadFromEnv() (*AppConfig, error) {
	cfg := &AppConfig{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "db"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", "postgres"),
		},
		Keycloak: KeycloakConfig{
			Host:          getEnv("KEYCLOAK_HOST", "keycloak"),
			Port:          getEnv("KEYCLOAK_INNER_PORT", "8080"),
			Client:        getEnv("KEYCLOAK_CLIENT", ""),
			Realm:         getEnv("KEYCLOAK_REALM", ""),
			ClientSecret:  getEnv("KEYCLOAK_CLIENT_SECRET", ""),
			Admin:         getEnv("KEYCLOAK_ADMIN", ""),
			AdminPassword: getEnv("KEYCLOAK_ADMIN_PASSWORD", ""),
			MasterRealm:   getEnv("KEYCLOAK_MASTER_REALM", "master"),
		},
		Kafka: KafkaConfig{
			Broker: getEnv("KAFKA_BROKER", "kafka"),
			Port:   getEnv("KAFKA_PORT", "9092"),
			Topic:  getEnv("KAFKA_TOPIC", "notifications"),
		},
		Jaeger: JaegerConfig{
			Host:      getEnv("JAEGER_HOST", "jaeger"),
			AgentPort: getEnv("JAEGER_AGENT_PORT", "6831"),
			SendPort:  getEnv("JAEGER_SEND_PORT", "14268"),
		},
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Telegram: TelegramConfig{
			APIKey: getEnv("TELEGRAM_API_KEY", ""),
		},
	}

	if cfg.Keycloak.Client == "" {
		return nil, fmt.Errorf("KEYCLOAK_CLIENT is required")
	}
	if cfg.Keycloak.Realm == "" {
		return nil, fmt.Errorf("KEYCLOAK_REALM is required")
	}
	if cfg.Keycloak.ClientSecret == "" {
		return nil, fmt.Errorf("KEYCLOAK_CLIENT_SECRET is required")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

