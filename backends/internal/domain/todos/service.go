package todos

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type (
	Repository interface {
		Create(ctx context.Context, inp CreateInput) (id string, err error)
		Get(ctx context.Context, id string) (todo Todo, err error)
		GetAll(ctx context.Context, config GetAllInput) (todos []Todo, err error)
		Update(ctx context.Context, inp UpdateInput) error
		MarkAsComplete(ctx context.Context, id string) error
		MarkAsNotComplete(ctx context.Context, id string) error
		Delete(ctx context.Context, id string) error
	}

	Service interface {
		Create(ctx context.Context, inp CreateInput) (id string, err error)
		Get(ctx context.Context, id string) (todo Todo, err error)
		GetAll(ctx context.Context, config GetAllInput) (todos []Todo, err error)
		Update(ctx context.Context, inp UpdateInput) error
		MarkAsComplete(ctx context.Context, id string) error
		MarkAsNotComplete(ctx context.Context, id string) error
		Delete(ctx context.Context, id string) error
	}

	service struct {
		repo Repository
		log  *zap.Logger
	}
)

func NewService(repo Repository, log *zap.Logger) Service {
	return &service{
		repo: repo,
		log:  log,
	}
}

func (s *service) Create(ctx context.Context, inp CreateInput) (id string, err error) {
	if len(inp.Title) < 6 || len(inp.Title) > 100 {
		return "", ErrInvalidTitle
	}
	if len(inp.Body) > 2000 {
		return "", ErrInvalidBody
	}
	if inp.Deadline.Before(time.Now()) {
		return "", ErrInvalidDeadline
	}
	id, err = s.repo.Create(ctx, inp)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *service) Get(ctx context.Context, id string) (todo Todo, err error) {
	todo, err = s.repo.Get(ctx, id)
	if err != nil {
		return todo, err
	}

	return todo, nil
}

func (s *service) GetAll(ctx context.Context, config GetAllInput) (todos []Todo, err error) {
	todos, err = s.repo.GetAll(ctx, config)
	if err != nil {
		return todos, err
	}

	return todos, nil
}

func (s *service) Update(ctx context.Context, inp UpdateInput) error {
	if len(inp.Title) < 6 || len(inp.Title) > 100 {
		return ErrInvalidTitle
	}
	if len(inp.Body) > 2000 {
		return ErrInvalidBody
	}
	if inp.Deadline.Before(time.Now()) {
		return ErrInvalidDeadline
	}
	err := s.repo.Update(ctx, inp)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) MarkAsComplete(ctx context.Context, id string) error {
	err := s.repo.MarkAsComplete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) MarkAsNotComplete(ctx context.Context, id string) error {
	err := s.repo.MarkAsNotComplete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
