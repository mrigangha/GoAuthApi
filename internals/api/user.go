package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mrigangha/GoAuthApi/internals/services"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
	RoleGuest Role = "guest"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:phone`
	Password string `json:"password"`
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	user := User{
		ID:    1,
		Name:  "Mrigangha",
		Email: "test@example.com",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Encode struct into JSON
	json.NewEncoder(w).Encode(user)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	// Always close body
	defer r.Body.Close()

	var user User

	// Decode JSON body into struct
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	hashedPassword, err := services.HashPassword(user.Password)
	if err != nil {
		// Handle error appropriately
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = hashedPassword
	fmt.Println("Received:", user.Password)

	token, err := services.GenerateJWT(user.ID, user.Email)
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}
	fmt.Println("TOken: ", token)
	data, err := services.VerifyJWT(token)
	fmt.Println("Data:", data)

	w.WriteHeader(http.StatusCreated)
}
