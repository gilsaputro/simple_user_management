package repository

import (
	"context"
	"database/sql"
	"time"
)

func (r *Repository) RegisterUser(ctx context.Context, input RegisterUserInput) (output RegisterUserOutput, err error) {
	// Begin transaction.
	tx, err := r.Db.Begin()
	if err != nil {
		return RegisterUserOutput{}, err
	}

	defer func() {
		// Rollback transaction if error.
		if err != nil {
			tx.Rollback()
			return
		}
	}()

	createdTime := time.Now().UTC()

	// Insert user.
	var userID int
	err = tx.QueryRow("INSERT INTO users(phone_number, full_name, password, created_at) VALUES($1, $2, $3, $4) RETURNING id",
		input.PhoneNumber, input.FullName, input.Password, createdTime).Scan(&userID)
	if err != nil {
		return RegisterUserOutput{}, err
	}

	// Commit Transaction if everything is ok.
	err = tx.Commit()
	if err != nil {
		return RegisterUserOutput{}, err
	}

	return RegisterUserOutput{
		UserID: int64(userID),
	}, nil
}

func (r *Repository) LoginUser(ctx context.Context, input LoginUserInput) (output LoginUserOutput, err error) {
	// query user.
	row := r.Db.QueryRow("SELECT id, phone_number, password FROM users WHERE phone_number = $1", input.PhoneNumber)

	// scan result.
	err = row.Scan(&output.UserID, &output.PhoneNumber, &output.Password)
	if err != nil {
		return LoginUserOutput{}, err
	}

	return output, nil
}

func (r *Repository) GetUser(ctx context.Context, input GetUserInput) (output GetUserOutput, err error) {
	row := r.Db.QueryRow("SELECT id, phone_number, full_name FROM users WHERE id = $1", input.UserID)

	err = row.Scan(&output.UserID, &output.PhoneNumber, &output.FullName)
	if err != nil {
		return GetUserOutput{}, err
	}

	return output, nil
}

func (r *Repository) UpdateUser(ctx context.Context, input UpdateUserInput) (output UpdateUserOutput, err error) {
	tx, err := r.Db.Begin()
	if err != nil {
		return UpdateUserOutput{}, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()

	// get user.
	var userInfo UpdateUserOutput
	row := tx.QueryRow("SELECT id, phone_number, full_name, password FROM users WHERE id = $1", input.UserID)
	err = row.Scan(&userInfo.UserID, &userInfo.PhoneNumber, &userInfo.FullName, &userInfo.Password)
	if err != nil {
		return UpdateUserOutput{}, err
	}

	// Update password.
	if input.Password != "" {
		userInfo.Password = input.Password
	}

	// Update phone number.
	if input.PhoneNumber != "" {
		userInfo.PhoneNumber = input.PhoneNumber
	}

	// Update full name.
	if input.FullName != "" {
		userInfo.FullName = input.FullName
	}

	updatedAt := time.Now().UTC()
	// Update user.
	_, err = tx.Exec("UPDATE users SET phone_number = $1, full_name = $2, password = $3, updated_at = $4 WHERE id = $5", userInfo.PhoneNumber, userInfo.FullName, userInfo.Password, updatedAt, userInfo.UserID)
	if err != nil {
		return UpdateUserOutput{}, err
	}

	// Commit transaction.
	err = tx.Commit()
	if err != nil {
		return UpdateUserOutput{}, err
	}

	return UpdateUserOutput{
		UserID:      userInfo.UserID,
		PhoneNumber: userInfo.PhoneNumber,
		FullName:    userInfo.FullName,
	}, nil
}

func (r *Repository) IncrementLoginCount(ctx context.Context, userID int) (err error) {
	createdTime := time.Now().UTC()

	// Check if the user already has a login entry
	var existingLoginCount int
	err = r.Db.QueryRow("SELECT login_count FROM users_login_history WHERE user_id = $1", userID).Scan(&existingLoginCount)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	// If the user has an existing login entry, update the login_count; otherwise, insert a new entry
	if err == nil {
		_, err = r.Db.Exec("UPDATE users_login_history SET login_count = login_count + 1 WHERE user_id = $1 RETURNING login_count", userID)
	} else {
		_, err = r.Db.Exec("INSERT INTO users_login_history(user_id, created_at) VALUES($1, $2)", userID, createdTime)
	}

	return err
}
