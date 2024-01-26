// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

type RepositoryInterface interface {
	RegisterUser(ctx context.Context, req RegisterUserInput) (RegisterUserOutput, error)
	LoginUser(ctx context.Context, req LoginUserInput) (LoginUserOutput, error)
	GetUser(ctx context.Context, req GetUserInput) (GetUserOutput, error)
	UpdateUser(ctx context.Context, req UpdateUserInput) (UpdateUserOutput, error)
	IncrementLoginCount(ctx context.Context, userID int) (err error)
}
