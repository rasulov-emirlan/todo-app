package postgres

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rasulov-emirlan/todo-app/backends/internal/storage/postgres/migrations"
)

type repository struct {
	conn *pgxpool.Pool

	usersRepository *usersRepository
	todosRepository *todosRepository
}

func NewRepository(url string, withMigrations bool) (*repository, error) {
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
		log.Println("Dfsfdfsdf")
		if err := migrations.Up(url); err != nil {
			return nil, err
		}
	}
	return &repository{
		conn:            conn,
		usersRepository: &usersRepository{conn},
		todosRepository: &todosRepository{conn},
	}, nil
}

func (r *repository) Close() error {
	r.conn.Close()
	return nil
}

func (r *repository) Users() *usersRepository {
	return r.usersRepository
}

func (r *repository) Todos() *todosRepository {
	return r.todosRepository
}
