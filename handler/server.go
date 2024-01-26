package handler

import (
	hash "github.com/SawitProRecruitment/UserService/pkg/hash"
	"github.com/SawitProRecruitment/UserService/pkg/token"
	"github.com/SawitProRecruitment/UserService/repository"
)

type Server struct {
	Repository repository.RepositoryInterface
	Hash       hash.HashMethod
	Token      token.TokenMethod
}

type NewServerOptions struct {
	Repository repository.RepositoryInterface
	Hash       hash.HashMethod
	Token      token.TokenMethod
}

func NewServer(opts NewServerOptions) *Server {
	return &Server{
		Repository: opts.Repository,
		Hash:       opts.Hash,
		Token:      opts.Token,
	}
}
