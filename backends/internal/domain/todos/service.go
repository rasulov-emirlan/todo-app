package todos

import (
	"context"

	"github.com/rasulov-emirlan/todo-app/backends/internal/domain/users"
	"github.com/rasulov-emirlan/todo-app/backends/pkg/logging"
	"github.com/rasulov-emirlan/todo-app/backends/pkg/validation"
)

type (
	Repository interface {
		Create(ctx context.Context, inp CreateInput) (id string, err error)
		Get(ctx context.Context, id string) (todo Todo, err error)
		GetAll(ctx context.Context, config GetAllInput) (todos []Todo, err error)
		// Should not update fields that are empty in UpdateInput
		Update(ctx context.Context, inp UpdateInput) error
		MarkAsComplete(ctx context.Context, id string) error
		MarkAsNotComplete(ctx context.Context, id string) error
		Delete(ctx context.Context, id string) error
	}

	UsersRepository interface {
		Get(ctx context.Context, id string) (users.User, error)
	}

	Service interface {
		Create(ctx context.Context, inp CreateInput) (id string, err error)
		Get(ctx context.Context, id string) (todo Todo, err error)
		GetAll(ctx context.Context, config GetAllInput) (todos []Todo, err error)
		// Will not update fields that are empty in UpdateInput.
		// But ID is required

		// userID represents a user that calls this service.
		// With that id we determine if user is allowed to use this service.
		Update(ctx context.Context, userID string, inp UpdateInput) error
		MarkAsComplete(ctx context.Context, userID, id string) error
		MarkAsNotComplete(ctx context.Context, userID, id string) error
		Delete(ctx context.Context, userID, id string) error
	}

	service struct {
		repo      Repository
		uRepo     UsersRepository
		log       *logging.Logger
		validator *validation.Validator
	}
)

func NewService(repo Repository, uRepo UsersRepository, logger *logging.Logger, validator *validation.Validator) Service {
	return &service{
		repo:      repo,
		uRepo:     uRepo,
		log:       logger,
		validator: validator,
	}
}

func (s *service) Create(ctx context.Context, inp CreateInput) (id string, err error) {
	if err := s.validator.ValidateStruct(inp); err != nil {
		return "", err
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

// TODO: find a better way of sending changesets
// maybe a map[customTypeForFields]any
// would be a good solution...or maybe it would be so bad
func (s *service) Update(ctx context.Context, userID string, inp UpdateInput) error {
	if err := s.validator.ValidateStruct(inp); err != nil {
		return err
	}
	ok, err := s.isAllowed(ctx, userID, inp.ID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotAllowed
	}
	err = s.repo.Update(ctx, inp)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) MarkAsComplete(ctx context.Context, userID, id string) error {
	ok, err := s.isAllowed(ctx, userID, id)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotAllowed
	}
	err = s.repo.MarkAsComplete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) MarkAsNotComplete(ctx context.Context, userID, id string) error {
	ok, err := s.isAllowed(ctx, userID, id)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotAllowed
	}
	err = s.repo.MarkAsNotComplete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) Delete(ctx context.Context, userID, id string) error {
	ok, err := s.isAllowed(ctx, userID, id)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotAllowed
	}
	err = s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) isAllowed(ctx context.Context, userID, todoID string) (bool, error) {
	u, err := s.uRepo.Get(ctx, userID)
	if err != nil {
		return false, err
	}
	if u.Role == users.RoleAdmin {
		return true, nil
	}
	t, err := s.repo.Get(ctx, todoID)
	if err != nil {
		return false, err
	}
	if t.Author.ID == u.ID {
		return true, nil
	}
	return false, nil
}
