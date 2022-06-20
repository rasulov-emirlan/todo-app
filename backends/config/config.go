package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	JWTsecret string `env:"JWT_SECRET" env-default:"secret"`
	Port      string `env:"PORT" env-default:":8080"`
	Log       struct {
		Level  string `env:"LOG_LEVEL" env-default:"info"`
		Output string `env:"LOG_OUTPUT" env-default:"stdout"`
	}
	Database struct {
		Host string `env:"POSTGRES_HOST" env-default:"localhost"`
		Port string `env:"POSTGRES_PORT" env-default:"5432"`
		User string `env:"POSTGRES_USERNAME" env-default:"postgres"`
		Pass string `env:"POSTGRES_PASSWORD" env-default:"postgres"`
		Name string `env:"POSTGRES_DB" env-default:"todos"`
	}
}

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
