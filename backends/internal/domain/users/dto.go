package users

import "github.com/golang-jwt/jwt"

type (
	UpdateInput struct {
		ID       string
		Username string
		Password string
	}

	SignInOutput struct {
		AccessKey  string
		RefreshKey string
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
