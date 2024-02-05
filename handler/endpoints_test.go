package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/SawitProRecruitment/UserService/pkg/hash"
	"github.com/SawitProRecruitment/UserService/pkg/token"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/labstack/echo/v4"
	"go.uber.org/mock/gomock"
)

func TestServer_RegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	mockHash := hash.NewMockHashMethod(ctrl)
	mockToken := token.NewMockTokenMethod(ctrl)
	type want struct {
		body string
		code int
	}
	tests := []struct {
		name     string
		body     string
		mockFunc func()
		want     want
	}{
		{
			name: "success flow",
			body: `{"fullName":"testing","password":"@Password1","phoneNumber":"+628123456789"}`,
			mockFunc: func() {
				mockHash.EXPECT().HashValue("@Password1").Return([]byte("123456"), nil)
				mockRepo.EXPECT().RegisterUser(gomock.Any(), repository.RegisterUserInput{
					FullName:    "testing",
					Password:    "123456",
					PhoneNumber: "+628123456789",
				}).Return(repository.RegisterUserOutput{
					UserID: 1,
				}, nil)
			},
			want: want{
				code: 200,
				body: `{"id":1}`,
			},
		},
		{
			name: "failed flow",
			body: `{"fullName":"testing","password":"@Password1","phoneNumber":"+628123456789"}`,
			mockFunc: func() {
				mockHash.EXPECT().HashValue("@Password1").Return([]byte("123456"), nil)
				mockRepo.EXPECT().RegisterUser(gomock.Any(), repository.RegisterUserInput{
					FullName:    "testing",
					Password:    "123456",
					PhoneNumber: "+628123456789",
				}).Return(repository.RegisterUserOutput{}, fmt.Errorf("some error"))
			},
			want: want{
				code: 500,
				body: `{"message":"some error"}`,
			},
		},
		{
			name: "failed flow duplicate phone number",
			body: `{"fullName":"testing","password":"@Password1","phoneNumber":"+628123456789"}`,
			mockFunc: func() {
				mockHash.EXPECT().HashValue("@Password1").Return([]byte("123456"), nil)
				mockRepo.EXPECT().RegisterUser(gomock.Any(), repository.RegisterUserInput{
					FullName:    "testing",
					Password:    "123456",
					PhoneNumber: "+628123456789",
				}).Return(repository.RegisterUserOutput{}, fmt.Errorf("duplicate key value violates unique constraint"))
			},
			want: want{
				code: 409,
				body: `{"message":"Phone number already registered"}`,
			},
		},
		{
			name: "failed flow while hashing password",
			body: `{"fullName":"testing","password":"@Password1","phoneNumber":"+628123456789"}`,
			mockFunc: func() {
				mockHash.EXPECT().HashValue("@Password1").Return([]byte(""), fmt.Errorf("some error"))
			},
			want: want{
				code: 500,
				body: `{"message":"some error"}`,
			},
		},
		{
			name:     "failed invalid request full name",
			body:     `{"fullName":"","password":"@Password1","phoneNumber":"+628123456789"}`,
			mockFunc: func() {},
			want: want{
				code: 400,
				body: `{"message":"full name is required"}`,
			},
		},
		{
			name:     "failed invalid request password",
			body:     `{"fullName":"testing","password":"","phoneNumber":"+628123456789"}`,
			mockFunc: func() {},
			want: want{
				code: 400,
				body: `{"message":"password is required"}`,
			},
		},
		{
			name:     "failed invalid request phone number",
			body:     `{"fullName":"testing","password":"@Password1","phoneNumber":""}`,
			mockFunc: func() {},
			want: want{
				code: 400,
				body: `{"message":"phone number is required"}`,
			},
		},
		{
			name:     "failed invalid request length full name",
			body:     `{"fullName":"a","password":"@Password1","phoneNumber":"+628123456789"}`,
			mockFunc: func() {},
			want: want{
				code: 400,
				body: `{"message":"full name length must be between 3 and 60 characters"}`,
			},
		},
		{
			name:     "failed invalid request length password",
			body:     `{"fullName":"testing","password":"@","phoneNumber":"+628123456789"}`,
			mockFunc: func() {},
			want: want{
				code: 400,
				body: `{"message":"password length must be between 6 and 64 characters"}`,
			},
		},
		{
			name:     "failed invalid request length phone number",
			body:     `{"fullName":"testing","password":"@Password1","phoneNumber":"+621"}`,
			mockFunc: func() {},
			want: want{
				code: 400,
				body: `{"message":"phone number length must be between 10 and 13 characters"}`,
			},
		},
		{
			name:     "failed invalid request phone number must start with +62",
			body:     `{"fullName":"testing","password":"@Password1","phoneNumber":"628123456789"}`,
			mockFunc: func() {},
			want: want{
				code: 400,
				body: `{"message":"phone number must start with +62"}`,
			},
		},
		{
			name:     "failed invalid request password must contain at least 1 uppercase, 1 lowercase, 1 number, and 1 symbol",
			body:     `{"fullName":"testing","password":"@password","phoneNumber":"+628123456789"}`,
			mockFunc: func() {},
			want: want{
				code: 400,
				body: `{"message":"password must contain at least 1 uppercase, 1 lowercase, 1 number, and 1 symbol"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewServer(NewServerOptions{
				Repository: mockRepo,
				Hash:       mockHash,
				Token:      mockToken,
			})
			tt.mockFunc()

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/register", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			ctx := e.NewContext(req, rec)

			handler.RegisterUser(ctx)

			if rec.Code != tt.want.code {
				t.Fatalf("RegisterUser status code got =%d, want %d \n", rec.Code, tt.want.code)
			}

			if !reflect.DeepEqual(tt.want.body, strings.ReplaceAll(string(rec.Body.Bytes()), "\n", "")) {
				t.Fatalf("Register Response body got =%s, want %s \n", string(rec.Body.Bytes()), tt.want.body)
			}
		})
	}
}

func TestServer_LoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	mockHash := hash.NewMockHashMethod(ctrl)
	mockToken := token.NewMockTokenMethod(ctrl)
	type want struct {
		body string
		code int
	}
	tests := []struct {
		name     string
		body     string
		mockFunc func()
		want     want
	}{
		{
			name: "success flow",
			body: `{"password":"@Password1","phoneNumber":"+628123456789"}`,
			mockFunc: func() {
				mockRepo.EXPECT().LoginUser(gomock.Any(), repository.LoginUserInput{
					PhoneNumber: "+628123456789",
				}).Return(repository.LoginUserOutput{
					UserID:   1,
					Password: "123456",
				}, nil)
				mockHash.EXPECT().CompareValue("123456", "@Password1").Return(true)
				mockToken.EXPECT().GenerateToken(gomock.Any()).Return("Bearer token", nil)
				mockRepo.EXPECT().IncrementLoginCount(gomock.Any(), 1).Return(nil)
			},
			want: want{
				code: 200,
				body: `{"id":1,"jwt":"Bearer token"}`,
			},
		},
		{
			name:     "failed flow invalid phone number",
			body:     `{"password":"@Password1","phoneNumber":""}`,
			mockFunc: func() {},
			want: want{
				code: 400,
				body: `{"message":"phone number is required"}`,
			},
		},
		{
			name:     "failed flow invalid password",
			body:     `{"password":"","phoneNumber":"+628123456789"}`,
			mockFunc: func() {},
			want: want{
				code: 400,
				body: `{"message":"password is required"}`,
			},
		},
		{
			name: "failed on repository login user",
			body: `{"password":"@Password1","phoneNumber":"+628123456789"}`,
			mockFunc: func() {
				mockRepo.EXPECT().LoginUser(gomock.Any(), repository.LoginUserInput{
					PhoneNumber: "+628123456789",
				}).Return(repository.LoginUserOutput{}, fmt.Errorf("some error"))
			},
			want: want{
				code: 500,
				body: `{"message":"some error"}`,
			},
		},
		{
			name: "failed on repository login user no rows in result set",
			body: `{"password":"@Password1","phoneNumber":"+628123456789"}`,
			mockFunc: func() {
				mockRepo.EXPECT().LoginUser(gomock.Any(), repository.LoginUserInput{
					PhoneNumber: "+628123456789",
				}).Return(repository.LoginUserOutput{}, fmt.Errorf("no rows in result set"))
			},
			want: want{
				code: 400,
				body: `{"message":"invalid phone number or password"}`,
			},
		},
		{
			name: "failed on repository login user invalid password",
			body: `{"password":"@Password1","phoneNumber":"+628123456789"}`,
			mockFunc: func() {
				mockRepo.EXPECT().LoginUser(gomock.Any(), repository.LoginUserInput{
					PhoneNumber: "+628123456789",
				}).Return(repository.LoginUserOutput{
					UserID:   1,
					Password: "123456",
				}, nil)
				mockHash.EXPECT().CompareValue("123456", "@Password1").Return(false)
			},
			want: want{
				code: 400,
				body: `{"message":"invalid password"}`,
			},
		},
		{
			name: "failed on token generate token",
			body: `{"password":"@Password1","phoneNumber":"+628123456789"}`,
			mockFunc: func() {
				mockRepo.EXPECT().LoginUser(gomock.Any(), repository.LoginUserInput{
					PhoneNumber: "+628123456789",
				}).Return(repository.LoginUserOutput{
					UserID:   1,
					Password: "123456",
				}, nil)
				mockHash.EXPECT().CompareValue("123456", "@Password1").Return(true)
				mockToken.EXPECT().GenerateToken(gomock.Any()).Return("", fmt.Errorf("some error"))
			},
			want: want{
				code: 500,
				body: `{"message":"some error"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewServer(NewServerOptions{
				Repository: mockRepo,
				Hash:       mockHash,
				Token:      mockToken,
			})
			tt.mockFunc()

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/login", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			ctx := e.NewContext(req, rec)

			handler.LoginUser(ctx)

			if rec.Code != tt.want.code {
				t.Fatalf("LoginUser status code got =%d, want %d \n", rec.Code, tt.want.code)
			}

			if !reflect.DeepEqual(tt.want.body, strings.ReplaceAll(string(rec.Body.Bytes()), "\n", "")) {
				t.Fatalf("LoginUser Response body got =%s, want %s \n", string(rec.Body.Bytes()), tt.want.body)
			}
		})
	}
}

