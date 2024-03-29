package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/pkg/token"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/labstack/echo/v4"
)

func (s *Server) RegisterUser(ctx echo.Context) error {
	var resp generated.RegisterResponse
	var req = generated.RegisterRequest{}
	var err error
	ctx.Bind(&req)
	err = validateRegisterRequest(req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	hashPassword, err := s.Hash.HashValue(req.Password)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	result, err := s.Repository.RegisterUser(ctx.Request().Context(), repository.RegisterUserInput{
		FullName:    req.FullName,
		Password:    string(hashPassword),
		PhoneNumber: req.PhoneNumber,
	})

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return ctx.JSON(http.StatusConflict, generated.ErrorResponse{
				Message: "Phone number already registered",
			})
		}
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	resp.Id = int(result.UserID)

	return ctx.JSON(http.StatusOK, resp)
}

func validateRegisterRequest(req generated.RegisterRequest) error {
	// validate request
	if req.FullName == "" {
		return fmt.Errorf("full name is required")
	}

	if req.Password == "" {
		return fmt.Errorf("password is required")
	}

	if req.PhoneNumber == "" {
		return fmt.Errorf("phone number is required")
	}

	// validate phone number length
	err := validatePhoneNumber(req.PhoneNumber)
	if err != nil {
		return err
	}

	// validate full name
	err = validateFullName(req.FullName)
	if err != nil {
		return err
	}

	// validate password
	err = validatePassword(req.Password)

	return err
}

func validatePassword(password string) error {
	// validate password
	if len(password) < 6 || len(password) > 64 {
		return fmt.Errorf("password length must be between 6 and 64 characters")
	}

	// validate password contains at least 1 uppercase, 1 lowercase, 1 number, and 1 symbol
	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)
	for _, c := range password {
		switch {
		case c >= 'A' && c <= 'Z':
			hasUpper = true
		case c >= 'a' && c <= 'z':
			hasLower = true
		case c >= '0' && c <= '9':
			hasNumber = true
		case c == '!' || c == '@' || c == '#' || c == '$' || c == '%' || c == '^' || c == '&':
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return fmt.Errorf("password must contain at least 1 uppercase, 1 lowercase, 1 number, and 1 symbol")
	}
	return nil
}

func validatePhoneNumber(phoneNumber string) error {
	// validate phone number
	if len(phoneNumber) < 10 || len(phoneNumber) > 13 {
		return fmt.Errorf("phone number length must be between 10 and 13 characters")
	}

	// validate phone number must start with +62
	if !strings.HasPrefix(phoneNumber, "+62") {
		return fmt.Errorf("phone number must start with +62")
	}
	return nil
}

func validateFullName(fullName string) error {
	// validate full name
	if len(fullName) < 3 || len(fullName) > 60 {
		return fmt.Errorf("full name length must be between 3 and 60 characters")
	}
	return nil
}

func (s *Server) LoginUser(ctx echo.Context) error {
	var resp generated.LoginResponse
	var req = generated.LoginRequest{}
	ctx.Bind(&req)

	if req.Password == "" {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "password is required",
		})
	}

	if req.PhoneNumber == "" {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "phone number is required",
		})
	}

	result, err := s.Repository.LoginUser(ctx.Request().Context(), repository.LoginUserInput{
		PhoneNumber: req.PhoneNumber,
	})
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
				Message: "invalid phone number or password",
			})
		}
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	val := s.Hash.CompareValue(result.Password, req.Password)
	if !val {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "invalid password",
		})
	}

	resp.Id = int(result.UserID)

	// generate token
	resp.Jwt, err = s.Token.GenerateToken(token.TokenBody{
		UserID: int(result.UserID),
	})

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	s.Repository.IncrementLoginCount(ctx.Request().Context(), int(result.UserID))

	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) GetMyProfile(ctx echo.Context) error {
	// get user id from middleware
	userID, err := getUserID(ctx)
	if err != nil {
		return ctx.JSON(http.StatusForbidden, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	result, err := s.Repository.GetUser(ctx.Request().Context(), repository.GetUserInput{
		UserID: userID,
	})

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	resp := generated.MyProfileResponse{
		Id:          result.UserID,
		Name:        result.FullName,
		PhoneNumber: result.PhoneNumber,
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) UpdateMyProfile(ctx echo.Context) error {
	var req = generated.UpdateMyProfileRequest{}
	ctx.Bind(&req)

	// get user id from middleware
	userID, err := getUserID(ctx)
	if err != nil {
		return ctx.JSON(http.StatusForbidden, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	repoRequest, err := validateUpdateRequest(req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	if repoRequest.Password != "" {
		hashPassword, err := s.Hash.HashValue(repoRequest.Password)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
				Message: err.Error(),
			})
		}
		repoRequest.Password = string(hashPassword)
	}

	output, err := s.Repository.UpdateUser(ctx.Request().Context(), repository.UpdateUserInput{
		UserID:      userID,
		FullName:    repoRequest.FullName,
		Password:    repoRequest.Password,
		PhoneNumber: repoRequest.PhoneNumber,
	})

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return ctx.JSON(http.StatusConflict, generated.ErrorResponse{
				Message: "Phone number already registered",
			})
		}
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	response := generated.MyProfileResponse{
		Id:          output.UserID,
		Name:        output.FullName,
		PhoneNumber: output.PhoneNumber,
	}

	return ctx.JSON(http.StatusOK, response)
}

func getUserID(ctx echo.Context) (int, error) {
	userID := ctx.Get("user_id")
	if userID == nil || userID.(int) <= 0 {
		return 0, fmt.Errorf("invalid token")
	}
	return userID.(int), nil
}

func validateUpdateRequest(req generated.UpdateMyProfileRequest) (repository.UpdateUserInput, error) {
	var resp repository.UpdateUserInput
	// validate request
	if req.FullName != nil {
		err := validateFullName(*req.FullName)
		if err != nil {
			return repository.UpdateUserInput{}, err
		}
		resp.FullName = *req.FullName
	}

	if req.Password != nil {
		err := validatePassword(*req.Password)
		if err != nil {
			return repository.UpdateUserInput{}, err
		}
		resp.Password = *req.Password
	}

	if req.PhoneNumber != nil {
		err := validatePhoneNumber(*req.PhoneNumber)
		if err != nil {
			return repository.UpdateUserInput{}, err
		}
		resp.PhoneNumber = *req.PhoneNumber
	}

	return resp, nil
}
