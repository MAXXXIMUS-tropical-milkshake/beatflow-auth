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
		Port string
	}

	Log struct {
		Level string
	}

	DB struct {
		URL string
	}

	TLS struct {
		Cert string
		Key  string
	}

	Auth struct {
		JWTSecret string
		TokenTTL  int
	}
)

func NewConfig() (*Config, error) {
	port := flag.String("port", "8080", "HTTP port")
	logLevel := flag.String("log_level", string(logger.InfoLevel), "logger level")
	dbURL := flag.String("db_url", "", "url for connection to database")

	// TLS
	cert := flag.String("cert", "", "path to cert file")
	key := flag.String("key", "", "path to key file")

	// JWT
	jwtSecret := flag.String("jwt_secret", "", "jwt secret")
	tokenTTL := flag.Int("token_ttl", 10, "token ttl")

	flag.Parse()

	cfg := &Config{
		HTTP: HTTP{
			Port: *port,
		},
		Log: Log{
			Level: *logLevel,
		},
		DB: DB{
			URL: *dbURL,
		},
		TLS: TLS{
			Cert: *cert,
			Key:  *key,
		},
		Auth: Auth{
			JWTSecret: *jwtSecret,
			TokenTTL:  *tokenTTL,
		},
	}

	return cfg, nil
}
