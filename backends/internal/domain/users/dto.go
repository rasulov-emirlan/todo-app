package users

import "github.com/golang-jwt/jwt"

type (
	SignUpInput struct {
		Username string `validate:"required,username"`
		Email    string `validate:"required,email"`
		Password string `validate:"required,password"`
	}

	UpdateInput struct {
		ID       string
		Username string
		Password string
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
