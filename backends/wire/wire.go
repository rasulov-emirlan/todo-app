//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/rasulov-emirlan/todo-app/backends/config"
	"github.com/rasulov-emirlan/todo-app/backends/internal/domain/todos"
	"github.com/rasulov-emirlan/todo-app/backends/internal/domain/users"
	"github.com/rasulov-emirlan/todo-app/backends/internal/storage/postgres"
	"github.com/rasulov-emirlan/todo-app/backends/internal/transport/resthttp"
	"github.com/rasulov-emirlan/todo-app/backends/pkg/logging"
	"github.com/rasulov-emirlan/todo-app/backends/pkg/validation"
)

func InitializeRepo(config config.Config, logger *logging.Logger) (*postgres.Repository, error) {
	wire.Build(postgres.NewRepository)
	return &postgres.Repository{}, nil
}

func InitializeLogger(config config.Config) (*logging.Logger, error) {
	wire.Build(logging.NewLogger)
	return &logging.Logger{}, nil
}

func InitializeValidator() (*validation.Validator, error) {
	wire.Build(validation.NewValidator)
	return &validation.Validator{}, nil
}

func InitializeRestApi(
	config config.Config,
	logger *logging.Logger,
	validator *validation.Validator,
	repository *postgres.Repository,
) (*resthttp.Server, error) {
	uS, err := users.NewService(repository.Users(), logger, validator, []byte(config.JWTsecret))
	if err != nil {
		return nil, err
	}
	tS := todos.NewService(repository.Todos(), repository.Users(), logger, validator)
	if err != nil {
		return nil, err
	}
	return resthttp.NewServer(config, logger, validator, uS, tS), nil
}
