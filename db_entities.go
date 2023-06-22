package main

import (
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `json:"id"`
	Alias     string    `josn:"alias"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	Password  string    `json:"password"`
}

func NewUserFromRequest(userRequest *CreateUserRequest) (*User, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	if string(encryptedPassword) == "" {
		return nil, errors.New("can't encrypt the password or the password was empty")
	}

	return &User{
		Alias:     userRequest.Alias,
		Name:      userRequest.Name,
		Email:     userRequest.Email,
		CreatedAt: time.Now().UTC(),
		Password:  string(encryptedPassword),
	}, nil
}

// TODO move to a utils file
func MapUserRowToUserEnt(rows *sql.Rows) (*User, error) {
	user := new(User)

	err := rows.Scan(
		&user.ID,
		&user.Alias,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	return user, err
}
