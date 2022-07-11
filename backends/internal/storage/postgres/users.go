package postgres

import (
	"context"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rasulov-emirlan/todo-app/backends/internal/domain/users"
	"github.com/rasulov-emirlan/todo-app/backends/pkg/logging"
)

type usersRepository struct {
	conn *pgxpool.Pool
	log *logging.Logger
}

func (r *usersRepository) Create(ctx context.Context, email, hashedPassword, username string) (id string, err error) {
	sql, args, err := sq.Insert("users").Columns(
		"email", "password", "role_id", "username", "created_at", "updated_at").
		Values(email, hashedPassword, 2, username, time.Now(), time.Now()).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return "", err
	}

	defer r.log.Sync()
	r.log.Debug("usersRepository: Create()", logging.String("sql", sql))

	conn, err := r.conn.Acquire(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Release()

	err = conn.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		r.log.Debug("usersRepository: Create()", logging.String("error", err.Error()))
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return id, users.ErrEmailIsTaken
			}
		}
	}
	return id, err
}

func (r *usersRepository) Get(ctx context.Context, id string) (user users.User, err error) {
	sql, args, err := sq.Select(`id, email, password, role_id, username, created_at, updated_at`).
		From("users").Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return user, err
	}

	defer r.log.Sync()
	r.log.Debug("usersRepository: Get()", logging.String("sql", sql))

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
	sql, args, err := sq.Select(`id, email, password, role_id, username, created_at, updated_at`).
		From("users").Where(sq.Eq{"email": email}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return user, err
	}

	defer r.log.Sync()
	r.log.Debug("usersRepository: GetByEmail()", logging.String("sql", sql))

	conn, err := r.conn.Acquire(ctx)
	if err != nil {
		return user, err
	}
	defer conn.Release()

	err = conn.QueryRow(ctx, sql, args...).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role,
		&user.Username, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		r.log.Debug("usersRepository: GetByEmail()", logging.String("error", err.Error()))
		return user, users.ErrNoSuchUser
	}

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

	defer r.log.Sync()
	r.log.Debug("usersRepository: Update()", logging.String("sql", sql))

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

	defer r.log.Sync()
	r.log.Debug("usersRepository: Delete()", logging.String("sql", sql))


	conn, err := r.conn.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, sql, args...)
	return err
}

var roleIds = map[int]users.Role{
	1: users.RoleAdmin,
	2: users.RoleUser,
}
