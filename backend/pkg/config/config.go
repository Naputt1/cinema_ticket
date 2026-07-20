package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Mongo    MongoConfig
	Redis    RedisConfig
	RabbitMQ RabbitMQConfig
	Google   GoogleOAuthConfig
	JWT      JWTConfig
	App      AppConfig
}

type ServerConfig struct {
	Port string
}

type MongoConfig struct {
	URI string
	DB  string
}

type RedisConfig struct {
	URL string
}

type RabbitMQConfig struct {
	URL string
}

type GoogleOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type JWTConfig struct {
	Secret        string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

type AppConfig struct {
	FrontendURL string
	AdminEmails []string
	Production  bool
	Env         string
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.SetDefault("SERVER_PORT", "8080")
	v.SetDefault("MONGO_URI", "mongodb://localhost:27017")
	v.SetDefault("MONGO_DB", "cinema_ticket")
	v.SetDefault("REDIS_URL", "redis://localhost:6379")
	v.SetDefault("RABBITMQ_URL", "amqp://guest:guest@localhost:5672")
	v.SetDefault("JWT_SECRET", "dev-secret-change-in-production")
	v.SetDefault("JWT_ACCESS_EXPIRY", "15m")
	v.SetDefault("JWT_REFRESH_EXPIRY", "168h")
	v.SetDefault("FRONTEND_URL", "http://localhost:5173")
	v.SetDefault("APP_ENV", "development")
	v.SetDefault("ADMIN_EMAILS", "")
	v.SetDefault("PRODUCTION", false)

	v.ReadInConfig()

	accessExpiry, _ := time.ParseDuration(v.GetString("JWT_ACCESS_EXPIRY"))
	refreshExpiry, _ := time.ParseDuration(v.GetString("JWT_REFRESH_EXPIRY"))

	adminEmails := []string{}
	if ae := v.GetString("ADMIN_EMAILS"); ae != "" {
		adminEmails = strings.Split(ae, ",")
	}

	return &Config{
		Server: ServerConfig{
			Port: v.GetString("SERVER_PORT"),
		},
		Mongo: MongoConfig{
			URI: v.GetString("MONGO_URI"),
			DB:  v.GetString("MONGO_DB"),
		},
		Redis: RedisConfig{
			URL: v.GetString("REDIS_URL"),
		},
		RabbitMQ: RabbitMQConfig{
			URL: v.GetString("RABBITMQ_URL"),
		},
		Google: GoogleOAuthConfig{
			ClientID:     v.GetString("GOOGLE_CLIENT_ID"),
			ClientSecret: v.GetString("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  v.GetString("GOOGLE_REDIRECT_URL"),
		},
		JWT: JWTConfig{
			Secret:        v.GetString("JWT_SECRET"),
			AccessExpiry:  accessExpiry,
			RefreshExpiry: refreshExpiry,
		},
		App: AppConfig{
			FrontendURL: v.GetString("FRONTEND_URL"),
			AdminEmails: adminEmails,
			Production:  v.GetBool("PRODUCTION"),
			Env:         v.GetString("APP_ENV"),
		},
	}, nil
}
