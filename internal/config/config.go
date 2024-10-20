package config

import (
	"flag"

	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/lib/logger"
)

type (
	Config struct {
		HTTP
		Log
		DB
		TLS
		Auth
	}

	HTTP struct {
		GRPCPort    string
		HTTPPort    string
		ReadTimeout int
	}

	Log struct {
		Level string
	}

	DB struct {
		URL           string
		RedisAddr     string
		RedisPassword string
		RedisDB       int
	}

	TLS struct {
		Cert string
		Key  string
	}

	Auth struct {
		JWTSecret       string
		AccessTokenTTL  int
		RefreshTokenTTL int
	}
)

func NewConfig() (*Config, error) {
	gRPCPort := flag.String("grpc_port", "localhost:50010", "GRPC Port")
	httpPort := flag.String("http_port", "localhost:8080", "HTTP Port")
	logLevel := flag.String("log_level", string(logger.InfoLevel), "logger level")
	dbURL := flag.String("db_url", "", "url for connection to database")
	readTimeout := flag.Int("read_timeout", 5, "read timeout")

	// TLS
	cert := flag.String("cert", "", "path to cert file")
	key := flag.String("key", "", "path to key file")

	// JWT
	jwtSecret := flag.String("jwt_secret", "", "jwt secret")
	accessTokenTTL := flag.Int("access_token_ttl", 2, "access token ttl")
	refreshTokenTTL := flag.Int("refresh_token_ttl", 14400, "refresh token ttl")

	// Redis
	redisAddr := flag.String("redis_addr", "localhost:6379", "redis address")
	redisPassword := flag.String("redis_password", "", "redis password")
	redisDB := flag.Int("redis_db", 0, "redis db")

	flag.Parse()

	cfg := &Config{
		HTTP: HTTP{
			GRPCPort:    *gRPCPort,
			HTTPPort:    *httpPort,
			ReadTimeout: *readTimeout,
		},
		Log: Log{
			Level: *logLevel,
		},
		DB: DB{
			URL:           *dbURL,
			RedisAddr:     *redisAddr,
			RedisPassword: *redisPassword,
			RedisDB:       *redisDB,
		},
		TLS: TLS{
			Cert: *cert,
			Key:  *key,
		},
		Auth: Auth{
			JWTSecret:       *jwtSecret,
			AccessTokenTTL:  *accessTokenTTL,
			RefreshTokenTTL: *refreshTokenTTL,
		},
	}

	return cfg, nil
}
