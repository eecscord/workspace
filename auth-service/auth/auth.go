package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/eecscord/workspace/auth-service/database"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux) error {
	mux.HandleFunc("/api/signin", handleSignIn)
	mux.HandleFunc("/api/signup", handleSignUp)
	return nil
}

func handleSignIn(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "POST":
			userSignIn(w, r)
			return
		default:
		  http.Error(w, errors.New("Only POST requests are allowed on this endpoint.").Error(), http.StatusBadRequest)
      return
	}
}

func handleSignUp(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "POST":
			userSignUp(w, r)
			return
		default:
		  http.Error(w, errors.New("Only POST requests are allowed on this endpoint.").Error(), http.StatusBadRequest)
      return
	}
}

func userSignIn(w http.ResponseWriter, r *http.Request) {
	credentials := Credentials{}
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if hashed password matches the one corresponding to the email.
	var hashedPassword string
	err = database.DB.QueryRow("select hashedPassword from users where email=@email", sql.Named("email", credentials.Email)).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(credentials.Password)); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	io.WriteString(w, "logged in!")
}

func userSignUp(w http.ResponseWriter, r *http.Request) {
	credentials := Credentials{}
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if email exists.
	rows := database.DB.QueryRow("SELECT email FROM users WHERE email = @email", sql.Named("email", credentials.Email))
	var email string
	if err = rows.Scan(&email); err != sql.ErrNoRows {
		http.Error(w, errors.New("Email already exists.").Error(), http.StatusBadRequest)
		return
	}

	// Hash the password using bcrypt.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, errors.New("Error hashing password.").Error(), http.StatusInternalServerError)
		log.Fatal(err.Error())
		return
	}

	// Store credentials into the database.
	_, err = database.DB.Query("INSERT INTO users(email, hashedPassword) VALUES (@email,@hashedPassword)", sql.Named("email", credentials.Email), sql.Named("hashedPassword", string(hashedPassword)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}