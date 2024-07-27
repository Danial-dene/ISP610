package users

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"tender-scraper/mail"
	"tender-scraper/database"
	"tender-scraper/types"
	"tender-scraper/config"
)

func RegisterUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		fmt.Println("Registering user")

		var user types.UserInfo
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if user.Email == "" {
			http.Error(w, "Email are required", http.StatusBadRequest)
			return
		}

		err = database.AddUser(db, user)
		if err != nil {
			fmt.Println("Failed to register user:", err)
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			return
		}

		cfg := config.LoadConfig()
		mail.SendEmailToUserRegister(cfg, user)
		fmt.Println("Email sent to user after register", user)
		if err != nil {
			fmt.Println("Failed to send email after user register:", err)
			http.Error(w, "Failed to send email user register", http.StatusInternalServerError)
			return
		}

		fmt.Println("User registered successfully")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}