func TestServer_GetMyProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	mockHash := hash.NewMockHashMethod(ctrl)
	mockToken := token.NewMockTokenMethod(ctrl)
	type want struct {
		body string
		code int
	}
	tests := []struct {
		name     string
		userID   int
		mockFunc func()
		want     want
	}{
		{
			name:   "success flow",
			userID: 1,
			mockFunc: func() {
				mockRepo.EXPECT().GetUser(gomock.Any(), repository.GetUserInput{
					UserID: 1,
				}).Return(repository.GetUserOutput{
					FullName:    "testing",
					UserID:      1,
					PhoneNumber: "+628123456789",
				}, nil)
			},
			want: want{
				code: 200,
				body: `{"id":1,"name":"testing","phoneNumber":"+628123456789"}`,
			},
		},
		{
			name:     "failed flow invalid user id",
			userID:   0,
			mockFunc: func() {},
			want: want{
				code: 403,
				body: `{"message":"invalid token"}`,
			},
		},
		{
			name:   "failed flow on repository get user",
			userID: 1,
			mockFunc: func() {
				mockRepo.EXPECT().GetUser(gomock.Any(), repository.GetUserInput{
					UserID: 1,
				}).Return(repository.GetUserOutput{}, fmt.Errorf("some error"))
			},
			want: want{
				code: 500,
				body: `{"message":"some error"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewServer(NewServerOptions{
				Repository: mockRepo,
				Hash:       mockHash,
				Token:      mockToken,
			})
			tt.mockFunc()

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/my-profile", nil)
			rec := httptest.NewRecorder()
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			ctx := e.NewContext(req, rec)
			if tt.userID > 0 {
				ctx.Set("user_id", tt.userID)
			}

			handler.GetMyProfile(ctx)

			if rec.Code != tt.want.code {
				t.Fatalf("GetMyProfile status code got =%d, want %d \n", rec.Code, tt.want.code)
			}

			if !reflect.DeepEqual(tt.want.body, strings.ReplaceAll(string(rec.Body.Bytes()), "\n", "")) {
				t.Fatalf("GetMyProfile Response body got =%s, want %s \n", string(rec.Body.Bytes()), tt.want.body)
			}
		})
	}
}

func TestServer_UpdateMyProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	mockHash := hash.NewMockHashMethod(ctrl)
	mockToken := token.NewMockTokenMethod(ctrl)
	type want struct {
		body string
		code int
	}
	tests := []struct {
		name     string
		userID   int
		body     string
		mockFunc func()
		want     want
	}{
		{
			name:   "success flow",
			userID: 1,
			body:   `{"fullName":"testing","password":"@Password1","phoneNumber":"+628123456789"}`,
			mockFunc: func() {
				mockHash.EXPECT().HashValue("@Password1").Return([]byte("123456"), nil)
				mockRepo.EXPECT().UpdateUser(gomock.Any(), repository.UpdateUserInput{
					UserID:      1,
					FullName:    "testing",
					PhoneNumber: "+628123456789",
					Password:    "123456",
				}).Return(repository.UpdateUserOutput{
					FullName:    "testing",
					UserID:      1,
					PhoneNumber: "+628123456789",
				}, nil)
			},
			want: want{
				code: 200,
				body: `{"id":1,"name":"testing","phoneNumber":"+628123456789"}`,
			},
		},
		{
			name:     "failed flow invalid password length",
			userID:   1,
			body:     `{"fullName":"testing","password":"@","phoneNumber":"+628123456789"}`,
			mockFunc: func() {},
			want: want{
				code: 400,
				body: `{"message":"password length must be between 6 and 64 characters"}`,
			},
		},
		{
			name:     "failed flow invalid password contain at least 1 uppercase, 1 lowercase, 1 number, and 1 symbol",
			userID:   1,
			body:     `{"fullName":"testing","password":"@password","phoneNumber":"+628123456789"}`,
			mockFunc: func() {},
			want: want{
				code: 400,
				body: `{"message":"password must contain at least 1 uppercase, 1 lowercase, 1 number, and 1 symbol"}`,
			},
		},
		{
			name:     "failed flow invalid full name length",
			userID:   1,
			body:     `{"fullName":"a","password":"@Password1","phoneNumber":"+628123456789"}`,
			mockFunc: func() {},
			want: want{
				code: 400,
				body: `{"message":"full name length must be between 3 and 60 characters"}`,
			},
		},
		{
			name:     "failed flow invalid phone number length",
			userID:   1,
			body:     `{"fullName":"testing","password":"@Password1","phoneNumber":"+621"}`,
			mockFunc: func() {},
			want: want{
				code: 400,
				body: `{"message":"phone number length must be between 10 and 13 characters"}`,
			},
		},
		{
			name:     "failed flow invalid phone number must start with +62",
			userID:   1,
			body:     `{"fullName":"testing","password":"@Password1","phoneNumber":"628123456789"}`,
			mockFunc: func() {},
			want: want{
				code: 400,
				body: `{"message":"phone number must start with +62"}`,
			},
		},
		{
			name:     "failed flow invalid user id",
			userID:   0,
			body:     `{"fullName":"testing","password":"@Password1","phoneNumber":"+628123456789"}`,
			mockFunc: func() {},
			want: want{
				code: 403,
				body: `{"message":"invalid token"}`,
			},
		},
		{
			name:   "failed flow on repository update user",
			userID: 1,
			body:   `{"fullName":"testing","password":"@Password1","phoneNumber":"+628123456789"}`,
			mockFunc: func() {
				mockHash.EXPECT().HashValue("@Password1").Return([]byte("123456"), nil)
				mockRepo.EXPECT().UpdateUser(gomock.Any(), repository.UpdateUserInput{
					UserID:      1,
					FullName:    "testing",
					PhoneNumber: "+628123456789",
					Password:    "123456",
				}).Return(repository.UpdateUserOutput{}, fmt.Errorf("some error"))
			},
			want: want{
				code: 500,
				body: `{"message":"some error"}`,
			},
		},
		{
			name:   "failed flow on repository update user duplicate phone number",
			userID: 1,
			body:   `{"fullName":"testing","password":"@Password1","phoneNumber":"+628123456789"}`,
			mockFunc: func() {
				mockHash.EXPECT().HashValue("@Password1").Return([]byte("123456"), nil)
				mockRepo.EXPECT().UpdateUser(gomock.Any(), repository.UpdateUserInput{
					UserID:      1,
					FullName:    "testing",
					PhoneNumber: "+628123456789",
					Password:    "123456",
				}).Return(repository.UpdateUserOutput{}, fmt.Errorf("duplicate key value violates unique constraint"))
			},
			want: want{
				code: 409,
				body: `{"message":"Phone number already registered"}`,
			},
		},
		{
			name:   "failed flow on repository update user while hashing password",
			userID: 1,
			body:   `{"fullName":"testing","password":"@Password1","phoneNumber":"+628123456789"}`,
			mockFunc: func() {
				mockHash.EXPECT().HashValue("@Password1").Return([]byte(""), fmt.Errorf("some error"))
			},
			want: want{
				code: 500,
				body: `{"message":"some error"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewServer(NewServerOptions{
				Repository: mockRepo,
				Hash:       mockHash,
				Token:      mockToken,
			})
			tt.mockFunc()

			e := echo.New()
			req := httptest.NewRequest(http.MethodPut, "/my-profile", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			ctx := e.NewContext(req, rec)
			if tt.userID > 0 {
				ctx.Set("user_id", tt.userID)
			}

			handler.UpdateMyProfile(ctx)

			if rec.Code != tt.want.code {
				t.Fatalf("UpdateMyProfile status code got =%d, want %d \n", rec.Code, tt.want.code)
			}

			if !reflect.DeepEqual(tt.want.body, strings.ReplaceAll(string(rec.Body.Bytes()), "\n", "")) {
				t.Fatalf("UpdateMyProfile Response body got =%s, want %s \n", string(rec.Body.Bytes()), tt.want.body)
			}
		})
	}
}
