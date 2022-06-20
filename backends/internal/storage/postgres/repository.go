package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rasulov-emirlan/todo-app/backends/internal/domain/users"
	"github.com/rasulov-emirlan/todo-app/backends/internal/storage/postgres/migrations"
)

type repository struct {
	conn *pgxpool.Pool

	usersRepository users.Repository
}

func NewRepository(url string, withMigrations bool) (*repository, error) {
	conn, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}
	if err := migrations.Up(url); err != nil {
		return nil, err
	}
	return &repository{
		conn:            conn,
		usersRepository: &usersRepository{conn},
	}, nil
}

func (r *repository) Close() error {
	r.conn.Close()
	return nil
}

func (r *repository) Users() users.Repository {
	return r.usersRepository
}
