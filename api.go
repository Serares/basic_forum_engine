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
	Error   string `json:"error"`
}

type CreateUserRequest struct {
	Alias    string `josn:"alias"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
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
	router.HandleFunc("/create_user", s.handleUserCreation).Methods("POST")

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
		writeJSONResponse(w, http.StatusMethodNotAllowed, BasicResponse{Message: "The request body can't be parsed", Error: err.Error()})
		return
	}

	if doesExist, _ := s.db.DoesUserEmailExist(req.Email); doesExist {
		writeJSONResponse(w, http.StatusConflict, BasicResponse{Message: "User with this email already exists"})
		return
	}

	userFromRequest, err := NewUserFromRequest(req)

	if err != nil {
		log.Fatal("New user from request failed to be mapped")
	}

	if err = s.db.CreateUser(userFromRequest); err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, BasicResponse{Message: "Trouble creating the user", Error: err.Error()})
		return
	}

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
