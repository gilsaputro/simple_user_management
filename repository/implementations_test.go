package repository

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"testing"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestRepository_RegisterUser(t *testing.T) {
	db, mockDB, _ := sqlmock.New()
	defer db.Close()
	tests := []struct {
		name       string
		mockFunc   func()
		input      RegisterUserInput
		wantOutput RegisterUserOutput
		wantErr    bool
	}{
		{
			name: "success",
			mockFunc: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery(regexp.QuoteMeta("INSERT INTO users(phone_number, full_name, password, created_at) VALUES($1, $2, $3, $4) RETURNING id")).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mockDB.ExpectCommit()
			},
			input: RegisterUserInput{
				FullName:    "John Doe",
				Password:    "password",
				PhoneNumber: "+6281234567890",
			},
			wantOutput: RegisterUserOutput{
				UserID: 1,
			},
		},
		{
			name: "error while begin transaction",
			mockFunc: func() {
				mockDB.ExpectBegin().WillReturnError(fmt.Errorf("some error"))
			},
			input:      RegisterUserInput{},
			wantOutput: RegisterUserOutput{},
			wantErr:    true,
		},
		{
			name: "error while insert user",
			mockFunc: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery(regexp.QuoteMeta("INSERT INTO users(phone_number, full_name, password, created_at) VALUES($1, $2, $3, $4) RETURNING id")).WillReturnError(fmt.Errorf("some error"))
			},
			input:      RegisterUserInput{},
			wantOutput: RegisterUserOutput{},
			wantErr:    true,
		},
		{
			name: "error while commit transaction",
			mockFunc: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery(regexp.QuoteMeta("INSERT INTO users(phone_number, full_name, password, created_at) VALUES($1, $2, $3, $4) RETURNING id")).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mockDB.ExpectCommit().WillReturnError(fmt.Errorf("some error"))
			},
			input:      RegisterUserInput{},
			wantOutput: RegisterUserOutput{},
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Repository{
				Db: db,
			}

			tt.mockFunc()

			gotOutput, err := r.RegisterUser(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("Repository.RegisterUser() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestRepository_LoginUser(t *testing.T) {
	db, mockDB, _ := sqlmock.New()
	defer db.Close()
	tests := []struct {
		name       string
		input      LoginUserInput
		mockFunc   func()
		wantOutput LoginUserOutput
		wantErr    bool
	}{
		{
			name: "success",
			input: LoginUserInput{
				PhoneNumber: "+6281234567890",
			},
			mockFunc: func() {
				mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, phone_number, password FROM users WHERE phone_number = $1")).
					WithArgs("+6281234567890").WillReturnRows(sqlmock.NewRows([]string{"id", "phone_number", "password"}).AddRow(1, "+6281234567890", "password"))
			},
			wantOutput: LoginUserOutput{
				UserID:      1,
				PhoneNumber: "+6281234567890",
				Password:    "password",
			},
		},
		{
			name: "error while query",
			input: LoginUserInput{
				PhoneNumber: "+6281234567890",
			},
			mockFunc: func() {
				mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, phone_number, password FROM users WHERE phone_number = $1")).WillReturnError(fmt.Errorf("some error"))
			},
			wantOutput: LoginUserOutput{},
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Repository{
				Db: db,
			}
			tt.mockFunc()
			gotOutput, err := r.LoginUser(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.LoginUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("Repository.LoginUser() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestRepository_GetUser(t *testing.T) {
	db, mockDB, _ := sqlmock.New()
	defer db.Close()
	tests := []struct {
		name       string
		mockFunc   func()
		input      GetUserInput
		wantOutput GetUserOutput
		wantErr    bool
	}{
		{
			name: "success",
			mockFunc: func() {
				mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, phone_number, full_name FROM users WHERE id = $1")).
					WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "phone_number", "full_name"}).AddRow(1, "+6281234567890", "John Doe"))
			},
			input: GetUserInput{
				UserID: 1,
			},
			wantOutput: GetUserOutput{
				UserID:      1,
				PhoneNumber: "+6281234567890",
				FullName:    "John Doe",
			},
		},
		{
			name: "error while query",
			mockFunc: func() {
				mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, phone_number, full_name FROM users WHERE id = $1")).
					WithArgs(1).WillReturnError(fmt.Errorf("some error"))
			},
			input:      GetUserInput{},
			wantOutput: GetUserOutput{},
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Repository{
				Db: db,
			}
			tt.mockFunc()
			gotOutput, err := r.GetUser(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("Repository.GetUser() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestRepository_UpdateUser(t *testing.T) {
	db, mockDB, _ := sqlmock.New()
	defer db.Close()
	tests := []struct {
		name       string
		mockFunc   func()
		input      UpdateUserInput
		wantOutput UpdateUserOutput
		wantErr    bool
	}{
		{
			name: "success",
			mockFunc: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, phone_number, full_name, password FROM users WHERE id = $1")).
					WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "phone_number", "full_name", "password"}).AddRow(1, "+6281234567890", "John Doe", "password"))
				mockDB.ExpectExec(regexp.QuoteMeta("UPDATE users SET phone_number = $1, full_name = $2, password = $3, updated_at = $4 WHERE id = $5")).WillReturnResult(sqlmock.NewResult(1, 1))
				mockDB.ExpectCommit()
			},
			input: UpdateUserInput{
				UserID:      1,
				PhoneNumber: "+6281234567890",
				FullName:    "John Doe",
				Password:    "password",
			},
			wantOutput: UpdateUserOutput{
				UserID:      1,
				PhoneNumber: "+6281234567890",
				FullName:    "John Doe",
			},
			wantErr: false,
		},
		{
			name: "error while begin transaction",
			mockFunc: func() {
				mockDB.ExpectBegin().WillReturnError(fmt.Errorf("some error"))
			},
			input:      UpdateUserInput{},
			wantOutput: UpdateUserOutput{},
			wantErr:    true,
		},
		{
			name: "error while query",
			mockFunc: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, phone_number, full_name, password FROM users WHERE id = $1")).
					WithArgs(1).WillReturnError(fmt.Errorf("some error"))
			},
			input:      UpdateUserInput{},
			wantOutput: UpdateUserOutput{},
			wantErr:    true,
		},
		{
			name: "error while update user",
			mockFunc: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, phone_number, full_name, password FROM users WHERE id = $1")).
					WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "phone_number", "full_name", "password"}).AddRow(1, "+6281234567890", "John Doe", "password"))
				mockDB.ExpectExec(regexp.QuoteMeta("UPDATE users SET phone_number = $1, full_name = $2, password = $3, updated_at = $4 WHERE id = $5")).WillReturnError(fmt.Errorf("some error"))
			},
			input:      UpdateUserInput{},
			wantOutput: UpdateUserOutput{},
			wantErr:    true,
		},
		{
			name: "error while commit transaction",
			mockFunc: func() {
				mockDB.ExpectBegin()
				mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, phone_number, full_name, password FROM users WHERE id = $1")).
					WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "phone_number", "full_name", "password"}).AddRow(1, "+6281234567890", "John Doe", "password"))
				mockDB.ExpectExec(regexp.QuoteMeta("UPDATE users SET phone_number = $1, full_name = $2, password = $3, updated_at = $4 WHERE id = $5")).WillReturnResult(sqlmock.NewResult(1, 1))
				mockDB.ExpectCommit().WillReturnError(fmt.Errorf("some error"))
			},
			input:      UpdateUserInput{},
			wantOutput: UpdateUserOutput{},
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Repository{
				Db: db,
			}
			tt.mockFunc()
			gotOutput, err := r.UpdateUser(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("Repository.UpdateUser() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestRepository_IncrementLoginCount(t *testing.T) {
	db, mockDB, _ := sqlmock.New()
	defer db.Close()
	tests := []struct {
		name     string
		mockFunc func()
		userID   int
		wantErr  bool
	}{
		{
			name: "success update login count",
			mockFunc: func() {
				mockDB.ExpectQuery(regexp.QuoteMeta("SELECT login_count FROM users_login_history WHERE user_id = $1")).
					WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"login_count"}).AddRow(1))
				mockDB.ExpectExec(regexp.QuoteMeta("UPDATE users_login_history SET login_count = login_count + 1 WHERE user_id = $1 RETURNING login_count")).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			userID:  1,
			wantErr: false,
		},
		{
			name: "success insert login count",
			mockFunc: func() {
				mockDB.ExpectQuery(regexp.QuoteMeta("SELECT login_count FROM users_login_history WHERE user_id = $1")).
					WithArgs(1).WillReturnError(sql.ErrNoRows)
				mockDB.ExpectExec(regexp.QuoteMeta("INSERT INTO users_login_history(user_id, created_at) VALUES($1, $2)")).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			userID:  1,
			wantErr: false,
		},
		{
			name: "error while query",
			mockFunc: func() {
				mockDB.ExpectQuery(regexp.QuoteMeta("SELECT login_count FROM users_login_history WHERE user_id = $1")).
					WithArgs(1).WillReturnError(fmt.Errorf("some error"))
			},
			userID:  1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Repository{
				Db: db,
			}

			tt.mockFunc()
			if err := r.IncrementLoginCount(context.Background(), tt.userID); (err != nil) != tt.wantErr {
				t.Errorf("Repository.IncrementLoginCount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
