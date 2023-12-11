package http

import (
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type AuthClaim struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	UserId   string `json:"user_id"`
	jwt.RegisteredClaims
}

func NewJwtConfig(secret string) echojwt.Config {
	return echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(AuthClaim)
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(401, map[string]interface{}{
				"status":  "error",
				"message": err.Error(),
			})
		},
		TokenLookup: "header:Authorization:TSTMY ",
	}
}

func NewJwtConfigMiddleware(secret string) echo.MiddlewareFunc {
	return echojwt.WithConfig(NewJwtConfig(secret))
}
