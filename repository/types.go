// This file contains types that are used in the repository layer.
package repository

type RegisterUserInput struct {
	FullName    string
	Password    string
	PhoneNumber string
}

type RegisterUserOutput struct {
	UserID int64
}

type LoginUserInput struct {
	PhoneNumber string
}

type LoginUserOutput struct {
	UserID      int
	PhoneNumber string
	Password    string
}

type GetUserInput struct {
	UserID int
}

type GetUserOutput struct {
	UserID      int
	PhoneNumber string
	FullName    string
}

type UpdateUserInput struct {
	UserID      int
	PhoneNumber string
	FullName    string
	Password    string
}

type UpdateUserOutput struct {
	UserID      int
	PhoneNumber string
	FullName    string
	Password    string
}
