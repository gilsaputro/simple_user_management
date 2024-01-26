package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/middleware"
	"github.com/SawitProRecruitment/UserService/pkg/hash"
	"github.com/SawitProRecruitment/UserService/pkg/token"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()
	server := newServer()
	e.Logger.SetLevel(log.DEBUG)
	e.Use(server.middleware.MiddlewareLogger)
	generated.RegisterHandlers(e, server.handler)
	e.Logger.Fatal(e.Start(":1323"))
}

type Server struct {
	handler    *handler.Server
	middleware *middleware.Server
	repository repository.RepositoryInterface
	hash       hash.HashMethod
	token      token.TokenMethod
}

func newServer() Server {
	s := Server{}

	// Init Repo
	{
		dbDsn := os.Getenv("DATABASE_URL")
		if dbDsn == "" {
			dbDsn = "postgres://postgres:postgres@localhost:5432/database?sslmode=disable"
		}
		s.repository = repository.NewRepository(repository.NewRepositoryOptions{
			Dsn: dbDsn,
		})
		fmt.Println("INIT REPO")
	}

	// Init Hash
	{
		cost := os.Getenv("HASH_COST")
		costInt, err := strconv.Atoi(cost)
		if cost == "" || err != nil || costInt < 4 || costInt > 31 {
			costInt = 10
		}

		s.hash = hash.NewHashMethod(costInt)
		fmt.Println("INIT HASH")
	}

	// Init Token
	{
		privateLocation := os.Getenv("PRIVATE_KEY_LOCATION")
		if privateLocation == "" {
			privateLocation = "./config/private_key.pem"
		}

		publicLocation := os.Getenv("PUBLIC_KEY_LOCATION")
		if publicLocation == "" {
			publicLocation = "./config/public_key.pem"
		}
		method, err := token.NewTokenMethod(
			token.NewTokenConfig{
				PrivateKeyLocation: privateLocation,
				PublicKeyLocation:  publicLocation,
				ExpinHour:          24,
			})
		if err != nil {
			panic(err)
		}
		s.token = method
		fmt.Println("INIT TOKEN")
	}

	// Init Middleware
	{
		s.middleware = middleware.NewMiddlewareServer(middleware.NewMiddlewareOptions{
			Token: s.token,
		})
		fmt.Println("INIT MIDDLEWARE")
	}

	// Init Handler
	{
		s.handler = handler.NewServer(handler.NewServerOptions{
			Repository: s.repository,
			Hash:       s.hash,
			Token:      s.token,
		})
		fmt.Println("INIT HANDLER")
	}

	return s
}
