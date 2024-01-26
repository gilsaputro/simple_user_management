package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/pkg/token"
	"github.com/labstack/echo/v4"
)

func (s *Server) MiddlewareLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Custom log message
		fmt.Println("Request received:", c.Request().Method, c.Request().URL.Path)

		var excludePaths = map[string]bool{
			"/login":    true,
			"/register": true,
		}

		if !excludePaths[c.Request().URL.Path] {
			// Call the next middleware or handler
			authHeader := c.Request().Header.Get("Authorization")

			// validate token
			auth := strings.Split(authHeader, "Bearer ")
			if len(auth) != 2 {
				return c.JSON(http.StatusForbidden, generated.ErrorResponse{
					Message: "Invalid Authorization Header",
				})
			}
			body, err := s.Token.ValidateToken(auth[1])
			if err != nil {
				return c.JSON(http.StatusForbidden, generated.ErrorResponse{
					Message: "Invalid Authorization Header",
				})
			}

			c.Set("user_id", body.UserID)
			return next(c)
		}
		// Call the next middleware or handler
		return next(c)
	}
}

type NewMiddlewareOptions struct {
	Token token.TokenMethod
}

type Server struct {
	Token token.TokenMethod
}

func NewMiddlewareServer(opt NewMiddlewareOptions) *Server {
	return &Server{
		Token: opt.Token,
	}
}
