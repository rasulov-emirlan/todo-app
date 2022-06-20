package postgres

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rasulov-emirlan/todo-app/backends/internal/domain/users"
)

type usersRepository struct {
	conn *pgxpool.Pool
}

func (r *usersRepository) Create(ctx context.Context, email, hashedPassword, username string) (id string, err error) {
	sql, args, err := sq.Insert("users").Columns("id",
		"email", "password", "role", "username", "created_at", "updated_at").
		Values(id, email, hashedPassword, username, time.Now(), time.Now()).
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

func (r *usersRepository) Get(ctx context.Context, id string) (user users.User, err error) {
	sql, args, err := sq.Select(`id, email, password, role, username, created_at, updated_at`).
		From("users").Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return user, err
	}

	conn, err := r.conn.Acquire(ctx)
	if err != nil {
		return user, err
	}
	defer conn.Release()

	err = conn.QueryRow(ctx, sql, args...).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role,
		&user.Username, &user.CreatedAt, &user.UpdatedAt,
	)

	return user, err
}

func (r *usersRepository) GetByEmail(ctx context.Context, email string) (user users.User, err error) {
	sql, args, err := sq.Select(`id, email, password, role, username, created_at, updated_at`).
		From("users").Where(sq.Eq{"email": email}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return user, err
	}

	conn, err := r.conn.Acquire(ctx)
	if err != nil {
		return user, err
	}
	defer conn.Release()

	err = conn.QueryRow(ctx, sql, args...).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role,
		&user.Username, &user.CreatedAt, &user.UpdatedAt,
	)

	return user, err
}

func (r *usersRepository) Update(ctx context.Context, inp users.UpdateInput) (err error) {
	sql, args, err := sq.Update("users").
		Set("password", inp.Password).
		Set("username", inp.Username).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": inp.ID}).
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

func (r *usersRepository) Delete(ctx context.Context, id string) (err error) {
	sql, args, err := sq.Delete("users").Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar).ToSql()
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
