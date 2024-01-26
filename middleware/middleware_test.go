package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SawitProRecruitment/UserService/pkg/token"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
)

func TestServer_MiddlewareLogger(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockToken := token.NewMockTokenMethod(ctrl)
	mid := NewMiddlewareServer(NewMiddlewareOptions{
		Token: mockToken,
	})
	e := echo.New()
	e.Use(mid.MiddlewareLogger)
	e.GET("/my-profile", func(c echo.Context) error {
		userID := c.Get("user_id")
		if userID == nil || userID.(int) <= 0 {
			return c.JSON(http.StatusUnauthorized, nil)
		}
		fmt.Println("user_id", userID)
		c.Response().Header().Set("result_user_id", fmt.Sprintf("%v", userID))
		return c.JSON(http.StatusOK, nil)
	})
	type want struct {
		code   int
		userID string
	}
	tests := []struct {
		name     string
		token    string
		mockFunc func()
		want     want
	}{
		{
			name:  "success",
			token: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			mockFunc: func() {
				mockToken.EXPECT().ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9").Return(token.TokenBody{
					UserID: 1,
				}, nil)
			},
			want: want{
				code:   http.StatusOK,
				userID: "1",
			},
		},
		{
			name:  "invalid token",
			token: "Bearer",
			mockFunc: func() {
			},
			want: want{
				code: http.StatusForbidden,
			},
		},
		{
			name:  "error validate token",
			token: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			mockFunc: func() {
				mockToken.EXPECT().ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9").Return(token.TokenBody{
					UserID: 1,
				}, fmt.Errorf("some error"))
			},
			want: want{
				code: http.StatusForbidden,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			req := httptest.NewRequest(http.MethodGet, "/my-profile", nil)
			req.Header.Add("Authorization", tt.token)
			w := httptest.NewRecorder()
			e.ServeHTTP(w, req)

			if w.Code != tt.want.code {
				t.Errorf("Server.MiddlewareLogger() = %v, want %v", w.Code, tt.want.code)
			}

			resultUserID := w.Header().Get("result_user_id")
			if resultUserID != tt.want.userID {
				t.Errorf("Server.MiddlewareLogger() = %v, want %v", resultUserID, tt.want.userID)
			}
		})
	}
}
