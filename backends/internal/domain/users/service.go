package users

import (
	"context"
	"errors"
	"strings"
	"time"

	// TODO: do something with all these imports
	// maybe create a wrapper for validation, jwt creation and hashing

	"github.com/golang-jwt/jwt"
	"github.com/rasulov-emirlan/todo-app/backends/pkg/logging"
	"github.com/rasulov-emirlan/todo-app/backends/pkg/validation"
	"golang.org/x/crypto/bcrypt"
)

type (
	Repository interface {
		Create(ctx context.Context, email, hashedPassword, username string) (id string, err error)
		Get(ctx context.Context, id string) (user User, err error)
		GetByEmail(ctx context.Context, email string) (user User, err error)
		Update(ctx context.Context, inp UpdateInput) error
		Delete(ctx context.Context, id string) error
	}

	Service interface {
		SignUp(ctx context.Context, inp SignUpInput) (SignInOutput, error)
		SignIn(ctx context.Context, email, password string) (SignInOutput, error)

		UnpackAccessKey(ctx context.Context, accessKey string) (JWTaccess, error)
		Refresh(ctx context.Context, refreshKey string) (SignInOutput, error)

		Update(ctx context.Context, inp UpdateInput) error
		Delete(ctx context.Context, id string) error
	}

	service struct {
		repo       Repository
		validation *validation.Validator
		log        *logging.Logger

		secretKey []byte
	}
)

func NewService(repo Repository, logger *logging.Logger, validator *validation.Validator, secretKey []byte) (Service, error) {
	return &service{
		repo:       repo,
		log:        logger,
		validation: validator,
		secretKey:  secretKey,
	}, nil
}

func (s *service) SignUp(ctx context.Context, inp SignUpInput) (SignInOutput, error) {
	defer s.log.Sync()
	s.log.Info("users: SignUp(): start")
	if err := s.validation.ValidateStruct(inp); err != nil {
		s.log.Debug(
			"users: SignUp(): invalid info was provided",
			// make sure not to log passwords anywhere
			logging.String("username", inp.Username),
			logging.String("email", inp.Email),
		)
		return SignInOutput{}, err
	}
	// Technicaly emails are case sensitive.
	// But we keep them all lower case
	// to make our lifes easier, synce a lot of
	// other services are not case sensitive
	inp.Email = strings.ToLower(inp.Email)
	passwordHash, err := hashPassword(inp.Password)
	if err != nil {
		s.log.Error("users: SignUp(): could not hash password", logging.String("error", err.Error()))
		return SignInOutput{}, err
	}
	_, err = s.repo.Create(ctx, inp.Email, passwordHash, inp.Username)
	if err != nil {
		s.log.Debug("users: SignUp(): could not create user in database", logging.String("error", err.Error()))
		return SignInOutput{}, err
	}

	return s.SignIn(ctx, inp.Email, inp.Password)
}

// TODO: add a remember me option
func (s *service) SignIn(ctx context.Context, email, password string) (SignInOutput, error) {
	defer s.log.Sync()
	s.log.Info("users: SignIn(): start")
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return SignInOutput{}, err
	}
	if err := comparePassword(password, user.PasswordHash); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			s.log.Debug("users: SignIn(): wrong password")
			return SignInOutput{}, ErrWrongPassword
		}
		// TODO: not sure what to do here. is it ok to return err from bcrypt?
		s.log.Error("users: SignIn(): could not compare passwords", logging.String("error", err.Error()))
		return SignInOutput{}, ErrWrongPassword
	}

	return generateKeys(user, s.secretKey, accessLifeTime, refreshLifeTime)
}

func (s *service) Refresh(ctx context.Context, refreshKey string) (SignInOutput, error) {
	defer s.log.Sync()
	s.log.Info("users: Refresh(): start")
	token, err := jwt.ParseWithClaims(refreshKey, &JWTrefresh{}, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})
	if err != nil {
		s.log.Debug("users: Refresh(): could not parse claims", logging.String("error", err.Error()))
		return SignInOutput{}, err
	}
	claims, ok := token.Claims.(*JWTrefresh)
	if !ok || !token.Valid {
		return SignInOutput{}, ErrInvalidRefreshKey
	}
	user, err := s.repo.Get(ctx, claims.ID)
	if err != nil {
		s.log.Debug(
			"users: Debug(): could not get user by id",
			logging.String("id", claims.ID),
			logging.String("error", err.Error()),
		)
		return SignInOutput{}, err
	}

	// TODO: add a remember me option or at least think about it
	return generateKeys(user, s.secretKey, accessLifeTime, refreshLifeTime)
}

func (s *service) Update(ctx context.Context, inp UpdateInput) error {
	defer s.log.Sync()
	s.log.Info("users: Update(): start")
	return s.repo.Update(ctx, inp)
}

func (s *service) Delete(ctx context.Context, id string) error {
	defer s.log.Sync()
	s.log.Info("users: Delete(): start")
	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.log.Debug(
			"users: Delete(): could not delete user",
			logging.String("id", id),
			logging.String("error", err.Error()),
		)
		return err
	}
	s.log.Info("users: Delete(): user was deleted", logging.String("id", id))
	return nil
}

func generateKeys(user User, secretKey []byte, accessEXP, refreshEXP time.Duration) (SignInOutput, error) {
	expAccess := time.Now().Add(accessEXP)
	expRefresh := time.Now().Add(refreshEXP)

	claimsAccess := JWTaccess{
		ID:   user.ID,
		Role: user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expAccess.Unix(),
		},
	}
	claimsRefresh := JWTrefresh{
		ID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expRefresh.Unix(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsAccess)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)

	accessKey, err := accessToken.SignedString(secretKey)
	if err != nil {
		return SignInOutput{}, err
	}
	refreshKey, err := refreshToken.SignedString(secretKey)
	if err != nil {
		return SignInOutput{}, err
	}

	return SignInOutput{
		AccessKey:  accessKey,
		RefreshKey: refreshKey,
	}, nil
}

func (s *service) UnpackAccessKey(ctx context.Context, accessKey string) (JWTaccess, error) {
	// TODO: idk. i thinkg this method might be used too often.
	// so maybe we should not log anything in here. Or log everything in
	// debug level :|
	token, err := jwt.ParseWithClaims(accessKey, &JWTaccess{}, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})
	if err != nil {
		return JWTaccess{}, err
	}
	claims, ok := token.Claims.(*JWTaccess)
	if !ok || !token.Valid {
		return JWTaccess{}, err
	}
	return *claims, nil
}
