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
}

type PGStorage struct {
	db *sql.DB
}

func CreateDbConnection() (*PGStorage, error) {
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", os.Getenv("PG_USER"), os.Getenv("PG_PASSWORD"), os.Getenv("PG_DB"))
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
		email VARCHAR
	)`, os.Getenv("USERS_TABLE"))

	_, err := d.db.Exec(query)

	if err != nil {
		log.Fatal("Error on table creation")
	}
	return err
}

func (d *PGStorage) CreateUser(u *User) error {
	query := fmt.Sprintf(`insert into %s (alias, name, email) values ($1, $2, $3)`, os.Getenv("USERS_TABLE"))

	_, err := d.db.Query(query)

	return err

}
