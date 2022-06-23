package postgres

import (
	"context"
	"log"
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
	sql, args, err := sq.
		Insert("todos").
		Columns("user_id, title, description, deadline, created_at, updated_at").
		Values(inp.UserID, inp.Title, inp.Body, inp.Deadline, time.Now(), nil).
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
	sql, args, err := sq.
		Select(
			`t.id, user_id, username, 
			email, role_id, u.created_at,
			title, description, deadline, 
			t.created_at, t.updated_at`,
		).
		From("todos AS t").Where(sq.Eq{"t.id": id}).
		InnerJoin("users AS u ON t.user_id = u.id").
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return todo, err
	}

	log.Println(sql)

	conn, err := r.conn.Acquire(ctx)
	if err != nil {
		return todo, err
	}
	defer conn.Release()

	var (
		author    users.User
		roleID    int
		deadline  pq.NullTime
		updatedAt pq.NullTime
	)

	err = conn.QueryRow(ctx, sql, args...).Scan(
		&todo.ID, &author.ID, &author.Username,
		&author.Email, &roleID, &author.CreatedAt,
		&todo.Title, &todo.Body, &deadline,
		&todo.CreatedAt, &updatedAt,
	)
	if err != nil {
		return todo, err
	}

	if updatedAt.Valid {
		todo.UpdatedAt = updatedAt.Time
	}
	if deadline.Valid {
		todo.Deadline = deadline.Time
	}
	author.Role = roleIds[roleID]
	todo.Author = &author
	return todo, nil
}

var sortingVariants = map[todos.SortBy]string{
	todos.SortByCreationASC:  "created_at ASC",
	todos.SortByCreationDESC: "created_at DESC",
	todos.SortByDeadlineASC:  "deadline ASC",
	todos.SortByDeadlineDESC: "deadline DESC",
}

func (r *todosRepository) GetAll(ctx context.Context, config todos.GetAllInput) ([]todos.Todo, error) {
	query := sq.
		Select(`id, user_id, title, description, deadline, created_at, updated_at`).
		From("todos").
		Limit(uint64(config.PageSize)).
		Offset(uint64(config.PageSize * config.Page))

	if len(config.UserID) != 0 {
		query.Where(sq.Eq{"user_id": config.UserID})
	}

	if sorting, ok := sortingVariants[config.SortBy]; ok {
		query.OrderBy(sorting)
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
	query := sq.
		Update("todos").
		Where(sq.Eq{"id": inp.ID})

	if len(inp.Title) != 0 {
		query.Set("title", inp.Title)
	}
	if len(inp.Body) != 0 {
		query.Set("body", inp.Body)
	}
	if !inp.Deadline.IsZero() {
		query.Set("deadline", inp.Deadline)
	}

	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
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
	sql, args, err := sq.
		Update("todos").
		Set("completed", true).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).ToSql()
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

func (r *todosRepository) MarkAsNotComplete(ctx context.Context, id string) error {
	sql, args, err := sq.
		Update("todos").
		Set("completed", false).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).ToSql()
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

func (r *todosRepository) Delete(ctx context.Context, id string) error {
	sql, args, err := sq.
		Delete("todos").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).ToSql()
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
