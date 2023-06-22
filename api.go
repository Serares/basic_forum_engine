package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ApiServer struct {
	port string
	db   Storage
}

type BasicResponse struct {
	Message string `json:"message"`
	Error   string
}

type CreateUserRequest struct {
	Alias string `josn:"alias"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewApiServer(port string, db Storage) *ApiServer {
	return &ApiServer{
		port: port,
		db:   db,
	}
}

func (s *ApiServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/", s.handleDefault).Methods("GET")

	fmt.Println("Server is listening on port: ", s.port)
	log.Fatal(http.ListenAndServe(s.port, router))
}

func (s *ApiServer) handleDefault(w http.ResponseWriter, r *http.Request) {
	successMessage := BasicResponse{
		Message: "All good",
	}

	writeJSONResponse(w, http.StatusOK, successMessage)
}

func (s *ApiServer) handleUserCreation(w http.ResponseWriter, r *http.Request) {
	req := new(CreateUserRequest)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONResponse(w, http.StatusMethodNotAllowed, BasicResponse{Message: "The request is not compliant", Error: err.Error()})
	}

	s.db.CreateUser()
	successMessage := BasicResponse{
		Message: "User created successfully",
	}
	writeJSONResponse(w, http.StatusAccepted, successMessage)
}

func writeJSONResponse(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}
