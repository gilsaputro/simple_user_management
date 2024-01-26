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

// This is just a test endpoint to get you started. Please delete this endpoint.
// (GET /hello)
func (s *Server) Hello(ctx echo.Context, params generated.HelloParams) error {
	fmt.Println("TEST")
	var resp generated.HelloResponse
	resp.Message = fmt.Sprintf("Hello User %d", params.Id)
	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) RegisterUser(ctx echo.Context) error {
	var resp generated.RegisterResponse
	var req = generated.RegisterRequest{}
	ctx.Bind(&req)

	err := validateRegisterRequest(req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	hashPassword, err := s.Hash.HashValue(req.Password)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	result, err := s.Repository.RegisterUser(ctx.Request().Context(), repository.RegisterUserInput{
		FullName:    req.FullName,
		Password:    string(hashPassword),
		PhoneNumber: req.PhoneNumber,
	})

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
				Message: "Phone number already registered",
			})
		}
		ctx.JSON(http.StatusInternalServerError, err)
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

	// validate fullname length
	if len(req.FullName) < 3 || len(req.FullName) > 60 {
		return fmt.Errorf("full name length must be between 3 and 60 characters")
	}

	// validate password
	if len(req.Password) < 6 || len(req.Password) > 64 {
		return fmt.Errorf("password length must be between 6 and 64 characters")
	}

	// validate phone number length
	if len(req.PhoneNumber) < 10 || len(req.PhoneNumber) > 13 {
		return fmt.Errorf("phone number length must be between 10 and 13 characters")
	}

	// validate phone number must start with +62
	if !strings.HasPrefix(req.PhoneNumber, "+62") {
		return fmt.Errorf("phone number must start with +62")
	}

	// validate password contains at least 1 uppercase, 1 lowercase, 1 number, and 1 symbol
	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)
	for _, c := range req.Password {
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
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	val := s.Hash.CompareValue(result.Password, req.Password)
	if !val {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "invalid phone number or password",
		})
	}

	resp.Id = int(result.UserID)

	// generate token
	resp.Jwt, err = s.Token.GenerateToken(token.TokenBody{
		UserID: int(result.UserID),
	})

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	err = s.Repository.IncrementLoginCount(ctx.Request().Context(), int(result.UserID))
	if err != nil {
		fmt.Println(1)
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) GetMyProfile(ctx echo.Context) error {
	// get user id from middleware
	userID := ctx.Get("user_id").(int)
	if userID <= 0 {
		return ctx.JSON(http.StatusForbidden, generated.ErrorResponse{
			Message: "invalid token",
		})
	}

	result, err := s.Repository.GetUser(ctx.Request().Context(), repository.GetUserInput{
		UserID: userID,
	})

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
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
	userID := ctx.Get("user_id").(int)
	if userID <= 0 {
		return ctx.JSON(http.StatusForbidden, generated.ErrorResponse{
			Message: "invalid token",
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
			return ctx.JSON(http.StatusInternalServerError, err)
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
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	response := generated.MyProfileResponse{
		Id:          output.UserID,
		Name:        output.FullName,
		PhoneNumber: output.PhoneNumber,
	}

	return ctx.JSON(http.StatusOK, response)
}

func validateUpdateRequest(req generated.UpdateMyProfileRequest) (repository.UpdateUserInput, error) {
	var resp repository.UpdateUserInput
	// validate request
	if req.FullName != nil {
		if len(*req.FullName) < 3 || len(*req.FullName) > 60 {
			return repository.UpdateUserInput{}, fmt.Errorf("full name length must be between 3 and 60 characters")
		}
		resp.FullName = *req.FullName
	}

	if req.Password != nil {
		if len(*req.Password) < 6 || len(*req.Password) > 64 {
			return repository.UpdateUserInput{}, fmt.Errorf("password length must be between 6 and 64 characters")
		}

		// validate password contains at least 1 uppercase, 1 lowercase, 1 number, and 1 symbol
		var (
			hasUpper   bool
			hasLower   bool
			hasNumber  bool
			hasSpecial bool
		)
		for _, c := range *req.Password {
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

			if hasUpper && hasLower && hasNumber && hasSpecial {
				break
			}
		}

		if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
			return repository.UpdateUserInput{}, fmt.Errorf("password must contain at least 1 uppercase, 1 lowercase, 1 number, and 1 symbol")
		}

		resp.Password = *req.Password
	}

	if req.PhoneNumber != nil {
		// validate phone number length
		if len(*req.PhoneNumber) < 10 || len(*req.PhoneNumber) > 13 {
			return repository.UpdateUserInput{}, fmt.Errorf("phone number length must be between 10 and 13 characters")
		}

		// validate phone number must start with +62
		if !strings.HasPrefix(*req.PhoneNumber, "+62") {
			return repository.UpdateUserInput{}, fmt.Errorf("phone number must start with +62")
		}

		resp.PhoneNumber = *req.PhoneNumber
	}

	return resp, nil
}
