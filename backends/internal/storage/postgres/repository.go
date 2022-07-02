package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rasulov-emirlan/todo-app/backends/internal/storage/postgres/migrations"
)

type Repository struct {
	conn *pgxpool.Pool

	usersRepository *usersRepository
	todosRepository *todosRepository
}

func NewRepository(url string, withMigrations bool) (*Repository, error) {
	conn, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		for i := 0; i < 60; i++ {
			time.Sleep(time.Second)
			conn, err = pgxpool.Connect(context.Background(), url)
			if err == nil {
				break
			}
		}
		if err != nil {
			return nil, err
		}
	}
	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}
	if withMigrations {
		if err := migrations.Up(url); err != nil {
			return nil, err
		}
	}
	return &Repository{
		conn:            conn,
		usersRepository: &usersRepository{conn},
		todosRepository: &todosRepository{conn},
	}, nil
}

func (r *Repository) Close() error {
	r.conn.Close()
	return nil
}

func (r *Repository) Users() *usersRepository {
	return r.usersRepository
}

func (r *Repository) Todos() *todosRepository {
	return r.todosRepository
}

func (r *Repository) Ping() error {
	return r.conn.Ping(context.Background())
}
