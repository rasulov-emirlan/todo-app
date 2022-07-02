package users

import (
	"context"
	"errors"
	"net/mail"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/rasulov-emirlan/todo-app/backends/pkg/log"

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
		SignUp(ctx context.Context, email, password, username string) (SignInOutput, error)
		SignIn(ctx context.Context, email, password string) (SignInOutput, error)

		UnpackAccessKey(ctx context.Context, accessKey string) (JWTaccess, error)
		Refresh(ctx context.Context, refreshKey string) (SignInOutput, error)

		Update(ctx context.Context, inp UpdateInput) error
		Delete(ctx context.Context, id string) error
	}

	service struct {
		repo Repository
		log  *log.Logger

		secretKey []byte
	}
)

func NewService(repo Repository, logger *log.Logger, secretKey []byte) Service {
	return &service{
		repo:      repo,
		log:       logger,
		secretKey: secretKey,
	}
}

func (s *service) SignUp(ctx context.Context, email, password, username string) (SignInOutput, error) {
	if l := len(password); l < 6 || l > 128 {
		return SignInOutput{}, ErrInvalidPassword
	}
	if l := len(username); l < 6 || l > 20 {
		return SignInOutput{}, ErrInvalidUsername
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return SignInOutput{}, ErrInvalidEmail
	}
	email = strings.ToLower(email)
	passwordHash, err := hashPassword(password)
	if err != nil {
		return SignInOutput{}, err
	}
	_, err = s.repo.Create(ctx, email, passwordHash, username)
	if err != nil {
		return SignInOutput{}, err
	}

	return s.SignIn(ctx, email, password)
}

func (s *service) SignIn(ctx context.Context, email, password string) (SignInOutput, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return SignInOutput{}, err
	}
	if err := comparePassword(password, user.PasswordHash); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return SignInOutput{}, ErrWrongPassword
		}
		return SignInOutput{}, err
	}

	return generateKeys(user, s.secretKey, accessLifeTime, refreshLifeTime)
}

func (s *service) Refresh(ctx context.Context, refreshKey string) (SignInOutput, error) {
	token, err := jwt.ParseWithClaims(refreshKey, &JWTrefresh{}, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})
	if err != nil {
		return SignInOutput{}, err
	}
	claims, ok := token.Claims.(*JWTrefresh)
	if !ok || !token.Valid {
		return SignInOutput{}, ErrInvalidRefreshKey
	}
	user, err := s.repo.Get(ctx, claims.ID)
	if err != nil {
		return SignInOutput{}, err
	}
	return generateKeys(user, s.secretKey, accessLifeTime, time.Duration(claims.ExpiresAt))
}

func (s *service) Update(ctx context.Context, inp UpdateInput) error {
	return s.repo.Update(ctx, inp)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
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
