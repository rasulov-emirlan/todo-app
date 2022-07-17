package config

import (
	"net/url"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		JWTsecret    string        `env:"JWT_SECRET" env-default:"secret"`
		Port         string        `env:"PORT" env-default:":8080"`
		WriteTimeout time.Duration `env:"WRITE_TIMEOUT" env-default:"15s"`
		ReadTimeout  time.Duration `env:"READ_TIMEOUT" env-default:"15s"`
		CORSorigins  string        `env:"CORS_ORIGINS" env-default:"*"`
		Log          struct {
			Level  string `env:"LOG_LEVEL" env-default:"debug"`
			Output string `env:"LOG_OUTPUT" env-default:"stdout"`
		}
		Database database
	}
	database struct {
		Host           string `env:"POSTGRES_HOST" env-default:"localhost"`
		Port           string `env:"POSTGRES_PORT" env-default:"5432"`
		User           string `env:"POSTGRES_USER" env-default:"postgres"`
		Pass           string `env:"POSTGRES_PASSWORD" env-default:"postgres"`
		Name           string `env:"POSTGRES_DB" env-default:"todos"`
		WithMigrations bool   `env:"WITH_MIGRATIONS" env-default:"false"`
	}
)

func LoadConfigs(filename string) (*Config, error) {
	var config Config
	switch len(filename) {
	case 0:
		err := cleanenv.ReadEnv(&config)
		if err != nil {
			return nil, err
		}
	default:
		err := cleanenv.ReadConfig(filename, &config)
		if err != nil {
			return nil, err
		}
	}
	return &config, nil
}

func (d database) URL() string {
	url := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(d.User, d.Pass),
		Host:   d.Host + ":" + d.Port,
		Path:   d.Name,
	}
	return url.String() + "?sslmode=disable"
}
