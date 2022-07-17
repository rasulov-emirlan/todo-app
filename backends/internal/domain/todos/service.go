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
	defer s.log.Sync()
	s.log.Info("todos: Create(): start")

	if err := s.validator.ValidateStruct(inp); err != nil {
		s.log.Debug(
			"todos: Create(): validation failed",
			logging.String("error", err.Error()),
		)
		return "", err
	}

	id, err = s.repo.Create(ctx, inp)
	if err != nil {
		s.log.Debug(
			"todos: Create(): could not create todo in db",
			logging.String("error", err.Error()),
		)
		return "", err
	}

	return id, nil
}

func (s *service) Get(ctx context.Context, id string) (todo Todo, err error) {
	defer s.log.Sync()
	s.log.Info("todos: Get(): start")

	todo, err = s.repo.Get(ctx, id)
	if err != nil {
		s.log.Debug(
			"todos: Get(): could not get todo from db",
			logging.String("error", err.Error()),
		)
		return todo, err
	}

	return todo, nil
}

func (s *service) GetAll(ctx context.Context, config GetAllInput) (todos []Todo, err error) {
	defer s.log.Sync()
	s.log.Info("todos: GetAll(): start")

	todos, err = s.repo.GetAll(ctx, config)
	if err != nil {
		s.log.Debug(
			"todos: GetAll(): could not get todos from db",
			logging.String("error", err.Error()),
		)
		return todos, err
	}

	return todos, nil
}

// TODO: find a better way of sending changesets
// maybe a map[customTypeForFields]any
// would be a good solution...or maybe it would be so bad
func (s *service) Update(ctx context.Context, userID string, inp UpdateInput) error {
	defer s.log.Sync()
	s.log.Info("todos: Update(): start")

	if err := s.validator.ValidateStruct(inp); err != nil {
		s.log.Debug(
			"todos: Update(): validation failed",
			logging.String("error", err.Error()),
		)
		return err
	}

	ok, err := s.isAllowed(ctx, userID, inp.ID)
	if err != nil {
		s.log.Debug(
			"todos: Update(): isAllowed returned error",
			logging.String("error", err.Error()),
		)
		return err
	}

	if !ok {
		s.log.Debug(
			"todos: Update(): user is not allowed",
			logging.String("userID", userID),
		)
		return ErrNotAllowed
	}

	err = s.repo.Update(ctx, inp)
	if err != nil {
		s.log.Debug(
			"todos: Update(): could not update todo in db",
			logging.String("error", err.Error()),
		)
		return err
	}

	return nil
}

func (s *service) MarkAsComplete(ctx context.Context, userID, id string) error {
	defer s.log.Sync()
	s.log.Info("todos: MarkAsComplete(): start")

	ok, err := s.isAllowed(ctx, userID, id)
	if err != nil {
		s.log.Debug(
			"todos: MarkAsComplete(): isAllowed returned error",
			logging.String("error", err.Error()),
		)
		return err
	}
	if !ok {
		s.log.Debug(
			"todos: MarkAsComplete(): user is not allowed",
			logging.String("userID", userID),
		)
		return ErrNotAllowed
	}

	err = s.repo.MarkAsComplete(ctx, id)
	if err != nil {
		s.log.Debug(
			"todos: MarkAsComplete(): could not mark todo as complete in db",
			logging.String("error", err.Error()),
		)
		return err
	}

	return nil
}

func (s *service) MarkAsNotComplete(ctx context.Context, userID, id string) error {
	defer s.log.Sync()
	s.log.Info("todos: MarkAsNotComplete(): start")

	ok, err := s.isAllowed(ctx, userID, id)
	if err != nil {
		s.log.Debug(
			"todos: MarkAsNotComplete(): isAllowed returned error",
			logging.String("error", err.Error()),
		)
		return err
	}
	if !ok {
		s.log.Debug(
			"todos: MarkAsNotComplete(): user is not allowed",
			logging.String("userID", userID),
		)
		return ErrNotAllowed
	}

	err = s.repo.MarkAsNotComplete(ctx, id)
	if err != nil {
		s.log.Debug(
			"todos: MarkAsNotComplete(): could not mark todo as not complete in db",
			logging.String("error", err.Error()),
		)
		return err
	}

	return nil
}

func (s *service) Delete(ctx context.Context, userID, id string) error {
	defer s.log.Sync()
	s.log.Info("todos: Delete(): start")

	ok, err := s.isAllowed(ctx, userID, id)
	if err != nil {
		s.log.Debug(
			"todos: Delete(): isAllowed returned error",
			logging.String("error", err.Error()),
		)
		return err
	}
	if !ok {
		s.log.Debug(
			"todos: Delete(): user is not allowed",
			logging.String("userID", userID),
		)
		return ErrNotAllowed
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		s.log.Debug(
			"todos: Delete(): could not delete todo from db",
			logging.String("error", err.Error()),
		)
		return err
	}

	return nil
}

// TODO: rewrite this function so we wont use repo calls at all
// We can get all the info we need from users accessKey so yeah
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
