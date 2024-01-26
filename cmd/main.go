package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/pkg/hash"
	"github.com/SawitProRecruitment/UserService/pkg/token"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()
	var server generated.ServerInterface = newServer()
	e.Logger.SetLevel(log.DEBUG)
	generated.RegisterHandlers(e, server)
	// Middleware: Logger
	e.Use(middlewareLogger)
	e.Logger.Fatal(e.Start(":1323"))
}

func middlewareLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Custom log message
		fmt.Println("Request received:", c.Request().Method, c.Request().URL.Path)

		// Call the next middleware or handler
		return next(c)
	}
}

func newServer() *handler.Server {
	dbDsn := os.Getenv("DATABASE_URL")
	if dbDsn == "" {
		dbDsn = "postgres://postgres:postgres@localhost:5432/database?sslmode=disable"
	}
	var repo repository.RepositoryInterface = repository.NewRepository(repository.NewRepositoryOptions{
		Dsn: dbDsn,
	})

	cost := os.Getenv("HASH_COST")
	costInt, err := strconv.Atoi(cost)
	if cost == "" || err != nil || costInt < 4 || costInt > 31 {
		costInt = 10
	}

	var hashMethod hash.HashMethod = hash.NewHashMethod(costInt)

	secret := os.Getenv("TOKEN_SECRET")
	var tokenMethod token.TokenMethod = token.NewTokenMethod(secret, 24)

	opts := handler.NewServerOptions{
		Repository: repo,
		Hash:       hashMethod,
		Token:      tokenMethod,
	}
	return handler.NewServer(opts)
}
