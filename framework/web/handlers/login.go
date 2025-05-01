package handlers

import (
	"encoding/json"
	"github.com/Andision/RoundWorld/framework/web/structures"
	"net/http"
	"strings"
)

type loginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// checkPassword is a mock function to check the password.
func checkPassword(username, password string) bool {
	if !strings.HasPrefix(username, "testAccount") {
		return false
	}
	if password != "TestAccountWillBeDeleted?JustHaveFun!" {
		return false
	}
	return true
}

// LoginHandler is the handler for the login API endpoint.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var data loginData

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if checkPassword(data.Username, data.Password) {
		token, err := structures.GenerateJwt(data.Username)
		if err != nil {
			http.Error(w, "Failed to generate JWT", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]string{
			"token": token,
		})

		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
	}
}
