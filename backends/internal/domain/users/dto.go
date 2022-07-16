package users

import "github.com/golang-jwt/jwt"

type (
	SignUpInput struct {
		Username string `validate:"required,gt=6,lt=20"`
		Email    string `validate:"required,email"`
		Password string `validate:"required,gt=6,lt=128"`
	}

	UpdateInput struct {
		ID       string `validate:"required"`
		Username string `validate:"username"`
		Password string `validate:"password"`
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
