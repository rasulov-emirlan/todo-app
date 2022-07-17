package postgres

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rasulov-emirlan/todo-app/backends/config"
	"github.com/rasulov-emirlan/todo-app/backends/internal/storage/postgres/migrations"
	"github.com/rasulov-emirlan/todo-app/backends/pkg/logging"
)

type Repository struct {
	conn *pgxpool.Pool

	usersRepository *usersRepository
	todosRepository *todosRepository
}

func NewRepository(cfg config.Config, logger *logging.Logger) (*Repository, error) {
	url := cfg.Database.URL()
	conn, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		for i := 0; i < 60; i++ {
			time.Sleep(time.Second)
			log.Println("Trying to connect to database at url:", url)
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
	if cfg.Database.WithMigrations {
		if err := migrations.Up(url); err != nil {
			return nil, err
		}
	}
	return &Repository{
		conn:            conn,
		usersRepository: &usersRepository{conn: conn, log: logger},
		todosRepository: &todosRepository{conn: conn, log: logger},
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
