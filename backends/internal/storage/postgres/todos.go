package postgres

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lib/pq"
	"github.com/rasulov-emirlan/todo-app/backends/internal/domain/todos"
	"github.com/rasulov-emirlan/todo-app/backends/internal/domain/users"
)

type todosRepository struct {
	conn *pgxpool.Pool
}

func (r *todosRepository) Create(ctx context.Context, inp todos.CreateInput) (id string, err error) {
	sql, args, err := sq.Insert("todos").Columns(
		"user_id", "title", "description", "deadline", "created_at", "updated_at").
		Values(inp.UserID, inp.Title, inp.Body, inp.Deadline, time.Now(), time.Now()).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return "", err
	}

	conn, err := r.conn.Acquire(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Release()

	err = conn.QueryRow(ctx, sql, args...).Scan(&id)
	return id, err
}

func (r *todosRepository) Get(ctx context.Context, id string) (todo todos.Todo, err error) {
	sql, args, err := sq.Select(`id, user_id, title, description, "deadline", created_at, updated_at`).
		From("todos").Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return todo, err
	}

	conn, err := r.conn.Acquire(ctx)
	if err != nil {
		return todo, err
	}
	defer conn.Release()

	var (
		authorId  string
		deadline  pq.NullTime
		updatedAt pq.NullTime
	)

	err = conn.QueryRow(ctx, sql, args...).Scan(&todo.ID, &authorId, &todo.Title, &todo.Body, &deadline, &todo.CreatedAt, &updatedAt)
	if err != nil {
		return todo, err
	}

	if updatedAt.Valid {
		todo.UpdatedAt = updatedAt.Time
	}
	if deadline.Valid {
		todo.Deadline = deadline.Time
	}
	todo.Author = &users.User{ID: authorId}
	return todo, nil
}

func (r *todosRepository) GetAll(ctx context.Context, config todos.GetAllInput) ([]todos.Todo, error) {
	query := sq.
		Select(`id, user_id, title, description, deadline, created_at, updated_at`).
		From("todos").
		Limit(uint64(config.PageSize)).
		Offset(uint64(config.PageSize * config.Page))

	switch config.SortBy {
	case todos.SortByCreationASC:
		query.OrderBy("created_at ASC")
	case todos.SortByCreationDESC:
		query.OrderBy("created_at DESC")
	case todos.SortByDeadlineASC:
		query.OrderBy("deadline ASC")
	case todos.SortByDeadlineDESC:
		query.OrderBy("deadline DESC")
	}

	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	conn, err := r.conn.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		authorId  string
		deadline  pq.NullTime
		updatedAt pq.NullTime
		todolist  = []todos.Todo{}
		todo      = todos.Todo{}
	)

	for rows.Next() {
		err = rows.Scan(
			&todo.ID,
			&authorId,
			&todo.Title,
			&todo.Body,
			&deadline,
			&todo.CreatedAt,
			&updatedAt)
		if err != nil {
			return nil, err
		}
		if updatedAt.Valid {
			todo.UpdatedAt = updatedAt.Time
		}
		if deadline.Valid {
			todo.Deadline = deadline.Time
		}
		todo.Author = &users.User{ID: authorId}
		todolist = append(todolist, todo)
	}

	return todolist, nil
}

func (r *todosRepository) Update(ctx context.Context, inp todos.UpdateInput) error {
	sql, args, err := sq.
		Update("todos").
		Set("title", inp.Title).
		Set("description", inp.Body).
		Set("deadline", inp.Deadline).
		Where(sq.Eq{"id": inp.ID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	conn, err := r.conn.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, sql, args...)
	return err
}

func (r *todosRepository) MarkAsComplete(ctx context.Context, id string) error {
	panic("not implemented")
}
func (r *todosRepository) MarkAsNotComplete(ctx context.Context, id string) error {
	panic("not implemented")
}
func (r *todosRepository) Delete(ctx context.Context, id string) error {
	panic("not implemented")
}
