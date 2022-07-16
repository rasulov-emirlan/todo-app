package users

import "github.com/golang-jwt/jwt"

type (
	SignUpInput struct {
		Email    string `validate:"required,email"`
		Username string `validate:"required,gt=6,lt=20"`
		Password string `validate:"required,gt=6,lt=128"`
	}

	UpdateInput struct {
		ID       string `validate:"required"`
		Username string `validate:"required,gt=6,lt=20"`
		Password string `validate:"required,gt=6,lt=128"`
	}

	SignInOutput struct {
		AccessKey  string `json:"accessKey"`
		RefreshKey string `json:"refreshKey"`
	}

	JWTaccess struct {
		ID   string `json:"userID"`
		Role Role   `json:"role"`

		jwt.StandardClaims
	}

	JWTrefresh struct {
		ID string `json:"userID"`

		jwt.StandardClaims
	}
)
