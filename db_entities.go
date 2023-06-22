package main

type User struct {
	ID    int    `json:"id"`
	Alias string `josn:"alias"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
