package users

import (
	"context"
	"errors"
	"strings"
	"time"

	// TODO: do something with all these imports
	// maybe create a wrapper for validation, jwt creation and hashing
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/golang-jwt/jwt"
	"github.com/rasulov-emirlan/todo-app/backends/pkg/logging"
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

		UnpackValidationErrors(err error) []string
	}

	service struct {
		repo       Repository
		validation *validator.Validate
		trans      ut.Translator
		log        *logging.Logger

		secretKey []byte
	}
)

func NewService(repo Repository, logger *logging.Logger, secretKey []byte) (Service, error) {
	en := en.New()
	uni := ut.New(en, en)

	trans, _ := uni.GetTranslator("en")
	v := validator.New()
	en_translations.RegisterDefaultTranslations(v, trans)

	err := v.RegisterTranslation("password", trans,
		func(ut ut.Translator) error {
			return ut.Add("password", "{0} can't be shorter than 6 characters or longer then 128", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, err := ut.T("password", fe.Field())
			if err != nil {
				logger.Fatal("Could not initialize users serivice", log.String("error", err.Error()))
			}
			return t
		},
	)
	if err != nil {
		return nil, err
	}
	err = v.RegisterTranslation("username", trans,
		func(ut ut.Translator) error {
			return ut.Add("username", "{0} can't be shorter than 6 characters or longer then 20", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, err := ut.T("username", fe.Field())
			if err != nil {
				logger.Fatal("Could not initialize users serivice", log.String("error", err.Error()))
			}
			return t
		},
	)
	if err != nil {
		return nil, err
	}

	err = v.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		l := len(fl.Field().String())
		return l >= 6 && l <= 128
	})
	if err != nil {
		return nil, err
	}
	err = v.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		l := len(fl.Field().String())
		return l >= 6 && l <= 20
	})
	if err != nil {
		return nil, err
	}

	return &service{
		repo:       repo,
		log:        logger,
		validation: v,
		trans:      trans,
		secretKey:  secretKey,
	}, nil
}

func (s *service) UnpackValidationErrors(err error) []string {
	v, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil
	}
	errs := make([]string, 0, len(v))
	for _, vv := range v {
		errs = append(errs, vv.Translate(s.trans))
	}
	return errs
}

func (s *service) SignUp(ctx context.Context, inp SignUpInput) (SignInOutput, error) {
	defer s.log.Sync()
	s.log.Info("users: SignUp(): start")
	if err := s.validation.Struct(inp); err != nil {
		s.log.Debug(
			"users: SignUp(): invalid info was provided",
			// make sure not to log passwords anywhere
			log.String("username", inp.Username),
			log.String("email", inp.Email),
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
		s.log.Error("users: SignUp(): could not hash password", log.String("error", err.Error()))
		return SignInOutput{}, err
	}
	_, err = s.repo.Create(ctx, inp.Email, passwordHash, inp.Username)
	if err != nil {
		s.log.Debug("users: SignUp(): could not create user in database", log.String("error", err.Error()))
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
		s.log.Error("users: SignIn(): could not compare passwords", log.String("error", err.Error()))
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
		s.log.Debug("users: Refresh(): could not parse claims", log.String("error", err.Error()))
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
			log.String("id", claims.ID),
			log.String("error", err.Error()),
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
			log.String("id", id),
			log.String("error", err.Error()),
		)
		return err
	}
	s.log.Info("users: Delete(): user was deleted", log.String("id", id))
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
