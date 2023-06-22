package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateUser(*User) error
	GetUserByEmail(string) (*User, error)
	DoesUserEmailExist(string) (bool, error)
}

type PGStorage struct {
	db *sql.DB
}

func CreateDbConnection() (*PGStorage, error) {
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", os.Getenv("PG_USER"), os.Getenv("PG_DB"), os.Getenv("PG_PASSWORD"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PGStorage{
		db: db,
	}, nil
}

func (d *PGStorage) createTable() error {
	query := fmt.Sprintf(`CREATE TABLE if not exists %s (
		id SERIAL PRIMARY KEY,
		alias VARCHAR,
		name VARCHAR,
		email VARCHAR,
		createdAt timestamp,
		password VARCHAR
	)`, os.Getenv("USERS_TABLE"))

	_, err := d.db.Exec(query)

	if err != nil {
		log.Fatal("Error on table creation", err)
	}
	return err
}

func (d *PGStorage) CreateUser(u *User) error {
	query := fmt.Sprintf(`insert into %s (alias, name, email, createdAt, password) values ($1, $2, $3, $4, $5)`, os.Getenv("USERS_TABLE"))

	_, err := d.db.Query(query, u.Alias, u.Name, u.Email, u.CreatedAt, u.Password)

	return err

}

func (d *PGStorage) DoesUserEmailExist(email string) (bool, error) {
	query := fmt.Sprintf(`select id from %s where email = $1`, os.Getenv("USERS_TABLE"))
	var userId int

	err := d.db.QueryRow(query, email).Scan(&userId)
	if err != nil {
		return false, err
	}

	return true, nil
}

// TODO select all columns to map the user
func (d *PGStorage) GetUserByEmail(email string) (*User, error) {
	query := fmt.Sprintf(`select id from %s 
	where email = $1`, os.Getenv("USERS_TABLE"))

	rows, err := d.db.Query(query, email)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return MapUserRowToUserEnt(rows)
	}

	return nil, fmt.Errorf("could not find user with email %s", email)
}
