package token

import (
	"reflect"
	"testing"

	"github.com/dgrijalva/jwt-go"
)

func TestTokenConfig_GenerateToken(t *testing.T) {
	type args struct {
		bodyGenerate TokenBody
	}
	tests := []struct {
		name            string
		tr              TokenConfig
		args            args
		mockFunc        func(string) string
		wantErrValidate bool
		wantErr         bool
		want            TokenBody
	}{
		{
			name: "success flow",
			tr: TokenConfig{
				Secret:        "my_secret_key",
				ExpTimeInHour: 1,
			},
			args: args{
				bodyGenerate: TokenBody{
					UserID: 1,
				},
			},
			mockFunc: func(s string) string {
				return s
			},
			wantErrValidate: false,
			wantErr:         false,
			want: TokenBody{
				UserID: 1,
			},
		},
		{
			name: "error validate invalid userid flow",
			tr: TokenConfig{
				Secret:        "my_secret_key",
				ExpTimeInHour: 1,
			},
			args: args{
				bodyGenerate: TokenBody{
					UserID: 1,
				},
			},
			mockFunc: func(s string) string {
				claims := jwt.MapClaims{
					"username": "username",
					"userid":   "abc",
				}

				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tkn, _ := token.SignedString([]byte("my_secret_key"))
				return tkn
			},
			wantErrValidate: true,
			wantErr:         false,
			want:            TokenBody{},
		},
		{
			name: "error validate invalid value flow",
			tr: TokenConfig{
				Secret:        "my_secret_key",
				ExpTimeInHour: 1,
			},
			args: args{
				bodyGenerate: TokenBody{
					UserID: 1,
				},
			},
			mockFunc: func(s string) string {
				claims := jwt.MapClaims{
					"username": "",
					"userid":   0,
				}

				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tkn, _ := token.SignedString([]byte("my_secret_key"))
				return tkn
			},
			wantErrValidate: true,
			wantErr:         false,
			want:            TokenBody{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tr.GenerateToken(tt.args.bodyGenerate)
			if (err != nil) != tt.wantErr {
				t.Errorf("TokenConfig.GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				req := tt.mockFunc(got)
				gotData, err := tt.tr.ValidateToken(req)
				if (err != nil) != tt.wantErrValidate {
					t.Errorf("TokenConfig.ValidateToken() error = %v, wantErr %v", err, tt.wantErrValidate)
					return
				}

				if !tt.wantErrValidate {
					if (tt.want) != gotData {
						t.Errorf("TokenConfig.ValidateToken() got = %v, want %v", gotData, tt.want)
						return
					}
				}
			}
		})
	}
}

func TestNewTokenMethod(t *testing.T) {
	type args struct {
		secret    string
		expinHour int64
	}
	tests := []struct {
		name string
		args args
		want TokenMethod
	}{
		{
			name: "success flow",
			args: args{
				secret:    "some_secret",
				expinHour: 1,
			},
			want: TokenConfig{
				Secret:        "some_secret",
				ExpTimeInHour: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTokenMethod(tt.args.secret, tt.args.expinHour); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTokenMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}
