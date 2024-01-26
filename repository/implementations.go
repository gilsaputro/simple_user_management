package repository

import (
	"context"
)

func (r *Repository) RegisterUser(ctx context.Context, input RegisterUserInput) (output RegisterUserOutput, err error) {
	tx, err := r.Db.Begin()
	if err != nil {
		return RegisterUserOutput{}, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()

	// Menambahkan user baru.
	var userID int
	err = tx.QueryRow("INSERT INTO users (phone_number, password, full_name) VALUES ($1, $2, $3) RETURNING id", input.PhoneNumber, input.Password, input.FullName).Scan(&userID)
	if err != nil {
		return RegisterUserOutput{}, err
	}
	// Commit transaksi jika semuanya berhasil.
	err = tx.Commit()
	if err != nil {
		return RegisterUserOutput{}, err
	}

	return RegisterUserOutput{
		UserID: int64(userID),
	}, nil
}

func (r *Repository) LoginUser(ctx context.Context, input LoginUserInput) (output LoginUserOutput, err error) {
	row := r.Db.QueryRow("SELECT id, phone_number, password FROM users WHERE phone_number = $1", input.PhoneNumber)

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

	// Update user.
	_, err = tx.Exec("UPDATE users SET phone_number = $1, full_name = $2, password = $3 WHERE id = $4", userInfo.PhoneNumber, userInfo.FullName, userInfo.Password, userInfo.UserID)
	if err != nil {
		return UpdateUserOutput{}, err
	}

	// Commit transaksi jika semuanya berhasil.
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
